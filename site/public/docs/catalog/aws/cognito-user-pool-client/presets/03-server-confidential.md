---
title: "Preset: Server-Side Confidential"
description: "**Rank**: 3"
type: "preset"
rank: "03"
presetSlug: "03-server-confidential"
componentSlug: "cognito-user-pool-client"
componentTitle: "Cognito User Pool Client"
provider: "aws"
icon: "package"
order: 3
---

# Preset: Server-Side Confidential

**Rank**: 3

## When to Use

- Traditional server-rendered web applications (the backend handles the OAuth exchange)
- Backends-for-frontends that authenticate users on the browser's behalf

## What It Provides

- Authorization Code grant with a generated client secret (confidential client; ForceNew)
- Rotating refresh tokens with a 30-second retry grace window
- End-user IP/user-agent propagated to threat protection (server-side flows otherwise hide the real client)
- Token revocation and user-enumeration protection

## What You Might Add

- Federated providers in `supportedIdentityProviders` -- reference `AwsCognitoIdentityProvider` resources
- `readAttributes` / `writeAttributes` to scope the attributes this application may touch
- `analyticsConfiguration` to publish sign-in events to a Pinpoint project
