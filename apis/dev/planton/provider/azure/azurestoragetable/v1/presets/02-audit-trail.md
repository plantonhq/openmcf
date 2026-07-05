# Audit Trail Table

This preset creates a table for append-heavy audit/event records --
the workload Table storage is economically unbeatable for: petabytes of
small entities, written once, point-read or range-scanned rarely, at
storage-account prices with zero provisioned throughput.

## When to Use

- Audit trails, change logs, compliance event records
- IoT telemetry and device history
- Any append-mostly dataset too large or too cold for a database

## Key Configuration Choices

- **Partition by time bucket or tenant** (an application-side choice):
  spreading writes across partitions avoids the single-partition
  throughput ceiling, and time-bucketed partitions make retention
  scans cheap
- **Writers get Storage Table Data Contributor on THIS table** -- the
  `table_id` output is the least-privilege grant scope; auditors get
  Reader

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<AuditTableName>` | Letter-start 3-63 alphanumerics | Your naming convention |
