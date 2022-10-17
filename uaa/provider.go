package uaa

import (
	"context"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/user"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {

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
			"uaa_user":   user.DataSource,
			"uaa_client": dataSourceClient(),
			"uaa_group":  dataSourceGroup(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"uaa_user":   user.Resource,
			"uaa_client": resourceClient(),
			"uaa_group":  resourceGroup(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		loginEndpoint:     d.Get("login_endpoint").(string),
		authEndpoint:      d.Get("auth_endpoint").(string),
		clientID:          d.Get("client_id").(string),
		clientSecret:      d.Get("client_secret").(string),
		caCert:            d.Get("ca_cert").(string),
		skipSslValidation: d.Get("skip_ssl_validation").(bool),
	}
	client, err := config.Client()
	if err != nil {
		return client, diag.FromErr(err)
	}
	return client, nil
}
