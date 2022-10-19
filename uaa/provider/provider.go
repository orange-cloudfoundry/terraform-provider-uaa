package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/client"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/group"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider/fields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/user"
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
	"uaa_client":        client.DataSource,
	"uaa_group":         group.DataSource,
	"uaa_identity_zone": identityzone.DataSource,
	"uaa_user":          user.DataSource,
}

var Resources = map[string]*schema.Resource{
	"uaa_user":   user.Resource,
	"uaa_client": client.Resource,
	"uaa_group":  group.Resource,
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := api.Config{
		LoginEndpoint:     d.Get(fields.LoginEndpoint.String()).(string),
		AuthEndpoint:      d.Get(fields.AuthEndpoint.String()).(string),
		ClientID:          d.Get(fields.ClientId.String()).(string),
		ClientSecret:      d.Get(fields.ClientSecret.String()).(string),
		CaCert:            d.Get(fields.CaCert.String()).(string),
		SkipSslValidation: d.Get(fields.SkipSslValidation.String()).(bool),
	}
	client, err := config.Client()
	if err != nil {
		return client, diag.FromErr(err)
	}
	return client, nil
}
