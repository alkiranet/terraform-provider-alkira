---
subcategory: "Release Notes"
page_title: "v1.1.1"
description: |-
    Release notes for v1.1.1
---

This patch release introduces several new resources (in BETA state)
and bug fixes.

NOTE: Some new features may not be enabled for all tenants. For
support and detailed usage of new features, please contact Alkira
Support Team.

## Resources & Data Sources

#### resource `alkira_connector_ipsec_adv` (**NEW**)

New resource to create IPSec connector with more options.

#### resource `alkira_connector_ipsec`

Several bug fixes and optimizations.

* Fix default value of `enable_tunnel_redundancy`.

#### resource `alkira_connector_versa_sdwan` (**NEW**)

#### resource `alkira_policy_nat`

Optimize the usage of address translation block.

* Set the default value of `src_addr_translation_match_and_invalidate` to `true`.
* Set the default value of `dst_addr_translation_advertise_to_connector` to `true`.
* Remove `src_addr_translation_bidirectional` that can't be changed by user.
* Remove `dst_addr_translation_bidirectional` that can't be changed by user.
* Add Optional new `allow_overlapping_translated_source_addresses`.


