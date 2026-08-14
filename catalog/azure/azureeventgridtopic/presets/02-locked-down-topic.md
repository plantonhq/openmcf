# Locked-Down Topic

This preset creates the hardened publish edge: Entra-only authentication (keys disabled), an IP allowlist, and a managed identity for secured delivery targets.

## When to Use

- Production topics carrying business events, once every publisher authenticates with Microsoft Entra ID
- Environments with known publisher egress (service tiers, CI runners) that an allowlist should pin

## Key Configuration Choices

- **`localAuthEnabled: false`** -- SAS keys stop working entirely; publishers need the EventGrid Data Sender role BEFORE this deploys, or publishing breaks at cutover
- **IP allowlist with public access on** -- the rules gate the public path; disabling public access instead would ignore the rules and require private endpoints
- **System-assigned identity** -- lets subscriptions deliver to identity-secured Event Hubs, queues, and dead-letter storage; grant the identity on those targets before creating the subscriptions

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-secure-events` | The topic's region-wide-unique name -- swap in your org prefix | Your naming convention |
| `<your-egress-cidr>` | The CIDR range your publishers egress from, e.g. `203.0.113.0/24` | Your network team / NAT gateway configuration |
