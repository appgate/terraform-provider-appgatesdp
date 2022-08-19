---
layout: "appgatesdp"
page_title: "Version compatibility"
sidebar_current: "docs-appgate-guide-version-compatibility"
description: |-
  Version compatibility
---

## Version compatibility

The terraform provider tries to maintain full support for all [supported version](https://www.appgate.com/support/software-defined-perimeter-support). Depending on your appgate SDP appliance version, the configuration is different.
You need to specify the `client_version` if you are not running the latest supported version.

The `client_version` tries to maintain backwards compatibility 2 versions back all the time.


|                         	|  client version 14 	| client version 15 	    | client version 16   | **client version 17**     |
|-------------------------	|--------------------	|-------------------	    |-------------------	|-------------------   |
| Appgate SDP 5.3.*     	| Full support  	|    |      	  |      |
| Appgate SDP 5.4.*     	| Partial support   	| Full support  	      |    |  |
| *Appgate SDP 5.5.*   	  | Partial support   	| Partial support   	    | Full support    |    |
| **Appgate SDP 6.0.***   | Partial support   	| Partial support   	    | Partial support     | **Full support**     |




####  Terraform 0.13+

##### Example configuration for `6.0.X`

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    appgatesdp = {
      source = "appgate/appgatesdp"
      version = "0.9.0"
    }
  }
}

provider "appgatesdp" {
  provider = "local"
  client_version = 17
}
```

##### Example configuration for `5.5.X`

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    appgatesdp = {
      source = "appgate/appgatesdp"
      version = "0.7.0"
    }
  }
}

provider "appgatesdp" {
  provider = "local"
  client_version = 16
}
```

##### Example configuration for `5.4.X`

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    appgatesdp = {
      source = "appgate/appgatesdp"
      version = "0.7.0"
    }
  }
}

provider "appgatesdp" {
  provider = "local"
  client_version = 15
}
```

##### Example configuration for `5.3.X`

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    appgatesdp = {
      source = "appgate/appgatesdp"
      version = "0.7.0"
    }
  }
}

provider "appgatesdp" {
  provider = "local"
  client_version = 14
}
```

##### Example configuration for `5.2.X`

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    appgatesdp = {
      source = "appgate/appgatesdp"
      version = "0.7.0"
    }
  }
}

provider "appgatesdp" {
  provider = "local"
  client_version = 13
}
```

For additional configuration options, see [example usage](https://registry.terraform.io/providers/appgate/appgatesdp/latest/docs#example-usage).
