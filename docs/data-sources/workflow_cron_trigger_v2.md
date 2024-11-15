---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_cron_trigger_v2"
sidebar_current: "docs-openstack-datasource-workflow-cron-trigger-v2"
description: |-
  Get information on a cron trigger.
---

# openstack\_workflow\_cron_trigger_v2

Use this data source to get the ID of an available cron trigger.

## Example Usage

```hcl
data "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name        = "cron_trigger_1"
  workflow_id = "e6100d91-1e85-454a-9f5e-00a4cec71abb"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Workflow client.

* `name` - (Optional) The name of the cron trigger.

* `workflow_id` - (Optional) The ID of the associated workflow.

* `project_id` - (Optional) The id of the project to retrieve the workflow.
    Requires admin privileges.

## Attributes Reference

`id` is set to the ID of the found cron trigger.
In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `workflow_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `pattern` - The cron pattern.
* `workflow_params` - The workflow type specific parameters.
* `workflow_input` - The workflow input values.
* `remaining_executions` - The remaining executions.
* `first_execution_time` - The first execution time of the trigger.
* `created_at` - The date the workflow was created.
