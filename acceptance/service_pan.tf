resource "alkira_service_pan" "test" {
  name         = "acceptance"
  cxp          = var.cxp
  license_type = "PAY_AS_YOU_GO"

  pan_username = "admin"
  pan_password = "Test12345678"

  panorama_enabled      = true
  panorama_device_group = "test"
  panorama_ip_addresses = ["172.16.0.8"]
  panorama_template     = "test"
  max_instance_count    = 1
  segment_ids           = [alkira_segment.test1.id, alkira_segment.test2.id]
  management_segment_id = alkira_segment.test1.id
  size                  = "SMALL"
  bundle                = "VM_SERIES_BUNDLE_2"
  version               = "10.2.2-h2"

  registration_pin_id     = "1234567890ABCDEF"
  registration_pin_value  = "1234567890ABCDEF"
  registration_pin_expiry = "2025-12-31"

  # master_key_enabled = true
  # master_key         = "1234567890ABCDEF"
  # master_key_expiry  = "2022-09-30"

  instance {
    name     = "tf-pan-instance-1"
    auth_key = "test"
  }

  # zones_to_groups {
  #   segment_id   = alkira_segment.seg1.id
  #   zone_name    = "Zone222"
  #   groups       = [alkira_group.test.name]
  # }
}
