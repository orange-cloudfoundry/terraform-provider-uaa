package client

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/client/fields"
	"github.com/jlpospisil/terraform-provider-uaa/util"
)

var clientSchema = map[string]*schema.Schema{
	fields.ClientId.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.ClientSecret.String(): {
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	},
	fields.AuthorizedGrantTypes.String(): {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.RedirectUri.String(): {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.Scope.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.ResourceIds.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.Authorities.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.AutoApprove.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.AccessTokenValidity.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	fields.RefreshTokenValidity.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	fields.AllowProviders.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
	fields.Name.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.TokenSalt.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.CreatedWith.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.ApprovalsDeleted.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	fields.RequiredUserGroups.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      util.ResourceStringHash,
	},
}

// dataSourceSchema is the same as the resource schema but only the client Id is required; all other fields are optional
var dataSourceSchema = mapSchemaForDataSource()

func mapSchemaForDataSource() map[string]*schema.Schema {

	dsSchema := map[string]*schema.Schema{}

	for k, v := range clientSchema {
		isClientId := k == fields.ClientId.String()
		dsSchema[k] = &schema.Schema{
			Type:     v.Type,
			Required: isClientId,
			Optional: !isClientId,
			Elem:     v.Elem,
		}
	}

	return dsSchema
}
