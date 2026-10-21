# Data sources look objects up by id or by exact name.
data "appgatesdp_site" "by_name" {
  name = "Default Site"
}

data "appgatesdp_site" "by_id" {
  id = "4c07bc67-57ea-42dd-b702-c2d6c45419fc"
}
