# CloudflareZeroTrustAccessApplication guide

Operational judgment for Access applications. The README covers what each field is; this covers how the pieces interact.

## The application is the door; policies are the guards

An Access application binds a protected resource (a hostname, an SSH target, a SaaS app) to the policies that decide who enters. The policies are standalone CloudflareZeroTrustAccessPolicy objects attached by reference in `policies` — evaluation follows precedence (lower first), then list order. An application with no policies is a locked door nobody can open; Cloudflare denies by default.

## Type decides which fields even exist

`self_hosted` (and ssh/vnc/rdp) fronts a hostname and requires `domain`; `saas` federates identity to an external app and ignores `domain` entirely; `bookmark` is just a launcher tile; `app_launcher` styles the portal itself; `infrastructure` and `mcp` target non-HTTP surfaces. Most of the spec's fields apply to only some types — the field comments say which, and setting a field the type ignores does nothing rather than failing.

## The domain must live in a zone this account controls — an ACTIVE one

For self-hosted types, Cloudflare serves the Access login at the application's domain — which only works when that hostname's zone is in this account and its traffic actually flows through Cloudflare. Access on a hostname whose DNS points elsewhere protects nothing.

The zone must also be ACTIVE (registrar-delegated to Cloudflare's nameservers), not merely added: the create is rejected with `400 access.api.error.invalid_request: domain does not belong to zone` when the domain sits on a PENDING zone, and the identical create succeeds the moment the zone is active (measured live, 2026-08-26). If a fresh zone was just added to the account, finish the nameserver delegation before creating self-hosted applications on it.

## Session duration is a trade, not a preference

`session_duration` caps how long a login lasts before re-authentication. Short durations (1h) suit admin panels; long ones (720h) suit low-risk internal tools where login friction costs more than it protects. The policy layer can override per policy; the application value is the ceiling users actually feel.

## SaaS applications mint real credentials

A `saas` application's outputs include the OIDC client secret Cloudflare generates — that is a live credential for your downstream app, marked sensitive in the outputs. Rotate it by re-keying the SaaS config, and treat every surface it lands on as secret storage.

## Adopting an existing application: the first apply touches three echoed toggles

Cloudflare's read API echoes `auto_redirect_to_identity`, `enable_binding_cookie`, and `options_preflight_bypass` as `false` on types that support them, while a manifest that never enabled them omits them entirely — so the first plan after adopting an existing self-hosted application shows an in-place update on exactly those attributes (measured live at provider v5.23.0). The apply converges to the same server state. Expected, harmless, once.
