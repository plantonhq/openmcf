# Policy-Anchored Read-Only Table

This preset creates a table whose external access rides a read-only
stored access policy. SAS tokens issued against the policy inherit its
window and query-only permissions -- and revoking the policy instantly
revokes every token anchored to it.

## When to Use

- Sharing reference data or reports with external consumers without
  identity federation
- Time-boxed read access that must be revocable in one action

## Key Configuration Choices

- **`permissions: r`** -- query-only; the letters follow Azure's strict
  order (r, a, u, d) and anything out of order is rejected by the spec
- **Table policies require the full window** -- start AND expiry are
  mandatory (unlike blob/share policies); the expiry is the revocation
  backstop
- **At most five policies per table** -- Azure's limit, enforced in the
  spec

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `SharedEntities` | Letter-start 3-63 alphanumerics | Your naming convention |
| `<window-start-rfc3339>` | Policy validity start, e.g. `2026-07-01T00:00:00Z` | Your sharing agreement |
| `<window-expiry-rfc3339>` | Policy validity end -- the revocation backstop | Your sharing agreement |
