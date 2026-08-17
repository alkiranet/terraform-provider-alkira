# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## Provider credentials (`password`, `api_key`) are now marked sensitive (AK-72007)

The provider-level `password` and `api_key` configuration inputs now carry `Sensitive: true`. Previously, a value set inline in the `provider "alkira"` block was rendered in cleartext in `terraform plan` / `terraform apply` output and captured verbatim in CI job logs; it now renders as `(sensitive value)`.

**Impact:** display-only, and only for the provider block. These are provider configuration inputs, not resource attributes, so they cannot be referenced from an `output` — nothing that plans today starts failing. Note this masks CLI output only; Terraform state still stores whatever the configuration provides, as it does for every provider.

## Bluecat service — anycast `ips`/`backup_cxps` changed to sets to stop spurious diffs (AK-69549)

`ips` and `backup_cxps` inside the `bdds_anycast` and `edge_anycast` blocks of `alkira_service_bluecat` are now `Set of String` instead of `List of String`. The API returns these values in a different order than configured; as ordered lists they re-hashed the surrounding anycast block, so every plan showed the block as removed and re-added even when nothing changed. As sets, order no longer matters.

A state migration (`SchemaVersion` 0 → 1) re-keys existing state from positional to set form on first refresh.

**Impact:** Bluecat services that previously showed a perpetual no-op diff on `bdds_anycast` / `edge_anycast` now plan clean. No HCL changes required — order within `ips` / `backup_cxps` is now irrelevant. The first plan after upgrade is idempotent; the migration runs automatically.
