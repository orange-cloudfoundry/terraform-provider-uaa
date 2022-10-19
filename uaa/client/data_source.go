package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/client/fields"
	"github.com/jlpospisil/terraform-provider-uaa/util"
)

var DataSource = &schema.Resource{
	Schema:      dataSourceSchema,
	ReadContext: readDataSource,
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	um := session.ClientManager()

	var (
		id     string
		client api.UAAClient
	)

	id = data.Get(fields.ClientId.String()).(string)
	client, err := um.FindByClientID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(client.ClientID)
	data.Set(fields.Scope.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.Scope)))
	data.Set(fields.Authorities.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.Authorities)))
	data.Set(fields.ResourceIds.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.ResourceIds)))
	data.Set(fields.AuthorizedGrantTypes.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AuthorizedGrantTypes)))
	data.Set(fields.RedirectUri.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.RedirectURI)))
	data.Set(fields.AutoApprove.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AutoApprove)))
	data.Set(fields.AllowProviders.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AllowedProviders)))
	data.Set(fields.RequiredUserGroups.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.RequiredUserGroups)))
	data.Set(fields.ClientId.String(), client.ClientID)
	data.Set(fields.AccessTokenValidity.String(), client.AccessTokenValidity)
	data.Set(fields.RefreshTokenValidity.String(), client.RefreshTokenValidity)
	data.Set(fields.Name.String(), client.Name)
	data.Set(fields.TokenSalt.String(), client.TokenSalt)
	data.Set(fields.CreatedWith.String(), client.CreatedWith)
	data.Set(fields.ApprovalsDeleted.String(), client.ApprovalsDeleted)

	return nil
}
