package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/client/fields"
	"github.com/jlpospisil/terraform-provider-uaa/util"
)

var Resource = &schema.Resource{
	Schema:        clientSchema,
	CreateContext: createResource,
	ReadContext:   readResource,
	UpdateContext: updateResource,
	DeleteContext: deleteResource,
}

func createResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	client := api.UAAClient{
		ClientID:             data.Get(fields.ClientId.String()).(string),
		ClientSecret:         data.Get(fields.ClientSecret.String()).(string),
		AuthorizedGrantTypes: util.ToStringsSlice(data.Get(fields.AuthorizedGrantTypes.String())),
		RedirectURI:          util.ToStringsSlice(data.Get(fields.RedirectUri.String())),
		ResourceIds:          util.ToStringsSlice(data.Get(fields.ResourceIds.String())),
		Authorities:          util.ToStringsSlice(data.Get(fields.Authorities.String())),
		AutoApprove:          util.ToStringsSlice(data.Get(fields.AutoApprove.String())),
		AllowedProviders:     util.ToStringsSlice(data.Get(fields.AllowProviders.String())),
		RequiredUserGroups:   util.ToStringsSlice(data.Get(fields.RequiredUserGroups.String())),
		Scope:                util.ToStringsSlice(data.Get(fields.Scope.String())),
		AccessTokenValidity:  data.Get(fields.AccessTokenValidity.String()).(int),
		RefreshTokenValidity: data.Get(fields.RefreshTokenValidity.String()).(int),
		Name:                 data.Get(fields.Name.String()).(string),
		TokenSalt:            data.Get(fields.TokenSalt.String()).(string),
		CreatedWith:          data.Get(fields.CreatedWith.String()).(string),
		ApprovalsDeleted:     data.Get(fields.ApprovalsDeleted.String()).(bool),
	}

	um := session.ClientManager()
	client, err := um.Create(client)
	if err != nil {
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("New client created: %# v", client)

	data.SetId(client.ClientID)

	return nil
}

func readResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	um := session.ClientManager()
	id := data.Id()

	client, err := um.GetClient(id)
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("Client with ID '%s' retrieved: %# v", id, client)

	if !client.HasDefaultScope() {
		data.Set(fields.Scope.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.Scope)))
	}

	if !client.HasDefaultAuthorites() {
		data.Set(fields.Authorities.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.Authorities)))
	}

	if !client.HasDefaultResourceIds() {
		data.Set(fields.ResourceIds.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.ResourceIds)))
	}

	data.Set(fields.ClientId.String(), client.ClientID)
	data.Set(fields.AuthorizedGrantTypes.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AuthorizedGrantTypes)))
	data.Set(fields.RedirectUri.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.RedirectURI)))
	data.Set(fields.AutoApprove.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AutoApprove)))
	data.Set(fields.AccessTokenValidity.String(), client.AccessTokenValidity)
	data.Set(fields.RefreshTokenValidity.String(), client.RefreshTokenValidity)
	data.Set(fields.AllowProviders.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.AllowedProviders)))
	data.Set(fields.Name.String(), client.Name)
	data.Set(fields.TokenSalt.String(), client.TokenSalt)
	data.Set(fields.CreatedWith.String(), client.CreatedWith)
	data.Set(fields.ApprovalsDeleted.String(), client.ApprovalsDeleted)
	data.Set(fields.RequiredUserGroups.String(), schema.NewSet(util.ResourceStringHash, util.ToInterface(client.RequiredUserGroups)))

	return nil
}

func updateResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	um := session.ClientManager()

	isModified := false
	name := util.GetChangedValueString(fields.Name.String(), &isModified, data)
	salt := util.GetChangedValueString(fields.TokenSalt.String(), &isModified, data)
	created := util.GetChangedValueString(fields.CreatedWith.String(), &isModified, data)
	providers := util.GetChangedValueStringList(fields.AllowProviders.String(), &isModified, data)
	grants := util.GetChangedValueStringList(fields.AuthorizedGrantTypes.String(), &isModified, data)
	uris := util.GetChangedValueStringList(fields.RedirectUri.String(), &isModified, data)
	scope := util.GetChangedValueStringList(fields.Scope.String(), &isModified, data)
	resources := util.GetChangedValueStringList(fields.ResourceIds.String(), &isModified, data)
	authorities := util.GetChangedValueStringList(fields.Authorities.String(), &isModified, data)
	groups := util.GetChangedValueStringList(fields.RequiredUserGroups.String(), &isModified, data)
	autoApprove := util.GetChangedValueStringList(fields.AutoApprove.String(), &isModified, data)
	accessTokenValidity := util.GetChangedValueInt(fields.AccessTokenValidity.String(), &isModified, data)
	refreshTokenValidity := util.GetChangedValueInt(fields.RefreshTokenValidity.String(), &isModified, data)
	approval := util.GetChangedValueBool(fields.ApprovalsDeleted.String(), &isModified, data)

	if isModified {
		client := api.UAAClient{
			ClientID:             id,
			AuthorizedGrantTypes: *grants,
			RedirectURI:          *uris,
			Scope:                *scope,
			ResourceIds:          *resources,
			Authorities:          *authorities,
			AutoApprove:          *autoApprove,
			AccessTokenValidity:  *accessTokenValidity,
			RefreshTokenValidity: *refreshTokenValidity,
			AllowedProviders:     *providers,
			Name:                 *name,
			TokenSalt:            *salt,
			CreatedWith:          *created,
			ApprovalsDeleted:     *approval,
			RequiredUserGroups:   *groups,
		}
		nclient, err := um.UpdateClient(&client)
		if err != nil {
			return diag.FromErr(err)
		}
		session.Log.DebugMessage("Client updated: %# v", nclient)
	}

	updateSecret, oldSecret, newSecret := util.GetResourceChange("client_secret", data)
	if updateSecret {
		err := um.ChangeSecret(id, oldSecret, newSecret)
		if err != nil {
			return diag.FromErr(err)
		}
		session.Log.DebugMessage("Secret for client with id '%s' updated.", id)
	}

	return nil
}

func deleteResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	um := session.ClientManager()
	_ = um.DeleteClient(id)

	return nil
}
