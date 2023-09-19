# resource "alkira_connector_versa_sdwan" "test1" {
#   name  = "tf-test-1"
#   cxp   = var.cxp
#   group = alkira_group.test1.name
#   size  = "SMALL"

#   versa_controller_host = "172.16.0.1"
#   local_id              = 1

#   remote_id                = 2
#   remote_public_shared_key = "123456ABCD"

#   versa_vos_device {
#     hostname                   = "dev1"
#     local_device_serial_number = "12345678"
#     version                    = "21.2.3-B"
#   }

#   vrf_segment_mapping {
#     segment_id    = alkira_segment.test1.id
#     vrf_name      = "test"
#     versa_bgp_asn = 1203403435
#   }
# }

# resource "alkira_connector_versa_sdwan" "test2" {
#   name  = "tf-test-2"
#   cxp   = var.cxp
#   group = alkira_group.test1.name
#   size  = "SMALL"

#   versa_controller_host = "172.16.0.1"
#   local_id              = 1
#   remote_id             = 2

#   versa_vos_device {
#     hostname                   = "dev1"
#     local_device_serial_number = "12345678"
#     version                    = "21.2.3-B"
#   }

#   versa_vos_device {
#     hostname                   = "dev2"
#     local_device_serial_number = "123456789"
#     version                    = "21.2.3-B"
#   }

#   vrf_segment_mapping {
#     segment_id    = alkira_segment.test1.id
#     vrf_name      = "test"
#     versa_bgp_asn = 1203403435
#   }
# }
