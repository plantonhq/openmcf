# CloudflareAuthenticatedOriginPullsCertificate guide

The judgment this guide protects you from: a provider defect makes one rotation path silently do nothing, the scope decides the blast radius, and deletion is slower than it looks.

## Never rotate only the private key

At provider v5.23.0 the zone-scoped upload has an empty Update and its `private_key` does not force replacement -- a plan that changes ONLY the key applies "successfully" while the API keeps the old key. The certificate presented to your origin and the key Cloudflare signs with silently diverge. The discipline: key and certificate always change together (a real re-issue), which forces the replacement the provider does honor. The hostname-scoped surface replaces on key changes correctly, but follow the same discipline everywhere -- a rotation habit that depends on remembering which scope is safe is not a habit.

## Scope is blast radius

`scope: zone` REPLACES Cloudflare's shared client certificate for the entire zone -- every origin pull presents it from then on. `scope: hostname` only uploads material; nothing changes until a `CloudflareAuthenticatedOriginPulls` association pins a hostname to the `certificate_id`. When in doubt, upload hostname-scoped and pin explicitly.

## Deletion is asynchronous

The API answers 200 with `pending_deletion` (and later `deleted`) before the record goes. Automation that deletes and immediately re-uploads the same certificate can race the pending state. Associations referencing a hostname certificate stop authenticating when it dies -- revert or re-point them BEFORE destroying the upload.

## Self-signed, byte-stable -- but always a LEAF

The origin validates this certificate, so self-signed pairs are the designed case -- but Cloudflare rejects CA-flagged uploads on the origin-pull surfaces (400 code 1412 "Missing leaf certificate.", measured live 2026-08-28), and a plain `openssl req -x509` mints CA:TRUE by default. Mint with an explicit leaf profile:

```
openssl req -x509 -newkey rsa:2048 -keyout client-key.pem -out client-cert.pem \
  -days 365 -nodes -subj "/CN=my-origin-pull-client" \
  -addext "basicConstraints=critical,CA:FALSE" \
  -addext "keyUsage=digitalSignature,keyEncipherment" \
  -addext "extendedKeyUsage=clientAuth"
```

Keep the PEM byte-stable (trailing newline included): the provider replaces on semantic certificate changes, and formatting churn is a replacement for nothing.

## Importing cannot restore the key

The API never returns `private_key`, and on the hostname surface the key forces replacement -- so the first apply after importing an existing upload REPLACES it under a NEW certificate id, and every association referencing the old id must be re-pointed. Prefer the create-new / re-point / delete-old path over import when adopting existing material.

## Pairs well with

- [CloudflareAuthenticatedOriginPulls](../cloudflareauthenticatedoriginpulls/README.md) -- pins hostnames to this upload via `value_from` on `certificate_id`.
- [CloudflareMtlsCertificate](../cloudflaremtlscertificate/README.md) -- the account-level CA side of per-hostname validation.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
