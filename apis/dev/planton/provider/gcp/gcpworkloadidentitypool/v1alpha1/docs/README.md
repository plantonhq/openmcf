# GCP Workload Identity Pools: The End of Service-Account Keys

## The Problem with Keys

The default way external systems authenticate to Google Cloud — a downloaded service-account JSON key — is the single most common GCP security failure. Keys are long-lived bearer credentials: anyone holding the file *is* the service account, forever, from anywhere. They leak through CI logs, laptop backups, git history, and container images; they demand rotation machinery nobody builds; and every security framework flags them on sight.

Workload Identity Federation replaces the key with a handshake. An external workload proves its identity to its *own* issuer (GitHub mints an OIDC token for a workflow run; AWS signs a caller identity; an enterprise IdP asserts a SAML subject; a private CA issues a client certificate), and GCP's Security Token Service exchanges that short-lived proof for short-lived Google credentials. Nothing long-lived exists to leak.

The pool is where that trust is anchored.

## Pool vs Provider: Why Two Nodes

The pool deliberately holds no issuer configuration. It is:

- **The principal namespace.** Every federated identity becomes a principal embedding the pool's resource name: `principal://iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/subject/<subject>`, or `principalSet://.../attribute.<name>/<value>` for groups of identities. IAM bindings target these strings, so the pool name is the stable handle access control depends on.
- **The trust boundary.** Disabling the pool cuts off every issuer attached to it at once — the emergency lever for a suspected federation compromise.
- **The container for providers.** One pool typically serves several issuers: the same boundary can trust a GitHub org and an AWS account, each through its own provider with its own attribute mapping and conditions.

Bundling the issuer into the pool would force one-pool-per-issuer and make the principal namespace churn with issuer changes. The provider is a genuinely separate, many-per-parent resource — which is exactly why it is a separate composable node (GcpWorkloadIdentityPoolProvider).

## Lifecycle Sharp Edges

**1. Identity is immutable.** `workloadIdentityPoolId`, the owning project, and `mode` can never change. Everything downstream — principals, provider parents, token audiences — embeds the pool name, so a recreate invalidates the entire federation surface. Choose pool IDs the way you choose domain names.

**2. Deletion is soft, and creation does NOT undelete.** A deleted pool sits in `DELETED` state for ~30 days: token exchanges fail, the ID stays reserved, and — unlike IAM custom roles — creating a pool with the soft-deleted ID fails outright rather than undeleting it. Restoring requires an explicit `UndeleteWorkloadIdentityPool` call outside the create path. The operational consequence: **never delete a pool to "reset" it.** Use `disabled: true` for shutoffs; it is instant and reversible.

**3. Disabled is the kill switch.** A disabled pool rejects all token exchanges and existing tokens stop granting access; re-enabling restores them. This is the right lever for incident response and rotation windows.

## Modes: Federation vs Trust Domain

`FEDERATION_ONLY` (the default) is the keyless-auth workhorse: external identities in, Google credentials out, no structure imposed on subjects, providers allowed.

`TRUST_DOMAIN` turns the pool into a SPIFFE-style trust domain that *assigns* managed identities to Google Cloud workloads (subjects shaped `ns/<namespace>/sa/<workload>`), issues mTLS workload certificates (via `inlineCertificateIssuanceConfig`), and can extend trust to foreign domains (`inlineTrustConfig`). Trust-domain pools cannot hold providers — the two modes are genuinely different machines sharing a resource type.

A third mode, `SYSTEM_TRUST_DOMAIN`, exists only for pools Google itself manages; it cannot be created and is rejected by the spec with an explaining message.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `workload_identity_pool_id` | ✅ `workloadIdentityPoolId` | Same validation as the API (4-32 chars, `gcp-` reserved) |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; falls back to the provider default project |
| `display_name` / `description` | ✅ | Max 32 / 256 chars |
| `disabled` | ✅ `disabled` | The kill switch |
| `mode` | ✅ `mode` | `FEDERATION_ONLY` default; `SYSTEM_TRUST_DOMAIN` rejected (Google-managed only) |
| `inline_certificate_issuance_config` | ✅ | Region→CA-pool map, key algorithm, lifetime, rotation window |
| `inline_trust_config` | ✅ | Foreign trust domains with PEM trust anchors |
| `name` / `state` (computed) | ✅ outputs | The composable handle + soft-delete visibility |
| `deletion_policy` | ❌ | A Terraform-provider-level abandon-vs-delete lever, not a property of the pool; Planton's lifecycle management owns this concern |
| `timeouts` | ❌ | Operation plumbing, not resource configuration |

The sibling resources for managed workload identities — pool *namespaces* and *managed identities* (with their attestation rules) — are deliberately not modeled. They serve the SPIFFE managed-identity niche, not keyless federation, and would earn their place as first-class kinds if that demand materializes rather than as bolted-on fields here.

## Composition: the Keyless Path End to End

A working keyless setup composes four nodes:

1. **GcpWorkloadIdentityPool** (this component) — the trust boundary.
2. **GcpWorkloadIdentityPoolProvider** — the issuer, referencing the pool's `workload_identity_pool_id` output.
3. **GcpServiceAccount** — the identity to impersonate, with `roles/iam.workloadIdentityUser` granted to the pool's principals.
4. **GcpProjectIamMember** — grants on the service account (or directly to `principal://` members for a service-account-free setup).

The pool's `name` output is the string every principal embeds; the provider's `name` output is the audience tokens are minted for. Both are exported precisely so downstream nodes reference them instead of re-deriving them.

## Operational Guidance

- **Pool-per-boundary, not pool-per-issuer**: group issuers that share a trust posture into one pool; separate pools only where you want independent kill switches and principal namespaces (e.g. prod vs ci).
- **Quota**: pools are a limited resource per project (default quota is low double digits) — another reason to prefer few well-named pools over many ad-hoc ones.
- **Auditing**: the `description` field is what a reviewer sees when walking trust boundaries; write it for them.
