# AliCloud Kubernetes Cluster

Deploys an ACK Managed Kubernetes cluster on Alibaba Cloud with a fully managed control plane, configurable Flannel or Terway CNI networking, optional RRSA for pod-level IAM, KMS secrets encryption, control plane logging, maintenance windows, and automatic version upgrades. Worker nodes are managed separately through AliCloudKubernetesNodePool. The cluster integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches, security groups, KMS keys, and SLS log projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ACK Managed Kubernetes Cluster** -- an `alicloud_cs_managed_kubernetes` with the selected specification tier, Kubernetes version, CNI configuration, and VSwitches for multi-AZ control plane placement
- **Cluster Addons** -- one addon per entry in the `addons` list (e.g., `flannel` or `terway-eniip` for networking, `csi-plugin` for storage, `logtail-ds` for logging); addons are only configurable at creation time
- **Control Plane Logging** -- created only when `logging` is configured; sends control plane component logs and optional audit logs to an SLS project
- **Maintenance Window** -- created only when `maintenanceWindow` is configured; restricts ACK updates and patches to the specified time window
- **Auto-Upgrade Policy** -- created only when `autoUpgrade.enabled` is `true` and a maintenance window is configured; automatically upgrades the cluster version based on the selected channel
- **NAT Gateway** -- created only when `newNatGateway` is `true` (default); provides internet access for cluster nodes
- **Internet-Facing SLB** -- created only when `slbInternetEnabled` is `true` (default); exposes the Kubernetes API server publicly
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **1-5 VSwitches** in distinct availability zones within a VPC. These determine the control plane and default worker node network placement. Provide VSwitch IDs directly or reference AliCloudVswitch Cloud Resources via ValueFromRef.
- **Dedicated pod VSwitches** (optional) -- required when using Terway CNI (`terway-eniip` addon). These VSwitches should be in the same AZs as the node VSwitches but use dedicated CIDR ranges to avoid IP exhaustion.
- **A security group** (optional) -- provide an existing security group ID or let ACK auto-create one. Set `isEnterpriseSecurityGroup: true` for advanced security groups supporting up to 65,536 rules.
- **A KMS key** (optional) -- for encrypting Kubernetes Secrets at rest via envelope encryption. Provide the key ID directly or reference an AliCloudKmsKey Cloud Resource via ValueFromRef. Immutable after creation.
- **An SLS project** (optional) -- for control plane and audit logging. If omitted, ACK auto-creates a project named `k8s-log-{cluster-id}`. Provide the project name directly or reference an AliCloudLogProject Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **AliCloud Kubernetes Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Terway** preset in the [Presets](#presets) tab to pre-populate a production-ready configuration with ENI-based pod networking.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudKubernetesCluster
metadata:
  name: platform-cluster
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vswitchIds:
    - value: "vsw-abc001"
    - value: "vsw-abc002"
  serviceCidr: "172.21.0.0/20"
  addons:
    - name: flannel
    - name: csi-plugin
    - name: csi-provisioner
  deletionProtection: true
```

```shell
planton apply -f ack-cluster.yaml
```

This creates an `ack.standard` Managed Kubernetes cluster with Flannel CNI, IPVS proxy mode, a /24 node CIDR mask, auto-created NAT gateway, and a public API server endpoint. RRSA, secrets encryption, control plane logging, and maintenance windows are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to VSwitches, a security group, and a KMS key deployed in the same InfraPipeline:

```yaml
spec:
  vswitchIds:
    - valueFrom:
        kind: AliCloudVswitch
        name: node-vswitch-a
        fieldPath: status.outputs.vswitch_id
    - valueFrom:
        kind: AliCloudVswitch
        name: node-vswitch-b
        fieldPath: status.outputs.vswitch_id
  securityGroupId:
    valueFrom:
      kind: AliCloudSecurityGroup
      name: cluster-sg
      fieldPath: status.outputs.security_group_id
  encryptionProviderKey:
    valueFrom:
      kind: AliCloudKmsKey
      name: secrets-key
      fieldPath: status.outputs.key_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitches, security group, and KMS key first, then provisions the ACK cluster with the resolved values.

## Key Configuration

These are the most important decisions when configuring an ACK cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CNI plugin** -- Choose between Flannel (overlay networking, set `podCidr`) and Terway (ENI-based, set `podVswitchIds`). Flannel is simpler and works with any VPC configuration. Terway assigns VPC-native IPs to pods for better performance and direct security group attachment, but requires dedicated pod VSwitches with sufficient IP capacity. This choice is immutable after creation. Add the corresponding addon (`flannel` or `terway-eniip`) in the `addons` list.

**Cluster specification** -- Set `clusterSpec` to `ack.standard` (free, basic SLA) or `ack.pro.small` (paid, enhanced SLA with managed node pools, topology-aware scheduling, and additional observability). Supports in-place upgrade from standard to professional.

**RRSA for pod-level IAM** -- Set `enableRrsa: true` to enable RAM Roles for Service Accounts via OIDC federation. This allows Kubernetes service accounts to assume RAM roles without static access keys in pods. Requires Kubernetes 1.22.3+. Once enabled, RRSA cannot be disabled.

**Addons** -- Addons are only configurable at cluster creation time. Common addons: `flannel` or `terway-eniip` (networking), `csi-plugin` and `csi-provisioner` (storage), `logtail-ds` (logging to SLS), `nginx-ingress-controller` or `alb-ingress-controller` (ingress), `arms-prometheus` (monitoring), `metrics-server` (autoscaling). Post-creation addon management requires separate resources.

**API server access** -- Set `slbInternetEnabled: false` to restrict the API server to VPC-internal access only. When public, use `customSan` to add custom Subject Alternative Names to the API server TLS certificate for additional domain names or IP addresses.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchIds` | `status.outputs.vswitch_id` |
| **AliCloudVswitch** (optional) | `podVswitchIds` | `status.outputs.vswitch_id` |
| **AliCloudSecurityGroup** (optional) | `securityGroupId` | `status.outputs.security_group_id` |
| **AliCloudKmsKey** (optional) | `encryptionProviderKey` | `status.outputs.key_id` |
| **AliCloudLogProject** (optional) | `logging.controlPlaneLogProject` | `status.outputs.project_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | ACK cluster ID | Node pool attachment, monitoring dashboards |
| `cluster_name` | Cluster name as created | kubectl context naming, audit references |
| `api_server_internet` | Public API server endpoint | kubectl configuration, CI/CD pipeline access |
| `api_server_intranet` | Private (VPC-internal) API server endpoint | VPC-internal kubectl access, private CI/CD runners |
| `vpc_id` | VPC ID computed from the provided VSwitches | Network auditing, related resource lookups |
| `security_group_id` | Security group used by cluster nodes | Additional security group rule configuration |
| `nat_gateway_id` | NAT gateway auto-created by the cluster | Network auditing (empty when `newNatGateway` is false) |
| `worker_ram_role_name` | RAM role name for worker nodes | Granting worker nodes access to ACR, SLS, or other services |
| `rrsa_oidc_issuer_url` | RRSA OIDC issuer URL | RAM role trust policies for pod-level IAM (empty when RRSA disabled) |
| `ram_oidc_provider_name` | RRSA OIDC provider name in RAM | OIDC federation configuration (empty when RRSA disabled) |
| `ram_oidc_provider_arn` | RRSA OIDC provider ARN | IAM policy conditions for pod identity (empty when RRSA disabled) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production with Terway** -- An `ack.pro.small` cluster with Terway ENI-based networking for VPC-native pod IPs, dedicated pod VSwitches, RRSA enabled, secrets encryption, and deletion protection. Start from the **Production Terway** preset.

**Development with Flannel** -- An `ack.standard` cluster with Flannel overlay networking, minimal configuration, and PostPaid billing. Suitable for development and testing environments. Start from the **Development Flannel** preset.

**Production with Flannel** -- An `ack.pro.small` cluster with Flannel networking, RRSA enabled, control plane logging, maintenance window, and deletion protection. Suitable for production workloads that do not require VPC-native pod IPs. Start from the **Production Flannel** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitches for control plane and worker node network placement
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- provides network access control for cluster nodes
- [**AliCloud KMS Key**](/cloud-catalog/ali-cloud-kms-key) -- provides a customer-managed key for Kubernetes secrets envelope encryption
- [**AliCloud Log Project**](/cloud-catalog/ali-cloud-log-project) -- provides the SLS project for control plane and audit logging