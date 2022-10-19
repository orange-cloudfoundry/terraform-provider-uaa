package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider/fields"
)

var Schema = map[string]*schema.Schema{
	fields.LoginEndpoint.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_LOGIN_URL", ""),
	},
	fields.AuthEndpoint.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_AUTH_URL", ""),
	},
	fields.ClientId.String(): {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_ID", ""),
	},
	fields.ClientSecret.String(): {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_SECRET", ""),
	},
	fields.CaCert.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CA_CERT", ""),
	},
	fields.SkipSslValidation.String(): {
		Type:        schema.TypeBool,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_SKIP_SSL_VALIDATION", "true"),
	},
}
