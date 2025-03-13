---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_workflow_v2"
sidebar_current: "docs-openstack-resource-workflow-workflow-v2"
description: |-
  Manages a V2 Workflow resource within OpenStack.
---

# openstack\_workflow\_workflow\_v2

Manages a V2 Workflow resource within OpenStack.

~> **Note:** This resource cannot be updated at this time. Any changes to the
workflow definition would require the deletion and recreation of the resource.

## Example Usage

```hcl
resource "openstack_workflow_workflow_v2" "workflow_1" {
  namespace = "some-namespace"
  scope     = "private"
  definition = <<EOF
    version: '2.0'

    workflow_echo:
      description: Simple workflow example
      type: direct
      tags:
        - echo

      input:
        - msg

      tasks:
        test:
          action: std.echo output="<% $.msg %>"
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Workflow client.

* `namespace` - The namespace of the workflow.

* `scope` - The scope of the workflow (e.g., `private` or `public`).

* `definition` - The workflow definition in Mistral v2 DSL.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Workflow.
* `region` - See Argument Reference above.
* `namespace` - See Argument Reference above.
* `scope` - See Argument Reference above.
* `definition` - See Argument Reference above.
* `name` - The name of the workflow, automatically assigned based on the `definition` content.
* `input` - A set of input parameters required for workflow execution, automatically assigned based on the `definition` content
* `tags` - The tags associated with the workflow, automatically assigned based on the `definition` content.
* `created_at` - The date the workflow was created.

## Import

Workflows can be imported using the `id`, e.g.

```bash
$ terraform import openstack_workflow_workflow_v2.workflow_1 53c3c098-c0b7-4fd9-ae93-7b4341fec0e5
```
