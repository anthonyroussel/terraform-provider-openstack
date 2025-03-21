package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
)

func resourceWorkflowCronTriggerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowCronTriggerV2Create,
		ReadContext:   resourceWorkflowCronTriggerV2Read,
		DeleteContext: resourceWorkflowCronTriggerV2Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"workflow_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ConflictsWith: []string{"workflow_name"},
			},

			"workflow_name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ConflictsWith: []string{"workflow_id"},
			},

			"workflow_input": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"workflow_params": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"pattern": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remaining_executions": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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

func resourceWorkflowCronTriggerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	createOpts := crontriggers.CreateOpts{
		Name:                d.Get("name").(string),
		WorkflowID:          d.Get("workflow_id").(string),
		WorkflowName:        d.Get("workflow_name").(string),
		WorkflowInput:       d.Get("workflow_input").(map[string]any),
		WorkflowParams:      d.Get("workflow_params").(map[string]any),
		Pattern:             d.Get("pattern").(string),
		RemainingExecutions: d.Get("remaining_executions").(int),
	}

	if createOpts.WorkflowID == "" && createOpts.WorkflowName == "" {
		return diag.Errorf(`Either "workflow_id" or "workflow_name" must be set`)
	}

	log.Printf("[DEBUG] openstack_workflow_cron_trigger_v2 create options: %#v", createOpts)

	crontrigger, err := crontriggers.Create(ctx, workflowClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create openstack_workflow_cron_trigger_v2: %s", err)
	}

	d.SetId(crontrigger.ID)
	d.Set("name", d.Get("name").(string))
	d.Set("workflow_id", d.Get("workflow_id").(string))
	d.Set("workflow_name", d.Get("workflow_name").(string))
	d.Set("workflow_input", d.Get("workflow_input").(map[string]any))
	d.Set("workflow_params", d.Get("workflow_params").(map[string]any))
	d.Set("pattern", d.Get("pattern").(string))
	d.Set("remaining_executions", d.Get("remaining_executions").(int))
	d.Set("project_id", d.Get("project_id").(string))

	return resourceWorkflowCronTriggerV2Read(ctx, d, meta)
}

func resourceWorkflowCronTriggerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	crontrigger, err := crontriggers.Get(ctx, workflowClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_workflow_cron_trigger_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_workflow_cron_trigger_v2 %s: %#v", d.Id(), crontrigger)

	d.Set("name", crontrigger.Name)
	d.Set("region", GetRegion(d, config))
	d.Set("workflow_id", crontrigger.WorkflowID)
	d.Set("workflow_name", crontrigger.WorkflowName)
	d.Set("workflow_input", crontrigger.WorkflowInput)
	d.Set("workflow_params", crontrigger.WorkflowParams)
	d.Set("pattern", crontrigger.Pattern)
	d.Set("remaining_executions", crontrigger.RemainingExecutions)
	d.Set("project_id", crontrigger.ProjectID)

	if err := d.Set("created_at", crontrigger.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for openstack_workflow_cron_trigger_v2 %s: %s", crontrigger.ID, err)
	}

	return nil
}

func resourceWorkflowCronTriggerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	err = crontriggers.Delete(ctx, workflowClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_workflow_cron_trigger_v2"))
	}

	return nil
}
