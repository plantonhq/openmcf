# AWS EKS Node Group

Deploys a managed EKS node group: an EC2 worker fleet AWS provisions,
health-checks, and rolls, registered to an `AwsEksCluster` -- with
launch mechanics from inline knobs or a referenced `AwsLaunchTemplate`,
surge-enabled rollouts, and managed node auto-repair.

## What Gets Created

When you deploy an AwsEksNodeGroup resource, Planton provisions:

- **Managed node group** — an `aws_eks_node_group` / `eks.NodeGroup`
  named from `metadata.name`, registered to the referenced cluster,
  assuming the referenced node role, spread across the referenced
  subnets, with the scaling bounds, labels, taints, update and repair
  configuration you declare
- Behind the scenes AWS manages an **EC2 Auto Scaling group** for the
  fleet (surfaced as the `asg_name` output) and, when SSH is enabled
  without explicit source groups, a **remote-access security group**
  (surfaced as `remote_access_sg_id`)

The node role is never modified: attach `AmazonEKSWorkerNodePolicy`,
`AmazonEC2ContainerRegistryReadOnly`, and `AmazonEKS_CNI_Policy` on the
referenced `AwsIamRole` itself (`managedPolicyArns`).

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An EKS cluster** (`AwsEksCluster`) for the nodes to register with.
- **A node IAM role** (`AwsIamRole`) trusting `ec2.amazonaws.com` with the three worker policies attached.
- **Subnets** (`AwsSubnet`) — typically the cluster VPC's private subnets.
- **A launch template** (`AwsLaunchTemplate`) if the pool needs custom AMI/IMDSv2/storage mechanics.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksNodeGroup
metadata:
  name: general
spec:
  region: us-west-2
  clusterName:
    valueFrom:
      kind: AwsEksCluster
      name: platform
      fieldPath: status.outputs.name
  nodeRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: eks-node-role
      fieldPath: status.outputs.role_arn
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  instanceTypes: [m6i.large]
  amiType: AL2023_x86_64_STANDARD
  scaling:
    minSize: 2
    maxSize: 5
    desiredSize: 2
```

```shell
planton apply -f node-group.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the cluster's. | Required; non-empty |
| `clusterName` | `string \| valueFrom` | The cluster nodes register with. Defaults to referencing an `AwsEksCluster` `name` output. | Required |
| `nodeRoleArn` | `string \| valueFrom` | The IAM role nodes assume (must carry the three worker policies). Defaults to referencing an `AwsIamRole` `role_arn` output. | Required |
| `subnetIds` | `string[] \| valueFrom` | Subnets nodes launch into. One subnet is a legitimate zonal topology; use ≥2 AZs for fault tolerance. | Required; ≥1 entry |
| `scaling` | `object` | `minSize` (≥0) / `maxSize` (≥1) / `desiredSize` (≥0, within bounds). min=desired=0 is a dormant pool. | Required |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `launchTemplate.launchTemplateId` | `string \| valueFrom` | — | Launch from an `AwsLaunchTemplate`. Excludes `instanceTypes`, `diskSizeGb`, and `remoteAccess` (the template owns them). |
| `launchTemplate.version` | `string` | `$Default` | Numeric pin, `$Default`, or `$Latest`. Changing it rolls the pool. |
| `instanceTypes` | `string[]` | `t3.medium` | EC2 types AWS may launch; several types is the Spot best practice. Create-only. |
| `amiType` | `string` | inferred | EKS AMI family (`AL2023_*`, `BOTTLEROCKET_*`, `WINDOWS_*`, `AL2_*`, `CUSTOM`). Create-only. |
| `capacityType` | `enum` | `on_demand` | `on_demand`, `spot`, or `capacity_block`. Create-only. |
| `diskSizeGb` | `int` | 20 (Linux) / 50 (Windows) | Root volume size; 100 is a comfortable production default. Create-only. |
| `remoteAccess.ec2SshKey` | `string` | — | EC2 key pair for SSH. Without source groups, AWS opens port 22 to the internet — always scope it. Immutable. |
| `remoteAccess.sourceSecurityGroupIds` | `string[] \| valueFrom` | — | Security groups allowed to SSH. Immutable. |
| `labels` | `map` | `{}` | Kubernetes node labels. Updates in place. |
| `taints` | `object[]` | `[]` | Up to 50 taints (`key`/`value`/`effect`) for dedicated capacity. Updates in place. |
| `updateConfig` | `object` | 1 node | Rollout budget: exactly one of `maxUnavailable` / `maxUnavailablePercentage` (1–100), plus `updateStrategy: DEFAULT\|MINIMAL` (MINIMAL surges: launch-before-terminate). |
| `nodeRepairConfig` | `object` | off | Managed auto-repair: `enabled`, parallelism/threshold bounds (count XOR percentage), per-condition overrides. |
| `version` | `string` | cluster version | Kubernetes minor of the nodes; pin during control-plane upgrades, bump to roll. |
| `releaseVersion` | `string` | latest for `version` | Exact EKS-optimized AMI release, for byte-identical fleets. |
| `forceUpdateVersion` | `bool` | `false` | Force version updates past unsatisfiable pod disruption budgets. |

## Examples

### Dedicated GPU pool from a hardened launch template

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksNodeGroup
metadata:
  name: gpu
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  nodeRoleArn:
    valueFrom: { kind: AwsIamRole, name: eks-node-role, fieldPath: status.outputs.role_arn }
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  launchTemplate:
    launchTemplateId:
      valueFrom: { kind: AwsLaunchTemplate, name: gpu-nodes, fieldPath: status.outputs.launch_template_id }
    version: $Default
  amiType: AL2023_x86_64_NVIDIA
  scaling:
    minSize: 0
    maxSize: 4
    desiredSize: 1
  labels:
    pool: gpu
  taints:
    - key: nvidia.com/gpu
      value: "true"
      effect: NO_SCHEDULE
```

### Spot pool with surge rollouts and auto-repair

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksNodeGroup
metadata:
  name: batch-spot
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  nodeRoleArn:
    valueFrom: { kind: AwsIamRole, name: eks-node-role, fieldPath: status.outputs.role_arn }
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  instanceTypes: [m6i.large, m5.large, m6a.large]
  capacityType: spot
  scaling:
    minSize: 0
    maxSize: 10
    desiredSize: 3
  taints:
    - key: node-lifecycle
      value: spot
      effect: NO_SCHEDULE
  updateConfig:
    maxUnavailablePercentage: 25
    updateStrategy: MINIMAL
  nodeRepairConfig:
    enabled: true
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `nodegroup_name` | The pool's name |
| `nodegroup_arn` | The pool's ARN — what access entries and IAM policies reference |
| `asg_name` | The EC2 Auto Scaling group AWS manages behind the pool — the hook for ASG-level tooling |
| `remote_access_sg_id` | The SSH security group AWS creates when remote access lacks explicit source groups; empty otherwise |

## Related Components

- [AwsEksCluster](/docs/catalog/aws/awsekscluster) — the control plane this pool registers with
- [AwsLaunchTemplate](/docs/catalog/aws/awslaunchtemplate) — hardened launch mechanics for the pool
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — the node role (carries its own worker policies)
- [AwsSubnet](/docs/catalog/aws/awssubnet) — where the nodes launch
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — SSH source scoping for remote access
