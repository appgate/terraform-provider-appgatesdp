# Uploads a license file. Changing the license replaces it (the previous
# license is deleted).
resource "appgatesdp_license" "this" {
  license = file("appgate.license")
}
