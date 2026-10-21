resource "appgatesdp_condition" "business_hours" {
  name       = "business-hours"
  expression = "return new Date().getHours() >= 8 && new Date().getHours() < 18;"
  tags       = ["terraform"]

  repeat_schedules = ["1h", "13:00"]

  remedy_methods = [{
    type    = "DisplayMessage"
    message = "Access is only allowed during business hours."
  }]
}
