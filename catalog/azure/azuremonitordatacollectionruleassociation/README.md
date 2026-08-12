# Overview

The **AzureMonitorDataCollectionRuleAssociation** component attaches ONE machine (an Azure virtual machine, a VM scale set, or an Arc-enabled server) to an Azure Monitor data collection rule or data collection endpoint. The association is how a machine enters monitoring: the Azure Monitor Agent on the target discovers its associations, downloads the referenced rule, and starts collecting.

## Purpose

- **Onboard machines without touching the rule**: the rule is the shared collection policy; each machine attaches with its own association and detaches independently -- fleets grow and shrink without a single edit to the rule.
- **Compose monitoring in charts**: a chart deploying a VM wires `target_resource_id` to the VM's `vm_id` output and `data_collection_rule_id` to the rule's output -- the machine is born monitored.
- **Private-link configuration access**: the endpoint form binds a machine to a Data Collection Endpoint so agent configuration flows through private networking.

## Key Features

- Full azurerm v5 surface: both binding arms (rule XOR endpoint, provider-enforced and spec-validated), the target as a no-default reference (VM, VMSS, or Arc server), description.
- The endpoint arm's fixed-name contract handled automatically: Azure mandates the name `configurationAccessEndpoint` for endpoint associations, and both engines apply it when the name is left unset.
- Chart-ready: `data_collection_rule_id` defaults its reference to the AzureMonitorDataCollectionRule kind's `data_collection_rule_id` output.

## Use Cases

- **Put a VM under a fleet rule**: one association per VM against the shared Linux or Windows baseline rule.
- **Layer workload monitoring**: a machine carries the fleet baseline association plus a workload-specific rule's association -- several rules, one machine.
- **Private configuration access**: bind machines in locked-down networks to a Data Collection Endpoint.

## Future Enhancements

- The `data_collection_endpoint_id` reference becomes a typed reference when the Data Collection Endpoint kind enters the catalog (it is a literal ARM id today).
