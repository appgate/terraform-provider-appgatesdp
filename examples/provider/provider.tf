terraform {
  required_providers {
    appgatesdp = {
      source  = "appgate/appgatesdp"
      version = "~> 2.0"
    }
  }
}

provider "appgatesdp" {
  url          = "https://controller.example.com:8443/admin"
  username     = var.username
  password     = var.password
  api_version  = "v25" # SDP 7.0; omit to negotiate with the controller
  pem_filepath = "controller-ca.pem"
}

variable "username" {
  type = string
}

variable "password" {
  type      = string
  sensitive = true
}
