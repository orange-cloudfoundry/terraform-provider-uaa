package identityzone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/api"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/fields"
)

func MapIdentityZone(identityZone *api.IdentityZone, d *schema.ResourceData) diag.Diagnostics {

	d.SetId(identityZone.Id)
	d.Set(fields.IsActive.String(), identityZone.IsActive)
	d.Set(fields.SubDomain.String(), identityZone.SubDomain)

	return nil
}
