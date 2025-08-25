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

# Create a connection
resource "make_connection" "gmail" {
  name     = "Gmail Connection"
  app_name = "gmail"
  team_id  = "your-team-id"
}

# Create a webhook
resource "make_webhook" "incoming" {
  name    = "Incoming Webhook"
  team_id = "your-team-id"
  active  = true
}

# Read an existing scenario
data "make_scenario" "existing" {
  id = "existing-scenario-id"
}

# Read an existing connection
data "make_connection" "existing_conn" {
  id = "existing-connection-id"
}

# Output the resource details
output "scenario_id" {
  value = make_scenario.example.id
}

output "connection_id" {
  value = make_connection.gmail.id
}

output "webhook_id" {
  value = make_webhook.incoming.id
}

output "webhook_url" {
  value = make_webhook.incoming.url
}

output "existing_scenario_name" {
  value = data.make_scenario.existing.name
}

output "existing_connection_verified" {
  value = data.make_connection.existing_conn.verified
}