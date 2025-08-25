package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScenarioResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccScenarioResourceConfig("example"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("make_scenario.test", "name", "Test Scenario example"),
					resource.TestCheckResourceAttr("make_scenario.test", "description", "Test scenario description"),
					resource.TestCheckResourceAttr("make_scenario.test", "active", "true"),
					resource.TestCheckResourceAttrSet("make_scenario.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "make_scenario.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccScenarioResourceConfig("updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("make_scenario.test", "name", "Test Scenario updated"),
				),
			},
		},
	})
}

func testAccScenarioResourceConfig(suffix string) string {
	return `
resource "make_scenario" "test" {
  name        = "Test Scenario ` + suffix + `"
  description = "Test scenario description"
  active      = true
}
`
}

func TestAccConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccConnectionResourceConfig("example"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("make_connection.test", "name", "Test Connection example"),
					resource.TestCheckResourceAttr("make_connection.test", "app_name", "gmail"),
					resource.TestCheckResourceAttrSet("make_connection.test", "id"),
					resource.TestCheckResourceAttrSet("make_connection.test", "verified"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "make_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConnectionResourceConfig(suffix string) string {
	return `
resource "make_connection" "test" {
  name     = "Test Connection ` + suffix + `"
  app_name = "gmail"
}
`
}

func TestAccWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWebhookResourceConfig("example"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("make_webhook.test", "name", "Test Webhook example"),
					resource.TestCheckResourceAttr("make_webhook.test", "active", "true"),
					resource.TestCheckResourceAttrSet("make_webhook.test", "id"),
					resource.TestCheckResourceAttrSet("make_webhook.test", "url"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "make_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWebhookResourceConfig(suffix string) string {
	return `
resource "make_webhook" "test" {
  name   = "Test Webhook ` + suffix + `"
  active = true
}
`
}