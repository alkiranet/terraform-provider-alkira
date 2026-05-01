---
subcategory: "Release Notes"
page_title: "v1.5.0"
description: |-
    Release notes for v1.5.0
---

# Alkira Terraform Provider v1.5.0 Release Notes

Release Date: 2026-05-01

## Overview

Version 1.5.0 introduces new resources, F5 ILB support, and deprecates username/password authentication in favor of API keys.

---

## New Resources

- **`alkira_service_bluecat`**: Full lifecycle management of Bluecat DNS/DHCP services. Supports BDDS and Edge instance deployment with anycast IP configuration (including backup CXPs), segment-based assignment, DNS/DHCP service options, and automatic instance ordering.

- **`alkira_connector_azure_vnet_third_party`**: Enables connectivity to Azure VNETs through third-party network appliances. Supports static route prefix list configuration, segment assignment, and size-based scaling.

- **`alkira_peering_gateway_azure_vnet_third_party_connector_attachment`**: Manages the attachment between a CXP Peering Gateway and an Azure VNET for third-party connector use cases.

Corresponding data sources are also available for lookup by name or ID.

---

## Enhancements

- **Provider:** Username/password authentication is now deprecated. Use `api_key` instead. API keys can be managed from Portal > Settings > User Management.
- **Aruba Edge Connector:** Added `scale_group_id` for scalegroup support.
- **GCP VPC Connector:** Added `export_all_subnets` field to `gcp_routing` schema. The provider now automatically derives this value from `vpc_subnet` presence.
- **AWS VPC Connector:** When `vpc_subnet` entries are specified, `exportAllSubnets` is now forced to `false` to prevent conflicting configuration.
- **Juniper SD-WAN Connector:** Added `MaxItems: 1` constraint to the `instance` block.
- **F5 Load Balancer (`alkira_service_f5_lb`):** Added Internal Load Balancer (ILB) support with new fields `ilb_service_group_name`, `ilb_implicit_group_id`, `lb_type`, and `instance_metadata`.
- **F5 vServer Endpoint:** Added ILB type support. `fqdn_prefix` and `port_ranges` are now optional (required only for ELB type).
- **Network Entity Scale Options:** Added `additional_tunnel_options_per_node` nested block.
- **CXP Peering Gateway:** Added computed `metadata` block. Data source now supports lookup by ID.
- **IP Reservation:** The `prefix` field is now computed, supporting server-assigned values.
- **Policy Prefix List:** The `prefixes` field is deprecated. Use `prefix` and `prefix_range` blocks instead.

---

## Bug Fixes

### Import Fixes

Multiple fixes related to `terraform import` to properly populate all fields and improve error handling for the following resources:

- `alkira_connector_gcp_vpc`
- `alkira_connector_ipsec`
- `alkira_service_pan`
- `alkira_segment_resource_share`
- `alkira_policy_prefix_list`
- `alkira_service_checkpoint`
- `alkira_service_cisco_ftdv`

Import now fails with a descriptive error on invalid resource IDs instead of silently producing incomplete state.

### State & Drift Fixes

- **GCP VPC Connector:** Fixed `export_all_subnets` to prevent breaking diff on upgrade. Fixed `prefix_list_ids` reordering diffs by switching to `TypeSet` with automatic state migration. Fixed `userInputPrefixes` to send empty array instead of `null`. Fixed regression when removing all `vpc_subnet` blocks.
- **AWS VPC Connector:** Fixed drift from backend-defaulted routing fields.
- **PAN Service:** Fixed perpetual diff for `global_protect_segment_options`.
- **Policy Prefix List:** Fixed deprecated `prefixes` field handling. Fixed spurious updates by switching to `TypeSet` with automatic state migration.
- **IPsec Connector:** Fixed `segment_options` flatten to handle `interface{}` type correctly.
- **Cisco FTDv:** Read now populates `segment_ids` from API response.
- **Versa SD-WAN Connector:** Fixed `version` field not being populated during Read.
- **Probe TCP:** Fixed `network_entity_id` being set incorrectly.
- **Aruba Edge Connector:** Fixed instance IDs not being saved to state. Fixed `credentialId` being reset when credential fields were updated.
- **Peering Gateway TGW Attachment:** Fixed missing DXGW fields in Read. Fixed infinite loop and FAILED state handling during create.
- **NAT Policy:** Fixed schema definition.

---

## Documentation

- Rewrote provider authentication documentation to recommend API key as the primary method.
- Added OUTBOUND NAT example to `policy_nat_rule` documentation.
- Updated F5 vServer Endpoint SNAT description.
- Fixed examples that referenced deprecated global CIDR lists.

---

## Upgrade Instructions

### From v1.4.4 to v1.5.0

1. **Authentication:**
   - If you use `username`/`password`, you will see a deprecation warning. Migrate to `api_key` at your convenience.

2. **State Migrations:**
   - `alkira_connector_gcp_vpc`: `prefix_list_ids` migrated from `TypeList` to `TypeSet`.
   - `alkira_policy_prefix_list`: `prefix` and `prefix_range` migrated from `TypeList` to `TypeSet`.
   - Run `terraform plan` after upgrading to verify no unexpected changes.

3. **GCP VPC Connector:**
   - Setting `export_all_subnets = false` without `vpc_subnet` blocks now produces a plan-time error. Either set `export_all_subnets = true` or add `vpc_subnet` blocks.

4. **One-Time State Refresh:**
   - Many Read function fixes in this release mean that fields previously missing from state will now correctly reflect API values. After upgrading, run `terraform plan` and expect one-time diffs on the following resources:
     - `alkira_service_pan`
     - `alkira_peering_gateway_aws_tgw_attachment`
     - `alkira_connector_versa_sdwan`
     - `alkira_service_checkpoint`
     - `alkira_service_cisco_ftdv`
   - These diffs are benign. Run `terraform apply` once to stabilize state.

5. **No Breaking Changes.** This release is backward-compatible. All state migrations are automatic.
