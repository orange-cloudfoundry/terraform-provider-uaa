package uaa

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

func resourceClient() *schema.Resource {

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
				Set:      resourceStringHash,
			},
			"redirect_uri": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"resource_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"authorities": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"autoapprove": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
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
				Set:      resourceStringHash,
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
				Set:      resourceStringHash,
			},
		},
	}
}

func resourceClientCreate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	client := uaaapi.UAAClient{
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
	session := meta.(*uaaapi.Session)
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

	if false == client.HasDefaultScope() {
		d.Set("scope", schema.NewSet(resourceStringHash, toInterface(client.Scope)))
	}

	if false == client.HasDefaultAuthorites() {
		d.Set("authorities", schema.NewSet(resourceStringHash, toInterface(client.Authorities)))
	}

	if false == client.HasDefaultResourceIds() {
		d.Set("resource_ids", schema.NewSet(resourceStringHash, toInterface(client.ResourceIds)))
	}

	d.Set("client_id", client.ClientID)
	d.Set("authorized_grant_types", schema.NewSet(resourceStringHash, toInterface(client.AuthorizedGrantTypes)))
	d.Set("redirect_uri", schema.NewSet(resourceStringHash, toInterface(client.RedirectURI)))
	d.Set("autoapprove", schema.NewSet(resourceStringHash, toInterface(client.Autoapprove)))
	d.Set("access_token_validity", client.AccessTokenValidity)
	d.Set("refresh_token_validity", client.RefreshTokenValidity)
	d.Set("allowedproviders", schema.NewSet(resourceStringHash, toInterface(client.Allowedproviders)))
	d.Set("name", client.Name)
	d.Set("token_salt", client.TokenSalt)
	d.Set("createdwith", client.CreatedWith)
	d.Set("approvals_deleted", client.ApprovalsDeleted)
	d.Set("required_user_groups", schema.NewSet(resourceStringHash, toInterface(client.RequiredUserGroups)))

	return nil
}

func resourceClientUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		newResource bool
		session     *uaaapi.Session
	)

	if m, ok := meta.(NewResourceMeta); ok {
		session = m.meta.(*uaaapi.Session)
		newResource = true
	} else {
		session = meta.(*uaaapi.Session)
		if session == nil {
			return fmt.Errorf("client is nil")
		}
		newResource = false
	}

	id := d.Id()
	um := session.ClientManager()

	if !newResource {
		u := false
		name := getChangedValueString("name", &u, d)
		salt := getChangedValueString("token_salt", &u, d)
		created := getChangedValueString("createdwith", &u, d)
		providers := getChangedValueStringList("allowedproviders", &u, d)
		grants := getChangedValueStringList("authorized_grant_types", &u, d)
		uris := getChangedValueStringList("redirect_uri", &u, d)
		scope := getChangedValueStringList("scopes", &u, d)
		resources := getChangedValueStringList("resource_ids", &u, d)
		authorities := getChangedValueStringList("authorities", &u, d)
		groups := getChangedValueStringList("required_user_groups", &u, d)
		autoapprove := getChangedValueStringList("autoapprove", &u, d)
		accestok := getChangedValueInt("access_token_validity", &u, d)
		refreshtok := getChangedValueInt("refresh_token_validity", &u, d)
		approval := getChangedValueBool("approvals_deleted", &u, d)

		if u {
			client := uaaapi.UAAClient{
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

		updateSecret, oldSecret, newSecret := getResourceChange("client_secret", d)
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
	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	um := session.ClientManager()
	um.DeleteClient(id)
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
