---
title: "Application Credentials"
description: "This preset creates a Secrets Manager secret holding an application's credentials as a JSON key/value document, encrypted under the AWS-managed key, with a 7-day recovery window."
type: "preset"
rank: "01"
presetSlug: "01-application-credentials"
componentSlug: "secrets-manager-secret"
componentTitle: "Secrets Manager Secret"
provider: "aws"
icon: "package"
order: 1
---

# Application Credentials

This preset creates a Secrets Manager secret holding an application's credentials as a JSON key/value document, encrypted under the AWS-managed key, with a 7-day recovery window.

## When to Use

- Database passwords, API keys, and tokens an application reads at runtime
- Any credential that today lives in an environment variable or config file
- The starting point for most secrets — add rotation or replication as needs grow

## Key Configuration Choices

- **Hierarchical name** (`<env>/<app>/credentials`) — Secrets Manager names support `/`, and path-style names keep IAM policies simple (`secretsmanager:GetSecretValue` on `arn:...:secret:prod/*`)
- **Managed-secret reference** (`$secret/...`) — the value is resolved just-in-time at deploy; plaintext never lives in the control plane
- **7-day recovery window** — recoverable after an accidental destroy without reserving the name for the default 30 days

## What You Get

The secret with one version staged `AWSCURRENT`. Applications read it with `GetSecretValue` by name or ARN (use the `secret_arn` output — AWS appends a random suffix, so the ARN is never derivable from the name).
