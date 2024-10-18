package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "engineer1" {
  name  = "John Doe"
  email = "john.doe@example.com"
}

resource "devops-bootcamp_engineer-resource" "engineer2" {
  name  = "Jane Smith"
  email = "jane.smith@example.com"
}

resource "devops-bootcamp_dev_resource" "test" {
  name = "Test Dev Group"
  engineers = [
    devops-bootcamp_engineer-resource.engineer1,
    devops-bootcamp_engineer-resource.engineer2
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes for engineer1
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.engineer1", "name", "John Doe"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.engineer1", "email", "john.doe@example.com"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer-resource.engineer1", "id"),
					// Verify attributes for engineer2
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.engineer2", "name", "Jane Smith"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer-resource.engineer2", "email", "jane.smith@example.com"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer-resource.engineer2", "id"),
					// Verify attributes for dev resource
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev_resource.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_dev_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer-resource" "engineer1" {
  name  = "Updated John Doe"
  email = "john.doe@example.com"
}

resource "devops-bootcamp_engineer-resource" "engineer2" {
  name  = "Updated Jane Smith"
  email = "jane.smith@example.com"
}

resource "devops-bootcamp_dev_resource" "test" {
  name = "Updated Test Dev Group"
  engineers = [
    devops-bootcamp_engineer-resource.engineer1,
    devops-bootcamp_engineer-resource.engineer2
  ]
}
`,
			},
		},
	})
}
