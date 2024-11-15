package openstack

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
)

func dataSourceWorkflowCronTriggerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowCronTriggerV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"workflow_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWorkflowCronTriggerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	listOpts := crontriggers.ListOpts{
		ProjectID: d.Get("project_id").(string),
	}

	name := d.Get("name").(string)
	if name != "" {
		listOpts.Name = &crontriggers.ListFilter{
			Filter: "eq",
			Value:  name,
		}
	}

	workflowID := d.Get("workflow_id").(string)
	if workflowID != "" {
		listOpts.WorkflowID = workflowID
	}

	allPages, err := crontriggers.List(workflowClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query cron triggers: %s", err)
	}

	allCronTriggers, err := crontriggers.ExtractCronTriggers(allPages)

	if err != nil {
		return diag.Errorf("Unable to retrieve cron triggers: %s", err)
	}

	if len(allCronTriggers) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var cronTrigger crontriggers.CronTrigger
	if len(allCronTriggers) > 1 {
		tflog.Debug(ctx, "Multiple results found", map[string]interface{}{"allCronTriggers": allCronTriggers})
		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}
	cronTrigger = allCronTriggers[0]

	dataSourceWorkflowCronTriggerV2Attributes(ctx, d, &cronTrigger, GetRegion(d, config))

	return nil
}

func dataSourceWorkflowCronTriggerV2Attributes(ctx context.Context, d *schema.ResourceData, cronTrigger *crontriggers.CronTrigger, region string) {
	d.SetId(cronTrigger.ID)
	d.Set("region", region)
	d.Set("name", cronTrigger.Name)
	d.Set("project_id", cronTrigger.ProjectID)

	if err := d.Set("created_at", cronTrigger.CreatedAt.Format(time.RFC3339)); err != nil {
		tflog.Debug(ctx, "Unable to set created_at for openstack_workflow_cron_trigger_v2", map[string]interface{}{
			"id":  cronTrigger.ID,
			"err": err,
		})
	}
}
