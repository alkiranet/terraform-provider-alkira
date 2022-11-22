# Create list resources for each under Management/Lists/

resource "alkira_policy_prefix_list" "tf_prefix_list" {
  name        = "tf-test-prefix-list"
  description = "terraform test policy prefix list"
  prefixes    = ["0.0.0.0/0"]
}

resource "alkira_list_community" "tf_test" {
  name        = "tf-test"
  description = "terraform test community list"
  values      = ["65512:20", "65512:21"]
}

resource "alkira_list_extended_community" "tf_test" {
  name        = "tf-test"
  description = "terraform test extended community list"
  values      = ["soo:65512:20", "soo:65512:21"]
}

resource "alkira_list_as_path" "tf_test" {
  name        = "tf-test"
  description = "terraform test as path list"
  values      = ["100 [2-5]00", "_6400_"]
}

resource "alkira_list_global_cidr" "tf_test" {
  name        = "tf-test"
  description = "terraform test global cidr list for cisco ftdv"
  values      = ["10.0.0.0/25"]
  cxp         = "US-WEST-1"
  tags        = ["CISCO_FTDV_FW"]
}
