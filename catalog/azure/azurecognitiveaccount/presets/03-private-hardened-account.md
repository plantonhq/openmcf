# Private Hardened Account

This preset locks the account down on every axis the service offers: deny-by-default network ACLs with explicit IP and subnet allowances, Entra-ID-only authentication (keys disabled), and restricted outbound traffic with an FQDN allowlist.

## When to Use

- Production AI workloads under a security review
- Data-sensitive scenarios where the account must not call arbitrary endpoints
- Organizations standardizing on token (Entra ID) auth over shared keys

## Key Configuration Choices

- **Deny by default** -- only the listed IP ranges and subnets reach the endpoint; `AZURE_SERVICES` bypass keeps trusted Azure services working
- **Keys off** -- `localAuthEnabled: false` empties both key outputs; applications authenticate with Entra ID tokens
- **Outbound restricted** -- the account can only call FQDNs on the allowlist (grounding data-loss prevention)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the account | `AzureResourceGroup` status outputs (`resource_group_name`) |
| `acme-openai-locked` | Example account name (also the subdomain) -- replace with your own globally unique name | Your naming convention |
| `203.0.113.0/24` | Example documentation range -- replace with the public IPv4 range allowed through the perimeter | Your network team |
| `<your-app-subnet-id>` | ARM ID of the subnet your app runs in | `AzureSubnet` status outputs (`subnet_id`), or reference it with valueFrom |
| `<your-allowed-outbound-fqdn>` | Hostname the account may call outbound | The service you ground on (e.g. your search endpoint) |
