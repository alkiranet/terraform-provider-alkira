resource "alkira_list_udr" "test1" {
  name               = "tf-test-1"
  description        = "terraform test UDR list 1"
  cloud_provider     = "AZURE"

  route {
    prefix = "10.0.0.0/24"
    description = "test route 1"
  }
}
