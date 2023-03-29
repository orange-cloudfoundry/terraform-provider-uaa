package uaa

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

func dataSourceGroup() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, meta interface{}) (err error) {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	gm := session.GroupManager()

	var (
		displayName string
		group       uaaapi.UAAGroup
	)

	displayName = d.Get("display_name").(string)

	group, err = gm.FindByDisplayName(displayName)
	if err != nil {
		return
	}

	d.SetId(group.ID)
	d.Set("description", group.Description)
	d.Set("zone_id", group.ZoneId)

	return
}
