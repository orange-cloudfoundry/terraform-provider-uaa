package uaatest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

const clientDataResource = `
data "uaa_client" "admin-client" {
    client_id = "admin"
}
`

const clientDataResourceNotFound = `
data "uaa_client" "admin-client2" {
    client_id = "does-not-exist"
}
`

func TestAccDataSourceClient_normal(t *testing.T) {
	ref := "data.uaa_client.admin-client"
	resource.Test(t,
		resource.TestCase{
			//PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: TestManager.ProviderFactories,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: clientDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceClientExists(ref),
						resource.TestCheckResourceAttr(ref, "client_id", "admin"),
						testCheckResourceSet(ref, "authorities", []string{
							"clients.read",
							"clients.secret",
							"clients.write",
							"password.write",
							"scim.read",
							"scim.write",
							"uaa.admin",
						}),
					),
				},
			},
		})
}

func TestAccDataSourceClient_notfound(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			//PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: TestManager.ProviderFactories,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config:      clientDataResourceNotFound,
					ExpectError: regexp.MustCompile(".*Client does-not-exist not found.*"),
				},
			},
		})
}

func checkDataSourceClientExists(resource string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		TestManager.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		client_id := rs.Primary.Attributes["client_id"]

		var (
			err    error
			client uaaapi.UAAClient
		)

		client, err = TestManager.UaaSession().ClientManager().FindByClientID(client_id)
		if err != nil {
			return err
		}
		if err := AssertSame(client.ClientID, id); err != nil {
			return err
		}

		return nil
	}
}
