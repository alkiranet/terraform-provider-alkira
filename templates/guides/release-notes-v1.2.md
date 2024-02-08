---
subcategory: "Release Notes"
page_title: "v1.2"
description: |-
    Release notes for v1.2
---

This release contains several new resources and enhancements across
the board. Some common improvements:

* Performance improvements with all Terraform commands.
* Better debug-ability.
* New architecture support of `LINUX/ARM64`.

## Generating configuration (**BETA**)

In this release, we focused on supporting the new [Terraform Config
generation](https://developer.hashicorp.com/terraform/language/import/generating-configuration). Most
resources should have tested with the new feature.

## Resources & Data Sources

#### resource `alkira_connector_ipsec_tunnel_profile` (**NEW**)

This new resource is required when specifying advanced options in
`alkira_connector_ipsec_adv`.

#### resource `alkira_connector_ipsec_adv`

After the initial introduction in `v1.1.1`, this resource has been
improved a lot in this release.

#### resource `alkira_connector_ipsec`

Several bug fixes and optimizations.

* Add optional `scale_group_id`.

#### resource `alkira_connector_remote_access` (**NEW**)

New resource to create Alkira Remote Access Connector.

#### resource `alkira_ip_reservation` (**NEW**)

New resource to create IP reservation.

#### resource `alkira_policy_nat_rule`

* Bug fixes and improvements.

#### resource `alkira_service_fortinet`

* Fix the problem that resource may generate a diff after the first
  apply.




