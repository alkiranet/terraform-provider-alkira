resource "alkira_list_community" "test" {
  name        = "test"
  description = "test community list"
  values      = ["65512:20", "65512:21"]
}
