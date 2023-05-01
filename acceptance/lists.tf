resource "alkira_policy_prefix_list" "test1" {
  name        = "acceptance-test-prefix-list"
  description = "terraform test policy prefix list"
  prefixes    = ["0.0.0.0/0"]
}

resource "alkira_policy_prefix_list" "test2" {
  name        = "acceptance-test-prefix-list-ranges"
  description = "terraform test policy prefix list"
  prefixes    = ["0.0.0.0/0"]

  prefix_range {
    prefix = "0.0.0.0/0"
    le     = 4
    ge     = 2
  }
}

resource "alkira_list_community" "test" {
  name        = "acceptance-test"
  description = "terraform test community list"
  values      = ["65512:20", "65512:21"]
}

resource "alkira_list_extended_community" "test" {
  name        = "acceptance-test"
  description = "terraform test extended community list"
  values      = ["soo:65512:20", "soo:65512:21"]
}

resource "alkira_list_as_path" "test" {
  name        = "acceptance-test"
  description = "terraform test as path list"
  values      = ["100 [2-5]00", "_6400_"]
}

resource "alkira_list_global_cidr" "ciscofdtv" {
  name        = "acceptance-test"
  description = "terraform test global cidr list for cisco ftdv"
  values      = ["10.0.0.0/25"]
  cxp         = "US-WEST-1"
  tags        = ["CISCO_FTDV_FW"]
}

resource "alkira_list_global_cidr" "checkpoint" {
  name        = "acceptance-checkpoint"
  description = "global cidr list for checkpoint"
  values      = ["172.16.1.0/24"]
  cxp         = "US-WEST-1"
}

data "alkira_policy_prefix_list" "prefix1" {
  name = "acceptance-data-prefix1"
}
