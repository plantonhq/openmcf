---
title: "EKS Autodiscovery"
description: "This preset installs the Cluster Autoscaler on EKS in the recommended posture: tag-based ASG auto-discovery, keyless AWS access via IRSA, the `least-waste` expander, and balanced multi-AZ node..."
type: "preset"
rank: "01"
presetSlug: "01-eks-autodiscovery"
componentSlug: "cluster-autoscaler"
componentTitle: "Cluster Autoscaler"
provider: "kubernetes"
icon: "package"
order: 1
---

# EKS Autodiscovery

This preset installs the Cluster Autoscaler on EKS in the recommended
posture: tag-based ASG auto-discovery, keyless AWS access via IRSA, the
`least-waste` expander, and balanced multi-AZ node groups. One
installation per cluster — the autoscaler leader-elects and owns the
cluster-wide scaling decision.

Note: for AWS clusters that want right-sized machines launched on demand
instead of pre-defined groups, `KubernetesKarpenter` is the modern
alternative; this component earns its keep where capacity is organized
as pre-defined Auto Scaling groups.

## When to Use

- EKS clusters whose node capacity is organized as EC2 Auto Scaling
  groups (managed or self-managed node groups)
- The 30-second choice for an autoscaler on EKS

## Key Configuration Choices

- **Tag-based auto-discovery** (`aws.autoDiscovery.clusterName`) — the
  recommended mode: every ASG tagged
  `k8s.io/cluster-autoscaler/enabled` +
  `k8s.io/cluster-autoscaler/<cluster-name>` is managed, so new node
  groups enroll by tagging alone (no manifest change)
- **IRSA (`aws.irsaRoleArn`)** — the autoscaler calls the Auto Scaling
  APIs without stored keys; preferred over static access keys
- **`expander: least-waste`** — the common production choice for picking
  among node groups (upstream's default is `random`)
- **`balanceSimilarNodeGroups: true`** — treats node groups with
  identical instance types and labels as one and keeps their sizes
  balanced, the pattern for one-ASG-per-zone multi-AZ setups
- **`kube-system` + `createNamespace: false`** — the upstream
  convention; the namespace already exists
- **`chartVersion: "9.59.0"`** (spec default) — keep the autoscaler's
  minor version aligned with the cluster's Kubernetes minor per
  upstream guidance

## Placeholders to Replace

| Placeholder                                         | Description                                 | Where to Find                          |
| --------------------------------------------------- | ------------------------------------------- | -------------------------------------- |
| `<aws-region>`                                      | AWS region of the Auto Scaling groups       | Your AWS account                       |
| `<eks-cluster-name>`                                | Cluster name carried by the discovery tags  | EKS console or `AwsEksCluster` outputs |
| `arn:aws:iam::123456789012:role/cluster-autoscaler` | IRSA role ARN — replace account id and name | IAM console                            |

## Related Presets

- **02-cluster-api** — Cluster API MachineDeployments (self-managed /
  multi-cloud clusters)
- **03-azure-vmss** — Azure VM scale sets with workload identity
