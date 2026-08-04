# AwsEksAccessEntry

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsEksAccessEntrySpec grants one IAM principal (a role or user) access
to an AwsEksCluster's Kubernetes API through EKS access entries -- the
modern access model that replaces hand-editing the aws-auth ConfigMap.
The cluster must have API authentication enabled
(access_config.authentication_mode "API" or "API_AND_CONFIG_MAP").

Authorization comes from either side, or both:
- policy_associations attach AWS-managed EKS access policies
  (AmazonEKSViewPolicy, ...AdminPolicy, ...ClusterAdminPolicy, ...),
  scoped to the whole cluster or to namespaces -- no Kubernetes RBAC
  objects needed.
- kubernetes_groups names groups your OWN RBAC bindings reference --
  the entry maps the principal onto them; you own the (Cluster)Role
  and bindings inside the cluster.

AWS keys the entry on (cluster, principal): one entry per principal
per cluster, and both keys are create-time immutable. Groups,
username, and policy associations update in place. EKS auto-creates
entries for the cluster creator and for managed node group / Fargate
roles -- model entries only for the principals you grant yourself.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksAccessEntry
metadata:
  name: awseksaccessentry-demo
spec:
  region: us-west-2
  clusterName:
    value: awsekscluster-demo
  principalArn:
    value: arn:aws:iam::123456789012:role/TeamViewerRole
  kubernetesGroups:
    - platform-viewers
  policyAssociations:
    - policyArn: arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy
      accessScope:
        type: cluster
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.clusterName` | `string \| valueFrom` | yes |  | AwsEksCluster (`status.outputs.name`) |
| `spec.principalArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.type` | `string` |  |  |  |
| `spec.kubernetesGroups` | `[]string` |  |  |  |
| `spec.userName` | `string` |  |  |  |
| `spec.policyAssociations` | `[]AwsEksAccessEntryPolicyAssociation` |  |  |  |
| `spec.policyAssociations[].policyArn` | `string` | yes |  |  |
| `spec.policyAssociations[].accessScope` | `AwsEksAccessEntryAccessScope` | yes |  |  |
| `spec.policyAssociations[].accessScope.type` | `string` | yes |  |  |
| `spec.policyAssociations[].accessScope.namespaces` | `[]string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the entry's cluster lives in. Must match the
cluster's region. Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.clusterName

`string | valueFrom` · required

The EKS cluster the principal gets access to. Reference an
AwsEksCluster's name output or pass a literal cluster name for a
cluster managed outside Planton. Create-only in AWS.

- references: AwsEksCluster (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEksCluster, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.principalArn

`string | valueFrom` · required

The IAM principal being granted access. Reference an AwsIamRole's
role_arn output, or pass a literal role or user ARN (IAM users work
as literals -- roles are the norm for team and workload access).
Create-only in AWS: one entry per principal per cluster.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.type

`string`

The entry type. Empty or "STANDARD" is the human/workload entry this
kind exists for. The node types -- "EC2", "EC2_LINUX", "EC2_WINDOWS",
"FARGATE_LINUX", "HYBRID_LINUX" -- exist for infrastructure roles
(EKS creates them automatically for managed node groups and Fargate
profiles; set one only when registering self-managed or hybrid
nodes). AWS forbids groups, username, and policy associations on
non-STANDARD entries (enforced below). Create-only in AWS.

### spec.kubernetesGroups

`[]string`

Kubernetes groups the principal is mapped onto -- your own RBAC
(Cluster)RoleBindings reference these group names; nothing is
created in-cluster for you. Groups may not start with the reserved
"system:" prefix. STANDARD entries only. Updates in place.

- rule: {"repeated":{"items":{"string":{"minLen":"1","maxLen":"253","notContains":" "}}}}

### spec.userName

`string`

The Kubernetes username the principal authenticates as, visible in
audit logs and usable in RBAC bindings. Empty lets AWS default it
(the principal ARN for users; a session-templated name for roles --
the right choice for audit trails, since it preserves the session
name). STANDARD entries only. Updates in place.

### spec.policyAssociations

`[]AwsEksAccessEntryPolicyAssociation`

AWS-managed EKS access policies attached to the principal,
cluster-scoped or namespace-scoped -- authorization without any
in-cluster RBAC objects. STANDARD entries only. Associations add,
change scope, and remove in place.

### spec.policyAssociations[].policyArn

`string` · required

The access policy's ARN. These are AWS-managed and account-less --
"arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy",
...EditPolicy, ...AdminPolicy (namespace-level admin),
...ClusterAdminPolicy (full cluster admin), plus service-specific
policies. `aws eks list-access-policies` shows the catalog; custom
policies do not exist.

- rule: {"required":true,"string":{"prefix":"arn:"}}

### spec.policyAssociations[].accessScope

`AwsEksAccessEntryAccessScope` · required

What the policy's permissions apply to: the whole cluster or a set
of namespaces.

- rule: {"required":true}
- rule: access scope type must be 'cluster' or 'namespace'
- rule: namespace-scoped associations must list at least one namespace
- rule: cluster-scoped associations must not list namespaces

### spec.policyAssociations[].accessScope.type

`string` · required

"cluster" applies the policy cluster-wide; "namespace" restricts it
to the listed namespaces.

- rule: {"required":true}

### spec.policyAssociations[].accessScope.namespaces

`[]string`

The namespaces the policy applies to when type is "namespace" --
e.g. give a team AmazonEKSAdminPolicy inside its own namespaces
only. Trailing "*" wildcards are allowed ("team-*").

- rule: {"repeated":{"items":{"string":{"minLen":"1","maxLen":"63"}}}}

## Validation Rules

- `type_valid`: type must be 'STANDARD', 'EC2', 'EC2_LINUX', 'EC2_WINDOWS', 'FARGATE_LINUX', or 'HYBRID_LINUX' when set
- `non_standard_type_excludes_access_fields`: kubernetes_groups, user_name, and policy_associations are only valid on STANDARD access entries
- `kubernetes_groups_not_reserved`: kubernetes_groups may not use the reserved 'system:' prefix

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEksAccessEntry, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.accessEntryArn` | `string` | access_entry_arn is the Amazon Resource Name of the entry -- arn:aws:eks:<region>:<account>:access-entry/<cluster>/<principal type>/<account>/<principal name>/<uuid>. |
| `status.outputs.principalArn` | `string` | principal_arn is the IAM principal the entry grants access to, as resolved at provisioning time. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.clusterName` | AwsEksCluster | `status.outputs.name` |
| `spec.principalArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
