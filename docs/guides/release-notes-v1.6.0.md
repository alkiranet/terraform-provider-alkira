---
subcategory: "Release Notes"
page_title: "v1.6.0"
description: |-
    Release notes for v1.6.0
---

# Alkira Terraform Provider v1.6.0 Release Notes

Release Date: 2026-09-23

## Overview

Version 1.6.0 adds the Prisma SD-WAN connector, Infoblox NIOS-X support, and Strata Cloud Manager with advanced routing for PAN firewalls. It also relaxes several fields that were required unnecessarily, and includes drift, import, and validation fixes across a number of resources.

---

## New Resources

- **`alkira_connector_prisma_sdwan`**: Connects a Palo Alto Prisma SD-WAN fabric to the Alkira CXP. Supports segment assignment, size-based scaling, and per-instance configuration. A matching data source is available for lookups by name.

- **`alkira_credential_prisma_sdwan`**: Stores the ION token and secret used to register Prisma SD-WAN devices.

---

## Enhancements

- **Infoblox Service (`alkira_service_infoblox`):** Added NIOS-X support. Each `instance` accepts a `platform` of `NIOS` or `NIOS_X` (defaults to `NIOS`) and a `join_token` for NIOS-X onboarding, and the service accepts a `size`. Fields that apply only to NIOS are now optional — `grid_master`, `shared_secret`, and `model`, `password`, and `type` within `instance` — so a NIOS-X-only service can omit them. `platform` cannot be changed after provisioning.
- **AWS VPC Connector (`alkira_connector_aws_vpc`):** The `aws_account_id` field is now optional. When omitted, the account is detected from the VPC.
- **PAN Service (`alkira_service_pan`):** Added Strata Cloud Manager (SCM) and advanced routing support with new fields `scm_enabled`, `scm_folder`, and `routing_type`. SCM requires `routing_type = "advanced"`, `license_type = "BRING_YOUR_OWN"`, and PAN-OS 10.2.3 or later. Advanced routing is only supported with SCM — `routing_type = "advanced"` without `scm_enabled` is rejected. SCM and Panorama are mutually exclusive. Changing any of these fields on an existing service forces replacement, and the platform rejects the change once the service is provisioned.
- **Cisco SD-WAN Connector (`alkira_connector_cisco_sdwan`):** Added an optional `allow_list` restricting which IPv4 CIDRs or addresses can reach the management interface, matching the Fortinet SD-WAN connector. Supported only on `CAT8000V` connectors.
- **Internet Application (`alkira_internet_application`):** Added a `target` type of `INTERNAL_DNS`, which resolves the target from an FQDN list rather than an address. `INTERNAL_DNS` targets require the new `policy_fqdn_list_id` and must omit `value`; `IP` and `ILB_NAME` targets still require `value` and must omit `policy_fqdn_list_id`.
- **Azure VNet Third Party Connector (`alkira_connector_azure_vnet_third_party`):** Added optional `scale_group_id`. This field cannot be updated after provisioning.
- **Fortinet Service (`alkira_service_fortinet`):** Added optional `alkira_admin_password`, used to authenticate against the FortiGate during first-time provisioning. After provisioning, the field is for record-keeping only — updating it does not rotate the password on deployed FortiGate instances. To change it, update the FortiGate side first, then update this field to match.
- **Advanced IPSec Connector (`alkira_connector_ipsec_adv`):** `customer_end_overlay_ip_reservation_id` in `gateway.tunnel` is now optional, so a tunnel can set `customer_end_overlay_ip` directly instead — useful when the customer-end IP falls outside the ranges available from a reservation. The two are now mutually exclusive and a configuration setting both fails at plan. When `customer_end_overlay_ip` is used, the CXP-end overlay reservation must be a `/32`, and the address must not overlap segment IP blocks or link-local, multicast, broadcast, or loopback ranges.

---

## Bug Fixes

### State & Drift Fixes

- **Bluecat Service (`alkira_service_bluecat`):** Fixed the perpetual diff on `bdds_anycast` and `edge_anycast`. The `ips` and `backup_cxps` fields are now sets rather than lists, so ordering no longer matters, with automatic state migration.
- **Bluecat Service (`alkira_service_bluecat`):** Repeated `bdds_anycast` or `edge_anycast` blocks are now rejected at plan time instead of being silently dropped. Only one block ever took effect. See Upgrade Instructions.
- **Segment Resource (`alkira_segment_resource`), Segment Resource Share (`alkira_segment_resource_share`):** Refresh no longer fails when a referenced segment is already marked for deletion, which previously blocked `terraform destroy`.

### Import & Read Fixes

- **Fortinet Service (`alkira_service_fortinet`):** Fixed `segment_ids` and `instances` not being populated on refresh and import. Previously only one instance was returned, regardless of how many existed.
- **Internet Application (`alkira_internet_application`):** Fixed `inbound_connector_type` not being populated on refresh and import.

### Validation Fixes

- **IPSec Connectors (`alkira_connector_ipsec`, `alkira_connector_ipsec_adv`):** `availability` in `routing_options` no longer accepts `PING`, which was never selectable in the Portal. Use `IPSEC_INTERFACE_PING`, which monitors the tunnel identically. See Upgrade Instructions.
- **Segment Resource (`alkira_segment_resource`), Segment Resource Share (`alkira_segment_resource_share`):** `segment_id` and `designated_segment_id` now reject a segment name at plan time instead of failing slowly against the API. Use the segment ID. See Upgrade Instructions.
- **Segment Resource Share (`alkira_segment_resource_share`):** `traffic_from_end` is now validated against its accepted values.
- **Network Entity Scale Options (`alkira_network_entity_scale_options`):** When `additional_tunnel_options_per_node` is configured, `additional_tunnels_per_node` must match the number of labels. A mismatch is now reported at plan time rather than drifting on every plan. See Upgrade Instructions.

---

## Documentation

- **PAN Service (`alkira_service_pan`):** Documented that credential fields are write-only and are not returned by the API, so they are not populated on import and must be set in configuration.
- **Fortinet Service (`alkira_service_fortinet`):** Added import syntax and example.
- **Bluecat Service (`alkira_service_bluecat`):** Documented instance set comparison and sensitive field redaction.
- **Segment Resource Share (`alkira_segment_resource_share`):** Documented the accepted values for `traffic_from_end` and related parameters.
- **GCP Interconnect Connector (`alkira_connector_gcp_interconnect`):** Documented `20LARGE`, `30LARGE`, `40LARGE`, and `50LARGE` sizes.
- **Azure VNet Third Party Connector (`alkira_connector_azure_vnet_third_party`):** Documented that `scale_group_id` cannot be updated after provisioning.

---

## Upgrade Instructions

### From v1.5.1 to v1.6.0

1. **Verify No PAN Services Are Marked for Replacement:**
   - The new `scm_enabled`, `scm_folder`, and `routing_type` fields force replacement when changed. Services created before this release read them as empty and plan clean.
   - As a precaution, run `terraform plan` and confirm no `alkira_service_pan` resource is marked for replacement before applying.
   - SCM cannot be enabled on an existing service — these fields are immutable after provisioning.

2. **Merge Repeated Bluecat AnyCast Blocks:**
   - Configurations with more than one `bdds_anycast` or `edge_anycast` block on the same service now fail at plan.
   - Combine them into a single block and list all AnyCast IPs in `ips`. Because only one block previously took effect, verify the result against the Portal rather than your previous configuration.

3. **Replace `availability = "PING"` on IPSec Connectors:**
   - Configurations using `availability = "PING"` in `routing_options` now fail validation. Replace it with `IPSEC_INTERFACE_PING`.
   - For connectors already storing `PING`, the next plan shows a single in-place change that resolves on apply.

4. **Use Segment IDs Rather Than Names:**
   - `segment_id` on `alkira_segment_resource` and `designated_segment_id` on `alkira_segment_resource_share` now reject segment names at plan time. Use `alkira_segment.<name>.id` or the numeric ID.
   - Leading zeros are also rejected — use `690`, not `0690`.

5. **Align `additional_tunnels_per_node` With Tunnel Option Labels:**
   - Where `additional_tunnel_options_per_node` is configured, `additional_tunnels_per_node` must equal the number of labels, or the plan now fails naming the block.
   - Such configurations previously showed drift on every plan.

6. **Review Diffs After the First Refresh:**
   - Read fixes in this release populate fields that were previously missing from state, so the first `terraform plan` after upgrading may show changes.
   - **`alkira_service_fortinet` (`segment_ids`, `instances`) — review before applying.** These fields were never refreshed before, so drift was not detected. A diff here may be a real difference between your configuration and the deployed service, and applying it will change the deployed resource. Confirm the plan matches your intent rather than applying it blindly.
   - `alkira_internet_application` (`inbound_connector_type`) — cosmetic; the value is read back from the platform and converges on the first apply.
   - The Bluecat state migration runs on first refresh and may show a one-time diff if the same AnyCast IP was listed twice.

7. **Compatibility.** Existing configurations continue to plan and apply unchanged, except as noted in steps 2 through 5 — each of which rejects a configuration that could not work correctly before. The Bluecat state migration is automatic; no other state migrations are required.
