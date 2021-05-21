---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_administrative_role"
sidebar_current: "docs-appgate-resource-administrative_role"
description: |-
   Create a new Administrative Role.
---

# appgatesdp_administrative_role

Create a new Administrative Role.

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

### Example with scope


```hcl

data "appgatesdp_site" "default_site" {
  site_name = "Default site"
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
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### privileges
Administrative privilege list.

* `type`: (Required)  Enum values: `All,View,Create,Edit,Tag,Delete,Revoke,Export,Upgrade,RenewCertificate,DownloadLogs,Test,GetUserAttributes,Backup,CheckStatus,Reevaluate`The type of the Privilege defines the possible administrator actions.
* `target`: (Required)  Enum values: `All,Appliance,Condition,CriteriaScript,Entitlement,AdministrativeRole,IdentityProvider,MfaProvider,IpPool,LocalUser,Policy,Site,DeviceScript,EntitlementScript,RingfenceRule,ApplianceCustomization,OtpSeed,TokenRecord,Blacklist,UserLicense,RegisteredDevice,AllocatedIp,SessionInfo,AuditLog,AdminMessage,GlobalSetting,CaCertificate,File,FailedAuthentication`The target of the Privilege defines the possible target objects for that type.
* `scope`:  (Optional) The scope of the Privilege. Only applicable to certain type-target combinations. Some types depends on the IdP/MFA type, such as GetUserAttributes. This field must be omitted if not applicable.
* `default_tags`:  (Optional) The items in this list would be added automatically to the newly created objects' tags. Only applicable on "Create" type and targets with tagging capability. This field must be omitted if not applicable.
### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_administrative_role d3131f83-10d1-4abc-ac0b-7349538e8300
```
