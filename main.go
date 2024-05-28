package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Harness-CLI",
		Usage: "CLI tool to interact with Harness",
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "Authenticate with Harness",
				Action: login,
			},
			{
				Name:   "init",
				Usage:  "Initialize Harness CLI in the project",
				Action: initProject,
			},
			{
				Name:   "deploy",
				Usage:  "Deploy the project using Harness",
				Action: deploy,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func login(c *cli.Context) error {
	fmt.Println("Logging in to Harness...")
	// Implement login functionality here
	return nil
}

func initProject(c *cli.Context) error {
	fmt.Println("Initializing project...")
	// Implement project initialization functionality here
	return nil
}

func deploy(c *cli.Context) error {
	fmt.Println("Deploying project...")
	// Implement deploy functionality here
	return nil
}
