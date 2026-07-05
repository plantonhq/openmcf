---
title: "Host-Based Routing"
description: "This preset routes one hostname to one service: requests whose Host header matches (say `api.example.com`) forward to the service's target group. Each service gets its own subdomain while sharing the..."
type: "preset"
rank: "02"
presetSlug: "02-host-based-routing"
componentSlug: "lb-listener-rule"
componentTitle: "LB Listener Rule"
provider: "aws"
icon: "package"
order: 2
---

# Host-Based Routing

This preset routes one hostname to one service: requests whose Host header
matches (say `api.example.com`) forward to the service's target group. Each
service gets its own subdomain while sharing the ALB, the listener, and --
via SNI -- per-domain certificates. Host rules are naturally disjoint, so
they rarely fight over priority the way path rules do.

## When to Use

- One ALB serving many subdomains (`api.`, `app.`, `admin.`), each backed by
  its own service
- Multi-tenant platforms where tenants get their own hostnames
  (`*.customers.example.com` as a wildcard value)
- Keeping services isolated by domain without one load balancer per service

## Key Configuration Choices

- **Host values support wildcards** -- `*.example.com` matches any single
  level (`*` spans characters, `?` matches one); several `values` entries OR
  together
- **DNS and certificates travel together** -- the hostname needs a DNS record
  pointing at the ALB and a certificate covering it on the listener
  (`certificateArn` or `additionalCertificateArns`)
- **`priority: 20`** -- explicit for predictability, though disjoint host
  rules can safely omit it and let AWS append

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name prefix for the rule resource | Your service's name |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<listener-resource-name>` | Name of the AwsLbListener to attach to | Your AwsLbListener manifest's `metadata.name` |
| `<service-hostname>` | Hostname to match (e.g., `api.example.com`) | Your DNS zone |
| `<target-group-resource-name>` | Name of the AwsLbTargetGroup receiving traffic | Your AwsLbTargetGroup manifest's `metadata.name` |

## Common Additions

- Add a `pathPattern` condition block to route only part of the domain
  (blocks AND together)
- Add a `host-header-rewrite` transform when the backend expects a different
  hostname than the public one

## Related Presets

- **01-path-based-routing** -- split by URL prefix instead of domain
- **03-canary-weighted** -- shift this route's traffic gradually between two groups
