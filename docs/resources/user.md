---
page_title: "Cloud Foundry UAA: uaa_user"
---

# User Resource

Provides a resource for managing Cloud Foundry UUA users. This resource provides extended functionality to attach additional UAA roles to the user.

## Example Usage

The following example creates a user and attaches additional UAA roles to grant administrator rights to that user.

```
resource "uaa_user" "admin-service-user" {
    
    name = "uaa-admin"
    password = "Passw0rd"
    
    given_name = "John"
    family_name = "Doe"

    groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user. This will also be the users login name
* `password` - (Optional) The user's password
* `origin` - (Optional) The user authentication origin. By default, this will be `UAA`. For users authenticated by LDAP this should be `ldap`
* `given_name` - (Optional) The given name of the user
* `family_name` - (Optional) The family name of the user
* `email` - (Optional) The email address of the user
* `groups` - (Optional) Any UAA `groups` / `roles` to associated the user with

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the User
* `email` - If not provided this attributed will be assigned the same value as the `name`, assuming that the username is the user's email address

