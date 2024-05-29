package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"harness/auth"
	"harness/defaults"
	"harness/globals"
	. "harness/share"
	"harness/utils"
	"os"
)

var Version = "Harness CLI : Hack Week 24 🎉"

type UserDataLoad struct {
	ApiKey    string `json:"apiKey"`
	AccountId string `json:"accountId"`
	BaseURL   string `json:"baseUrl"`
	UserId    string `json:"userId"`
}

var asciiArt = `
██╗  ██╗ █████╗ ██████╗ ███╗   ██╗███████╗███████╗███████╗     ██████╗ ██████╗ ███╗   ███╗███╗   ███╗ █████╗ ███╗   ██╗██████╗ ███████╗██████╗ 
██║  ██║██╔══██╗██╔══██╗████╗  ██║██╔════╝██╔════╝██╔════╝    ██╔════╝██╔═══██╗████╗ ████║████╗ ████║██╔══██╗████╗  ██║██╔══██╗██╔════╝██╔══██╗
███████║███████║██████╔╝██╔██╗ ██║█████╗  ███████╗███████╗    ██║     ██║   ██║██╔████╔██║██╔████╔██║███████║██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██╔══██║██╔══██║██╔══██╗██║╚██╗██║██╔══╝  ╚════██║╚════██║    ██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██╔══██║██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
██║  ██║██║  ██║██║  ██║██║ ╚████║███████╗███████║███████║    ╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║  ██║██║ ╚████║██████╔╝███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝     ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝
                                                                                                                                               
 Welcome to the Harness Commander, a utility to interact with the Harness Platform and seamlessly 
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
		Usage:                "Welcome to Harness Commander, a utility to interact with Harness Platform to seamlessly deploy your applications.",
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
				Usage: "Initialize Harness Commander CLI in the project",
				Action: func(context *cli.Context) error {
					if err := LoadCredentials(); err != nil {
						fmt.Println(err, "\nPlease log in first.")
						return auth.Login(context)
					}
					return cliWrapper(InitProject, context)
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy your project using Harness",
				Action: func(context *cli.Context) error {
					if err := LoadCredentials(); err != nil {
						fmt.Println(err, "\nPlease log in first.")
						return auth.Login(context)
					}
					return cliWrapper(DeployProject, context)
				},
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

func LoadCredentials() error {
	exactFilePath := utils.GetUserHomePath() + "/" + defaults.SECRETS_STORE_PATH
	credsJson, err := os.ReadFile(exactFilePath)
	if err != nil {
		return fmt.Errorf("error reading creds file: %v", err)
	}
	var UserData UserDataLoad
	err = json.Unmarshal(credsJson, &UserData)
	if err != nil {
		return fmt.Errorf("error unmarshaling creds file: %v", err)
	}
	globals.ApiKey = UserData.ApiKey
	globals.AccountId = UserData.AccountId
	globals.BaseURL = UserData.BaseURL
	globals.UserId = UserData.UserId
	return nil
}
