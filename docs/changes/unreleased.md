# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.7.0.md`).

## Advanced IPSec connector — imported tunnels using an overlay reservation now plan cleanly (AK-76265)

For an `alkira_connector_ipsec_adv` tunnel that uses `customer_end_overlay_ip_reservation_id`, the provider no longer writes the backend-derived `customer_end_overlay_ip` into state. Previously both values landed in state, so `terraform plan -generate-config-out` emitted both lines and the generated config failed with "tunnel customer_end_overlay_ip and customer_end_overlay_ip_reservation_id are mutually exclusive".

**Impact:** importing such a connector now plans `1 to import, 0 to change` with no hand edits. For existing reservation-based tunnels, the next refresh clears `customer_end_overlay_ip` in state; this is not a plan diff and makes no API change. Tunnels that set `customer_end_overlay_ip` directly are unchanged. Anything that reads `customer_end_overlay_ip` as an attribute of a reservation-based tunnel now gets an empty string.
