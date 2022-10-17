package user

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"github.com/terraform-providers/terraform-provider-uaa/util"
)

var Resource = &schema.Resource{
	Schema: Schema,
	Create: resourceUserCreate,
	Read:   resourceUserRead,
	Update: resourceUserUpdate,
	Delete: resourceUserDelete,
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	name := d.Get("name").(string)
	password := d.Get("password").(string)
	origin := d.Get("origin").(string)
	givenName := d.Get("given_name").(string)
	familyName := d.Get("family_name").(string)

	email := name
	if val, ok := d.GetOk("email"); ok {
		email = val.(string)
	} else {
		d.Set("email", email)
	}

	um := session.UserManager()
	user, err := um.CreateUser(name, password, origin, givenName, familyName, email)
	if err != nil {
		return err
	}
	session.Log.DebugMessage("New user created: %# v", user)

	d.SetId(user.ID)
	return resourceUserUpdate(d, uaa.NewResourceMeta{
		Meta: meta,
	})
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	um := session.UserManager()
	id := d.Id()

	user, err := um.GetUser(id)
	if err != nil {
		d.SetId("")
		return err
	}
	session.Log.DebugMessage("User with GUID '%s' retrieved: %# v", id, user)

	d.Set("name", user.Username)
	d.Set("origin", user.Origin)
	d.Set("given_name", user.Name.GivenName)
	d.Set("family_name", user.Name.FamilyName)
	d.Set("email", user.Emails[0].Value)

	var groups []interface{}
	for _, g := range user.Groups {
		if !um.IsDefaultGroup(g.Display) {
			groups = append(groups, g.Display)
		}
	}
	d.Set("groups", schema.NewSet(util.ResourceStringHash, groups))

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {

	var (
		newResource bool
		session     *uaaapi.Session
	)

	if m, ok := meta.(uaa.NewResourceMeta); ok {
		session = m.Meta.(*uaaapi.Session)
		newResource = true
	} else {
		session = meta.(*uaaapi.Session)
		if session == nil {
			return fmt.Errorf("client is nil")
		}
		newResource = false
	}

	id := d.Id()
	um := session.UserManager()

	if !newResource {

		updateUserDetail := false
		u, _, name := util.GetResourceChange("name", d)
		updateUserDetail = updateUserDetail || u
		u, _, givenName := util.GetResourceChange("given_name", d)
		updateUserDetail = updateUserDetail || u
		u, _, familyName := util.GetResourceChange("family_name", d)
		updateUserDetail = updateUserDetail || u
		u, _, email := util.GetResourceChange("email", d)
		updateUserDetail = updateUserDetail || u
		if updateUserDetail {
			user, err := um.UpdateUser(id, name, givenName, familyName, email)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("User updated: %# v", user)
		}

		updatePassword, oldPassword, newPassword := util.GetResourceChange("password", d)
		if updatePassword {
			err := um.ChangePassword(id, oldPassword, newPassword)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("Password for user with id '%s' and name %s' updated.", id, name)
		}
	}

	old, new := d.GetChange("groups")
	rolesToDelete, rolesToAdd := util.GetListChanges(old, new)

	if len(rolesToDelete) > 0 || len(rolesToAdd) > 0 {
		err := um.UpdateRoles(id, rolesToDelete, rolesToAdd, d.Get("origin").(string))
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {

	session := meta.(*uaaapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	um := session.UserManager()
	um.DeleteUser(id) //nolint error is authorized here to allow not existing to be deleted without error

	return nil
}
