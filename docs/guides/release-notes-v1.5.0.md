---
subcategory: "Release Notes"
page_title: "v1.5.0"
description: |-
    Release notes for v1.5.0
---

# Alkira Terraform Provider v1.5.0 Release Notes

Release Date: 2026-05-01

## Overview

Version 1.5.0 introduces three new resources, F5 Internal Load Balancer and Aruba Edge scalegroup support, deprecates username/password authentication in favor of API keys, and includes 30+ bug fixes across import, state drift, and connector reliability.

---

## New Resources

- **`alkira_service_bluecat`**: Full lifecycle management of Bluecat DNS/DHCP services. Supports BDDS and Edge instance deployment with anycast IP configuration (including backup CXPs), segment-based assignment, DNS/DHCP service options, and automatic instance ordering. Instances can be managed with hostname-based identification, and the resource handles license credential provisioning automatically.

- **`alkira_connector_azure_vnet_third_party`**: Enables connectivity to Azure VNETs through third-party network appliances. Supports static route prefix list configuration, segment assignment, and size-based scaling. Works in conjunction with the CXP Peering Gateway for third-party VNET integration.

- **`alkira_peering_gateway_azure_vnet_third_party_connector_attachment`**: Manages the attachment between a CXP Peering Gateway and an Azure VNET for third-party connector use cases. Supports linking Azure VNET resources to peering gateways with state tracking.

Corresponding data sources are also available for lookup by name or ID.

---

## Enhancements

### Authentication

- **Provider:** Username/password authentication is now deprecated. Use `api_key` instead. Existing username/password configurations will continue to work but will produce a deprecation warning. API keys can be managed from Portal > Settings > User Management.

### Connectors

- **Aruba Edge Connector:** Added `scale_group_id` for scalegroup support.
- **GCP VPC Connector:** Added `export_all_subnets` field to `gcp_routing` schema. The provider now automatically derives this value from `vpc_subnet` presence, preventing conflicting configurations.
- **AWS VPC Connector:** When `vpc_subnet` entries are specified, `exportAllSubnets` is now forced to `false` to prevent conflicting configuration.
- **Juniper SD-WAN Connector:** Added `MaxItems: 1` constraint to the `instance` block. Only one instance per connector is supported.

### Services

- **F5 Load Balancer (`alkira_service_f5_lb`):** Added Internal Load Balancer (ILB) support. New fields: `ilb_service_group_name`, `ilb_implicit_group_id`, and `lb_type` in `segment_options`. Computed `instance_metadata` now includes tunnel and routing details.
- **F5 vServer Endpoint (`alkira_service_f5_vserver_endpoint`):** Added ILB type support. New fields: `destination_endpoint_ip_addresses`, `destination_endpoint_port_ranges`, `instance_metadata`. The `fqdn_prefix` and `port_ranges` fields are now optional (required only for ELB type).
- **Network Entity Scale Options:** Added `additional_tunnel_options_per_node` nested block.

### Peering

- **CXP Peering Gateway:** Added computed `metadata` block with ILB IP address, VNET details, and subscription information. Data source now supports lookup by ID in addition to name.

### Other

- **IP Reservation:** The `prefix` field is now computed, supporting server-assigned prefix values.
- **Policy Prefix List:** The `prefixes` field is deprecated. Use the `prefix` and `prefix_range` blocks instead.

---

## Bug Fixes

### Import Fixes

- **GCP VPC Connector:** Fixed import to populate `gcp_region`, `vpc_subnet`, and internal subnet ID.
- **IPsec Connector:** Fixed import to populate `segment_options` and `routing_options`.
- **PAN Firewall:** Fixed missing `segment_options` during import.
- **Segment Resource Share:** Fixed import to populate `designated_segment_id`.
- **Policy Prefix List:** Fixed import to populate `prefix` and `prefix_range` data.
- **Checkpoint, Cisco FTDv:** Fixed segment name-to-ID conversion during import.
- **General:** Import now fails with a descriptive error on invalid resource IDs instead of silently producing incomplete state.

### State & Drift Fixes

- **GCP VPC Connector:** Fixed `export_all_subnets` to prevent breaking diff on provider upgrade. Fixed `prefix_list_ids` reordering diffs by switching to `TypeSet` with automatic state migration. Fixed `userInputPrefixes` to send empty array instead of `null`. Fixed regression when removing all `vpc_subnet` blocks.
- **AWS VPC Connector:** Fixed drift from backend-defaulted routing fields.
- **PAN Service:** Fixed perpetual diff for `global_protect_segment_options`. Read now populates `global_protect_enabled`, `panorama_ip_addresses`, and credential ID fields.
- **Policy Prefix List:** Fixed deprecated `prefixes` field handling with backward compatibility. Fixed spurious updates by switching to `TypeSet` with automatic state migration.
- **IPsec Connector:** Fixed `segment_options` flatten to handle `interface{}` type correctly.
- **Cisco FTDv:** Read now populates `segment_ids` from API response.
- **Versa SD-WAN Connector:** Fixed `version` field not being populated during Read due to a typo in the field key.
- **Probe TCP:** Fixed `network_entity_id` being set as a pointer instead of a value.

### Connector & Resource Fixes

- **Aruba Edge Connector:** Fixed instance IDs not being saved to state after apply. Fixed `credentialId` being incorrectly reset when credential fields were updated.
- **Peering Gateway TGW Attachment:** Fixed missing DXGW fields in Read (`type`, `peer_direct_connect_gateway_id`, `peer_allowed_prefixes`). Fixed infinite loop and FAILED state handling during create (now times out after 5 minutes with a clear error).
- **Credentials:** Fixed credential name prefix validation to avoid API validation errors.
- **NAT Policy:** Fixed schema definition.

### Error Handling

- **All Resources:** Delete error messages now include resource name and ID for easier debugging.

---

## Documentation

- Rewrote provider authentication documentation to recommend API key as the primary method.
- Updated Global CIDR List constraints documentation.
- Added OUTBOUND NAT example to `policy_nat_rule` documentation.
- Updated F5 vServer Endpoint SNAT description to indicate ILB with SNAT is allowed.
- Fixed examples that referenced deprecated global CIDR lists to use policy prefix lists instead.

---

## Upgrade Instructions

### From v1.4.4 to v1.5.0

1. **Authentication:**
   - If you use `username`/`password` in your provider configuration, you will see a deprecation warning. Migrate to `api_key` at your convenience. No immediate action required.

2. **Automatic State Migrations:**
   - `alkira_connector_gcp_vpc`: `prefix_list_ids` migrated from `TypeList` to `TypeSet`.
   - `alkira_policy_prefix_list`: `prefix` and `prefix_range` migrated from `TypeList` to `TypeSet`.
   - Run `terraform plan` after upgrading to verify no unexpected changes.

3. **GCP VPC Connector — `export_all_subnets` Validation:**
   - Setting `export_all_subnets = false` without `vpc_subnet` blocks now produces a plan-time error. Either set `export_all_subnets = true` or add `vpc_subnet` blocks.
   - When `export_all_subnets` is not explicitly set, the provider automatically derives the correct value based on `vpc_subnet` presence.

4. **One-Time State Refresh:**
   - Many Read function fixes in this release mean that fields previously missing from state will now correctly reflect API values. After upgrading, the first `terraform plan` may show one-time diffs on the following resources. These diffs are benign — run `terraform apply` once to stabilize state:
     - `alkira_service_pan`: `global_protect_segment_options`, `global_protect_enabled`, credential IDs
     - `alkira_peering_gateway_aws_tgw_attachment`: `type`, `peer_direct_connect_gateway_id`, `peer_allowed_prefixes`
     - `alkira_connector_versa_sdwan`: `version`
     - `alkira_service_checkpoint`: `management_server` fields
     - `alkira_service_cisco_ftdv`: `firepower_management_center` fields
     - Resources for which import was fixed may also show one-time diffs as Read functions now populate previously missing fields.

5. **No Breaking Changes.** This release is backward-compatible. All state migrations are automatic.

---

## Infrastructure

- Upgraded Go to 1.26.
- Upgraded `google.golang.org/grpc` from v1.72.1 to v1.79.3.
- Upgraded `alkira-client-go` to v1.59.0.
- Added build and Slack notification workflow for main and release branches.
- Fixed build versioning to use latest tag instead of ancestor tag.
