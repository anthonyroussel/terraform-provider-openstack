// Copyright (c) HashiCorp, Inc.

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
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccWorkflowV2WorkflowBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2WorkflowExists("openstack_workflow_workflow_v2.workflow_1", &workflow),
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

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Workflow not found")
		}

		*workflow = *found

		return nil
	}
}

const TestAccWorkflowV2WorkflowBasic = `
resource "openstack_workflow_workflow_v2" "workflow_1" {
	namespace = "my_namespace"
	scope     = "private"
	definition = <<EOF
    version: '2.0'

    my_workflow_resource:
      description: Simple echo example

      input:
        - my_arg1
        - my_arg2

      tags:
        - echo

      tasks:
        echo:
          action: std.echo
          input:
            output:
              my_arg1: <% $.my_arg1 %>
              my_arg2: <% $.my_arg2 %>
	EOF
}
`
