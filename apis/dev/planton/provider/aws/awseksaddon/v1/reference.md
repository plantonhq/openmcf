# AwsEksAddon

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsEksAddonSpec installs an EKS MANAGED add-on on an AwsEksCluster:
cluster software (vpc-cni, coredns, kube-proxy, the EBS/EFS CSI
drivers, Pod Identity agent, marketplace add-ons, ...) whose lifecycle
-- install, version upgrades, configuration, removal -- AWS manages
through the EKS control plane instead of Helm charts you operate
yourself.

The add-on composes onto its neighbors instead of embedding them: the
cluster attaches by reference (status.outputs.name), and the IAM
identity the add-on's pods run as -- when it needs one beyond the node
role -- is a referenced AwsIamRole, wired either through IRSA
(service_account_role_arn, requires the cluster's OIDC provider) or
through EKS Pod Identity (pod_identity_associations, requires the
eks-pod-identity-agent add-on). This component never modifies a role
it merely references.

One add-on name per cluster: AWS keys the add-on on
(cluster, addon_name), so a second resource with the same pair
conflicts. The add-on name and namespace config are create-time
immutable; version, configuration, and the IAM wiring update in place.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksAddon
metadata:
  name: awseksaddon-demo
spec:
  region: us-west-2
  clusterName:
    value: awsekscluster-demo
  addonName: coredns
  resolveConflictsOnCreate: OVERWRITE
  resolveConflictsOnUpdate: OVERWRITE
  configurationValues: '{"replicaCount":2}'
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.clusterName` | `string \| valueFrom` | yes |  | AwsEksCluster (`status.outputs.name`) |
| `spec.addonName` | `string` | yes |  |  |
| `spec.addonVersion` | `string` |  |  |  |
| `spec.resolveConflictsOnCreate` | `string` |  |  |  |
| `spec.resolveConflictsOnUpdate` | `string` |  |  |  |
| `spec.configurationValues` | `string` |  |  |  |
| `spec.serviceAccountRoleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.podIdentityAssociations` | `[]AwsEksAddonPodIdentityAssociation` |  |  |  |
| `spec.podIdentityAssociations[].roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.podIdentityAssociations[].serviceAccount` | `string` | yes |  |  |
| `spec.preserve` | `bool` |  |  |  |
| `spec.namespaceConfig` | `AwsEksAddonNamespaceConfig` |  |  |  |
| `spec.namespaceConfig.namespace` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the add-on's cluster lives in. Must match the
cluster's region. Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.clusterName

`string | valueFrom` · required

The EKS cluster the add-on installs on. Reference an AwsEksCluster's
name output or pass a literal cluster name for a cluster managed
outside Planton. Create-only in AWS.

- references: AwsEksCluster (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEksCluster, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.addonName

`string` · required

The add-on to install, by its EKS catalog name: AWS-built add-ons
("vpc-cni", "coredns", "kube-proxy", "aws-ebs-csi-driver",
"aws-efs-csi-driver", "eks-pod-identity-agent", "snapshot-controller",
"amazon-cloudwatch-observability", ...) or a marketplace add-on's
vendor-prefixed name. `aws eks describe-addon-versions` lists what the
cluster's Kubernetes version supports. Create-only in AWS -- changing
the name replaces the add-on.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"100"}}

### spec.addonVersion

`string`

The add-on version to run, e.g. "v1.18.1-eksbuild.3". Empty installs
the AWS default version for the cluster's Kubernetes version -- the
never-goes-stale choice; pin a version for byte-identical clusters
and controlled upgrades. Version changes update in place (the add-on
rolls its own pods).

### spec.resolveConflictsOnCreate

`string`

How to resolve conflicts when the add-on's Kubernetes resources
already exist on the cluster at install time -- typically because
the cluster bootstrapped self-managed copies of vpc-cni / coredns /
kube-proxy at creation. "OVERWRITE" adopts and overwrites the
existing install (the standard way to migrate a self-managed add-on
to a managed one); "NONE" (the default when empty) fails the install
with a conflict instead. AWS accepts only these two at create time
-- "PRESERVE" exists only for updates.

### spec.resolveConflictsOnUpdate

`string`

How to resolve drift when an update finds fields changed out-of-band
(e.g. someone kubectl-edited the add-on's config). "OVERWRITE"
restores the managed values, "PRESERVE" keeps the hand-made changes,
"NONE" (the default when empty) fails the update on conflict.

### spec.configurationValues

`string`

Custom configuration for the add-on as a single JSON document (e.g.
'{"replicaCount":3}' for coredns). Each add-on publishes its own
schema -- `aws eks describe-addon-configuration` shows what is
configurable. Empty keeps every add-on default. Updates in place.

### spec.serviceAccountRoleArn

`string | valueFrom`

The IAM role the add-on's service account assumes via IRSA. Requires
an IAM OIDC provider for the cluster (AwsIamOidcProvider on the
cluster's oidc_issuer_url output); without it AWS rejects the
install. Empty means the add-on's pods fall back to the node role's
permissions. Prefer pod_identity_associations on new clusters --
Pod Identity needs no per-cluster OIDC provider. Reference an
AwsIamRole's role_arn output or pass a literal ARN. Updates in place.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.podIdentityAssociations

`[]AwsEksAddonPodIdentityAssociation`

EKS Pod Identity wiring: bind the add-on's service account(s) to IAM
roles through the Pod Identity agent (the modern, no-OIDC-provider
alternative to IRSA). The role must trust
"pods.eks.amazonaws.com"; the eks-pod-identity-agent add-on must be
installed on the cluster. Updates in place.

### spec.podIdentityAssociations[].roleArn

`string | valueFrom` · required

The IAM role the service account assumes. It must trust
"pods.eks.amazonaws.com" and carry the permissions the add-on's
documentation requires (e.g. the EBS CSI driver policy). Reference
an AwsIamRole's role_arn output or pass a literal ARN.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.podIdentityAssociations[].serviceAccount

`string` · required

The Kubernetes service account (in the add-on's namespace) the role
binds to -- each add-on documents its service account name(s), e.g.
"ebs-csi-controller-sa" for the EBS CSI driver.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.preserve

`bool`

Keep the add-on's Kubernetes resources on the cluster when this
resource is deleted -- AWS stops managing them but leaves them
running (they become self-managed). Off by default: deletion removes
the software. Turn it on when handing an add-on's lifecycle back to
cluster operators without an outage.

### spec.namespaceConfig

`AwsEksAddonNamespaceConfig`

Install the add-on into a custom namespace instead of its default.
Create-only in AWS: changing the namespace requires removing and
re-creating the add-on. Only some add-ons support this.

### spec.namespaceConfig.namespace

`string` · required

The Kubernetes namespace to install the add-on into. Must be a valid
RFC 1123 DNS label. Create-only in AWS.

- rule: {"required":true,"string":{"pattern":"^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$"}}

## Validation Rules

- `addon_version_format`: addon_version must be a full add-on version like 'v1.18.1-eksbuild.3' (aws eks describe-addon-versions lists them)
- `resolve_conflicts_on_create_valid`: resolve_conflicts_on_create must be 'NONE' or 'OVERWRITE' when set (PRESERVE is update-only in AWS)
- `resolve_conflicts_on_update_valid`: resolve_conflicts_on_update must be 'NONE', 'OVERWRITE', or 'PRESERVE' when set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEksAddon, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.addon_arn` | `string` | addon_arn is the Amazon Resource Name of the add-on -- arn:aws:eks:<region>:<account>:addon/<cluster>/<addon-name>/<uuid>. |
| `status.outputs.addon_name` | `string` | addon_name is the EKS catalog name the add-on was installed under (e.g. "vpc-cni"). |
| `status.outputs.addon_version` | `string` | addon_version is the version actually running -- the resolved AWS default when the spec left addon_version empty, else the pinned version. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.clusterName` | AwsEksCluster | `status.outputs.name` |
| `spec.serviceAccountRoleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.podIdentityAssociations[].roleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
