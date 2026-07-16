---
title: "Domain with Configuration Set"
description: "A verified domain identity that inherits a default configuration set's delivery, tracking, and event-publishing rules."
type: "preset"
rank: "02"
presetSlug: "02-domain-with-config-set"
componentSlug: "ses-email-identity"
componentTitle: "SES Email Identity"
provider: "aws"
icon: "package"
order: 2
---

# Domain with Configuration Set

A verified domain identity that inherits a default configuration set's delivery,
tracking, and event-publishing rules.

## When to Use

- Production senders where TLS posture, suppression, and event destinations are
  defined once on the set and shared across identities

## What It Configures

- Domain identity with Easy DKIM (2048-bit)
- **`configurationSet`** reference to an `AwsSesConfigurationSet` output

## What to Customize

- Replace `example.com`, `<aws-region>`, and `<configuration-set-name>`
- Deploy the configuration set first, then this identity
