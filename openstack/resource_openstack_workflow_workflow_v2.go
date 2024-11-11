// Copyright (c) HashiCorp, Inc.

package openstack

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/verify"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: can be updated through the API,
				// FIXME: only one workflow at a time is supported at this time.
				ValidateFunc: verify.ValidStringIsYAML,
				StateFunc: func(v interface{}) string {
					template, _ := verify.NormalizeYAMLString(v)
					return template
				},
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

func resourceWorkflowWorkflowV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	log.Printf("[DEBUG] openstack_workflow_workflow_v2 create options: %#v", createOpts)

	workflowList, err := workflows.Create(ctx, workflowClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create openstack_workflow_workflow_v2: %s", err)
	}

	workflow := workflowList[len(workflowList)-1]

	d.SetId(workflow.ID)
	d.Set("name", workflow.Name)
	d.Set("scope", d.Get("scope").(string))
	d.Set("namespace", d.Get("namespace").(string))
	d.Set("definition", definition)
	d.Set("input", workflow.Input)
	d.Set("tags", workflow.Tags)
	d.Set("scope", workflow.Scope)
	d.Set("project_id", workflow.ProjectID)
	d.Set("created_at", workflow.CreatedAt)

	return resourceWorkflowWorkflowV2Read(ctx, d, meta)
}

func resourceWorkflowWorkflowV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	workflow, err := workflows.Get(ctx, workflowClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_workflow_workflow_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_workflow_workflow_v2 %s: %#v", d.Id(), workflow)

	d.Set("region", GetRegion(d, config))
	d.Set("name", workflow.Name)
	d.Set("scope", workflow.Scope)
	d.Set("namespace", workflow.Namespace)
	d.Set("definition", workflow.Definition)
	d.Set("input", workflow.Input)
	d.Set("tags", workflow.Tags)
	d.Set("scope", workflow.Scope)
	d.Set("project_id", workflow.ProjectID)
	d.Set("created_at", workflow.CreatedAt)

	return nil
}

func resourceWorkflowWorkflowV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
