package openstack

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"

	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/verify"
)

func resourceWorkflowWorkbookV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowWorkbookV2Create,
		ReadContext:   resourceWorkflowWorkbookV2Read,
		DeleteContext: resourceWorkflowWorkbookV2Delete,
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

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceWorkflowWorkbookV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	definition, err := verify.NormalizeYAMLString(d.Get("definition").(string))
	if err != nil {
		return diag.Errorf("Unable to parse definition: %s", err)
	}

	createOpts := workbooks.CreateOpts{
		Scope:      d.Get("scope").(string),
		Namespace:  d.Get("namespace").(string),
		Definition: strings.NewReader(definition),
	}

	workbook, err := workbooks.Create(ctx, workflowClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create openstack_workflow_workbook_v2: %s", err)
	}

	d.SetId(workbook.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("scope", d.Get("scope").(string))
	d.Set("namespace", d.Get("namespace").(string))
	d.Set("definition", definition)
	d.Set("name", workbook.Name)
	d.Set("tags", workbook.Tags)
	d.Set("scope", workbook.Scope)
	d.Set("project_id", workbook.ProjectID)

	if err := d.Set("created_at", workbook.CreatedAt.Format(time.RFC3339)); err != nil {
		tflog.Debug(ctx, "Unable to set created_at for openstack_workflow_workbook_v2", map[string]interface{}{
			"id":  workbook.ID,
			"err": err,
		})
	}

	return resourceWorkflowWorkbookV2Read(ctx, d, meta)
}

func resourceWorkflowWorkbookV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	workbook, err := workbooks.Get(ctx, workflowClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_workflow_workbook_v2"))
	}

	tflog.Debug(ctx, "Retrieved openstack_workflow_workbook_v2", map[string]interface{}{
		"id":       d.Id(),
		"workbook": workbook,
	})

	d.Set("region", GetRegion(d, config))
	d.Set("name", workbook.Name)
	d.Set("scope", workbook.Scope)
	d.Set("namespace", workbook.Namespace)
	d.Set("definition", workbook.Definition)
	d.Set("tags", workbook.Tags)
	d.Set("scope", workbook.Scope)
	d.Set("project_id", workbook.ProjectID)

	if err := d.Set("created_at", workbook.CreatedAt.Format(time.RFC3339)); err != nil {
		tflog.Debug(ctx, "Unable to set created_at for openstack_workflow_workflow_v2", map[string]interface{}{
			"id":  workbook.ID,
			"err": err,
		})
	}

	return nil
}

func resourceWorkflowWorkbookV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	err = workbooks.Delete(ctx, workflowClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_workflow_workbook_v2"))
	}

	return nil
}
