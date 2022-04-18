resource "alkira_credential_checkpoint" "tf_test_checkpoint" {
  name     = "tf-test-checkpoint"
  password = "Ak12345678"
  sic_keys = ["88888888888", "999999999"]
}
