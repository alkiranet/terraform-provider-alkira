resource "alkira_credential_gcp_vpc" "gcp" {
  name                 = "tftest"
  auth_provider        = "https://www.googleapis.com/oauth2/v1/certs"
  auth_uri             = "https://accounts.google.com/o/oauth2/auth"

  client_email         = "xxx@xxx.gcp.iam.gserviceaccount.com"
  client_id            = "123455678"
  client_x509_cert_url = "https://www.googleapis.com/robot/v1/metadata/x509/xxxx.iam.gserviceaccount.com"

  private_key          = "PEM KEY"
  private_key_id       = "PEM KEY ID"

  project_id           = "test"
  token_uri            = "https://oauth2.googleapis.com/token"
  type                 = "service_account"
}


