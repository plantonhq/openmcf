# HTTP-Validated Certificate

A managed certificate for an apex domain (example.com itself): DNS forbids CNAME at the apex, so the domain routes by A record and Azure proves ownership by serving an HTTP token through it.

## When to Use

- Serving a naked domain from a Container App (apex domains route to the environment's static IP by A record)

## Key Configuration Choices

- The domain must already route to the app over HTTP for the validation to succeed: an A record to the environment's `static_ip_address` output, plus the `asuid.{domain}` TXT record
- Everything else matches the CNAME preset: bind the domain first, certificate attaches asynchronously, one hostname per certificate

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `example.com` | Replace with the apex domain | Your domain registrar account |
| `example-com` | Replace with the certificate's resource name (derive from the domain) | -- |
| `<your-environment>` | The AzureContainerAppEnvironment resource's metadata name | Your Planton resource inventory |

## Related Presets

- `01-cname-validated` -- the standard subdomain flow
