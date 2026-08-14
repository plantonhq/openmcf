# CloudflareOriginCaCertificate guide

Operational judgment for Origin CA certificates. The README covers what each field is; this covers how the pieces interact.

## Revoke is not delete

Destroying this resource revokes the certificate. A just-revoked certificate may still answer GET for a window, so an automation that treats "the id exists" as "the cert is live" will lie. Wait for a real 404 before considering it gone, and rotate the origin's installed cert *before* revoking the old one or you will break origin TLS.

## The CSR is write-only

The API never returns the CSR, and `requested_validity` is not returned after creation. An imported certificate lands without those fields; a post-import plan that wants to re-assert them is expected, not drift. Prefer the generated-key path (omit `spec.csr`) unless the private key must never touch Planton.

## Generated key vs BYO CSR

Omit `csr` and the module mints the key (RSA for `origin-rsa`, ECDSA for `origin-ecc`) and exports it as a sensitive `private_key` output — that is the one-click path, and you must store that output or you cannot install the cert on the origin. Supply a CSR and the key never leaves your control; `private_key` is empty and you already have the key.

## Hostnames are strings, not a zone reference

The Origin CA API takes hostname strings, not a zone id. They should belong to a zone on this account, but this kind will not look that up. A wildcard (`*.example.com`) does not include the apex; list both if the origin serves both.
