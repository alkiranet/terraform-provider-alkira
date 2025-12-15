# Alkira Terraform Provider v1.3.3 Release Notes

Release Date: 2025-12-15

## Overview

Version 1.3.3 introduces three new BETA resources, BETA enhancements for F5 BGP and AWS TGW Direct Connect Gateway support, performance enhancements, bug fixes, and documentation improvements.

---

## 🎉 New Features (BETA)

### New Resources

#### 1. BlueCat DNS Service (`alkira_service_bluecat`) - BETA

Full-featured BlueCat DNS service integration with support for both BDDS and EDGE instances.

**Key Features:**
- BDDS (BlueCat DNS/DHCP Server) instance support
- EDGE instance support for distributed deployments
- Anycast configuration for high availability
- Flexible deployment options: BDDS-only, EDGE-only, or hybrid
- Comprehensive documentation with 5 deployment examples

**Example:**
```hcl
resource "alkira_service_bluecat" "example" {
  name                = "bluecat-service"
  cxp                 = "US-WEST"
  global_cidr_list_id = alkira_list_global_cidr.dns.id
  license_type        = "BRING_YOUR_OWN"
  segment_ids         = [alkira_segment.default.id]
  service_group_name  = "dns-services"

  instance {
    name = "bdds-primary"
    type = "BDDS"

    bdds_options {
      hostname       = "bdds.example.com"
      model          = "cBDDS50"
      version        = "9.4.0"
      client_id      = "client-001"
      activation_key = "YOUR_ACTIVATION_KEY"
    }
  }
}
```

**Documentation:** `docs/resources/service_bluecat.md`

---

#### 2. Juniper SD-WAN Connector (`alkira_connector_juniper_sdwan`) - BETA

Native support for Juniper SD-WAN (Session Smart Router) connectivity.

**Key Features:**
- Juniper SSR version management
- VRF mapping with Alkira segment integration
- GRE tunnel protocol support
- Registration key-based instance authentication
- Internet exit control via `advertise_default_route`

**Example:**
```hcl
resource "alkira_connector_juniper_sdwan" "example" {
  name                = "juniper-connector"
  cxp                 = "US-EAST"
  size                = "SMALL"
  juniper_ssr_version = "6.3.4"
  availability_zone   = 0

  instance {
    hostname         = "juniper-host"
    registration_key = "YOUR_REGISTRATION_KEY"
  }

  juniper_ssr_vrf_mapping {
    segment_id              = alkira_segment.production.id
    advertise_default_route = false
  }
}
```

**Documentation:** `docs/resources/connector_juniper_sdwan.md`

---

#### 3. Network Entity Scale Options (`alkira_network_entity_scale_options`) - BETA

Fine-grained control over connector and service scaling for enhanced performance.

**Key Features:**
- Configurable scale options for any connector or service
- Additional tunnels per node configuration
- Additional nodes configuration
- Segment-based and zone-based scaling
- Prevents race conditions with proper dependency tracking

**Example:**
```hcl
resource "alkira_network_entity_scale_options" "example" {
  name        = "scale-options"
  entity_id   = alkira_service_fortinet.example.id
  entity_type = "SERVICE"

  segment_scale_options {
    additional_tunnels_per_node = 2
    segment_id                  = alkira_segment.production.id
    zone_name                   = "ZoneA"
  }
}
```

**Documentation:** `docs/resources/network_entity_scale_options.md`

---

## ✨ Feature Enhancements (BETA)

### F5 Load Balancer Service - BGP Support (BETA)

- **BGP Support:** Added BGP routing capabilities for F5 services
- **Error Handling:** Improved error handling and reporting for vServer update operations

**Note:** F5 BGP support is currently in BETA. Test thoroughly before production use.

### AWS Transit Gateway Peering - Direct Connect Gateway Support (BETA)

**Direct Connect Gateway Association Proposals** - BETA

Added comprehensive support for AWS Direct Connect Gateway association proposals with automated state management.

**New Fields:**
- `direct_connect_gateway_association_proposal_id` - Proposal identifier
- `direct_connect_gateway_association_proposal_state` - Current proposal state
- `direct_connect_gateway_association_proposal_created_at` - Creation timestamp
- `direct_connect_gateway_association_proposal_updated_at` - Last update timestamp

**Documentation:** `docs/resources/peering_gateway_aws_tgw_attachment.md`

**Note:** Direct Connect Gateway association proposal support is currently in BETA. Test thoroughly before production use.

---

## ⚡ Performance Enhancements

Improved provider performance and reliability with new configuration options:

```hcl
provider "alkira" {
  portal   = "tenant.portal.alkira.com"
  api_key  = "YOUR_API_KEY"

  # Performance options
  serialization_enabled  = true    # Improved API call handling
  serialization_timeout  = 300     # Timeout in seconds
}
```

**Key Improvements:**
- Enhanced API call handling for better concurrency control
- Improved error diagnostics and reporting
- Configurable timeout controls

**Documentation:** `docs/index.md`

---

## 🐛 Bug Fixes

### Critical Fixes

- **Azure VNet Connector:** Made `subnet_cidr` required to prevent provisioning failures
- **Cisco FTDv:** Fixed slice assignment bug preventing runtime panics

### Configuration Fixes

- **IPSec Connectors:** Made `bgpAuthKeyAlkira` optional in standard and advanced connectors
- **Service Zone Groups:** Made `groups` optional in Checkpoint, Fortinet, and PAN service configurations
- **Size Validations:** Removed restrictive size validations to support new instance sizes

### Azure Enhancements

- **ExpressRoute Connector:** Added `implicit_group_id` computed field
- **PAN Instance Validation:** Added validation ensuring min/max instance counts are equal for Azure (no autoscale support)

---

## ⚠️ Breaking Changes

### Azure VNet Connector - Required Field

**Change:** `subnet_cidr` is now a required field in `alkira_connector_azure_vnet`

**Impact:** Configurations missing `subnet_cidr` will fail validation

**Migration:**
```hcl
resource "alkira_connector_azure_vnet" "example" {
  # ... other fields ...

  # REQUIRED - must be specified
  subnet_cidr = "10.0.1.0/24"
}
```

---

## 📚 Documentation Improvements

### Enhanced Documentation

- **Project README:** Comprehensive setup, development, and contribution guidelines
- **Group Resource:** Enhanced with multiple examples including service groups and policy integration
- **IPSec Connector:** Expanded with 8+ configuration examples:
  - BGP authentication
  - HA configurations
  - Multi-site deployments
  - Policy-based routing
  - Advanced crypto options
  - Dynamic gateway configurations

### PAN Service Documentation Updates

- **auth_key Clarification:** Must be generated from Panorama CLI (not web interface)
- **panorama_template:** Clarified support for both Template and Template Stack
- **pan_license_key:** Updated description to match UI terminology

### Other Updates

- Azure ExpressRoute documentation updated with `implicit_group_id` field
- All new resources include comprehensive examples and usage notes

---

## 🏗️ Infrastructure & Development

- Added Slack notifications for failed CI/CD workflows
- Updated golangci-lint configuration with comprehensive rules
- Code quality improvements and linting fixes across codebase

---

## 📦 Upgrade Instructions

### From v1.3.2 to v1.3.3

1. **Review Azure VNet Connectors:**
   - Ensure all `alkira_connector_azure_vnet` resources have `subnet_cidr` specified
   - Add the field if missing

2. **Test BETA Features:**
   - All new resources are marked BETA
   - F5 BGP support is BETA
   - AWS TGW Direct Connect Gateway association proposals are BETA
   - Test thoroughly in non-production environments before production use

3. **Review Documentation:**
   - Check updated examples for IPSec, PAN, and Group resources
   - Review new provider configuration options

---

## 🙏 Contributors

Thank you to all contributors who made this release possible!

---

**Note:** Features marked as BETA are production-ready but may undergo changes based on user feedback and evolving requirements. Please report any issues or feedback via GitHub Issues.
