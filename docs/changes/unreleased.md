# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## Azure VNET connector — `peering_gateway_cxp_id` is now computed, so an omitted value keeps the backend's gateway (AK-75338)

`peering_gateway_cxp_id` on `alkira_connector_azure_vnet` is now `Computed` as well as `Optional`. When the configuration omits it, the provider keeps the gateway ID from state and sends it on update.

**Impact:** updates to a provisioned connector whose configuration omits `peering_gateway_cxp_id` no longer fail with `400 The CXP Peering Gateway of connector '<name>' cannot be updated after it is provisioned.` Configurations that set it explicitly are unaffected. Removing the attribute from a configuration retains the last known value rather than clearing it. No state migration is required.
