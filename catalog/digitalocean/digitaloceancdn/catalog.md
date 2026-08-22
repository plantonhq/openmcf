# CDN on DigitalOcean

Puts a Spaces bucket behind DigitalOcean's global edge network, so files load fast everywhere -- with an optional custom subdomain and managed TLS certificate. The endpoint itself is free: CDN delivery is included with the Spaces subscription. Integrates with Planton's Provider Connections for DigitalOcean API token management; the origin bucket and the certificate are wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **CDN Endpoint** -- the edge distribution fronting your referenced bucket, with your cache TTL and (optionally) your custom domain and certificate

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A bucket** -- the DigitalOceanBucket whose content the CDN serves.
- **A certificate (optional)** -- a DigitalOceanCertificate, required only with a custom domain.

### DigitalOcean Account

- Nothing extra: the endpoint is free and bandwidth draws from the bucket's existing transfer allowance.

## After You Deploy

Serve content from the `endpoint` output instead of the bucket's own domain. For a custom domain, point a CNAME at the endpoint -- and remember objects must be publicly readable in the bucket for the edge to serve them.
