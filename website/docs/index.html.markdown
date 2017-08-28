---
layout: "uaa"
page_title: "Provider: Cloud Foundry UAA"
sidebar_current: "docs-uaa-index"
description: |-
  The Cloud Foundry UAA provider is used to manage users, clients and roles managed by a UAA service.
---

# UAA Provider

The UAA provider is used to manage OAuth users, clients and roles managed by a Cloud Foundry [User Authentication and Authorization](https://github.com/cloudfoundry/uaa) (UAA) service. This provider manages only the UAA service and does not depend on the availability of a running Cloud Foundry environment (i.e. the Cloud Controller service).

Use the navigation to the left to read about the available resources.

## Example Usage

```
# Set the variable values in *.tfvars file
# or using -var="api_url=..." CLI option

variable "uaa_auth_url" {}
variable "uaa_login_url" {}
variable "uaa_client_secret" {}

# Configure the Cloud Foundry UAA Provider

provider "uaa" {
    
    login_url = "${var.uaa_login_url}"
    auth_url = "${var.uaa_auth_url}"

    client_id = "admin"
    client_secret = "${var.uaa_client_secret}"

    skip_ssl_validation = true
}
```

## Argument Reference

The following arguments are supported:

* `login_url` - (Required) Login/token endpoint (e.g. https://login.local.pcfdev.io). This can also be specified
  with the `UAA_LOGIN_URL` shell environment variable.

* `auth_url` - (Required) Authorization endpoint (e.g. https://uaa.local.pcfdev.io). This can also be specified
  with the `UAA_AUTH_URL` shell environment variable.

* `client_id` - (Optional) The UAA admin client ID. Defaults to "admin". This can also be specified
  with the `UAA_CLIENT_ID` shell environment variable.

  > Managing UAA resources require a UAA client with an admin privileges.

* `client_secret` - (Required) This secret of the UAA client. This can also be specified
  with the `UAA_CLIENT_SECRET` shell environment variable.

* `skip_ssl_validation` - (Optional) Skip verification of the API endpoint - Not recommended!. Defaults to "false". This can also be specified
  with the `UAA_SKIP_SSL_VALIDATION` shell environment variable.
