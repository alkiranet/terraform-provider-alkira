---
subcategory: "Release Notes"
page_title: "v1.5.0"
description: |-
    Release notes for v1.5.0
---

# Alkira Terraform Provider v1.5.0 Release Notes

Release Date: 2026-05-20

## Overview

Version 1.5.0 introduces new resources, F5 ILB support, and deprecates username/password authentication in favor of API keys.

---

## New Resources

- **`alkira_service_bluecat`**: Full lifecycle management of Bluecat DNS/DHCP services. Supports BDDS and Edge instance deployment with anycast IP configuration (including backup CXPs), segment-based assignment, DNS/DHCP service options, and automatic instance ordering.

- **`alkira_connector_azure_vnet_third_party`**: Connects to an Azure VNET through a third-party CXP peering gateway. Requires an existing peering gateway attachment, and supports static route prefix list configuration, segment assignment, and size-based scaling.

- **`alkira_peering_gateway_azure_vnet_third_party_connector_attachment`**: Manages the attachment between a CXP Peering Gateway and an Azure VNET for third-party connector use cases.

---

## Enhancements

- **Provider:** Username/password authentication is now deprecated. Use `api_key` instead. API keys can be managed from Portal > Settings > User Management. Note that API keys have a short expiry — if you use Terraform in CI/CD pipelines, consider automating API key rotation.
- **GCP VPC Connector (`alkira_connector_gcp_vpc`):** Added `export_all_subnets` field to `gcp_routing` schema. The provider now automatically derives this value from `vpc_subnet` presence.
- **Azure VNet Connector (`alkira_connector_azure_vnet`):** The `customer_asn` field is now `Optional + Computed`. In VGW mode, the backend auto-assigns an ASN when omitted, and the provider no longer reports drift for the backend-assigned value.
- **Juniper SD-WAN Connector (`alkira_connector_juniper_sdwan`):** Added `MaxItems: 1` constraint to the `instance` block.
- **F5 Load Balancer (`alkira_service_f5_lb`):** Added Internal Load Balancer (ILB) support with new fields `ilb_service_group_name`, `ilb_implicit_group_id`, `lb_type`, and `instance_metadata`.
- **F5 vServer Endpoint (`alkira_service_f5_vserver_endpoint`):** Added ILB type support. `fqdn_prefix` and `port_ranges` are now optional (required only for ELB type).
- **Network Entity Scale Options (`alkira_network_entity_scale_options`):** Added `additional_tunnel_options_per_node` nested block.
- **CXP Peering Gateway (`alkira_peering_gateway_cxp`):** Added computed `metadata` block. Data source now supports lookup by ID.
- **IP Reservation (`alkira_ip_reservation`):** The `prefix` field is now computed, supporting server-assigned values.
- **Policy Prefix List (`alkira_policy_prefix_list`):** Improved backward compatibility for the `prefixes` field. Use `prefix` and `prefix_range` blocks for new configurations.

---

## Bug Fixes

### Import Fixes

Multiple fixes related to `terraform import` to properly populate all fields and improve error handling for the following resources:

- `alkira_connector_gcp_vpc`
- `alkira_service_pan`
- `alkira_service_checkpoint`
- `alkira_service_cisco_ftdv`

Import now fails with a descriptive error on invalid resource IDs instead of silently producing incomplete state.

### State & Drift Fixes

- **GCP VPC Connector (`alkira_connector_gcp_vpc`):** Fixed `export_all_subnets` to prevent breaking diff on upgrade. Fixed `prefix_list_ids` reordering diffs so that order no longer matters, with automatic state migration. Fixed `userInputPrefixes` to send empty array instead of `null`. Fixed regression when removing all `vpc_subnet` blocks.
- **Azure VNet Connector (`alkira_connector_azure_vnet`):** Fixed perpetual diff on `customer_asn` in VGW mode when the backend auto-assigns the ASN value.
- **PAN Service (`alkira_service_pan`):** Fixed perpetual diff for `global_protect_segment_options`. Fixed `segment_options` perpetual diff caused by backend-injected ALKIRA_MGMT_ZONE.
- **Checkpoint (`alkira_service_checkpoint`):** Fixed `management_server` perpetual diff caused by sensitive fields not returned by the API. The block now preserves item order, with automatic state migration. Sensitive fields (password) are now carried from prior state.
- **Cisco FTDv (`alkira_service_cisco_ftdv`):** Fixed `firepower_management_center` perpetual diff caused by sensitive fields not returned by the API. The block now preserves item order, with automatic state migration. Sensitive fields (username, password) are now carried from prior state. Read now populates `segment_ids` from API response.
- **Policy Prefix List (`alkira_policy_prefix_list`):** Fixed deprecated `prefixes` field handling. Fixed spurious updates so that prefix order no longer matters, with automatic state migration. Fixed prefixes without descriptions not appearing in portal NAT-rule view.
- **Versa SD-WAN Connector (`alkira_connector_versa_sdwan`):** Fixed `version` field not being populated during Read.
- **Probe TCP (`alkira_probe_tcp`):** Fixed `network_entity_id` being set incorrectly.
- **Aruba Edge Connector (`alkira_connector_aruba_edge`):** Fixed instance IDs not being saved to state. Fixed `credentialId` being reset when credential fields were updated.
- **Peering Gateway TGW Attachment (`alkira_peering_gateway_aws_tgw_attachment`):** Fixed missing DXGW fields in Read. Fixed infinite loop and FAILED state handling during create.

---

## Documentation

- Added OUTBOUND NAT example to `policy_nat_rule` documentation.
- Updated F5 vServer Endpoint SNAT description.
- Fixed examples that referenced deprecated global CIDR lists.

---

## Upgrade Instructions

### From v1.4.4 to v1.5.0

1. **Authentication:**
   - If you use `username`/`password`, you will see a deprecation warning. Migrate to `api_key` at your convenience.

2. **State Migrations:**
   - `alkira_connector_gcp_vpc`: `prefix_list_ids` changed so that item order no longer matters.
   - `alkira_policy_prefix_list`: `prefix` and `prefix_range` changed so that item order no longer matters.
   - `alkira_service_checkpoint`: `management_server` changed to preserve item order and fix sensitive field handling.
   - `alkira_service_cisco_ftdv`: `firepower_management_center` changed to preserve item order and fix sensitive field handling.
   - All migrations are automatic. Run `terraform plan` after upgrading to verify no unexpected changes.

3. **GCP VPC Connector:**
   - Setting `export_all_subnets = false` without `vpc_subnet` blocks now produces a plan-time error. Either set `export_all_subnets = true` or add `vpc_subnet` blocks. This is not a breaking change — this configuration was already rejected by the API with a 400 error. The provider now catches it earlier at plan time.

4. **One-Time State Refresh:**
   - Many Read function fixes in this release mean that fields previously missing from state will now correctly reflect API values. After upgrading, run `terraform plan` and expect one-time diffs on the following resources:
     - `alkira_service_pan`
     - `alkira_peering_gateway_aws_tgw_attachment`
     - `alkira_connector_versa_sdwan`
     - `alkira_service_checkpoint`: `management_server.password` will appear as added
     - `alkira_service_cisco_ftdv`: `firepower_management_center.username` and `password` will appear as added
   - These diffs are benign. Run `terraform apply` once to stabilize state.

5. **No Breaking Changes.** This release is backward-compatible. All state migrations are automatic.
