package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccWorkflowV2CronTrigger_basic(t *testing.T) {
	var crontrigger crontriggers.CronTrigger

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2CronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccWorkflowV2CronTriggerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2CronTriggerExists("openstack_workflow_cron_trigger_v2.trigger_1", &crontrigger),
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

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("CronTrigger not found")
		}

		*crontrigger = *found

		return nil
	}
}

const TestAccWorkflowV2CronTriggerBasic = `
resource "openstack_workflow_cron_trigger_v2" "my_cron_trigger" {
  name        = "my_cron_trigger"
  workflow_id = "428ce188-9881-4784-b36d-ef20659feced"
  pattern     = "0 5 * * *"

  workflow_input = {
    param1 = "val1"
    param2 = "val2"
  }
}
`
