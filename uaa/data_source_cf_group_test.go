package uaa

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"regexp"
	"testing"
)

const groupDataResource = `
data uaa_group "uaa-admin" {
	display_name = "uaa.admin"
}
`

const groupDataResourceNotFound = `
data uaa_group "not-found" {
	display_name = "not.found"
}
`

func TestGroupDataSourceClient_normal(t *testing.T) {
	ref := "data.uaa_group.uaa-admin"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			Steps: []resource.TestStep{
				{
					Config: groupDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", "uaa.admin"),
						resource.TestCheckResourceAttr(ref, "description", "Act as an administrator throughout the UAA"),
						resource.TestCheckResourceAttr(ref, "zone_id", "uaa"),
					),
				},
			},
		})
}

func TestGroupDataSourceClient_notFound(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			Steps: []resource.TestStep{
				{
					Config:      groupDataResourceNotFound,
					ExpectError: regexp.MustCompile(".*Group not.found not found.*"),
				},
			},
		})
}

func checkDataSourceGroupExists(resource string) resource.TestCheckFunc {

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
		display_name := rs.Primary.Attributes["display_name"]

		var (
			err   error
			group uaaapi.UAAGroup
		)

		group, err = session.GroupManager().FindByDisplayName(display_name)
		if err != nil {
			return err
		}
		if err := assertSame(id, group.ID); err != nil {
			return err
		}

		return nil
	}
}
