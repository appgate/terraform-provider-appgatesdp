# Registers this collective as a replication target of another collective
# (API v24+). The token comes from appgatesdp_replication_target on the
# replication source side. Destroying the resource unregisters it.
resource "appgatesdp_replication_source" "this" {
  token = var.replication_token
}

variable "replication_token" {
  type      = string
  sensitive = true
}
