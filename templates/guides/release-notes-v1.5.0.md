---
subcategory: "Release Notes"
page_title: "v1.5.0"
description: |-
    Release notes for v1.5.0
---

# Alkira Terraform Provider v1.5.0 Release Notes

Release Date: 2026-04-28

## Overview

Version 1.5.0 is a feature release that introduces new resources (Bluecat DNS service, F5 Internal Load Balancer, VNET Third Party Connector), adds scalegroup support for the Aruba Edge connector. It also includes bug fixes across import and documentaion improvements.

---

## New Resources

- **`alkira_service_bluecat`**: Support for Bluecat DNS service, including BDDS/Edge instances and anycast configuration.
- **`alkira_connector_vnet_third_party`**: Support for Azure VNET Third Party Connector.
- **`alkira_connector_vnet_third_party_attachment`**: Support for VNET Third Party Connector Attachment.
- **`alkira_service_f5_ilb`**: Terraform support for F5 Internal Load Balancer service.

---

## Enhancements

### Connectors

- **Aruba Edge Connector:** Added scalegroup support, enabling scale group configuration for Aruba Edge connectors.
- **GCP VPC Connector:** Added `export_all_subnets` field to `gcp_routing` schema, allowing control over subnet export behavior.
- **AWS VPC Connector:** When `vpc_subnet` entries are specified, `exportAllSubnets` is now forced to `false` to prevent conflicting configuration.
- **Network Entity Scale Options:** Added additional tunnel option per node for scale options. Removed deprecated `networkEntitySubType` field.
- **IP Reservation:** The `prefix` field can now be set by the backend, supporting server-assigned prefix values.

---

## Bug Fixes

### Connector & Resource Fixes

- **GCP VPC Connector:** Fixed import to populate `gcp_region`, `vpc_subnet`, and internal subnet ID.
- **IPsec Connector:** Fixed import to populate `segment_options` and `routing_options`.
- **PAN Firewall:** Fixed missing `segment_options` during import.
- **Segment Resource Share:** Fixed import to populate `designated_segment_id`.
- **Policy Prefix List:** Fixed import to populate `prefix` and `prefix_range` data.
- **Aruba Edge Connector:** Fixed instance IDs not being saved to state after apply. Fixed `credentialId` being incorrectly reset when credential fields were updated.
- **Credentials:** Fixed credential name prefix validation to avoid using input fields that trigger API validation errors.
- **NAT Policy:** Fixed schema definition.
- **General:** Import now fails clearly with a descriptive error on invalid resource IDs instead of silently producing incomplete state.

### State & Drift Fixes

- **GCP VPC Connector:** Fixed `export_all_subnets` schema to prevent breaking diff on provider upgrade. Fixed `prefix_list_ids` reordering diffs by switching to `TypeSet` with automatic state migration. Fixed `userInputPrefixes` to send empty array instead of `null`.
- **AWS VPC Connector:** Fixed drift from backend-defaulted routing fields. Fixed `export_all_subnets` regression when `vpc_subnet` blocks are removed.
- **PAN Service:** Fixed perpetual diff for `global_protect_segment_options`.
- **Policy Prefix List:** Fixed deprecated `prefixes` field handling. Fixed `TypeList` spurious updates.
- **Prefix List:** Fixed prefix reordering on deletion by switching to `TypeSet` with automatic state migration.
- **IPsec Connector:** Fixed `segment_options` flatten to handle `interface{}` type correctly.
- **Field Names:** Fixed inconsistent field naming across multiple resources.

### Error Handling

- **Delete Operations:** Error messages now include resource name and ID for easier debugging.

---

## Documentation

- Updated Global CIDR List constraints documentation.
- Added OUTBOUND NAT example to `policy_nat_rule` documentation.
- Updated F5 vServer Endpoint SNAT description to indicate ILB with SNAT is allowed.
- Fixed examples that referenced deprecated global CIDR lists to use policy prefix lists instead.

---

## Upgrade Instructions

### From v1.4.4 to v1.5.0

1. **State Migrations:**
   - `alkira_connector_gcp_vpc`: Automatic state migration for `prefix_list_ids` from `TypeList` to `TypeSet`.
   - `alkira_policy_prefix_list`: Automatic state migration for prefix entries from `TypeList` to `TypeSet`.
   - `alkira_service_bluecat`: `bdds_options` and `edge_options` changed from `TypeSet` to `TypeList`.
   - Run `terraform plan` after upgrading to verify no unexpected changes.

2. **Re-import Recommended:**
   - If you previously imported any of the following resources and noticed missing fields, re-import them:
     - `alkira_connector_gcp_vpc`
     - `alkira_connector_ipsec`
     - `alkira_connector_aruba_edge`
     - `alkira_service_pan`
     - `alkira_segment_resource_share`
     - `alkira_policy_prefix_list`

3. **No Breaking Changes:** This release is fully backward-compatible. All state migrations are automatic.
