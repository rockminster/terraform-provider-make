# Terraform Provider for Make.com

A Terraform provider for managing Make.com (formerly Integromat) resources using the terraform-plugin-framework.

## Features

- **Modern Architecture**: Built using terraform-plugin-framework for enhanced performance and developer experience
- **Make.com Integration**: Manage Make.com scenarios, connections, and other resources
- **Terraform Best Practices**: Follows Terraform provider development best practices

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```hcl
terraform {
  required_providers {
    make = {
      source = "registry.terraform.io/rockminster/make"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Make.com Provider
provider "make" {
  api_token = var.make_api_token
  # base_url = "https://api.make.com/"  # Optional
}

# Create a scenario
resource "make_scenario" "example" {
  name        = "My Terraform Scenario"
  description = "A scenario managed by Terraform"
  active      = true
  team_id     = "your-team-id"
}

# Read an existing scenario
data "make_scenario" "existing" {
  id = "existing-scenario-id"
}
```

## Provider Configuration

### Environment Variables

You can configure the provider using environment variables:

- `MAKE_API_TOKEN` - Make.com API token
- `MAKE_BASE_URL` - Base URL for Make.com API (defaults to https://api.make.com/)

### Provider Block

```hcl
provider "make" {
  api_token = "your-api-token"  # Can also use MAKE_API_TOKEN env var
  base_url  = "https://api.make.com/"  # Optional
}
```

## Available Resources

### make_scenario

Manages Make.com scenarios.

#### Example Usage

```hcl
resource "make_scenario" "example" {
  name        = "My Scenario"
  description = "Example scenario"
  active      = true
  team_id     = "team-123"
}
```

#### Arguments

- `name` (Required) - Name of the scenario
- `description` (Optional) - Description of the scenario
- `active` (Optional) - Whether the scenario is active
- `team_id` (Optional) - Team ID where the scenario belongs

#### Attributes

- `id` - Scenario identifier

## Available Data Sources

### make_scenario

Reads information about an existing Make.com scenario.

#### Example Usage

```hcl
data "make_scenario" "example" {
  id = "scenario-id-123"
}
```

#### Arguments

- `id` (Required) - Scenario identifier

#### Attributes

- `name` - Name of the scenario
- `description` - Description of the scenario
- `active` - Whether the scenario is active
- `team_id` - Team ID where the scenario belongs

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Development Commands

- `make build` - Build the provider binary
- `make install` - Install the provider locally for testing
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests
- `make docs` - Generate documentation
- `make lint` - Run linters
- `make fmt` - Format code
- `make check` - Run all checks (fmt, lint, test)
- `make clean` - Clean build artifacts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run `make check` to ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Make.com API Reference

This provider is designed to work with the Make.com API. For more information about available API endpoints and functionality, refer to the [Make.com API documentation](https://www.make.com/en/api-documentation).

## Roadmap

- [ ] Implement comprehensive scenario management
- [ ] Add connection resource management
- [ ] Add team and organization management
- [ ] Add webhook management
- [ ] Add data store management
- [ ] Add comprehensive error handling and validation
- [ ] Add comprehensive test coverage
