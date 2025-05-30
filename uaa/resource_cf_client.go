package uaa

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

func resourceClient() *schema.Resource {

	return &schema.Resource{
		Create: resourceClientCreate,
		Read:   resourceClientRead,
		Update: resourceClientUpdate,
		Delete: resourceClientDelete,

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"authorized_grant_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"redirect_uri": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"scope": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"authorities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"autoapprove": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"access_token_validity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"refresh_token_validity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allowedproviders": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_salt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"createdwith": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"approvals_deleted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"required_user_groups": {
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

	if !client.HasDefaultScope() {
		err := d.Set("scope", schema.NewSet(resourceStringHash, toInterface(client.Scope)))
		if err != nil {
			return fmt.Errorf("error while setting scope data to client '%s'", err)
		}
	}

	if !client.HasDefaultAuthorites() {
		err := d.Set("authorities", schema.NewSet(resourceStringHash, toInterface(client.Authorities)))
		if err != nil {
			return fmt.Errorf("error while setting authorities data to client '%s'", err)
		}
	}

	if !client.HasDefaultResourceIds() {
		err := d.Set("resource_ids", schema.NewSet(resourceStringHash, toInterface(client.ResourceIds)))
		if err != nil {
			return fmt.Errorf("error while setting resource_ids data to client '%s'", err)
		}
	}

	err = d.Set("client_id", client.ClientID)
	if err != nil {
		return fmt.Errorf("error setting client_id: %s", err)
	}
	err = d.Set("authorized_grant_types", schema.NewSet(resourceStringHash, toInterface(client.AuthorizedGrantTypes)))
	if err != nil {
		return fmt.Errorf("error while setting authorized_grant_types data to client '%s'", err)
	}
	err = d.Set("redirect_uri", schema.NewSet(resourceStringHash, toInterface(client.RedirectURI)))
	if err != nil {
		return fmt.Errorf("error while setting redirect_uri data to client '%s'", err)
	}
	err = d.Set("autoapprove", schema.NewSet(resourceStringHash, toInterface(client.Autoapprove)))
	if err != nil {
		return fmt.Errorf("error while setting autoapprove data to client '%s'", err)
	}
	err = d.Set("access_token_validity", client.AccessTokenValidity)
	if err != nil {
		return fmt.Errorf("error while setting access_token_validity data to client '%s'", err)
	}
	err = d.Set("refresh_token_validity", client.RefreshTokenValidity)
	if err != nil {
		return fmt.Errorf("error while setting refresh_token_validity data to client '%s'", err)
	}
	err = d.Set("allowedproviders", schema.NewSet(resourceStringHash, toInterface(client.Allowedproviders)))
	if err != nil {
		return fmt.Errorf("error while setting allowedproviders data to client '%s'", err)
	}
	err = d.Set("name", client.Name)
	if err != nil {
		return fmt.Errorf("error while setting name data to client '%s'", err)
	}
	err = d.Set("token_salt", client.TokenSalt)
	if err != nil {
		return fmt.Errorf("error while setting token_salt data to client '%s'", err)
	}
	err = d.Set("createdwith", client.CreatedWith)
	if err != nil {
		return fmt.Errorf("error while setting createdwith data to client '%s'", err)
	}
	err = d.Set("approvals_deleted", client.ApprovalsDeleted)
	if err != nil {
		return fmt.Errorf("error while setting approvals_deleted data to client '%s'", err)
	}
	err = d.Set("required_user_groups", schema.NewSet(resourceStringHash, toInterface(client.RequiredUserGroups)))
	if err != nil {
		return fmt.Errorf("error while setting required_user_groups data to client '%s'", err)
	}

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
