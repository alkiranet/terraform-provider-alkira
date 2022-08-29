resource "alkira_segment" "test" {
  name  = "testoci"
  asn   = "65513"
  cidrs = ["10.16.1.0/24"]
}

resource "alkira_group" "test" {
  name        = "testoci"
  description = "test group"
}

resource "alkira_connector_oci_vcn" "test" {
  name           = "test"
  oci_region     = "us-sanjose-1"
  cxp            = "US-WEST"
  enabled        = true
  primary        = true
  vcn_id         = "ocid1.vcn.oc1.us-sanjose-1.XXXX" # fill in proper vcn_id
  vcn_cidr       = ["172.24.0.0/16"]
  credential_id  = alkira_credential_oci_vcn.test.id # using credential_oci_vcn
  group          = alkira_group.test.name
  segment_id     = alkira_segment.test.id
  size           = "SMALL"
}
