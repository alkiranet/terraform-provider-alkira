#
# Create segments assuming this is the completely empty tenant network
#
resource "alkira_segment" "segment1" {
  name = "seg1"
  asn  = "65513"
  cidr = "10.16.1.0/24"
}


#
# Create the credential to store the access to the AWS account that
# VPCs belongs two. In this example, both VPCs belong to this AWS
# account.
#
resource "alkira_credential_azure_vnet" "customer_azure" {
  name            = "customer-azure"
  application_id  = ""
  secret_key      = ""
  subscription_id = ""
  tenant_id       = ""
}


#
# Create AZURE-VNET connector for the first VNET and attach it with
# segment 1
#
resource "alkira_connector_azure_vnet" "connector_vnet1" {
  name           = "customer-vnet1"
  azure_region   = "westus2"
  azure_vnet_id  = "/subscriptions/XXXX/resourceGroups/Test/providers/Microsoft.Network/virtualNetworks/customer-vnet1"
  credential_id  = alkira_credential_azure_vnet.customer_azure.id
  cxp            = "US-WEST"
  group          = "test"
  segment        = alkira_segment.segment1.name
  size           = "SMALL"
}
