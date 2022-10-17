package uaa

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"github.com/terraform-providers/terraform-provider-uaa/util"
)

func ResourceGroup() *schema.Resource {

	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

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
				ForceNew: true,
				Optional: true,
				Default:  "uaa",
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	zoneId := d.Get("zone_id").(string)

	gm := session.GroupManager()
	group, err := gm.CreateGroup(displayName, description, zoneId)
	if err != nil {
		return err
	}
	session.Log.DebugMessage("New group created: %# v", group)

	d.SetId(group.ID)

	return nil
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	gm := session.GroupManager()
	id := d.Id()
	zoneId := d.Get("zone_id").(string)

	group, err := gm.GetGroup(id, zoneId)
	if err != nil {
		d.SetId("")
		return err
	}
	session.Log.DebugMessage("Group with GUID '%s' retrieved: %# v", id, group)

	d.Set("description", group.Description)
	d.Set("zone_id", group.ZoneId)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	gm := session.GroupManager()
	updateGroup := false
	changed, _, displayName := util.GetResourceChange("display_name", d)
	updateGroup = updateGroup || changed
	changed, _, description := util.GetResourceChange("description", d)
	updateGroup = updateGroup || changed
	changed, _, zoneId := util.GetResourceChange("zone_id", d)
	updateGroup = updateGroup || changed

	if updateGroup {
		group, err := gm.UpdateGroup(id, displayName, description, zoneId)
		if err != nil {
			return err
		}
		session.Log.DebugMessage("Group updated: %# v", group)
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	zoneId := d.Get("zone_id").(string)
	gm := session.GroupManager()
	err = gm.DeleteGroup(id, zoneId) //nolint error is authorized here to allow not existing to be deleted without error
	if err != nil {
		return err
	}

	return
}
