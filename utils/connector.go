package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"harness/globals"
	"os"
	"strings"
	"time"
)

func DockerConnector(c *cli.Context) (bool, error) {
	dockerConnName := promptUser("Enter Docker Connector Name", "DockerConnector")
	progressBar(dockerConnName)
	color.Set(color.FgGreen)
	fmt.Printf("\nDocker Connector : %s is present in the account\n", dockerConnName)
	color.Unset()
	globals.DockerConnector = dockerConnName
	return true, nil
}

func K8sConnector(c *cli.Context) (bool, error) {
	K8sConnName := promptUser("Enter K8s Connector Name", "K8sConnector")
	progressBar(K8sConnName)
	color.Set(color.FgGreen)
	fmt.Printf("\nK8s Connector %s is present in the account\n", K8sConnName)
	color.Unset()
	return true, nil
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
func progressBar(connectorName string) {
	barLength := 20
	spinChars := []string{"|", "/", "-", "\\"}
	for i := 0; i <= barLength; i++ {
		spinner := spinChars[i%len(spinChars)]
		progress := strings.Repeat("=", i) + strings.Repeat(" ", barLength-i)
		color.Set(color.FgCyan)
		fmt.Printf("\r[%s] %s %s", progress, spinner, "Checking... `"+connectorName+"` is present in the account...")
		color.Unset()
		time.Sleep(time.Millisecond * 65)
	}
}
