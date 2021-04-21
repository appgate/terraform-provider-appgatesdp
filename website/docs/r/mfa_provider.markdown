---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_mfa_provider"
sidebar_current: "docs-appgate-resource-mfa_provider"
description: |-
   Create a new MFA Provider.
---

# appgatesdp_mfa_provider

Create a new MFA Provider.

## Example Usage

```hcl

resource "appgatesdp_mfa_provider" "mfa" {
   name = "hello world"
   port = 1812
   type = "Radius"
   shared_secret = "helloworld"
   challenge_shared_secret = "secretString"
   hostnames = [
      "mfa.company.com"
   ]

   tags = [
      "terraform",
      "api-created"
   ]
}

```

## Argument Reference

The following arguments are supported:


* `type`: (Required) The type of the MFA Provider. "DefaultTimeBased" is built-in, a new one cannot be created.
* `hostnames`: (Required) Hostnames/IP addresses to connect.
* `port`: (Required) Port to connect.
* `input_type`: (Optional) The input type used in the client to enter the MFA code. 
 * "Masked" - The input is masked the same way as a password field.
 * "Numeric" - The input is marked as a numeric input.
 * "Text" - The input is handled as a regular plain text field.

* `shared_secret`: (Optional) Radius shared secret to authenticate to the server.
* `authentication_protocol`: (Optional) Radius protocol to use while authenticating users.
* `timeout`: (Optional) Timeout in seconds before giving up on response.
* `mode`: (Optional) Defines the multi-factor authentication flow for RADIUS.
 * "OneFactor" - The input from the user is sent as password and the response is used for result.
 * "Challenge" - Before prompting the user, Controller sends a challenge request to the RADIUS server
 using "challengeSharedSecret" or the user password. Data from the response is used with user input to
 send the second RADIUS authentication request.
 * "Push" - "challengeSharedSecret" or the user password is sent to RADIUS which triggers an external
 authentication flow. When the external authentication flow returns success, the MFA attempt is
 authenticated.

* `use_user_password`: (Optional) -> If enabled, the Client will send the cached password instead of using challengeSharedSecret" to initiate the multi-factor authentication.
* `challenge_shared_secret`: (Optional) -> Password sent to RADIUS to initiate multi-factor authentication. Required if "useUserPassword" is not enabled.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### hostnames
Hostnames/IP addresses to connect.

### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_mfa_provider d3131f83-10d1-4abc-ac0b-7349538e8300
```
