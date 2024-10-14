package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devops-bootcamp_engineer" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineer.#", "3"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineer.0.id", "UCS24"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineer.0.name", "Ryan"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "engineer.0.email", "ryan@ferrets.com"),
				),
			},
		},
	})
}
