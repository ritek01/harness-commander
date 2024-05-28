package utils

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"harness/defaults"
	"harness/netclient"
	. "harness/share"
	. "harness/types"
	"io"
	"os"
	"strings"
)

func WriteToFile(text string, filename string, overwrite bool) {
	exactFilePath := GetUserHomePath() + "/" + filename
	var permissions = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if overwrite {
		permissions = os.O_APPEND | os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	f, err := os.OpenFile(exactFilePath, permissions, 0644)
	if overwrite {
		f.WriteString("")
	}
	f.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}

func ReadFromFile(filepath string) (s string, r []byte) {
	var _fileContents = ""

	file, _ := os.OpenFile(filepath, os.O_RDONLY, 0644)
	defer file.Close()

	byteValue, readError := io.ReadAll(file)
	if readError != nil {
		fmt.Println("Error reading file:", readError)
		return "", nil
	}
	_fileContents = string(byteValue)

	return _fileContents, byteValue
}

func SaveCredentials(c *cli.Context, showWelcome bool) (err error) {
	baseURL := c.String("base-url")
	if baseURL == "" {
		baseURL = CliRequestData.BaseUrl
	}
	if CliRequestData.BaseUrl == "" {
		baseURL = defaults.HARNESS_PROD_URL
	}
	authCredentials := SecretStore{
		ApiKey:    CliRequestData.AuthToken,
		AccountId: CliRequestData.Account,
		BaseURL:   baseURL,
		UserId:    CliRequestData.UserId,
	}
	jsonObj, err := json.MarshalIndent(authCredentials, "", "  ")
	if err != nil {
		fmt.Println("Error creating secrets json:", err)
		return
	}

	WriteToFile(string(jsonObj), defaults.SECRETS_STORE_PATH, true)
	if showWelcome {
		println(GetColoredText("Login successfully done. Yay!", color.FgGreen))
		println(GetColoredText("Get ready to harness the power of Harness!", color.FgHiMagenta))
	}
	return nil
}

func GetUserHomePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get user's home directory:", err)
		return ""
	}
	return homeDir
}

func GetColoredText(text string, textColor color.Attribute) string {
	colored := color.New(textColor).SprintFunc()
	return colored(text)
}

func PromptAccountDetails(ctx *cli.Context) bool {
	promptConfirm := false

	if len(CliRequestData.Account) == 0 {
		promptConfirm = true
		CliRequestData.Account = TextInput("Harness Account Id that you wish to login to : ")
	}

	if len(CliRequestData.AuthToken) == 0 {
		promptConfirm = true
		CliRequestData.AuthToken = TextInput("Provide your api-key : ")
	}
	return promptConfirm
}

func GetAccountDetails(ctx *cli.Context) error {
	var baseURL = GetNGBaseURL(ctx)
	accountsEndpoint := defaults.ACCOUNTS_ENDPOINT + CliRequestData.Account
	url := GetUrlWithQueryParams("", baseURL, accountsEndpoint, map[string]string{
		"accountIdentifier": CliRequestData.Account,
	})
	resp, err := netclient.Get(url, CliRequestData.AuthToken)
	if err != nil {
		println(GetColoredText("Could not log in: Did you provide correct credentials?", color.FgRed))
		fmt.Printf("Response code: %s \n", resp.Code)
		return err
	}
	return nil
}

func GetNGBaseURL(c *cli.Context) string {
	baseURL := c.String("base-url")
	if baseURL == "" {
		if CliRequestData.BaseUrl == "" {
			baseURL = defaults.HARNESS_PROD_URL
		} else {
			baseURL = CliRequestData.BaseUrl
		}
	}

	baseURL = strings.TrimRight(baseURL, "/")
	baseURL = baseURL + defaults.NG_BASE_URL
	return baseURL
}

func TextInput(question string) string {
	var text = ""
	prompt := &survey.Input{
		Message: question,
	}
	err := survey.AskOne(prompt, &text, survey.WithValidator(survey.Required))
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
	return text
}
func GetUrlWithQueryParams(environment string, service string, endpoint string, queryParams map[string]string) string {
	if len(queryParams) > 0 {
		params := ""
		totalItems := len(queryParams)
		currentIndex := 0
		for k, v := range queryParams {
			currentIndex++
			if v != "" {
				if currentIndex == totalItems {
					params = params + k + "=" + v
				} else {
					params = params + k + "=" + v + "&"
				}
			}

		}
		lastChar := params[len(params)-1]
		if lastChar == '&' {
			params = strings.TrimSuffix(params, string('&'))
		}
		return fmt.Sprintf("%s/%s?%s", service, endpoint, params)
	} else {
		return fmt.Sprintf("%s/%s", service, endpoint)
	}
}

func GetUserDetails(ctx *cli.Context) error {
	var baseURL = GetNGBaseURL(ctx)
	url := GetUrlWithQueryParams("", baseURL, defaults.USER_INFO_ENDPOINT, map[string]string{
		"accountIdentifier": CliRequestData.Account,
	})
	resp, err := netclient.Get(url, CliRequestData.AuthToken)
	if err != nil {
		println(GetColoredText("Could not log in: Did you provide correct credentials?", color.FgRed))
		fmt.Printf("Response code: %s \n", resp.Code)
		return err
	}
	dataJSON, err := json.Marshal(resp.Data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return err
	}
	var currentUserInfo UserInfo
	err = json.Unmarshal(dataJSON, &currentUserInfo)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}
	CliRequestData.UserId = currentUserInfo.Email
	return nil
}
