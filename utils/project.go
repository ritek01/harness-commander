package utils

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"harness/globals"
	"harness/netclient"
	. "harness/types"
)

func getProjectIdByName(project []ProjectDetails, projectName string) string {
	for _, p := range project {
		if p.Identifier == projectName {
			return p.Identifier
		}
	}
	return ""
}

func getProjects(c *cli.Context, orgName string) []ProjectDetails {
	url := fmt.Sprintf("%s/projects?accountIdentifier=%s&orgIdentifier=%s&pageSize=1000", GetNGBaseURL(c), globals.AccountId, orgName)
	resp, err := netclient.GetNew(url, globals.ApiKey)
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch projects", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projects ProjectListBody
	err = json.Unmarshal(byteData, &projects)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projectDetails []ProjectDetails

	for _, p := range projects.Projects {
		projectDetails = append(projectDetails, p.Project)
	}
	return projectDetails
}

func CheckProjectExistsAndCreate(c *cli.Context, orgName string, projectName string) (bool, error) {
	projects := getProjects(c, orgName)
	projectId := getProjectIdByName(projects, projectName)

	if len(projectId) == 0 {
		fmt.Println("Project not found. Creating project...")
		url := fmt.Sprintf("%s/projects?accountIdentifier=%s&orgIdentifier=%s", GetNGBaseURL(c), globals.AccountId, orgName)
		resp, err := netclient.PostNew(url, globals.ApiKey, ProjectBody{
			Project: ProjectDetails{
				OrgIdentifier: orgName,
				Identifier:    projectName,
				Name:          projectName,
				Color:         "#5dc22f",
				Modules:       []string{"CD"},
				Description:   "",
			}})
		if err != nil || resp.Status != "SUCCESS" {
			log.Fatal("Failed to create project", err)
		}
		log.Info("Project created successfully")
	} else {
		fmt.Printf("Project '%s' already exists. Do you want to use it? (y/n): ", projectName)
		var useProject string
		fmt.Scanln(&useProject)

		if useProject != "y" {
			fmt.Println("Deployment aborted.")
			return false, nil
		}
	}
	return true, nil
}
