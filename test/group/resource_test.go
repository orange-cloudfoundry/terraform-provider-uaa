package group

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"regexp"
	"testing"
)

const ref = "uaa_group.new-group"
const originalDisplayName = "new.group.for.tests"
const originalDescription = "A group used for testing group resource functionality"
const updatedDisplayName = "updated.display.name"
const updatedDescription = "An updated description for the group resource"
const defaultZoneId = "uaa"
const updatedZoneId = "test-zone"

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
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckGroupDestroy(originalDisplayName, updatedZoneId),
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
				{
					Config: createTestGroupResource(updatedDisplayName, updatedDescription, updatedZoneId),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", updatedDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", updatedZoneId),
					),
				},
			},
		})
}

func TestGroupResource_createError(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckGroupDestroy(ref, defaultZoneId),
			Steps: []resource.TestStep{
				{
					Config:      createTestGroupResource("", originalDescription, defaultZoneId),
					ExpectError: regexp.MustCompile("The argument \"display_name\" is required, but no definition was found."),
				},
			},
		})
}

func testAccCheckGroupDestroy(id, zoneId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		gm := util.UaaSession().GroupManager()
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
