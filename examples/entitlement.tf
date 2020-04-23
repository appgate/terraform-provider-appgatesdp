
# data "appgate_site" "default_site" {
#   site_name = "Default site"
# }

# data "appgate_condition" "always" {
#   condition_name = "Always"
# }



# resource "appgate_entitlement" "ping_entitlement" {
#   name = "test entitlement"
#   site = data.appgate_site.default_site.id
#   # site = appgate_site.gbg_site.id
#   conditions = [
#     data.appgate_condition.always.id
#   ]

#   actions {
#     subtype = "icmp_up"
#     action  = "allow"
#     # https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
#     types = ["0-16"]
#     hosts = [
#       "10.0.0.1",
#       "10.0.0.0/24",
#       "hostname.company.com",
#       "dns://hostname.company.com",
#       "aws://security-group:accounting"
#     ]
#   }

# }
