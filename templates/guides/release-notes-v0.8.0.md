---
subcategory: "Release Notes"
page_title: "v0.8.0"
description: |-
    Release notes for v0.8.0
---

Major release with multiple new resources and enhancements across
various resources.

### resource `alkira_connector_oci_vcn` (**NEW**)

New resource is added to support managing OCI VCN connector through
Terraform.

### resource `credential_oci_vcn` (**NEW**)

New credential resource is added to support managing OCI VCN
credential.

### resource `alkira_cloudvisor_account` (**NEW**)

New resource to CloudVisor account.

### resource `alkira_policy_nat` (**NEW**)

New resource to manage NAT policy.

### resource `alkira_policy_nat_rule` (**NEW**)

New resource to manage NAT policy rule to work with NAT policy.

### resource `alkira_list_as_path` (**NEW**)

New resource to manage AS Path list.

### resource `alkira_list_community` (**NEW**)

New resource to manage community list.

### resource `alkira_list_extended_community` (**NEW**)

New resource to manage extended community list.

### resource `alkira_list_global_cidr` (**NEW**)

New resource to manage extended global CIDR list.

### resource `alkira_connector_cisco_sdwan`

* Add new identifier `type` to specify the type of Cisco SDWAN.
* `vedge` is required now.
* `vrf_segment_mapping` is required now.
* `hostname` and `cloud_init_file` are required now.

### resource `alkira_service_pan`

* Add new identifier `bundle`.
