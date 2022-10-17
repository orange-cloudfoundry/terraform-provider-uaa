package group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/group/fields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
)

var DataSource = &schema.Resource{
	Schema:      groupSchema,
	ReadContext: readDataSource,
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*uaaapi.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	gm := session.GroupManager()

	var (
		displayName string
		group       *uaaapi.UAAGroup
	)

	displayName = data.Get(fields.DisplayName.String()).(string)
	zoneId := data.Get(fields.ZoneId.String()).(string)

	group, err := gm.FindByDisplayName(displayName, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(group.ID)
	data.Set("description", group.Description)
	data.Set("zone_id", group.ZoneId)

	return nil
}
