terraform {
  required_providers {
    alkira = {
      source = "alkiranet/alkira"
    }
  }
}

provider "alkira" {
}

locals {
  cxp = "US-WEST-1"
}
