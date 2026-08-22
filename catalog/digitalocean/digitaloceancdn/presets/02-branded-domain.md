# Branded Domain CDN

This preset serves bucket content from your own subdomain over TLS: the edge endpoint plus a managed certificate and a custom domain, so users see `assets.example.com` instead of a digitaloceanspaces.com host.

## When to Use

- Production asset delivery under a brand domain
- Any consumer-facing URL that must survive an origin or endpoint change (re-point the CNAME, keep the URL)

## Key Configuration Choices

- **Certificate by reference, by NAME** -- Let's Encrypt renewals rotate certificate UUIDs; the stable name is the only safe handle (and validation requires the certificate whenever a custom domain is set).
- **CNAME the domain at the endpoint output** -- DNS is yours to point; the endpoint hostname is in the outputs.
- **Moderate TTL (3600)** -- a starting point; raise it for fingerprinted content.

## What You Get

A branded, TLS-served edge distribution -- free on the endpoint side, with the certificate auto-renewing on DigitalOcean's side.
