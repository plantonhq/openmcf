---
title: "Presets"
description: "Ready-to-deploy configuration presets for SES Configuration Set"
type: "preset-list"
componentSlug: "ses-configuration-set"
componentTitle: "SES Configuration Set"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-minimal-transactional"
    rank: "01"
    title: "Minimal Transactional"
    excerpt: "A production-minded baseline: require TLS for delivery, honor account-level suppression for bounces and complaints, and enable reputation metrics."
  - slug: "02-cloudwatch-events"
    rank: "02"
    title: "CloudWatch Event Destination"
    excerpt: "Transactional posture plus a CloudWatch event destination for zero-infrastructure bounce and complaint metrics — alarm on `AWS/SES` without standing up SNS or Firehose."
---

# SES Configuration Set Presets

Ready-to-deploy configuration presets for SES Configuration Set. Each preset is a complete manifest you can copy, customize, and deploy.
