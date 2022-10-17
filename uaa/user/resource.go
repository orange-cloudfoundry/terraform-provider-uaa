package user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/api"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/user/fields"
	"github.com/terraform-providers/terraform-provider-uaa/util"
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
	diagErr := updateResource(ctx, data, uaa.NewResourceMeta{
		Meta: i,
	})
	if diagErr != nil {
		return diagErr
	}

	return nil
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

	var (
		newResource bool
		session     *api.Session
	)

	if m, ok := i.(uaa.NewResourceMeta); ok {
		session = m.Meta.(*api.Session)
		newResource = true
	} else {
		session = i.(*api.Session)
		if session == nil {
			return diag.Errorf("client is nil")
		}
		newResource = false
	}

	id := data.Id()
	um := session.UserManager()

	if !newResource {

		updateUserDetail := false
		u, _, name := util.GetResourceChange(fields.Name.String(), data)
		updateUserDetail = updateUserDetail || u
		u, _, givenName := util.GetResourceChange(fields.GivenName.String(), data)
		updateUserDetail = updateUserDetail || u
		u, _, familyName := util.GetResourceChange(fields.FamilyName.String(), data)
		updateUserDetail = updateUserDetail || u
		u, _, email := util.GetResourceChange(fields.Email.String(), data)
		updateUserDetail = updateUserDetail || u
		if updateUserDetail {
			user, err := um.UpdateUser(id, name, givenName, familyName, email)
			if err != nil {
				return diag.FromErr(err)
			}
			session.Log.DebugMessage("User updated: %# v", user)
		}

		updatePassword, oldPassword, newPassword := util.GetResourceChange("password", data)
		if updatePassword {
			err := um.ChangePassword(id, oldPassword, newPassword)
			if err != nil {
				return diag.FromErr(err)
			}
			session.Log.DebugMessage("Password for user with id '%s' and name %s' updated.", id, name)
		}
	}

	oldUser, newUser := data.GetChange(fields.Groups.String())
	rolesToDelete, rolesToAdd := util.GetListChanges(oldUser, newUser)

	if len(rolesToDelete) > 0 || len(rolesToAdd) > 0 {
		err := um.UpdateRoles(id, rolesToDelete, rolesToAdd, data.Get(fields.Origin.String()).(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
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
