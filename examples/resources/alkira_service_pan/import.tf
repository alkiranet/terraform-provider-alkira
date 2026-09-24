#
# A configuration sufficient to import an existing PAN service.
#
# Every argument below is Required by the schema, so all of them must be
# present before `terraform import` runs. The four credential arguments and
# the instance auth values are never returned by the API, so their values
# must come from you -- use the same ones the service was created with.
#
resource "alkira_service_pan" "example" {
  name                  = "pan-fw-1"
  cxp                   = "US-WEST"
  license_type          = "PAY_AS_YOU_GO"
  max_instance_count    = 1
  segment_ids           = [alkira_segment.example.id]
  management_segment_id = alkira_segment.example.id
  size                  = "SMALL"
  version               = "10.2.3"

  # Write-only on the API side. Import cannot read these back.
  pan_username           = var.pan_username
  pan_password           = var.pan_password
  registration_pin_id    = var.pan_registration_pin_id
  registration_pin_value = var.pan_registration_pin_value

  # Required block. `auth_key` and `auth_code` are also write-only.
  instance {
    name     = "pan-fw-1-instance-1"
    auth_key = var.pan_instance_auth_key
  }
}
