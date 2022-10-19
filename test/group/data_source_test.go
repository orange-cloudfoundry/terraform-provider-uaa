package group

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
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

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		util.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		displayName := rs.Primary.Attributes["display_name"]
		zoneId := rs.Primary.Attributes["zone_id"]

		var (
			err   error
			group *api.UAAGroup
		)

		group, err = util.UaaSession().GroupManager().FindByDisplayName(displayName, zoneId)
		if err != nil {
			return err
		}
		if err := util.AssertSame(group.ID, id); err != nil {
			return err
		}

		return nil
	}
}
