---
subcategory: "Release Notes"
page_title: "v1.3.0"
description: |-
    Release notes for v1.3.0
---

This major release contains several new resources and enhancements
across the board.

## Resources & Data Sources

#### resource `alkira_connector_aruba_edge`

* Deprecate field `segment_id`.
* Deprecate field `gateway_bgp_asn`.

Both fields should be specified inside `vrf_mapping` block now.

#### resource `alkira_connector_aws_dx` (**NEW**)

New resource for supporting AWS DirectConnect.

#### resource `alkira_connector_azure_expressroute`

* Add optional field `description`.

#### resource `alkira_connector_gcp_interconnect` (**NEW**)

New resource for supporting GCP InerConnect.

#### resource `alkira_connector_gcp_vpc`

* Update supported `size`.

#### resource `alkira_connector_ipsec`

* Fix documentation of various fields.
* Update examples.

#### resource `alkira_connector_ipsec_adv`

* Fix documentation of various fields.
* Update examples.

#### resource `alkira_connector_remote_access`

* New optional field `fallback_to_tcp` to support TCP fallback.

#### resource `alkira_policy_prefix_list`

* Fix the reading problem of `prefix_range` block during update and
  import.

#### resource `alkira_service_f5` (**NEW**)(**BETA**)

Resource for support F5 load balancer. This is the first version with
basic functionality support.

#### resource `alkira_service_zscaler`

* Fix the crash when payload is truncated or imcomplete.

