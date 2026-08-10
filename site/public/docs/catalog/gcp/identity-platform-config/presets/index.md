---
title: "Presets"
description: "Ready-to-deploy configuration presets for Identity Platform Config"
type: "preset-list"
componentSlug: "identity-platform-config"
componentTitle: "Identity Platform Config"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-email-password"
    rank: "01"
    title: "Email + Password"
    excerpt: "The minimal working sign-in surface: users register and sign in with an email address and a password. No external credentials, no IdP consoles."
  - slug: "02-google-signin"
    rank: "02"
    title: "Google Sign-in"
    excerpt: "Email/password plus \"Sign in with Google\" — the standard consumer-app pair. Google handles the second flow's credentials; your app never sees a password for those users."
  - slug: "03-enterprise-saml"
    rank: "03"
    title: "Enterprise SAML"
    excerpt: "Inbound SAML SSO from a corporate identity provider (Okta, Azure AD, Ping, ...) with MFA enabled and the app's domain authorized — the enterprise sign-in posture."
---

# Identity Platform Config Presets

Ready-to-deploy configuration presets for Identity Platform Config. Each preset is a complete manifest you can copy, customize, and deploy.
