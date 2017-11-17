---
layout: "uaa"
page_title: "Cloud Foundry UAA: uaa_client"
sidebar_current: "docs-uaa-resource-client"
description: |-
  Provides a Cloud Foundry UAA Client resource.
---

# uaa\_client

Provides a resource for managing Cloud Foundry UUA clients.

## Example Usage

The following example creates a client.

```
resource "uaa_client" "admin-service-client" {
    client_id = "admin-client"
    client_secret = "mysecret"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    scope = ["uaa.admin", "openid"]
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) Client identifier, unique within identity zone.
* `authorized_grant_types` - (Required) List of grant types that can be used to obtain a token with this client. Can include authorization_code, password, implicit, and/or client_credentials.
* `redirect_uri` - (Required) Allowed URI pattern for redirect during authorization. Wildcard patterns can be specified using the Ant-style pattern.
* `scope` - (Optional) Scopes allowed for the client.
* `resource_ids` - (Optional) Resources the client is allowed access to.
* `authorities` - (Optional) Scopes which the client is able to grant when creating a client.
* `autoapprove` - (Optional) Scopes that do not require user approval.
* `access_token_validity` - (Optional) time in seconds to access token expiration after it is issued.
* `refresh_token_validity` - (Optional) time in seconds to refresh token expiration after it is issued.
* `allowedproviders` - (Optional) A list of origin keys (alias) for identity providers the client is limited to.
* `name` - (Optional) A human readable name for the client.
* `token_salt` - (Optional) A random string used to generate the client's revokation key. Change this value to revoke all active tokens for the client.
* `createdwith` - (Optional) What scope the bearer token had when client was created.
* `approvals_deleted` - (Optional) Were the approvals deleted for the client, and an audit event sent.
* `required_user_groups` - (Optional) A list of group names.
* `client_secret` - (Required if the client allows authorization_code or client_credentials grant type) A secret string used for authenticating as this client.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the Client
