resource "alkira_billing_tag" "test1" {
  name = "acceptance-test1"
}

resource "alkira_billing_tag" "test2" {
  name = "acceptance-test2"
}

resource "alkira_billing_tag" "test3" {
  name = "acceptance-test3"
}

data "alkira_billing_tag" "tag1" {
  name = "acceptance-data-tag1"
}
