package group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/group/fields"
)

var groupSchema = map[string]*schema.Schema{
	fields.DisplayName.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.Description.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.ZoneId.String(): {
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Default:  "uaa",
	},
}
