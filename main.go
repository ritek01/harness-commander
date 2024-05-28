package main

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"harness/auth"
	. "harness/share"
	"os"
)

var Version = "Harness CLI : Hack Week 24"

var asciiArt = `

 ██╗  ██╗  █████╗  ██████╗  ███╗   ██╗ ███████╗  ██████╗  ██████╗        █████╗  ██╗      ██╗
 ██║  ██║ ██╔══██╗ ██╔══██╗ ████╗  ██║ ██╔════╝ ██╔════╝ ██╔════╝       ██╔══██╗ ██║      ██║
 ███████║ ███████║ ██████╔╝ ██╔██╗ ██║ █████╗   ╚█████╗  ╚█████╗        ██║  ╚═╝ ██║   	  ██║
 ██╔══██║ ██╔══██║ ██╔══██╗ ██║ ╚████║ ██╔══╝    ╚═══██╗  ╚═══██╗       ██║  ██╗ ██║   	  ██║
 ██║  ██║ ██║  ██║ ██║  ██║ ██║  ╚███║ ███████╗ ██████╔╝ ██████╔╝       ╚█████╔╝ ███████║ ██║
 ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═╝   ╚══╝ ╚══════╝ ╚═════╝  ╚═════╝         ╚════╝  ╚══════╝ ╚═╝

 Welcome to the new Harness CLI utility to interact with the Harness Platform and seamlessly 
 deploy your applications.
`

type cliFnWrapper func(ctx *cli.Context) error

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Println(cCtx.App.Version)
	}
}

func main() {
	color.Set(color.FgBlue)
	fmt.Println(asciiArt)
	color.Unset()

	app := &cli.App{
		Name:                 "harness",
		Version:              Version,
		Usage:                "Welcome to New Harness CLI utility to interact with Harness Platform to seamlessly deploy your applications.",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "Authenticate with Harness",
				Action: func(context *cli.Context) error {
					return cliWrapper(func(context *cli.Context) error {
						return auth.Login(context)
					}, context)
				},
			},
			{
				Name:  "init",
				Usage: "Initialize Harness CLI in the project",
				Action: func(context *cli.Context) error {
					return cliWrapper(func(context *cli.Context) error {
						return InitProject(context)
					}, context)
				},
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

func cliWrapper(fn cliFnWrapper, ctx *cli.Context) error {
	if CliRequestData.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if CliRequestData.Json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return fn(ctx)
}

func beforeAction(globalFlags []cli.Flag) {
	altsrc.InitInputSourceWithContext(globalFlags, altsrc.NewYamlSourceFromFlagFunc("load"))
}
