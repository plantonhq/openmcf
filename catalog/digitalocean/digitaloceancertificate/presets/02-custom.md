# Custom Certificate

This preset creates an SSL certificate from user-provided PEM content. Use when you have a certificate from an enterprise CA, a purchased certificate, or a certificate issued outside of Let's Encrypt. The leaf certificate, private key, and optional intermediate chain are supplied directly.

## When to Use

- Enterprise or purchased SSL certificates
- Certificates from private CAs or internal PKI
- Wildcard or EV certificates not available via Let's Encrypt
- Migrating existing certificates to DigitalOcean

## Key Configuration Choices

- **Custom branch** (`custom`) -- setting this branch makes the certificate a custom one; you supply the PEM content, and there is no auto-renewal.
- **Leaf certificate** (`leafCertificate`) -- PEM-encoded server certificate; required.
- **Private key** (`privateKey`) -- PEM-encoded private key; must match the leaf certificate. It is a secret: DigitalOcean never returns it, and Planton state stores only a hash.
- **Certificate chain** (`certificateChain`) -- optional; include intermediate(s) if your CA provides them (e.g., DigiCert, Sectigo), ordered from issuing CA upward.
- **No auto-renewal** -- replace the certificate before expiry by re-applying with the new material; every field is create-only, and the replacement is created before the old certificate is destroyed, so load balancers referencing the name never observe a gap.

## Placeholders to Replace

The preset ships with a THROWAWAY self-signed pair (issued for `planton-e2e.invalid`) so it validates and dry-runs as-is -- it secures nothing and must be replaced.

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `leafCertificate` value | PEM content of the leaf/server certificate | Your CA or certificate issuer |
| `privateKey` value | PEM content of the private key | Generated with the cert or from your CA |
| `certificateChain` (add if needed) | PEM content of intermediate certificates | Your CA; often included in the cert package |
| `my-custom-cert` | Human-readable certificate identifier | Choose a name; used when referencing in load balancers |

## Related Presets

- **01-lets-encrypt** -- Use when you can use Let's Encrypt for free auto-renewal
