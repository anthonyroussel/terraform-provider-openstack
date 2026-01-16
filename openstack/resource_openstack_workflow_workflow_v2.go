package openstack

import (
	"context"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/verify"
)

func resourceWorkflowWorkflowV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowWorkflowV2Create,
		ReadContext:   resourceWorkflowWorkflowV2Read,
		DeleteContext: resourceWorkflowWorkflowV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Input
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public",
					"private",
				}, false),
			},

			"definition": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidStringIsYAML),
				// TODO: can be updated through the API,
				// FIXME: only one workflow at a time is supported at this time.
			},

			// Computed
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"input": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWorkflowWorkflowV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	definition, err := verify.NormalizeYAMLString(d.Get("definition").(string))
	if err != nil {
		return diag.Errorf("Unable to parse definition: %s", err)
	}

	createOpts := workflows.CreateOpts{
		Scope:      d.Get("scope").(string),
		Namespace:  d.Get("namespace").(string),
		Definition: strings.NewReader(definition),
	}

	workflowList, err := workflows.Create(ctx, workflowClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create openstack_workflow_workflow_v2: %s", err)
	}

	workflow := workflowList[len(workflowList)-1]

	d.SetId(workflow.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("definition", definition)

	return resourceWorkflowWorkflowV2Read(ctx, d, meta)
}

func resourceWorkflowWorkflowV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	workflow, err := workflows.Get(ctx, workflowClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_workflow_workflow_v2"))
	}

	tflog.Debug(ctx, "Retrieved openstack_workflow_workflow_v2 resource", map[string]any{
		"id":       d.Id(),
		"workflow": workflow,
	})

	d.Set("region", GetRegion(d, config))
	d.Set("name", workflow.Name)
	d.Set("scope", workflow.Scope)
	d.Set("namespace", workflow.Namespace)
	d.Set("definition", workflow.Definition)
	d.Set("input", workflow.Input)
	d.Set("tags", workflow.Tags)
	d.Set("project_id", workflow.ProjectID)

	if err := d.Set("created_at", workflow.CreatedAt.Format(time.RFC3339)); err != nil {
		tflog.Debug(ctx, "Unable to set created_at for openstack_workflow_workflow_v2", map[string]any{
			"id":  workflow.ID,
			"err": err,
		})
	}

	return nil
}

func resourceWorkflowWorkflowV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	err = workflows.Delete(ctx, workflowClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_workflow_workflow_v2"))
	}

	return nil
}
