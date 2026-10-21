# Turns an activated appliance into an additional controller and waits until
# it reports controller_ready. Destroying it disables the controller function.
resource "appgatesdp_appliance_controller_activation" "second" {
  appliance_id = appgatesdp_appliance.controller_2.id
  wait_timeout = "30m"

  controller = {
    enabled = true
  }

  admin_interface = {
    hostname      = "controller-2.example.com"
    https_ciphers = ["ECDHE-RSA-AES256-GCM-SHA384", "ECDHE-RSA-AES128-GCM-SHA256"]
  }
}
