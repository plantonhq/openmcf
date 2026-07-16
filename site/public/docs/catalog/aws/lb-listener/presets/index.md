---
title: "Presets"
description: "Ready-to-deploy configuration presets for LB Listener"
type: "preset-list"
componentSlug: "lb-listener"
componentTitle: "LB Listener"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-https-forward"
    rank: "01"
    title: "HTTPS Forward"
    excerpt: "This preset creates the workhorse ALB listener: HTTPS on 443, terminating TLS with an ACM certificate and forwarding everything to one target group. It is the entry point most architectures start..."
  - slug: "02-http-redirect-to-https"
    rank: "02"
    title: "HTTP Redirect to HTTPS"
    excerpt: "This preset creates the canonical port-80 listener: every plain-HTTP request gets a permanent redirect to HTTPS on 443, preserving the host, path, and query string. No target group is involved -- the..."
  - slug: "03-oidc-protected"
    rank: "03"
    title: "OIDC-Protected HTTPS"
    excerpt: "This preset creates an HTTPS listener that requires a login before any request reaches a target: an `authenticate-oidc` action redirects unauthenticated browsers to your identity provider (Okta,..."
---

# LB Listener Presets

Ready-to-deploy configuration presets for LB Listener. Each preset is a complete manifest you can copy, customize, and deploy.
