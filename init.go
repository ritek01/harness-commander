package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var tempFilePath = "project_info.tmp"

func InitProject(c *cli.Context) error {

	color.Set(color.FgYellow)
	fmt.Println("Initializing project...")
	color.Unset()
	framework, language := detectProjectFramework()
	fmt.Printf("Detected framework: %s\n", framework)
	fmt.Printf("Detected language: %s\n", language)
	fmt.Print("Is this correct? (y/n) : ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return err
	}

	if response != "y" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the framework: ")
		framework, _ = reader.ReadString('\n')
		fmt.Print("Enter the language: ")
		language, _ = reader.ReadString('\n')
		framework = strings.TrimSpace(framework)
		language = strings.TrimSpace(language)
	}

	err = saveProjectInfo(framework, language)
	if err != nil {
		return err
	}

	fmt.Println("Project initialized with the following framework and language:")
	fmt.Printf("Framework: %s\n", framework)
	fmt.Printf("Language: %s\n", language)

	progressBar()
	color.Set(color.FgGreen)
	fmt.Println("\nHarness project initialized successfully!")
	color.Unset()

	return nil
}

func progressBar() {
	barLength := 20
	spinChars := []string{"|", "/", "-", "\\"}
	for i := 0; i <= barLength; i++ {
		spinner := spinChars[i%len(spinChars)]
		progress := strings.Repeat("=", i) + strings.Repeat(" ", barLength-i)
		color.Set(color.FgCyan)
		fmt.Printf("\r[%s] %s %s", progress, spinner, "Initializing...")
		color.Unset()
		time.Sleep(time.Millisecond * 100)
	}
}

func detectProjectFramework() (string, string) {
	type Framework struct {
		Name     string
		Language string
	}

	frameworks := map[string]Framework{
		"pom.xml":          {"Spring Boot", "Java"},
		"package.json":     {"Node.js", "JavaScript"},
		"build.gradle":     {"Spring Boot", "Java"},
		"requirements.txt": {"Flask", "Python"},
		"Gemfile":          {"Ruby on Rails", "Ruby"},
		"go.mod":           {"Go Module", "Go"},
		"Cargo.toml":       {"Cargo", "Rust"},
		"composer.json":    {"Composer", "PHP"},
		"CMakeLists.txt":   {"CMake", "C++"},
		"Makefile":         {"Make", "C/C++"},
		"mix.exs":          {"Elixir Mix", "Elixir"},
		"project.clj":      {"Leiningen", "Clojure"},
		"build.sbt":        {"SBT", "Scala"},
		"pubspec.yaml":     {"Dart", "Dart"},
	}

	for file, framework := range frameworks {
		if _, err := os.Stat(file); err == nil {
			return framework.Name, framework.Language
		}
	}

	return "Unknown", "Unknown"
}

func saveProjectInfo(framework, language string) error {
	file, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("framework=%s\nlanguage=%s\n", framework, language))
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %v", err)
	}

	return nil
}
