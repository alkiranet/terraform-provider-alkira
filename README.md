![unit-test](https://github.com/alkiranet/terraform-provider-alkira/actions/workflows/unit-test.yaml/badge.svg)
![lint](https://github.com/alkiranet/terraform-provider-alkira/actions/workflows/lint.yaml/badge.svg)
![acceptance-preprod](https://github.com/alkiranet/terraform-provider-alkira/actions/workflows/acceptance-preprod.yml/badge.svg)
![acceptance-prod](https://github.com/alkiranet/terraform-provider-alkira/actions/workflows/acceptance-prod.yml/badge.svg)

# Terraform Provider for Alkira

The Terraform provider for Alkira enables full lifecycle management of Alkira Cloud Services Exchange resources.

- **Website**: [www.alkira.com](http://www.alkira.com)
- **Documentation**: [Terraform Registry](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs)

## Table of Contents

- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Building the Provider](#building-the-provider)
- [Testing](#testing)
- [Code Quality](#code-quality)
- [Contributing](#contributing)
  - [Adding New Resources](#adding-new-resources)
  - [Creating Examples](#creating-examples)
  - [Creating Documentation Templates](#creating-documentation-templates)
  - [Generating Documentation](#generating-documentation)
- [Project Structure](#project-structure)
- [CI/CD](#cicd)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.0 or later
- [Go](https://golang.org/doc/install) 1.23.8
- [Git](https://git-scm.com/downloads)

## Getting Started

To use the provider in your Terraform configurations:

```hcl
terraform {
  required_providers {
    alkira = {
      source  = "alkiranet/alkira"
      version = "~> 1.0"
    }
  }
}

provider "alkira" {
  portal   = "your_tenant_name.portal.alkira.com"
  username = "your_name@email.com"
  password = "your_password"
}
```

For detailed provider configuration and authentication options, see the [provider documentation](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs).

## Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/alkiranet/terraform-provider-alkira.git
cd terraform-provider-alkira
```

### 2. Install Required Tools

#### Install tfplugindocs

[tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) is used to generate provider documentation automatically. Install it using:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

Verify installation:

```bash
tfplugindocs --version
```

#### Install golangci-lint

[golangci-lint](https://golangci-lint.run/) is a fast Go linters aggregator used for code quality checks.

**macOS:**
```bash
brew install golangci-lint
```

**Linux:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**Windows:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Verify installation:

```bash
golangci-lint --version
```

### 3. Configure Private Dependencies

The provider depends on the private `alkira-client-go` library. Configure Go to access private Alkira repositories:

```bash
export GOPRIVATE=github.com/alkiranet
```

Add this to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it persistent.

### 4. Install Dependencies

```bash
make vendor
```

This command runs `go mod tidy` and `go mod vendor` to fetch and vendor all dependencies.

**Note:** If you need to make changes to the `alkira-client-go` library, see [Making Changes to alkira-client-go](#making-changes-to-alkira-client-go) for the required workflow to maintain code quality.

## Building the Provider

Build the provider binary:

```bash
make build
```

This will:
1. Format all Go files in the `alkira/` package
2. Build the binary to `bin/terraform-provider-alkira`

### Cross-Platform Builds

To build release binaries for multiple platforms:

```bash
make release
```

This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Testing

### Unit Tests

Run all unit tests:

```bash
make test
```

Run tests with verbose output:

```bash
make test VERBOSE=-v
```

Run a specific test:

```bash
go test ./alkira/... -run TestAlkiraSegment
```

### Acceptance Tests

Acceptance tests run against a live Alkira environment and require valid credentials:

```bash
# Set environment variables
export ALKIRA_PORTAL="tenant.portal.alkira.com"
export ALKIRA_USERNAME="your_username"
export ALKIRA_PASSWORD="your_password"

# Run acceptance tests (if configured)
# Note: These are typically run via GitHub Actions
```

Acceptance test configurations are located in the `acceptance/` directory.

## Code Quality

### Formatting

Format all Go code:

```bash
make fmt
```

This runs `go fmt` on all files in the `alkira/` package.

### Linting

The project uses [golangci-lint](https://golangci-lint.run/) for comprehensive code quality checks. The configuration in `.golangci.yml` is tailored for Terraform providers with 13 enabled linters.

**Run linting (checks only new/changed code):**

```bash
make lint
```

**Auto-fix issues where possible:**

```bash
make lint-fix
```

**Lint all code including legacy code:**

```bash
make lint-all
```

**What gets checked:**
- Error handling (error wrapping, naming conventions)
- Bug detection (nil checks, unreachable code, ineffectual assignments)
- Code quality (formatting, spelling, unnecessary conversions)
- Resource leaks (HTTP body close, SQL connections, loop variable capture)
- Best practices (staticcheck rules, vet analysis)

**Linting strategy:**
- CI enforces linting on new/changed code only (non-disruptive)
- Auto-installable via Makefile (no manual setup required)
- Configured to avoid breaking backwards compatibility
- Excludes overly opinionated style rules

The linting setup focuses on catching real bugs while allowing incremental improvements to the codebase.


## Contributing

We welcome contributions! Please follow these guidelines when contributing to the provider.

### Adding New Resources

When adding a new Terraform resource or data source:

#### 1. Implement the Resource

Create resource files in the `alkira/` directory following the naming convention:

**For simple resources (single file):**
- `resource_alkira_<type>.go` - Contains schema, CRUD operations, and inline helpers

**For complex resources (split files):**
- `resource_alkira_<type>.go` - Schema definition and CRUD operations
- `resource_alkira_<type>_helper.go` - Request generation, expand/flatten functions

**For data sources:**
- `data_source_alkira_<type>.go` - Read-only data source implementation

**Key patterns to follow:**
- Use the 5-return-value pattern for API calls: `response, provState, err, valErr, provErr`
- Always call `d.SetId()` before checking validation/provisioning errors
- Implement matching `expand<Type>()` and `flatten<Type>()` functions for nested structures
- Add proper error handling for general, validation, and provisioning errors
- Support import via resource ID using pass-through

#### 2. Write Tests

Create unit tests in the same directory:

```bash
# Test file naming
resource_alkira_<type>_test.go
```

Test naming convention:
```go
func TestAlkira<ResourceType>_<TestName>(t *testing.T) {
    // Test implementation
}
```

#### 3. Create Examples

Add example Terraform configurations to demonstrate resource usage:

```bash
examples/resources/alkira_<resource_name>/resource.tf
```

Example structure:
```hcl
resource "alkira_<resource_name>" "example" {
  name        = "example-resource"
  description = "Example resource configuration"

  # Include all required and common optional parameters
  # Add comments explaining complex configurations
}
```

For data sources:
```bash
examples/data-sources/alkira_<data_source_name>/data-source.tf
```

**Important**: Examples are used by tfplugindocs to generate documentation automatically.

#### 4. Create Import Examples

If the resource supports import (most do), create an import example:

```bash
examples/resources/alkira_<resource_name>/import.sh
```

Content:
```bash
# Import using resource ID
terraform import alkira_<resource_name>.example <resource_id>
```

### Creating Documentation Templates

Documentation templates control how the final documentation is rendered. Templates are only needed for resources that require custom documentation formatting beyond the standard schema documentation.

**When to create a template:**
- Custom usage instructions or warnings
- Complex configuration examples requiring explanation
- Special prerequisites or version requirements
- Non-standard import procedures

**When NOT to create a template:**
- Standard resources with straightforward schemas (tfplugindocs generates these automatically)
- Resources following standard patterns

#### Template Location

Create templates in the `templates/` directory:

```bash
templates/resources/<resource_name>.md.tmpl        # For resources
templates/data-sources/<data_source_name>.md.tmpl  # For data sources
```

#### Template Structure

```markdown
---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

## Example Usage

{{ tffile "examples/resources/alkira_<resource_name>/resource.tf" }}

<!-- Add custom content here -->
<!-- For example, special notes, prerequisites, version requirements -->

{{ .SchemaMarkdown | trimspace }}

## Import

Import is supported using the following syntax:

{{ codefile "shell" .ImportFile }}
```

**Template variables:**
- `{{.Name}}` - Resource name
- `{{.Type}}` - Type (Resource or Data Source)
- `{{.ProviderName}}` - Provider name
- `{{.Description}}` - Resource description from schema
- `{{.SchemaMarkdown}}` - Auto-generated schema documentation
- `{{.ImportFile}}` - Import example file path

**Template functions:**
- `tffile` - Includes Terraform configuration file
- `codefile` - Includes code file with syntax highlighting
- `plainmarkdown` - Converts to plain markdown
- `trimspace` - Removes leading/trailing whitespace
- `prefixlines` - Adds prefix to each line

### Generating Documentation

**CRITICAL: NEVER manually edit files in the `docs/` directory.** All documentation is auto-generated from:
- Resource/data source schema definitions
- Example files in `examples/`
- Template files in `templates/` (when present)

To generate documentation:

```bash
make doc
```

This runs `tfplugindocs generate`, which:
1. Scans all resource and data source schemas
2. Reads example files from `examples/`
3. Applies templates from `templates/` (if they exist)
4. Generates markdown files in `docs/resources/` and `docs/data-sources/`
5. Creates/updates `docs/index.md` for the provider

**Workflow:**
1. Implement resource/data source with proper schema descriptions
2. Create example Terraform configurations
3. (Optional) Create custom template if needed
4. Run `make doc` to generate documentation
5. Review generated files in `docs/` but **never edit them directly**
6. Commit both source files and generated documentation

### Making Changes to alkira-client-go

The provider depends on the private `alkira-client-go` library. When changes are needed to the client library, follow this workflow to ensure code quality:

**IMPORTANT: Never edit files directly in the `/vendor` folder of this repository.**

#### Workflow

1. **Make changes in the alkira-client-go repository**
   ```bash
   cd /path/to/alkira-client-go
   # Make your changes to the client library
   ```

2. **Run linting and tests in alkira-client-go**
   ```bash
   # Run linting to ensure code quality
   make lint

   # Run all tests
   make test

   # Ensure all checks pass before proceeding
   ```

3. **Commit and push changes to alkira-client-go**
   ```bash
   git add .
   git commit -m "Description of changes"
   git push
   ```

4. **Update vendor folder in terraform-provider-alkira**
   ```bash
   cd /path/to/terraform-provider-alkira

   # Update go.mod to reference the latest main branch
   go get github.com/alkiranet/alkira-client-go@main

   # Update vendor directory
   make vendor
   ```

5. **Verify changes in terraform-provider-alkira**
   ```bash
   # Run linting to ensure vendored code meets quality standards
   make lint

   # Run tests to ensure compatibility
   make test

   # Build to verify no compilation errors
   make build
   ```

**Rationale**: This workflow ensures that the vendored code in `/vendor` has been properly linted and tested in its source repository before being integrated. This maintains code quality standards across both repositories.

### Pull Request Guidelines

Before submitting a pull request:

1. **Run tests**: `make test`
2. **Run linting**: `make lint`
3. **Format code**: `make fmt`
4. **Generate docs**: `make doc`
5. **Build successfully**: `make build`
6. **Commit all changes**: Include generated documentation in your commit

Your PR should include:
- Resource/data source implementation
- Unit tests
- Example configurations
- Generated documentation (from `make doc`)
- Updated CHANGELOG (if applicable)

## Project Structure

```
terraform-provider-alkira/
├── .github/
│   └── workflows/          # GitHub Actions workflows
├── acceptance/             # Acceptance test configurations
├── alkira/                 # Provider implementation
│   ├── data_source_*.go    # Data source implementations
│   ├── resource_*.go       # Resource implementations
│   ├── resource_*_helper.go # Helper functions for complex resources
│   ├── resource_*_test.go  # Unit tests
│   ├── helper.go           # Shared utility functions
│   └── provider.go         # Provider configuration
├── bin/                    # Compiled binaries
├── docs/                   # Auto-generated documentation (DO NOT EDIT)
│   ├── data-sources/       # Data source documentation
│   ├── resources/          # Resource documentation
│   └── index.md            # Provider documentation
├── examples/               # Example Terraform configurations
│   ├── data-sources/       # Data source examples
│   ├── resources/          # Resource examples
│   └── provider/           # Provider configuration examples
├── templates/              # Documentation templates (optional)
│   ├── data-sources/       # Data source templates
│   ├── resources/          # Resource templates
│   └── index.md.tmpl       # Provider documentation template
├── vendor/                 # Vendored dependencies
├── .golangci.yml           # Linting configuration
├── GNUmakefile             # Build automation
├── README.md               # This file
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── main.go                 # Provider entry point
```

### Key Directories

- **alkira/**: All provider source code (174 files)
- **examples/**: Terraform configurations used for documentation and testing
- **templates/**: Custom documentation templates (only when needed)
- **docs/**: Auto-generated documentation (never edit manually)
- **acceptance/**: Live integration test configurations

## CI/CD

The repository uses GitHub Actions for continuous integration:

### Workflows

1. **unit-test.yaml** - Runs on every push/PR to main/dev
   - Formats code
   - Runs unit tests
   - Triggers: Push, PR, Manual

2. **lint.yaml** - Code quality checks on every push/PR to main/dev
   - Runs golangci-lint on new/changed code
   - Auto-installs linter in CI
   - Triggers: Push, PR, Manual

3. **acceptance-preprod.yml** - Acceptance tests against pre-production environment
   - Requires Alkira credentials
   - Tests real Terraform configurations

4. **acceptance-prod.yml** - Acceptance tests against production environment
   - Requires Alkira credentials
   - Tests real Terraform configurations

5. **build-latest.yaml** - Builds binaries from latest code
   - Multi-platform builds

6. **build-release.yaml** - Builds and publishes releases
   - Triggered on tags
   - Uses GoReleaser

### Local Testing Before Push

Always run locally before pushing:

```bash
make fmt      # Format code
make lint     # Check code quality
make test     # Run unit tests
make build    # Verify build
make doc      # Generate documentation
```

## Additional Resources

- [Terraform Plugin SDK](https://developer.hashicorp.com/terraform/plugin/sdkv2)
- [Alkira Documentation](https://registry.terraform.io/providers/alkiranet/alkira/latest/docs)
- [Alkira Client Go](https://github.com/alkiranet/alkira-client-go) (private)


## Support

For issues and questions:
- GitHub Issues: [terraform-provider-alkira/issues](https://github.com/alkiranet/terraform-provider-alkira/issues)
- Alkira Support Portal

---

**Maintained by Alkira Engineering Team**
