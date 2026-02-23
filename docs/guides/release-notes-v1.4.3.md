---
subcategory: "Release Notes"
page_title: "v1.4.3"
description: |-
    Release notes for v1.4.3
---

# Alkira Terraform Provider v1.4.3 Release Notes

Release Date: 2026-02-23

## Overview

Version 1.4.3 is a bug fix release focused on improving `terraform import` reliability across multiple resources, fixing state drift issues, and enhancing error handling.

---

## Bug Fixes

### Import Fixes

Multiple fixes related to `terraform import` to properly populate all fields and improve error handling for the following resources:

- `alkira_connector_gcp_vpc`
- `alkira_connector_ipsec`
- `alkira_connector_aruba_edge`
- `alkira_service_pan`
- `alkira_segment_resource_share`
- `alkira_policy_prefix_list`

### State & Drift Fixes

- **Prefix List:** Fixed prefix reordering on deletion by switching to `TypeSet`, preventing false drift detection. Includes automatic state migration.
- **Internet Applications:** Fixed `protocol` field not being set during import/refresh
- **Field Name Inconsistencies:** Fixed inconsistent field naming across multiple resources

---

## Enhancements

- **Internet Applications:** Added traffic policy warning when using Internet Application resources to clarify policy requirements

---

## Documentation Improvements

- **Global CIDR List:** Updated constraints documentation

---

## Upgrade Instructions

### From v1.4.2 to v1.4.3

1. **Prefix List State Migration:**
   - The `alkira_policy_prefix_list` resource has an automatic state migration from `TypeList` to `TypeSet` for prefix entries
   - Run `terraform plan` after upgrading to verify no unexpected changes
   - This eliminates false drift caused by prefix reordering

2. **Re-import Recommended:**
   - If you previously imported any of the following resources and noticed missing fields, re-import them to populate all attributes:
     - `alkira_connector_gcp_vpc`
     - `alkira_connector_ipsec`
     - `alkira_service_pan`
     - `alkira_segment_resource_share`
     - `alkira_policy_prefix_list`

3. **No Breaking Changes:** This is a fully backward-compatible patch release.
