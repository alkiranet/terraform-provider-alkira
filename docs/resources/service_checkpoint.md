---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alkira_service_checkpoint Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage checkpoint services
---

# alkira_service_checkpoint (Resource)

Manage checkpoint services

## Example Usage

```terraform
resource "alkira_service_checkpoint" "test1" {
  auto_scale         = "OFF"
  cxp                = "US-WEST"
  credential_id      = alkira_credential_checkpoint.tf_test_checkpoint.id
  license_type       = "PAY_AS_YOU_GO"
  max_instance_count = 2
  min_instance_count = 2
  name               = "testname"
  segment_names      = [alkira_segment.test-seg-1.name]
  size               = "LARGE"
  tunnel_protocol    = "IPSEC"
  version            = "R80.30"

  segment_options {
    segment_id = alkira_segment.test-seg-1.id
    zone_name  = "DEFAULT"
    groups     = [alkira_group.test.name]
  }

  management_server {
    configuration_mode  = "MANUAL"
    global_cidr_list_id = 22
    ips                 = ["10.2.0.3"]
    reachability        = "PRIVATE"
    segment_id          = alkira_segment.test-seg-1.id
    type                = "SMS"
    user_name           = "admin"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **credential_id** (String) ID of Checkpoint Firewall credential managed by credential resource.
- **cxp** (String) CXP region.
- **license_type** (String) Checkpoint license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.
- **management_server** (Block Set, Min: 1) (see [below for nested schema](#nestedblock--management_server))
- **max_instance_count** (Number) The maximum number of Checkpoint Firewall instances that should be deployed when auto-scale is enabled. Note that auto-scale is not supported with Checkpoint at this time. `max_instance_count` must be greater than or equal to `min_instance_count`.
- **name** (String) Name of the Checkpoint Firewall service.
- **segment_names** (List of String) The names of the segments associated with the service.
- **segment_options** (Block Set, Min: 1) The segment options as used by your checkpoint firewall. (see [below for nested schema](#nestedblock--segment_options))
- **size** (String) The size of the service, one of `SMALL`, `MEDIUM`, `LARGE`.
- **version** (String) The version of the Checkpoint Firewall.

### Optional

- **auto_scale** (String) Indicate if `auto_scale` should be enabled for your checkpointfirewall. `ON` and `OFF` are accepted values. `OFF` is the default if field is omitted
- **billing_tag_ids** (List of Number) Billing tag IDs to associate with the service.
- **description** (String) The description of the checkpoint service.
- **id** (String) The ID of this resource.
- **instances** (Block Set) An array containing properties for each Checkpoint Firewall instance that needs to be deployed. The number of instances should be equal to `max_instance_count`. (see [below for nested schema](#nestedblock--instances))
- **min_instance_count** (Number) The minimum number of Checkpoint Firewall instances that should be deployed at any point in time.
- **pdp_ips** (List of String) The IPs of the PDP Brokers.
- **tunnel_protocol** (String) Tunnel Protocol, default to `IPSEC`, could be either `IPSEC` or `GRE`.

<a id="nestedblock--management_server"></a>
### Nested Schema for `management_server`

Required:

- **configuration_mode** (String) The configuration_mode specifies whether the firewall is to be automatically configured by Alkira or not. To automatically configure the firewall Alkira needs access to the CheckPoint management server. If you choose to use manual configuration Alkira will provide the customer information about the checkpoint instances so that you can manually configure the firewall.
- **global_cidr_list_id** (Number) The ID of the global cidr list to be associated with the management server.
- **ips** (List of String) Management server IPs.

Optional:

- **domain** (String) Management server domain.
- **reachability** (String) This option specifies whether the management server is publicly reachable or not. If the reachability is private then you need to provide the segment to be used to access the management server. Default value is `PUBLIC`.
- **segment_id** (Number) The ID of the segment to be used to access the management server.
- **type** (String) The type of the management server.
- **user_name** (String) The user_name of the management server.


<a id="nestedblock--segment_options"></a>
### Nested Schema for `segment_options`

Required:

- **groups** (List of String) The list of Groups associated with the zone.
- **segment_id** (Number) The ID of the segment.
- **zone_name** (String) The name of the associated zone.


<a id="nestedblock--instances"></a>
### Nested Schema for `instances`

Required:

- **name** (String) The name of the Checkpoint Firewall instance.

