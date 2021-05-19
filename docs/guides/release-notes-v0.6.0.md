---
subcategory: "Release Notes"
page_title: "v0.6.0"
description: |-
    Release notes for v0.6.0
---

Major release focusing on IPSec connector.

### resource `alkira_connector_ipsec`

Big upgrade with support of various new features:

* Argument `vpn_mode` is required now for defining IPSec connector.

* New `policy_options` is added to support policy based routing if
  `vpn_mode` is set to `POLICY_BASED`.

* New `routing_options` is added to support route based routing if
  `vpn_mode` is set to `ROUTE_BASED`.

* Better validation was added to validate arguments.

* There should be only one `segment_options` block allowed, since
  IPSec connector only supports `1` segment for now.

* `advanced` of `endpoint` arguments are added for testing for now.


### resource `alkira_billing_tags`

Added new optional argument `description` that could be used to
describe the billing tag.
