resource "alkira_service_pan" "test1" {
  name                    = "PanFwTest"
  pan_password            = "Ak12345678"
  pan_username            = "admin"
  bundle                  = "PAN_VM_300_BUNDLE_2"
  cxp                     = "US-WEST"
  global_protect_enabled  = false
  license_type            = "PAY_AS_YOU_GO"
  max_instance_count      = 1
  segment_ids             = [alkira_segment.test1.id]
  management_segment_id   = alkira_segment.test1.id
  size                    = "SMALL"
  type                    = "VM-300"
  version                 = "9.1.3"

  panorama_enabled        = true
  panorama_device_group   = "alkira-test"
  panorama_ip_addresses   = ["172.16.0.8"]
  panorama_template       = "test"

  registration_pin_id     = "1234567890ABCDEF"
  registration_pin_value  = "1234567890ABCDEF"
  registration_pin_expiry = "2023-07-30"

  master_key_enabled      = true
  master_key              = "1234567890ABCDEF"
  master_key_expiry       = "2023-08-01"

  global_protect_segment_options {
    segment_id            = (alkira_segment.test1.id)
    remote_user_zone_name = "RandomZoneName"
    portal_fqdn_prefix    = "randomprefix"
    service_group_name    = "RandomServiceGroupName"
  }

  # You can add more instance blocks. Make sure to set max_instance_count.
  instance {
    name      = "tf-pan-instance-1"
    auth_key  = "tenant-pan-auth-code" 
    auth_code = "tenant-pan-auth-code"
    global_protect_segment_options {
      segment_id      = (alkira_segment.test1.id)
      portal_enabled  = true
      gateway_enabled = true
      prefix_list_id  = alkira_policy_prefix_list.tf_prefix_list.id
    }
  }

  zones_to_groups {
    segment_id  = alkira_segment.test1.id
    zone_name   = "Zone1"
    groups      = [alkira_group.test.name]
  }
}
