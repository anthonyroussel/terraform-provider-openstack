package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		CheckDestroy:      testAccCheckWorkflowV2WorkflowDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2WorkflowBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
