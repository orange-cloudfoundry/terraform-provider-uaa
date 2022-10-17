package user_test

//import (
//	"fmt"
//	"github.com/terraform-providers/terraform-provider-uaa/uaa"
//	"github.com/terraform-providers/terraform-provider-uaa/util"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//
//	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
//)
//
//const userDataResource = `
//
//data "uaa_user" "admin-user" {
//    name = "admin"
//}
//`
//
//func TestAccDataSourceUser_normal(t *testing.T) {
//
//	ref := "data.uaa_user.admin-user"
//
//	resource.Test(t,
//		resource.TestCase{
//			//PreCheck:          func() { uaa.testAccPreCheck(t) },
//			ProviderFactories: uaa.testAccProvidersFactories,
//			Steps: []resource.TestStep{
//				{
//					Config: userDataResource,
//					Check: resource.ComposeTestCheckFunc(
//						checkDataSourceUserExists(ref),
//						resource.TestCheckResourceAttr(
//							ref, "name", "admin"),
//					),
//				},
//			},
//		})
//}
//
//func checkDataSourceUserExists(resource string) resource.TestCheckFunc {
//
//	return func(s *terraform.State) error {
//
//		session := uaa.testAccProvider.Meta().(*uaaapi.Session)
//
//		rs, ok := s.RootModule().Resources[resource]
//		if !ok {
//			return fmt.Errorf("user '%s' not found in terraform state", resource)
//		}
//
//		session.Log.DebugMessage(
//			"terraform state for resource '%s': %# v",
//			resource, rs)
//
//		id := rs.Primary.ID
//		name := rs.Primary.Attributes["name"]
//
//		var (
//			err  error
//			user uaaapi.UAAUser
//		)
//
//		user, err = session.UserManager().FindByUsername(name)
//		if err != nil {
//			return err
//		}
//		if err := util.AssertSame(user.ID, id); err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
