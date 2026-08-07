# AzureMssqlFailoverGroup - Pulumi Module

Pulumi implementation for the AzureMssqlFailoverGroup deployment
component.

## Architecture

```
mssql.FailoverGroup (one cross-region failover group + listener outputs)
```

## Key Design Decisions

- **Listener FQDNs are composed, not read** --
  `{name}.database.windows.net` and
  `{name}.secondary.database.windows.net` are built from the group name
  because Azure does not return them; the group name is therefore also a
  DNS label.
- **`grace_minutes` pairs only with Automatic failover** (CEL-enforced
  at the spec); Manual sends no grace value.
- **Databases are set only when non-empty** -- an empty failover group
  is legal and deploys cleanly.
- **`readonly_endpoint_failover_policy_enabled` unset deploys the
  provider's Disabled default**, keeping both engines identical.
- **PARITY-EXCEPTION on tag shape** versus the Terraform module
  (documented in both engines) -- output-neutral.

## Provider

Built via the shared `pulumiazureprovider.Get` builder -- static client
secret, keyless web identity, or ambient chain, resolved from the stack
input. Never construct the provider inline.
