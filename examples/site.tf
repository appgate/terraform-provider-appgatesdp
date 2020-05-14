
# resource "appgate_site" "gbg_site" {
#   name       = "Gothenburg site"
#   short_name = "gbg"
#   tags = [
#     "developer",
#     "api-created"
#   ]

#   notes = "This object has been created for test purposes."

#   network_subnets = [
#     "10.0.0.0/16"
#   ]
#   default_gateway {
#     enabled_v4       = false
#     enabled_v6       = false
#     excluded_subnets = []
#   }

#   name_resolution {

#     dns_resolvers {
#       name            = "DNS Resolver 1"
#       update_interval = 13
#       servers = [
#         "8.8.8.8",
#         "1.1.1.1"
#       ]
#       search_domains = [
#         "hostname.dns",
#         "foo.bar"
#       ]
#     }

#     aws_resolvers {
#       name = "AWS Resolver 1"
#       regions = [
#         "eu-central-1",
#         "eu-west-1"
#       ]
#       update_interval    = 59
#       vpcs               = []
#       vpc_auto_discovery = true
#       use_iam_role       = true
#       access_key_id      = "string1"
#       secret_access_key  = "string2"
#       https_proxy                     = "username:password@appgate.com:4443"
#       resolve_with_master_credentials = true

#       assumed_roles                   = []
#     }

#     azure_resolvers {
#       name            = "Azure Resolver 1"
#       update_interval = 30
#       subscription_id = "string1"
#       tenant_id       = "string2"
#       client_id       = "string3"
#       secret_id       = "string4"
#     }

#     esx_resolvers {
#       name            = "ESX Resolver 1"
#       update_interval = 120
#       hostname        = "string1"
#       username        = "string2"
#       password        = "secret_password"
#     }

#     gcp_resolvers {
#       name            = "GCP Resolver 1"
#       update_interval = 360
#       project_filter  = "string1"
#       instance_filter = "string2"
#     }

#   }

# }

