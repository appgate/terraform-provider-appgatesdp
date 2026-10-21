# Global settings always exist; destroying this resource resets them to the
# controller defaults.
resource "appgatesdp_global_settings" "this" {
  administration_token_expiration = 60
  audit_log_persistence_mode      = "Default"
  claims_token_expiration         = 1440
  entitlement_token_expiration    = 180
  vpn_certificate_expiration      = 525600
}
