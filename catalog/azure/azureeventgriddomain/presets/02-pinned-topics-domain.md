# Pinned-Topics Domain

This preset creates the governance-posture domain: topics exist only as explicitly declared AzureEventgridDomainTopic resources, and publishing is Entra-only.

## When to Use

- SaaS products with a real tenant registry -- tenant onboarding should be an auditable IaC act, not a side effect of someone creating a subscription
- Domains carrying data-isolation or billing boundaries, where "which streams exist and why" must be answerable from code review

## Key Configuration Choices

- **Both lifecycle flags off** -- a subscription against an undeclared topic fails loudly, and topics never vanish because subscriptions churned; declare each stream with the AzureEventgridDomainTopic kind referencing this domain's `domain_id` output
- **`localAuthEnabled: false`** -- the domain key (which authorizes publishing to EVERY tenant's topic) stops working; the publishing service tier authenticates with Entra ID instead
- **Onboard topic and subscriptions together** -- events published to a topic with no subscription are dropped, not queued

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-governed-events` | The domain's region-wide-unique name -- swap in your org prefix | Your naming convention |
