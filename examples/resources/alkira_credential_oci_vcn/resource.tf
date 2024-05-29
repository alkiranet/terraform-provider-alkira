resource "alkira_credential_oci_vcn" "example" {
  name        = "example"
  user_ocid   = "ocid1.user.oc1..axxxxx677oenm6cu2qcl46rhtq"
  fingerprint = "XX:XX:XX:XX:XX"
  key         = "PRIVATE_KEY"
  tenant_ocid = "xxxxxxxxxxxxxx" # Find this information from your OCI account
}



