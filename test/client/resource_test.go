package client

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"regexp"
	"testing"
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: clientResource,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "mysecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
				resource.TestStep{
					Config: clientResourceUpdateSecret,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "newsecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: clientResourceWithScope,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref),
						testAccCheckValidSecret(ref, "mysecret"),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
						util.TestCheckResourceSet(ref, "scope", []string{"openid", "uaa.admin"}),
					),
				},
			},
		})
}

func TestAccClient_createError(t *testing.T) {
	clientid := "my-name2"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckClientDestroy(clientid),
			Steps: []resource.TestStep{
				resource.TestStep{
					Config:      clientResourceWithoutSecret,
					ExpectError: regexp.MustCompile(".*Client secret is required for client_credentials.*"),
				},
			},
		})
}

func testAccCheckValidSecret(resource, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		id := rs.Primary.ID
		auth := util.UaaSession().AuthManager()
		if _, err := auth.GetClientToken(id, secret); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckClientExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}
		util.UaaSession().Log.DebugMessage("terraform state for resource '%s': %# v", resource, rs)

		id := rs.Primary.ID
		um := util.UaaSession().ClientManager()

		// check client exists
		_, err := um.GetClient(id)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckClientDestroy(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		um := util.UaaSession().ClientManager()
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
