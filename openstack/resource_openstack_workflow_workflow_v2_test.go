package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
)

func TestAccWorkflowV2Workflow_basic(t *testing.T) {
	var workflow workflows.Workflow

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckWorkflow(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccWorkflowV2WorkflowBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2WorkflowExists("openstack_workflow_workflow_v2.workflow_1", &workflow),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "name", "workflow_echo_resource"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "namespace", "some-namespace"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "input", "msg"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workflow_v2.workflow_1", "definition"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "scope", "private"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workflow_v2.workflow_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_workflow_v2.workflow_1", "created_at"),
				),
			},
		},
	})
}

func testAccCheckWorkflowV2WorkflowDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	workflowClient, err := config.WorkflowV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_workflow_workflow_v2" {
			continue
		}

		_, err := workflows.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Workflow still exists")
		}
	}

	return nil
}

func testAccCheckWorkflowV2WorkflowExists(n string, workflow *workflows.Workflow) resource.TestCheckFunc {
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

		found, err := workflows.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Workflow not found")
		}

		*workflow = *found

		return nil
	}
}

const TestAccWorkflowV2WorkflowBasic = `
resource "openstack_workflow_workflow_v2" "workflow_1" {
  namespace = "some-namespace"
  scope     = "private"
  definition = <<EOF
    version: '2.0'

    workflow_echo_resource:
      description: Simple workflow example
      type: direct
      tags:
        - echo

      input:
        - msg

      tasks:
        test:
          action: std.echo output="<% $.msg %>"
  EOF
}
`
