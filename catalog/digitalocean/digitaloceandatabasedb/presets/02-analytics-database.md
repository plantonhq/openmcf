# Analytics Sidecar Database

This preset adds a second logical database to the same cluster for reporting or analytics tables, keeping them out of the transactional application database while sharing the cluster's compute.

## When to Use

- Separating analytics/reporting tables from transactional data without paying for a second cluster
- Giving a BI tool or ETL job its own blast radius (drop `analytics` without touching `orders`)
- Environments where schema experiments need a disposable home

## Key Configuration Choices

- **Same cluster, second database** -- logical databases are free; isolation comes from separate namespaces and separate users, not separate hardware.
- **Disposable by design** -- both fields are create-only, so deleting this resource drops exactly this database and nothing else on the cluster.

## What You Get

A second logical database on the shared cluster. Pair it with a read-only-ish DigitalOceanDatabaseUser for the analytics consumer.
