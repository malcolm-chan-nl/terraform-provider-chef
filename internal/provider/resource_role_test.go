package provider

import (
	"fmt"
	"reflect"
	"testing"

	chefc "github.com/go-chef/chef"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRole_basic(t *testing.T) {
	var role chefc.Role

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccRoleCheckDestroy(&role),
		Steps: []resource.TestStep{
			{
				Config: testSuffixRender(testAccRoleConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccRoleCheckExists("chef_role.test", &role),
					func(s *terraform.State) error {

						if expected := "terraform-acc-test-basic-" + testSuffix; role.Name != expected {
							return fmt.Errorf("wrong name; expected %v, got %v", expected, role.Name)
						}
						if expected := "Terraform Acceptance Tests"; role.Description != expected {
							return fmt.Errorf("wrong description; expected %v, got %v", expected, role.Description)
						}

						expectedRunListStrings := []string{
							"recipe[terraform@1.0.0]",
							"recipe[consul]",
							"role[foo]",
						}
						expectedRunList := chefc.RunList(expectedRunListStrings)
						if !reflect.DeepEqual(role.RunList, expectedRunList) {
							return fmt.Errorf("wrong runlist; expected %#v, got %#v", expectedRunList, role.RunList)
						}

						expectedAttributes := map[string]interface{}{
							"terraform_acc_test": true,
						}
						if !reflect.DeepEqual(role.DefaultAttributes, expectedAttributes) {
							return fmt.Errorf("wrong default attributes; expected %#v, got %#v", expectedAttributes, role.DefaultAttributes)
						}
						if !reflect.DeepEqual(role.OverrideAttributes, expectedAttributes) {
							return fmt.Errorf("wrong override attributes; expected %#v, got %#v", expectedAttributes, role.OverrideAttributes)
						}
						if !reflect.DeepEqual(role.EnvRunList, expectedAttributes) {
							return fmt.Errorf("wrong environment run list; expected %#v, got %#v", expectedAttributes, role.EnvRunList)
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccRoleCheckExists(rn string, role *chefc.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role id not set")
		}

		client := testAccProvider.Meta().(*chefClient)
		gotRole, err := client.Roles.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting role: %s", err)
		}

		*role = *gotRole

		return nil
	}
}

func testAccRoleCheckDestroy(role *chefc.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*chefClient)
		_, err := client.Roles.Get(role.Name)
		if err == nil {
			return fmt.Errorf("role still exists")
		}
		if _, ok := err.(*chefc.ErrorResponse); !ok {
			// A more specific check is tricky because Chef Server can return
			// a few different error codes in this case depending on which
			// part of its stack catches the error.
			return fmt.Errorf("got something other than an HTTP error (%v) when getting role", err)
		}

		return nil
	}
}

const testAccRoleConfig_basic = `
resource "chef_role" "test" {
  name = "terraform-acc-test-basic-{{.}}"
  description = "Terraform Acceptance Tests"
  default_attributes_json = <<EOT
{
     "terraform_acc_test": true
}
EOT
  override_attributes_json = <<EOT
{
     "terraform_acc_test": true
}
EOT
  env_run_list_json = <<EOT
{
     "terraform_acc_test": true
}
EOT

  run_list = ["terraform@1.0.0", "recipe[consul]", "role[foo]"]
}
`
