# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## Bluecat service — anycast `ips`/`backup_cxps` changed to sets to stop spurious diffs (AK-69549)

`ips` and `backup_cxps` inside the `bdds_anycast` and `edge_anycast` blocks of `alkira_service_bluecat` are now `Set of String` instead of `List of String`. The API returns these values in a different order than configured; as ordered lists they re-hashed the surrounding anycast block, so every plan showed the block as removed and re-added even when nothing changed. As sets, order no longer matters.

A state migration (`SchemaVersion` 0 → 1) re-keys existing state from positional to set form on first refresh.

**Impact:** Bluecat services that previously showed a perpetual no-op diff on `bdds_anycast` / `edge_anycast` now plan clean. No HCL changes required — order within `ips` / `backup_cxps` is now irrelevant. The first plan after upgrade is idempotent; the migration runs automatically.

## IPSec connectors — `PING` removed from `routing_options.availability` (AK-73307)

`availability` inside the `routing_options` block of `alkira_connector_ipsec` and `alkira_connector_ipsec_adv` no longer accepts `PING`. The valid values are `IKE_STATUS` and `IPSEC_INTERFACE_PING`, matching the two options the portal offers. `PING` was never selectable in the UI, and a connector configured with it rendered in the portal as if `IKE Status` were selected while the backend still held `PING` — an inconsistent state that a subsequent save from the portal would silently rewrite. `PING` and `IPSEC_INTERFACE_PING` resolve to the same tunnel probe on the backend, so switching between them does not change how availability is monitored.

**Impact:** configurations that set `availability = "PING"` now fail validation with an error naming the two accepted values; replace it with `IPSEC_INTERFACE_PING`. Configurations that omit `availability` are unaffected, and for connectors whose stored value is `PING` the next plan shows a single in-place change to `IPSEC_INTERFACE_PING` that converges on apply. No state migration is required. Note that `terraform plan -refresh-only` reports a stored `PING` as-is without flagging or correcting it — a normal apply is needed to update the connector.
