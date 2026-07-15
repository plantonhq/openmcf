---
title: "Presets"
description: "Ready-to-deploy configuration presets for SES Email Identity"
type: "preset-list"
componentSlug: "ses-email-identity"
componentTitle: "SES Email Identity"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-domain-easy-dkim"
    rank: "01"
    title: "Domain with Easy DKIM"
    excerpt: "Verifies a sending domain with AWS-managed DKIM keys. The stack output `dkim_tokens` carries three CNAME names to publish via `AwsRoute53DnsRecord`."
  - slug: "02-domain-with-config-set"
    rank: "02"
    title: "Domain with Configuration Set"
    excerpt: "A verified domain identity that inherits a default configuration set's delivery, tracking, and event-publishing rules."
---

# SES Email Identity Presets

Ready-to-deploy configuration presets for SES Email Identity. Each preset is a complete manifest you can copy, customize, and deploy.
