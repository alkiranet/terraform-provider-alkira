---
page_title: "alkira_service_f5_lb Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  F5 Load Balancer Service. (BETA)
---

# alkira_service_f5_lb (Resource)

F5 Load Balancer Service. (**BETA**)

F5 Load Balancer Service supports two `license_type`.

`license_type`: `BRING_YOUR_OWN`
```terraform
resource "alkira_service_f5_lb" "example-lb" {
  name                = "example-lb"
  description         = "example-lb description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  instance {
    deployment_type     = "LTM_DNS"
    hostname_fqdn       = "examplelb.hostname"
    license_type        = "BRING_YOUR_OWN"
    name                = "example-lb-instance-1"
    version             = "17.1.1.1-0.0.2"
    deployment_option   = "TWO_BOOT_LOCATION"
    f5_registration_key = "key"
    f5_username         = "admin"
    f5_password         = "verysecretpassword"

  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  service_group_name = "example-service-group"
  size               = "LARGE"
}
```
or `license_type`: `PAY_AS_YOU_GO`
```terraform
resource "alkira_service_f5_lb" "example-lb-4" {
  name                = "example-lb-1"
  description         = "example-lb-1 description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  prefix_list_id      = alkira_list_prefix_list.example-prefix-list.id
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.4"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-4-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  service_group_name = "example-service-group-4"
  size               = "2LARGE"
}
``` 
 User can add multiple `instances` 
 ```terraform
resource "alkira_service_f5_lb" "example-lb-1" {
  name                = "example-lb-1"
  description         = "example-lb-1 description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  prefix_list_id      = alkira_list_prefix_list.example-prefix-list.id
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.1"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-1-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.1"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-1-instance-2"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  service_group_name = "example-service-group-1"
  size               = "2LARGE"
}
```
 User can also add configure multiple segments with `segment_options`
 ```terraform
resource "alkira_service_f5_lb" "example-lb-2" {
  name                = "example-lb-2"
  description         = "example-lb-2 description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  prefix_list_id      = alkira_list_prefix_list.example-prefix-list.id
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.2"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-2-instance-1"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  instance {
    deployment_type = "GOOD"
    hostname_fqdn   = "examplelb.hostname.2"
    license_type    = "PAY_AS_YOU_GO"
    name            = "example-lb-2-instance-2"
    version         = "17.1.1.1-0.0.2"
    f5_password     = "passwordispassword"
    f5_username     = "admin"

  }
  segment_ids = [alkira_segment.example-segment.id, alkira_segment.example-segment-1.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
  }
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment-1.id
  }
  service_group_name = "example-service-group-2"
  size               = "2LARGE"
}
```
 User can also configure BGP with `bgp_options_advertise_to_cxp_prefix_list_id` and `availability_zone`
 ```terraform
resource "alkira_service_f5_lb" "example-lb" {
  name                = "example-lb"
  description         = "example-lb description."
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.example-global-cidr.id
  instance {
    deployment_type     = "LTM_DNS"
    hostname_fqdn       = "examplelb.hostname"
    license_type        = "BRING_YOUR_OWN"
    name                = "example-lb-instance-1"
    version             = "17.1.1.1-0.0.2"
    deployment_option   = "TWO_BOOT_LOCATION"
    f5_registration_key = "key"
    f5_username         = "admin"
    f5_password         = "verysecretpassword"
    availability_zone   = 0
  }
  segment_ids = [alkira_segment.example-segment.id]
  segment_options {
    elb_nic_count = 2
    segment_id    = alkira_segment.example-segment.id
    bgp_options_advertise_to_cxp_prefix_list_id = alkira_policy_prefix_list.example.id
  }
  service_group_name = "example-service-group"
  size               = "LARGE"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) CXP on which the service should be provisioned.
- `global_cidr_list_id` (Number) ID of global CIDR list from which subnets will be allocated for the external network interfaces of instances. These interfaces host the public IP addresses needed for virtual IPs.
- `instance` (Block List, Min: 1) An array containing the properties for each F5 load balancer instance. (see [below for nested schema](#nestedblock--instance))
- `name` (String) Name of the service.
- `segment_ids` (Set of String) IDs of segments associated with the service.
- `segment_options` (Block Set, Min: 1) The segment options as used by your F5 Load Balancer. (see [below for nested schema](#nestedblock--segment_options))
- `service_group_name` (String) Name of the service group to be associated with the service.
- `size` (String) Size of the service, one of `SMALL`, `MEDIUM`, `LARGE` `2LARGE`, `5LARGE`.

### Optional

- `billing_tag_ids` (Set of Number) IDs of billing tags to associate with the service.
- `description` (String) Description of the service.
- `prefix_list_id` (Number) ID of prefix list to use for IP allowlist

### Read-Only

- `id` (String) The ID of this resource.
- `implicit_group_id` (Number) The ID of implicit group automaticaly created with the connector.
- `provision_state` (String) The provisioning state of the resource.

<a id="nestedblock--instance"></a>
### Nested Schema for `instance`

Required:

- `deployment_type` (String) The deployment type used for the F5 load balancer instance.The value could be one of `GOOD`, `BETTER`, `BEST` or `LTM_DNS`. Type `GOOD`, `BETTER` and `BEST` is only applicable when `license_type` is `PAY_AS_YOU_GO`. `LTM_DNS` is only applicable when `license_type` `BRING_YOUR_OWN`.
- `hostname_fqdn` (String) The FQDN defined in route 53.
- `license_type` (String) The type of license used for the F5 load balancer instance. Can be one of `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`
- `name` (String) Name of the F5 load balancer instance.
- `version` (String) The version of the F5 load balancer.

Optional:

- `availability_zone` (Number) Availability Zone of F5 Instance. Only used when bgp_options_advertise_to_cxp_prefix_list_id is provided
- `credential_id` (String) ID of the F5 load balancer credential. If the `credential_id` is not passed, `f5_username` and `f5_password` is required to create new credentials.
- `f5_password` (String, Sensitive) Password for the F5 load balancer. This can also be set by `ALKIRA_F5_PASSWORD` environment variable.
- `f5_registration_key` (String, Sensitive) Registration key for the F5 load balancer. Only required if `license_type` is `BRING_YOUR_OWN`. This can also be set by `ALKIRA_F5_REGISTRATION_KEY` environment variable.
- `f5_username` (String, Sensitive) Username for the F5 load balancer. Username is `admin` for AWS CXP and `akadmin`  for Azure CXP any other value will be rejected. This can also be set by `ALKIRA_F5_USERNAME` environment variable.
- `registration_credential_id` (String) ID of the F5 load balancer registration credential. If the `registration_credential_id` is not passed, `f5_registration_key` is required to create new credentials. Only required if `license_type` is `BRING_YOUR_OWN`.

Read-Only:

- `id` (Number) ID of the F5 load balancer instance.


<a id="nestedblock--segment_options"></a>
### Nested Schema for `segment_options`

Required:

- `elb_nic_count` (Number) Number of NICs to allocate for the segment.
- `segment_id` (String) ID of the segment.

Optional:

- `bgp_options_advertise_to_cxp_prefix_list_id` (Number) ID of prefix list used to advertise prefixes from F5 Load Balancer

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_service_f5_lb.example SERVICE_ID
```
