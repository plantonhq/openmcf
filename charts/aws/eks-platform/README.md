# AWS EKS Platform

A production Kubernetes platform, not just a cluster: EKS with
access-entry authentication (the auditable successor to the aws-auth
ConfigMap), compute as a managed node group or EKS Auto Mode behind one
toggle, the core add-ons as first-class versioned resources, IRSA wiring
through an OIDC provider, and the IAM roles — every piece a referenceable
node in the graph rather than a side effect.

## Architecture

```
        AwsIamRole (cluster role: eks.amazonaws.com + AmazonEKSClusterPolicy)
                │
        AwsEksCluster (authenticationMode API, audit+authenticator logs,
                │      public/private endpoint dials)
                │
   ┌────────────┼───────────────────────────────┐
   │            │                               │
 auto_mode_enabled: false          auto_mode_enabled: true
   │            │                               │
 AwsEksNodeGroup│                    autoMode: general-purpose + system
 (AL2023, on-demand,                 pools — AWS manages compute, storage,
  25% surge updates)                 LB, networking, and DNS; NO add-on
   │            │                    resources are created
 AwsEksAddon ×4 │
 (vpc-cni, coredns, kube-proxy, pod-identity-agent)
                │
 AwsIamRole (node role: shared by both compute shapes)
                │
 AwsEksAccessEntry (cluster-admin for your IAM principal; toggle)
 AwsIamOidcProvider (IRSA: issuer URL ← cluster output; toggle)
```

The Auto Mode toggle deliberately changes MORE than the compute: add-ons
are only created in the node-group shape, because Auto Mode manages
networking and DNS itself and a second owner for those components would
fight it.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Cluster role | `AwsIamRole` | The control plane's identity (`AmazonEKSClusterPolicy`) |
| Node role | `AwsIamRole` | The worker identity, shared by both compute shapes |
| Cluster | `AwsEksCluster` | Control plane: API auth mode, endpoints, logs, Auto Mode arm |
| Node group | `AwsEksNodeGroup` | Managed workers — AL2023, on-demand, surge updates (classic shape) |
| Add-ons | `AwsEksAddon` ×4 | vpc-cni, coredns, kube-proxy, pod-identity-agent (classic shape) |
| Admin access entry | `AwsEksAccessEntry` | Cluster-admin for your IAM principal (conditional) |
| OIDC provider | `AwsIamOidcProvider` | The IRSA foundation (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for the cluster and companions | `us-east-1` | string |
| `cluster_name` | Name prefix and the cluster's own name | `platform` | string |
| `kubernetes_version` | Pinned Kubernetes minor (upgrades = bump + redeploy) | `1.31` | string |
| `subnet_ids` | Private subnets across two-plus AZs | placeholders | list |
| `endpoint_public_access` | Public (IAM-authenticated) API endpoint | `true` | bool |
| `endpoint_private_access` | In-VPC API endpoint (keep on) | `true` | bool |
| `auto_mode_enabled` | EKS Auto Mode instead of node group + add-ons | `false` | bool |
| `node_instance_types` | Node group instance types, preference-ordered | `m6i.large` | list |
| `node_min_size` / `node_max_size` / `node_desired_size` | Node scaling bounds and seed | `2` / `6` / `2` | number |
| `node_disk_size_gb` | Root volume per node | `100` | number |
| `admin_access_entry_enabled` | Cluster-admin entry for your principal | `true` | bool |
| `admin_principal_arn` | The IAM role/user to grant admin | placeholder | string |
| `oidc_provider_enabled` | Register the OIDC issuer for IRSA | `true` | bool |

## Choosing the compute shape

- **Node group (default)**: you pick instance types and scaling bounds,
  you own upgrades (surge-rolled at 25%), you pay plain EC2 prices. The
  right default for steady fleets and teams that tune their nodes.
- **Auto Mode** (`auto_mode_enabled: true`): AWS launches and retires
  nodes per pod demand, manages storage/LB/networking/DNS, patches and
  rotates nodes itself — zero node operations, for a management premium on
  the EC2 rate. The right shape when operating nodes is not where you want
  to spend people. Flipping an EXISTING cluster between shapes is a
  workload migration, not a toggle flip: bring the new compute up, drain,
  then remove the old.

## After deploying

1. **Connect**: `aws eks update-kubeconfig --name <cluster_name> --region
   <region>` with the `admin_principal_arn` identity, then `kubectl get
   nodes`.
2. The useful join points:
   - `AwsEksCluster` → `status.outputs.name` (every satellite's cluster
     argument), `status.outputs.endpoint`,
     `status.outputs.oidc_issuer_url`
   - `AwsIamOidcProvider` → `status.outputs.provider_arn` (workload role
     trust policies) and `status.outputs.provider_url` (their conditions)
   - `AwsIamRole` (node role) → keep it minimal; workload permissions go
     through IRSA or Pod Identity, never here

## Giving a workload AWS permissions (the IRSA loop)

One `AwsIamRole` per workload, trusting the provider for exactly one
ServiceAccount:

```yaml
trustPolicy:
  Version: "2012-10-17"
  Statement:
    - Effect: Allow
      Principal:
        Federated: <AwsIamOidcProvider status.outputs.provider_arn>
      Action: sts:AssumeRoleWithWebIdentity
      Condition:
        StringEquals:
          <provider_url>:sub: system:serviceaccount:<namespace>:<serviceaccount>
          <provider_url>:aud: sts.amazonaws.com
```

Annotate the ServiceAccount with the role's ARN
(`eks.amazonaws.com/role-arn`) and the pod's SDK calls carry that role —
no node-wide grants, no static keys. The pod-identity-agent add-on ships
too, so the newer Pod Identity path (associations on an `AwsEksAddon` or
via the console) is equally available per workload.

## Day-2 guidance

- **More teams, narrower grants**: add `AwsEksAccessEntry` resources per
  team role with `AmazonEKSEditPolicy`/`AmazonEKSViewPolicy`, and
  namespace-scope them with `accessScope: {type: namespace, namespaces:
  [...]}`.
- **Upgrades**: bump `kubernetes_version` and redeploy — control plane
  first, then the node group surge-rolls; add-ons resolve their matching
  default versions. One minor at a time, as Kubernetes requires.
- **Spot capacity**: add a second `AwsEksNodeGroup` with `capacityType:
  spot` and several instance types for interruption-tolerant workloads,
  rather than mixing spot into the baseline group.
- **Cluster autoscaling**: install Karpenter or cluster-autoscaler via
  IRSA (the loop above); the scaling bounds here then become its limits.
- **Stateful workloads**: add the `aws-ebs-csi-driver` add-on as another
  `AwsEksAddon` with a Pod Identity association or IRSA role — the same
  first-class shape as the core four.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
