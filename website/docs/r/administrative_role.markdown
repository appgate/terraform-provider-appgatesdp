---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_administrative_role"
sidebar_current: "docs-appgate-resource-administrative_role"
description: |-
   Create a new Administrative Role.
---

# appgatesdp_administrative_role

Create a new Administrative Role.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.0.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl



resource "appgatesdp_administrative_role" "test_administrative_role" {
  name  = "tf-admin"
  notes = "hello world"
  tags = [
    "terraform"
  ]
  privileges {
    type         = "Create"
    target       = "Entitlement"
    default_tags = ["api-created"]
  }
}

```
## Example with data source
```hcl
## Example with data source

data "appgatesdp_site" "default_site" {
  site_name = "Default Site"
}
resource "appgatesdp_administrative_role" "administrative_role_with_scope" {
  name = "tf-admin-with-scope"
  tags = [
    "terraform"
  ]
  privileges {
    type   = "View"
    target = "Site"
    scope {
      ids  = [data.appgatesdp_site.default_site.id]
      tags = ["builtin"]
    }
  }
}


```


## Argument Reference

The following arguments are supported:


* `privileges`: (Required) Administrative privilege list.
* `administrative_role_id`: (Optional) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### privileges
Administrative privilege list.

* `type`: (Required)  Enum values: `All,View,Create,Edit,Tag,Delete,Revoke,Export,Upgrade,RenewCertificate,DownloadLogs,Test,GetUserAttributes,Backup,CheckStatus,Reevaluate,Reboot`The type of the Privilege defines the possible administrator actions.
* `target`: (Required)  Enum values: `All,Appliance,Condition,CriteriaScript,Entitlement,AdministrativeRole,IdentityProvider,MfaProvider,IpPool,LocalUser,Policy,Site,DeviceClaimScript,EntitlementScript,RingfenceRule,ApplianceCustomization,TrustedCertificate,UserClaimScript,OtpSeed,Fido2Device,TokenRecord,Blacklist,License,UserLicense,RegisteredDevice,AllocatedIp,SessionInfo,AuditLog,AdminMessage,GlobalSetting,CaCertificate,File,FailedAuthentication,AutoUpdate,ClientConnection`The target of the Privilege defines the possible target objects for that type.
* `scope`:  (Optional) The scope of the Privilege. Only applicable to certain type-target combinations. Some types depends on the IdP/MFA type, such as GetUserAttributes. This field must be omitted if not applicable.
* `default_tags`:  (Optional) The items in this list would be added automatically to the newly created objects' tags. Only applicable on "Create" type and targets with tagging capability. This field must be omitted if not applicable.
* `functions`:  (Optional) Privilege for changing Appliance Functions. Only applicable on "`AssignFunction`" type with Appliance or All target. This field must be omitted if not applicable.

#### scope

The scope of the Privilege. Only applicable to certain type-target combinations. Some types depend on the IdP/MFA type, such as GetUserAttributes. This field must be omitted if not applicable.

* `all`:  (Computed) 'If "true", all objects are accessible. For example, "type: Edit - target: Condition - scope.all: true" means the administrator can edit all Conditions in the system.'
* `ids`:  (Optional) Specific object IDs this Privilege would have access to.
* `tags`:  (Optional) Object tags this privilege would have access to.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_administrative_role.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
