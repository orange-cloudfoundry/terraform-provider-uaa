package uaa

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

func dataSourceClient() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceClientRead,

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_grant_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
			"redirect_uri": {
				Type:     schema.TypeSet,
				Optional: true,
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

func dataSourceClientRead(d *schema.ResourceData, meta interface{}) (err error) {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	um := session.ClientManager()

	var (
		id     string
		client uaaapi.UAAClient
	)

	id = d.Get("client_id").(string)
	client, err = um.FindByClientID(id)
	if err != nil {
		return
	}

	d.SetId(client.ClientID)
	err = d.Set("scope", schema.NewSet(resourceStringHash, toInterface(client.Scope)))
	if err != nil {
		return fmt.Errorf("error while setting scope data to client '%s'", err)
	}
	err = d.Set("authorities", schema.NewSet(resourceStringHash, toInterface(client.Authorities)))
	if err != nil {
		return fmt.Errorf("error while authorities scope data to client '%s'", err)
	}
	err = d.Set("resource_ids", schema.NewSet(resourceStringHash, toInterface(client.ResourceIds)))
	if err != nil {
		return fmt.Errorf("error while setting resource_ids data to client '%s'", err)
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
	err = d.Set("allowedproviders", schema.NewSet(resourceStringHash, toInterface(client.Allowedproviders)))
	if err != nil {
		return fmt.Errorf("error while setting allowedproviders data to client '%s'", err)
	}
	err = d.Set("required_user_groups", schema.NewSet(resourceStringHash, toInterface(client.RequiredUserGroups)))
	if err != nil {
		return fmt.Errorf("error while setting required_user_groups data to client '%s'", err)
	}
	err = d.Set("client_id", client.ClientID)
	if err != nil {
		return fmt.Errorf("error while setting client_id data to client '%s'", err)
	}
	err = d.Set("access_token_validity", client.AccessTokenValidity)
	if err != nil {
		return fmt.Errorf("error while setting access_token_validity data to client '%s'", err)
	}
	err = d.Set("refresh_token_validity", client.RefreshTokenValidity)
	if err != nil {
		return fmt.Errorf("error while setting refresh_token_validity data to client '%s'", err)
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
	return
}
