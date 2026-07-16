# Administrative Change Alert

This preset creates a resource-group-scoped alert on administrative
operations at critical/error severity that succeeded -- a broad safety net
for "something significant just changed in this resource group". Narrow it
to a specific `operationName` (for example
`Microsoft.Compute/virtualMachines/delete`) when you want to alert on one
kind of change; leave it broad to catch anything at the chosen severity.

## When to Use

- Guarding a production resource group against unexpected control-plane
  changes
- Feeding a change-audit channel that records who did what, when

## Key Configuration Choices

- **`category: ADMINISTRATIVE`** -- control-plane operations (create,
  update, delete, RBAC changes)
- **`levels` + `statuses`** -- narrow to the severity and outcome you care
  about; add `operationName` to target a single operation
- **Resource-group scope** -- watches every resource in the group; scope at
  a specific resource id to watch just one

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the alert in | The resource group's `status.outputs.resource_group_name` |
| `<watched-resource-group>` | The resource group the alert watches | The watched group's Planton resource name |
| `<action-group>` | The AzureMonitorActionGroup to notify | Your action group's Planton resource name |
