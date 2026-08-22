# CloudflareZeroTrustAccessApplication guide

Operational judgment for Access applications. The README covers what each field is; this covers how the pieces interact.

## The application is the door; policies are the guards

An Access application binds a protected resource (a hostname, an SSH target, a SaaS app) to the policies that decide who enters. The policies are standalone CloudflareZeroTrustAccessPolicy objects attached by reference in `policies` — evaluation follows precedence (lower first), then list order. An application with no policies is a locked door nobody can open; Cloudflare denies by default.

## Type decides which fields even exist

`self_hosted` (and ssh/vnc/rdp) fronts a hostname and requires `domain`; `saas` federates identity to an external app and ignores `domain` entirely; `bookmark` is just a launcher tile; `app_launcher` styles the portal itself; `infrastructure` and `mcp` target non-HTTP surfaces. Most of the spec's fields apply to only some types — the field comments say which, and setting a field the type ignores does nothing rather than failing.

## The domain must live in a zone this account controls

For self-hosted types, Cloudflare serves the Access login at the application's domain — which only works when that hostname's zone is in this account and its traffic actually flows through Cloudflare. Access on a hostname whose DNS points elsewhere protects nothing.

## Session duration is a trade, not a preference

`session_duration` caps how long a login lasts before re-authentication. Short durations (1h) suit admin panels; long ones (720h) suit low-risk internal tools where login friction costs more than it protects. The policy layer can override per policy; the application value is the ceiling users actually feel.

## SaaS applications mint real credentials

A `saas` application's outputs include the OIDC client secret Cloudflare generates — that is a live credential for your downstream app, marked sensitive in the outputs. Rotate it by re-keying the SaaS config, and treat every surface it lands on as secret storage.
