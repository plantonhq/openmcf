---
title: "Presets"
description: "Ready-to-deploy configuration presets for MemoryDB User"
type: "preset-list"
componentSlug: "memorydb-user"
componentTitle: "MemoryDB User"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-password-auth"
    rank: "01"
    title: "Password-Authenticated Application User"
    excerpt: "This preset creates one MemoryDB ACL identity that authenticates with a password in the AUTH command. Each application gets its own user with an access string scoping exactly which keys and commands..."
  - slug: "02-iam-auth"
    rank: "02"
    title: "IAM-Authenticated Application User"
    excerpt: "This preset creates a MemoryDB user that authenticates with short-lived IAM-signed tokens instead of a password — no long-lived secret exists anywhere. The connecting workload signs an auth token..."
---

# MemoryDB User Presets

Ready-to-deploy configuration presets for MemoryDB User. Each preset is a complete manifest you can copy, customize, and deploy.
