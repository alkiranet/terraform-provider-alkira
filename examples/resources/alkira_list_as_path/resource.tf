resource "alkira_list_as_path" "test" {
  name        = "test"
  description = "test as path list"
  values      = ["100 [2-5]00", "_6400_"]
}
