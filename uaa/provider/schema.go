package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var Schema = map[string]*schema.Schema{
	"login_endpoint": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_LOGIN_URL", ""),
	},
	"auth_endpoint": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_AUTH_URL", ""),
	},
	"client_id": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_ID", ""),
	},
	"client_secret": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_SECRET", ""),
	},
	"ca_cert": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CA_CERT", ""),
	},
	"skip_ssl_validation": {
		Type:        schema.TypeBool,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_SKIP_SSL_VALIDATION", "true"),
	},
}
