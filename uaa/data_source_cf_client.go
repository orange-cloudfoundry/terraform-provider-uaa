package uaa

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-uaa/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

func DataSourceClient() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceClientRead,

		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_grant_types": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      util.ResourceStringHash,
			},
			"redirect_uri": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
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
	d.Set("scope", schema.NewSet(util.ResourceStringHash, toInterface(client.Scope)))
	d.Set("authorities", schema.NewSet(util.ResourceStringHash, toInterface(client.Authorities)))
	d.Set("resource_ids", schema.NewSet(util.ResourceStringHash, toInterface(client.ResourceIds)))
	d.Set("authorized_grant_types", schema.NewSet(util.ResourceStringHash, toInterface(client.AuthorizedGrantTypes)))
	d.Set("redirect_uri", schema.NewSet(util.ResourceStringHash, toInterface(client.RedirectURI)))
	d.Set("autoapprove", schema.NewSet(util.ResourceStringHash, toInterface(client.Autoapprove)))
	d.Set("allowedproviders", schema.NewSet(util.ResourceStringHash, toInterface(client.Allowedproviders)))
	d.Set("required_user_groups", schema.NewSet(util.ResourceStringHash, toInterface(client.RequiredUserGroups)))
	d.Set("client_id", client.ClientID)
	d.Set("access_token_validity", client.AccessTokenValidity)
	d.Set("refresh_token_validity", client.RefreshTokenValidity)
	d.Set("name", client.Name)
	d.Set("token_salt", client.TokenSalt)
	d.Set("createdwith", client.CreatedWith)
	d.Set("approvals_deleted", client.ApprovalsDeleted)
	return
}
