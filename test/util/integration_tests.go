package util

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider"
	"os"
	"testing"
)

var uaaProvider = provider.Provider()

var ProviderFactories = map[string]func() (*schema.Provider, error){

	"uaa": func() (*schema.Provider, error) {
		return uaaProvider, nil
	},
}

func UaaSession() *api.Session {
	return uaaProvider.Meta().(*api.Session)
}

func IntegrationTestPreCheck(t *testing.T) {

	if !testAccEnvironmentSet() {
		t.Fatal("Acceptance environment has not been set.")
	}
}

func testAccEnvironmentSet() bool {

	loginEndpoint := os.Getenv("UAA_LOGIN_URL")
	authEndpoint := os.Getenv("UAA_AUTH_URL")
	clientID := os.Getenv("UAA_CLIENT_ID")
	clientSecret := os.Getenv("UAA_CLIENT_SECRET")

	if len(loginEndpoint) == 0 || len(authEndpoint) == 0 || len(clientID) == 0 || len(clientSecret) == 0 {
		envVars := "UAA_LOGIN_URL, UAA_AUTH_URL, UAA_CLIENT_ID, UAA_CLIENT_SECRET"
		fmt.Println(envVars + " must be set when running tests.")
		return false
	}
	return true
}
