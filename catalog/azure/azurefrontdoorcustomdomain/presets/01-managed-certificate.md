# Managed-Certificate Domain

This preset creates a custom domain with Azure's managed certificate --
the zero-maintenance TLS posture: Azure issues, hosts, and auto-rotates
a free DV certificate for the exact hostname.

## When to Use

- Any single (non-wildcard) hostname up to 64 characters where a
  domain-validated certificate is sufficient -- the right default for
  most sites
- When you want zero certificate operations: no vault, no secret, no
  renewal calendar

## Key Configuration Choices

- **`tls: {}`** -- an empty TLS block deploys MANAGED_CERTIFICATE (the
  documented default); a secret reference is forbidden with it (there is
  nothing to bring)
- **`dnsZoneId` set** -- with Azure DNS, Azure watches the zone for the
  validation record; the TXT record itself is an AzureDnsRecord you
  manage in the zone
- **`domainName` vs `hostName`** -- the ARM resource name (hyphens, no
  dots) vs the actual hostname; keep them visibly related

## Validation Workflow (after deploy)

1. Read the domain's `validation_token` output.
2. Publish it as a TXT record at `_dnsauth.www.example.com`.
3. Azure flips the domain to approved (minutes).
4. CNAME `www.example.com` to the endpoint's `host_name` output.
5. Attach the domain to a route via `custom_domain_ids`.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<dns-zone-resource-name>` | The AzureDnsZone's Planton resource name | Your DNS composition |
| `hostName` (example value) | Your real hostname | Your DNS zone |
