---
title: "Rotated Database Password"
description: "This preset creates a KMS-encrypted secret with 30-day automatic rotation through a rotation Lambda — the production posture for database credentials."
type: "preset"
rank: "02"
presetSlug: "02-rotated-database-password"
componentSlug: "secrets-manager-secret"
componentTitle: "Secrets Manager Secret"
provider: "aws"
icon: "package"
order: 2
---

# Rotated Database Password

This preset creates a KMS-encrypted secret with 30-day automatic rotation through a rotation Lambda — the production posture for database credentials.

## When to Use

- Database master or application passwords in production
- Any credential your compliance posture requires rotating on a schedule
- Credentials shared with other AWS accounts (the customer-managed key is what makes cross-account reads possible)

## Key Configuration Choices

- **Customer-managed KMS key** — the AWS-managed key cannot be granted cross-account, and a dedicated key gives you rotation-independent audit and revocation
- **`automaticallyAfterDays: 30`** — AWS derives the rotation window; use `scheduleExpression` (`rate(...)`/`cron(...)`) when you need precise control
- **`rotateImmediately: true`** — the bootstrap value is replaced by a rotation-function-issued credential as soon as the secret deploys

## Prerequisites

The rotation Lambda must implement the four rotation steps and grant Secrets Manager invoke permission (`principal: secretsmanager.amazonaws.com`) — AWS's rotation function templates cover the common database engines. Both the KMS key and the Lambda are wired here as references to their Planton resources.
