# Alkira Terraform Provider v1.5.0 — Upgrade Verification Report

**Date**: 2026-05-01
**Baseline**: v1.4.4 | **Upgrade target**: v1.5.0 candidate (`release/2026-03-16-64`)
**Environments**: Azure CXP (USEAST-AZURE-1) and AWS CXP (US-WEST)

---

## Summary

| Metric | Value |
|---|---|
| Resources tested | 22 |
| Zero drift | 17 |
| Pre-existing perpetual diffs (v1.4.4) | 4 (`segment_options` system zones) |
| NEW perpetual diffs (release-64) | 2 (`management_server` / `firepower_management_center` cycling) |
| Replacements or destroys | 0 |
| Backward compatibility broken | No |

---

## Perpetual Diffs

### Pre-existing in v1.4.4 (NOT a regression)

Read functions populate `segment_options` from the API including system-generated zones (DEFAULT, ALKIRA_MGMT_ZONE). User config has no `segment_options` block. Apply removes them, Read adds them back. Verified by running create + second plan entirely on v1.4.4 — same diff appears.

| Resource | Zone | Fix |
|---|---|---|
| `alkira_connector_ipsec` | DEFAULT | `dd25dc37` (AK-67049) |
| `alkira_service_checkpoint` | DEFAULT | Same |
| `alkira_service_cisco_ftdv` | DEFAULT | Same |
| `alkira_service_pan` | ALKIRA_MGMT_ZONE | Same |

### NEW in release-64 (commit `9a131d9a` AK-63243)

The `management_server` (Checkpoint) and `firepower_management_center` (Cisco FTDv) blocks cycle on every plan. Does NOT occur on v1.4.4.

| Resource | Block | Cause |
|---|---|---|
| `alkira_service_checkpoint` | `management_server` | Deflate field names corrected + `credential_id` regeneration |
| `alkira_service_cisco_ftdv` | `firepower_management_center` | Same |

**Root cause**: Commit `9a131d9a` fixed deflate field names (`"segment"` -> `"segment_id"`, `"user_name"` -> `"username"`) to match the schema. In v1.4.4, mismatched keys were silently ignored by Terraform — fields were never written to state. Now they are, changing the `TypeSet` hash. Combined with `credential_id` being regenerated on each update, the hash never stabilizes.

**Fix**: Change `management_server`/`firepower_management_center` from `TypeSet` to `TypeList` with `MaxItems: 1` (avoids hash-based comparison).

### Verification Evidence

| Test | Provider | management_server diff? | segment_options diff? |
|---|---|---|---|
| Create + plan (same version) | v1.4.4 | No | Yes (pre-existing) |
| Create v1.4.4, plan release-64 | release-64 | **Yes (new)** | Yes (pre-existing) |
| Apply release-64, plan again | release-64 | **Yes (perpetual)** | Yes (perpetual) |

---

## Clean Resources (17)

| Resource Type | Tested |
|---|---|
| `alkira_segment` | 2 instances |
| `alkira_group` | 2 instances |
| `alkira_billing_tag` | 1 instance |
| `alkira_list_global_cidr` | 2 instances |
| `alkira_connector_internet_exit` | 1 instance |
| `alkira_internet_application` | 1 instance |
| `alkira_policy` | 1 instance |
| `alkira_policy_rule` | 2 instances (protocol=any, protocol=tcp) |
| `alkira_policy_rule_list` | 1 instance |
| `alkira_policy_nat_rule` | 1 instance |
| `alkira_policy_prefix_list` | 3 instances (with desc, without desc, with ranges) |

---

## Positive Findings

- **`alkira_policy_prefix_list`**: Hash function change (now includes `description`) did NOT cause drift. Three variants tested — all clean.
- **`alkira_policy_rule`**: Protocol field stable. No drift.
- **`alkira_service_pan`**: `global_protect_segment_options` fix is working correctly — no drift on that field.

---

## Code-Level Verification (no live infrastructure available)

Read function diffs between v1.4.4 and release-64 were inspected for all remaining resources.

### No Read changes — zero drift risk

| Resource | Change in release-64 |
|---|---|
| `alkira_connector_aws_vpc` | Delete error message only |
| `alkira_connector_aws_dx` | Delete error message only |
| `alkira_connector_azure_vnet` | Delete error message only |
| `alkira_segment_resource_share` | Delete error message only |

### One-time drift expected (stabilizes after single apply)

| Resource | Change | Root Cause |
|---|---|---|
| `alkira_connector_versa_sdwan` | `version` field newly populated | Typo fix: deflate key `"vesion"` -> `"version"`. v1.4.4 silently dropped this field (key mismatch with schema). |
| `alkira_probe_tcp` | `network_entity_id` value corrected | `&probe.NetworkEntity.ID` (pointer) -> `probe.NetworkEntity.ID` (value). May have stored incorrect value in state. |
| `alkira_peering_gateway_aws_tgw_attachment` | 3 computed fields newly populated | Read now sets `type`, `peer_direct_connect_gateway_id`, `peer_allowed_prefixes`. Not `TypeSet` — no perpetual diff risk. |

### Behavioral change (already analyzed)

| Resource | Change |
|---|---|
| `alkira_connector_gcp_vpc` | `export_all_subnets` auto-derived from `vpc_subnet` presence. Validation tightened: `export_all_subnets=false` without `vpc_subnet` now errors at plan time (previously failed at API). |

---

## Release Recommendation: GO WITH CAUTION

### Must-fix before release

1. **Checkpoint `management_server` perpetual diff** — change to `TypeList` with `MaxItems: 1`. Commit `9a131d9a`.
2. **Cisco FTDv `firepower_management_center` perpetual diff** — same fix, same commit.

### Nice-to-have

3. Cherry-pick `dd25dc37` (AK-67049) — fixes pre-existing `segment_options` system zone diffs.

### After fix

Re-run verification, confirm second `terraform plan` returns clean for all resources.
