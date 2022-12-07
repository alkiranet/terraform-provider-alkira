resource "alkira_policy_prefix_list" "test" {
  name        = "acceptance-test"
  description = "terraform test policy prefix list"
  prefixes    = ["0.0.0.0/0"]
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

resource "alkira_list_global_cidr" "test" {
  name        = "acceptance-test"
  description = "terraform test global cidr list for cisco ftdv"
  values      = ["10.0.0.0/25"]
  cxp         = "US-WEST-1"
  tags        = ["CISCO_FTDV_FW"]
}
