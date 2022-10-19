package client

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"regexp"
	"testing"
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: clientDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceClientExists(ref),
						resource.TestCheckResourceAttr(ref, "client_id", "admin"),
						util.TestCheckResourceSet(ref, "authorities", []string{
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
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

		util.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		client_id := rs.Primary.Attributes["client_id"]

		var (
			err    error
			client api.UAAClient
		)

		client, err = util.UaaSession().ClientManager().FindByClientID(client_id)
		if err != nil {
			return err
		}
		if err := util.AssertSame(client.ClientID, id); err != nil {
			return err
		}

		return nil
	}
}
