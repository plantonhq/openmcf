# Overview

The **AzureMonitorDataCollectionRule** component deploys an Azure Monitor data collection rule (DCR) -- the routing table of Azure Monitor. A rule declares what telemetry the Azure Monitor Agent collects (Linux syslog, Windows event logs, performance counters, Prometheus metrics, IIS logs, custom log files, extension telemetry, Event Hub imports), where it lands (Log Analytics workspaces, Azure Monitor metrics, Event Hubs, storage), and how the two wire together (data flows, optionally with an ingestion-time KQL transformation).

## Purpose

- **One rule, any number of machines**: the rule is the reusable collection policy; machines attach to it with `AzureMonitorDataCollectionRuleAssociation` resources and detach without touching the rule.
- **Collect precisely, pay precisely**: XPath filters on Windows events, facility/severity filters on syslog, and ingestion-time KQL transformations drop noise BEFORE it bills at the workspace.
- **Custom logs as first-class telemetry**: stream declarations give a text or JSON log file a schema, and the rule lands it in a Log Analytics custom table.

## Key Features

- Full azurerm v5 surface: all ten data-source arms (syslog, performance counters, Windows event logs, extensions, IIS logs, custom log files, Prometheus forwarders, Windows Firewall logs, platform telemetry, Event Hub import) and all eight destination arms (Log Analytics, Azure Monitor metrics, Event Hub, Event Hub direct, Azure Monitor workspace, storage blob, storage blob direct, storage table direct).
- The provider's contracts front-loaded as validation: at least one destination, at least one data flow, exactly the provider's vocabularies (syslog facilities and levels, log-file timestamp formats, stream column types), unique custom-stream names, identity flavor rules.
- Chart-ready: the resource group, workspace, Event Hub, and storage-account references are `valueFrom`-wirable outputs of their own kinds; the rule's own `data_collection_rule_id` output is what associations and workspaces reference.

## Use Cases

- **Linux fleet to a workspace**: syslog (auth, daemon) plus CPU/memory performance counters into Log Analytics -- the canonical Azure Monitor Agent shape.
- **Windows fleet**: filtered event logs (errors and warnings only, via XPath) plus performance counters.
- **Custom application logs**: a JSON or text log file with a declared schema, transformed and landed in a custom workspace table.

## Future Enhancements

- The `data_collection_endpoint_id` and `monitor_account_id` references become typed references when the Data Collection Endpoint and Azure Monitor Workspace kinds enter the catalog (they are literal ARM ids today).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
