resource "alkira_list_policy_fqdn" "test" {
  name               = "test"
  description        = "test policy fqdn list"
  fqdns              = ["test.alkira.com"]
  list_dns_server_id = alkira_list_dns_server.test.id
}
