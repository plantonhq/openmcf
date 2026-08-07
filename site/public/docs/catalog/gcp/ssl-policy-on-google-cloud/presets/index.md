---
title: "Presets"
description: "Ready-to-deploy configuration presets for SSL Policy on Google Cloud"
type: "preset-list"
componentSlug: "ssl-policy-on-google-cloud"
componentTitle: "SSL Policy on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-modern-tls12"
    rank: "01"
    title: "Modern TLS 1.2 Baseline"
    excerpt: "The recommended production posture: a TLS 1.2 floor with the MODERN cipher profile. Attach it to every internet-facing HTTPS proxy unless a stricter regime demands RESTRICTED."
  - slug: "02-restricted-strict"
    rank: "02"
    title: "Restricted High-Security Policy"
    excerpt: "The strictest predefined posture: the RESTRICTED cipher profile with a TLS 1.2 floor. For frontends where security review outranks client reach."
  - slug: "03-custom-cipher-list"
    rank: "03"
    title: "Custom Cipher Allowlist"
    excerpt: "Hand-pick the exact cipher suites the load balancer may negotiate — for security reviews that specify an allowlist rather than a named profile."
---

# SSL Policy on Google Cloud Presets

Ready-to-deploy configuration presets for SSL Policy on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
