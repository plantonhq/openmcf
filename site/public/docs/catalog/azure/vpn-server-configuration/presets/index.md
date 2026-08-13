---
title: "Presets"
description: "Ready-to-deploy configuration presets for VPN Server Configuration"
type: "preset-list"
componentSlug: "vpn-server-configuration"
componentTitle: "VPN Server Configuration"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-entra-remote-workforce"
    rank: "01"
    title: "Entra ID Remote Workforce"
    excerpt: "This preset authenticates remote users with Entra ID (Azure AD): sign-in rides your tenant's conditional access and MFA, revoking a person is an identity operation, and no certificates are..."
  - slug: "02-certificate-auth-policy-groups"
    rank: "02"
    title: "Certificate Auth with Policy Groups"
    excerpt: "This preset authenticates clients against your own root certificate and segments them into policy groups by the certificate's common name -- engineering and contractors here. A point-to-site gateway..."
---

# VPN Server Configuration Presets

Ready-to-deploy configuration presets for VPN Server Configuration. Each preset is a complete manifest you can copy, customize, and deploy.
