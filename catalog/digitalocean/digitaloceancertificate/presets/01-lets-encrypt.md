# Let's Encrypt Certificate

This preset creates a free, auto-renewing SSL certificate from Let's Encrypt via DigitalOcean. Supports multiple domains and wildcards. DigitalOcean handles renewal automatically; use the certificate name in load balancers for HTTPS termination. Ideal for production websites.

## When to Use

- Production HTTPS for public websites
- Cost-free SSL with automatic renewal
- Multiple domains or subdomains on one certificate
- Load balancer SSL termination

## Key Configuration Choices

- **Let's Encrypt branch** (`letsEncrypt`) -- setting this branch makes the certificate a Let's Encrypt one; DigitalOcean performs the ACME validation and always auto-renews.
- **Domains** (`domains`) -- list of FQDNs to include (e.g., `example.com`, `www.example.com`). Wildcards are supported (e.g., `*.example.com`). Every domain must already be managed by DigitalOcean DNS in the same account -- issuance fails otherwise.
- **Certificate name** (`certificateName`) -- the certificate's stable identity. Reference certificates BY NAME everywhere (load balancers do this natively): the certificate's UUID rotates on every auto-renewal, the name never does.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `example.com`, `www.example.com` | Domains to include in the certificate | Your DigitalOcean-DNS-managed domain and desired subdomains |
| `my-lets-encrypt-cert` | Human-readable certificate identifier | Choose a descriptive name; used in load balancer config |

## Related Presets

- **02-custom** -- Use when you have an existing certificate (e.g., from enterprise CA, purchased cert)
