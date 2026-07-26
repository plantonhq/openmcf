# Sharded analytics preset

The capacity posture: four shards, each a disjoint slice of the data
carried by two replicas — eight ClickHouse hosts a Distributed table
queries in parallel. This is the shape for datasets or write rates a
single shard cannot carry; if durability is the concern rather than
capacity, the production preset's single-shard/three-replica shape is
simpler and cheaper.

Sharding is a real commitment. Table design must route data (a
Distributed table over ReplicatedMergeTree locals, with a sharding
key you choose); rebalancing existing data after changing
`shards` is a manual migration, not a spec edit. The
`cluster_name` here ("warehouse") is what your DDL targets — 
`ON CLUSTER 'warehouse'` — and the operator caps it at 15 characters
because it becomes part of every generated child name.

Inter-host traffic for distributed queries authenticates with an
operator-generated shared secret, and replicas of the same shard
never share a Kubernetes node. The first thing to change is the
shard count — size it from your data, not from this example — then
`disk_size` per host and the memory pair, which is what ClickHouse
converts into speed.

See [03-sharded-analytics.yaml](./03-sharded-analytics.yaml) for the
manifest.
