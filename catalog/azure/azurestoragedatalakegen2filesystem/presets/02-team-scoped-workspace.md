# Team-Scoped Workspace Filesystem

This preset creates a self-service analytics workspace owned by an
Entra group -- the team writes freely inside its filesystem, and
membership changes happen in Entra, never in storage config.

## When to Use

- Per-team scratch/workspace areas in a shared lake account
- Data-science sandboxes where the team manages its own layout

## Key Configuration Choices

- **`group` hands root ownership to an Entra group** -- the owning
  group's rwx entry plus the DEFAULT inheritance entry means everything
  the team creates stays team-writable
- **`properties` values are base64-encoded** (Azure's requirement) --
  `YW5hbHl0aWNz` is "analytics"
- **NOTE for managed identities**: object IDs are PRINCIPAL ids, not
  client ids -- an ACL naming the client id silently never matches

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The HNS AzureStorageAccount's Planton resource name | Your lake composition |
| `<team-entra-group-object-id>` | The team's Entra group object ID | Entra ID -> Groups -> the team -> Object ID |
