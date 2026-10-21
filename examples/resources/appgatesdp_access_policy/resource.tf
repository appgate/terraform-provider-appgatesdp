resource "appgatesdp_access_policy" "developers" {
  name         = "developers"
  expression   = "return claims.user.groups.indexOf('developers') >= 0;"
  entitlements = [appgatesdp_entitlement.git.id]
}
