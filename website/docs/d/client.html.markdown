---
layout: "uaa"
page_title: "Cloud Foundry UAA: uaa_client"
sidebar_current: "docs-uaa-datasource-client"
description: |-
  Get information on a Cloud Foundry UAA Client.
---

# uaa\_client

Gets information on a Cloud Foundry UAA Client.

## Example Usage

The following example looks up a client named 'myclient'.

```
data "uaa_client" "myclient" {
    client_id = "myclient"
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) The client_id of the client to look up

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the client
* `client_id` - Client identifier, unique within identity zone.
* `authorized_grant_types` - List of grant types that can be used to obtain a token with this client.
* `redirect_uri` - Allowed URI pattern for redirect during authorization.
* `scope` - Scopes allowed for the client.
* `resource_ids` - Resources the client is allowed access to.
* `authorities` - Scopes which the client is able to grant when creating a client.
* `autoapprove` - Scopes that do not require user approval.
* `access_token_validity` - time in seconds to access token expiration after it is issued.
* `refresh_token_validity` - time in seconds to refresh token expiration after it is issued.
* `allowedproviders` - A list of origin keys (alias) for identity providers the client is limited to.
* `name` - A human readable name for the client.
* `createdwith` - What scope the bearer token had when client was created.
* `token_salt` - A random string used to generate the client's revokation key.
* `approvals_deleted` - Were the approvals deleted for the client, and an audit event sent.
* `required_user_groups` - A list of group names.
