package uaa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

const userDataResource = `

data "uaa_user" "admin-user" {
    name = "admin"
}
`

func TestAccDataSourceUser_normal(t *testing.T) {

	ref := "data.uaa_user.admin-user"

	resource.Test(t,
		resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{

				resource.TestStep{
					Config: userDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceUserExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", "admin"),
					),
				},
			},
		})
}

func checkDataSourceUserExists(resource string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		session := testAccProvider.Meta().(*uaaapi.Session)

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("user '%s' not found in terraform state", resource)
		}

		session.Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		name := rs.Primary.Attributes["name"]

		var (
			err  error
			user uaaapi.UAAUser
		)

		user, err = session.UserManager().FindByUsername(name)
		if err != nil {
			return err
		}
		if err := assertSame(id, user.ID); err != nil {
			return err
		}

		return nil
	}
}
