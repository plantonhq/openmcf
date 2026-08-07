# AWS EKS Add-on: Managed Cluster Software

## What This Component Is

An EKS managed add-on is cluster software AWS installs and operates
through the EKS control plane: the core networking trio (vpc-cni,
coredns, kube-proxy), storage drivers (EBS/EFS CSI), the Pod Identity
agent, observability agents, and marketplace software. `AwsEksAddon`
models one add-on on one cluster -- clusters run many, each with its
own version lifecycle, which is exactly why the add-on is a first-class
node rather than a list on the cluster.

AWS keys the add-on on (cluster, addon_name): one resource per pair.
The name and namespace config are create-time immutable; version,
configuration, conflict handling, and IAM wiring update in place.

## The Bootstrap Conflict, Modeled Honestly

A cluster created with `bootstrapSelfManagedAddons` enabled (the AWS
default) already runs self-managed copies of vpc-cni, coredns, and
kube-proxy. Installing the managed add-on over one of those fails with
`ConfigurationConflict` unless `resolveConflictsOnCreate: OVERWRITE`
adopts the existing install -- the exact migration path AWS built the
flag for.

The conflict enums are asymmetric in AWS -- create accepts
`NONE`/`OVERWRITE`, update additionally accepts `PRESERVE` (keep
out-of-band edits) -- and the spec enforces the split with two separate
CEL rules, so `PRESERVE` at create time fails validation instead of the
deploy.

## Two Identity Paths, Both by Reference

Add-ons that call AWS APIs (CSI drivers, observability agents) need an
IAM identity beyond the node role:

- **IRSA** (`serviceAccountRoleArn`): the classic path; requires an
  `AwsIamOidcProvider` for the cluster's `oidc_issuer_url` output, and
  the role trusts that provider.
- **EKS Pod Identity** (`podIdentityAssociations`): the modern path; the
  role trusts `pods.eks.amazonaws.com`, no per-cluster OIDC provider,
  one role usable across clusters. Requires the `eks-pod-identity-agent`
  add-on on the cluster.

Either way the role is a referenced `AwsIamRole` carrying its own
policies -- this component never modifies a role it merely references.

## Version Discipline

Empty `addonVersion` follows the AWS default for the cluster's
Kubernetes version -- manifests never go stale, and the resolved
version surfaces in the `addon_version` output. Pinning
(`v1.18.1-eksbuild.3` form, CEL-validated) freezes the add-on for
byte-identical clusters; bumping the pin rolls the add-on's own pods
(never the nodes).

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`tags` per add-on beyond the identity set** -- custom user tags are
  a platform-wide concern, not per-component scope.
- **Add-on catalog validation** -- `addonName` accepts any catalog
  name (including vendor-prefixed marketplace add-ons) rather than
  hardcoding today's AWS list; AWS validates availability against the
  cluster version at create time with a clear error.

## Failure Modes Worth Knowing

- A failed create (conflict without `OVERWRITE`) taints the resource;
  the retry deletes and re-creates the add-on, purging prior add-on
  configuration -- AWS's documented behavior.
- `coredns` reports DEGRADED on a cluster with zero nodes (its
  replicas cannot schedule); that is a scheduling truth, not an install
  failure. `kube-proxy` and `vpc-cni` (DaemonSets) reach ACTIVE on an
  empty cluster.
- IRSA without the OIDC provider fails at create with an AWS-side
  validation error; the provider requirement is on the spec comment.
