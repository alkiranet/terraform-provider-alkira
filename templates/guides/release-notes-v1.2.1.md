---
subcategory: "Release Notes"
page_title: "v1.2.1"
description: |-
    Release notes for v1.2.1
---

A quick patch release to address use cases around `implicit_group`.


The `implicit_group` will be automatically created along with the most
connector resources and it could be used directly in `alkira_policy`.

The following connectors already support `implicit_group`:

* alkira_connector_azure_expressroute
* alkira_connector_akamai_prolexic
* alkira_connector_gcp_vpc
* alkira_connector_ipsec
* alkira_connector_azure_vnet
* alkira_connector_aws_vpc
* alkira_connector_ipsec
* alkira_connector_ipsec_adv
* alkira_connector_aruba_edge
* alkira_connector_cisco_sdwan
* alkira_connector_internet_exit
* alkira_connector_oci_vcn

Once the above connectors are created, an `implicit_group_id` computed
field will be provided from the resource.


## DATA SOURCES

Data sources for above connectors have been updated as well to show
the field.


### `alkira_ip_reservation`

The data source was added to allow the existing IP reservation to be
referred directly.


