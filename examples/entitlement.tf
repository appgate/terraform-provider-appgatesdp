
# data "appgate_site" "default_site" {
#   site_name = "Default site"
# }

# data "appgate_condition" "always" {
#   condition_name = "Always"
# }



# resource "appgate_entitlement" "ping_entitlement" {
#   name = "test entitlement"
#   site = data.appgate_site.default_site.id
#   conditions = [
#     data.appgate_condition.always.id
#   ]

#   tags = [
#     "terraform",
#     "api-created"
#   ]
#   disabled = true

#   condition_logic = "and"
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

#   app_shortcuts {
#     name       = "ping"
#     url        = "https://www.google.com"
#     color_code = 5
#   }

#   app_shortcut_scripts = [
#     "313464a6-9dcb-4c6e-90fc-28dceaecb0a1"
#   ]

# }
