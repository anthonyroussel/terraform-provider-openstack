package openstack

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
)

func TestAccWorkflowWorkflowV2DataSource_basic(t *testing.T) {
	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowWorkflowV2Create()
		if err != nil {
			t.Fatal(err)
		}
		workflowID = workflow.ID
		defer testAccWorkflowWorkflowV2Delete(workflow.ID) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWorkflow(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowWorkflowV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "id", workflowID),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "name", "my_workflow"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "input", "my_arg1, my_arg2"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_v2.workflow_1", "definition"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_v2.workflow_1", "scope", "private"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_v2.workflow_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_v2.workflow_1", "created_at"),
				),
			},
		},
	})
}

func testAccWorkflowWorkflowV2Delete(workflowID string) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	client, err := config.WorkflowV2Client(context.TODO(), osRegionName)
	if err != nil {
		return err
	}

	err = workflows.Delete(context.TODO(), client, workflowID).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

func testAccWorkflowWorkflowV2Create() (*workflows.Workflow, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := config.WorkflowV2Client(context.TODO(), osRegionName)
	if err != nil {
		return nil, err
	}

	createWorkflowOpts := workflows.CreateOpts{
		Scope:      "private",
		Namespace:  "my_namespace",
		Definition: strings.NewReader(testAccWorkflowWorkflowV2DataSourceBasicDefinition),
	}

	workflows, err := workflows.Create(context.TODO(), client, createWorkflowOpts).Extract()
	if err != nil {
		return nil, err
	}

	workflow := workflows[len(workflows)-1]

	return &workflow, nil
}

const testAccWorkflowWorkflowV2DataSourceBasic = `
data "openstack_workflow_v2" "workflow_1" {
	name      = "my_workflow"
	namespace = "my_namespace"
}
`

const testAccWorkflowWorkflowV2DataSourceBasicDefinition = `
version: '2.0'

my_workflow:
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
`
