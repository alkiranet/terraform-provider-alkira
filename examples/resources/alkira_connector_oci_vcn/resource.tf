resource "alkira_connector_oci_vcn" "test" {
  name           = "test"
  oci_region     = "us-sanjose-1"
  cxp            = "US-WEST"
  enabled        = true
  primary        = true
  vcn_id         = "ocid1.vcn.oc1.us-sanjose-1.XXXX"
  vcn_cidr       = ["172.24.0.0/16"]
  credential_id  = alkira_credential_oci_vcn.test.id
  group          = alkira_group.test.name
  segment_id     = alkira_segment.test.id
  size           = "SMALL"
}
