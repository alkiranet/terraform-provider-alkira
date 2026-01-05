---
subcategory: "Release Notes"
page_title: "v1.4.2"
description: |-
    Release notes for v1.4.2
---

## Bug Fixes

- **segment_options**: Fixed issue where `groups` field sent `null` instead of empty array when not specified, causing API validation errors. This affected `alkira_service_pan`, `alkira_service_checkpoint`, and `alkira_service_fortinet` resources.
