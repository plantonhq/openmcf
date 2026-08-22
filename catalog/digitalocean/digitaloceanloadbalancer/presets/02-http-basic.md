# HTTP Load Balancer (explicit Droplets)

This preset creates a regional HTTP load balancer forwarding port 80 to port 8080 on two named Droplets. No TLS. Suitable for development, staging, or an internal service sitting behind a CDN.

## When to Use

- HTTP-only environments
- Explicit Droplet membership (as opposed to tag-based)
- A first balancer while a certificate is still being issued

## Key Configuration Choices

- **HTTP to HTTP** (`entryPort: 80`, `targetPort: 8080`) -- no certificate required.
- **Explicit Droplets** -- `dropletIds` references two `DigitalOceanDroplet` resources. Mutually exclusive with `dropletTag`.
- **VPC** -- a reference to a `DigitalOceanVpc`. The Droplets should live in the same VPC.

## Related Presets

- **01-https-ssl-termination** -- Use for production HTTPS with tag-based targeting
