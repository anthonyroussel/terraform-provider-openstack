package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2WorkflowBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2WorkflowExists(t.Context(), "openstack_workflow_workflow_v2.workflow_1", &workflow),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "name", "hello_workflow"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_workflow_v2.workflow_1", "input", "message"),
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

func testAccCheckWorkflowV2WorkflowDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		workflowClient, err := config.WorkflowV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack workflow client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_workflow_workflow_v2" {
				continue
			}

			_, err := workflows.Get(ctx, workflowClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Workflow still exists")
			}
		}

		return nil
	}
}

func testAccCheckWorkflowV2WorkflowExists(ctx context.Context, n string, workflow *workflows.Workflow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		workflowClient, err := config.WorkflowV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack workflow client: %w", err)
		}

		found, err := workflows.Get(ctx, workflowClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Workflow not found")
		}

		*workflow = *found

		return nil
	}
}

const testAccWorkflowV2WorkflowBasic = `
resource "openstack_workflow_workflow_v2" "workflow_1" {
  namespace = "my_namespace"
  scope     = "private"
  definition = <<EOF
    version: '2.0'

    hello_workflow:
      description: Simple echo example

      input:
        - message

      tags:
        - echo

      tasks:
        echo:
          action: std.echo
          input:
            output:
              my_message: <% $.message %>
  EOF
}
`
