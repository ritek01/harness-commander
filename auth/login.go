package auth

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	. "harness/utils"
)

func Login(ctx *cli.Context) error {
	color.Set(color.FgYellow)
	fmt.Println("Welcome to All New Harness CLI!")
	color.Unset()
	PromptAccountDetails(ctx)
	err := SaveCredentials(ctx, false)
	if err != nil {
		return err
	}
	loginError := GetAccountDetails(ctx)

	if loginError != nil {
		return loginError

	}
	err = GetUserDetails(ctx)
	if err != nil {
		return err
	}
	err = SaveCredentials(ctx, true)
	if err != nil {
		return err
	}

	return nil
}
