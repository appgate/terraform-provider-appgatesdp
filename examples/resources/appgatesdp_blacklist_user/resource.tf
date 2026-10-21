resource "appgatesdp_blacklist_user" "alice" {
  user_distinguished_name = "CN=alice,OU=ldap"
  reason                  = "Device reported lost"
}
