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
)

func DeployProject(c *cli.Context) error {
	color.Set(color.FgYellow)
	fmt.Println("Deploying project...")
	color.Unset()

	framework, language, err := loadProjectInfo()
	if err != nil {
		return err
	}

	fmt.Printf("Loaded framework: %s\n", GetColoredText(framework, color.FgCyan))
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
			color.Set(color.FgGreen)
			fmt.Println("🐋 Dockerfile created.")
			color.Unset()
			hasDockerfile = true
		}
	}

	err = saveDockerfileInfo(hasDockerfile)
	if err != nil {
		return err
	}

	fmt.Print("Do you want to proceed deploying using Harness? (y/n): ")
	var proceed string
	fmt.Scanln(&proceed)

	if proceed != "y" {
		fmt.Println("Deployment aborted.")
		return nil
	}

	orgName := promptUser("Org Name (default)", "default")
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

	fmt.Print("Do you want to use the Harness Code Repository for code hosting? (y/n): ")
	var useHarnessRepo string
	fmt.Scanln(&useHarnessRepo)

	if useHarnessRepo == "y" {
		err = UploadToHarnessCodeRepo()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Feature not supported yet.")
	}
	fmt.Println("Deployment process initialized.")

	return nil
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
		return "", "", fmt.Errorf("failed to open temp file: %v", err)
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
