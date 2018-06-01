package uaa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"regexp"
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
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: clientDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceClientExists(ref),
						resource.TestCheckResourceAttr(ref, "client_id", "admin"),
						testCheckResourceSet(ref, "authorities", []string{
							"clients.read",
							"password.write",
							"clients.secret",
							"clients.write",
							"uaa.admin",
							"scim.write",
							"scim.read",
						}),
					),
				},
			},
		})
}

func TestAccDataSourceClient_notfound(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
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

		session := testAccProvider.Meta().(*uaaapi.Session)

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		session.Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		client_id := rs.Primary.Attributes["client_id"]

		var (
			err    error
			client uaaapi.UAAClient
		)

		client, err = session.ClientManager().FindByClientID(client_id)
		if err != nil {
			return err
		}
		if err := assertSame(id, client.ClientID); err != nil {
			return err
		}

		return nil
	}
}
