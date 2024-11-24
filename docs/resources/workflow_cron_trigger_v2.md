---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_cron_trigger_v2"
sidebar_current: "docs-openstack-workflow-cron-trigger-v2"
description: |-
  Manages a V2 cron trigger resource within OpenStack.
---

# openstack\_workflow\_cron\_trigger\_v2

Manages a V2 cron trigger resource within OpenStack.

## Example Usage

```hcl
resource "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name        = "cron_trigger_1"
  workflow_id = "428ce188-9881-4784-b36d-ef20659feced"
  pattern     = "0 5 * * *"

  workflow_input = {
    param1 = "val1"
    param2 = "val2"
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

* `workflow_input` - (Optional) Map of input parameters to pass to the workflow
    upon execution. Changing this creates a new cron trigger.

* `workflow_params` - (Optional) Map of additional parameters for the workflow.
    Changing this creates a new cron trigger.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the cron trigger.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `workflow_id` - See Argument Reference above.
* `pattern` - See Argument Reference above.
* `workflow_input` - See Argument Reference above.
* `workflow_params` - See Argument Reference above.

## Import

Cron triggers can be imported using the `id`, e.g.

```
$ terraform import openstack_workflow_cron_trigger_v2.cron_trigger_1 bae24970-d96e-4ed0-80c1-b798cb2208c6
```
