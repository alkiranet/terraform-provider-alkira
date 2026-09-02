# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## Provider credentials (`password`, `api_key`) are now marked sensitive (AK-72007)

The provider-level `password` and `api_key` configuration inputs now carry `Sensitive: true`. Previously, a value set inline in the `provider "alkira"` block was rendered in cleartext in `terraform plan` / `terraform apply` output and captured verbatim in CI job logs; it now renders as `(sensitive value)`.

**Impact:** display-only, and only for the provider block. These are provider configuration inputs, not resource attributes, so they cannot be referenced from an `output` — nothing that plans today starts failing. Note this masks CLI output only; Terraform state still stores whatever the configuration provides, as it does for every provider.

## Bluecat service — anycast `ips`/`backup_cxps` changed to sets to stop spurious diffs (AK-69549)

`ips` and `backup_cxps` inside the `bdds_anycast` and `edge_anycast` blocks of `alkira_service_bluecat` are now `Set of String` instead of `List of String`. The API returns these values in a different order than configured; as ordered lists they re-hashed the surrounding anycast block, so every plan showed the block as removed and re-added even when nothing changed. As sets, order no longer matters.

A state migration (`SchemaVersion` 0 → 1) re-keys existing state from positional to set form on first refresh.

**Impact:** Bluecat services that previously showed a perpetual no-op diff on `bdds_anycast` / `edge_anycast` now plan clean. No HCL changes required — order within `ips` / `backup_cxps` is now irrelevant. The first plan after upgrade is idempotent; the migration runs automatically.

## IPSec connectors — `PING` removed from `routing_options.availability` (AK-73307)

`availability` inside the `routing_options` block of `alkira_connector_ipsec` and `alkira_connector_ipsec_adv` no longer accepts `PING`. The valid values are `IKE_STATUS` and `IPSEC_INTERFACE_PING`, matching the two options the portal offers. `PING` was never selectable in the UI, and a connector configured with it rendered in the portal as if `IKE Status` were selected while the backend still held `PING` — an inconsistent state that a subsequent save from the portal would silently rewrite. `PING` and `IPSEC_INTERFACE_PING` resolve to the same tunnel probe on the backend, so switching between them does not change how availability is monitored.

**Impact:** configurations that set `availability = "PING"` now fail validation with an error naming the two accepted values; replace it with `IPSEC_INTERFACE_PING`. Configurations that omit `availability` are unaffected, and for connectors whose stored value is `PING` the next plan shows a single in-place change to `IPSEC_INTERFACE_PING` that converges on apply. No state migration is required. Note that `terraform plan -refresh-only` reports a stored `PING` as-is without flagging or correcting it — a normal apply is needed to update the connector.

## Fortinet service — `segment_ids` and `instances` are now populated on refresh and import (AK-73611)

`alkira_service_fortinet` builds `segment_ids` in its Read path, because the API stores a service's segments by name while the schema stores them by ID. That conversion produced a `[]int`, which the SDK rejects against the declared `Set of String` element type, and it also pre-sized its result slice before appending, prefixing the value with one zero per segment. The rejection returned an error that was discarded, so Read never updated `segment_ids` while reporting success. Read now uses the same `convertSegmentNamesToSegmentIds` helper as the Infoblox, Bluecat and Zscaler services, writes the segment IDs as strings, and checks the error.

Read also dropped all but one `instances` entry. It walks the instances already tracked in state, then walks the API response for any it did not match, and that second walk stopped after the first unmatched instance. A service with more instances than state knew about landed exactly one of them in `instances`, which is `Required`.

**Impact:** for a service created through Terraform, `segment_ids` was written to state at create time and the failed refresh left it untouched, so there was no spurious diff, but drift on the attribute went undetected. Drift is now detected: a service whose segments were reassigned outside Terraform shows a one-time in-place diff on the next plan that restores the configured segments. The same applies to a service whose instances were changed outside Terraform. After `terraform import`, where state starts empty, `segment_ids` stayed empty and `instances` held a single entry, so the first plan asked to add every segment and every remaining instance back; both are now populated during the import. Configurations already in sync with the portal are unaffected. No HCL changes and no state migration are required; `segment_ids` remains `Set of String`, matching the six other services that key segments by ID.

## Segment resource and segment resource share — `segment_id` and `designated_segment_id` reject a segment name (AK-74335)

`segment_id` on `alkira_segment_resource` and `designated_segment_id` on `alkira_segment_resource_share` now validate that the value is a segment ID. Both fields are documented as IDs, but they accepted any alphanumeric string, and the provider handed the value straight to `GET /segments/<value>`. A segment name there returns a 500 that the client retries five times with escalating backoff, so an apply sat for over two minutes before failing on a type-conversion error from the backend. The value is now rejected during `terraform plan`, before any API call, with an error that names the field and points at `alkira_segment.example.id`. Leading zeros are rejected for the same reason a name is: `GET /segments/0690` succeeds and answers with id `690`, which Read writes back to state against a config that still says `0690`.

`alkira_segment_resource_share` also stopped refreshing partway through when the segment lookup failed. Read returned at that point, and Terraform counts a warning as a successful refresh, so `end_a_segment_resource_ids`, `end_b_segment_resource_ids`, the route limits, the traffic fields and `policy_rule_list_id` kept whatever state already held. Read now skips only `designated_segment_id` and refreshes the rest. The lookup stays non-fatal on purpose: the share's own `GetById` asks for a resource marked for deletion while the segment get-by-name does not, so making it fatal would abort the refresh that `terraform destroy` runs first.

**Impact:** configurations that reference a segment by ID, including `alkira_segment.example.id`, are unaffected, and importing either resource continues to plan clean. No state migration is required.

One group does need an edit: anyone whose HCL holds a segment name in either field. That configuration planned clean because a plan only compares, so the mismatch went unnoticed. For `alkira_segment_resource_share`, a name is also what `terraform import` wrote into `designated_segment_id` on builds between the two AK-64221 fixes, so HCL written to match an import from that window holds one. Validation now rejects it, and the failure blocks `terraform plan`, `terraform apply` and `terraform destroy` for the whole root module until the value changes to `alkira_segment.<name>.id` or the segment's numeric ID.

## Segment IDs are validated before the segment lookup, across every resource that takes one (AK-74389)

`getSegmentNameById` turns a `segment_id` from configuration into `GET /segments/<value>`. It validated that value with `validateReferenceId`, whose pattern accepts letters, so a segment **name** reached the backend, which answers a non-numeric id with a 500. The client retries a 500 five times with escalating backoff, so an apply spent around four minutes on nine requests before failing on `giving up after 6 attempt(s): retryable status code: 500`. The value is now checked against the segment ID pattern first, and a name fails in about a second with an error naming the value and pointing at `alkira_segment.example.id`. Leading zeros are rejected for the same reason they are on `alkira_segment_resource`: the backend accepts `GET /segments/0690` and answers with id `690`, which Read then writes back to state against the config that produced it.

This covers the 37 call sites that reach `getSegmentNameById` directly, the six services that reach it through `convertSegmentIdsToSegmentNames`, and `segment_options.segment_id`, which ran the same lookup inline.

Three code paths discarded the lookup error and sent an empty segment name to the API: the `global_protect_segment_options` and `global_protect_segment_options_instance` blocks of `alkira_service_pan`, and the `firepower_management_center` block of `alkira_service_cisco_ftdv`. A fourth, `alkira_connector_remote_access`, dropped every segment name and continued with none. All four now report the failure.

**Impact:** configurations that reference segments by ID are unaffected. A configuration that passes a segment name already failed; it now fails in a second with a message that says what to change, instead of after four minutes with a backend status code. For the four paths above, a failure that previously produced a silently wrong API request now stops the apply. No HCL changes and no state migration are required.
