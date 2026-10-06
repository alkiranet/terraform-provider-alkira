---
subcategory: "Release Notes"
page_title: "v1.7.0"
description: |-
    Release notes for v1.7.0
---

# Alkira Terraform Provider v1.7.0 Release Notes

## Overview

Version 1.7.0 marks credential and key attributes across the provider as sensitive, so they are redacted in plan and apply output. It extends the management-access allow-list to five more connectors and services, and fixes several fields that produced a permanent diff when left out of a configuration.

---

## Enhancements

- **Provider and all credential-bearing resources:** Passwords, API keys, pre-shared keys, BGP authentication keys and license keys are now marked sensitive and are redacted in plan and apply output. This covers `provider` (`password`, `api_key`), the AWS, Azure and GCP credential resources, and the IPSec, AWS Direct Connect, GCP Interconnect, Cisco SD-WAN, Fortinet SD-WAN, Check Point, Cisco FTDv, Fortinet, Infoblox and PAN resources. See Upgrade Instructions.

- **Management-access allow-list extended:** `allow_list` is now available on `alkira_connector_aruba_edge`, `alkira_connector_versa_sdwan`, `alkira_service_pan`, `alkira_service_fortinet` and `alkira_service_checkpoint`, joining `alkira_connector_cisco_sdwan`. When set, only the listed IPv4 CIDRs or addresses can reach the management interface of the service instances.

---

## Bug Fixes

### State & Drift Fixes

- **Azure VNet Connector (`alkira_connector_azure_vnet`):** Fixed an error on update when `peering_gateway_cxp_id` was omitted from the configuration. Where the platform had assigned a peering gateway, the provider tried to clear it on the next update, and the update was rejected because the gateway cannot be changed after provisioning. The field is now computed, so an omitted value keeps the assigned gateway.

- **Infoblox Service (`alkira_service_infoblox`):** Fixed a permanent diff on `grid_master.external`. The value is derived by the platform from whether `ip` is set, so a configured value was never applied. The attribute is now computed and **deprecated** — remove it from your configuration. See Upgrade Instructions.

### Import & Read Fixes

- **NAT Policy (`alkira_policy_nat`) and NAT Rule (`alkira_policy_nat_rule`):** Fixed `category` and `direction` not being read back into state, so they are now populated on refresh and import.

- **Import:** Resource IDs are now validated before the read is attempted, so an invalid ID fails with a clear error instead of producing incomplete state.

---

## Documentation

- **Connectors with `scale_group_id`:** Documented that `scale_group_id` can only be set at create time and cannot be changed after provisioning.

---

## Upgrade Instructions

### From v1.6.0 to v1.7.0

1. **Mark Outputs That Reference Credentials as Sensitive:**
   - Credential and key attributes are now sensitive. An output that references one of them without `sensitive = true` fails with `Output refers to sensitive values`, and the plan stops until it is corrected.
   - Add `sensitive = true` to any such output:
     ```hcl
     output "fw_password" {
       value     = alkira_service_fortinet.example.password
       sensitive = true
     }
     ```
   - Values are redacted in plan and apply output. State is unchanged — it was never encrypted and still is not, so continue to protect your state file.

2. **Remove `grid_master.external` From Infoblox Configurations:**
   - The attribute is deprecated. The platform derives it from whether `ip` is set, so any configured value was already being ignored.
   - Leaving it in place produces a deprecation warning on every plan. Removing it changes nothing about the deployed service.

3. **One-Time State Refresh:**
   - The NAT policy read fix populates `category` and `direction`, which were previously missing from state. After upgrading, run `terraform plan` and expect a one-time diff on `alkira_policy_nat` and `alkira_policy_nat_rule`.
   - Review the diff before applying rather than assuming it is cosmetic — these fields were never refreshed before, so a difference may reflect a real change made outside Terraform.

4. **Compatibility.** Existing configurations continue to plan and apply unchanged, except for outputs referencing newly sensitive values as described in step 1. No state migrations are required.
