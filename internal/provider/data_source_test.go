package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScenarioDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScenarioDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.make_scenario.test", "name", "Test Scenario"),
					resource.TestCheckResourceAttr("data.make_scenario.test", "description", "Test scenario description"),
					resource.TestCheckResourceAttr("data.make_scenario.test", "active", "true"),
				),
			},
		},
	})
}

func testAccScenarioDataSourceConfig() string {
	return `
resource "make_scenario" "test" {
  name        = "Test Scenario"
  description = "Test scenario description"
  active      = true
}

data "make_scenario" "test" {
  id = make_scenario.test.id
}
`
}

func TestAccConnectionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.make_connection.test", "name", "Test Connection"),
					resource.TestCheckResourceAttr("data.make_connection.test", "app_name", "gmail"),
					resource.TestCheckResourceAttrSet("data.make_connection.test", "verified"),
				),
			},
		},
	})
}

func testAccConnectionDataSourceConfig() string {
	return `
resource "make_connection" "test" {
  name     = "Test Connection"
  app_name = "gmail"
}

data "make_connection" "test" {
  id = make_connection.test.id
}
`
}

func TestAccTeamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.make_team.test", "name", "Test Team"),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig() string {
	return `
resource "make_team" "test" {
  name = "Test Team"
}

data "make_team" "test" {
  id = make_team.test.id
}
`
}

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.make_organization.test", "name", "Test Organization"),
				),
			},
		},
	})
}

func testAccOrganizationDataSourceConfig() string {
	return `
resource "make_organization" "test" {
  name = "Test Organization"
}

data "make_organization" "test" {
  id = make_organization.test.id
}
`
}

func TestAccDataStoreDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataStoreDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.make_data_store.test", "name", "Test Data Store"),
					resource.TestCheckResourceAttr("data.make_data_store.test", "description", "Test data store description"),
				),
			},
		},
	})
}

func testAccDataStoreDataSourceConfig() string {
	return `
resource "make_data_store" "test" {
  name        = "Test Data Store"
  description = "Test data store description"
}

data "make_data_store" "test" {
  id = make_data_store.test.id
}
`
}
