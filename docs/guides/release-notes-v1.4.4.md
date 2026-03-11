---
subcategory: "Release Notes"
page_title: "v1.4.4"
description: |-
    Release notes for v1.4.4
---

# Alkira Terraform Provider v1.4.4 Release Notes

Release Date: 2026-03-11

## Overview

Version 1.4.4 is a bug fix release focused on improving `terraform import` reliability, credential management, and provisioning error visibility.

---

## Bug Fixes

### Import Fixes

- **Policy Rule:** Fixed a `ConflictsWith` error that occurred during `terraform import` for `alkira_policy_rule` resources. Fields `src_ip`, `dst_ip`, `src_prefix_list_id`, `dst_prefix_list_id`, and `internet_application_id` are now only set when non-empty, preventing conflicts with mutually exclusive field definitions.

### Credential Management

- **Aruba Edge:** Fixed an issue where `credentialId` was incorrectly reset when credential fields were updated. The credential ID is now preserved when already set, and instance IDs are correctly saved to state after apply.

### Error Messages

- **Provisioning:** Provisioning failure errors now include a detailed reason when available, making it easier to diagnose configuration issues without contacting support.

---

## Upgrade Instructions

### From v1.4.3 to v1.4.4

1. **No Breaking Changes:** This is a fully backward-compatible patch release.

2. **Re-import Recommended:**
   - If you previously imported `alkira_policy_rule` and encountered `ConflictsWith` errors, re-import the resource to populate all attributes correctly.
