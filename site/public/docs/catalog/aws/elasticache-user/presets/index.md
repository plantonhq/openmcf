---
title: "Presets"
description: "Ready-to-deploy configuration presets for ElastiCache User"
type: "preset-list"
componentSlug: "elasticache-user"
componentTitle: "ElastiCache User"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-password-auth"
    rank: "01"
    title: "Password-Authenticated Application User"
    excerpt: "This preset creates one Redis/Valkey RBAC identity that authenticates with a password in the AUTH command. Each application gets its own user with an access string scoping exactly which keys and..."
  - slug: "02-iam-auth"
    rank: "02"
    title: "IAM-Authenticated User"
    excerpt: "This preset creates a Redis/Valkey RBAC identity that authenticates with a short-lived IAM-signed token instead of a long-lived password. No secret material lives in the manifest — clients sign..."
  - slug: "03-disabled-default"
    rank: "03"
    title: "Disabled Default User"
    excerpt: "This preset creates the mandatory \"default\" user every ElastiCache user group must contain. With the access string switched \"off\", unauthenticated clients are rejected outright — the standard..."
---

# ElastiCache User Presets

Ready-to-deploy configuration presets for ElastiCache User. Each preset is a complete manifest you can copy, customize, and deploy.
