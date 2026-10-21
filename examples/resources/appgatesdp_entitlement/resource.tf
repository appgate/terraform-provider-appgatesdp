data "appgatesdp_site" "default" {
  name = "Default Site"
}

resource "appgatesdp_entitlement" "git" {
  name       = "git"
  site       = data.appgatesdp_site.default.id
  conditions = []

  actions = [{
    subtype = "tcp_up"
    action  = "allow"
    hosts   = ["git.internal.example.com"]
    ports   = ["443", "22"]
  }]
}
