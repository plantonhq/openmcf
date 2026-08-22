# HTTPS Load Balancer with SSL Termination

This preset creates a load balancer that terminates TLS on port 443 and forwards traffic to backend Droplets over HTTP on port 80. An HTTP rule on port 80 plus `redirectHttpToHttps` sends browsers to HTTPS. Tag-based targeting attaches every Droplet carrying the `web` tag. Health checks keep unhealthy backends out of rotation.

## When to Use

- Production web applications that need HTTPS
- TLS termination at the load balancer (certificate managed as a `DigitalOceanCertificate`)
- Tag-based scaling: add or remove Droplets with the tag to change membership

## Key Configuration Choices

- **HTTPS to HTTP** (`entryPort: 443`, `entryProtocol: https`, `targetPort: 80`) -- TLS terminates at the balancer; backends serve plain HTTP.
- **Certificate by name** -- the `certificateName` reference resolves to `DigitalOceanCertificate.status.outputs.certificate_id`, which at the pinned provider is the certificate NAME (UUIDs rotate on Let's Encrypt renewal).
- **HTTP redirect** -- the port-80 rule plus `redirectHttpToHttps` sends cleartext visitors to 443.
- **Tag-based targeting** (`dropletTag: web`) -- membership follows the tag; do not also set `dropletIds`.
- **VPC** -- a reference to a `DigitalOceanVpc`. Omit it to use the region's default VPC.

## Related Presets

- **02-http-basic** -- Use when HTTPS is not required (dev/staging)
