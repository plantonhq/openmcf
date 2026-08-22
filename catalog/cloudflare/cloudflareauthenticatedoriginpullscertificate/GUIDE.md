# CloudflareAuthenticatedOriginPullsCertificate guide

The judgment this guide protects you from: a provider defect makes one rotation path silently do nothing, the scope decides the blast radius, and deletion is slower than it looks.

## Never rotate only the private key

At provider v5.23.0 the zone-scoped upload has an empty Update and its `private_key` does not force replacement -- a plan that changes ONLY the key applies "successfully" while the API keeps the old key. The certificate presented to your origin and the key Cloudflare signs with silently diverge. The discipline: key and certificate always change together (a real re-issue), which forces the replacement the provider does honor. The hostname-scoped surface replaces on key changes correctly, but follow the same discipline everywhere -- a rotation habit that depends on remembering which scope is safe is not a habit.

## Scope is blast radius

`scope: zone` REPLACES Cloudflare's shared client certificate for the entire zone -- every origin pull presents it from then on. `scope: hostname` only uploads material; nothing changes until a `CloudflareAuthenticatedOriginPulls` association pins a hostname to the `certificate_id`. When in doubt, upload hostname-scoped and pin explicitly.

## Deletion is asynchronous

The API answers 200 with `pending_deletion` (and later `deleted`) before the record goes. Automation that deletes and immediately re-uploads the same certificate can race the pending state. Associations referencing a hostname certificate stop authenticating when it dies -- revert or re-point them BEFORE destroying the upload.

## Self-signed, byte-stable

The origin validates this certificate, so self-signed pairs are the designed case -- mint them with openssl and keep the PEM byte-stable (trailing newline included): the provider replaces on semantic certificate changes, and formatting churn is a replacement for nothing.

## Pairs well with

- [CloudflareAuthenticatedOriginPulls](../cloudflareauthenticatedoriginpulls/README.md) -- pins hostnames to this upload via `value_from` on `certificate_id`.
- [CloudflareMtlsCertificate](../cloudflaremtlscertificate/README.md) -- the account-level CA side of per-hostname validation.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
