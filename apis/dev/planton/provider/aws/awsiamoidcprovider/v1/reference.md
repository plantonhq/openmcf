# AwsIamOidcProvider

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsIamOidcProviderSpec defines an IAM OpenID Connect (OIDC) identity provider.

An IAM OIDC provider is the trust anchor that lets an external OIDC issuer's
short-lived tokens be exchanged for AWS credentials via STS
AssumeRoleWithWebIdentity -- no long-lived access keys required. It is the
missing link in two of the most common keyless patterns:
  - EKS IRSA (IAM Roles for Service Accounts): register the cluster's OIDC
    issuer here, then a Kubernetes ServiceAccount annotated with an IAM role
    ARN can assume that role. Point `url` at an AwsEksCluster's
    status.outputs.oidc_issuer_url to wire this with a single reference.
  - CI/CD federation (GitHub Actions, GitLab, etc.): register the CI provider's
    issuer so pipelines assume a deploy role without storing AWS secrets.

The resulting provider ARN (status.outputs.provider_arn) is what an AwsIamRole
trust policy references as a `Federated` principal, closing the loop:
issuer -> AwsIamOidcProvider -> AwsIamRole.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamOidcProvider
metadata:
  name: github-actions-oidc-demo
spec:
  region: us-west-2
  # GitHub Actions' OIDC issuer: a real, publicly reachable, well-known-CA issuer.
  # Because its CA is publicly trusted, thumbprintList is omitted and AWS derives it.
  # (For EKS IRSA, set url.valueFrom to an AwsEksCluster's status.outputs.oidc_issuer_url;
  # see presets/01-eks-irsa.yaml.)
  url:
    value: "https://token.actions.githubusercontent.com"
  # sts.amazonaws.com is the audience AWS expects for GitHub Actions web-identity tokens.
  clientIdList:
    - sts.amazonaws.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.url` | `string \| valueFrom` | yes |  | AwsEksCluster (`status.outputs.oidc_issuer_url`) |
| `spec.clientIdList` | `[]string` | yes |  |  |
| `spec.thumbprintList` | `[]string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region used to configure the provider. IAM is a global service, so
the OIDC provider is account-wide regardless of region; this value only
selects the regional STS/IAM endpoint the deploy talks to.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.url

`string | valueFrom` · required

url is the URL of the OIDC identity provider, corresponding to the issuer
(`iss`) claim in the tokens it mints. It must be a valid HTTPS URL with no
query or fragment (e.g. "https://oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED"
for EKS, or "https://token.actions.githubusercontent.com" for GitHub Actions).
AWS allows at most one OIDC provider per unique URL per account.
This field is ForceNew: changing it requires replacing the provider.

Reference an AwsEksCluster to enable IRSA without copying the issuer by hand:
this defaults to that cluster's status.outputs.oidc_issuer_url.

- references: AwsEksCluster (`status.outputs.oidc_issuer_url`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEksCluster, name: <that resource's name>, fieldPath: status.outputs.oidc_issuer_url}} -- a bare string does not parse

### spec.clientIdList

`[]string` · required

client_id_list is the set of client IDs (also called audiences) that are
allowed to authenticate against this provider -- the value tokens carry in
their `aud` claim. For EKS IRSA this is "sts.amazonaws.com"; for GitHub
Actions it is also typically "sts.amazonaws.com". At least one is required,
each must be 1-255 characters, and duplicates are not allowed.

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"string":{"minLen":"1","maxLen":"255"}}}}

### spec.thumbprintList

`[]string`

thumbprint_list is an optional set of SHA-1 fingerprints of the OIDC
provider's root CA certificate(s), each exactly 40 hexadecimal characters.
Leave this empty for issuers backed by a well-known certificate authority
(EKS, GitHub Actions, Google, etc.): AWS secures the TLS connection using its
trusted CA store and derives the thumbprint for you. Supply thumbprints only
for issuers whose CA is not publicly trusted. Duplicates are not allowed.

AWS quirk: once thumbprints are set they cannot be cleared in place -- the
update API rejects an empty list. Going from explicit thumbprints back to
"let AWS derive them" requires replacing the provider (delete + recreate).

- rule: {"repeated":{"unique":true,"items":{"string":{"pattern":"^[0-9a-fA-F]{40}$"}}}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsIamOidcProvider, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.provider_arn` | `string` | provider_arn is the Amazon Resource Name (ARN) of the IAM OIDC provider (e.g. "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED"). This is the value an AwsIamRole trust policy references as a `Federated` principal to grant web-identity (IRSA / CI federation) access. |
| `status.outputs.provider_url` | `string` | provider_url is the issuer URL AWS stored for this provider, with the "https://" scheme stripped (e.g. "oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED"). It matches the `<provider-url>` segment of provider_arn and is the value used to build the `<provider-url>:sub` / `<provider-url>:aud` condition keys in a role's trust policy. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.url` | AwsEksCluster | `status.outputs.oidc_issuer_url` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsIamRole | `spec.oidcTrust.providerArn` | `status.outputs.provider_arn` |
| AwsIamRole | `spec.oidcTrust.providerUrl` | `status.outputs.provider_url` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
