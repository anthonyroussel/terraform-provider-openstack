package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkflowV2WorkflowDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWorkflow(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2WorkflowDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "name", "hello_workflow"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "input", "message"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "definition"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "scope", "private"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "created_at"),
				),
			},
		},
	})
}

const testAccWorkflowV2WorkflowDataSourceBasic = `
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

data "openstack_workflow_workflow_v2" "workflow_1" {
	name      = openstack_workflow_workflow_v2.workflow_1.name
	namespace = openstack_workflow_workflow_v2.workflow_1.namespace
}
`
