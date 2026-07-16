---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Secret"
type: "preset-list"
componentSlug: "front-door-secret"
componentTitle: "Front Door Secret"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-rotating-byo-certificate"
    rank: "01"
    title: "Rotating Bring-Your-Own Certificate"
    excerpt: "This preset creates a Front Door secret wrapping a Key Vault certificate by its VERSIONLESS id -- the rotation-follows-latest posture: when Key Vault renews or you upload a new version, Front Door..."
  - slug: "02-pinned-certificate-version"
    rank: "02"
    title: "Pinned Certificate Version"
    excerpt: "This preset creates a Front Door secret pinned to ONE exact Key Vault certificate version -- Front Door keeps serving that version until the secret itself is replaced, no matter what rotation happens..."
---

# Front Door Secret Presets

Ready-to-deploy configuration presets for Front Door Secret. Each preset is a complete manifest you can copy, customize, and deploy.
