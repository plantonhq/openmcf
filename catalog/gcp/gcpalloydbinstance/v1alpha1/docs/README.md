# GcpAlloydbInstance — Design & Research

## What this component is

`GcpAlloydbInstance` models one AlloyDB instance (`google_alloydb_instance`) on an existing cluster. The bundled primary lives in `GcpAlloydbCluster`; this kind exists so read pools and advanced instance types have their own lifecycle.

## Design notes

- **Cluster by full path** — `cluster` ref resolves `GcpAlloydbCluster.status.outputs.cluster_id`.
- **READ_POOL default** — presets and middleware default to READ_POOL; CEL enforces `read_pool_config.node_count` for that type.
- **Machine XOR** — `cpu_count` vs `machine_type`, same pattern as the cluster's primary instance block.
- **Public IP arms** — flattened from TF `network_config` (`enable_public_ip`, `enable_outbound_public_ip`, `authorized_external_networks`).
- **API enablement** — both engines enable `alloydb.googleapis.com` (cluster module currently does not).

## Deliberately unmodeled

- **`activation_policy`, `gce_zone`, labels/annotations beyond Planton attribution** — operational levers outside the 90/10 surface or handled by platform labels.
- **`deletion_policy` (ABANDON)** — conflicts with managed destroy semantics.

## Composition map

- `cluster` ← `GcpAlloydbCluster.status.outputs.cluster_id`
- READ_POOL instances depend on a PRIMARY existing on the same cluster (bundled or otherwise).
