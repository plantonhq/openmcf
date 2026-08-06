# Single Instance

This preset declares the smallest useful MongoDB cluster: one replica
set (`rs0`) with a single member, small storage, no backups. It is a
single point of failure by construction — the development posture,
cheap and fast to create. Applications connect at the exported
`kube_endpoint` (`dev-mongo-rs0.percona-mongo.svc.cluster.local:27017`
with `?replicaSet=rs0`), authenticating with the admin credential from
the operator-managed `dev-mongo-secrets` Secret (the outputs point at
it).

## When to Use

- Development and preview environments, feature branches, throwaway
  databases
- Anywhere losing the data is acceptable — there is no second member
  and no backup

## Key Configuration Choices

- **One replica set, `size: 1`** — a single mongod, still running as a
  replica set (`rs0`) so drivers and the oplog behave exactly as they
  will in production. No majority exists, so there is no failover
- **`unsafe.replset_size: true`** — the operator REJECTS replica sets
  smaller than 3 members without this explicit opt-in; the flag is the
  spec making the unsafe topology visible instead of silent
- **`storage.size: 5Gi`** — deliberately small; the operator applies
  grows to live PVCs and rejects shrinks, so starting small is safe
- **Small resources with an explicit memory limit** — WiredTiger sizes
  its cache from the memory limit, so even a dev instance should
  declare one
- **No backup block** — omitted backups are a deliberate choice the
  spec makes visible, not a silent default
- **Namespace = the operator's namespace** — the default operator
  posture watches its OWN namespace; a database declared elsewhere is
  never reconciled

## Placeholders to Replace

None — this preset deploys as-is once KubernetesPerconaMongoOperator
runs in the `percona-mongo` namespace (rename `metadata.name` per
database; every derived object follows it).

## Related Presets

- **02-replica-set** — the production shape: 3 members, resources,
  PodDisruptionBudget, zone anti-affinity
