package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"harness/globals"
	"net/http"
	"os"
	_ "os"
)

type Pipeline struct {
	Pipeline struct {
		Name              string   `yaml:"name"`
		Identifier        string   `yaml:"identifier"`
		ProjectIdentifier string   `yaml:"projectIdentifier"`
		OrgIdentifier     string   `yaml:"orgIdentifier"`
		Tags              struct{} `yaml:"tags"`
		Properties        struct {
			Ci struct {
				Codebase struct {
					RepoName string `yaml:"repoName"`
					Build    string `yaml:"build"`
				} `yaml:"codebase"`
			} `yaml:"ci"`
		} `yaml:"properties"`
		Stages []struct {
			Stage struct {
				Name        string `yaml:"name"`
				Identifier  string `yaml:"identifier"`
				Description string `yaml:"description"`
				Type        string `yaml:"type"`
				Spec        struct {
					CloneCodebase bool `yaml:"cloneCodebase"`
					Platform      struct {
						Os   string `yaml:"os"`
						Arch string `yaml:"arch"`
					} `yaml:"platform"`
					Runtime struct {
						Type string   `yaml:"type"`
						Spec struct{} `yaml:"spec"`
					} `yaml:"runtime"`
					Execution struct {
						Steps []struct {
							Step struct {
								Type       string `yaml:"type"`
								Name       string `yaml:"name"`
								Identifier string `yaml:"identifier"`
								Spec       struct {
									ConnectorRef    string   `yaml:"connectorRef"`
									Repo            string   `yaml:"repo"`
									Tags            []string `yaml:"tags"`
									Image           string   `yaml:"image,omitempty"`
									Shell           string   `yaml:"shell,omitempty"`
									Command         string   `yaml:"command,omitempty"`
									ImagePullPolicy string   `yaml:"imagePullPolicy,omitempty"`
								} `yaml:"spec"`
							} `yaml:"step"`
						} `yaml:"steps"`
					} `yaml:"execution"`
				} `yaml:"spec"`
			} `yaml:"stage"`
		} `yaml:"stages"`
	} `yaml:"pipeline"`
}

func CreatePipeline() {
	// Load the base YAML
	baseYAML, err := os.ReadFile("base_pipeline.yaml")
	if err != nil {
		fmt.Printf("Error reading base pipeline YAML file: %v\n", err)
		return
	}

	var pipeline Pipeline
	err = yaml.Unmarshal(baseYAML, &pipeline)
	if err != nil {
		fmt.Printf("Error unmarshalling base pipeline YAML: %v\n", err)
		return
	}

	pipeline.Pipeline.ProjectIdentifier = globals.ProjectId
	pipeline.Pipeline.OrgIdentifier = globals.OrgId
	pipeline.Pipeline.Properties.Ci.Codebase.RepoName = globals.RepoName
	pipeline.Pipeline.Stages[0].Stage.Spec.Execution.Steps[0].Step.Spec.ConnectorRef = globals.DockerConnector

	modifiedYAML, err := yaml.Marshal(&pipeline)
	if err != nil {
		fmt.Printf("Error marshalling modified pipeline YAML: %v\n", err)
		return
	}

	err = os.WriteFile("modified_pipeline.yaml", modifiedYAML, 0644)
	if err != nil {
		fmt.Printf("Error writing modified pipeline YAML file: %v\n", err)
		return
	}

	err = CreateHarnessPipeline(modifiedYAML)
	if err != nil {
		fmt.Printf("Error creating pipeline in Harness: %v\n", err)
		return
	}
	color.Set(color.FgGreen)
	fmt.Println("Pipeline created successfully!")
	color.Unset()
	PipelineUrl := "https://app.harness.io/ng/account/" + globals.AccountId + "/home/orgs/" + globals.OrgId + "/projects/" + globals.ProjectId + "/pipelines/" + pipeline.Pipeline.Identifier + "/pipeline-studio/"
	fmt.Printf("Pipeline Url : %s\n", GetColoredText(PipelineUrl, color.FgCyan))
}

func CreateHarnessPipeline(pipelineYAML []byte) error {
	// Create the base URL
	baseURL := "https://app.harness.io/pipeline/api/pipelines/v2"

	// Create the URL with the query parameters
	reqURL := fmt.Sprintf("%s?accountIdentifier=%s&orgIdentifier=%s&projectIdentifier=%s",
		baseURL, globals.AccountId, globals.OrgId, globals.ProjectId)

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(pipelineYAML))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/yaml")
	req.Header.Set("x-api-key", globals.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			return fmt.Errorf("failed to parse error response body: %v", err)
		}
		return fmt.Errorf("failed to create pipeline, status code: %d, error message: %v", resp.StatusCode, responseBody)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}
	return nil
}
