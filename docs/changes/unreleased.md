# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## Bluecat service — anycast `ips`/`backup_cxps` changed to sets to stop spurious diffs (AK-69549)

`ips` and `backup_cxps` inside the `bdds_anycast` and `edge_anycast` blocks of `alkira_service_bluecat` are now `Set of String` instead of `List of String`. The API returns these values in a different order than configured; as ordered lists they re-hashed the surrounding anycast block, so every plan showed the block as removed and re-added even when nothing changed. As sets, order no longer matters.

A state migration (`SchemaVersion` 0 → 1) re-keys existing state from positional to set form on first refresh.

**Impact:** Bluecat services that previously showed a perpetual no-op diff on `bdds_anycast` / `edge_anycast` now plan clean. No HCL changes required — order within `ips` / `backup_cxps` is now irrelevant. The first plan after upgrade is idempotent; the migration runs automatically.
