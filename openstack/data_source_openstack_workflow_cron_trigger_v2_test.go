package openstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowV2CronTriggerDataSource_basic(t *testing.T) {
	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowV2WorkflowCreate()
		if err != nil {
			t.Fatal(err)
		}
		workflowID = workflow.ID
		defer testAccWorkflowV2WorkflowDelete(workflow.ID) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWorkflow(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2CronTriggerDataSourceBasic(workflowID),
			},
			{
				Config: testAccWorkflowV2CronTriggerDataSourceSource(workflowID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "name", "my_workflow"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_id", workflowID),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "pattern", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_params", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_input", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "remaining_executions", "0"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "first_execution_time", "2022-01-01 01:01:01"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "created_at"),
				),
			},
		},
	})
}

func testAccWorkflowV2CronTriggerDataSourceBasic(workflowID string) string {
	return fmt.Sprintf(`
%s

data "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
	name = openstack_workflow_cron_trigger_v2.cron_trigger_1.name
}
`, testAccWorkflowV2CronTriggerDataSourceSource(workflowID))
}

func testAccWorkflowV2CronTriggerDataSourceSource(workflowID string) string {
	return fmt.Sprintf(`
resource "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name        = "cron_trigger_1"
  workflow_id = "%s"
  pattern     = "0 5 * * *"

  workflow_input = {
    my_arg1 = "value1"
    my_arg2 = "value2"
  }
}
`, workflowID)
}
