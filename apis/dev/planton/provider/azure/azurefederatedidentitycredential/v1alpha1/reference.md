# AzureFederatedIdentityCredential

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFederatedIdentityCredentialSpec** defines the configuration for
creating a federated identity credential: a keyless trust rule on a
user-assigned managed identity that lets an external workload exchange its
own OIDC token for Azure credentials -- no client secret, nothing to
rotate, nothing to leak.

The trust model is a three-way match. An external workload presents a
token; Azure AD accepts the exchange only when the token's `iss` claim
matches this credential's issuer, its `sub` claim matches the subject, and
its `aud` claim matches the audience. Each credential encodes exactly one
such rule; an identity carries one credential per external workload that
should be able to act as it (Azure allows up to 20 per identity).

The two headline flows this unlocks:
- **CI without secrets (GitHub Actions and friends)**: a workflow's OIDC
  token is exchanged for the identity's credentials, so pipelines deploy to
  Azure with no stored service-principal secret.
- **AKS workload identity**: a Kubernetes service account's projected token
  is exchanged for the identity's credentials, so pods reach Key Vault,
  Storage, or any RBAC-granted resource without node-level credentials.

A credential conveys no permissions by itself -- it only authenticates the
external workload AS the identity. What that identity may do is granted
separately through role assignments; compose this kind with
AzureUserAssignedIdentity (the parent) and AzureRoleAssignment (the
grants).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFederatedIdentityCredential
metadata:
  name: test-federated-credential
  org: test-org
  env: dev
spec:
  name: github-main-branch
  userAssignedIdentity:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-identity
  issuer:
    value: https://token.actions.githubusercontent.com
  subject: repo:test-org/platform:ref:refs/heads/main
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.userAssignedIdentity` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.issuer` | `string \| valueFrom` | yes |  | AzureAksCluster (`status.outputs.oidc_issuer_url`) |
| `spec.subject` | `string` | yes |  |  |
| `spec.audience` | `string` | yes | `api://AzureADTokenExchange` |  |

## Field Details

### spec.name

`string` · required

The credential's name -- its ARM resource name under the parent
identity, shown in the portal's "Federated credentials" tab. Must be
unique within the parent identity (3-120 characters: alphanumeric,
hyphens, and underscores). Name it after the external workload it
trusts, e.g. "github-main-branch" or "aks-payments-serviceaccount".
Changing the name replaces the credential (delete + create).

- rule: {"required":true,"string":{"minLen":"3","maxLen":"120"}}

### spec.userAssignedIdentity

`string | valueFrom` · required

The user-assigned managed identity this credential is written on -- the
identity the external workload will act as after the token exchange.
Takes the identity's full ARM resource ID; defaults to referencing an
AzureUserAssignedIdentity's identity_id output in composed environments.
Changing the parent replaces the credential (it is a child resource of
the identity, and the resource group is derived from this ID).

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.issuer

`string | valueFrom` · required

The OIDC issuer URL of the external identity provider -- the value the
incoming token's `iss` claim must equal, exactly. Must be a full URL
(scheme, host, and any path; no trailing-slash forgiveness -- Azure
matches the string byte for byte). Takes a literal URL for external
issuers, or a reference for issuers that are themselves resources in the
environment; defaults to referencing an AzureAksCluster's
oidc_issuer_url output -- the workload-identity composition where the
trusted issuer only exists once the cluster is deployed. Well-known
literal issuers:
- GitHub Actions: "https://token.actions.githubusercontent.com"
- GitLab: "https://gitlab.com"
Updatable in place.

- references: AzureAksCluster (`status.outputs.oidc_issuer_url`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureAksCluster, name: <that resource's name>, fieldPath: status.outputs.oidc_issuer_url}} -- a bare string does not parse

### spec.subject

`string` · required

The identifier of the external workload -- the value the incoming
token's `sub` claim must equal, exactly. Each issuer has its own subject
format; the trust is only as narrow as this string, so prefer the most
specific form the issuer offers:
- GitHub Actions (branch): "repo:{owner}/{repo}:ref:refs/heads/{branch}"
- GitHub Actions (environment): "repo:{owner}/{repo}:environment:{env}"
- GitHub Actions (tags): "repo:{owner}/{repo}:ref:refs/tags/{tag}"
- AKS workload identity: "system:serviceaccount:{namespace}:{serviceaccount}"
No wildcards: one credential trusts one subject. Add one credential per
branch/environment/service-account that needs access. Updatable in place.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.audience

`string` · required · optional (explicit presence)

The audience the incoming token must be issued for -- the value of its
`aud` claim. Defaults to "api://AzureADTokenExchange", the audience
Azure AD's token-exchange endpoint expects and the value every standard
client (azure-sdk, azure/login, AKS workload identity webhook) requests;
override only when a provider genuinely mints a different audience.
Azure accepts exactly one audience per credential today (the ARM API
models it as a single-element list). Updatable in place.

- default: `api://AzureADTokenExchange`
- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFederatedIdentityCredential, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.federated_identity_credential_id` | `string` | The full Azure Resource Manager ID of the federated identity credential. Format: {identity-id}/federatedIdentityCredentials/{name} where {identity-id} is the parent user-assigned identity's ARM ID. |
| `status.outputs.name` | `string` | The credential's ARM resource name under its parent identity, as deployed (what the portal lists in the identity's "Federated credentials" tab). |
| `status.outputs.user_assigned_identity_id` | `string` | The ARM ID of the parent user-assigned identity this credential is written on -- the identity the trusted external workload acts as. |
| `status.outputs.issuer` | `string` | The OIDC issuer URL the trust matches against the incoming token's `iss` claim, as deployed. |
| `status.outputs.subject` | `string` | The workload identifier the trust matches against the incoming token's `sub` claim, as deployed. |
| `status.outputs.audience` | `string` | The audience the trust matches against the incoming token's `aud` claim, as deployed -- "api://AzureADTokenExchange" unless overridden. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.userAssignedIdentity` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.issuer` | AzureAksCluster | `status.outputs.oidc_issuer_url` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
