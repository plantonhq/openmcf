# HTTPS Redirect Frontend

The production-standard role for a target HTTP proxy: port 80 exists only to bounce clients to HTTPS. The proxy points at a redirect-only URL map (its `defaultUrlRedirect` sets `httpsRedirect: true` with a 301), while the real application is served by the target HTTPS proxy sibling on port 443.

## When to Use

- Any internet-facing HTTPS load balancer — browsers still try port 80 first
- Completing the pair: one `GcpGlobalForwardingRule` on port 80 pointing here, one on port 443 pointing at the HTTPS proxy, both sharing a single static `GcpGlobalAddress`

## Remix Notes

- Reference the redirect `GcpUrlMap` via `valueFrom` instead of a literal self-link.
- The redirect URL map needs no backends at all — `defaultUrlRedirect` with `httpsRedirect: true` and `stripQuery: false` is the whole configuration.
- Keep the proxy's name aligned with its HTTPS sibling (e.g. `web-http` / `web-https`) so an operator can spot the pair at a glance.
