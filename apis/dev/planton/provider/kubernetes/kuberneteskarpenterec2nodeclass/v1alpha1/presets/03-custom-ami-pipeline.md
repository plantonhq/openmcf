# Custom AMI Pipeline

This preset declares a machine template for organizations that build
their own node AMIs: the image is resolved from an SSM parameter your
pipeline maintains, so AMI rollout is controlled by the pipeline
publishing a new image id — not by an upstream "latest" alias. Kubelet
reservations and image-GC thresholds are tuned explicitly, and EC2
resources carry ownership tags.

## When to Use

- Organizations with a golden-image pipeline (hardening, agents,
  pre-pulled images) that publishes approved AMI ids to SSM
- Compliance regimes requiring every node image to pass an internal
  approval step before rollout

## Key Configuration Choices

- **`ssmParameter` selector term** — the custom-AMI-pipeline arm of AMI
  selection; when the parameter's value changes, nodes drift to the new
  AMI, making the pipeline the single control point for rollout
- **`amiFamily: AL2023`** — REQUIRED with ssm-parameter selection (the
  bootstrap/user-data family cannot be inferred without an alias); set
  it to the family your pipeline builds from
- **Kubelet tuning** — `maxPods: 110` overrides the ENI-based default;
  `kubeReserved`/`systemReserved` keep Kubernetes components and OS
  daemons from being starved by workload pods; image GC always runs
  above 85% disk usage and never below 80% (the high threshold must
  exceed the low one)
- **`tags`** — applied to the EC2 resources Karpenter creates
  (instances, launch templates, volumes) for ownership and cost
  attribution; note the `kubernetes.io/cluster/*` and `karpenter.sh/*`
  keys are controller-owned and rejected
- **Discovery-tag subnets/SGs and Karpenter-managed `role`** — same as
  the other presets

## Placeholders to Replace

| Placeholder                  | Description                                    | Where to Find                                 |
| ---------------------------- | ---------------------------------------------- | --------------------------------------------- |
| `<ssm-parameter-name>`       | SSM parameter holding the approved AMI id      | Your image pipeline / SSM Parameter Store     |
| `<eks-cluster-name>`         | Cluster name in the discovery tag value        | EKS console; tags on your VPC subnets and SGs |
| `<karpenter-node-role-name>` | IAM role name nodes assume                     | IAM console                                   |
| `<owning-team>`              | Team tag value for EC2 resources               | Your tagging policy                           |
| `<cost-center>`              | Cost-center tag value for EC2 resources        | Your tagging policy                           |

## Related Presets

- **01-al2023-standard** — EKS-optimized AL2023 without a custom pipeline
- **02-bottlerocket-encrypted** — hardened Bottlerocket with encrypted volumes
