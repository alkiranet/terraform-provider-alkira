# ZTA Profile Data Source Example
data "alkira_zta_profile" "example" {
  name = "corporate-zta-profile"
}

# Using the ZTA profile ID in a resource
resource "alkira_connector_remote_access" "zta_connector" {
  name           = "zta-remote-access"
  description    = "Remote access connector with ZTA profile"
  cxp            = "US-WEST"
  group          = alkira_group.remote_access.name
  segment_id     = alkira_segment.remote_access.id
  size           = "SMALL"

  # Reference the ZTA profile from the data source
  zta_profile_id = data.alkira_zta_profile.example.id

  endpoint {
    name                = "zta-endpoint"
    customer_gateway_ip = "203.0.113.50"
    preshared_keys      = ["zta-preshared-key"]
  }
}
