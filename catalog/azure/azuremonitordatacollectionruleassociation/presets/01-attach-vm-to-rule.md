# Attach VM to Rule

This preset attaches one Azure virtual machine to a data collection rule -- the standard onboarding step that puts a machine under a fleet's collection policy.

## When to Use

- Every time a machine joins a rule's fleet: one association per machine per rule
- Inside charts that deploy a VM -- wire the association to the chart's VM and the shared rule, and the machine is born monitored

## Key Configuration Choices

- **The target is an explicit reference** -- `targetResourceId` names the kind and field path because several kinds can be attached (VMs, scale sets, Arc servers); swap the kind and `*_id` output for other targets
- **The rule reference uses the default wiring** -- `dataCollectionRuleId` defaults to the `AzureMonitorDataCollectionRule` kind's `data_collection_rule_id` output, so only the resource name is needed
- **Name after the rule** -- `linux-baseline-assoc` reads correctly in the machine's association listing; generic names guarantee archaeology later

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-vm>` | The Planton name of your `AzureVirtualMachine` resource | Planton console (or replace `valueFrom` with `value:` and a literal VM ARM ID) |
| `<your-data-collection-rule>` | The Planton name of your `AzureMonitorDataCollectionRule` resource | Planton console (or replace `valueFrom` with `value:` and a literal rule ARM ID) |
