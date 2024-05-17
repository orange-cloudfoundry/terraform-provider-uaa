package uaa

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/orange-cloudfoundry/terraform-provider-uaa/uaa/uaaapi"
)

func resourceUser() *schema.Resource {

	return &schema.Resource{

		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"origin": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "uaa",
			},
			"given_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"family_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceStringHash,
			},
		},
	}
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
	return resourceUserUpdate(d, NewResourceMeta{meta})
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
	d.Set("groups", schema.NewSet(resourceStringHash, groups))

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {

	var (
		newResource bool
		session     *uaaapi.Session
	)

	if m, ok := meta.(NewResourceMeta); ok {
		session = m.meta.(*uaaapi.Session)
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
		u, _, name := getResourceChange("name", d)
		updateUserDetail = updateUserDetail || u
		u, _, givenName := getResourceChange("given_name", d)
		updateUserDetail = updateUserDetail || u
		u, _, familyName := getResourceChange("family_name", d)
		updateUserDetail = updateUserDetail || u
		u, _, email := getResourceChange("email", d)
		updateUserDetail = updateUserDetail || u
		if updateUserDetail {
			user, err := um.UpdateUser(id, name, givenName, familyName, email)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("User updated: %# v", user)
		}

		updatePassword, oldPassword, newPassword := getResourceChange("password", d)
		if updatePassword {
			err := um.ChangePassword(id, oldPassword, newPassword)
			if err != nil {
				return err
			}
			session.Log.DebugMessage("Password for user with id '%s' and name %s' updated.", id, name)
		}
	}

	old, cur := d.GetChange("groups")
	rolesToDelete, rolesToAdd := getListChanges(old, cur)

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
