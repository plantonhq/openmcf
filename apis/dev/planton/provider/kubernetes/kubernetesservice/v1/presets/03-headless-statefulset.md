# Headless Service for StatefulSet Peers

This preset creates a headless service (`clusterIP: None`): no virtual IP is allocated and DNS returns the pod IPs directly, with each pod also getting its own stable name (`pod-0.db-peers.<namespace>.svc.cluster.local`). This is the governing Service a StatefulSet's `serviceName` points at, and the tool of choice whenever clients must address each pod individually.

## When to Use

- StatefulSet peer discovery — databases and quorum systems (PostgreSQL, Cassandra, Kafka, etcd) whose members must find each other by stable per-pod names
- Client-side load balancing where the client wants the full endpoint list, not one proxied VIP
- Any protocol where "connect to a specific replica" matters more than "connect to any replica"

## Key Configuration Choices

- **`headless: true`** — the spec-level way to get `clusterIP: None`; never set `cluster_ip_address` to "None". Headless is incompatible with `node_port`/`load_balancer` types and with a static cluster IP (validation rejects both combinations)
- **`publish_not_ready_addresses: true`** — DNS publishes pod addresses even before pods report Ready. Essential here: quorum members must discover each other DURING startup, before any of them can pass a readiness probe. Leave it false on ordinary traffic-serving services, where it would send real traffic to unready pods
- **`selector`** — matches the pods by their `app` label; for a Planton-managed StatefulSet this is the workload's `metadata.name`
- **`target_port`** — headless services do no proxying, so keep it equal to `port` (or omit it)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace the StatefulSet pods run in | Your namespace management |
| `<your-statefulset-name>` | Value of the `app` label on the StatefulSet's pods | The StatefulSet manifest, or `kubectl get pods --show-labels` |

The port `5432`/name `peer` pair is a working example (PostgreSQL) — replace with your system's peer port.

## Related Presets

- **01-cluster-ip-app** — a normal virtual-IP service for clients that just want "any replica"
