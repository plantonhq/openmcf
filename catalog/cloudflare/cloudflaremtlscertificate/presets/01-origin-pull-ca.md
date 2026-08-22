# Origin-pull CA upload

A CA certificate uploaded to the account mTLS store -- the trust anchor that per-hostname Authenticated Origin Pulls and zone TLS CA associations reference. Self-signed CAs are the normal case: your infrastructure validates these, not the public.

## When to Use

- Setting up per-hostname Authenticated Origin Pulls (the CA validates the client certificates)
- Scoping a CA to specific hostnames through `CloudflareZoneTlsSettings` CA associations
- Any consumer that needs YOUR trust material at the account level

## Key Configuration Choices

- **ca: true** -- trust material for validating clients; no private key travels (the signing key stays in your PKI). Use `ca: false` with a `private_key` only when Cloudflare must PRESENT the certificate itself.
- **No private_key** -- deliberately absent on CA uploads; uploading one anyway spreads a secret for nothing.
- **Rotation is replacement** -- every field is create-only; upload the new CA, re-point consumers, then destroy this one.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account | Cloudflare Dashboard -> account home -> right-hand panel (Account ID) |
| `certificates` | The CA certificate (or chain) in PEM form | Your PKI -- `openssl req -x509 -new ...` mints a self-signed CA |

## Related Presets

None yet -- a leaf variant (ca false, with key) belongs with a Workers mTLS binding that consumes it.
