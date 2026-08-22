# CloudflareCustomSslCertificate guide

The judgment this guide protects you from: this kind is the manual-renewal path on a paid plan, rotation changes the certificate's identity, and the trust wall is Cloudflare's, not yours.

## Business plan or above, publicly trusted only

Cloudflare enforces both walls at create, not here: a zone below Business is rejected (402/403), and a certificate not issued by a CA on Cloudflare's trust list -- including every self-signed certificate -- is rejected regardless of plan. If you control the origin and just need Cloudflare-to-origin authentication, you want `CloudflareAuthenticatedOriginPullsCertificate` (self-signed is normal there), not this kind.

## You own the renewal now

Universal SSL and `CloudflareCertificatePack` renew themselves. A custom certificate expires on YOUR calendar -- watch the `expires_on` output and rotate before it. Rotation is replacement: the upload is destroyed and re-created, the `certificate_id` changes, and anything referencing the old ID must follow. Cloudflare keeps serving the previous certificate until the replacement deploys, so a timely rotation is not an outage.

## legacy_custom vs sni_custom

`legacy_custom` (the API default) works on every TLS client but occupies the zone's single legacy slot -- a second legacy upload is rejected. `sni_custom` requires SNI (every modern browser) and allows multiple uploads. Changing `type` on an existing upload replaces it. Prefer `sni_custom` unless you have measured legacy-client traffic.

## Priority is not manageable

At provider v5.23.0 the certificate `priority` is read-only -- the v4 reprioritization surface was deliberately dropped. If you need a specific certificate served preferentially, control it by what you upload, not by an ordering knob.

## The policy asymmetry

The API accepts the geo policy as `policy` and returns it as a separate read-only field (`policy_restrictions`). The modules send exactly what you write; a non-empty re-plan on this field after apply is the provider's own normalization class, not your manifest.

## Destroy leaves the zone on its fallback

Destroy is a real delete and deployment states settle asynchronously (status may read pending briefly). The zone falls back to Universal SSL / Advanced certificates -- verify one is active BEFORE destroying, or visitors get TLS handshake failures.

## Pairs well with

- [CloudflareCertificatePack](../cloudflarecertificatepack/README.md) -- the managed-renewal alternative; read the boundary above before you pick.
- [CloudflareZoneTlsSettings](../cloudflarezonetlssettings/README.md) -- minimum TLS version and Universal SSL posture on the same zone.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
