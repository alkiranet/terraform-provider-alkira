# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## AWS VPC connector — filter system-generated default routing from state (AK-65432)

The provider now reads routing back from the API into Terraform state on every refresh of `alkira_connector_aws_vpc`, and filters out backend-tagged auto-generated entries. The Alkira backend tags system-generated routing with `defaultRouting: true`; the provider drops these so they never enter state.

`ConflictsWith` between `vpc_cidr` and `vpc_subnet` removed. Both can now be specified together. The only backend constraint is that a SUBNET prefix in `vpc_subnet` cannot be a subset of a CIDR prefix in `vpc_cidr` on the same connector.

**Impact — most users:** none. New connectors created without routing in HCL no longer drift. New connectors with explicit routing work as before. `terraform import` populates routing in state correctly.

**Impact — legacy connectors** (created before backend `defaultRouting` tagging): routing pulls into state on first refresh. If HCL has no routing block, the next plan shows a diff for `vpc_cidr` / `vpc_subnet` / `vpc_route_table`. Add the routing block to HCL to clear it. The apply that follows is idempotent — backend generates no provisioning tasks.
