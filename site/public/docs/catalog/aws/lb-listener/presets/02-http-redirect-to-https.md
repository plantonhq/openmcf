---
title: "HTTP Redirect to HTTPS"
description: "This preset creates the canonical port-80 listener: every plain-HTTP request gets a permanent redirect to HTTPS on 443, preserving the host, path, and query string. No target group is involved -- the..."
type: "preset"
rank: "02"
presetSlug: "02-http-redirect-to-https"
componentSlug: "lb-listener"
componentTitle: "LB Listener"
provider: "aws"
icon: "package"
order: 2
---

# HTTP Redirect to HTTPS

This preset creates the canonical port-80 listener: every plain-HTTP request
gets a permanent redirect to HTTPS on 443, preserving the host, path, and
query string. No target group is involved -- the ALB answers the redirect
itself. Deploy it alongside an HTTPS listener (see **01-https-forward**) on
the same ALB.

## When to Use

- Every internet-facing ALB that serves browsers -- users type `http://`,
  and search engines follow it
- Enforcing HTTPS-only access without touching application code

## Key Configuration Choices

- **`HTTP_301`** -- a permanent redirect that browsers cache, so repeat
  visitors skip the extra round trip; use `HTTP_302` only while testing
  redirect behavior
- **Only protocol and port overridden** -- the untouched components
  (`#{host}`, `#{path}`, `#{query}`) pass through, so
  `http://example.com/a?b=c` lands on `https://example.com/a?b=c`
- **No certificate** -- TLS material lives on the 443 listener; this one
  never sees encrypted traffic

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<alb-resource-name>` | Name of the AwsAlb resource to attach to | Your AwsAlb manifest's `metadata.name` |

## Common Additions

- Add a `Strict-Transport-Security` response header on the *HTTPS* listener
  (`httpHeaders.response.strictTransportSecurity`) so browsers stop sending
  plain HTTP at all after the first visit

## Related Presets

- **01-https-forward** -- the 443 listener this one redirects to
- **03-oidc-protected** -- an HTTPS listener with a login step in front
