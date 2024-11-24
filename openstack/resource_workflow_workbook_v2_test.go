package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"
)

func TestAccWorkflowV2Workbook_basic(t *testing.T) {
	var workbook workbooks.Workbook

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckWorkflow(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2WorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccWorkflowV2WorkbookBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2WorkbookExists("openstack_workflow_workbook_v2.workbook_1", &workbook),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workbook_v2.workbook_1", "name", "my_workbook_resource"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workbook_v2.workbook_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workbook_v2.workbook_1", "definition"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workbook_v2.workbook_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workbook_v2.workbook_1", "scope", "private"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workbook_v2.workbook_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workbook_v2.workbook_1", "created_at"),
				),
			},
		},
	})
}

func testAccCheckWorkflowV2WorkbookDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	workflowClient, err := config.WorkflowV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_workflow_workbook_v2" {
			continue
		}

		_, err := workbooks.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Workbook still exists")
		}
	}

	return nil
}

func testAccCheckWorkflowV2WorkbookExists(n string, workbook *workbooks.Workbook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		workflowClient, err := config.WorkflowV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack workflow client: %s", err)
		}

		found, err := workbooks.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Workbook not found")
		}

		*workbook = *found

		return nil
	}
}

const TestAccWorkflowV2WorkbookBasic = `
resource "openstack_workflow_workbook_v2" "workbook_1" {
  namespace = "my_namespace"
  scope     = "private"
  definition = <<EOF
    version: '2.0'

    name: my_workbook
    description: My workbook
    tags:
      - test

    workflows:
      test:
        description: Simple workflow example
        type: direct
        input:
          - msg

        tasks:
          test:
            action: std.echo output="<%% $.msg %%>
  EOF
}
`
