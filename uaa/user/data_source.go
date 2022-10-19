package user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/user/fields"
)

var DataSource = &schema.Resource{
	Schema:      dataSourceSchema,
	ReadContext: readDataSource,
}

var dataSourceSchema = map[string]*schema.Schema{
	fields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	um := session.UserManager()
	name := data.Get(fields.Name.String()).(string)

	user, err := um.FindByUsername(name)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(user.ID)

	return nil
}
