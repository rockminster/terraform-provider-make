terraform {
  required_providers {
    make = {
      source = "registry.terraform.io/rockminster/make"
      version = "~> 0.1.0"
    }
  }
}

provider "make" {
  api_token = var.make_api_token
  # base_url = "https://api.make.com/"  # Optional, defaults to this
}

variable "make_api_token" {
  description = "Make.com API token"
  type        = string
  sensitive   = true
}

# Create a scenario
resource "make_scenario" "example" {
  name        = "Example Scenario"
  description = "This is an example scenario created by Terraform"
  active      = true
  team_id     = "your-team-id"
}

# Read an existing scenario
data "make_scenario" "existing" {
  id = "existing-scenario-id"
}

# Output the scenario details
output "scenario_id" {
  value = make_scenario.example.id
}

output "existing_scenario_name" {
  value = data.make_scenario.existing.name
}