resource "alkira_policy_prefix_list" "ranges" {
  name        = "range-based-prefixes"
  description = "Prefix list with CIDR ranges"

  prefix_range {
    prefix      = "10.1.0.0/16"
    le          = 20
    ge          = 18
    description = "Flexible range for branch offices"
  }

}

