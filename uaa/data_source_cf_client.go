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
	d.Set("scope", schema.NewSet(resourceStringHash, toInterface(client.Scope)))
	d.Set("authorities", schema.NewSet(resourceStringHash, toInterface(client.Authorities)))
	d.Set("resource_ids", schema.NewSet(resourceStringHash, toInterface(client.ResourceIds)))
	d.Set("authorized_grant_types", schema.NewSet(resourceStringHash, toInterface(client.AuthorizedGrantTypes)))
	d.Set("redirect_uri", schema.NewSet(resourceStringHash, toInterface(client.RedirectURI)))
	d.Set("autoapprove", schema.NewSet(resourceStringHash, toInterface(client.Autoapprove)))
	d.Set("allowedproviders", schema.NewSet(resourceStringHash, toInterface(client.Allowedproviders)))
	d.Set("required_user_groups", schema.NewSet(resourceStringHash, toInterface(client.RequiredUserGroups)))
	d.Set("client_id", client.ClientID)
	d.Set("access_token_validity", client.AccessTokenValidity)
	d.Set("refresh_token_validity", client.RefreshTokenValidity)
	d.Set("name", client.Name)
	d.Set("token_salt", client.TokenSalt)
	d.Set("createdwith", client.CreatedWith)
	d.Set("approvals_deleted", client.ApprovalsDeleted)
	return
}
