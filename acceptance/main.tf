terraform {
  required_providers {
    alkira = {
      source = "alkiranet/alkira"
    }
  }
}

provider "alkira" {
  portal   = "terraform.preprod.alkira3.net"
}
