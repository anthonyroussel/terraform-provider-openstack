package openstack

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
)

func TestAccWorkflowV2CronTrigger_basic(t *testing.T) {
	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowV2WorkflowCreate()
		if err != nil {
			t.Fatal(err)
		}
		workflowID = workflow.ID
		defer testAccWorkflowV2WorkflowDelete(workflow.ID) //nolint:errcheck
	}

	var crontrigger crontriggers.CronTrigger

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckWorkflow(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2CronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2CronTriggerBasic(workflowID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2CronTriggerExists("openstack_workflow_cron_trigger_v2.cron_trigger_1", &crontrigger),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "id"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "name", "cron_trigger_1"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_id", workflowID),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "pattern", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_input.%", "2"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "created_at"),
				),
			},
		},
	})
}

func testAccCheckWorkflowV2CronTriggerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	workflowClient, err := config.WorkflowV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_workflow_cron_trigger_v2" {
			continue
		}

		_, err := crontriggers.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("CronTrigger still exists")
		}
	}

	return nil
}

func testAccCheckWorkflowV2CronTriggerExists(n string, crontrigger *crontriggers.CronTrigger) resource.TestCheckFunc {
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

		found, err := crontriggers.Get(context.TODO(), workflowClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("CronTrigger not found")
		}

		*crontrigger = *found

		return nil
	}
}

func testAccWorkflowV2CronTriggerBasic(workflowID string) string {
	return fmt.Sprintf(`
resource "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name        = "cron_trigger_1"
  workflow_id = "%s"
  pattern     = "0 5 * * *"

  workflow_input = {
    my_arg1 = "val1"
    my_arg2 = "val2"
  }
}
`, workflowID)
}
