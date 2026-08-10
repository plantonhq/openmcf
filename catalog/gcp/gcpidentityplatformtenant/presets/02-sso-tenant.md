# SSO Tenant

A tenant whose customer brings their own corporate IdP: users
authenticate against the customer's SAML directory (Okta, Azure AD,
Ping, ...), and no passwords live in the tenant at all.

## What it configures

- One `inboundSamlConfigs` entry scoped to this tenant only — the
  customer's SAML connection touches no other tenant.
- `spConfig` with BOTH `callbackUri` (https) and `spEntityId` — required
  at tenant level, unlike the project level where it is optional.
- `deletionPolicy: PREVENT` — this tenant holds a customer's user pool.

## Adjust before deploying

- **idpConfig** — copy `idpEntityId`, `ssoUrl`, and the X.509 signing
  certificate from the customer IdP's metadata XML (public certificate,
  not a secret).
- **spConfig.callbackUri** — must be `https://`; register the same URI
  with the customer's IdP. Use a per-tenant `spEntityId` so each
  customer's IdP sees a distinct service provider.
- **name** — `saml.`-prefixed and immutable; a durable per-customer slug
  (e.g. `saml.acme-okta`).

## When to choose something else

If the customer's IdP speaks OIDC rather than SAML, use
`oauthIdpConfigs` instead — note that at tenant level the OIDC
`displayName` is required and there is no `responseType` selection. For
customers without an IdP, the **B2B Tenant** preset's password sign-up
is the simpler start.
