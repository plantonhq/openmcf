# Azure Monitor Data Collection Rule -- Operational Guide

Judgment calls that matter when you run DCR-based collection in production.

## Design rules around workloads, not machines

A rule is a reusable collection policy -- "Linux web tier: auth syslog + CPU/memory counters to the ops workspace" -- and machines attach to it by association. Resist the one-rule-per-machine anti-pattern: it multiplies objects without adding control (a machine can carry several associations when it genuinely needs several policies). A fleet-wide baseline rule plus small workload-specific rules composes better than one giant rule every team edits.

## Filter at the rule, not at the workspace

Everything a flow lands in Log Analytics bills per GB at ingestion. The rule gives you three progressively stronger filters that all run BEFORE billing: XPath queries on Windows event logs (take Critical/Error/Warning, skip Information), facility/severity selection on syslog (auth and daemon, not `*`), and `transform_kql` on the flow (drop columns, filter rows, reshape). A `"*"` facility list with an unfiltered flow is the classic surprise workspace bill.

## The 60-second rule for InsightsMetrics

Performance counters can flow to two different places with different contracts: `Microsoft-Perf` (into the workspace's Perf table, any 1-1800s sampling) and `Microsoft-InsightsMetrics` (into Azure Monitor metrics via the `azure_monitor_metrics` destination, sampling EXACTLY 60s -- Azure rejects anything else at deploy time). When a rule serves both, use two performance-counter sources with their own sampling rather than compromising one.

## Custom logs: the schema is a contract

A stream declaration's columns are the custom table's schema. Include a `TimeGenerated` datetime column (or produce one in `transform_kql`) or the workspace stamps arrival time, which skews every time-based query. Changing a declared schema later means coordinating the rule, any transformation, and queries against the table -- version custom stream names (`Custom-MyAppLogs2`) rather than mutating a live schema in place.

## Kind is a one-way door once set

A rule created without `kind` accepts every platform's sources and can adopt a kind later in place -- but once set, ANY change to it (including clearing) replaces the rule, and every association on the old rule dies with it. Leave `kind` unset unless you need what only a kind unlocks (the `*_direct` destinations require `AgentDirectToStore`).

## When nothing lands, walk the chain in order

Four links must hold: the machine carries an association to this rule; the Azure Monitor Agent runs on the machine (the association creates fine without it -- collection simply doesn't start); the data source's `streams` and the flow's `streams` match token-for-token; the flow's `destinations` name an existing destination. The portal's DCR **Data sources** blade and the agent's local troubleshooter (`/var/opt/microsoft/azuremonitoragent/log/mdsd.err` on Linux) localize which link broke. Records typically lag 3-5 minutes on a healthy chain -- do not debug earlier.

## Regional placement matters less than you think, until it does

Machines in any region can associate with a rule, but the rule performs its ingestion processing in ITS region -- co-locate the rule with its primary destination workspace to avoid cross-region processing latency, and expect Azure to reject cross-cloud bindings outright.
