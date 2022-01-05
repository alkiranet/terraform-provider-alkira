resource "alkira_list_extended_community" "test" {
  name        = "test"
  description = "test extended community list"
  values      = ["soo:65512:20", "soo:65512:21"]
}
