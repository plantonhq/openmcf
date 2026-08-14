---
title: "Presets"
description: "Ready-to-deploy configuration presets for Secrets Manager Secret"
type: "preset-list"
componentSlug: "secrets-manager-secret"
componentTitle: "Secrets Manager Secret"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-application-credentials"
    rank: "01"
    title: "Application Credentials"
    excerpt: "This preset creates a Secrets Manager secret holding an application's credentials as a JSON key/value document, encrypted under the AWS-managed key, with a 7-day recovery window."
  - slug: "02-rotated-database-password"
    rank: "02"
    title: "Rotated Database Password"
    excerpt: "This preset creates a KMS-encrypted secret with 30-day automatic rotation through a rotation Lambda — the production posture for database credentials."
  - slug: "03-multi-region-replicated"
    rank: "03"
    title: "Multi-Region Replicated"
    excerpt: "This preset creates a secret replicated to additional regions with a resource policy restricting reads to a named application role."
---

# Secrets Manager Secret Presets

Ready-to-deploy configuration presets for Secrets Manager Secret. Each preset is a complete manifest you can copy, customize, and deploy.
