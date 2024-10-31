package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowV2CronTrigger_importBasic(t *testing.T) {
	resourceName := "openstack_workflow_cron_trigger_v2.my_cron_trigger"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2CronTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccWorkflowV2CronTriggerBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
