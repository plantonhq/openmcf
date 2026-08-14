# A Record

This preset writes the everyday private name entry: an A record answering `db.<your-zone>` with a private IPv4 address, resolvable from every virtual network linked to the zone.

## When to Use

- Naming databases, internal APIs, and other services on private addresses
- Replacing hosts-file entries and hand-run `az network private-dns record-set` commands with declared resources
- Any name whose address you manage (auto-registered VM records are the service's own lifecycle -- do not declare those)

## Key Configuration Choices

- **The zone is a reference** -- `privateDnsZoneId` points at an `AzurePrivateDnsZone` component by name (deploy order follows automatically); pass `value: <arm-id>` for a zone managed outside Planton
- **All of the name's addresses in ONE record** -- multiple entries round-robin; Azure caps an A record set at 20 (spec-enforced)
- **TTL follows change cadence** -- 300 (the platform default) suits most services; drop to 60 for failover-sensitive names, raise to 3600+ for stable ones

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-azure-private-dns-zone-resource-name>` | The AzurePrivateDnsZone component's resource name | Your Planton catalog (or use `value:` with the zone's ARM ID from `az network private-dns zone show`) |
