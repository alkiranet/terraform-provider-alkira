variable "meraki_vmx_1_auth_token" {
  type      = string
  sensitive = true
}

variable "meraki_vmx_2_auth_token" {
  type      = string
  sensitive = true
}

resource "alkira_connector_meraki_sdwan" "test1" {
  cxp   = "AZURE-WESTUS2"
  group = alkira_group.test.name
  name  = "meraki-sdwan-connector1"
  size  = "SMALL"

  target_segment {
    segment_id               = alkira_segment.test1.id
    meraki_bgp_asn           = 65500
    advertise_on_prem_routes = true
    advertise_default_route  = true
  }

  instances {
    host_name  = "meraki-vmx-1"
    auth_token = var.meraki_vmx_1_auth_token
  }

  instances {
    host_name  = "meraki-vmx-2"
    auth_token = var.meraki_vmx_2_auth_token
  }
}
