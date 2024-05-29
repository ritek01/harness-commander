package utils

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"harness/globals"
	"harness/netclient"
	. "harness/types"
	"log"
)

func getOrganisations(c *cli.Context) []OrgDetails {
	url := fmt.Sprintf("%s/aggregate/organizations?accountIdentifier=%s&pageSize=1000", GetNGBaseURL(c), globals.AccountId)
	resp, err := netclient.GetNew(url, globals.ApiKey)
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch organisations", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch organisations", err)
	}
	var orgListBody OrgListBody
	err = json.Unmarshal(byteData, &orgListBody)
	if err != nil {
		log.Fatal("Failed to fetch organisations", err)
	}
	var details []OrgDetails

	for _, o := range orgListBody.Organisations {
		details = append(details, o.Org.Org)
	}
	return details
}
func findOrgIdByName(org []OrgDetails, orgName string) string {
	for _, p := range org {
		if p.Name == orgName {
			return p.Identifier
		}
	}
	return ""
}
func createHarnessOrg(c *cli.Context, orgName string) error {
	url := fmt.Sprintf("%s/organizations?accountIdentifier=%s", GetNGBaseURL(c), globals.AccountId)
	//payload := fmt.Sprintf(`{"Identifier": "%s", "Name": "%s"}`, orgName, orgName)
	resp, err := netclient.PostNew(url, globals.ApiKey, OrgBody{
		Org: OrgDetails{
			Identifier:  orgName,
			Name:        orgName,
			Description: "",
		}})
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to create organization", err)
	}
	return nil
}

func CheckOrgExistsAndCreate(c *cli.Context, orgName string) (bool, error) {
	org := getOrganisations(c)
	orgId := findOrgIdByName(org, orgName)
	if len(orgId) == 0 {
		fmt.Printf("\nOrganization '%s' does not exist. Do you want to create it? (y/n): ", orgName)
		var createOrg string
		fmt.Scanln(&createOrg)

		if createOrg != "y" {
			fmt.Println("Deployment aborted.")
			return false, nil
		}
		err := createHarnessOrg(c, orgName)
		if err != nil {
			return false, err
		}
	} else {
		fmt.Printf("\nOrganization '%s' already exists. Do you want to use it? (y/n): ", orgName)
		var useOrg string
		fmt.Scanln(&useOrg)

		if useOrg != "y" {
			fmt.Println("Deployment aborted.")
			return false, nil
		}
	}
	return true, nil
}
