---
subcategory: "Release Notes"
page_title: "v1.5.0"
description: |-
    Release notes for v1.5.0
---

# Alkira Terraform Provider v1.5.0 Release Notes

Release Date: 2026-04-23

## Overview

Version 1.5.0 is a feature release that introduces new resources (Bluecat DNS service, F5 Internal Load Balancer, VNET Third Party Connector), adds scalegroup support for the Aruba Edge connector, improves credential security with WriteOnly attributes, and includes 40+ bug fixes across import, state drift, and connector reliability.

---

## New Resources

- **`alkira_service_bluecat`**: Full support for Bluecat DNS service, including BDDS/Edge instances, anycast configuration, and segment assignment.
- **`alkira_connector_vnet_third_party`**: Support for Azure VNET Third Party Connector.
- **`alkira_connector_vnet_third_party_attachment`**: Support for VNET Third Party Connector Attachment.
- **`alkira_service_f5_ilb`**: Terraform support for F5 Internal Load Balancer service.

---

## Enhancements

### Connectors

- **Aruba Edge Connector:** Added scalegroup support, enabling scale group configuration for Aruba Edge connectors.
- **GCP VPC Connector:** Added `export_all_subnets` field to `gcp_routing` schema, allowing control over subnet export behavior.
- **AWS VPC Connector:** When `vpc_subnet` entries are specified, `exportAllSubnets` is now forced to `false` to prevent conflicting configuration.

### Services

- **F5 LB Service & vServer Endpoint:** Added metadata fields for F5 LB service and F5 vServer endpoint resources.
- **Network Entity Scale Options:** Added additional tunnel option per node for scale options. Removed deprecated `networkEntitySubType` field.

### Credentials

- **WriteOnly Attributes:** Sensitive credential fields now use Terraform's `WriteOnly` attribute for the following resources, preventing sensitive values from being stored in state:
  - `alkira_credential_aws_vpc`
  - `alkira_credential_azure_vnet`
  - `alkira_credential_gcp_vpc`
  - `alkira_credential_oci_vcn`

### Peering Gateway

- **CXP Peering Gateway:** Updated to support Third Party VNET Connector.

### Validation

- **Bluecat Service:** Added `MinItems` validation to `segment_ids` to catch missing segment assignment at plan time.
- **IP Reservation:** The `prefix` field can now be set by the backend, supporting server-assigned prefix values.

---

## Bug Fixes

### Import Fixes

- **GCP VPC Connector:** Fixed import to populate `gcp_region`, `vpc_subnet`, and internal subnet ID.
- **IPsec Connector:** Fixed import to populate `segment_options` and `routing_options`.
- **PAN Firewall:** Fixed missing `segment_options` during import.
- **Segment Resource Share:** Fixed import to populate `designated_segment_id`.
- **Policy Prefix List:** Fixed import to populate `prefix` and `prefix_range` data.
- **Policy Rule:** Fixed `ConflictsWith` error during import by only setting fields when non-empty.
- **Bluecat Service:** Fixed multiple import issues including field mapping and hostname handling.
- **Internet Application:** Fixed `protocol` field not being set during import/refresh.
- **General:** Import now fails clearly with a descriptive error on invalid resource IDs instead of silently producing incomplete state.

### State & Drift Fixes

- **GCP VPC Connector:** Fixed `export_all_subnets` schema to prevent breaking diff on provider upgrade. Fixed `prefix_list_ids` reordering diffs by switching to `TypeSet` with automatic state migration. Fixed `userInputPrefixes` to send empty array instead of `null`.
- **AWS VPC Connector:** Fixed Read function to populate routing and TGW fields. Fixed drift from backend-defaulted routing fields. Fixed `export_all_subnets` regression when `vpc_subnet` blocks are removed.
- **PAN Service:** Fixed perpetual diff for `global_protect_segment_options`.
- **Bluecat Service:** Fixed sensitive fields showing as diff. Fixed first-apply diff caused by computed `id` and `name` fields (now checks hostname). Fixed `TypeList` spurious updates. Fixed instance ordering/index issue during updates.
- **Policy Prefix List:** Fixed deprecated `prefixes` field handling. Fixed `TypeList` spurious updates.
- **Prefix List:** Fixed prefix reordering on deletion by switching to `TypeSet` with automatic state migration.
- **IPsec Connector:** Fixed `segment_options` flatten to handle `interface{}` type correctly.
- **segment_options:** Fixed handling when one segment option group is `null` but another is used.
- **Field Names:** Fixed inconsistent field naming across multiple resources.

### Connector & Resource Fixes

- **Aruba Edge Connector:** Fixed instance IDs not being saved to state after apply. Fixed `credentialId` being incorrectly reset when credential fields were updated.
- **Peering Gateway TGW Attachment:** Fixed missing DXGW fields in Read. Fixed infinite loop and FAILED state handling during create.
- **Credentials:** Fixed credential name prefix validation to avoid using input fields that trigger API validation errors.
- **NAT Policy:** Fixed schema definition.
- **Bluecat Service:** Made instance `name` field read-only (API sets this to hostname).

### Error Handling

- **Provisioning:** Failure errors now include detailed reason when available.
- **Delete Operations:** Error messages now include resource name and ID for easier debugging.

---

## Documentation

- Updated Global CIDR List constraints documentation.
- Added OUTBOUND NAT example to `policy_nat_rule` documentation.
- Updated F5 vServer Endpoint SNAT description to indicate ILB with SNAT is allowed.
- Fixed examples that referenced deprecated global CIDR lists to use policy prefix lists instead.

---

## Infrastructure

- Upgraded Go version to 1.26.
- Upgraded `google.golang.org/grpc` from v1.72.1 to v1.79.3.
- Bumped `alkira-client-go` to v1.59.0.
- Added build and Slack notification workflow for main and release branches.
- Added PR checks for `release/*` branches.
- Fixed build versioning to use latest tag instead of ancestor tag.

---

## Upgrade Instructions

### From v1.4.4 to v1.5.0

1. **State Migrations:**
   - `alkira_connector_gcp_vpc`: Automatic state migration for `prefix_list_ids` from `TypeList` to `TypeSet`.
   - `alkira_policy_prefix_list`: Automatic state migration for prefix entries from `TypeList` to `TypeSet`.
   - `alkira_service_bluecat`: `bdds_options` and `edge_options` changed from `TypeSet` to `TypeList`.
   - Run `terraform plan` after upgrading to verify no unexpected changes.

2. **Credential Resources:**
   - Sensitive fields on `alkira_credential_aws_vpc`, `alkira_credential_azure_vnet`, `alkira_credential_gcp_vpc`, and `alkira_credential_oci_vcn` now use `WriteOnly` attributes. These values will no longer appear in state files.

3. **Re-import Recommended:**
   - If you previously imported any of the following resources and noticed missing fields, re-import them:
     - `alkira_connector_gcp_vpc`
     - `alkira_connector_ipsec`
     - `alkira_connector_aruba_edge`
     - `alkira_service_pan`
     - `alkira_segment_resource_share`
     - `alkira_policy_prefix_list`
     - `alkira_policy_rule`

4. **No Breaking Changes:** This release is fully backward-compatible. All state migrations are automatic.
