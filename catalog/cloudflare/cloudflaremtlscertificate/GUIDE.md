# CloudflareMtlsCertificate guide

The judgment this guide protects you from: everything on this kind is create-only, the certificate ID is the identity consumers hold, and destroying an upload out of order strands them.

## Immutable by design: rotate by replace

Every field forces replacement at the API -- there is no in-place update. The rotation discipline is a three-step dance, in this order:

1. Upload the NEW certificate (a new resource, or change the value and let the apply replace).
2. Re-point every consumer (zone TLS CA associations, Workers mTLS bindings, AOP rows) at the new `certificate_id`.
3. Destroy the OLD upload.

Skipping step 2 before step 3 leaves consumers referencing a deleted certificate -- client validation starts failing at the edge, not at apply time.

## CA vs leaf: the ca flag decides the consumer

`ca: true` is trust material for VALIDATING clients (Authenticated Origin Pulls, API Shield). It needs no private key -- uploading one anyway spreads a secret for nothing. `ca: false` is a leaf Cloudflare PRESENTS as a client (Workers mTLS bindings); that one needs its key. The flag must be stated explicitly and cannot change after upload.

## Self-signed is the point

Unlike zone-facing certificates, nothing here needs public trust -- your origin or your policy validates these. Mint CAs with your own PKI (or openssl) and keep the signing key outside Cloudflare entirely; the store only ever needs the certificate side of a CA.

## The key never comes back

The API never returns `private_key`. An imported or refreshed resource re-asserts it from configuration, which means your secret store is the real system of record -- an upload whose key you lost cannot be re-presented, only replaced.

## Pairs well with

- [CloudflareZoneTlsSettings](../cloudflarezonetlssettings/README.md) -- the CA hostname associations that scope this CA to hostnames.
- [CloudflareAuthenticatedOriginPulls](../cloudflareauthenticatedoriginpulls/README.md) -- the zone-level enablement the CA validates for.
- [CloudflareWorker](../cloudflareworker/README.md) -- mTLS bindings that present leaf uploads from this store.
