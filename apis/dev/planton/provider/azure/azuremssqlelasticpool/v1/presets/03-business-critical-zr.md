# Zone-Redundant Business Critical Pool

This preset creates a 4-vCore Business Critical elastic pool with
zone-redundant replicas: local-SSD storage for the lowest I/O latency,
built-in synchronous replicas for the fastest failover, spread across
availability zones to survive a zone outage.

## When to Use

- Latency-sensitive OLTP fleets where I/O time dominates query time
- Databases whose availability target justifies the Business Critical
  premium

## Key Configuration Choices

- **`BC_Gen5`** -- Business Critical runs compute and storage together
  on local SSD with an Always On availability group under the hood
- **`zoneRedundant: true`** -- the replicas land in different zones;
  failover keeps working through a zone loss
- **`minCapacity: 0.5`** -- each database keeps half a vCore reserved,
  trading oversubscription for tail-latency protection

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `<server-region>` | The parent server's region (must match) | The server's spec |
| `<pool-name>` | The pool's name on the server | Your convention |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
