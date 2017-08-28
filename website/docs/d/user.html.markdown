---
layout: "uaa"
page_title: "Cloud Foundry UAA: uaa_user"
sidebar_current: "docs-uaa-datasource-user"
description: |-
  Get information on a Cloud Foundry UAA User.
---

# uaa\_user

Gets information on a Cloud Foundry UAA user.

## Example Usage

The following example looks up an user named 'myuser'. 

```
data "uaa_user" "myuser" {
    name = "myuser"    
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user to look up

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the user
