# Terraform Provider for Make.com

This is a Go-based Terraform provider for Make.com (formerly Integromat) built using the terraform-plugin-framework for modern provider development.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Prerequisites and Environment Setup
- Ensure Go >= 1.21 is installed. Current environment has Go 1.24.6 which meets requirements.
- Add Go binary path to PATH: `export PATH=$PATH:~/go/bin`
- Install Terraform for documentation generation:
  ```bash
  cd /tmp
  curl -LO https://releases.hashicorp.com/terraform/1.9.8/terraform_1.9.8_linux_amd64.zip
  unzip terraform_1.9.8_linux_amd64.zip
  sudo mv terraform /usr/local/bin/
  ```

### Build and Development Workflow
- Build the provider: `make build` -- takes ~25 seconds on first run with dependency downloads, ~3 seconds subsequent runs. NEVER CANCEL. Set timeout to 60+ seconds for safety.
- Format code: `make fmt` -- takes ~0.2 seconds
- Lint code: `make lint` -- takes ~8 seconds after golangci-lint installation. NEVER CANCEL. Set timeout to 120+ seconds.
  - First run installs golangci-lint to ~/go/bin/ - ensure PATH includes this directory
  - May show lint errors for unused test helper functions (expected in template)
- Run unit tests: `make test` -- takes ~4 seconds. NEVER CANCEL. Set timeout to 60+ seconds.
  - Note: Currently no actual test implementations exist
- Generate documentation: `make docs` -- takes ~8 seconds. NEVER CANCEL. Set timeout to 120+ seconds.
  - Requires Terraform to be installed
- Install provider locally: `make install` -- takes ~1 second
- Run acceptance tests: `make testacc` -- takes ~4 seconds. NEVER CANCEL. Set timeout to 180+ minutes.
  - Note: Currently no actual acceptance test implementations exist
  - When implemented, these tests may take much longer as they create real resources
- Run all checks: `make check` -- runs fmt, lint, and test sequentially. NEVER CANCEL. Set timeout to 240+ seconds.
- Clean build artifacts: `make clean`

### Testing Your Changes
After making code changes, always run this validation sequence:
1. Build and install the provider:
   ```bash
   make build
   make install
   ```
2. Test with a local Terraform configuration:
   ```bash
   mkdir -p /tmp/test-config
   cd /tmp/test-config
   cat > main.tf << 'EOF'
   terraform {
     required_providers {
       make = {
         source = "registry.terraform.io/rockminster/make"
         version = "~> 0.1.0"
       }
     }
   }
   
   provider "make" {
     api_token = "test-token"
     base_url  = "https://api.make.com/"
   }
   
   resource "make_scenario" "test" {
     name        = "Test Scenario"
     description = "A test scenario"
     active      = false
     team_id     = "test-team"
   }
   
   resource "make_connection" "test" {
     name     = "Test Connection"
     app_name = "http"
     team_id  = "test-team"
   }
   EOF
   terraform init
   terraform plan
   ```
3. Verify the plan shows the expected resource configuration
4. Run the complete check suite: `make check`

## Key Project Structure

### Repository Root
```
.
├── README.md              # Main documentation
├── LICENSE                # MIT License
├── Makefile              # Build automation
├── main.go               # Provider entry point
├── version.go            # Version information
├── tools.go              # Build tool dependencies
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── .goreleaser.yml       # Release configuration
├── .gitignore            # Git ignore rules
├── internal/             # Internal provider code
├── examples/             # Terraform usage examples
└── docs/                 # Generated documentation (after make docs)
```

### Core Provider Code
- `internal/provider/provider.go` - Main provider implementation
- `internal/provider/client.go` - Make.com API client implementation
- `internal/provider/scenario_resource.go` - Scenario resource implementation
- `internal/provider/scenario_data_source.go` - Scenario data source implementation
- `internal/provider/connection_resource.go` - Connection resource implementation  
- `internal/provider/connection_data_source.go` - Connection data source implementation
- `internal/provider/webhook_resource.go` - Webhook resource implementation
- `internal/provider/provider_test.go` - Provider test helpers
- `internal/provider/client_test.go` - API client tests
- `internal/provider/resource_test.go` - Resource acceptance tests

### Example Configurations
- `examples/provider/provider.tf` - Full provider usage example
- `examples/resources/make_scenario/resource.tf` - Scenario resource example
- `examples/resources/make_connection/resource.tf` - Connection resource example
- `examples/resources/make_webhook/resource.tf` - Webhook resource example
- `examples/data-sources/make_scenario/data-source.tf` - Scenario data source example
- `examples/data-sources/make_connection/data-source.tf` - Connection data source example

## Common Issues and Solutions

### Lint Failures
- golangci-lint not in PATH: Run `export PATH=$PATH:~/go/bin` first
- Unused test functions: Expected for template code in `provider_test.go`. To avoid general lint failures, configure golangci-lint to ignore unused code warnings in test files by updating `.golangci.yml` or using `//nolint:unused` comments on the relevant functions.
- Do not ignore general lint failures. Instead, suppress only the unused function warnings in `provider_test.go` until actual tests are implemented.

### Documentation Generation
- Terraform not found: Install Terraform as shown in Prerequisites section
- Missing docs: Run `make docs` to regenerate

### Provider Testing
- "No test files" messages: Expected when no actual tests are implemented yet
- Provider not found during terraform init: Run `make install` first

## Validation Requirements

Before committing changes, always:
1. Run `make fmt && make test` to ensure code formatting and tests pass (skip lint if unused functions cause failures)
2. Test provider installation and basic Terraform workflow as shown above
3. If adding new resources or data sources, add corresponding examples
4. If modifying provider schema, regenerate docs with `make docs`
5. For clean code, address lint issues except for expected unused test helper functions

## API Integration Notes

This provider integrates with the Make.com API. For actual functionality testing:
- Set `MAKE_API_TOKEN` environment variable with valid API token
- Set `MAKE_BASE_URL` if using different API endpoint (defaults to https://api.make.com/)
- Be aware that acceptance tests create real resources and may incur costs

## Development Timeouts

CRITICAL: When running any make commands, use these minimum timeout values:
- `make build`: 60+ seconds (initial build ~25s, subsequent ~1s, 50% safety buffer)
- `make lint`: 120+ seconds (installs golangci-lint on first run, currently fails on unused functions)
- `make test`: 60+ seconds (minimal currently, may grow with test additions)  
- `make docs`: 120+ seconds (includes Terraform operations)
- `make testacc`: 180+ minutes (acceptance tests can be very long-running when implemented)
- `make check`: 240+ seconds (combination of fmt, lint, test - may fail on lint currently)

NEVER CANCEL builds or long-running commands. If a command appears stuck, wait the full timeout period before investigating.