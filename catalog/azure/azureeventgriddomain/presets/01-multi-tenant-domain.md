# Multi-Tenant Domain

This preset creates a CloudEvents domain on Azure's auto-managed topic lifecycle -- topics materialize when their first subscription arrives and vanish with the last one.

## When to Use

- Multi-tenant eventing where subscription creation IS tenant onboarding (no separate topic registry to maintain)
- Internal platforms starting frictionless -- the lifecycle flags update in place, so you can tighten to the pinned posture later

## Key Configuration Choices

- **Auto-managed lifecycle (the defaults, sent explicitly)** -- zero ceremony per tenant; the trade-off is that "which streams exist" is answered by Azure, not your IaC
- **CloudEvents input schema** -- domain-wide and create-only: every tenant's events arrive in the same envelope, chosen deliberately on day one
- **Defaults left on** -- public network access and key auth stay enabled for the fast integration path; harden with an IP allowlist and `localAuthEnabled: false` once publishers are on Entra ID

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-tenant-events` | The domain's region-wide-unique name -- swap in your org prefix | Your naming convention |
