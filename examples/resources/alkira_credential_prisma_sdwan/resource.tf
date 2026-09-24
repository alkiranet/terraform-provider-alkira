resource "alkira_credential_prisma_sdwan" "example" {
  name       = "prisma-sdwan-cred-1"
  ion_token  = var.prisma_ion_token
  ion_secret = var.prisma_ion_secret
}
