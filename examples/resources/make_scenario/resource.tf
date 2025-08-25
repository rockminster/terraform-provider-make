resource "make_scenario" "example" {
  name        = "My Terraform Scenario"
  description = "A scenario managed by Terraform"
  active      = true
  team_id     = "team-123"
}