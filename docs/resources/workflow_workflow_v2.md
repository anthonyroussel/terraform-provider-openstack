---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_workflow_v2"
sidebar_current: "docs-openstack-resource-workflow-workflow-v2"
description: |-
  Manages a Mistral V2 Workflow resource within OpenStack.
---

# openstack\_workflow\_workflow\_v2

Manages a Mistral V2 Workflow resource within OpenStack.

~> **Note:** This resource is immutable. Any changes to the workflow
definition requires deleting and recreating the resource.

~> **Note:** Each workflow resource MUST define exactly **one workflow**
because Terraform can track only one ID per resource. To manage multiple
workflows, declare a separate resource for each workflow.

## Example Usage

```hcl
resource "openstack_workflow_workflow_v2" "workflow_1" {
	namespace = "my_namespace"
	scope     = "private"
	definition = <<EOF
    version: '2.0'

    hello_workflow:
      description: Simple echo example

      input:
        - message

      tags:
        - echo

      tasks:
        echo:
          action: std.echo
          input:
            output:
              my_message: <% $.message %>
	EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Workflow client.

* `namespace` - (Optional) The namespace of the workflow.

* `scope` - (Required) The scope of the workflow (e.g., `private` or `public`).

* `definition` - (Required) The workflow definition in Mistral v2 DSL.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the workflow.
* `region` - See Argument Reference above.
* `namespace` - See Argument Reference above.
* `scope` - See Argument Reference above.
* `definition` - See Argument Reference above.
* `name` - The name of the workflow, automatically assigned based on the `definition` content.
* `input` - A set of input parameters required for workflow execution, automatically assigned based on the `definition` content.
* `tags` - The tags associated with the workflow, automatically assigned based on the `definition` content.
* `created_at` - The date the workflow was created.

## Import

Workflows can be imported using the `id`, e.g.

```shell
terraform import openstack_workflow_workflow_v2.workflow_1 53c3c098-c0b7-4fd9-ae93-7b4341fec0e5
```
