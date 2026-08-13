---
title: "Presets"
description: "Ready-to-deploy configuration presets for Identity Platform Tenant"
type: "preset-list"
componentSlug: "identity-platform-tenant"
componentTitle: "Identity Platform Tenant"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-b2b-tenant"
    rank: "01"
    title: "B2B Tenant"
    excerpt: "One isolated user pool for one customer organization — the standard B2B SaaS shape. Users sign up with email/password inside the tenant and never mix with any other customer's pool."
  - slug: "02-sso-tenant"
    rank: "02"
    title: "SSO Tenant"
    excerpt: "A tenant whose customer brings their own corporate IdP: users authenticate against the customer's SAML directory (Okta, Azure AD, Ping, ...), and no passwords live in the tenant at all."
---

# Identity Platform Tenant Presets

Ready-to-deploy configuration presets for Identity Platform Tenant. Each preset is a complete manifest you can copy, customize, and deploy.
