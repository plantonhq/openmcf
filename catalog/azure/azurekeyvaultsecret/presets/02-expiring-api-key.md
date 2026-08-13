# Expiring API Key

This preset stores a third-party API key with explicit activation and expiry attributes, making the credential's validity window -- and the rotation it implies -- auditable infrastructure.

## When to Use

- Partner and vendor API keys issued with contractual validity periods
- Credentials under a rotation policy where near-expiry monitoring should fire before things break
- Compliance regimes that audit secrets for missing expiration dates (Azure Policy has a built-in for exactly this)

## Key Configuration Choices

- **Expiry is advisory, monitored -- not enforced**: an RBAC-permitted reader can still fetch an expired value; the attribute's job is to make rotation visible and auditable
- **Update the dates when the key rotates** -- a value change creates a new secret version, and the new version should carry the new window
- **A `rotation` tag** documents cadence where cost/governance tooling looks

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-key-vault>` | The AzureKeyVault the secret lives in | The vault component's name |
| `<your-search-service>` | The example value source (an AzureSearchService whose admin key is stored) -- swap the whole `value.valueFrom` block for your real source | The producing component's name |

Adjust the two RFC 3339 instants to the key's real validity window before deploying.
