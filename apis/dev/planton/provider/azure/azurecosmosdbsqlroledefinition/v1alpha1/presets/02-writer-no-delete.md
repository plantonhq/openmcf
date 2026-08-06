# Writer Without Delete

This preset creates the role the built-ins cannot express: full read
access plus create/replace/upsert -- but never delete. Ingest
pipelines, event processors, and application workloads write documents
all day and have no business destroying them; reserving deletion for
operators (or a retention job with its own grant) turns a whole class
of bugs and compromises into authorization errors.

## When to Use

- Application workloads and ingest pipelines that append and update but
  must never destroy data
- Environments where deletion is an operator-only or retention-job-only
  action with its own audited grant
- Anywhere the built-in Data Contributor (which includes delete) is
  more power than the workload needs

## Key Configuration Choices

- **Deletion is absent, not carved out** -- Cosmos data-plane RBAC is
  allow-only, so this role lists every item operation EXCEPT
  `items/delete`; the wildcard `items/*` would silently grant deletion
- **The full read surface rides along** -- writers read what they
  write; splitting read and write grants doubles assignments for no
  security gain here
- **Account-wide assignable scope** -- narrow the list to database
  paths to pre-authorize where the role may ever be granted

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name (a GLOBAL_DOCUMENT_DB account) | Your Cosmos composition |

## Downstream Wiring

```yaml
# On an AzureCosmosdbSqlRoleAssignment
roleDefinitionId:
  valueFrom:
    kind: AzureCosmosdbSqlRoleDefinition
    name: my-writer-no-delete
    fieldPath: status.outputs.role_definition_id
```
