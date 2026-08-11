# Priority Failover

This preset deploys an active/passive front tuned for fast failover: a Priority profile (traffic holds on the lowest-priority healthy endpoint) with fast-interval probing and a short TTL, cutting worst-case failover to well under a minute.

## When to Use

- Active/passive multi-region setups where the standby serves only when the primary fails
- Disaster-recovery fronts where failover speed matters more than probe cost

## Key Configuration Choices

- **Priority routing** -- give the primary endpoint priority 1, standbys 2, 3, ... (on the endpoint resources); traffic moves only on health
- **The fast-failover trio** -- the 10-second probe interval (billed extra per endpoint), one tolerated failure, and a 30-second TTL put detection + cache drain near ~50 seconds worst-case
- **The timeout is explicit** -- the fast interval narrows it to 5-9 (spec-enforced, mirroring Azure); the default 10 no longer fits

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-globally-unique-dns-label>` | The trafficmanager.net label -- globally unique across ALL of Azure | Your organization's naming convention (e.g. `contoso-app-dr`) |
