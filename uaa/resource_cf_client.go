package uaa

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-uaa/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-uaa/uaa/api"
)

func ResourceClient() *schema.Resource {

	return &schema.Resource{
		Create: resourceClientCreate,
		Read:   resourceClientRead,
		Update: resourceClientUpdate,
		Delete: resourceClientDelete,

		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"authorized_grant_types": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"redirect_uri": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"resource_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"authorities": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"autoapprove": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"access_token_validity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"refresh_token_validity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allowedproviders": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_salt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"createdwith": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"approvals_deleted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"required_user_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
		},
	}
}

func resourceClientCreate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*api.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	client := api.UAAClient{
		ClientID:             d.Get("client_id").(string),
		ClientSecret:         d.Get("client_secret").(string),
		AuthorizedGrantTypes: toStrings(d.Get("authorized_grant_types")),
		RedirectURI:          toStrings(d.Get("redirect_uri")),
		ResourceIds:          toStrings(d.Get("resource_ids")),
		Authorities:          toStrings(d.Get("authorities")),
		Autoapprove:          toStrings(d.Get("autoapprove")),
		Allowedproviders:     toStrings(d.Get("allowedproviders")),
		RequiredUserGroups:   toStrings(d.Get("required_user_groups")),
		Scope:                toStrings(d.Get("scope")),
		AccessTokenValidity:  d.Get("access_token_validity").(int),
		RefreshTokenValidity: d.Get("refresh_token_validity").(int),
		Name:                 d.Get("name").(string),
		TokenSalt:            d.Get("token_salt").(string),
		CreatedWith:          d.Get("createdwith").(string),
		ApprovalsDeleted:     d.Get("approvals_deleted").(bool),
	}

	um := session.ClientManager()
	client, err := um.Create(client)
	if err != nil {
		return err
	}
	session.Log.DebugMessage("New client created: %# v", client)

	d.SetId(client.ClientID)
	return resourceClientUpdate(d, NewResourceMeta{meta})
}

func resourceClientRead(d *schema.ResourceData, meta interface{}) error {
	session := meta.(*api.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	um := session.ClientManager()
	id := d.Id()

	client, err := um.GetClient(id)
	if err != nil {
		d.SetId("")
		return err
	}
	session.Log.DebugMessage("Client with ID '%s' retrieved: %# v", id, client)

	if !client.HasDefaultScope() {
		d.Set("scope", schema.NewSet(util.ResourceStringHash, toInterface(client.Scope)))
	}

	if !client.HasDefaultAuthorites() {
		d.Set("authorities", schema.NewSet(util.ResourceStringHash, toInterface(client.Authorities)))
	}

	if !client.HasDefaultResourceIds() {
		d.Set("resource_ids", schema.NewSet(util.ResourceStringHash, toInterface(client.ResourceIds)))
	}

	d.Set("client_id", client.ClientID)
	d.Set("authorized_grant_types", schema.NewSet(util.ResourceStringHash, toInterface(client.AuthorizedGrantTypes)))
	d.Set("redirect_uri", schema.NewSet(util.ResourceStringHash, toInterface(client.RedirectURI)))
	d.Set("autoapprove", schema.NewSet(util.ResourceStringHash, toInterface(client.Autoapprove)))
	d.Set("access_token_validity", client.AccessTokenValidity)
	d.Set("refresh_token_validity", client.RefreshTokenValidity)
	d.Set("allowedproviders", schema.NewSet(util.ResourceStringHash, toInterface(client.Allowedproviders)))
	d.Set("name", client.Name)
	d.Set("token_salt", client.TokenSalt)
	d.Set("createdwith", client.CreatedWith)
	d.Set("approvals_deleted", client.ApprovalsDeleted)
	d.Set("required_user_groups", schema.NewSet(util.ResourceStringHash, toInterface(client.RequiredUserGroups)))

	return nil
}

func resourceClientUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		newResource bool
		session     *api.Session
	)

	if m, ok := meta.(NewResourceMeta); ok {
		session = m.Meta.(*api.Session)
		newResource = true
	} else {
		session = meta.(*api.Session)
		if session == nil {
			return fmt.Errorf("client is nil")
		}
		newResource = false
	}

	id := d.Id()
	um := session.ClientManager()

	if !newResource {
		u := false
		name := util.GetChangedValueString("name", &u, d)
		salt := util.GetChangedValueString("token_salt", &u, d)
		created := util.GetChangedValueString("createdwith", &u, d)
		providers := util.GetChangedValueStringList("allowedproviders", &u, d)
		grants := util.GetChangedValueStringList("authorized_grant_types", &u, d)
		uris := util.GetChangedValueStringList("redirect_uri", &u, d)
		scope := util.GetChangedValueStringList("scopes", &u, d)
		resources := util.GetChangedValueStringList("resource_ids", &u, d)
		authorities := util.GetChangedValueStringList("authorities", &u, d)
		groups := util.GetChangedValueStringList("required_user_groups", &u, d)
		autoapprove := util.GetChangedValueStringList("autoapprove", &u, d)
		accestok := util.GetChangedValueInt("access_token_validity", &u, d)
		refreshtok := util.GetChangedValueInt("refresh_token_validity", &u, d)
		approval := util.GetChangedValueBool("approvals_deleted", &u, d)

		if u {
			client := api.UAAClient{
				ClientID:             id,
				AuthorizedGrantTypes: *grants,
				RedirectURI:          *uris,
				Scope:                *scope,
				ResourceIds:          *resources,
				Authorities:          *authorities,
				Autoapprove:          *autoapprove,
				AccessTokenValidity:  *accestok,
				RefreshTokenValidity: *refreshtok,
				Allowedproviders:     *providers,
				Name:                 *name,
				TokenSalt:            *salt,
				CreatedWith:          *created,
				ApprovalsDeleted:     *approval,
				RequiredUserGroups:   *groups,
			}
			nclient, err := um.UpdateClient(&client)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("Client updated: %# v", nclient)
		}

		updateSecret, oldSecret, newSecret := util.GetResourceChange("client_secret", d)
		if updateSecret {
			err := um.ChangeSecret(id, oldSecret, newSecret)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("Secret for client with id '%s' updated.", id)
		}
	}

	return nil
}

func resourceClientDelete(d *schema.ResourceData, meta interface{}) error {
	session := meta.(*api.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	um := session.ClientManager()
	um.DeleteClient(id) //nolint error is authorized here to allow not existing to be deleted without error
	return nil
}

func toStrings(data interface{}) (res []string) {
	for _, val := range data.(*schema.Set).List() {
		res = append(res, val.(string))
	}
	return
}

func toInterface(data []string) (res []interface{}) {
	for _, val := range data {
		res = append(res, val)
	}
	return
}
