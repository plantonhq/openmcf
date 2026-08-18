# Application Database

This preset creates one logical database for one application inside an existing managed cluster, referencing the cluster by name. Pair it with a DigitalOceanDatabaseUser of the same service and a DigitalOceanDatabaseConnectionPool pointing at this database name.

## When to Use

- Giving each application its own logical database on a shared cluster
- Replacing use of the cluster's built-in `defaultdb` for real workloads
- Chart composition: cluster + database + user + pool as one deployable unit

## Key Configuration Choices

- **Cluster by reference** (`valueFrom`) -- wires the database to a DigitalOceanDatabaseCluster in the same chart or environment; swap in a literal UUID for an existing cluster.
- **The name is the identity** -- `databaseName` is how clients, users, and pools address this database. Choose it like an API: renaming later REPLACES the database and drops its data.

## What You Get

A logical database visible in the cluster's Users & Databases tab, addressable as `orders` in connection strings and pool configurations.
