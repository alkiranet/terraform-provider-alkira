---
subcategory: "Release Notes"
page_title: "v1.6.0"
description: |-
    Release notes for v1.6.0
---

# Alkira Terraform Provider v1.6.0 Release Notes

## Overview

Version 1.6.0 adds the Prisma SD-WAN connector, Infoblox NIOS-X support, and Strata Cloud Manager with advanced routing on the PAN firewall service. It also relaxes several fields that were required but did not need to be, most notably `aws_account_id` on the AWS VPC connector. The remainder is drift and import fixes — including the Bluecat AnyCast blocks, which no longer produce a perpetual no-op diff — together with stricter, earlier validation on fields that previously failed slowly against the API.

This release also fixes a logging issue that could write IPSec pre-shared keys into provider debug logs. See Upgrade Instructions.

---

## New Resources

- **`alkira_connector_prisma_sdwan`**: Connects a Palo Alto Prisma SD-WAN fabric to the Alkira CXP. Supports segment assignment, size-based scaling and per-instance configuration. A matching data source is available for lookups by name or ID.

- **`alkira_credential_prisma_sdwan`**: Stores the Prisma SD-WAN service account credentials referenced by the connector.

---

## Enhancements

- **Infoblox Service (`alkira_service_infoblox`):** Added support for NIOS-X instances. Each entry in `instance` accepts a `platform` of `NIOS` or `NIOS_X`, defaulting to `NIOS`, along with a `join_token` for NIOS-X onboarding, and the service accepts a `size`. Fields that only apply to NIOS are now optional rather than required — `grid_master`, `shared_secret`, and `model`, `password` and `type` within `instance` — so a NIOS-X-only service can omit them. `platform` cannot be changed once the service is provisioned, and that is now rejected at plan time rather than failing during apply.

- **AWS VPC Connector (`alkira_connector_aws_vpc`):** `aws_account_id` is now optional. When omitted, the account is detected from the VPC. Existing configurations that set it are unaffected.

- **PAN Service (`alkira_service_pan`):** Added support for Strata Cloud Manager and advanced routing through three new optional fields — `scm_enabled`, `scm_folder` and `routing_type`. SCM mode requires `routing_type = "advanced"` and PAN-OS 10.2.3 or later; advanced routing alone requires 10.2 or later. SCM and Panorama are mutually exclusive. All three force replacement when changed, and the platform rejects changes to them on an already-provisioned service. See Upgrade Instructions.

- **Cisco SD-WAN Connector (`alkira_connector_cisco_sdwan`):** Added an optional `allow_list` of IPv4 CIDRs or addresses that may reach the management interface of the connector's instances, bringing it to parity with the Fortinet SD-WAN connector.

- **Internet Application (`alkira_internet_application`):** Added an optional `policy_fqdn_list_id` to the `target` block, allowing a target to be expressed as an FQDN list. `value` inside `target` is now optional, since a target defined by `policy_fqdn_list_id` does not need one.

- **Azure VNet Third Party Connector (`alkira_connector_azure_vnet_third_party`):** Added an optional `scale_group_id` for associating the connector with a scale group. It cannot be updated after provisioning.

- **Fortinet Service (`alkira_service_fortinet`):** Added an optional `alkira_admin_password` on the credential resource.

- **Advanced IPSec Connector (`alkira_connector_ipsec_adv`):** `customer_end_overlay_ip_reservation_id` inside `gateway.tunnel` is now optional rather than required, so tunnels that do not reserve a customer-end overlay IP no longer need a placeholder value.

---

## Bug Fixes

### Security

- **IPSec Connectors (`alkira_connector_ipsec`):** A shared helper logged the full contents of every string list it converted, and `preshared_keys` passed through it. With provider logging enabled, pre-shared keys were written to the log in cleartext. The helper now logs only the number of entries. If you have run this provider with `TF_LOG` set and shared or retained those logs, treat the affected pre-shared keys as exposed and rotate them.

### Drift & Spurious Diffs

- **Bluecat Service (`alkira_service_bluecat`):** `ips` and `backup_cxps` inside the `bdds_anycast` and `edge_anycast` blocks are now `Set of String` instead of `List of String`. The API returns these values in a different order than configured, and as ordered lists they re-hashed the surrounding block, so every plan showed the AnyCast block as removed and re-added even when nothing had changed. Order no longer matters. A state migration re-keys existing state on first refresh; it runs automatically and requires no action.

- **Bluecat Service (`alkira_service_bluecat`):** Repeated `bdds_anycast` or `edge_anycast` blocks are now rejected at plan time instead of being silently discarded. Previously only one block took effect — and because the blocks are unordered, which one survived was arbitrary — so a configuration with several blocks never described the infrastructure it produced. See Upgrade Instructions.

### Import & Read Fixes

- **Fortinet Service (`alkira_service_fortinet`):** `segment_ids` and `instances` are now populated on refresh and import. `segment_ids` was built as a list of integers, which the schema rejected as a set of strings, and the resulting error was discarded — so Read silently never updated the field. `instances` stopped after the first entry the state did not already know about, landing exactly one instance no matter how many existed. After `terraform import`, both fields now arrive complete.

- **Internet Application (`alkira_internet_application`):** `inbound_connector_type` is now set during Read, so it is populated on refresh and after import.

### Validation

- **IPSec Connectors (`alkira_connector_ipsec`, `alkira_connector_ipsec_adv`):** `availability` inside `routing_options` no longer accepts `PING`. The accepted values are `IKE_STATUS` and `IPSEC_INTERFACE_PING`, matching the two options the portal offers. `PING` was never selectable in the UI, and a connector holding it rendered as though `IKE Status` were selected while the backend still stored `PING`. `PING` and `IPSEC_INTERFACE_PING` resolve to the same tunnel probe, so availability monitoring is unchanged. See Upgrade Instructions.

- **Segment Resource (`alkira_segment_resource`) and Segment Resource Share (`alkira_segment_resource_share`):** `segment_id` and `designated_segment_id` now reject a segment name at plan time. Both fields are documented as IDs, but any alphanumeric string was accepted and handed straight to the API, where a name returns an error that the client retries five times — so an apply sat for more than two minutes before failing with an unhelpful message. The value is now validated before any API call, with an error naming the field and pointing at `alkira_segment.example.id`. See Upgrade Instructions.

- **Segment Resource Share (`alkira_segment_resource_share`):** `traffic_from_end` is now validated against its accepted values.

- **Network Entity Scale Options (`alkira_network_entity_scale_options`):** When `additional_tunnel_options_per_node` labels are configured, `additional_tunnels_per_node` must now equal the number of labels. The API derives the count from the labels, so a differing value — including the zero produced by omitting the field — diverged from what the API stored and drifted on every plan. The mismatch is now reported at plan time. See Upgrade Instructions.

- **Segment Resource (`alkira_segment_resource`) and Segment Resource Share (`alkira_segment_resource_share`):** Read no longer aborts when the segment lookup fails. A segment already marked for deletion could not be resolved by name, which aborted the refresh and blocked `terraform destroy`. Read now warns, keeps the existing `segment_id` in state, and refreshes the remaining fields.

---

## Documentation

- **PAN Service (`alkira_service_pan`):** Documented that credential fields are write-only and are not returned by the API, so they are not populated by `terraform import` and must be supplied in configuration.

- **Fortinet Service (`alkira_service_fortinet`):** Added `terraform import` documentation and an import script.

- **Bluecat Service (`alkira_service_bluecat`):** Documented that instance blocks are compared as a set, and that sensitive instance fields are redacted rather than returned by the API.

- **Segment Resource Share (`alkira_segment_resource_share`):** Documented the accepted values for `traffic_from_end` and the other share parameters.

- **GCP Interconnect Connector (`alkira_connector_gcp_interconnect`):** Documented `20LARGE`, `30LARGE`, `40LARGE` and `50LARGE` as accepted sizes.

- **Azure VNet Third Party Connector (`alkira_connector_azure_vnet_third_party`):** Documented that `scale_group_id` cannot be updated after provisioning.

---

## Upgrade Instructions

### From v1.5.1 to v1.6.0

1. **Rotate IPSec Pre-Shared Keys If You Have Kept Debug Logs.**
   - Earlier versions wrote `preshared_keys` in cleartext to the provider log whenever logging was enabled. If you have run `terraform` with `TF_LOG` set and those logs were retained, shared with support, or stored in CI artifacts, rotate the affected pre-shared keys.
   - No configuration change is required; the logging itself is fixed in this release.

2. **PAN Services Are Not Affected by the New SCM Fields.**
   - `scm_enabled`, `scm_folder` and `routing_type` force replacement when changed. They are not returned for services created before this release, so existing services read them as empty and plan clean.
   - As a precaution, run `terraform plan` after upgrading and confirm no `alkira_service_pan` resource is marked for replacement before applying.
   - Enabling SCM on an existing service is not supported in place: the platform rejects changes to these fields once a service is provisioned.

3. **Bluecat AnyCast Blocks Must Not Be Repeated.**
   - If your configuration has more than one `bdds_anycast` or more than one `edge_anycast` block on the same `alkira_service_bluecat`, `terraform plan` now fails.
   - Merge them into a single block and list every AnyCast IP in `ips`. Only one of the blocks was ever taking effect, so the merged block should list the addresses you intended all along — check it against the portal rather than against your previous configuration.

4. **Replace `availability = "PING"` on IPSec Connectors.**
   - Configurations setting `availability = "PING"` inside `routing_options` now fail validation. Replace it with `IPSEC_INTERFACE_PING`, which probes the tunnel the same way.
   - For connectors whose stored value is `PING`, the next plan shows a single in-place change that converges on apply. Note that `terraform plan -refresh-only` reports the stored value as-is without flagging it; a normal apply is needed to update the connector.

5. **Use Segment IDs, Not Names.**
   - `segment_id` on `alkira_segment_resource` and `designated_segment_id` on `alkira_segment_resource_share` now reject anything that is not a segment ID, and the failure blocks `terraform plan`, `apply` and `destroy` for the whole root module until it is corrected.
   - Replace a segment name with `alkira_segment.<name>.id` or the segment's numeric ID. A name never worked — it failed slowly against the API — but it did pass validation, so it may be present in configurations written to match an earlier import.
   - Leading zeros are also rejected. `0690` previously resolved to segment `690` and then produced a permanent diff, because Read wrote `690` back against a configuration that still said `0690`. Remove the leading zero.

6. **Match `additional_tunnels_per_node` to Your Tunnel Option Labels.**
   - If a `segment_scale_options` block sets `additional_tunnel_options_per_node`, `additional_tunnels_per_node` must equal the number of labels, or `terraform plan` now fails naming the block.
   - Such configurations were already drifting on every plan, because the API derived the count from the labels and the stored value never matched.

7. **One-Time State Refresh.**
   - The Fortinet and Internet Application Read fixes populate fields that were previously missing from state. After upgrading, run `terraform plan` and expect one-time diffs on:
     - `alkira_service_fortinet` (`segment_ids`, `instances`)
     - `alkira_internet_application` (`inbound_connector_type`)
   - The Bluecat state migration also runs on first refresh and may show a one-time diff on `bdds_anycast` / `edge_anycast` if the same AnyCast IP was listed twice, since duplicates collapse when the list becomes a set.
   - These diffs are benign. Run `terraform apply` once to stabilize state.

8. **Compatibility.** Existing configurations continue to plan and apply unchanged, with the exceptions listed in steps 3 through 6 — each of which rejects a configuration that could not work correctly before. The Bluecat state migration is automatic; no other state migrations are required.
