package uaa

import (
	"fmt"
	"testing"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

const ldapUserResource = `

resource "uaa_user" "manager1" {
    name = "manager1@acme.com"
    origin = "ldap"
}
`

const userResourceWithGroups = `

resource "uaa_user" "admin-service-user" {
    name = "cf-admin"
	password = "qwerty"
	given_name = "Build"
	family_name = "User"
    groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]
}
`

const userResourceWithGroupsUpdate = `

resource "uaa_user" "admin-service-user" {
    name = "cf-admin"
	password = "asdfg"
	email = "cf-admin@acme.com"
    groups = [ "cloud_controller.admin", "clients.admin", "uaa.admin", "doppler.firehose" ]
}
`

func TestAccUser_LdapOrigin_normal(t *testing.T) {

	ref := "uaa_user.manager1"
	username := "manager1@acme.com"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckUserDestroy(username),
			Steps: []resource.TestStep{

				resource.TestStep{
					Config: ldapUserResource,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", username),
						resource.TestCheckResourceAttr(
							ref, "origin", "ldap"),
						resource.TestCheckResourceAttr(
							ref, "email", username),
					),
				},
			},
		})
}

// TODO: this test isn't working with test containers...
func TestAccUser_WithGroups_normal(t *testing.T) {

	ref := "uaa_user.admin-service-user"
	username := "cf-admin"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      testAccCheckUserDestroy(username),
			Steps: []resource.TestStep{

				resource.TestStep{
					Config: userResourceWithGroups,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", username),
						resource.TestCheckResourceAttr(
							ref, "password", "qwerty"),
						resource.TestCheckResourceAttr(
							ref, "email", username),
						testCheckResourceSet(ref, "groups", []string{
							"cloud_controller.admin",
							"scim.read",
							"scim.write",
						}),
					),
				},

				resource.TestStep{
					Config: userResourceWithGroupsUpdate,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", "cf-admin"),
						resource.TestCheckResourceAttr(
							ref, "password", "asdfg"),
						resource.TestCheckResourceAttr(
							ref, "email", "cf-admin@acme.com"),
						testCheckResourceSet(ref, "groups", []string{
							"clients.admin",
							"cloud_controller.admin",
							"doppler.firehose",
							"uaa.admin",
						}),
					),
				},
			},
		})
}

func testAccCheckUserExists(resource string) resource.TestCheckFunc {

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
		attributes := rs.Primary.Attributes

		um := session.UserManager()
		user, err := um.GetUser(id)
		if err != nil {
			return err
		}

		session.Log.DebugMessage(
			"retrieved user for resource '%s' with id '%s': %# v",
			resource, id, user)

		if err := assertEquals(attributes, "name", user.Username); err != nil {
			return err
		}
		if err := assertEquals(attributes, "origin", user.Origin); err != nil {
			return err
		}
		if err := assertEquals(attributes, "given_name", user.Name.GivenName); err != nil {
			return err
		}
		if err := assertEquals(attributes, "family_name", user.Name.FamilyName); err != nil {
			return err
		}
		if err := assertEquals(attributes, "email", user.Emails[0].Value); err != nil {
			return err
		}

		var groups []interface{}
		for _, g := range user.Groups {
			if !um.IsDefaultGroup(g.Display) {
				groups = append(groups, g.Display)
			}
		}
		if err := assertSetEquals(attributes, "groups", groups); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckUserDestroy(username string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		session := testAccProvider.Meta().(*uaaapi.Session)
		um := session.UserManager()
		if _, err := um.FindByUsername(username); err != nil {
			switch err.(type) {
			case *errors.ModelNotFoundError:
				return nil
			default:
				return err
			}
		}
		return fmt.Errorf("user with username '%s' still exists in cloud foundry", username)
	}
}
