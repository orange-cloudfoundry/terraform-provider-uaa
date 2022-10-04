---
page_title: "Cloud Foundry UAA: uaa_group"
---

# User Resource

Provides a resource for managing Cloud Foundry UUA groups.

## Example Usage

The following example creates a group.

```
resource uaa_group "group" {
    display_name = "resource.read"
    description  = "Read API resources"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The name of the group to look up
* `description` - (Optional) The description of the group
* `zone_id` - (Optional) The identity zone that the group belongs to

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the group
