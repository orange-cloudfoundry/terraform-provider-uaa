package uaa

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"testing"
)

const ref = "uaa_group.new-group"
const originalDisplayName = "new.group.for.tests"
const originalDescription = "A group used for testing group resource functionality"
const updatedDisplayName = "updated.display.name"
const updatedDescription = "An updated description for the group resource"
const defaultZoneId = "uaa"
const updatedZoneId = "not-uaa"

func createTestGroupResourceAttr(attribute, value string) string {
	if attribute == "" || value == "" {
		return ""
	}
	return `	` + attribute + ` = "` + value + `"`
}

func createTestGroupResource(displayName, description, zoneId string) string {
	return `resource uaa_group "new-group" {
		` + createTestGroupResourceAttr("display_name", displayName) + `
		` + createTestGroupResourceAttr("description", description) + `
		` + createTestGroupResourceAttr("zone_id", zoneId) + `
	}`
}

func TestGroupResource_normal(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckGroupDestroy(originalDisplayName, defaultZoneId),
			Steps: []resource.TestStep{
				{
					Config: createTestGroupResource(originalDisplayName, originalDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", originalDisplayName),
						resource.TestCheckResourceAttr(ref, "description", originalDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", defaultZoneId),
					),
				},
				{
					Config: createTestGroupResource(originalDisplayName, updatedDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", originalDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", defaultZoneId),
					),
				},
				{
					Config: createTestGroupResource(updatedDisplayName, updatedDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", updatedDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", defaultZoneId),
					),
				},
				// TODO: figure out a more consistent way to run these tests to ensure updatedZoneId exists in the
				// 			test instance.  This passes if it exists, but would fail for others who run this since
				//			it's pointed at a real UAA instance and it spins up with only the default identity zone.
				//{
				//	Config: createTestGroupResource(updatedDisplayName, updatedDescription, updatedZoneId),
				//	Check: resource.ComposeTestCheckFunc(
				//		checkDataSourceGroupExists(ref),
				//		resource.TestCheckResourceAttrSet(ref, "id"),
				//		resource.TestCheckResourceAttr(ref, "display_name", updatedDisplayName),
				//		resource.TestCheckResourceAttr(ref, "description", updatedDescription),
				//		resource.TestCheckResourceAttr(ref, "zone_id", updatedZoneId),
				//	),
				//},
			},
		})
}

func testAccCheckGroupDestroy(id, zoneId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		session := testAccProvider.Meta().(*uaaapi.Session)
		gm := session.GroupManager()
		if _, err := gm.FindByDisplayName(id, zoneId); err != nil {
			switch err.(type) {
			case *errors.ModelNotFoundError:
				return nil
			default:
				return err
			}
		}
		return fmt.Errorf("group with id '%s' still exists in cloud foundry", id)
	}
}
