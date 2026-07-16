---
title: "Apex-to-WWW HTTPS Redirect"
description: "A catch-all URL map whose only job is redirecting bare-apex (`example.com`) and unmatched traffic to `www.example.com` over HTTPS — the usual first step before attaching real backends behind a global..."
type: "preset"
rank: "03"
presetSlug: "03-apex-redirect"
componentSlug: "url-map"
componentTitle: "URL Map"
provider: "gcp"
icon: "package"
order: 3
---

# Apex-to-WWW HTTPS Redirect

A catch-all URL map whose only job is redirecting bare-apex (`example.com`) and unmatched traffic to `www.example.com` over HTTPS — the usual first step before attaching real backends behind a global external load balancer.

## When to Use

- You own both apex and www DNS records but want all browsers on `www`
- HTTP-to-HTTPS upgrade at the load balancer before TLS termination on a target HTTPS proxy
- A placeholder URL map while backend services are still being provisioned

## Remix Notes

- `httpsRedirect: true` upgrades the scheme; pair with a target HTTPS proxy and managed certificate on the VIP.
- Use `prefixRedirect` instead of `hostRedirect` when you need path rewriting (e.g. `/old/*` → `/new/*`).
- Once real backends exist, replace `defaultUrlRedirect` with `defaultService` or add `hostRules`/`pathMatchers` — redirects and backends are mutually exclusive at each routing level.
