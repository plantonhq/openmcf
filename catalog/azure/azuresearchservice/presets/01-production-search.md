# Production Search Service

The production posture: a `standard` service with three replicas (the
99.9% read-write SLA floor), Entra RBAC enabled alongside API keys,
and a system identity for keyless indexer connections.

## When to Use

- Production search workloads with availability requirements
- Estates moving toward Entra (RBAC) data-plane auth without breaking
  key-based clients
- Indexers that should reach data sources via managed identity

## Key Configuration Choices

- `replicaCount: 3` -- Azure's SLA floor for read-write workloads
  (2 replicas give read-only SLA; 1 gives none). Resizes in place.
- `authenticationFailureMode: http401WithBearerChallenge` -- setting
  it is what turns on RBAC-alongside-keys; the mode names how
  unauthenticated Entra calls are answered.
- `identity.type: SYSTEM_ASSIGNED` -- grant this identity data access
  (e.g. Storage Blob Data Reader) so indexers need no connection
  strings.
- `sku: standard` upgrades in place only to standard2/standard3 --
  any other SKU change replaces the service (and the name is held by
  DNS until deletion completes).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group-name>` | The resource group to create the service in | Portal -> Resource groups |
| `acme-search-prod` | Your globally-unique service name | It becomes {name}.search.windows.net |

## Related Presets

- `02-dev-basic-search` -- the cheap dedicated tier for development.
- `03-semantic-rag-search` -- semantic ranking for RAG applications.
