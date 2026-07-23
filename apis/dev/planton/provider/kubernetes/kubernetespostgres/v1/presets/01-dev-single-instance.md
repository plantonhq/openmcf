# Dev Single Instance

This preset declares the smallest useful PostgreSQL cluster: one
instance, small storage, a fresh `app` database with an
operator-generated password, no backups. It is a single point of
failure by construction — the development posture, cheap and fast to
create. Applications connect at the exported `kube_endpoint`
(`dev-db-rw.dev-db.svc.cluster.local:5432`) with the credentials from
the `dev-db-app` Secret.

## When to Use

- Development and preview environments, feature branches, throwaway
  databases
- Anywhere losing the data is acceptable — there is no replica and no
  backup

## Key Configuration Choices

- **`instances: 1`** — one primary, no replicas: no failover. The first
  production step is 3 instances (see 02-production-ha)
- **`storage.size: 5Gi`** — deliberately small; sizes can only GROW
  (the operator applies grows to live PVCs and rejects shrinks), so
  starting small is safe
- **initdb bootstrap with no declared password** — the operator
  generates the owner credential into the `dev-db-app` Secret; the
  outputs point applications at it (the Secret also carries ready-made
  `uri` / `jdbc-uri` connection strings)
- **No backup block** — omitted backups are a deliberate choice the spec
  makes visible, not a silent default
- **`enable_pdb: false`** — node drains never block on a development
  database (production keeps the default true, which protects the
  primary from voluntary eviction)

## Placeholders to Replace

None — this preset deploys as-is (rename `metadata.name` per database;
every derived object follows it).

## Related Presets

- **02-production-ha** — 3 instances, synchronous replication, hard
  anti-affinity, keyless S3 backups
- **03-s3-compatible-backups** — backups to MinIO/R2/any S3-compatible
  endpoint with declared keys
