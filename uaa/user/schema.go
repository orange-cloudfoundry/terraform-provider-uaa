package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/user/fields"
	"github.com/terraform-providers/terraform-provider-uaa/util"
)

var Schema = map[string]*schema.Schema{

	fields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.Password.String(): {
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	},
	fields.Origin.String(): {
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Default:  "uaa",
	},
	fields.GivenName.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.FamilyName.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.Email.String(): {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
	},
	fields.Groups.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Set: util.ResourceStringHash,
	},
}
