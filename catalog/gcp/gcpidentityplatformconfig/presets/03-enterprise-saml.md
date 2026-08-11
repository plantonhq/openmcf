# Enterprise SAML

Inbound SAML SSO from a corporate identity provider (Okta, Azure AD,
Ping, ...) with MFA enabled and the app's domain authorized — the
enterprise sign-in posture.

## What it configures

- One `inboundSamlConfigs` entry (`saml.okta-prod`) carrying both sides
  of the SAML exchange: the IdP's entity ID, SSO URL, and signing
  certificate (`idpConfig`), and this project's callback URI and SP
  entity ID (`spConfig`).
- `mfa.state: ENABLED` with `PHONE_SMS` — users may enroll a second
  factor.
- `authorizedDomains` pinned to the app's domain.

## Adjust before deploying

- **idpConfig** — copy `idpEntityId`, `ssoUrl`, and the X.509 signing
  certificate from the IdP's metadata XML (public certificate, not a
  secret).
- **spConfig.callbackUri** — must be `https://`; register the same URI
  with the IdP. `spEntityId` is what the IdP knows this project as.
- **name** — `saml.`-prefixed and immutable; pick a durable slug per
  IdP connection (e.g. `saml.okta-prod`).

## When to choose something else

SAML at the project level signs everyone into ONE shared user pool. For
B2B SaaS where each customer brings their own IdP, put the SAML config
on a per-customer **GcpIdentityPlatformTenant** instead (arm
`multiTenant.allowTenants` here first).
