resource "alkira_list_as_path" "tf_test" {
  name        = "tf-test"
  description = "terraform test as path list"
  values      = ["100 [2-5]00", "_6400_"]
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

resource "alkira_list_global_cidr" "tf_test" {
  name        = "tf-test"
  description = "terraform test global cidr list"
  values      = ["172.16.1.0/24", "10.1.0.0/24"]
  cxp         = "US-WEST-1"
}
