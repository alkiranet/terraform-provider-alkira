# terraform-provider-alkira

## 1. Primary Responsibility

The public Terraform provider (`alkiranet/alkira` on the Terraform Registry) that lets customers manage their entire Alkira Cloud Services Exchange (CSX) tenant as code. It translates HCL resource definitions into calls against the customer's Alkira portal API, and reconciles Terraform state with what actually exists in the tenant.

It exists because network-as-a-service customers overwhelmingly operate through IaC pipelines: they define connectors, segments, policies, and firewall services in Terraform alongside their AWS/Azure/GCP infrastructure. The provider is the **primary programmatic front door to Alkira for customers** — for many tenants, more configuration flows through Terraform than through the portal UI.

Business capability owned: customer-facing infrastructure-as-code lifecycle (create/read/update/delete + optional provisioning) for all tenant network resources.

## 2. Key Responsibilities

- Expose Alkira tenant resources (connectors, segments, policies, services, credentials, lists, groups, etc.) as Terraform resources and data sources.
- Translate HCL schemas to/from Alkira API JSON payloads, preserving customer state block ordering and computed fields.
- Authenticate to the customer's portal (API key preferred; username/password deprecated).
- Optionally drive **provisioning** as part of `terraform apply` (the `provision = true` provider flag), waiting for the tenant-network provision request to succeed before reporting the resource as created/updated.
- Detect drift between Terraform state and the live tenant on `terraform plan`/`refresh`.
- Ship generated documentation (tfplugindocs) and examples to the Terraform Registry.

## 3. Important Interfaces

The provider registers ~71 resources and ~42 data sources. Grouped by category rather than listing each:

| Category | Prefix / Examples | What it manages |
|---|---|---|
| Cloud connectors | `alkira_connector_aws_vpc`, `_azure_vnet`, `_gcp_vpc`, `_oci_vcn`, `_aws_tgw`, `_azure_vhub`, `_aws_dx`, `_gcp_interconnect`, `_azure_expressroute` | On-ramps from cloud provider networks into the CXP |
| SD-WAN / edge connectors | `_cisco_sdwan`, `_versa_sdwan`, `_vmware_sdwan`, `_fortinet_sdwan`, `_juniper_sdwan`, `_aruba_edge`, `_prisma_sdwan` (client-side) | Branch/SD-WAN fabric integration |
| Site / other connectors | `_ipsec`, `_ipsec_adv`, `_ipsec_tunnel_profile`, `_internet_exit`, `_remote_access`, `_akamai_prolexic` | Site VPN, internet breakout, ZTNA remote access, DDoS scrubbing |
| Network core | `alkira_segment`, `alkira_group`, `alkira_group_user`, `alkira_group_direct_inter_connector`, `alkira_internet_application`, `alkira_ip_reservation`, `alkira_byoip_prefix` | Segmentation, grouping, internet apps, addressing |
| Policy | `alkira_policy`, `_policy_rule`, `_policy_rule_list`, `_policy_prefix_list`, `_policy_nat`, `_policy_nat_rule`, `_policy_routing`, `_policy_inter_cxp_routing` | Traffic/NAT/routing policy objects |
| Lists | `alkira_list_as_path`, `_community`, `_extended_community`, `_global_cidr`, `_dns_server`, `_udr`, `_policy_fqdn` | Reusable match lists |
| Marketplace services | `alkira_service_pan`, `_checkpoint`, `_fortinet`, `_cisco_ftdv`, `_zscaler`, `_infoblox`, `_f5_lb`, `_f5_vserver_endpoint`, `_bluecat` | Firewall / DNS / LB service insertion |
| Credentials | `alkira_credential_aws_vpc`, `_azure_vnet`, `_gcp_vpc`, `_oci_vcn`, `_ssh_key_pair` | Cloud + service credentials stored in Alkira |
| Peering gateways | `alkira_peering_gateway_cxp`, `_aws_tgw`, `_aws_tgw_attachment`, `_azure_vnet_third_party_connector_attachment` | Cloud-native peering constructs |
| Observability / misc | `alkira_flow_collector`, `alkira_probe_http/https/tcp`, `alkira_billing_tag`, `alkira_cloudvisor_account`, `alkira_network_entity_scale_options`, `alkira_segment_resource(_share)` | Flow export, health probes, billing, CloudVisor, scaling |

Provider-level configuration (in `provider "alkira" {}`):
- `portal` — tenant portal hostname (`<tenant>.portal.alkira.com`)
- `api_key` (recommended) or `username`/`password` (deprecated)
- `provision` (bool, default false) — whether applies also trigger tenant provisioning
- `validation` (bool) — asynchronous validations
- `serialization_enabled` / `serialization_timeout` — serializes API calls through the client to avoid concurrency issues on the backend

## 4. Dependencies

**Incoming (who uses it):**
- Customers' Terraform CLI / Terraform Cloud / Atlantis pipelines (pulled from the public Terraform Registry `alkiranet/alkira`).
- Internal acceptance-test pipelines (GitHub Actions `acceptance-preprod` / `acceptance-prod` against live environments).
- Alkira SE/solution teams building reference architectures (`examples/`, `acceptance/`).

**Outgoing:**
- `alkira-client-go` — the only path to the backend. The provider itself contains no raw HTTP; every CRUD goes through `AlkiraClient` / `AlkiraAPI[T]`.
- Via the client, everything hits the customer's **portal API** (`https://<tenant>.portal.alkira.com/api/...`), which is served by the **Alkira Backend Services**.

```
customer HCL
    |
terraform CLI  --gRPC plugin protocol-->  terraform-provider-alkira
                                                 |
                                          alkira-client-go
                                                 |  HTTPS (api-key / basic auth)
                                <tenant>.portal.alkira.com/api
                                                 |
                                  Alkira Backend Services
```

## 5. Data Ownership

- **Owns no data.** The Alkira Backend Services are the source of truth for all resource configuration; Terraform state (customer-side, in their backend of choice) is a cache of it.
- Drift is resolved at `plan`/`refresh` time by reading live resources back through the client (`GetById` uses `includeMarkedForDeletion=true`).
- Resources carry a computed `provision_state` attribute (SUCCESS/FAILED/…) reflecting the last provisioning outcome — this is state *about* the backend, stored in the customer's tfstate.
- No databases, no Kafka, no persistent storage of any kind in this repo.

## 6. Deployment

- **Released to the public Terraform Registry** as `alkiranet/alkira`. Release = push a `v*` git tag → GitHub Actions `build-release.yaml` runs **goreleaser** (multi-OS/arch zips, SHA256SUMS, GPG-signed) → Registry picks up the GitHub release. Build failures notify a Slack channel.
- Semantic versioning; current line is v1.x (e.g. v1.5.x). Customers pin with `version = "~> 1.0"`.
- The provider pins a specific `alkira-client-go` version in `go.mod` (often a pseudo-version ahead of the last client tag) and **vendors** dependencies.
- CI: unit tests, lint (golangci), dependency review, and **live acceptance tests against preprod and prod** portals.
- Compatibility concerns: the backend API evolves continuously; an old provider against a new backend (or vice versa in preprod) can mis-serialize fields. Releases are gated by the `review-terraform-release-notes` process — backward-compatibility review of schema changes, ForceNew additions, defaults, and the client-go bump. **A published release cannot be un-shipped**; customers auto-upgrade via loose version constraints.

## 7. Failure Impact

**Severity: High** (customer-facing, but no data-plane impact).

- A bad release breaks customers' CI/CD pipelines: `terraform plan/apply` errors, spurious diffs, or — worst case — an unintended **ForceNew** that plans destroy/recreate of a production connector (traffic-affecting if applied).
- State corruption bugs (wrong type stored, missing fields) can strand customers with resources they can neither update nor cleanly delete.
- Because the provider is pulled from the public registry with `~>` constraints, a bad minor release propagates to customers immediately on their next `terraform init -upgrade`.
- Blast radius is limited to the control plane: existing tenant networks keep forwarding traffic. Impact is on the ability to *change* infrastructure, plus the destroy/recreate risk above.

## 8. Common Production Issues

| Issue | Observable symptom |
|---|---|
| **Schema drift / perpetual diff** | `terraform plan` shows changes on every run with no config edits — usually API returning fields in a different shape/order than stored, or a new backend field not handled in the `set` helpers. |
| **ForceNew surprise** | Plan shows `-/+ destroy and then create replacement` for a field the customer thought was updatable. Per repo convention, ForceNew is added *only* when the backend silently overwrites the field — mistakes here are release blockers. |
| **API incompatibility after backend release** | Apply fails with 4xx/unmarshal errors: backend added a required field, changed an enum, or renamed a property that the pinned client-go version doesn't know about. |
| **Provision-vs-apply mismatch** | With `provision = false`, apply succeeds but nothing changes on the network until someone provisions from the portal; customers report "terraform said OK but the connector is down." With `provision = true`, applies block up to 240 min waiting on the provision request and fail with `provision request <id> failed/timed out`. |
| **`provision_state = FAILED` in state** | Resource exists in Alkira config DB but provisioning failed; next plan shows the computed field flapping. |
| **Instance-ID type bugs** | `json.Number` vs `int` mismatches in connector instance blocks → panics or block reordering diffs (a recurring class of bug; see CLAUDE.md conventions). |
| **429 / rate limiting on big applies** | Slow applies; client retries honoring `Retry-After`. Serialization (`serialization_enabled`, default on) exists precisely because parallel Terraform operations overwhelmed backend concurrency handling. |
| **Auth failures** | `failed to login to portal (4xx)` — expired API key, deprecated username/password path, or portal URL typo. |

## 9. Debugging Guide

Customer reports "terraform apply fails" or "drift detected":

1. **Collect versions first**: provider version (`terraform providers` / lock file), Terraform version, and tenant portal. Map provider version → pinned alkira-client-go version via that tag's `go.mod`.
2. **Reproduce with `TF_LOG=DEBUG` (or `TF_LOG_PROVIDER=DEBUG`)**. The client logs every request with a request ID and lines like `client-create(<requestId>)`, `ALKIRA-PROVISION: true/false`, and provision wait loops (`waiting for provision request <id> to finish (state: ...)`). Grab the request ID and provision request ID.
3. **Classify the failure**:
   - HTTP 4xx on CRUD → payload/schema mismatch or RBAC/auth → compare the DEBUG-logged JSON body against the current backend API contract.
   - Unmarshal error → backend response shape changed; likely needs a client-go fix.
   - `provision request ... failed` → backend provisioning failure, not a provider bug → pivot to backend.
4. **Backend side**: search **Coralogix** (production env) for Alkira Backend Services access and provisioning logs by tenant + timestamp + URI (e.g. `/tenantnetworks/<id>/awsvpcconnectors`). The provision request ID from step 2 is the join key into the backend provisioning logs.
5. **Drift reports**: run `terraform plan -refresh-only`; diff the API GET response (DEBUG log) against tfstate. Perpetual diffs almost always live in a resource's `set*` helper or a missing `Computed`/`DiffSuppress`.
6. **Regression suspicion**: check the release notes / changelog between the customer's previous and current provider version (the `review-terraform-release-notes` skill audits exactly this: schema, ForceNew, defaults, client bump). Downgrade-pin the provider as mitigation.
7. **ForceNew/destroy plans**: never let a customer apply until confirmed intended; check whether the field is genuinely immutable on the backend (silent-overwrite check).

## 10. Ownership Boundaries

**Owns:**
- Terraform schemas, CRUD glue, expand/set helpers, state migrations for every `alkira_*` resource/data source.
- Registry documentation (`docs/`, generated by tfplugindocs from `templates/` + `examples/`).
- Release quality gate for the customer IaC surface (including deciding which client-go version ships).

**Does not own:**
- The HTTP client, auth, retry, serialization, and provision-wait mechanics — that is **alkira-client-go**.
- The API contract or provisioning behavior — that is the **Alkira Backend Services**.
- Customer Terraform state storage or their pipeline tooling.

**Neighboring repos:**
- `alkira-client-go` — sibling library, versioned in lockstep-ish; nearly every provider feature needs a client change first.

## 11. Glossary

- **Resource** — a Terraform-managed Alkira object (`alkira_segment`, …) with full CRUD.
- **Data source** — read-only lookup (usually by name) used to reference existing objects.
- **ForceNew** — schema flag: changing the field requires destroy + recreate. On network resources this is traffic-affecting, so it's added only when the backend silently ignores in-place updates.
- **Provision flag** — provider-level `provision = true`: every apply also triggers and waits on a tenant-network provision request (deploys config onto CXPs). Off by default; without it, applies only change the *intended* config.
- **provision_state** — computed per-resource attribute recording the outcome of the last provision (SUCCESS/FAILED).
- **Serialization** — client-side queue forcing one API call at a time to the portal (default on) to avoid backend write races during parallel Terraform graphs.
- **CXP** — Alkira Cloud Exchange Point, the PoP where tenant networks are provisioned.
- **Connector** — on-ramp from a cloud VPC/VNet, SD-WAN fabric, or site into the tenant's CXP network.
- **Segment** — an isolated routing domain within the tenant network.
- **Drift** — difference between tfstate and live tenant config, surfaced at plan/refresh.

## 12. Repository Health

| Dimension | Rating | Reasoning |
|---|---|---|
| Complexity | **Medium** | ~200 files but each follows one strict resource/helper/data-source pattern; the hard parts are per-resource expand/set edge cases, not architecture. |
| Coupling | **High** | Tightly coupled to both alkira-client-go structs and the live backend API contract; most changes are three-way (provider + client + backend). |
| Tech debt | **Medium** | Deprecated username/password auth still supported; legacy 5-return-value error plumbing (err, errVal, errProv) from the client leaks everywhere; a few resources still on old patterns vs. the documented reference pattern. |
| Documentation | **High** | Registry docs generated per resource, examples/ + acceptance/ HCL, a detailed CLAUDE.md with conventions, and a live-testing guide. |
| Test maturity | **Medium** | Good unit coverage of helpers/schemas plus real acceptance runs against preprod/prod in CI, but acceptance coverage is a subset of 71 resources and depends on live environments. |
