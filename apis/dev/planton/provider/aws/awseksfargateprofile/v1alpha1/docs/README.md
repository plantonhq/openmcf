# AWS EKS Fargate Profile: Serverless Pod Placement

## What This Component Is

A Fargate profile is the routing rule that sends matching Kubernetes
pods to AWS Fargate -- per-pod serverless compute with no EC2 instances
behind it. `AwsEksFargateProfile` models one profile; clusters run up
to ten, each targeting different namespaces or label sets, which is why
the profile is a first-class node rather than a field on the cluster.

Everything attaches by reference: the cluster by its `name` output, the
pod execution role as a referenced `AwsIamRole`, and the subnets as
`AwsSubnet` nodes.

## Selector Semantics

A pod runs on Fargate when it matches ANY of the profile's selectors
(up to 5, the AWS quota, enforced in CEL). Within one selector the
namespace must match AND the pod must carry every listed label (up to
5 pairs). Namespaces and label values accept `*` and `?` wildcards, so
`team-*` covers a namespace-per-team layout with one selector.

Matching happens at admission: pods created before the profile existed
stay where they are until rescheduled.

## Everything Is Create-Only

AWS makes the entire profile immutable -- name, cluster, role, subnets,
selectors. Any spec change replaces the profile. Two operational
consequences, both documented on the spec:

- **Migrations need overlap**: create the replacement profile first (a
  second profile with the new shape), then remove the old one, so
  matching pods never lose a scheduling target.
- **AWS serializes profile operations** per cluster: one create or
  delete at a time. Parallel profile changes queue; both engines simply
  wait.

## The Networking Constraint

Fargate pods get ENIs in the profile's subnets, and AWS rejects any
subnet whose route table carries an internet-gateway route -- private
subnets only. Outbound traffic (image pulls from public registries!)
needs a NAT gateway or VPC endpoints. A profile in an isolated subnet
with no egress path will accept pods that then fail to pull images --
an AWS-side behavior worth knowing before the first deploy.

## The Role Carries Its Own Policy

The pod execution role needs a trust policy for
`eks-fargate-pods.amazonaws.com` and
`AmazonEKSFargatePodExecutionRolePolicy`. Both belong on the
`AwsIamRole` itself -- this module never modifies a role it merely
references. AWS retries around IAM eventual consistency at create time
("Misconfigured PodExecutionRole Trust Policy" resolves itself within
the create timeout when the role is genuinely correct).

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`tags` beyond the identity set** -- custom user tags are a
  platform-wide concern, not per-component scope.
- **Fargate pod sizing** -- CPU/memory come from the pod's own
  requests (Fargate rounds up to its size table); that is workload
  configuration, not infrastructure.
