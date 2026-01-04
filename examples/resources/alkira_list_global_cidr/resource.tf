resource "alkira_list_global_cidr" "test" {
  name        = "test"
  description = "test global cidr list"
  values      = ["172.16.1.0/24", "10.1.0.0/24"]
  cxp         = "US-WEST"
  tags        = ["INFOBLOX", "CHKPFW", "CISCO_FTDV_FW"]
}
