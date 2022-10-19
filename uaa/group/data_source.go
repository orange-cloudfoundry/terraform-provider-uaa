package group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/group/fields"
)

var DataSource = &schema.Resource{
	Schema:      groupSchema,
	ReadContext: readDataSource,
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	gm := session.GroupManager()

	displayName := data.Get(fields.DisplayName.String()).(string)
	zoneId := data.Get(fields.ZoneId.String()).(string)

	group, err := gm.FindByDisplayName(displayName, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(group.ID)
	data.Set(fields.Description.String(), group.Description)
	data.Set(fields.ZoneId.String(), group.ZoneId)

	return nil
}
