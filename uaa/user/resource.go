package user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/user/fields"
	"github.com/jlpospisil/terraform-provider-uaa/util"
)

var Resource = &schema.Resource{
	Schema:        userSchema,
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

	name := data.Get(fields.Name.String()).(string)
	password := data.Get(fields.Password.String()).(string)
	origin := data.Get(fields.Origin.String()).(string)
	givenName := data.Get(fields.GivenName.String()).(string)
	familyName := data.Get(fields.FamilyName.String()).(string)

	email := name
	if val, ok := data.GetOk(fields.Email.String()); ok {
		email = val.(string)
	} else {
		data.Set(fields.Email.String(), email)
	}

	um := session.UserManager()
	user, err := um.CreateUser(name, password, origin, givenName, familyName, email)
	if err != nil {
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("New user created: %# v", user)

	data.SetId(user.ID)

	return updateClientRoles(um, data)
}

func readResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	um := session.UserManager()
	id := data.Id()

	user, err := um.GetUser(id)
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}
	session.Log.DebugMessage("User with GUID '%s' retrieved: %# v", id, user)

	data.Set(fields.Name.String(), user.Username)
	data.Set(fields.Origin.String(), user.Origin)
	data.Set(fields.GivenName.String(), user.Name.GivenName)
	data.Set(fields.FamilyName.String(), user.Name.FamilyName)
	data.Set(fields.Email.String(), user.Emails[0].Value)

	var groups []interface{}
	for _, g := range user.Groups {
		if !um.IsDefaultGroup(g.Display) {
			groups = append(groups, g.Display)
		}
	}
	data.Set(fields.Groups.String(), schema.NewSet(util.ResourceStringHash, groups))

	return nil
}

func updateResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	um := session.UserManager()

	isModified := false
	name := util.GetChangedValueString(fields.Name.String(), &isModified, data)
	givenName := util.GetChangedValueString(fields.GivenName.String(), &isModified, data)
	familyName := util.GetChangedValueString(fields.FamilyName.String(), &isModified, data)
	email := util.GetChangedValueString(fields.Email.String(), &isModified, data)

	if isModified {
		user, err := um.UpdateUser(id, *name, *givenName, *familyName, *email)
		if err != nil {
			return diag.FromErr(err)
		}
		session.Log.DebugMessage("User updated: %# v", user)
	}

	updatePassword, oldPassword, newPassword := util.GetResourceChange(fields.Password.String(), data)
	if updatePassword {
		err := um.ChangePassword(id, oldPassword, newPassword)
		if err != nil {
			return diag.FromErr(err)
		}
		session.Log.DebugMessage("Password for user with id '%s' and name %s' updated.", id, name)
	}

	return updateClientRoles(um, data)
}

func deleteResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	id := data.Id()
	um := session.UserManager()
	_ = um.DeleteUser(id)

	return nil
}

func updateClientRoles(um *api.UserManager, data *schema.ResourceData) diag.Diagnostics {

	origin := data.Get(fields.Origin.String()).(string)
	oldRoles, newRoles := data.GetChange(fields.Groups.String())
	rolesToDelete, rolesToAdd := util.GetListChanges(oldRoles, newRoles)

	if err := um.UpdateRoles(data.Id(), rolesToDelete, rolesToAdd, origin); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
