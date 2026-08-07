# AzureFederatedIdentityCredential — Research Notes

Design notes for the federated identity credential component: what the
resource is, how the trust model works, and how the spec maps to the
provider surface.

## What a Federated Identity Credential Is

Workload identity federation is Azure AD's mechanism for letting an
**external** workload authenticate as an Azure managed identity without any
Azure-issued secret. The external system (GitHub, an AKS cluster, GitLab,
anything that issues OIDC tokens) gives its workload a signed JWT; the
workload presents that JWT to Azure AD's token endpoint; Azure AD validates
the JWT against the issuer's published signing keys and, if a federated
identity credential on the target identity matches the token's claims,
issues Azure credentials for that identity.

A federated identity credential is one such matching rule, stored as an ARM
child resource of a user-assigned managed identity:

```
{identity-id}/federatedIdentityCredentials/{name}
```

### The three-way claim match

Azure AD performs an exact string match on three claims of the incoming
token:

| Credential field | Token claim | Semantics |
|------------------|-------------|-----------|
| `issuer` | `iss` | The external identity provider's OIDC issuer URL. Azure AD fetches `{issuer}/.well-known/openid-configuration` to discover signing keys. Exact match — scheme, host, path, trailing slash. |
| `subject` | `sub` | The specific workload within the issuer. Exact match, **no wildcards**. |
| `audience` | `aud` | Who the token was minted for. `api://AzureADTokenExchange` is the value Azure AD's exchange endpoint expects and every standard client requests. |

Because matching is exact, the trust is precisely as narrow as the subject
string. One credential = one trusted workload. Azure allows up to 20
credentials per identity, so an identity accumulates one rule per branch,
environment, or service account that acts as it.

### What a credential does NOT do

A credential only **authenticates** the external workload as the identity.
It conveys zero permissions. What the identity may do is granted separately
through Azure RBAC (`AzureRoleAssignment`). The complete keyless story is a
three-kind composition:

1. `AzureUserAssignedIdentity` — the identity itself
2. `AzureRoleAssignment` — what the identity may do
3. `AzureFederatedIdentityCredential` — who may act as the identity

## Headline Flows

### GitHub Actions (keyless CI/CD)

GitHub's OIDC issuer is `https://token.actions.githubusercontent.com`. A
workflow requests a token (`permissions: id-token: write`), and `azure/login`
exchanges it. The subject encodes what is running:

- `repo:{owner}/{repo}:ref:refs/heads/{branch}` — push/workflow runs on a branch
- `repo:{owner}/{repo}:environment:{env}` — jobs targeting a protected environment
- `repo:{owner}/{repo}:ref:refs/tags/{tag}` — tag builds
- `repo:{owner}/{repo}:pull_request` — pull-request runs

The environment form composes best with GitHub's own controls (required
reviewers on the environment gate the Azure access, not just the branch
protection).

### AKS workload identity

An AKS cluster with the OIDC issuer enabled publishes its own issuer URL
(exported by the cluster). The workload-identity webhook projects a
service-account token into annotated pods; the Azure SDK exchanges it. The
subject is always:

```
system:serviceaccount:{namespace}:{serviceaccount}
```

This replaces the deprecated AAD Pod Identity approach and is the reason
this kind is the AKS-workload-identity unlock.

## Provider Surface Map (azurerm v4)

`azurerm_federated_identity_credential`
(`internal/services/managedidentity/federated_identity_credential_resource.go`):

| azurerm argument | Spec field | Notes |
|------------------|-----------|-------|
| `name` (required, ForceNew) | `name` | ARM name under the parent identity; 3-120 chars |
| `user_assigned_identity_id` (required, ForceNew) | `user_assigned_identity` | The v4-canonical argument; `parent_id` is a deprecated alias removed in v5 |
| `resource_group_name` (deprecated, computed) | — (derived) | v4 derives it from the parent ID; v5 removes it. The spec deliberately has no resource-group field — it would restate derivable state that could then contradict the parent |
| `issuer` (required) | `issuer` | In-place update. Modeled literal-or-reference: a literal URL for external providers, or a reference defaulting to an `AzureAksCluster`'s `oidc_issuer_url` output for the workload-identity composition |
| `subject` (required) | `subject` | In-place update |
| `audience` (required, list, MaxItems 1) | `audience` (single string, defaulted) | ARM's wire shape is a single-element list; the spec models the one value directly and defaults it to `api://AzureADTokenExchange` |

Engine parity notes:

- The Pulumi provider (`armmsi.FederatedIdentityCredential`) flattens the
  audience to a single string and still requires `resource_group_name` +
  `parent_id` explicitly; the module derives the resource group by parsing
  the resolved parent ARM ID so both engines take identical spec input.
- Both providers serialize credential writes per parent identity (ARM
  rejects concurrent writes on one identity), so multiple credentials on the
  same identity apply sequentially by design.

### Deliberately not modeled (recorded reasons)

- **`resource_group_name`** — deprecated in azurerm v4, removed in v5,
  derivable from the parent identity's ARM ID. Modeling it would add
  contradictable redundant state.
- **Flexible federated identity credentials (claims-matching expressions,
  wildcard subjects)** — an Azure AD preview not modeled by azurerm v4 at
  all (the provider models exactly issuer/subject/audience). Revisit when
  the provider grows the surface.
- **System-assigned identity parents** — ARM only supports federated
  credentials on *user-assigned* identities; there is nothing to model.

## Operational Notes

- **Propagation**: a newly written credential can take a short time (up to a
  few minutes in ARM's documentation) to become effective for token
  exchange; retries in standard clients absorb this.
- **Deletion**: removing the credential immediately revokes the external
  workload's ability to obtain NEW tokens; already-issued Azure tokens live
  out their (typically 1h) lifetime.
- **Issuer key rotation**: Azure AD discovers signing keys from the issuer's
  OIDC metadata on each validation; external key rotation needs no Azure
  change.
- **Naming**: the credential name is only unique per identity, but naming it
  after the trusted workload ("github-main-branch") keeps the portal's
  Federated credentials tab self-documenting.
