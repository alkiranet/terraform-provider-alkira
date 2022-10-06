resource "alkira_credential_checkpoint" "tf_test_checkpoint" {
  name                       = "tf-test-checkpoint"
  password                   = "Ak12345678"
  management_server_password = "MGMTPSWD111"
  sic_keys                   = ["AAAAA88888888888", "BBBBB999999999"]
}
