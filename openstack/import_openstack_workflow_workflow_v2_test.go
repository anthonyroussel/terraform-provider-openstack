package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowV2Workflow_importBasic(t *testing.T) {
	resourceName := "openstack_workflow_workflow_v2.workflow_1"

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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
