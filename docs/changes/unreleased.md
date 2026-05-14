# Unreleased

User-visible changes since the last release. Renamed to the release version when cut (e.g. `v1.6.0.md`).

## PAN service — sensitive credentials + plan-time validation (AK-66660)

**Sensitive fields.** `pan_password`, `pan_license_key`, `master_key`, `registration_pin_id`, `registration_pin_value`, and per-instance `auth_key` / `auth_code` on `alkira_service_pan` are now `Sensitive: true`. Values are redacted in `plan`, `apply`, and `show` output.

**Plan-time validation.** New `CustomizeDiff` rules:
- `pan_username` must be `admin` (AWS / AWS China / GCP CXPs) or `akadmin` (Azure / Azure China). Other CXP types accept either.
- `master_key` and `master_key_expiry` required when `master_key_enabled = true` on Create or on a false→true transition.
- Multi-instance configurations must use unique non-empty `name` values.

**Bootstrap-only warning.** `apply` emits a `Warning` when `master_key`, `registration_pin_*`, or per-instance `auth_*` fields change post-create. The backend consumes these at PAN VM bootstrap only — changes are saved to state but not pushed to running devices. To rotate, destroy and recreate.

**Impact.** None for valid existing configurations. Configurations with a mismatched `pan_username` or multi-instance blocks with duplicate / empty names start failing at plan time; the backend already rejected these at apply.

**State file on disk is unchanged.** `Sensitive: true` redacts CLI output only. `terraform.tfstate` continues to store credential values in plaintext.
