#
# PAN FW service managed via Strata Cloud Manager (SCM) with PAN-OS
# Advanced Routing (AR).
#
# Required when scm_enabled = true:
#   - license_type    = "BRING_YOUR_OWN"   (PAYG not supported with SCM)
#   - routing_type    = "advanced"
#   - version         >= "10.2.3"
#   - registration_pin_id / registration_pin_value
#
# Mutually exclusive: scm_enabled and panorama_enabled cannot both be true.
# panorama_template / panorama_device_group / panorama_ip_addresses are
# NOT used in SCM mode.
#
resource "alkira_service_pan" "scm_example" {
  name                  = "pan-scm-example"
  cxp                   = "US-WEST"
  size                  = "SMALL"
  type                  = "VM-300"
  version               = "11.1.6"
  license_type          = "BRING_YOUR_OWN"
  pan_license_key       = "EXAMPLE"
  max_instance_count    = 1
  segment_ids           = [alkira_segment.test1.id]
  management_segment_id = alkira_segment.test1.id

  pan_password = "EXAMPLE"
  pan_username = "admin"

  registration_pin_id    = "EXAMPLE"
  registration_pin_value = "EXAMPLE"

  # SCM + Advanced Routing
  scm_enabled  = true
  scm_folder   = "prod/edge/site-a"
  routing_type = "advanced"

  instance {
    name      = "pan-scm-instance-1"
    auth_code = "tenant-pan-auth-code"
  }
}
