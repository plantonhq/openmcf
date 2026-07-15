---
title: "Application Data Access Point"
description: "Enforced POSIX identity (1000:1000) with a dedicated `/app-data` root directory, auto-created on first mount. The standard shape for giving one application least-privilege access to a shared file..."
type: "preset"
rank: "01"
presetSlug: "01-app-data"
componentSlug: "efs-access-point"
componentTitle: "EFS Access Point"
provider: "aws"
icon: "package"
order: 1
---

# Application Data Access Point

Enforced POSIX identity (1000:1000) with a dedicated `/app-data` root directory, auto-created on first mount. The standard shape for giving one application least-privilege access to a shared file system.

## When to Use

- ECS tasks or Lambda functions that need a confined slice of a shared EFS file system
- Any application that should not manage POSIX permissions itself
- Multi-tenant file systems where each app gets its own directory + identity

## What It Configures

- **POSIX identity 1000:1000** — every file operation uses this UID/GID, regardless of what the client runs as
- **Root directory `/app-data`** — the application sees this path as `/`; everything else is invisible
- **Creation info** — EFS creates `/app-data` owned by 1000:1000 with mode 0755 on first mount if it does not exist

## What to Customize

- Replace `<file-system-name>` with your AwsElasticFileSystem resource name
- Change the UID/GID to your application's runtime identity
- Tighten `permissions` (e.g., `"0750"`) for group-private data
- One access point per application — deploy this preset once per app with a distinct `path`
