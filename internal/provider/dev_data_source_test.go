package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devops-bootcamp_dev" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devops-bootcamp_dev.test", "dev.#", "2"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_dev.test", "dev.0.name", "dev_ferrets"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_dev.test", "dev.0.engineers.#", "1"),
					resource.TestCheckResourceAttr("data.devops-bootcamp_dev.test", "dev.0.engineers.0.email", "ryan@ferrets.com"),
				),
			},
		},
	})
}
