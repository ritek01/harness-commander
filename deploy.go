package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"harness/defaults"
	"harness/globals"
	. "harness/utils"
	"os"
	"strings"
	"time"
)

func DeployProject(c *cli.Context) error {
	color.Set(color.FgYellow)
	fmt.Println("Deploying project...")
	color.Unset()

	framework, language, err := loadProjectInfo()
	if err != nil {
		return err
	}
	ProjectprogressBar("Loading Project Info...")
	fmt.Printf("\nLoaded framework: %s\n", GetColoredText(framework, color.FgCyan))
	fmt.Printf("Loaded language: %s\n", GetColoredText(language, color.FgCyan))

	hasDockerfile, err := checkDockerfile()
	if err != nil {
		return err
	}

	if hasDockerfile {
		color.Set(color.FgGreen)
		fmt.Println("Awesome! 🐋 Dockerfile found.")
		color.Unset()
	} else {
		fmt.Print("No Dockerfile found. Would you like to create one? (y/n) : ")
		var response string
		fmt.Scanln(&response)

		if response == "y" {
			err = createDockerfile(framework, language)
			if err != nil {
				return err
			}
			ProjectprogressBar("Creating Dockerfile...")
			color.Set(color.FgGreen)
			fmt.Println("\n🐋 Dockerfile created.")
			color.Unset()
			hasDockerfile = true
		}
	}

	err = saveDockerfileInfo(hasDockerfile)
	if err != nil {
		return err
	}

	fmt.Print("\nDo you want to proceed deploying using Harness? (y/n): ")
	var proceed string
	fmt.Scanln(&proceed)

	if proceed != "y" {
		fmt.Println("Deployment aborted.")
		return nil
	}

	orgName := promptUser("\nOrg Name (default)", "default")
	projectName := promptUser("Project Name (default_project)", "default_project")

	_, err = CheckOrgExistsAndCreate(c, orgName)
	if err != nil {
		return err
	}
	_, err = CheckProjectExistsAndCreate(c, orgName, projectName)
	if err != nil {
		return err
	}

	globals.OrgId = orgName
	globals.ProjectId = projectName

	_, err = DockerConnector(c)
	if err != nil {
		return err
	}

	fmt.Print("\n\nDo you want to use the Harness Code Repository for code hosting? (y/n): ")
	var useHarnessRepo string
	fmt.Scanln(&useHarnessRepo)

	if useHarnessRepo == "y" {
		fmt.Printf("\nAlready using Harness Code Repository for code hosting ? (y/n): ")
		var useHarnessCodeRepo string
		fmt.Scanln(&useHarnessCodeRepo)
		if useHarnessCodeRepo == "y" {
			fmt.Println("\nProvide the repository Name: ")
			repoName := promptUser("Repository Name", "default_repo")
			globals.RepoName = repoName
		} else {
			err = UploadToHarnessCodeRepo()
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("Feature not supported yet.")
	}
	CreatePipeline()
	RemoveTempFiles()
	color.Set(color.FgGreen)
	fmt.Println("Deployment Done.")
	color.Unset()
	color.Set(color.FgYellow)
	fmt.Println("Thank you for using Harness Commander CLI.")
	color.Unset()
	return nil
}

func RemoveTempFiles() {
	err := os.Remove(defaults.TEMPFILEPATH)
	if err != nil {
		_ = fmt.Errorf("failed to delete temp file: %v", err)
	}
	err = os.Remove("base_pipeline.yaml")
	if err != nil {
		_ = fmt.Errorf("failed to delete base pipeline file: %v", err)
	}
	err = os.Remove("modified_pipeline.yaml")
	if err != nil {
		_ = fmt.Errorf("failed to delete modified pipeline file: %v", err)
	}
}

func promptUser(prompt, defaultValue string) string {
	fmt.Printf("%s: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func loadProjectInfo() (string, string, error) {
	file, err := os.Open(defaults.TEMPFILEPATH)
	if err != nil {
		return "", "", fmt.Errorf("please run 'harness init' first")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var framework, language string

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			switch parts[0] {
			case "framework":
				framework = parts[1]
			case "language":
				language = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("failed to read temp file: %v", err)
	}

	return framework, language, nil
}

func checkDockerfile() (bool, error) {
	if _, err := os.Stat("Dockerfile"); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check Dockerfile: %v", err)
	}
}

func createDockerfile(framework, language string) error {
	// TODO : uncomment this block
	if framework != "Spring Boot" || language != "Java" {
		return fmt.Errorf("unsupported framework or language: %s (%s)", framework, language)
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(`# Dockerfile for Spring Boot (Java)
FROM eclipse-temurin:17-jdk-focal

WORKDIR /app

COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN ./mvnw dependency:go-offline

COPY src ./src

CMD ["./mvnw", "spring-boot:run"]
`)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile: %v", err)
	}

	return nil
}

func saveDockerfileInfo(hasDockerfile bool) error {
	file, err := os.Open(defaults.TEMPFILEPATH)
	if err != nil {
		return fmt.Errorf("failed to open temp file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	dockerfileExists := false
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 && parts[0] == "dockerfile" {
			lines = append(lines, fmt.Sprintf("dockerfile=%t", hasDockerfile))
			dockerfileExists = true
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read temp file: %v", err)
	}

	if !dockerfileExists {
		lines = append(lines, fmt.Sprintf("dockerfile=%t", hasDockerfile))
	}

	file, err = os.OpenFile(defaults.TEMPFILEPATH, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open temp file for writing: %v", err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to temp file: %v", err)
		}
	}

	return nil
}

func ProjectprogressBar(info string) {
	barLength := 20
	spinChars := []string{"|", "/", "-", "\\"}
	for i := 0; i <= barLength; i++ {
		spinner := spinChars[i%len(spinChars)]
		progress := strings.Repeat("=", i) + strings.Repeat(" ", barLength-i)
		color.Set(color.FgCyan)
		fmt.Printf("\r[%s] %s %s", progress, spinner, info)
		color.Unset()
		time.Sleep(time.Millisecond * 100)
	}
}
