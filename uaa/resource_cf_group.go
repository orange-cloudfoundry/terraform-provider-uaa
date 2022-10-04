package uaa

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

func resourceGroup() *schema.Resource {

	return &schema.Resource{
		Read: resourceGroupRead,
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	displayname := d.Get("display_name").(string)
	zoneId := d.Get("zone_id").(string)

	gm := session.GroupManager()
	group, err := gm.
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	gm := session.GroupManager()
	id := d.Id()

	group, err := gm.GetGroup(id)
	if err != nil {
		d.SetId("")
		return err
	}
	session.Log.DebugMessage("Group with GUID '%s' retrieved: %# v", id, group)

	d.Set("description", group.Description)
	d.Set("zone_id", group.ZoneId)

	return nil
}
