package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "test" {
  name  = "John Doe"
  email = "john.doe@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "name", "John Doe"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "email", "john.doe@example.com"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer-resource.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_engineer-resource.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "test" {
  name  = "Jane Doe"
  email = "jane.doe@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated attributes
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "name", "Jane Doe"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.test", "email", "jane.doe@example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
