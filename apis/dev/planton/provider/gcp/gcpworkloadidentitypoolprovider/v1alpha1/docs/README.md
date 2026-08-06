# GCP Workload Identity Pool Providers: Where Federation Trust Is Configured

## The Provider Is the Contract

A Workload Identity Pool declares that a trust boundary exists; the provider declares *who* is trusted and *on what terms*. Every keyless authentication into Google Cloud walks through a provider:

1. The external workload obtains proof from its own issuer — an OIDC token, AWS credentials, a SAML assertion, or a client certificate.
2. GCP's Security Token Service validates that proof against the provider's issuer configuration.
3. The provider's `attribute_mapping` translates the issuer's claims into Google attributes; its `attribute_condition` decides whether this particular credential is accepted at all.
4. Out comes a short-lived Google credential for `principal://…/subject/<mapped subject>`.

The provider's full resource name is the **audience** of the whole handshake: tokens are minted *for* this provider, and anything consuming web-identity credentials (CI pipelines, deployment tooling, provider configurations) carries this exact string. That is why it is the primary stack output.

## The Four Issuer Types

Exactly one issuer arm is configured per provider, and the type is fixed for the provider's lifetime (the API rejects cross-type updates):

- **`oidc`** — the workhorse. Any OpenID Connect issuer: GitHub Actions, GitLab CI, Kubernetes cluster issuers, Buildkite, HashiCorp Vault, custom issuers. Signing keys come from the issuer's `.well-known` discovery document, or inline `jwksJson` for issuers unreachable from the internet.
- **`aws`** — trusts one AWS account by ID; EC2/Lambda/ECS workloads federate with their native AWS credentials, no token plumbing needed.
- **`saml`** — trusts an enterprise IdP via its metadata XML; the classic choice where the identity source is ADFS/Okta/PingFederate asserting non-interactive workloads.
- **`x509`** — trusts a certificate authority; clients presenting certificates chaining to the configured anchors federate with pure mTLS, no token issuer at all. The trust store takes root/intermediate anchors plus optional chain-building intermediates.

## Attribute Mapping: the Security-Critical Part

`attribute_mapping` is a map of Google attribute → CEL expression over the issuer's claims (`assertion.*`):

- `google.subject` (required for OIDC; defaulted for the other issuers) becomes the principal's identity — what appears in IAM bindings and audit logs.
- `attribute.<name>` entries create IAM-targetable *groups*: `principalSet://…/attribute.repository/my-org/my-repo` matches every credential whose mapped `repository` equals that value.

Two rules keep mappings safe:

1. **Map only claims the issuer actually signs.** Mapping `assertion.email` from an issuer that lets users set arbitrary emails hands out identity spoofing.
2. **Grant to the narrowest principal set that works.** Prefer `attribute.repository` over `attribute.repository_owner` over pool-wide `principalSet` grants.

## attribute_condition: Mandatory for Multi-Tenant Issuers

GitHub's OIDC issuer signs tokens for **every repository on github.com**. A provider trusting `https://token.actions.githubusercontent.com` with no condition accepts all of them. The condition is what scopes trust:

```
assertion.repository_owner == "my-org"
```

This is the single most common federation misconfiguration; the component's docs, presets, and examples all carry a condition for GitHub-style issuers.

## Lifecycle Sharp Edges

- **Immutability**: pool, provider ID, project, and issuer *type* are fixed at create. Mapped attributes, conditions, and the issuer arm's contents update in place. Note the API constraint that mapped attributes cannot be *removed* on update (only added or redefined) — a redefinition to a constant is the practical way to retire one.
- **Soft delete without undelete-on-create**: like pools, deleted providers linger ~30 days in `DELETED` state, their IDs reserved, and a create against the soft-deleted ID fails. Prefer `disabled: true`.
- **Disabling** stops new token exchanges but does not revoke already-issued Google credentials — they age out on their own TTL (an hour by default).

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `workload_identity_pool_id` | ✅ ref → GcpWorkloadIdentityPool | Resolved from the pool's outputs |
| `workload_identity_pool_provider_id` | ✅ | Same validation as the API |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject |
| `display_name` / `description` / `disabled` | ✅ | |
| `attribute_mapping` / `attribute_condition` | ✅ | OIDC's `google.subject` requirement validated pre-deploy |
| `aws` / `oidc` / `saml` / `x509` | ✅ oneof | Full arms including x509 trust store with intermediates |
| `name` / `state` (computed) | ✅ outputs | The audience handle + soft-delete visibility |
| `deletion_policy` | ❌ | Terraform-provider abandon-vs-delete lever; Planton's lifecycle management owns this concern |
| `timeouts` | ❌ | Operation plumbing, not resource configuration |

## Composition

The provider references its pool by the pool's `workload_identity_pool_id` output and the owning project by `projectId`. Downstream, its `name` output feeds:

- **Token minting**: CI systems set the token audience to `//iam.googleapis.com/<name>`.
- **Web-identity provider configurations**: keyless deploy credentials consume `<name>` as their audience verbatim.
- **IAM grants**: `principal://` and `principalSet://` members embed the *pool* name with attributes this provider mapped — grant `roles/iam.workloadIdentityUser` on a GcpServiceAccount for impersonation, or roles directly for a service-account-free setup.
