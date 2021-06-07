package uaa

import "github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

// resourceStringHash -
func resourceStringHash(v interface{}) int {
	return hashcode.String(v.(string))
}
