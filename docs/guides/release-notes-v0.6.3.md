---
subcategory: "Release Notes"
page_title: "v0.6.3"
description: |-
    Release notes for v0.6.3
---

This release fixes the following problems.

* Support `HTTP_PROXY` or `HTTPS_PROXY` for initializing provider.
* Fix update functionality for all credential resources.


### resource `alkira_service_pan`

* Add new variable `tunnel_protocol`.
* Update documentation and examples.

* Fix variable `billing_tags` to `billing_tag_ids` to be clearer.
* Fix variable `segment` to `segment_ids` to be clearer.
* Fix variable `management_segment` to `management_segment_id` to be clearer.


### resource `alkira_connector_internet`

* Update documentation.
