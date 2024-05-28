package utils

import (
	"fmt"
	"harness/globals"
	"harness/netclient"
	. "harness/types"
	"os"
	"os/exec"
	"path/filepath"
)

func UploadToHarnessCodeRepo() error {
	if !isGitInitialized() {
		fmt.Println("Git is not initialized. Initializing...")
		if err := initializeGitRepo(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %v", err)
		}
	} else {
		fmt.Println("Git is already initialized.")
	}

	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	folderName, err := getCurrentDirectoryName()
	spacePath := globals.AccountId + "/" + globals.OrgId + "/" + globals.ProjectId
	// Define repo details
	repoDetails := CreateRepoRequest{
		DefaultBranch: branch,
		Description:   "My new Harness repo",
		GitIgnore:     "",
		IsPublic:      false,
		License:       "",
		Uid:           folderName,
		Readme:        false,
		ParentRef:     spacePath,
	}

	respBody, err := createHarnessRepo(spacePath, repoDetails)
	if err != nil {
		fmt.Println("Error creating Harness repository:", err)
		return nil
	}
	fmt.Println("Repository created successfully:", respBody)

	repoURL := "https://git.harness.io/" + spacePath + "/" + folderName + ".git"
	if err := addRemoteAndPush(repoURL, branch); err != nil {
		return fmt.Errorf("failed to push to Harness repository: %v", err)
	}

	fmt.Println("Project successfully pushed to Harness Code Repository.")
	return nil
}

func isGitInitialized() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func initializeGitRepo() error {
	cmds := [][]string{
		{"git", "init"},
		{"git", "add", "."},
		{"git", "commit", "-m", "Project Init"},
	}

	for _, cmd := range cmds {
		if err := runCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return string(output), nil
}

func addRemoteAndPush(repoURL, branch string) error {
	cmds := [][]string{
		{"git", "remote", "add", "code", repoURL},
		{"git", "push", "code", branch},
	}

	for _, cmd := range cmds {
		if err := runCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if name == "git" && (args[0] == "push" || args[0] == "pull") {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_USERNAME=%s", globals.UserId))
		cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_PASSWORD=%s", globals.ApiKey))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %s %v failed: %v", name, args, err)
	}
	return nil
}

func createHarnessRepo(spacePath string, repoDetails CreateRepoRequest) (ResponseBody, error) {
	reqUrl := fmt.Sprintf("https://app.harness.io/gateway/code/api/v1/repos?routingId=%s&space_path=%s", globals.AccountId, spacePath)
	contentType := "application/json"

	respBody, err := netclient.Post(reqUrl, globals.ApiKey, repoDetails, contentType, nil)
	if err != nil {
		return ResponseBody{}, fmt.Errorf("failed to create Harness repository: %v", err)
	}

	return respBody, nil
}

func getCurrentDirectoryName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(dir), nil
}
