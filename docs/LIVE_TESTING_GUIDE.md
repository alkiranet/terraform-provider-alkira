# Live Testing Guide for Terraform Provider Development

This guide provides detailed instructions for testing Terraform provider resources, bug fixes, and new features. It covers the complete testing workflow including setup, execution, and cleanup.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Understanding Dev Overrides](#understanding-dev-overrides)
4. [Test Directory Structure](#test-directory-structure)
5. [Building the Provider](#building-the-provider)
6. [Testing Phases](#testing-phases)
7. [Debug Logging](#debug-logging)
8. [Common Test Scenarios](#common-test-scenarios)
9. [Resource Leak Prevention](#resource-leak-prevention)
10. [Cleanup Procedures](#cleanup-procedures)
11. [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before starting any live testing, ensure you have:

- **Go 1.26+** installed (check with `go version`)
- **Terraform 1.0+** installed (check with `terraform version`)
- **Alkira Portal credentials** (username/password or API key)
- **Access to an Alkira tenant** for testing (use a development tenant when possible)
- **Text editor** for modifying Terraform configurations
- **Terminal** with access to bash/zsh/fish

---

## Environment Setup

### 1. Initial Setup: Configure Dev Overrides

Dev overrides tell Terraform to use your local build instead of the published registry version.

#### Manual Configuration Methods

**Method 1: Edit ~/.terraformrc directly**

Create or edit `~/.terraformrc` (Unix/Linux/macOS) or `%APPDATA%\terraform.rc` (Windows):

**Dev overrides enabled** (for local testing):
```hcl
provider_installation {
  dev_overrides {
    "alkiranet/alkira" = "/path/to/terraform-provider-alkira/"
  }
  direct {}
}
```

**Dev overrides disabled** (for registry testing):
```hcl
provider_installation {
  direct {}
}
```

**Method 2: Use environment variable with alternate config file**

Create two separate config files and switch between them using `TF_CLI_CONFIG_FILE`:

```bash
# Create config for dev overrides
cat > ~/.terraformrc-dev << 'EOF'
provider_installation {
  dev_overrides {
    "alkiranet/alkira" = "/path/to/terraform-provider-alkira/"
  }
  direct {}
}
EOF

# Create config for registry
cat > ~/.terraformrc-registry << 'EOF'
provider_installation {
  direct {}
}
EOF

# Use dev overrides
export TF_CLI_CONFIG_FILE=~/.terraformrc-dev

# Use registry
export TF_CLI_CONFIG_FILE=~/.terraformrc-registry

# Unset to use default ~/.terraformrc
unset TF_CLI_CONFIG_FILE
```

**Tip:** Add the export to your shell profile (`~/.zshrc`, `~/.bashrc`) to make it persistent across sessions.

### 2. Set Up Alkira Credentials

Create a `terraform.tfvars` file in your test directory (never commit this):

```hcl
# Alkira Provider Configuration
alkira_portal = "your-portal.alkira.com"  # Do NOT include https://
alkira_username = "your-username"
alkira_password = "your-password"
# OR use API key instead:
# alkira_api_key = "your-api-key"

# Optional: Control provisioning behavior
alkira_provision = false  # Default: false (async operations)
alkira_async_validation = false  # Default: false
```

**IMPORTANT:** The portal value must be the hostname only, without `https://`. Including the protocol causes the provider to hang silently.

### 3. Environment Variables (Optional)

You can also set credentials via environment variables:

```bash
export ALKIRA_PORTAL="your-portal.alkira.com"
export ALKIRA_USERNAME="your-username"
export ALKIRA_PASSWORD="your-password"
# OR
export ALKIRA_API_KEY="your-api-key"
```

---

## Understanding Dev Overrides

### What Are Dev Overrides?

Dev overrides are a Terraform feature that lets you use a locally-built provider binary instead of downloading one from the Terraform Registry. This is essential for provider development.

### When to Use Each Mode

| Mode | Use When | How to Enable |
|------|----------|---------------|
| **Dev Overrides** | Testing your code changes, debugging new features | Edit `~/.terraformrc` to add `dev_overrides` block, or set `TF_CLI_CONFIG_FILE=~/.terraformrc-dev` |
| **Registry** | Testing against published version, brownfield upgrade testing | Edit `~/.terraformrc` to use only `direct {}`, or set `TF_CLI_CONFIG_FILE=~/.terraformrc-registry` |

### Switching Between Modes

**When switching TO dev overrides (registry → local):**
1. Edit `~/.terraformrc` to add the `dev_overrides` block (or set `TF_CLI_CONFIG_FILE`)
2. Build the provider: `go build -o terraform-provider-alkira`
3. Remove `.terraform` and `.terraform.lock.hcl` from test directory
4. Proceed with testing

**When switching TO registry (local → published):**
1. Edit `~/.terraformrc` to remove the `dev_overrides` block (or set `TF_CLI_CONFIG_FILE`)
2. Remove `.terraform` and `.terraform.lock.hcl` from test directory
3. Run `terraform init` to download registry provider
4. Proceed with testing

---

## Test Directory Structure

Create a dedicated test directory for each ticket or feature:

```
test/
├── AK-12345-feature-name/
│   ├── main.tf                 # Your resource configuration
│   ├── provider.tf             # Provider block (uses variables)
│   ├── variables.tf            # Variable declarations
│   ├── terraform.tfvars        # Credential values (NEVER commit)
│   ├── generated.tf            # Generated config from import
│   ├── terraform.tfstate       # State file (NEVER commit)
│   ├── tf-create.log           # Debug logs from operations
│   ├── tf-plan.log
│   ├── tf-apply.log
│   └── TEST_PLAN.md            # Document your test plan
```

### Example Test Directory Setup

```bash
# Create test directory
mkdir -p test/ak-65421-credential-writeonly
cd test/ak-65421-credential-writeonly

# Create provider.tf
cat > provider.tf << 'EOF'
terraform {
  required_providers {
    alkira = {
      source  = "alkiranet/alkira"
      version = ">= 1.0.0"
    }
  }
}

provider "alkira" {
  portal = var.alkira_portal
  username = var.alkira_username
  password = var.alkira_password
}
EOF

# Create variables.tf
cat > variables.tf << 'EOF'
variable "alkira_portal" {
  type      = string
  sensitive = true
}

variable "alkira_username" {
  type      = string
  sensitive = true
}

variable "alkira_password" {
  type      = string
  sensitive = true
}
EOF

# Create terraform.tfvars with your credentials
cat > terraform.tfvars << 'EOF'
alkira_portal = "your-portal.alkira.com"
alkira_username = "your-username"
alkira_password = "your-password"
EOF
```

---

## Building the Provider

### Development Build (For Dev Overrides)

Build in the repository root (NOT in `bin/`):

```bash
cd /path/to/terraform-provider-alkira
go build -o terraform-provider-alkira
```

This creates `terraform-provider-alkira` binary in the repo root, which dev overrides will use.

### Production Build (For Release)

Use the Makefile:

```bash
make build
```

This creates a production binary in `bin/` with version information.

### Rebuild After Code Changes

**IMPORTANT:** After any code change, you MUST rebuild:

```bash
cd /path/to/repo
go build -o terraform-provider-alkira
```

Terraform does not automatically reload the binary. Rebuild before running each test operation.

---

## Testing Phases

Complete testing involves multiple phases. Not all phases apply to every change.

### Phase 1: Greenfield CRUD (Fresh Resource)

Tests creating a new resource from scratch with your code changes.

#### 1.1 Create

```bash
# Always use debug logging
TF_LOG=DEBUG TF_LOG_PATH=./tf-create.log terraform apply -auto-approve
```

**Verify:**
- [ ] Apply succeeds
- [ ] Resource ID assigned
- [ ] `terraform state show` displays correct values
- [ ] Check raw state: `cat terraform.tfstate | python3 -m json.tool`
- [ ] Debug log shows API request: `grep "REQ:" tf-create.log`

#### 1.2 Read / Refresh

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-refresh.log terraform refresh
terraform plan
```

**Verify:**
- [ ] Refresh succeeds
- [ ] Plan shows **"No changes"**
- [ ] State shows correct values after refresh

#### 1.3 Update

Modify a non-sensitive field in your config (e.g., `name` or `description`):

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-update.log terraform apply -auto-approve
```

**Verify:**
- [ ] Plan shows ONLY the changed field
- [ ] Resource ID stays same (not recreated)
- [ ] API request includes all required fields: `grep "REQ:" tf-update.log`

#### 1.4 Idempotency Check

```bash
terraform plan
```

**Verify:**
- [ ] Shows **"No changes"**

#### 1.5 Delete

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-destroy.log terraform destroy -auto-approve
terraform state list
```

**Verify:**
- [ ] Destroy succeeds
- [ ] `terraform state list` is empty
- [ ] Debug log shows Delete API call

### Phase 2: Import Testing

Tests importing existing resources into Terraform state.

#### 2.1 Create Resource (if not existing)

Create a resource via Terraform or manually via API/portal. Record the resource ID.

#### 2.2 Import with Generate Config

Create an import-only config:

```hcl
# import.tf
import {
  to = alkira_credential_aws_vpc.test
  id = "12345"  # Replace with actual resource ID
}
```

Remove any existing resource blocks and state:

```bash
rm -f terraform.tfstate*
```

Generate config from import:

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-import-generate.log \
  terraform plan -generate-config-out=generated.tf
```

**Verify:**
- [ ] `generated.tf` created with resource block
- [ ] Required fields populated (`name`, `type`, etc.)
- [ ] WriteOnly/sensitive fields show as `null # sensitive`

#### 2.3 Apply the Import

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-import-apply.log \
  terraform apply -auto-approve
```

**Verify:**
- [ ] Import succeeds
- [ ] `terraform state show` displays resource

#### 2.4 Fill in Missing Fields (Critical Test)

For resources with WriteOnly fields, edit `generated.tf` to add the missing sensitive values. Then:

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-plan-after-fill.log terraform plan
```

**Critical Verification:**
- [ ] Plan shows **"No changes"** (no spurious diff)
- [ ] No update API call in logs: `! grep -i "PUT\|POST" tf-plan-after-fill.log`

### Phase 3: Brownfield Upgrade Testing

Tests upgrading from registry provider to your local build.

#### 3.1 Create with Registry Provider

```bash
# Switch to registry - edit ~/.terraformrc to use only direct {}
# OR use: export TF_CLI_CONFIG_FILE=~/.terraformrc-registry
rm -rf .terraform .terraform.lock.hcl
terraform init

# Create resource
terraform apply -auto-approve
```

#### 3.2 Upgrade to Local Build

```bash
# Switch to local build - edit ~/.terraformrc to add dev_overrides block
# OR use: export TF_CLI_CONFIG_FILE=~/.terraformrc-dev
go build -o terraform-provider-alkira
rm -rf .terraform .terraform.lock.hcl

# Check for spurious diffs
terraform plan
```

**Critical Verification:**
- [ ] Plan shows **"No changes"** (upgrade is seamless)

### Phase 4: Environment Variable Testing

For resources supporting environment variable fallback:

```bash
# Unset any existing values
unset ALKIRA_API_KEY

# Create with env vars
export AK_AWS_ACCESS_KEY_ID="test-key"
export AK_AWS_SECRET_ACCESS_KEY="test-secret"
terraform apply -auto-approve

# Test missing field error
unset AK_AWS_ACCESS_KEY_ID
terraform plan  # Should show clear error
```

---

## Debug Logging

### Always Enable Debug Logging

```bash
TF_LOG=DEBUG TF_LOG_PATH=./tf-operation.log terraform <command>
```

### Analyzing Debug Logs

**Search for API requests:**
```bash
grep "REQ:" tf-operation.log
```

**Search for API responses:**
```bash
grep "RSP:" tf-operation.log
```

**Search for specific fields:**
```bash
grep -i "aws_access_key" tf-operation.log
```

**Check for errors:**
```bash
grep -i "error\|fatal" tf-operation.log
```

### Log File Naming Convention

Use descriptive names:
- `tf-create.log` - Initial resource creation
- `tf-refresh.log` - Refresh/Read operation
- `tf-update.log` - Update operation
- `tf-destroy.log` - Deletion
- `tf-import-generate.log` - Import with config generation
- `tf-plan.log` - Plan operation output

---

## Common Test Scenarios

### Scenario 1: Testing Schema Changes

**Change:** Adding a new field to a resource

1. Add field to schema
2. Build provider
3. Create new resource: `terraform apply`
4. Verify field is in state
5. Update field: change value in config, `terraform apply`
6. Verify update works
7. Test import: `terraform plan -generate-config-out=generated.tf`

### Scenario 2: Testing WriteOnly Fields

**Change:** Making sensitive fields WriteOnly

1. Create resource with sensitive values
2. Check `terraform.tfstate` - sensitive fields should be `null`
3. `terraform state show` - should show `(sensitive)`
4. Update non-sensitive field - should work without re-entering secrets
5. Import existing resource - generated config has `null` for sensitive fields
6. Fill in sensitive fields - plan should show "No changes"

### Scenario 3: Testing Bug Fixes

1. Create minimal config to reproduce bug
2. Run with debug logging: `TF_LOG=DEBUG TF_LOG_PATH=./tf-before-fix.log terraform apply`
3. Apply fix
4. Rebuild: `go build -o terraform-provider-alkira`
5. Run with fresh log: `TF_LOG=DEBUG TF_LOG_PATH=./tf-after-fix.log terraform apply`
6. Compare logs to verify fix

### Scenario 4: Testing Default Values

**For Optional + Computed fields:**

1. Create resource WITHOUT specifying the field
2. Apply and check state: value should match API default
3. Plan again: should show "No changes"

**For Optional + Default fields:**

1. Verify default matches ALL existing API states
2. Create resource without field: default applies
3. Import existing resource: no diff if it has default value

### Scenario 5: Testing State Migrations

1. Create resource with OLD schema (registry provider)
2. Note state format
3. Switch to NEW schema (local build with migration)
4. Run `terraform plan`
5. Verify migration triggered automatically
6. Plan shows "No changes" after migration

---

## Resource Leak Prevention

### CRITICAL RULE: Never Delete State Before Resources

**WRONG:**
```bash
rm terraform.tfstate  # Resources still exist in API! Now orphaned!
```

**CORRECT:**
```bash
terraform destroy -auto-approve  # Remove from API first
rm terraform.tfstate*            # Then remove state
```

### Tracking Resource IDs

Always record the IDs of resources you create:

```bash
# After apply, save resource IDs
terraform state list > test-resources.txt
terraform state show alkira_credential_aws_vpc.test | grep "id ="
```

### Verification Steps

**After destroy:**
```bash
terraform state list  # Should be empty
```

**If state is accidentally lost:**
1. Use Alkira Portal to find orphaned resources
2. Manually delete via Portal or API
3. Document what was cleaned up

---

## Cleanup Procedures

### End of Testing Session

```bash
# 1. Destroy all resources
terraform destroy -auto-approve

# 2. Verify state is empty
terraform state list

# 3. Clean up local files (optional - keeping logs is useful)
# rm -f terraform.tfstate*
# rm -rf .terraform .terraform.lock.hcl

# 4. Verify dev overrides are still enabled (if desired)
cat ~/.terraformrc  # Check that dev_overrides block is present
# OR: echo $TF_CLI_CONFIG_FILE  # Check which config file is active
```

### What to Keep vs Delete

**Keep for debugging:**
- `*.log` files (especially on failure)
- `TEST_PLAN.md` with results
- `terraform.tfstate.backup` (for recovery)

**Delete after successful destroy:**
- `terraform.tfstate`
- `terraform.tfstate.backup`
- `.terraform/` directory
- `.terraform.lock.hcl`

**Never delete until resources are gone:**
- Any state files while resources exist

---

## Troubleshooting

### Issue: "Provider not found"

**Cause:** Dev overrides not configured or binary not built

**Solution:**
```bash
# Check dev overrides
cat ~/.terraformrc
# OR check which config file is active
echo $TF_CLI_CONFIG_FILE

# Build the provider
cd /path/to/repo
go build -o terraform-provider-alkira
```

### Issue: Plan shows spurious changes after upgrade

**Cause:** State migration needed or schema incompatibility

**Solution:**
```bash
# Check what changed
terraform plan -out=tfplan

# If it's a known migration, apply to update state
terraform apply tfplan
```

### Issue: Import generates wrong config

**Cause:** Read function returning incorrect values

**Solution:**
```bash
# Check what API is returning
TF_LOG=DEBUG TF_LOG_PATH=./tf-import.log terraform plan -generate-config-out=generated.tf
grep "RSP:" tf-import.log  # See API response
```

### Issue: Provider hangs on apply

**Cause:** Portal URL includes `https://` or network issue

**Solution:**
```bash
# Check provider config - remove https:// from portal value
# Should be: portal = "my-portal.alkira.com"
# NOT: portal = "https://my-portal.alkira.com"
```

### Issue: WriteOnly fields still showing diffs

**Cause:** Plan/Update logic not handling WriteOnly correctly

**Solution:**
```bash
# Check raw config vs state
terraform show -json | jq '.values.root_module.resources[0].values'

# Verify GetRawConfig() is used in generateRequest
```

### Issue: "Version mismatch" errors

**Cause:** Lock file conflicts between registry and dev overrides

**Solution:**
```bash
# Verify dev overrides are configured correctly
cat ~/.terraformrc

# Remove lock file and cache
rm -rf .terraform .terraform.lock.hcl
terraform init  # Or just proceed with dev overrides
```

---

## Quick Reference Commands

```bash
# Build provider
cd /path/to/repo && go build -o terraform-provider-alkira

# Toggle dev overrides - edit ~/.terraformrc or use:
export TF_CLI_CONFIG_FILE=~/.terraformrc-dev       # Enable dev overrides
export TF_CLI_CONFIG_FILE=~/.terraformrc-registry  # Use registry
unset TF_CLI_CONFIG_FILE                           # Use default ~/.terraformrc

# Initialize test directory
terraform init

# Plan with debug
TF_LOG=DEBUG TF_LOG_PATH=./tf-plan.log terraform plan

# Apply with debug
TF_LOG=DEBUG TF_LOG_PATH=./tf-apply.log terraform apply -auto-approve

# Generate config from import
terraform plan -generate-config-out=generated.tf

# Show resource state
terraform state show <resource_name>

# List all resources in state
terraform state list

# Refresh state from API
terraform refresh

# Destroy resources
terraform destroy -auto-approve

# Search logs for API calls
grep "REQ:" tf-operation.log
grep "RSP:" tf-operation.log
```

---

## Appendix: Example Test Plan Template

```markdown
# TEST_PLAN.md

## Ticket: AK-XXXXX - Brief Description

## Changes Under Test
- List code changes being tested
- Link to PR/commit if available

## Test Environment
- Portal: your-portal.alkira.com
- Test Directory: test/ak-xxxxx-description
- Dev Override Config: ~/.terraformrc (or TF_CLI_CONFIG_FILE path)

## Prerequisites
- Dev overrides configured in ~/.terraformrc (or TF_CLI_CONFIG_FILE set)
- Provider built: `go build -o terraform-provider-alkira`
- Test credentials available

## Test Cases

### Phase 1: Greenfield CRUD
- [ ] Create succeeds
- [ ] State shows correct values
- [ ] Refresh shows "No changes"
- [ ] Update works correctly
- [ ] Destroy succeeds

### Phase 2: Import
- [ ] Import generates correct config
- [ ] Fill in fields: No spurious diff

### Phase 3: Upgrade
- [ ] Registry create → local upgrade: No changes

## Results
- Date: YYYY-MM-DD
- Tester: Your Name
- Status: PASSED / FAILED
- Notes: Any issues found
```

---

For questions or issues with this guide, please contact the Terraform Provider team.
