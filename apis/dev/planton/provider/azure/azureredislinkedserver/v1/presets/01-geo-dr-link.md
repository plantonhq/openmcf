# Geo-Replication DR Link

This preset links two Premium caches into the standard disaster-recovery
pair: the primary serves reads and writes while continuously replicating
to a warm secondary in another region.

## When to Use

- Regional-outage resilience for caches whose data is expensive or slow
  to rebuild
- Read locality: the secondary serves reads in its own region while
  linked
- Any workload with a documented RTO that a cold cache rebuild would
  blow through

## Key Configuration Choices

- **Every field references the caches** -- ids AND the location derive
  from the composed AzureRedisCache resources, so nothing is
  hand-repeated
- **Deleting the link IS the failover** -- unlinking makes the secondary
  writable; point applications at the
  `geo_replicated_primary_host_name` output (not either cache's own
  hostname) so failovers need no connection-string change
- **Azure's requirements** -- both caches PREMIUM, different regions,
  secondary at least as large as the primary; the secondary rejects
  writes while linked

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-cache-resource-name>` | The primary AzureRedisCache's Planton resource name | Your cache composition |
| `<secondary-cache-resource-name>` | The secondary AzureRedisCache's Planton resource name | Your DR-region composition |
