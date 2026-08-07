# AWS EKS Access Entry: Cluster Access as a First-Class Grant

## What This Component Is

An access entry is EKS's modern answer to "who can reach this cluster's
Kubernetes API": a per-principal grant stored in EKS itself, replacing
the fragile aws-auth ConfigMap (one bad edit of which could lock a whole
team out). `AwsEksAccessEntry` models one grant -- one principal on one
cluster, which is AWS's own key for the resource.

The cluster must run with `accessConfig.authenticationMode: API` or
`API_AND_CONFIG_MAP`; CONFIG_MAP-only clusters reject entries at create
time.

## Two Authorization Paths

Authentication (the entry) and authorization (what the principal may
do) are separate, and the spec models both paths:

- **Access policies** (`policyAssociations`): AWS-managed policies --
  View/Edit/Admin/ClusterAdmin plus service-specific ones -- scoped
  `cluster`-wide or to `namespace` sets. No RBAC objects in-cluster;
  `aws eks list-access-policies` shows the catalog (custom policies do
  not exist).
- **RBAC groups** (`kubernetesGroups`): the entry maps the principal
  onto group names; you own the (Cluster)RoleBindings that give those
  groups power. Reserved `system:` groups are rejected at validation --
  AWS forbids them, and `system:masters` nostalgia belongs to the
  ClusterAdminPolicy instead.

Both can combine on one entry (policies for baseline, groups for
custom fine-grained rules).

## The Folded Associations

`aws_eks_access_policy_association` is a separate provider resource,
but it is pure per-principal glue -- the same (cluster, principal) pair
as its entry, referenced by nothing else, meaningless alone -- so
associations fold into the entry's spec. Both engines still materialize
each association as its own provider resource keyed by the policy name
(unique per entry: AWS allows one association per policy per
principal), so adding, re-scoping, or removing one diffs in place and
never touches its siblings.

## Node Types, Modeled Truthfully

`type` defaults to STANDARD -- the human/workload entry this kind
exists for. The node types (EC2, EC2_LINUX, EC2_WINDOWS, FARGATE_LINUX,
HYBRID_LINUX) register self-managed or hybrid nodes, and AWS forbids
groups/username/associations on them. The Terraform provider encodes
none of that -- it lets the API reject the deploy twenty minutes in --
so the spec enforces it in CEL instead. EKS auto-creates node-type
entries for managed node groups and Fargate profiles; model one only
for capacity EKS does not manage.

## Username Defaulting Is a Feature

Empty `userName` lets AWS template the username from the session
(`arn:...:role/X/{{SessionName}}` shape), which preserves the actual
session name in audit logs -- usually better than a fixed string.
Set it only when RBAC bindings need a predictable identity.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`tags` beyond the identity set** -- custom user tags are a
  platform-wide concern, not per-component scope.
- **Policy-ARN catalog validation** -- the spec checks the `arn:`
  shape, not membership in today's policy list; AWS validates the
  association at create time with a clear error, and the catalog grows
  faster than any hardcoded list would.
