# AzureMonitorDataCollectionRuleAssociation Pulumi Module

## Overview

Attaches ONE machine (an Azure VM, a VM scale set, or an Arc-enabled server) to an Azure Monitor data collection rule or data collection endpoint. The association is how a machine enters monitoring: the Azure Monitor Agent on the target discovers its associations, downloads the referenced rule, and starts collecting. Removing the association detaches the machine without touching the rule.

## Resources Created

- `monitoring.DataCollectionRuleAssociation` -- the attachment, an extension resource scoped under the target machine's ARM ID

## Outputs

- `data_collection_rule_association_id` -- the association's ARM resource ID (scoped under the target)
- `data_collection_rule_association_name` -- the association's name on the target

## Behavior Notes

- **Exactly one binding** -- a rule OR an endpoint, never both (spec CEL mirrors the provider's ExactlyOneOf).
- **`name` is required for rule bindings and left unset for endpoint bindings** -- Azure mandates the fixed name `configurationAccessEndpoint` for endpoint associations, and the provider applies it as the default.
- **The target and name are ForceNew**; swapping the bound rule or endpoint updates in place.
- **The association creates fine on a machine without the Azure Monitor Agent** -- collection simply starts when the agent arrives; the association is pure configuration.
- **No tags** -- ARM extension resources are untagged; the provider carries no tags argument.
- **Billing**: free. The telemetry the rule collects is billed at its destinations.
