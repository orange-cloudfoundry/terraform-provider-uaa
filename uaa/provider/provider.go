package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/user"
)

func Provider() *schema.Provider {

	return &schema.Provider{
		Schema:               Schema,
		DataSourcesMap:       DataSources,
		ResourcesMap:         Resources,
		ConfigureContextFunc: configureContext,
	}
}

var DataSources = map[string]*schema.Resource{
	"uaa_user":   user.DataSource,
	"uaa_client": uaa.DataSourceClient(),
	"uaa_group":  uaa.DataSourceGroup(),
}

var Resources = map[string]*schema.Resource{
	"uaa_user":   user.Resource,
	"uaa_client": uaa.ResourceClient(),
	"uaa_group":  uaa.ResourceGroup(),
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := uaaapi.Config{
		LoginEndpoint:     d.Get("login_endpoint").(string),
		AuthEndpoint:      d.Get("auth_endpoint").(string),
		ClientID:          d.Get("client_id").(string),
		ClientSecret:      d.Get("client_secret").(string),
		CaCert:            d.Get("ca_cert").(string),
		SkipSslValidation: d.Get("skip_ssl_validation").(bool),
	}
	client, err := config.Client()
	if err != nil {
		return client, diag.FromErr(err)
	}
	return client, nil
}
