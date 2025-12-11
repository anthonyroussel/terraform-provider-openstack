---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_cron_trigger_v2"
sidebar_current: "docs-openstack-workflow-cron-trigger-v2"
description: |-
  Manages a Mistral V2 Cron Trigger resource within OpenStack.
---

# openstack\_workflow\_cron\_trigger\_v2

Manages a Mistral V2 Cron Trigger resource within OpenStack.

A cron trigger schedules the execution of a workflow using a cron-like pattern.

## Example Usage

```hcl
data "openstack_workflow_workflow_v2" "hello_workflow" {
  name = "hello_workflow"
}

resource "openstack_workflow_cron_trigger_v2" "hello_cron_trigger" {
  name          = "hello_cron_trigger"
  workflow_id   = data.openstack_workflow_workflow_v2.hello_workflow.id
  pattern       = "0 5 * * *"

  workflow_input = {
    message = "Hello, OpenStack!"
  }

  workflow_params = {
    priority = "high"
    notify   = ["mistral@openstack.org"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Workflow client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new cron trigger.

* `name` - (Required) The name of the cron trigger. Changing this creates a new
    cron trigger.

* `workflow_id` - (Required) The ID of the workflow to be executed by this cron
    trigger. Changing this creates a new cron trigger.

* `pattern` - (Required) A cron-like schedule pattern indicating when the
    workflow should be executed. Changing this creates a new cron trigger.

* `workflow_input` - (Optional) Map of input parameters passed to the workflow
    upon execution. Changing this creates a new cron trigger.

* `workflow_params` - (Optional) Map of additional workflow parameters.
    Changing this creates a new cron trigger.

* `remaining_executions` - (Optional) Maximum number of executions for the cron trigger.
    If omitted, the trigger runs indefinitely.
    This value is applied only at creation time and automatically decreases as the
    trigger executes. To avoid perpetual diffs and state drift, the provider does
    not track updates to this value after creation.
    Changing this creates a new cron trigger.

* `first_execution_time` - (Optional) The date of the first execution.
    Must be in RFC3339 format (e.g. `2025-12-12T19:54:12Z`). If omitted, the trigger starts
    at the next scheduled cron interval.
    Changing this creates a new cron trigger.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the cron trigger.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Cron Trigger.
* `name` - See Argument Reference above.
* `workflow_id` - See Argument Reference above.
* `pattern` - See Argument Reference above.
* `workflow_input` - See Argument Reference above.
* `workflow_params` - See Argument Reference above.
* `remaining_executions` - See Argument Reference above.
* `first_execution_time` - See Argument Reference above.
* `created_at` - The time at which cron trigger was created.

## Import

Cron triggers can be imported using the `id`, e.g.

```shell
terraform import openstack_workflow_cron_trigger_v2.cron_trigger_1 bae24970-d96e-4ed0-80c1-b798cb2208c6
```
