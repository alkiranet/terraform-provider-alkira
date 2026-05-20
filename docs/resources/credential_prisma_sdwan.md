# alkira_credential_prisma_sdwan

Provides Prisma SD-WAN connector instance credential resource.

## Example Usage

```hcl
resource "alkira_credential_prisma_sdwan" "example" {
  name       = "prisma-sdwan-cred-1"
  ion_token  = var.prisma_ion_token
  ion_secret = var.prisma_ion_secret
}
```

## Argument Reference

* `name` - (Required) The name of the credential.
* `ion_token` - (Required, Sensitive) The ION token for Prisma SD-WAN device registration.
* `ion_secret` - (Required, Sensitive) The ION secret for Prisma SD-WAN device registration.

## Attribute Reference

* `id` - The ID of the credential.

## Import

```shell
terraform import alkira_credential_prisma_sdwan.example credential-id
```
