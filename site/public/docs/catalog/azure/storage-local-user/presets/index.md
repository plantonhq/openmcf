---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Local User"
type: "preset-list"
componentSlug: "storage-local-user"
componentTitle: "Storage Local User"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-partner-key-auth"
    rank: "01"
    title: "Partner SFTP User with Key Authentication"
    excerpt: "This preset onboards one exchange partner with SSH public-key authentication and a full-access grant on the partner's own container -- the per-partner isolation pattern that lets one account serve..."
  - slug: "02-password-dropbox"
    rank: "02"
    title: "Password-Authenticated Drop Box"
    excerpt: "This preset creates a write-only \"drop box\" identity: a sender uploads files over SFTP with a generated password but can never read, list, or delete anything -- the receive-only posture for inbound..."
---

# Storage Local User Presets

Ready-to-deploy configuration presets for Storage Local User. Each preset is a complete manifest you can copy, customize, and deploy.
