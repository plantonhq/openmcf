# Wildcard Certificate

This preset creates a Google-managed certificate covering an apex domain
and its wildcard under one certificate. Wildcards require DNS
authorization — load-balancer authorization cannot validate them.

## When to Use

- Serving many subdomains (`app.example.com`, `api.example.com`, ...)
  behind one certificate
- Multi-tenant platforms issuing tenant subdomains dynamically

## Key Configuration Choices

- **Apex + wildcard in one certificate** — one renewal lifecycle for both.
- **A single DNS authorization covers both entries** — an authorization
  for `example.com` validates `example.com` AND `*.example.com`.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |
| `<dns-authorization-resource-name>` | Name of the GcpCertManagerDnsAuthorization resource | Your DNS authorization manifest |

The sample domains `example.com` / `*.example.com` are realistic
placeholders for the pattern-validated `domains` entries — replace them
with your domain.

## Related Presets

- **01-managed-dns-auth** — single-domain managed certificate
- **03-self-managed-pem** — bring-your-own PEM certificate
