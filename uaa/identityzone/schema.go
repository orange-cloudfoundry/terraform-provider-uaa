package identityzone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/clientsecretpolicyfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/configfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/fields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/samlconfigfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/samlkeyfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/tokenpolicyfields"
)

var Schema = map[string]*schema.Schema{
	fields.Id.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.IsActive.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	fields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.SubDomain.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.Config.String(): {
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: ConfigSchema,
		},
	},
}

var ConfigSchema = map[string]*schema.Schema{
	configfields.ClientSecretPolicy.String(): {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: ClientSecretPolicySchema,
		},
	},
	configfields.TokenPolicy.String(): {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: TokenPolicySchema,
		},
	},
	configfields.Saml.String(): {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: SamlConfigSchema,
		},
	},
}

var ClientSecretPolicySchema = map[string]*schema.Schema{
	clientsecretpolicyfields.MaxLength.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinLength.String(): {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	clientsecretpolicyfields.MinUpperCaseChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	clientsecretpolicyfields.MinLowerCaseChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	clientsecretpolicyfields.MinDigits.String(): {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	clientsecretpolicyfields.MinSpecialChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
}

var TokenPolicySchema = map[string]*schema.Schema{
	tokenpolicyfields.AccessTokenTtl.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	tokenpolicyfields.RefreshTokenTtl.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	tokenpolicyfields.IsJwtRevocable.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	tokenpolicyfields.IsRefreshTokenUnique.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	tokenpolicyfields.RefreshTokenFormat.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "jwt",
	},
	tokenpolicyfields.ActiveKeyId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
}

var SamlConfigSchema = map[string]*schema.Schema{
	samlconfigfields.ActiveKeyId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	samlconfigfields.AssertionTtlSeconds.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	samlconfigfields.Certificate.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	samlconfigfields.DisableInResponseToCheck.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.EntityId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	samlconfigfields.IsAssertionSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.IsRequestSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	//samlconfigfields.Key.String(): {
	//	Type:     schema.TypeList,
	//	Optional: true,
	//	Elem:     &schema.Resource{
	//		// how do we model this when the property name is dynamic?
	//		// do we take it in with an extra name property and handle in the mapper?
	//	},
	//},
	samlconfigfields.WantAssertionSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.WantAuthRequestSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
}

var SamlConfigKeySchema = map[string]*schema.Schema{
	samlkeyfields.Certificate.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
}

// The only required field for looking up an existing identity zone is the `id`.  All other fields should be optional
// and computed.  We can iterate over the resource schema and change those properties to avoid managing two schemas
// that are otherwise identical.
var dataSourceSchema = mapSchemaForDataSource(Schema)

func mapSchemaForDataSource(originalSchema map[string]*schema.Schema) map[string]*schema.Schema {

	dsSchema := map[string]*schema.Schema{}

	for k, v := range originalSchema {
		isZoneId := k == fields.Id.String()
		dsSchema[k] = &schema.Schema{
			Type:     v.Type,
			Required: isZoneId,
			Computed: !isZoneId,
			Elem:     v.Elem,
		}
		if v.Type == schema.TypeList {
			elemSchema := v.Elem.(*schema.Resource).Schema
			dsSchema[k].Elem = &schema.Resource{
				Schema: mapSchemaForDataSource(elemSchema),
			}
		}
	}

	return dsSchema
}
