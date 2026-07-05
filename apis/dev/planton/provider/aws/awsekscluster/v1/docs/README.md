# AWS EKS Cluster: The Control Plane Node

## What This Component Is

An EKS cluster is AWS's managed Kubernetes control plane: the API server,
etcd, and the controllers that make a cluster a cluster. `AwsEksCluster`
models exactly that -- the control plane and its cluster-level posture --
and nothing else. Compute is a separate concern: explicit
`AwsEksNodeGroup` fleets reference this cluster by name, or EKS Auto Mode
(configured here, because it is a cluster capability) hands the data
plane to AWS entirely.

That boundary is the composition story. One cluster serves many node
pools with independent lifecycles; identity flows outward through the
exported OIDC issuer (`AwsIamOidcProvider` + `AwsIamRole` = IRSA); and
the cluster resource itself stays a stable, rarely-edited foundation
node.

## The Cluster Never Touches Its Role

The cluster role is a referenced `AwsIamRole`, and the reference is the
entire relationship: the module never attaches policies to a role it did
not create. `AmazonEKSClusterPolicy` belongs on the role itself
(`managedPolicyArns` -- one line). A role missing the policy fails
cluster creation at AWS with a clear error, which is the correct
behavior: the role's configuration is the role's responsibility. The same
principle applies to the node role on `AwsEksNodeGroup` and the Auto Mode
node role here.

## Endpoint Exposure Is a Pair, Not a Toggle

AWS models API-server reachability as two independent switches, and so
does the spec:

- `endpointPublicAccess` (AWS default true) -- internet reachability,
  scoped down with `publicAccessCidrs`.
- `endpointPrivateAccess` (AWS default false) -- in-VPC reachability
  through private ENIs.

Production clusters almost always want private access on; fully private
clusters turn public off AND private on -- disabling public alone would
leave the API server unreachable. The spec keeps both explicit so the
manifest reads as the posture it deploys.

## Auto Mode: One Toggle for AWS's Trio

EKS Auto Mode spans three AWS settings -- compute, block storage, and
elastic load balancing -- that the EKS API requires to be enabled or
disabled together. The spec expresses that constraint structurally:
`autoMode.enabled` is one switch, and both IaC modules expand it into the
three blocks AWS expects. A configuration where the trio disagrees is
unrepresentable rather than merely validated.

Auto Mode and managed node groups are alternative compute models. Both
are first-class here, but a cluster runs one or the other in practice:
Auto Mode provisions its own capacity, so `AwsEksNodeGroup` resources
attached to an Auto Mode cluster would fight it.

## One-Way Doors, Called Out

Several cluster decisions cannot be walked back, and the spec comments
say so where the decision is made:

- **Secrets encryption** (`kmsKeyArn`): enabling later works; disabling
  or re-keying a live cluster does not (AWS forces replacement).
- **`authenticationMode`**: `CONFIG_MAP` -> `API_AND_CONFIG_MAP` -> `API`
  is a one-way migration.
- **`controlPlaneEgressMode`**: reverting `CUSTOMER_ROUTED` to
  `AWS_MANAGED` replaces the cluster.
- **Create-only fields**: the name, the cluster role, `ipFamily`,
  `serviceIpv4Cidr`, `bootstrapSelfManagedAddons`, and
  `accessConfig.bootstrapClusterCreatorAdminPermissions`.

## Versioning Discipline

`version` accepts any 1.24+ minor and is deliberately future-proof (the
validation never needs relaxing as Kubernetes advances). Leaving it unset
takes AWS's current default. EKS upgrades one minor at a time and never
downgrades; node groups pin their own version, so the operational order
is control plane first, then roll each pool. `upgradeSupportType:
STANDARD` opts out of extended support's surcharge for teams that upgrade
on schedule.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`outpost_config`** -- EKS on AWS Outposts requires on-premises
  Outposts hardware; a niche, hardware-locked deployment model.
- **`remote_network_config`** -- hybrid nodes (on-premises machines
  joining the cluster) serve a small population; the CIDR plumbing is
  additive when demand appears.
- **`control_plane_scaling_config`** -- provisioned control-plane tiers
  (`tier-xl` and up) tune very-large-cluster API throughput; the standard
  tier serves the overwhelming majority.
- **`encryption_config.resources`** -- folded away: `secrets` is the only
  value the EKS API accepts, so the spec models the KMS key and nothing
  else.

## Immutability and Naming

The cluster name comes from `metadata.name` (AWS limit 100 characters,
truncated deterministically). Name and role are create-only; everything
else either updates in place or is a documented one-way door. Deletion
respects `deletionProtection` at the EKS API itself -- stronger than any
IaC-side guard.

## Dual-Engine Implementation

`AwsEksCluster` ships both a Terraform/OpenTofu module and a Pulumi (Go)
module at behavioral parity: the same endpoint-pair semantics
(proto-optional public access passes through as a real tri-state), the
same Auto Mode expansion, the same encryption one-way handling, and the
same seven outputs. One engine-level exception exists and is marked in
both modules: `control_plane_egress_mode` is not yet modeled by
pulumi-aws (v7.35.0), so only the Terraform module implements it until
the SDK catches up -- stack outputs are unaffected.
