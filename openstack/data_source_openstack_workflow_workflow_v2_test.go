package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowV2WorkflowDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWorkflow(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2WorkflowDataSourceBasic,
			},
			{
				Config: testAccWorkflowV2WorkflowDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "name", "my_workflow_datasource"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "input", "my_arg1, my_arg2"),
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

    my_workflow_datasource:
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

func testAccWorkflowV2WorkflowDataSourceSource() string {
	return fmt.Sprintf(`
%s
data "openstack_workflow_workflow_v2" "workflow_1" {
	name      = openstack_workflow_workflow_v2.workflow_1.name
	namespace = openstack_workflow_workflow_v2.workflow_1.namespace
}
`, testAccWorkflowV2WorkflowDataSourceBasic)
}
