# Corporate IdP

This preset creates the account's federation trust anchor from your
identity provider's published metadata — the one-time setup that lets
corporate sign-ins assume IAM roles.

## When to Use

- Connecting Okta, Entra ID, Google Workspace, or any SAML 2.0 IdP to
  the account
- Replacing long-lived IAM user credentials with federated,
  short-lived role sessions

## What You Get

- The IAM SAML provider, created from the pasted metadata XML (a
  public document — issuer, endpoints, signing certificates)
- The provider ARN as an output for role trust policies, plus the
  certificate expiry date to rotate by

## Customize

- Wire roles to it: a trust policy with
  `Principal: {Federated: <provider_arn>}`, action
  `sts:AssumeRoleWithSAML`, and the `SAML:aud` condition pinned to
  AWS's sign-in endpoint
- Configure the IdP's AWS app with `role_arn,provider_arn` pairs —
  both sides must agree
- The name is write-once: renaming replaces the provider and breaks
  trust policies naming the old ARN
