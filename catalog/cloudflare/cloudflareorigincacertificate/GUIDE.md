# CloudflareOriginCaCertificate guide

Operational judgment for Origin CA certificates. The README covers what each field is; this covers how the pieces interact.

## Revoke is not delete

Destroying this resource revokes the certificate — and the revoked certificate keeps answering GET 200 with its full body indefinitely, with `revoked_at` set (measured live). It never becomes a 404, so an automation that treats "the id exists" as "the cert is live" will lie forever: the honest liveness check is `revoked_at` being empty, not the id answering. Rotate the origin's installed cert *before* revoking the old one or you will break origin TLS.

## The CSR is write-only

The API never returns the CSR, and `requested_validity` is not returned after creation. An imported certificate lands without those fields; a post-import plan that wants to re-assert them is expected, not drift. Prefer the generated-key path (omit `spec.csr`) unless the private key must never touch Planton.

## Generated key vs BYO CSR

Omit `csr` and the module mints the key (RSA for `origin-rsa`, ECDSA for `origin-ecc`) and exports it as a sensitive `private_key` output — that is the one-click path, and you must store that output or you cannot install the cert on the origin. Supply a CSR and the key never leaves your control; `private_key` is empty and you already have the key.

## Hostnames are strings, not a zone reference

The Origin CA API takes hostname strings, not a zone id — but it validates every hostname against the account's ACTIVE zones. A hostname under a pending (undelegated) zone is rejected with a misleading 400 code 1010 "This zone is either not part of your account" even when the zone IS on the account (measured live); activate the zone first. A wildcard (`*.example.com`) does not include the apex; list both if the origin serves both.

Hostname order does not matter: Cloudflare stores the list in its own sorted order and the IaC modules send it sorted, so reordering a manifest's hostnames never re-issues the certificate. The first hostname you list becomes the generated CSR's common name — put the primary name first for readability.

## Adopting an existing certificate

A BYO-CSR certificate imports cleanly as the bare certificate id. A generated-key certificate does NOT round-trip: the key/CSR helpers ship no importer upstream, so the first apply after import mints a fresh key and CSR and re-issues the certificate. That replace is the honest adoption — plan for one origin cert rotation when adopting.
