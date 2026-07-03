# Generic External OIDC Issuer

This preset trusts any OIDC-compliant external system -- GitLab, a
self-hosted CI, another cloud's workload identity, an internal token service
-- to authenticate as a managed identity. Azure AD only needs the issuer to
publish standard OIDC discovery metadata
(`{issuer}/.well-known/openid-configuration`) over HTTPS; it validates
incoming tokens against the issuer's published signing keys on every
exchange.

The trust is the exact three-way claim match: the token's `iss`, `sub`, and
`aud` must equal this credential's issuer, subject, and audience. Consult the
external system's documentation for the `sub` format its tokens carry --
e.g. GitLab (issuer `https://gitlab.com`) mints subjects like
`project_path:{group}/{project}:ref_type:branch:ref:{branch}`.

## When to Use

- GitLab CI, Buildkite, or self-hosted pipeline systems deploying to Azure
  without stored secrets
- Cross-system automation where the caller already has an OIDC identity
- Standards-based integrations you want reviewable as infrastructure

## Key Configuration Choices

- **Audience stays the default** unless the external system cannot mint
  `api://AzureADTokenExchange` -- Azure AD accepts the exchange only when the
  token's `aud` matches, and the default is what standard clients request
- **Exact issuer URL** -- the `iss` comparison is character-exact including
  any trailing slash; copy it from the system's OIDC discovery document

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-assigned-identity-arm-id>` | The parent identity's ARM ID (or use `valueFrom` against an `AzureUserAssignedIdentity`) | The identity's `status.outputs.identity_id` |
| `<external-oidc-issuer-url>` | The system's OIDC issuer URL | Its `/.well-known/openid-configuration` document |
| `<external-subject>` | The `sub` claim its tokens carry | The external system's OIDC documentation |
