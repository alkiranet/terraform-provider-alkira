---
subcategory: "Release Notes"
page_title: "v1.5.1"
description: |-
    Release notes for v1.5.1
---

# Alkira Terraform Provider v1.5.1 Release Notes

## Overview

Version 1.5.1 adds the Azure Virtual Hub connector and an inter-CXP routing policy resource, enables API serialization by default, and includes a number of import/refresh (drift) fixes across existing resources.

---

## New Resources

~> **Beta:** The resources below are released as beta features and are available on request until the full release. Contact Alkira to have them enabled for your tenant.

- **`alkira_connector_azure_vhub`** (Beta): Connects an Azure Virtual Hub (vWAN) to the Alkira CXP. Supports segment assignment and size-based scaling. A matching data source is available for lookups by name or ID.

- **`alkira_policy_inter_cxp_routing`** (Beta): Manages inter-CXP route policies. Supports AS-path filtering and prepending, community and extended-community tagging, prefix / group / segment-resource matching, composite match conditions, and drop-all rules.

---

## Enhancements

- **Provider:** API serialization is now enabled by default. The `serialization_enabled` attribute now defaults to `true` (previously `false`), and can be controlled with the `ALKIRA_API_SERIALIZATION_ENABLED` environment variable (which now parses both `true` and `false`). Explicit configuration always takes precedence. See Upgrade Instructions.
- **Internet Application (`alkira_internet_application`):** Added an optional `description` field.
- **AWS Direct Connect Connector (`alkira_connector_aws_dx`):** Added `loopback_prefixes` for loopback IP allocation across Direct Connect instances with tunnel scale options.
- **Provisioning:** Resources now emit a non-fatal `CONFIGURATION CHANGES SKIPPED` warning on update when the resource is in a `FAILED` provision state. In that state the backend re-provisions the previously saved configuration and skips configuration changes until the resource recovers.

---

## Bug Fixes

### Import & Read Fixes

Multiple fixes to properly populate fields on `terraform import` / refresh for the following resources:

- **Infoblox Service (`alkira_service_infoblox`):** `name` and `segment_ids` are now populated.
- **Remote Access Connector (`alkira_connector_remote_access`):** `authentication_mode` and `subnet` are now populated correctly on import.
- **VMware SD-WAN Connector (`alkira_connector_vmware_sdwan`):** `vmware_sdwan_segment_name` is now populated on refresh.
- **PAN Service (`alkira_service_pan`):** `pan_credential_name` is now populated in Read.
- **F5 vServer Endpoint (`alkira_service_f5_vserver_endpoint`):** Destination endpoint port ranges and IP addresses are now populated in Read.
- **OCI VCN Connector (`alkira_connector_oci_vcn`):** VCN routing configuration (`vcn_cidr`, `vcn_subnet`, route tables) is now populated in Read.
- **Versa SD-WAN Connector (`alkira_connector_versa_sdwan`):** `global_tenant_id` and `versa_controller_host` are now populated in Read.
- **Akamai Prolexic Connector (`alkira_connector_akamai_prolexic`):** Tunnel/overlay configuration is now populated in Read, removing refresh drift.

---

## Documentation

- **NAT Policy (`alkira_policy_nat`):** Clarified `included_group_ids` / `excluded_group_ids` — Terraform requires group IDs (including implicit connector groups), even where the UI allows selecting connectors directly.
- **IP Reservation (`alkira_ip_reservation`):** Clarified the `first_ip_assignment` documentation — the backend retains this value once set, so omitting it from the configuration defers to the stored value rather than clearing it.

---

## Upgrade Instructions

### From v1.5.0 to v1.5.1

1. **API Serialization Enabled by Default:**
   - `serialization_enabled` now defaults to `true`, so provider API calls are serialized unless you opt out. This changes API call timing and throughput but does not affect resource state.
   - To restore the previous behavior, set `serialization_enabled = false` in the provider block, or set `ALKIRA_API_SERIALIZATION_ENABLED=false`.

2. **One-Time State Refresh:**
   - Several Read fixes in this release populate fields that were previously missing from state. After upgrading, run `terraform plan` and expect one-time diffs on the following resources:
     - `alkira_service_infoblox` (`name`, `segment_ids`)
     - `alkira_connector_remote_access` (`authentication_mode`, `subnet`)
     - `alkira_connector_vmware_sdwan` (`vmware_sdwan_segment_name`)
     - `alkira_service_pan` (`pan_credential_name`)
     - `alkira_service_f5_vserver_endpoint` (destination endpoints)
     - `alkira_connector_oci_vcn` (VCN routing fields)
     - `alkira_connector_versa_sdwan` (`global_tenant_id`, `versa_controller_host`)
   - These diffs are benign. Run `terraform apply` once to stabilize state.

3. **No Breaking Changes.** This release is backward-compatible. No state migrations are required.
