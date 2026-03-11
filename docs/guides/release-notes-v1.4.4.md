---
subcategory: "Release Notes"
page_title: "v1.4.4"
description: |-
    Release notes for v1.4.4
---

# Alkira Terraform Provider v1.4.4 Release Notes

Release Date: 2026-03-11

## Overview

Version 1.4.4 is a bug fix release focused on improving `terraform import` reliability, fixing state drift issues, and enhancing error handling.

---

## Bug Fixes

### Import Fixes

- **Policy Rule:** Fixed a `ConflictsWith` error that occurred during `terraform import` for `alkira_policy_rule` resources.
- **PAN FW:** Fixed missing `segment_options` not being populated during `terraform import` for `alkira_service_pan` resources.
- **Import Validation:** Fixed import error handling across all resources to fail clearly when an invalid ID is provided, rather than incorrectly reporting `Import successful!`.

### State & Drift Fixes

- **Prefix List:** Fixed prefix reordering on deletion by switching `prefix` and `prefix_range` fields from `TypeList` to `TypeSet`, preventing false drift detection. Includes automatic state migration.

### Credential Management

- **Aruba Edge:** Fixed an issue where `credentialId` was incorrectly reset when credential fields were updated. The credential ID is now preserved when it is already set and only credential field values change.

### Error Messages

- **Provisioning:** Updated provisioning failure error messages to include detailed information to help diagnose configuration issues.

---

## Upgrade Instructions

### From v1.4.3 to v1.4.4

1. **Prefix List State Migration:**
   - The `alkira_policy_prefix_list` resource has an automatic state migration for `prefix` and `prefix_range` fields
   - Run `terraform plan` after upgrading to verify no unexpected changes

2. **Re-import Recommended:**
   - If you previously imported `alkira_service_pan` and noticed missing `segment_options`, re-import the resource to populate all attributes

3. **No Breaking Changes:** This is a fully backward-compatible patch release.
