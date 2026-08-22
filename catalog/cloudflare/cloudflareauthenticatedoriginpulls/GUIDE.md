# CloudflareAuthenticatedOriginPulls guide

The judgment this guide protects you from: this kind is only HALF of a security control, its destroy does not turn anything off, and the association's remove path depends on state you might not expect.

## The origin side is your job

Enabling AOP makes Cloudflare PRESENT a client certificate on origin pulls. Nothing is protected until the origin REQUIRES and VALIDATES that certificate (nginx `ssl_verify_client`, Apache `SSLVerifyClient`, an ALB mTLS listener). An enabled zone with a non-validating origin is security theater -- traffic that bypasses Cloudflare still reaches the origin unchallenged.

## Destroy turns nothing off

The zone-wide toggle has NO delete at Cloudflare: destroying this resource drops it from state and abandons the live value. If AOP should be OFF after teardown, set `zone_enabled: false` and apply BEFORE destroying. An abandoned enabled-toggle on an origin that later stops validating certificates fails silently in the wrong direction -- Cloudflare keeps presenting, nobody keeps checking.

An association's "delete" is a revert write: the provider sends `enabled: null` (with the certificate id it holds in state) to void the association. That works -- but it is a WRITE, and it needs the state to still hold the association's certificate id.

## Presence semantics: unset means "leave it alone"

`zone_enabled` unset does not mean false -- it means the module does not manage the toggle at all (associations can be managed independently). Explicitly `false` asserts OFF. The same presence logic drives the association rows' `enabled`: unset is sent as TRUE, because Cloudflare treats null as "void the association" and a declared row is meant to exist. Set `enabled: false` for present-but-inactive.

## One resource per hostname, by provider law

The provider hard-fails any association resource whose config list holds more than one hostname -- at apply, not at plan. The module fans your rows out one-per-hostname automatically; keep hostnames unique in the list (duplicate keys collapse in the fan-out).

## Pairs well with

- [CloudflareAuthenticatedOriginPullsCertificate](../cloudflareauthenticatedoriginpullscertificate/README.md) -- upload the client certificate a row pins; wire `certificate_id` via `value_from`.
- [CloudflareMtlsCertificate](../cloudflaremtlscertificate/README.md) -- the account-level CA validating per-hostname client certificates.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
