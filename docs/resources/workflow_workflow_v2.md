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

## Example Usage

```hcl
resource "openstack_workflow_workflow_v2" "workflow_1" {
  namespace = "my_namespace"
  scope     = "private"
  definition = <<EOF
    version: '2.0'

    create_vm:
      description: Simple workflow example

      input:
        - vm_name
        - image_ref
        - flavor_ref
      output:
        vm_id: "{{ _.vm_id }}"
        vm_status: <% $.vm_status %>

      tasks:
        create_server:
          action: nova.servers_create name=<% $.vm_name %> image=<% $.image_ref %> flavor=<% $.flavor_ref %>
          publish:
            vm_id: <% task().result.id %>
          on-success:
            - wait_for_instance

        wait_for_instance:
          action: nova.servers_find id={{ _.vm_id }} status='ACTIVE'
          retry:
            delay: 5
            count: 15
          publish:
            vm_status: "{{ task().result.status }}"
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a VPN service. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    service.

* `name` - (Optional) The name of the service. Changing this updates the name of
    the existing service.

* `tenant_id` - (Optional) The owner of the service. Required if admin wants to
    create a service for another project. Changing this creates a new service.

* `description` - (Optional) The human-readable description for the service.
    Changing this updates the description of the existing service.

* `admin_state_up` - (Optional) The administrative state of the resource. Can either be up(true) or down(false).
    Changing this updates the administrative state of the existing service.

* `subnet_id` - (Optional) SubnetID is the ID of the subnet. Default is null.

* `router_id` - (Required) The ID of the router. Changing this creates a new service.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `status` - Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.
* `external_v6_ip` - The read-only external (public) IPv6 address that is used for the VPN service.
* `external_v4_ip` - The read-only external (public) IPv4 address that is used for the VPN service.
* `description` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

Workflows can be imported using the `id`, e.g.

```
$ terraform import openstack_workflow_workflow_v2.workflow_1 53c3c098-c0b7-4fd9-ae93-7b4341fec0e5
```
