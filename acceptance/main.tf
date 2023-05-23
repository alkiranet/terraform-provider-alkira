terraform {
  required_providers {
    alkira = {
      source = "alkiranet/alkira"
    }
  }
}

provider "alkira" {
  portal = "terraform.preprod.alkira3.net"
}

locals {
  cxp = "US-WEST-1"
}
