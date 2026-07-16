---
title: "IAP-Protected Internal Tool"
description: "Zero-trust access to an internal tool without a VPN: the backend service sits behind the public global load balancer, but Identity-Aware Proxy authenticates every request against Google identities..."
type: "preset"
rank: "03"
presetSlug: "03-iap-protected-internal-tool"
componentSlug: "backend-service"
componentTitle: "Backend Service"
provider: "gcp"
icon: "package"
order: 3
---

# IAP-Protected Internal Tool

Zero-trust access to an internal tool without a VPN: the backend service sits behind the public global load balancer, but Identity-Aware Proxy authenticates every request against Google identities before it reaches the backends.

## When to Use

- Admin dashboards, internal consoles, and staging environments that a distributed team reaches over the public internet
- Replacing VPN-gated access with identity-gated access

## Remix Notes

- Who gets in is IAM, not this spec: grant `roles/iap.httpsResourceAccessor` on the IAP-protected resource to the users/groups that may pass. Deploying IAP with no grants locks everyone out (fail-closed).
- The Google-managed OAuth client is right unless you need a branded consent screen — then set `oauth2ClientId`/`oauth2ClientSecret` together (the secret is reference-only secret material).
- IAP needs the frontend to terminate HTTPS (the target proxy's concern); the LB→backend leg here can stay HTTP.
- Backends should trust only IAP's signed headers (`x-goog-iap-jwt-assertion`) for identity — never a plain `X-Forwarded-For` check.
