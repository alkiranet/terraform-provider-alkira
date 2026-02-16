---
page_title: "alkira_service_bluecat Resource - terraform-provider-alkira"
subcategory: ""
description: |-
  Provide Bluecat service resource (BETA).
---

# alkira_service_bluecat (Resource)

Provide Bluecat service resource (**BETA**).


## ANYCAST

When Anycast is configured for Bluecat services, Alkira automatically 
generates system-managed routing policies(`alkira_policy_routes`) and 
prefix lists(`alkira_policy_prefix_list`). These resources are essential 
for proper Anycast functionality and have specific lifecycle and 
scope characteristics. These system generated resources use the prefix 
`ALK-SYSTEM-GENERATED-BLUECAT`.

These route policies and prefix lists cannot be deleted or modified
directly, their lifecycle is bound by the Bluecat services that are
configured on the network and their anycast configuration

Anycast configuration operates at the service level, not the instance level:
`bdds_anycast` is applied to all BDDS instances in the service and
`edge_anycast` is applied to all EDGE instances. If anycast is configured for 
BDDS instances it must be configured for EDGE instances also and vice versa.

## Instances

Bluecat services support flexible instance configurations.
There can be BDDS instances only services, 
EDGE instances only services or a
hybrid scenario with both BDDS and EDGE instances.


## Example Usage

The example uses resource `alkira_segment` and
`alkira_list_global_cidr`.

This is a minimal example with only BDDS instance and no anycast

```terraform
resource "alkira_service_bluecat" "minimal" {
  name                = "bluecat-minimal"
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.basic.id
  segment_ids         = [alkira_segment.default.id]
  service_group_name  = "dns-basic"

  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "basic-client"
      activation_key = "BASIC1234567890ABCDEF"
    }
  }
}
```

You can configure `bdds_anycast` for the above BDDS instance similar to the following example

```terraform
resource "alkira_service_bluecat" "bdds_only" {
  name                = "bluecat-bdds-only"
  cxp                 = "US-WEST"
  description         = "Bluecat service with BDDS instances only"
  global_cidr_list_id = alkira_list_global_cidr.dns_allowed.id
  segment_ids         = [alkira_segment.corp.id]
  service_group_name  = "dns-services"

  bdds_anycast {
    ips         = ["10.0.100.10"]
    backup_cxps = ["US-EAST"]
  }

  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-primary"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "bdds-client-001"
      activation_key = "ABCD1234EFGH5678IJKL9012"
    }
  }

  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-secondary"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "bdds-client-002"
      activation_key = "MNOP3456QRST7890UVWX1234"
    }
  }
}
```

Similar configuration can be created for a service with only EDGE instances.

```terraform
resource "alkira_service_bluecat" "edge_only" {
  name                = "bluecat-edge-only"
  cxp                 = "EU-CENTRAL"
  description         = "Bluecat service with Edge instances only"
  global_cidr_list_id = alkira_list_global_cidr.branch_dns.id
  segment_ids         = [alkira_segment.branch.id, alkira_segment.dmz.id]
  service_group_name  = "edge-dns-services"

  edge_anycast {
    ips         = ["172.16.50.10"]
    backup_cxps = ["US-WEST"]
  }

  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-branch-01"
      version     = "4.1.2"
      config_data = "CONFIG_DATA_STRING_BRANCH_01_ENCODED_BASE64"
    }
  }

  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-branch-02"
      version     = "4.1.2"
      config_data = "CONFIG_DATA_STRING_BRANCH_02_ENCODED_BASE64"
    }
  }

  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dmz"
      version     = "4.0.5"
      config_data = "CONFIG_DATA_STRING_DMZ_ENCODED_BASE64"
    }
  }
}
```

You can a create a hybrid configuration with both BDDS and EDGE instances with anycast enabled.

```terraform
resource "alkira_service_bluecat" "hybrid_deployment" {
  name                = "bluecat-hybrid"
  cxp                 = "US-EAST"
  description         = "Hybrid Bluecat deployment with both BDDS and Edge instances"
  global_cidr_list_id = alkira_list_global_cidr.global_dns.id
  segment_ids         = [alkira_segment.production.id]
  service_group_name  = "hybrid-dns-services"

  billing_tag_ids = [
    alkira_billing_tag.dns_infrastructure.id,
    alkira_billing_tag.production.id
  ]

  bdds_anycast {
    ips         = ["192.168.10.50"]
    backup_cxps = ["US-WEST"]
  }

  edge_anycast {
    ips         = ["192.168.20.50"]
    backup_cxps = ["US-WEST"]
  }

  # Core BDDS instances for centralized management
  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-core-01"
      model          = "cBDDS50"
      version        = "9.5.1"
      client_id      = "enterprise-core-001"
      activation_key = "CORE1234ABCD5678EFGH9012IJKL"
    }
  }

  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-core-02"
      model          = "cBDDS50"
      version        = "9.5.1"
      client_id      = "enterprise-core-002"
      activation_key = "CORE5678MNOP9012QRST3456UVWX"
    }
  }

  # Edge instances for distributed locations
  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dc-east"
      version     = "4.2.0"
      config_data = "EDGE_DC_EAST_CONFIG_BASE64_ENCODED_STRING"
    }
  }

  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-dc-west"
      version     = "4.2.0"
      config_data = "EDGE_DC_WEST_CONFIG_BASE64_ENCODED_STRING"
    }
  }
}
```

The following is a sample production deployment.

```terraform
resource "alkira_service_bluecat" "production" {
  name                = "bluecat-production"
  cxp                 = "ASIA-PACIFIC"
  description         = "Production Bluecat service for enterprise DNS"
  global_cidr_list_id = alkira_list_global_cidr.enterprise.id
  segment_ids         = [alkira_segment.production.id]
  service_group_name  = "enterprise-dns"

  billing_tag_ids = [
    alkira_billing_tag.networking.id,
    alkira_billing_tag.production.id,
    alkira_billing_tag.asia_pacific.id
  ]

  bdds_anycast {
    ips = ["10.100.1.10"]
    backup_cxps = ["US-WEST"]
  }

  edge_anycast {
    ips = ["10.200.1.10"]
    backup_cxps = ["US-WEST"]
  }

  # Primary BDDS for enterprise services
  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-ent-01"
      model          = "cBDDS50"
      version        = "9.5.2"
      client_id      = "enterprise-asia-001"
      activation_key = "ENT_ASIA_PRIMARY_KEY_2024_001"
    }
  }

  # Secondary BDDS for redundancy
  instance {
    type = "BDDS"
    
    bdds_options {
      hostname       = "bdds-ent-02"
      model          = "cBDDS50"
      version        = "9.5.2"
      client_id      = "enterprise-asia-002"
      activation_key = "ENT_ASIA_SECONDARY_KEY_2024_002"
    }
  }

  # Edge for distributed locations
  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-primary"
      version     = "4.2.1"
      config_data = "ASIA_PRIMARY_EDGE_CONFIG_BASE64"
    }
  }

  # Edge for backup services
  instance {
    type = "EDGE"
    
    edge_options {
      hostname    = "edge-backup"
      version     = "4.2.1"
      config_data = "ASIA_BACKUP_EDGE_CONFIG_BASE64"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cxp` (String) The CXP where the service should be provisioned.
- `global_cidr_list_id` (Number) The ID of the global cidr list to be associated with the Bluecat service.
- `instance` (Block List, Min: 1) The properties pertaining to each individual instance of the Bluecat service. (see [below for nested schema](#nestedblock--instance))
- `name` (String) Name of the Bluecat service.
- `segment_ids` (Set of String) IDs of segments associated with the service.
- `service_group_name` (String) The name of the service group to be associated with the service. A service group represents the service in traffic policies, route policies and when configuring segment resource shares.

### Optional

- `bdds_anycast` (Block Set) Defines the AnyCast configuration for BDDS type instances (see [below for nested schema](#nestedblock--bdds_anycast))
- `billing_tag_ids` (Set of Number) Billing tags to be associated with the resource. (see resource `alkira_billing_tag`).
- `description` (String) The description of the Bluecat service.
- `edge_anycast` (Block Set) Defines the AnyCast configuration for EDGE type instances. (see [below for nested schema](#nestedblock--edge_anycast))

### Read-Only

- `id` (String) The ID of this resource.
- `license_type` (String) Bluecat license type, only `BRING_YOUR_OWN` is supported right now.
- `provision_state` (String) The provision state of the resource.
- `service_group_id` (Number) The ID of the service group to be associated with the service. A service group represents the service in traffic policies, route policies and when configuring segment resource shares.
- `service_group_implicit_group_id` (Number) The ID of the implicit group to be associated with the service.

<a id="nestedblock--instance"></a>
### Nested Schema for `instance`

Required:

- `type` (String) The type of the Bluecat instance that is to be provisioned. The value could be `BDDS`, and `EDGE`.

Optional:

- `bdds_options` (Block List, Max: 1) Defines the options required when instance type is BDDS. bdds_options must be populated if type of instance is BDDS (see [below for nested schema](#nestedblock--instance--bdds_options))
- `edge_options` (Block List, Max: 1) Defines the options required when instance type is EDGE. edge_options must be populated if type of instance is EDGE (see [below for nested schema](#nestedblock--instance--edge_options))

Read-Only:

- `id` (Number) The ID of the Bluecat instance.
- `name` (String) The name of the Bluecat instance. This is set to hostname

<a id="nestedblock--instance--bdds_options"></a>
### Nested Schema for `instance.bdds_options`

Required:

- `activation_key` (String, Sensitive) The license activationKey of the Bluecat BDDS instance.
- `client_id` (String) The license clientId of the Bluecat BDDS instance.
- `hostname` (String) The host name of the instance.
- `model` (String) The model of the Bluecat BDDS instance.
- `version` (String) The version of the Bluecat BDDS instance to be used. Please check Alkira Portal for all supported versions

Read-Only:

- `license_credential_id` (String) The license credential ID of the BDDS instance.


<a id="nestedblock--instance--edge_options"></a>
### Nested Schema for `instance.edge_options`

Required:

- `config_data` (String) The Base64 encoded configuration data generated on Bluecat Edge portal.
- `hostname` (String) The host name of the Edge instance. This should match what was configured on the bluecat edge portal.
- `version` (String) The version of the Bluecat Edge instance to be used. Please check Alkira Portal for all supported versions

Read-Only:

- `credential_id` (String) The credential ID of the Edge instance.



<a id="nestedblock--bdds_anycast"></a>
### Nested Schema for `bdds_anycast`

Optional:

- `backup_cxps` (List of String) The `backup_cxps` to be used when the current Bluecat service is not available. It also needs to have a configured Bluecat service in order to take advantage of this feature. It is NOT required that the `backup_cxps` should have a configured Bluecat service before it can be designated as a backup.
- `ips` (List of String) The IPs to be used for AnyCast. The IPs used for AnyCast MUST NOT overlap the CIDR of `alkira_segment` resource associated with the service.


<a id="nestedblock--edge_anycast"></a>
### Nested Schema for `edge_anycast`

Optional:

- `backup_cxps` (List of String) The `backup_cxps` to be used when the current Bluecat service is not available. It also needs to have a configured Bluecat service in order to take advantage of this feature. It is NOT required that the `backup_cxps` should have a configured Bluecat service before it can be designated as a backup.
- `ips` (List of String) The IPs to be used for AnyCast. The IPs used for AnyCast MUST NOT overlap the CIDR of `alkira_segment` resource associated with the service.

## Import

Import is supported using the following syntax:

```shell
terraform import alkira_service_bluecat.example SERVICE_ID
```
