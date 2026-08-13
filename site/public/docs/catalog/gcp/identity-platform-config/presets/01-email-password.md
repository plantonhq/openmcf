---
title: "Email + Password"
description: "The minimal working sign-in surface: users register and sign in with an email address and a password. No external credentials, no IdP consoles."
type: "preset"
rank: "01"
presetSlug: "01-email-password"
componentSlug: "identity-platform-config"
componentTitle: "Identity Platform Config"
provider: "gcp"
icon: "package"
order: 1
---

# Email + Password

The minimal working sign-in surface: users register and sign in with an
email address and a password. No external credentials, no IdP consoles.

## What it configures

- `signIn.email.enabled: true` with `passwordRequired: true` — classic
  email/password accounts. The flag is sent explicitly, so this manifest
  actively enables the method.

## Adjust before deploying

- **The target project** — the first deploy permanently initializes
  Identity Platform on it (billing required, no de-initialize). Set
  `projectId` deliberately or rely on the provider default.
- **passwordRequired** — set to `false` to allow passwordless email-link
  sign-in alongside (or instead of) passwords.
- Add `authorizedDomains` once your app's domain is known — OAuth
  redirects and hosted flows only work from authorized domains.

## When to choose something else

Password accounts carry credential-reset and breach-reuse burden. For
consumer apps, add the **Google Sign-in** preset's provider; for
workforce or B2B apps, the **Enterprise SAML** preset delegates
authentication to the customer's IdP entirely.
