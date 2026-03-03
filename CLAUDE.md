# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## IMPORTANT: Privacy Guidelines

**NEVER mention AI, automation tools, or assistant technologies in any commits, pull requests, documentation, code comments, or any other project artifacts.**

All work should appear as standard human-written code contributions. Do not include co-authorship attributions to automated tools.

### Strict Rules

1. **No attribution references**: Never include any references to automated assistance, language models, or assistant tools in any part of the codebase
2. **Commit messages**: Write conventional, professional commit messages without any special attribution or co-authorship tags
3. **Pull request descriptions**: Describe changes clearly and professionally as if written directly by the developer
4. **Code comments**: Document code with standard inline comments and Javadoc without referencing how the code was created
5. **README files**: Technical documentation should be straightforward without mentioning development tools or assistance
6. **Review comments**: Feedback on pull requests should be written in natural developer voice

## Overview

This is a Terraform provider for Alkira, a network infrastructure management platform. It's built using the Terraform Plugin SDK v2 and interacts with the Alkira API through the `alkira-client-go` library. The provider is maintained by the Alkira engineering team.

**Requirements:**
- Go 1.23.8
- Terraform 1.0+

## Common Commands

### Building
```bash
make build          # Format and build provider binary to bin/
go build -o bin/terraform-provider-alkira
```

### Testing
```bash
make test           # Run all unit tests in alkira/
make test VERBOSE=-v  # Run with verbose output
go test ./alkira/... -run TestAlkiraSegment  # Run specific test
```

### Code Quality
```bash
make fmt            # Format all Go files in alkira/
make lint           # Run golangci-lint on new/changed code
make lint-fix       # Auto-fix linting issues where possible
make lint-all       # Lint all code including legacy code
make doc            # Generate provider documentation using tfplugindocs
```

### Vendoring
```bash
make vendor         # Tidy and vendor dependencies (requires GOPRIVATE=github.com/alkiranet)
```

**Note:** Use `/jira` to access the Jira workflow for working on tickets.

### GitHub and Git Workflows

Use the `gh` CLI for all GitHub-related operations. This ensures consistent interaction with the repository and provides better integration with CI/CD workflows.

#### Viewing Pull Requests

```bash
# View PR summary
gh pr view <PR_NUMBER>

# View PR with specific fields
gh pr view <PR_NUMBER> --json title,body,author,state,files,additions,deletions

# View PR diff
gh pr diff <PR_NUMBER>

# List PRs
gh pr list --state open --limit 20
gh pr list --author "@me"  # Your own PRs
```

#### Creating and Managing PRs

```bash
# Create a PR from current branch
gh pr create --title "AK-XXXXX: Description" --body "Description of changes"

# Create PR with specific base and title
gh pr create --base dev --title "AK-XXXXX: Fix issue" --body "Details..."

# Update PR description
gh pr edit <PR_NUMBER> --body "Updated description"

# Request review from specific users
gh pr edit <PR_NUMBER> --add-reviewer <username>
```

#### Code Review

```bash
# View PR comments
gh pr view <PR_NUMBER> --json comments --jq '.comments[]'

# Add a review comment
gh pr review <PR_NUMBER> --comment -b "My review comment"

# Approve a PR
gh pr review <PR_NUMBER> --approve

# Request changes
gh pr review <PR_NUMBER> --request-changes -b "Issues found..."

# Dismiss a review
gh pr review <PR_NUMBER> --dismiss
```

#### Checking PR Status

```bash
# Check CI status for a specific PR
gh pr checks <PR_NUMBER>

# View status of relevant PRs (current context)
gh pr status

# View merged PRs
gh pr list --state merged --limit 10
```

#### Commit and Branch Management

```bash
# View recent commits for a PR
gh pr view <PR_NUMBER> --json commits --jq '.commits[] | {message, author}'

# Cherry-pick commits (for release branch maintenance)
gh cherry-pick <commit-sha>  # Requires gh-cherry-pick extension
```

#### Git Branch Operations

```bash
# Create a new branch from dev
git checkout dev && git pull
git checkout -b task/AK-XXXXX

# Push branch and create PR in one command
git push -u origin task/AK-XXXXX
gh pr create --base dev
```

#### Useful Aliases (Optional)

Add to `~/.gitconfig` under `[alias]`:
```ini
# Pull Request aliases
pr-view = "!f() { gh pr view $1; }; f"
pr-diff = "!f() { gh pr diff $1; }; f"
pr-open = "!gh pr list --state open --limit 20"
pr-co = "!f() { gh pr checkout $1; }; f"

# Jira-style branch creation
new-task = "!f() { git checkout dev && git pull && git checkout -b task/AK-$1; }; f"
```

### Managing alkira-client-go Dependency

The provider depends on the private `alkira-client-go` library. **NEVER edit files directly in the `/vendor` folder.**

When changes are needed to the client library:

1. **Make changes in alkira-client-go repository first**
   - Implement changes in the alkira-client-go repo
   - Run `make lint` and `make test` in alkira-client-go
   - Ensure all quality checks pass before proceeding
   - Commit and push changes to the main branch

2. **Update vendor in terraform-provider-alkira**
   ```bash
   go get github.com/alkiranet/alkira-client-go@main
   make vendor
   ```

3. **Verify integration**
   ```bash
   make lint   # Ensure vendored code meets quality standards
   make test   # Ensure compatibility
   make build  # Verify compilation
   ```

**Rationale:** This ensures vendored code has been properly linted and tested in its source repository before integration, maintaining consistent code quality across both repositories.

## Architecture

### File Organization Pattern

The codebase follows a consistent naming convention with 174 files in the `alkira/` package:

**Simple resources** (single file):
- `resource_alkira_<type>.go` - Contains schema + CRUD operations + inline helpers

**Complex resources** (main + helper):
- `resource_alkira_<type>.go` - Schema definition + CRUD operations
- `resource_alkira_<type>_helper.go` - Request generation, expand/flatten functions, business logic

**Data sources**:
- `data_source_alkira_<type>.go` - Read-only lookups (typically by name to retrieve ID)

**Tests**:
- `resource_alkira_<type>_test.go` - Unit tests collocated with source
- `acceptance/<type>.tf` - Acceptance tests with live Terraform configurations

### Resource Implementation Pattern

All resources follow this structure:

1. **Schema Definition**: Using `schema.Resource` with field definitions
2. **CRUD Operations**: Create, Read, Update, Delete functions
3. **Helper Functions**:
   - `expand<Type>()` - Convert Terraform schema → API request structures
   - `flatten<Type>()` - Convert API response → Terraform state
   - `generate<Type>Request()` - Build complete API request objects

### AlkiraClient Interaction

The provider uses a consistent pattern for API calls:

```go
client := m.(*alkira.AlkiraClient)  // From provider metadata
api := alkira.NewConnector<Type>(client)  // Type-specific API client

// All Create/Update operations return 5 values:
response, provState, err, valErr, provErr := api.Create(request)
```

**Five return values to handle:**
1. `response` - API response object containing resource ID
2. `provState` - Provisioning state (if provisioning enabled)
3. `err` - General API/network errors (check first)
4. `valErr` - Validation errors (if async validation enabled)
5. `provErr` - Provisioning errors (if provisioning enabled)

Always check errors in order: `err` → `valErr` → `provErr`

### Provider Configuration

The provider supports two authentication methods:

**Option 1: Username/Password**
```hcl
provider "alkira" {
  portal   = "https://portal.example.alkira.com"
  username = "user@example.com"
  password = "secret"
}
```

**Option 2: API Key**
```hcl
provider "alkira" {
  portal  = "https://portal.example.alkira.com"
  api_key = "your-api-key"
}
```

**Environment variables:**
- `ALKIRA_PORTAL` - Portal URL
- `ALKIRA_USERNAME`, `ALKIRA_PASSWORD` - Basic auth
- `ALKIRA_API_KEY` - API key auth
- `ALKIRA_PROVISION` - Enable provisioning (default: false)
- `ALKIRA_ASYNC_VAL` - Enable async validation (default: false)
- `ALKIRA_API_SERIALIZATION_ENABLED` - Enable API call serialization (default: false)
- `ALKIRA_API_SERIALIZATION_TIMEOUT` - Serialization timeout in seconds (default: 120)

### Error Handling Patterns

The codebase implements three-level error handling:

```go
response, provState, err, valErr, provErr := api.Create(request)

if err != nil {
    return diag.FromErr(err)  // API/network failure - FATAL
}

d.SetId(strconv.Itoa(response.Id))  // Always set ID before checking validation/provisioning

// Check validation errors (only if client.Validate enabled)
if client.Validate && valErr != nil {
    var diags diag.Diagnostics
    readDiags := resourceRead(ctx, d, m)
    if readDiags.HasError() {
        diags = append(diags, readDiags...)
    }

    diags = append(diags, diag.Diagnostic{
        Severity: diag.Error,  // FATAL - terminates execution
        Summary:  "VALIDATION (CREATE) FAILED",
        Detail:   fmt.Sprintf("%s", valErr),
    })

    return diags
}

// Check provisioning errors (only if client.Provision enabled)
if client.Provision {
    d.Set("provision_state", provState)

    if provState == "FAILED" {
        return diag.Diagnostics{{
            Severity: diag.Warning,  // NON-FATAL - allows resource creation with warning
            Summary:  "PROVISION (CREATE) FAILED",
            Detail:   fmt.Sprintf("%s", provErr),
        }}
    }
}
```

**Error severity levels:**
- General errors (`err`): Fatal - terminates immediately
- Validation errors (`valErr`): Fatal when validation is enabled - uses `diag.Error`
- Provisioning errors (`provErr`): Non-fatal when provisioning is enabled - uses `diag.Warning`

Note: For Delete operations, validation errors use `diag.Warning` instead of `diag.Error`.

### Common Utility Functions

Located in [helper.go](alkira/helper.go):

- `convertTypeListToStringList()` - Convert Terraform TypeList to []string
- `convertTypeSetToIntList()` - Convert Terraform TypeSet to []int
- `expandSegmentOptions()` - Expand segment-to-zone mappings
- `randomNameSuffix()` - Generate random suffix for test resources
- `convertInputTimeToEpoch()` - Convert time strings to epoch

### Testing Conventions

**Unit Tests** (in `alkira/` directory):
- Test schema validation and structure
- Test expand/flatten helper functions
- Test request generation logic
- Run with: `make test` or `go test ./alkira/...`
- Naming: `TestAlkira<ResourceType>_<TestName>`

**Acceptance Tests** (in `acceptance/` directory):
- Require live Alkira credentials and portal
- Test actual Terraform configurations against real API
- Run via GitHub Actions workflows (acceptance-preprod.yml, acceptance-prod.yml)

### Key Patterns to Follow

1. **Expand/Flatten Pattern**: Always create matching expand/flatten functions for nested structures
2. **Import Support**: All resources use `importWithReadValidation(resourceRead)` wrapper for better error handling on invalid IDs (added via #450, AK-64152)
3. **Conflict Resolution**: Use `ConflictsWith` for mutually exclusive fields
4. **CustomizeDiff**: Used to reset FAILED provisioning states on configuration changes
5. **Computed Fields**: Mark fields as `Computed: true` if set by API responses
6. **Required vs Optional**: Use `Required`, `Optional`, or `Optional + Computed` appropriately

### Resource Lifecycle

When implementing new resources:

1. Define schema in `resource_alkira_<type>.go`
2. Implement Create, Read, Update, Delete functions
3. Create helper file if structure is complex (10+ nested fields)
4. Add expand/flatten functions for nested blocks
5. Write unit tests for schema and helpers
6. Add example Terraform configuration to `examples/`
7. Generate documentation with `make doc`

### Code Quality Standards

The repository uses golangci-lint with a configuration optimized for Terraform providers:

**Enabled linters (13):**
- Critical bug detection: `govet`, `staticcheck`, `ineffassign`, `unused`
- Code quality: `misspell`, `unconvert`, `wastedassign`
- Error handling: `errname`, `errorlint`, `nilerr`
- Resource leaks: `bodyclose`, `sqlclosecheck`, `copyloopvar`

**Explicitly disabled:**
- Style linters that would break compatibility: `revive`, `stylecheck`, `varnamelen`
- Function complexity linters (Terraform resources are inherently complex)
- Named return prohibition (5-return-value pattern uses named returns)

**Special exclusions:**
- `errcheck` disabled - cannot exclude Terraform SDK methods (`d.Set()`, `d.SetId()`) in golangci-lint
- Test files have relaxed complexity rules
- Resource files allow longer functions (common pattern)
- Staticcheck rules disabled: S1009 (explicit nil checks), ST* naming rules

**CI enforcement:**
- Linting runs on all PRs via GitHub Actions
- Only new/changed code is checked (non-disruptive to legacy code)
- Auto-fixes available via `make lint-fix`

### Debugging and Development

When debugging issues:
- Check the 5-return-value pattern is correctly handled
- Verify expand/flatten functions are symmetric
- Ensure SetId() is called before checking validation/provisioning errors
- Look for ConflictsWith constraints on fields
- Check that Computed fields aren't being set in schema incorrectly
- Run `make lint` before submitting PRs to catch common issues

#### Running Terraform Commands with Logging

When testing provider changes, enable debug logging to capture detailed API interactions:

```bash
# Enable debug logging with auto-approve (non-interactive)
TF_LOG=DEBUG TF_LOG_PATH=./tf.log terraform apply -auto-approve
TF_LOG=DEBUG TF_LOG_PATH=./tf-plan.log terraform plan
TF_LOG=DEBUG TF_LOG_PATH=./tf-import.log terraform import alkira_connector_gcp_vpc.test 16609
```

**Key flags:**
- `TF_LOG=DEBUG` - Enables verbose debug output
- `TF_LOG_PATH=<file>` - Writes logs to specified file instead of stdout
- `-auto-approve` - Skips interactive approval prompt (useful for testing/CI)

### Live Testing Bug Fixes (Brownfield & Greenfield)

When testing schema changes, new fields, or bug fixes against a live Alkira environment, always test **both** brownfield (existing customers upgrading) and greenfield (fresh import/new resources) scenarios. This catches backward-compatibility breaks that unit tests cannot detect.

**Prerequisites from user:**
- Portal URL and credentials (username/password or API key)
- Resource ID(s) to test against (e.g., connector ID for imports)
- API response JSON for the resource (to understand the actual state on the backend, e.g., whether a field is `null`, `false`, or `true`)
- Any cloud-specific IDs if creating new resources (e.g., GCP credentials, VPC IDs)

#### Environment Setup

The dev_overrides in `~/.terraformrc` points to the repo root:
```
dev_overrides {
  "alkiranet/alkira" = "/users/asim/alkira-workspace/terraform-provider-alkira/"
}
```

When dev_overrides is active, Terraform uses the binary at the repo root (`terraform-provider-alkira`), skipping `terraform init`. Build it with:
```bash
go build -o terraform-provider-alkira   # builds at repo root for dev_overrides
```

To use the **published registry version** instead, temporarily replace `~/.terraformrc`:
```
provider_installation {
  direct {}
}
```
Then run `terraform init` to pull the registry version. **Always restore the original `~/.terraformrc` afterwards.**

#### Test Directory Structure

Create a test directory under `test/` for each test scenario:
```bash
mkdir -p test/<descriptive_name>
```

Write a minimal `main.tf` with provider config and import block:
```hcl
terraform {
  required_providers {
    alkira = {
      source = "alkiranet/alkira"
    }
  }
}

provider "alkira" {
  portal    = "<portal_url>"
  username  = "<username>"
  password  = "<password>"
  provision = false
}

import {
  to = alkira_<resource_type>.<name>
  id = "<resource_id>"
}
```

#### Always Use Debug Logging

Every terraform command must use `TF_LOG=DEBUG TF_LOG_PATH=./tf-<operation>.log` so API request/response payloads can be inspected after the fact:
```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-plan.log terraform plan
TF_LOG=DEBUG TF_LOG_PATH=./tf-apply.log terraform apply -auto-approve
TF_LOG=DEBUG TF_LOG_PATH=./tf-refresh.log terraform refresh
TF_LOG=DEBUG TF_LOG_PATH=./tf-generate.log terraform plan -generate-config-out=generated.tf
```

#### Test 1: Brownfield (Provider Upgrade — No Breaking Diffs)

This tests that existing customers upgrading from the published provider to the new build see **no unexpected diffs** in `terraform plan`.

**Steps:**

1. **Disable dev_overrides** in `~/.terraformrc` (use `direct {}` only)
2. **Init + import with registry provider:**
   ```bash
   terraform init
   TF_LOG=DEBUG TF_LOG_PATH=./tf-generate-registry.log terraform plan -generate-config-out=generated.tf
   TF_LOG=DEBUG TF_LOG_PATH=./tf-apply-registry.log terraform apply -auto-approve
   ```
3. **Remove the import block** from `main.tf` (import is done, state exists)
4. **Restore dev_overrides** in `~/.terraformrc`
5. **Build the fixed binary:**
   ```bash
   cd /Users/asim/alkira-workspace/terraform-provider-alkira
   go build -o terraform-provider-alkira
   ```
6. **Plan with fixed build against registry-imported state:**
   ```bash
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-fixed.log terraform plan
   ```
7. **Expected result: `No changes.`** If there's a diff, the fix has a backward-compatibility problem.
8. **Refresh and re-plan** to verify state is updated cleanly:
   ```bash
   TF_LOG=DEBUG TF_LOG_PATH=./tf-refresh-fixed.log terraform refresh
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-post-refresh.log terraform plan
   ```
9. **Verify state** has the new field populated correctly:
   ```bash
   terraform state show <resource_address>
   # Or check raw JSON:
   cat terraform.tfstate | python3 -m json.tool | grep -A 20 "<field_name>"
   ```

#### Test 2: Greenfield (Fresh Import with Fixed Build)

This tests that the original bug is fixed — new imports/generate-config include the field properly.

**Steps:**

1. **Clean up** any previous state:
   ```bash
   rm -f terraform.tfstate terraform.tfstate.backup generated.tf
   ```
2. **Ensure dev_overrides is active** and fixed binary is built
3. **Add import block** back to `main.tf`
4. **Generate config:**
   ```bash
   TF_LOG=DEBUG TF_LOG_PATH=./tf-generate-fixed.log terraform plan -generate-config-out=generated.tf
   ```
5. **Verify** the generated config includes the new/fixed field with the correct value
6. **Apply the import:**
   ```bash
   TF_LOG=DEBUG TF_LOG_PATH=./tf-apply-fixed.log terraform apply -auto-approve
   ```
7. **Verify state** has the field correctly populated from the API response
8. **Plan again** to confirm no drift:
   ```bash
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-post-apply.log terraform plan
   ```
9. **Expected result: `No changes.`**

#### Test 3: Mutation Tests (Explicit Field Changes)

After brownfield or greenfield setup, test that the field responds correctly to explicit user changes:

1. **Explicit value matching state** — set the field to the same value already in state:
   ```bash
   # Add export_all_subnets = false to config (if state has false)
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-explicit-match.log terraform plan
   # Expected: No changes
   ```

2. **Explicit value differing from state** — set to the opposite value:
   ```bash
   # Change export_all_subnets = true in config (state has false)
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-explicit-diff.log terraform plan
   # Expected: Diff detected. For bool fields with mutual exclusion (like export_all_subnets + vpc_subnet),
   # CustomizeDiff should reject invalid combinations.
   ```

3. **Remove field from config** — after it was explicitly set, remove it:
   ```bash
   # Delete export_all_subnets line from config
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-field-removed.log terraform plan
   # Expected: No changes (Computed defers to state value)
   ```

4. **Update unrelated field without the new field in config** — this catches silent value changes:
   ```bash
   # Change description (or any other field) while export_all_subnets is NOT in config
   TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-update-other.log terraform plan
   # Expected: Only the changed field appears in the diff, NOT the new field
   TF_LOG=DEBUG TF_LOG_PATH=./tf-apply-update-other.log terraform apply -auto-approve
   # CRITICAL: Check the API PUT payload to confirm the new field's value is preserved:
   grep -i "<field_json_name>" tf-apply-update-other.log
   ```

This last test is especially important for `Optional + Computed` fields with `GetRawConfig` defaults. If the default logic isn't scoped correctly (e.g., only during Create), updating any other field can silently flip the new field's value in the API request.

#### What to Verify in Logs

Check the debug logs for:
- **API GET response** during Read: Confirm the field value matches what the API actually returns
- **API PUT/POST payload** during Create/Update: Confirm the field is sent with the correct value, especially during Updates where the field is not in the user's config
- **State diff** during Plan: If unexpected diffs appear, compare the config value vs state value vs API value

```bash
# Search for API payloads in logs
grep -i "request body\|response body\|REQ:\|RSP:" tf-apply.log | grep -i "<field_json_name>"
```

#### Edge Case Analysis Checklist

Before considering a fix complete, trace through **every combination** of these dimensions:

**API state variations** — check the actual API JSON for test resources:
- Field is `null` (never explicitly set on the backend)
- Field is `false` / zero value
- Field is `true` / non-zero value
- Different resources may have different API states — test multiple if they vary

**User config variations:**
- Field explicitly set to `true`
- Field explicitly set to `false`
- Field omitted from config entirely
- Entire parent block omitted (e.g., no `gcp_routing` at all)

**Lifecycle operations** — the code path differs for each:
- **Create**: `generateRequest` builds the full request; `GetRawConfig` has a real user config
- **Read/Import**: Only `setOptions`/flatten runs; `GetRawConfig` may return null
- **Update**: `generateRequest` rebuilds the full request; `GetRawConfig` has user config but the plan also includes prior state for unset Computed fields
- **CustomizeDiff**: Runs during plan; uses `ResourceDiff` not `ResourceData`

**Key pitfall for `Optional + Computed` booleans:**
- The SDK cannot distinguish "user set false" from "user didn't set it" — both are `false` in `d.Get()`
- Use `d.GetRawConfig()` to check the actual HCL, but **scope it correctly**:
  - During **Create** (`d.Id() == ""`): safe to check and apply defaults
  - During **Update** (`d.Id() != ""`): the plan already carries the correct state value for unset Computed fields — do NOT override with `GetRawConfig`, or you'll silently change the value
  - During **Read/Import**: `generateRequest` is not called, so irrelevant

#### Schema Design Rules

When adding a new field that the API already returns for existing resources:

| Scenario | Schema | Why |
|----------|--------|-----|
| Field has a sensible server-side default, user may not set it | `Optional: true, Computed: true` | Defers to state/API when user omits it — no breaking diff on upgrade |
| Field must always have a user-specified value | `Required: true` | Forces user to set it explicitly |
| Field has a fixed default that matches ALL existing API states | `Optional: true, Default: <value>` | Safe only if the default matches what the API returns for every existing resource |

**Critical rule:** Never use `Optional + Default` if the default value could differ from what the API returns for existing resources. This creates a diff on upgrade where the state (from API) says one value but the config (from Default) says another. Use `Optional + Computed` instead — the helper/expand function can still provide a default for new resources, while existing resources defer to their API-stored value.

**`GetRawConfig` pattern for `Optional + Computed` bool defaults:**
```go
// Only during Create — Update already has correct state values
if d.Id() == "" {
    rawConfig := d.GetRawConfig()
    blockRaw := rawConfig.GetAttr("block_name")
    if !blockRaw.IsNull() && blockRaw.IsKnown() && blockRaw.LengthInt() > 0 {
        fieldVal := blockRaw.Index(cty.NumberIntVal(0)).GetAttr("field_name")
        if fieldVal.IsNull() {
            result.Field = defaultValue // User didn't set it
        }
    }
}
```

#### Cross-Version Compatibility Testing

Every bug fix or schema change must be validated against the **published registry version** to ensure existing customers can upgrade without breaking. The core principle: state written by the old provider must be readable by the new provider without unexpected diffs.

**Why this matters:**
- Customers have existing `.tfstate` files written by the published provider
- When they upgrade, `terraform plan` runs the new provider's Read against the old state
- Any mismatch between what the new schema expects vs what the old state contains = breaking diff
- Breaking diffs are especially dangerous when they trigger API calls that the backend rejects

**The testing flow:**

1. **Registry provider (old)** → import resource → creates state without the new field
2. **Fixed build (new)** → plan against that state → must show `No changes`
3. **Fixed build (new)** → fresh import → must include the new field correctly

This sequence simulates exactly what happens when a customer upgrades their provider version. The brownfield test (step 2) catches backward-compatibility breaks. The greenfield test (step 3) confirms the original bug is fixed.

**When to test multiple resources:**
- If the API response varies across resources (e.g., one returns `null` for a field, another returns `false`), test both. Fetch the API JSON for each resource to understand the actual backend state before choosing test targets:
  ```bash
  # Check the actual API response for a resource to understand its state
  # The user should provide the API JSON or you can check the debug logs
  grep -A 5 "exportAllSubnets" tf-apply-registry.log
  ```

**Versioned binary management:**
- Registry version: disable dev_overrides, run `terraform init`
- Fixed build: enable dev_overrides, run `go build -o terraform-provider-alkira`
- Always back up `~/.terraformrc` before modifying: `cp ~/.terraformrc ~/.terraformrc.bak`
- Always restore after testing: `cp ~/.terraformrc.bak ~/.terraformrc`

#### Cleanup

After testing, clean up the test directory state files but keep `main.tf` and `generated.tf` as reference:
```bash
rm -f terraform.tfstate terraform.tfstate.backup
rm -rf .terraform .terraform.lock.hcl
```
