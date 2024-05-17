package uaa

import (
	"fmt"
	"regexp"
	"testing"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

const clientResource = `
resource "uaa_client" "client1" {
    client_id = "my-name"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "mysecret"
}
`

const clientResourceUpdateSecret = `
resource "uaa_client" "client1" {
    client_id = "my-name"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "newsecret"
}
`

const clientResourceWithoutSecret = `
resource "uaa_client" "client2" {
    client_id = "my-name2"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
}
`

const clientResourceWithScope = `
resource "uaa_client" "client3" {
    client_id = "my-name-scope"
    scope = ["uaa.admin", "openid"]
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "mysecret"
}
`

func TestAccClient_normal(t *testing.T) {
	ref := "uaa_client.client1"
	clientid := "my-name"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				{
					Config: clientResource,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "mysecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						testCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						testCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
				{
					Config: clientResourceUpdateSecret,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "newsecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						testCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						testCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
			},
		})
}

func TestAccClient_scope(t *testing.T) {
	ref := "uaa_client.client3"
	clientid := "my-name-scope"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				{
					Config: clientResourceWithScope,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "mysecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						testCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						testCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
						testCheckResourceSet(ref, "scope", []string{"openid", "uaa.admin"}),
					),
				},
			},
		})
}

func TestAccClient_createError(t *testing.T) {
	clientid := "my-name2"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				{
					Config:      clientResourceWithoutSecret,
					ExpectError: regexp.MustCompile(".*Client secret is required for client_credentials.*"),
				},
			},
		})
}

func testAccCheckValidSecret(resource, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		session := testAccProvider.Meta().(*uaaapi.Session)
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		id := rs.Primary.ID
		auth := session.AuthManager()
		if _, err := auth.GetClientToken(id, secret); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckClientExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		session := testAccProvider.Meta().(*uaaapi.Session)
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}
		session.Log.DebugMessage("terraform state for resource '%s': %# v", resource, rs)

		id := rs.Primary.ID
		um := session.ClientManager()

		// check client exists
		_, err := um.GetClient(id)
		if err != nil {
			return err
		}
		return nil
	}
}

func testCheckResourceSet(ref string, attr string, values []string) resource.TestCheckFunc {

	lTests := make([]resource.TestCheckFunc, 0)

	for i, cVal := range values {
		lKey := fmt.Sprintf("%s.%d", attr, i)
		lTests = append(lTests, resource.TestCheckResourceAttr(ref, lKey, cVal))
	}

	return resource.ComposeTestCheckFunc(lTests...)
}

func testAccCheckClientDestroy(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		session := testAccProvider.Meta().(*uaaapi.Session)
		um := session.ClientManager()
		if _, err := um.FindByClientID(id); err != nil {
			switch err.(type) {
			case *errors.ModelNotFoundError:
				return nil
			default:
				return err
			}
		}
		return fmt.Errorf("client with id '%s' still exists in cloud foundry", id)
	}
}
