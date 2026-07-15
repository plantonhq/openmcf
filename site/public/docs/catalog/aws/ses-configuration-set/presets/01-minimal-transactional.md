---
title: "Minimal Transactional"
description: "A production-minded baseline: require TLS for delivery, honor account-level suppression for bounces and complaints, and enable reputation metrics."
type: "preset"
rank: "01"
presetSlug: "01-minimal-transactional"
componentSlug: "ses-configuration-set"
componentTitle: "SES Configuration Set"
provider: "aws"
icon: "package"
order: 1
---

# Minimal Transactional

A production-minded baseline: require TLS for delivery, honor account-level
suppression for bounces and complaints, and enable reputation metrics.

## When to Use

- Transactional mail (password resets, receipts, alerts) where TLS and
  suppression hygiene matter from day one

## What It Configures

- **`tlsPolicy: REQUIRE`** — deliver only over TLS-protected connections
- **`suppressedReasons: [BOUNCE, COMPLAINT]`** — skip sends to suppressed addresses
- **`reputationMetricsEnabled: true`** — publish bounce/complaint rates to CloudWatch

## What to Customize

- Replace `<aws-region>` with your sending region
- Attach event destinations in a follow-on manifest when you need bounce feedback loops
