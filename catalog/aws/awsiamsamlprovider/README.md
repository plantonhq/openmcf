<p align="center">
  <img src="logo.svg" alt="AWS IAM SAML Provider" width="80"/>
</p>

# AWS IAM SAML Provider

Manage an [IAM SAML identity provider](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_saml.html)
— the account's trust anchor for SAML 2.0 federation from your
corporate identity provider (Okta, Entra ID, Google Workspace, ...).

## What Gets Managed

- **The provider** (named from `metadata.name` — WRITE-ONCE at AWS: a
  rename replaces the provider and invalidates every role trust
  policy naming its ARN), created from the IdP's
  **samlMetadataDocument** — the XML the IdP publishes, carrying its
  issuer and public signing certificates. A PUBLIC trust document,
  not a secret.
- The metadata updates IN PLACE — certificate rotations are ordinary
  updates, and the `valid_until` output is the rotate-by date.

IAM roles then trust the provider via `sts:AssumeRoleWithSAML` in
their trust policies ([AWS IAM Role](../awsiamrole)), and federated
users sign in through the IdP straight into assumed roles.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
