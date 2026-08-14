# CNAME Alias

This preset gives a service a stable name over a moving target: a CNAME answering `service.<your-zone>` with another hostname, so consumers keep one name while the infrastructure behind it changes.

## When to Use

- A stable name over a failover pair's active member, a private endpoint's generated FQDN, or any hostname that changes with redeployments
- Decoupling application configuration from infrastructure names -- apps hold the alias, operators move the target

## Key Configuration Choices

- **The target can be a reference** -- point `cname` at another component's hostname output with `valueFrom` when the target is minted at deploy time; a literal for externally-known names
- **Low TTL for moving targets** -- 60 keeps a target change visible within a minute; the point of an alias is usually agility
- **No CNAME at the zone apex** -- DNS itself forbids it (spec names must be a child label); private DNS has no alias-record workaround, so apex names use A records with referenced addresses

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-azure-private-dns-zone-resource-name>` | The AzurePrivateDnsZone component's resource name | Your Planton catalog |
| `<target-hostname>` | The hostname this alias answers with | The target service's documentation or its component's outputs |
