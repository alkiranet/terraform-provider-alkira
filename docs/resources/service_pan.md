---
page_title: "alkira_service_pan Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Manage Palo Alto Firewall service.
  When panorama_enabled is set to true, pan_username and pan_password are required.
---

# alkira_service_pan (Resource)

Manage Palo Alto Firewall service.

When `panorama_enabled` is set to `true`, `pan_username` and `pan_password` are required.

## Credentials

The PAN service is backed by three separate entries in the Alkira credentials vault, all
managed by this resource:

- **PAN service credential** — `pan_username`, `pan_password`, `pan_license_key`. Updated
  in the vault when these fields change. (`pan_password` and `pan_license_key` are
  marked sensitive; `pan_username` is constrained to `admin` / `akadmin` and is not
  considered secret.)
- **PAN registration credential** — `registration_pin_id`, `registration_pin_value`,
  `registration_pin_expiry`. Created with the resource; the provider does not update
  this credential on subsequent applies.
- **PAN master key credential** — `master_key`, `master_key_expiry`, gated by
  `master_key_enabled`. Created with the resource; the provider does not update this
  credential on subsequent applies.
- **Per-instance bootstrap credential** — `instance.auth_key`, `instance.auth_code`,
  `instance.auth_expiry` (one per `instance` block). Used by Panorama to authenticate
  the firewall VM during license activation; consumed during PAN-OS bootstrap and not
  returned by the API thereafter.

`pan_username` is constrained to `admin` or `akadmin` (backend constraint), and the
per-CXP rule is enforced at plan time:

- Azure / Azure China CXPs require `pan_username = "akadmin"`.
- AWS / AWS China / GCP CXPs require `pan_username = "admin"`.
- Other CXP types (OCI, on-prem) accept either value.

### Credential update semantics

All three credential blocks above are **consumed at instance bootstrap only**. The
Alkira backend pushes their values into the PAN VM during initial provisioning (via the
PAN-bootstrap Ansible playbook); there is no in-place rotation flow.

What this means in practice:

- Editing `pan_username`, `pan_password`, or `pan_license_key` and running
  `terraform apply` rotates the value in the Alkira credentials vault. It does **not**
  reconfigure already-provisioned PAN devices. New autoscale-out instances will pick up
  the latest value at their bootstrap; existing instances keep their original
  credentials.
- Editing `master_key`, `master_key_expiry`, `registration_pin_id`,
  `registration_pin_value`, `registration_pin_expiry`, or any per-instance
  `auth_key` / `auth_code` / `auth_expiry` and running `terraform apply` updates
  Terraform state but does not touch the vault or the deployment. The provider emits
  a `Warning` diagnostic on `apply` so the no-op is visible. (The warning fires at
  apply time, not at plan — `terraform plan` will still report a normal-looking
  diff before the no-op is detected.) To rotate any of these on a running deployment,
  destroy and recreate the resource (or rotate out-of-band on the PAN device itself,
  accepting that the vault and the device will then disagree).

This is an architectural property of the Alkira platform, not a provider limitation.

> **State file on disk is plaintext.** `Sensitive: true` only redacts CLI output. The
> values above are still stored verbatim in `terraform.tfstate`. Use a backend with
> encryption at rest and restrict access; do not commit state to source control.

## Import Behavior

Credential field values are write-only — the Alkira API never returns them. After
`terraform import`, `terraform plan` shows a diff for every credential field because
state has empty strings while configuration has the user-supplied values. Run
`terraform apply` once to reconcile state:

- The PAN service credential fields (`pan_username`, `pan_password`, `pan_license_key`)
  are pushed back to the vault on apply (see "Credential update semantics" — vault only,
  not the running device).
- The registration and master-key credential fields are recorded in Terraform state but
  no vault update is issued by the provider.

Subsequent plans show no changes. The imported state for the vault-only and
bootstrap-only fields is the configuration value, not whatever is actually in the vault
or on the deployed PAN device — there is no read path that could verify them.

## Example Usage

```terraform
variable "pan_password" {
  type      = string
  sensitive = true
}

variable "pan_registration_pin_id" {
  type      = string
  sensitive = true
}

variable "pan_registration_pin_value" {
  type      = string
  sensitive = true
}

variable "pan_master_key" {
  type      = string
  sensitive = true
}

variable "pan_instance_auth_key" {
  type      = string
  sensitive = true
}

variable "pan_instance_auth_code" {
  type      = string
  sensitive = true
}

resource "alkira_service_pan" "test1" {
  name                   = "PanFwTest"
  bundle                 = "VM_SERIES_BUNDLE_1"
  cxp                    = "US-WEST"
  global_protect_enabled = false
  license_type           = "PAY_AS_YOU_GO"
  max_instance_count     = 1
  segment_ids            = [alkira_segment.test1.id]
  management_segment_id  = alkira_segment.test1.id
  size                   = "SMALL"
  type                   = "VM-300"
  version                = "9.1.3"

  panorama_enabled      = true
  panorama_device_group = "alkira-test"
  panorama_ip_addresses = ["172.16.0.8"]
  panorama_template     = "test"

  # PAN Panorama credentials. pan_username must be "admin" (AWS/GCP) or
  # "akadmin" (Azure) per backend constraint.
  pan_password = var.pan_password
  pan_username = "admin"

  registration_pin_id     = var.pan_registration_pin_id
  registration_pin_value  = var.pan_registration_pin_value
  registration_pin_expiry = "2030-07-30"

  master_key_enabled = true
  master_key         = var.pan_master_key
  master_key_expiry  = "2030-08-01"

  global_protect_segment_options {
    segment_id            = (alkira_segment.test1.id)
    remote_user_zone_name = "RandomZoneName"
    portal_fqdn_prefix    = "randomprefix"
    service_group_name    = "RandomServiceGroupName"
  }

  # You can add more instance blocks. Make sure to change "max_instance_count".
  instance {
    name      = "tf-pan-instance-1"
    auth_key  = var.pan_instance_auth_key
    auth_code = var.pan_instance_auth_code
    global_protect_segment_options {
      segment_id      = (alkira_segment.test1.id)
      portal_enabled  = true
      gateway_enabled = true
      prefix_list_id  = alkira_policy_prefix_list.tf_prefix_list.id
    }
  }

  segment_options {
    segment_id = alkira_segment.segment.id
    zone_name  = "DEFAULT"
    groups     = [alkira_group.group.name, alkira_group.group1.name]
  }

  segment_options {
    segment_id = alkira_segment.segment1.id
    zone_name  = "zonename1"
    groups     = [alkira_group.group2.name]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the service should be provisioned.
- `instance` (Block List, Min: 1) (see [below for nested schema](#nestedblock--instance))
- `license_type` (String) PAN license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.
- `management_segment_id` (Number) Management Segment ID.
- `max_instance_count` (Number) Max number of Panorama instances for auto scale. Note: For Azure CXPs, this must equal `min_instance_count` as Azure does not support AutoScale.
- `name` (String) Name of the PAN service.
- `pan_password` (String, Sensitive) PAN Panorama password.
- `pan_username` (String) PAN Panorama username. For AWS and GCP, username must be `admin`. For Azure, it must be `akadmin`. The per-CXP rule is enforced at plan time via `CustomizeDiff`.
- `registration_pin_id` (String, Sensitive) PAN Registration PIN ID.
- `registration_pin_value` (String, Sensitive) PAN Registration PIN Value.
- `segment_ids` (Set of Number) IDs of segments associated with the service.
- `size` (String) The size of the service, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`.
- `version` (String) The version of the PAN firewall. Please check Alkira Portal for all supported versions.

### Optional

- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `bundle` (String) The software image bundle that would be used forPAN instance deployment. This is applicable for licenseType`PAY_AS_YOU_GO` only. If not provided, the default`PAN_VM_300_BUNDLE_2` would be used. However `PAN_VM_300_BUNDLE_2`is legacy bundle and is not supported on AWS. It is recommendedto use `VM_SERIES_BUNDLE_1` and `VM_SERIES_BUNDLE_2` (supports Global Protect).
- `description` (String) The description of the service.
- `global_protect_enabled` (Boolean) Enable global protect option or not. Default is `false`
- `global_protect_segment_options` (Block Set) Segment options for segments that are already associated with the service. Options should apply. If `global_protect_enabled` is set to false, `global_protect_segment_options` shound not be included in your request. (see [below for nested schema](#nestedblock--global_protect_segment_options))
- `license_sub_type` (String) PAN sub license type, either `CREDIT_BASED` or `MODEL_BASED`. (BETA)
- `master_key` (String, Sensitive) Master Key for PAN instances. Consumed at instance bootstrap only; see resource overview for credential update semantics.
- `master_key_enabled` (Boolean) Enable Master Key for PAN instances or not. It's default to `false`.
- `master_key_expiry` (String) PAN Master Key Expiry. The date should be in format of `YYYY-MM-DD`, e.g. `2000-01-01`.
- `min_instance_count` (Number) Minimal number of Panorama instances for auto scale. Default value is `0`. Note: For Azure CXPs, this must equal `max_instance_count` as Azure does not support AutoScale.
- `pan_license_key` (String, Sensitive) PAN Licensing API Key.
- `panorama_device_group` (String) Panorama device group.
- `panorama_enabled` (Boolean) Enable Panorama or not. Default value is `false`.
- `panorama_ip_addresses` (List of String) Panorama IP addresses.
- `panorama_template` (String) Panorama Template or Panorama Template Stack.
- `registration_pin_expiry` (String) PAN Registration PIN Expiry. The date should be in format of `YYYY-MM-DD`, e.g. `2000-01-01`.
- `segment_options` (Block Set) The segment options as used by your PAN firewall. (see [below for nested schema](#nestedblock--segment_options))
- `tunnel_protocol` (String) Tunnel Protocol, default to `IPSEC`, could be either `IPSEC` or `GRE`.
- `type` (String) The type of the PAN firewall. Either 'VM-300', 'VM-500' or 'VM-700'

### Read-Only

- `id` (String) The ID of this resource.
- `pan_credential_id` (String) ID of PAN credential.
- `pan_credential_name` (String) Name of PAN credential.
- `pan_master_key_credential_id` (String) ID of PAN master key credential.
- `pan_registration_credential_id` (String) ID of PAN Registration credential.
- `provision_state` (String) The provision state of the service.

<a id="nestedblock--instance"></a>
### Nested Schema for `instance`

Optional:

- `auth_code` (String, Sensitive) PAN instance auth code. Only required when `license_type` is `BRING_YOUR_OWN`.
- `auth_expiry` (String) PAN Auth Expiry. The date should be in format of `YYYY-MM-DD`, e.g. `2000-01-01`.
- `auth_key` (String, Sensitive) PAN instance auth key (VM-series bootstrap auth key). This is only required when `panorama_enabled` is set to `true`. **IMPORTANT:** The auth key MUST be generated from the Panorama CLI only. Auth keys generated using the Panorama web interface are NOT supported by Alkira and may cause provisioning to fail.
- `enable_traffic` (Boolean) Enable traffic on the PAN instance. Default value is `true`.
- `global_protect_segment_options` (Block Set) These options should be set only when global protect is enabled on service. These are set per segment. It is expected that on a segment where global protect is enabled at least 1 instance should be set with portal_enabled and at least one with gateway_enabled. It can be on the same instance or a different instance under the segment. (see [below for nested schema](#nestedblock--instance--global_protect_segment_options))
- `name` (String) The name of the PAN instance.

Read-Only:

- `credential_id` (String) ID of PAN instance credential.
- `id` (Number) The ID of the PAN instance.

<a id="nestedblock--instance--global_protect_segment_options"></a>
### Nested Schema for `instance.global_protect_segment_options`

Required:

- `gateway_enabled` (Boolean) indicates if the Global Protect Gateway is enabled on this PAN instance
- `portal_enabled` (Boolean) indicates if the GlobalProtect Portal is enabled on this PAN instance
- `prefix_list_id` (Number) Prefix List with Client IP Pool.
- `segment_id` (String) The segment ID for Global Protect options.



<a id="nestedblock--global_protect_segment_options"></a>
### Nested Schema for `global_protect_segment_options`

Required:

- `portal_fqdn_prefix` (String) Prefix for the global protect portal FQDN, this would be prepended to customer specific alkira domain For Example: if prefix is abc and tenant name is example then the FQDN would be abc.example.gpportal.alkira.com
- `remote_user_zone_name` (String) Firewall security zone is created using the zone name for remote user sessions.
- `segment_id` (String) The name of the segment to which the global protect options should apply
- `service_group_name` (String) The name of the service group. A group with the same name will be created.


<a id="nestedblock--segment_options"></a>
### Nested Schema for `segment_options`

Required:

- `segment_id` (String) The ID of the segment.
- `zone_name` (String) The name of the associated firewall zone.

Optional:

- `groups` (List of String) The list of groups associated with the zone.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_service_pan.example SERVICE_ID
```
