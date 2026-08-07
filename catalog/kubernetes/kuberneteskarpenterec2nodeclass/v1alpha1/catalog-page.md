# Kubernetes Karpenter EC2 Node Class

Declares a Karpenter EC2NodeClass — the AWS-level machine template
NodePools launch instances from: which AMIs to boot, which subnets and
security groups to attach, which IAM identity nodes assume, disk layout,
kubelet configuration, and instance metadata posture. One NodeClass is
typically shared by several NodePools — the pools differ in constraints
and taints; the class is the common "how a node is built". The spec
holds 100% fidelity with the upstream `karpenter.k8s.aws/v1`
EC2NodeClass CRD, with the CRD's own validation rules mirrored so
mistakes surface at validate time.

## What Gets Created

- **EC2NodeClass** (cluster-scoped, named after `metadata.name`) — the
  `karpenter.k8s.aws/v1` custom resource carrying the AMI, subnet, and
  security-group selector terms, the node IAM identity, block device
  mappings, kubelet configuration, IMDS options, tags, and extra user
  data

## Prerequisites

- A Karpenter installation (KubernetesKarpenter) on an EKS or
  EKS-compatible cluster — the EC2NodeClass CRD does not exist before it
- A node IAM role (`role`, recommended — Karpenter manages the instance
  profile, needing `iam:PassRole` on the controller) or a pre-existing
  instance profile, registered in the cluster's access configuration so
  nodes can join
- Subnets and security groups selectable by the declared terms — tagging
  them `karpenter.sh/discovery: <cluster>` is the EKS convention

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKarpenterEc2NodeClass
metadata:
  name: default-al2023
spec:
  amiSelectorTerms:
    - alias: al2023@v20240807
  subnetSelectorTerms:
    - tags:
        karpenter.sh/discovery: my-eks-cluster
  securityGroupSelectorTerms:
    - tags:
        karpenter.sh/discovery: my-eks-cluster
  role: KarpenterNodeRole-my-eks-cluster
```

NodePools then reference the class through `node_class_ref`, and every
node they launch is built from this template. Pin AMI versions in
production — `latest` drifts nodes when a new AMI ships.

## Stack Outputs

| Output | Description |
|---|---|
| `node_class_name` | Name of the cluster-scoped EC2NodeClass — the value NodePools reference through their `node_class_ref` |

## Next Steps

Create KubernetesKarpenterNodePool resources referencing the class — a
general-purpose pool first, then dedicated spot/GPU pools sharing the
same template. Harden the template as it becomes load-bearing: an
encrypted gp3 root volume via `block_device_mappings`, kubelet reserves
and eviction thresholds via `kubelet`, and the IMDS defaults (IMDSv2
required, hop limit 1) left in place. For custom-AMI pipelines, switch
the AMI term to `ssm_parameter` with an explicit `ami_family` and let
the pipeline advance the parameter.
