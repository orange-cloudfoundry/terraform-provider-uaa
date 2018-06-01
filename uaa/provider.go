package uaa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider -
func Provider() terraform.ResourceProvider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"login_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_LOGIN_URL", ""),
			},
			"auth_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_AUTH_URL", ""),
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_ID", ""),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_CLIENT_SECRET", ""),
			},
			"ca_cert": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_CA_CERT", ""),
			},
			"skip_ssl_validation": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UAA_SKIP_SSL_VALIDATION", "true"),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"uaa_user":   dataSourceUser(),
			"uaa_client": dataSourceClient(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"uaa_user":   resourceUser(),
			"uaa_client": resourceClient(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		loginEndpoint:     d.Get("login_endpoint").(string),
		authEndpoint:      d.Get("auth_endpoint").(string),
		clientID:          d.Get("client_id").(string),
		clientSecret:      d.Get("client_secret").(string),
		caCert:            d.Get("ca_cert").(string),
		skipSslValidation: d.Get("skip_ssl_validation").(bool),
	}
	return config.Client()
}
