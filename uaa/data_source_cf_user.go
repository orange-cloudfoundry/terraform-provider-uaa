package uaa

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

func dataSourceUser() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) (err error) {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	um := session.UserManager()

	var (
		name string
		user uaaapi.UAAUser
	)

	name = d.Get("name").(string)

	user, err = um.FindByUsername(name)
	if err != nil {
		return
	}

	d.SetId(user.ID)
	return
}
