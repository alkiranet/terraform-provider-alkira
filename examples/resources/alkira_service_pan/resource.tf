resource "alkira_service_pan" "test1" {
  name                   = "test1-update"
  credential_id          = alkira_credential_pan.tf_test_pan.id
  cxp                    = "US-WEST"
  global_protect_enabled = "true"
  license_type           = "PAY_AS_YOU_GO"
  panorama_enabled       = false
  panorama_device_group  = "alkira-test"
  panorama_ip_addresses  = ["172.16.0.8"]
  panorama_template      = "test"
  max_instance_count     = 1
  segment_ids            = [alkira_segment.test1.id, alkira_segment.test2.id]
  management_segment_id  = alkira_segment.test1.id
  size                   = "SMALL"
  type                   = "VM-300"
  version                = "9.0.5-xfr"

  global_protect_segment_options {
    segment_id            = (alkira_segment.test1.id)
    remote_user_zone_name = "doesn't matter"
    portal_fqdn_prefix    = "also doesn't matter"
    service_group_name    = "still doesn't matter"
  }

  instance {
    name          = "tf-pan-instance-1"
    credential_id = alkira_credential_pan_instance.tf_test_pan_instance.id
    global_protect_segment_options {
      segment_id      = (alkira_segment.test1.id)
      portal_enabled  = true
      gateway_enabled = true
      prefix_list_id  = 548
    }
  }

  zones_to_groups {
    segment_id = alkira_segment.test1.id
    zone_name  = "Zone11"
    groups     = [alkira_group.test.name]
  }
}
