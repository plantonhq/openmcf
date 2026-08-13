# AzureMonitorDataCollectionRule Terraform Module

## Overview

Creates an Azure Monitor data collection rule (DCR) -- the routing table declaring what telemetry the Azure Monitor Agent collects (Linux syslog, Windows event logs, performance counters, Prometheus metrics, IIS logs, custom log files, extension telemetry, Event Hub imports), where it lands (Log Analytics workspaces, Azure Monitor metrics, Event Hubs, storage blobs/tables), and how the two wire together (data flows, optionally with an ingestion-time KQL transformation). Machines attach to the rule with separate `AzureMonitorDataCollectionRuleAssociation` resources.

## Resources Created

- `azurerm_monitor_data_collection_rule` -- the rule with its data sources, destinations, data flows, stream declarations, and identity

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMonitorDataCollectionRuleSpec fields; the resource group, workspace, Event Hub, storage-account, and DCE references arrive as resolved literals

## Outputs

- `data_collection_rule_id` -- the rule's ARM resource ID
- `data_collection_rule_name` -- the rule's resource name
- `immutable_id` -- the identifier agents and the ingestion API address the rule by
- `identity_principal_id` -- the system-assigned identity's principal (empty when no identity is configured)

## Behavior Notes

- **Names wire the rule together** -- flows reference destinations by their rule-local names, and destination names share ONE namespace across all arms; Azure enforces both at deploy time.
- **`kind` omitted means all platforms**; once set, changing (or clearing) it forces a new rule (provider lifecycle). Platform compatibility (Linux forbids `windows_event_log`, Windows forbids `syslog`, the `*_direct` destinations require `AgentDirectToStore`) is enforced by Azure at deploy time -- the provider performs no early check.
- **Custom streams need a DCE**: stream declarations (and the log_file / Prometheus sources that ride them) require `data_collection_endpoint_id`; Azure rejects the rule without one. Custom stream names must start with `Custom-`.
- **The Event Hub data import is a single block by design** -- Azure's rule model carries exactly one, and the provider silently uses only the first entry of its list.
- **`description`, `extension_json`, `consumer_group`, and the three flow transform fields are sent only when set** -- the provider validates non-empty strings.
- **`sampling_frequency_in_seconds` must be exactly 60** for streams targeting `Microsoft-InsightsMetrics` (Azure enforces at deploy time).
- **Name, resource group, and region are ForceNew**; everything else updates in place (except `kind`, per above).
- **Billing**: the rule object itself is free; you pay for the telemetry it lands (workspace ingestion, storage, Event Hub throughput).

## Required Permissions

The deploying principal needs `Microsoft.Insights/dataCollectionRules/*` (Monitoring Contributor covers it) plus read on the referenced destinations.
