package group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/group/fields"
	"github.com/jlpospisil/terraform-provider-uaa/util"
)

var Resource = &schema.Resource{
	Schema:        groupSchema,
	CreateContext: createResource,
	ReadContext:   readResource,
	UpdateContext: updateResource,
	DeleteContext: deleteResource,
}

func createResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	displayName := data.Get(fields.DisplayName.String()).(string)
	description := data.Get(fields.Description.String()).(string)
	zoneId := data.Get(fields.ZoneId.String()).(string)

	gm := session.GroupManager()
	group, err := gm.CreateGroup(displayName, description, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("New group created: %# v", group)

	data.SetId(group.ID)

	return nil
}

func readResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	gm := session.GroupManager()
	id := data.Id()
	zoneId := data.Get(fields.ZoneId.String()).(string)

	group, err := gm.GetGroup(id, zoneId)
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("Group with GUID '%s' retrieved: %# v", id, group)

	data.Set(fields.Description.String(), group.Description)
	data.Set(fields.ZoneId.String(), group.ZoneId)

	return nil
}

func updateResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	gm := session.GroupManager()
	updateGroup := false
	changed, _, displayName := util.GetResourceChange(fields.DisplayName.String(), data)
	updateGroup = updateGroup || changed
	changed, _, description := util.GetResourceChange(fields.Description.String(), data)
	updateGroup = updateGroup || changed
	changed, _, zoneId := util.GetResourceChange(fields.ZoneId.String(), data)
	updateGroup = updateGroup || changed

	if updateGroup {
		group, err := gm.UpdateGroup(id, displayName, description, zoneId)
		if err != nil {
			return diag.FromErr(err)
		}
		session.Log.DebugMessage("Group updated: %# v", group)
	}

	return nil
}

func deleteResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	zoneId := data.Get(fields.ZoneId.String()).(string)
	gm := session.GroupManager()
	err := gm.DeleteGroup(id, zoneId) //nolint error is authorized here to allow not existing to be deleted without error
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
