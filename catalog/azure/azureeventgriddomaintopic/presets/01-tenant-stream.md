# Tenant Stream

This preset declares one tenant's event stream inside an Event Grid domain -- the onboarding move of the pinned-topics posture.

## When to Use

- Every tenant onboarding on a pinned-topics domain: one declared topic per tenant, created with the tenant's subscriptions
- Converting an auto-managed domain to governed streams -- declaring a topic pins it regardless of the domain's lifecycle flags

## Key Configuration Choices

- **The domain reference uses the default wiring** -- `domainId` defaults to the `AzureEventgridDomain` kind's `domain_id` output, so only the resource name is needed
- **The name is the publisher contract** -- publishers stamp exactly this value into the event's topic field; it is create-only, so treat a rename like an API version change
- **Pair with subscriptions in the same chart** -- events published to a topic with no subscription are dropped, not queued

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-eventgrid-domain>` | The Planton name of your `AzureEventgridDomain` resource | Planton console (or replace `valueFrom` with `value:` and a literal domain ARM ID) |
| `customer-fabrikam` | The stream's name -- swap in your tenant's stable slug | Your tenant registry |
