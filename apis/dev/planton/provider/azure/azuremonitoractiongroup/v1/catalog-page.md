# Azure Monitor Action Group

Deploys an Azure Monitor action group -- the notification and automation hub alerts fire into. Email, SMS, voice, mobile push, webhooks, Azure Functions, Logic Apps, Automation runbooks, Event Hubs, ITSM systems, and ARM-role fan-out, all behind one referenceable node.

## What Gets Created

When you deploy an AzureMonitorActionGroup resource, Planton provisions:

- **Action group** -- an `azurerm_monitor_action_group` (global -- notifications keep flowing during regional outages) carrying your receiver lists
- **Azure Tags** -- Planton-derived governance tags merged with your own (your values win)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** (can reference an AzureResourceGroup resource)

## Quick Start

Create a file `action-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorActionGroup
metadata:
  name: platform-oncall
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMonitorActionGroup.platform-oncall
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: observability-rg
      fieldPath: status.outputs.resource_group_name
  shortName: PltOnCall
  emailReceivers:
    - name: oncall-email
      emailAddress: oncall@example.com
      useCommonAlertSchema: true
```

Deploy it:

```bash
planton pulumi up --manifest action-group.yaml
```

## Common Configurations

- **Keyless webhooks**: add `aadAuth` (an Entra application object ID) to a webhook receiver instead of baking a secret into the URL
- **Role fan-out**: an `armRoleReceivers` entry with a built-in role GUID (Owner: `8e3af657-a8ff-443c-a75c-2fe8c4bcb635`) notifies whoever holds the role -- no address list to maintain
- **Maintenance silence**: `enabled: false` swallows every alert without touching alert rules
- **Machine payloads**: set `useCommonAlertSchema: true` on anything parsed by software

## Key Outputs

| Output | Use |
| --- | --- |
| `action_group_id` | What metric alerts and scheduled query alerts reference in their actions |
