# AliCloudKubernetesCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudKubernetesClusterSpec defines the configuration for an Alibaba Cloud
ACK Managed Kubernetes cluster.

ACK (Alibaba Cloud Container Service for Kubernetes) provides a fully managed
control plane. Worker nodes are managed separately through node pools
(AliCloudKubernetesNodePool), which have their own lifecycle.

This component wraps a single provider resource:
  Terraform: alicloud_cs_managed_kubernetes
  Pulumi:    cs.ManagedKubernetes

Networking supports two CNI plugins:
  - Flannel (overlay): set pod_cidr and add the "flannel" addon.
  - Terway (ENI-based): set pod_vswitch_ids and add the "terway-eniip" addon.

Addons are only configurable at cluster creation time. Post-creation addon
management requires alicloud_cs_kubernetes_addon (Terraform) or equivalent.

Billing: The cluster control plane itself is free for ack.standard; only
ack.pro.small incurs a management fee. Worker node costs are determined by
the node pool instance types.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudKubernetesCluster
metadata:
  name: alicloudkubernetescluster-demo
spec:
  region: cn-hangzhou
  vswitchIds:
    - value: vsw-aaa111
    - value: vsw-bbb222
  serviceCidr: "172.21.0.0/20"
  podCidr: "172.20.0.0/16"
  addons:
    - name: flannel
    - name: csi-plugin
    - name: csi-provisioner
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.name` | `string` |  |  |  |
| `spec.version` | `string` |  |  |  |
| `spec.clusterSpec` | `string` |  | `ack.standard` |  |
| `spec.clusterDomain` | `string` |  |  |  |
| `spec.vswitchIds` | `[]string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.podCidr` | `string` |  |  |  |
| `spec.podVswitchIds` | `[]string \| valueFrom` |  |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.serviceCidr` | `string` | yes |  |  |
| `spec.proxyMode` | `string` |  | `ipvs` |  |
| `spec.nodeCidrMask` | `int32` |  | `24` |  |
| `spec.newNatGateway` | `bool` |  | `true` |  |
| `spec.slbInternetEnabled` | `bool` |  | `true` |  |
| `spec.securityGroupId` | `string \| valueFrom` |  |  | AliCloudSecurityGroup (`status.outputs.security_group_id`) |
| `spec.isEnterpriseSecurityGroup` | `bool` |  |  |  |
| `spec.enableRrsa` | `bool` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.encryptionProviderKey` | `string \| valueFrom` |  |  | AliCloudKmsKey (`status.outputs.key_id`) |
| `spec.customSan` | `string` |  |  |  |
| `spec.addons` | `[]AliCloudKubernetesAddon` |  |  |  |
| `spec.addons[].name` | `string` | yes |  |  |
| `spec.addons[].config` | `string` |  |  |  |
| `spec.addons[].version` | `string` |  |  |  |
| `spec.addons[].disabled` | `bool` |  |  |  |
| `spec.logging` | `AliCloudKubernetesClusterLogging` |  |  |  |
| `spec.logging.controlPlaneLogProject` | `string \| valueFrom` |  |  | AliCloudLogProject (`status.outputs.project_name`) |
| `spec.logging.controlPlaneLogTtl` | `string` |  | `30` |  |
| `spec.logging.controlPlaneLogComponents` | `[]string` |  |  |  |
| `spec.logging.auditLogEnabled` | `bool` |  |  |  |
| `spec.logging.auditLogSlsProject` | `string` |  |  |  |
| `spec.maintenanceWindow` | `AliCloudKubernetesClusterMaintenanceWindow` |  |  |  |
| `spec.maintenanceWindow.enable` | `bool` |  |  |  |
| `spec.maintenanceWindow.maintenanceTime` | `string` |  |  |  |
| `spec.maintenanceWindow.duration` | `string` |  |  |  |
| `spec.maintenanceWindow.weeklyPeriod` | `string` |  |  |  |
| `spec.autoUpgrade` | `AliCloudKubernetesClusterAutoUpgrade` |  |  |  |
| `spec.autoUpgrade.enabled` | `bool` |  |  |  |
| `spec.autoUpgrade.channel` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.timezone` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for the cluster.
Must match the region of the VSwitches in vswitch_ids.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.name

`string`

Cluster name. 1-63 characters; must start with a letter or digit.
If omitted, defaults to metadata.name in the IaC modules.

- rule: name must be between 1 and 63 characters when set

### spec.version

`string`

Kubernetes version (e.g., "1.28", "1.30").
If omitted, the provider uses the latest stable version.

### spec.clusterSpec

`string` · optional (explicit presence)

Cluster specification that determines the control plane SLA and features.
"ack.standard" -- free managed control plane with basic SLA.
"ack.pro.small" -- professional managed cluster with enhanced SLA,
  managed node pools, topology-aware scheduling, and more.
Supports in-place upgrade from ack.standard to ack.pro.small.
Default: "ack.standard"

- default: `ack.standard`
- rule: cluster_spec must be one of: ack.standard, ack.pro.small

### spec.clusterDomain

`string`

Cluster-internal domain name used for Kubernetes service discovery.
Default: "cluster.local"
Immutable after creation.

### spec.vswitchIds

`[]string | valueFrom` · required

VSwitch IDs for control plane placement. 1-5 VSwitches in distinct
availability zones for high availability.
These VSwitches also serve as the default worker node VSwitches.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.podCidr

`string`

Pod network CIDR block for Flannel CNI.
Required when using the "flannel" network addon; ignored for Terway.
Must not overlap the VPC CIDR, service_cidr, or node CIDR.
Immutable after creation.
Example: "172.20.0.0/16"

### spec.podVswitchIds

`[]string | valueFrom`

VSwitch IDs for pod ENI allocation when using Terway CNI.
Required when using the "terway-eniip" network addon; ignored for Flannel.
Should be in the same AZs as vswitch_ids but use dedicated CIDR ranges
to avoid IP exhaustion on the node VSwitches.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.serviceCidr

`string` · required

Service network CIDR block for Kubernetes ClusterIP services.
Must not overlap the VPC CIDR, pod CIDR, or node CIDR.
Immutable after creation.
Example: "172.21.0.0/20"

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.proxyMode

`string` · optional (explicit presence)

kube-proxy mode for service routing.
"ipvs" (default) -- IPVS-based load balancing with O(1) connection processing.
"iptables" -- legacy iptables-based routing.
Immutable after creation.
Default: "ipvs"

- default: `ipvs`
- rule: proxy_mode must be one of: iptables, ipvs

### spec.nodeCidrMask

`int32` · optional (explicit presence)

Node CIDR mask that controls how many pods can run on each node.
A /24 mask gives 256 addresses per node (~253 pods); /26 gives 64.
Range: 24-28. Default: 24.
Immutable after creation.

- default: `24`
- rule: {"int32":{"lte":28,"gte":24}}

### spec.newNatGateway

`bool` · optional (explicit presence)

Whether ACK should automatically create a NAT gateway for the cluster
VPC. Set to false when you manage your own AliCloudNatGateway component
to avoid creating a duplicate NAT gateway.
Default: true

- default: `true`

### spec.slbInternetEnabled

`bool` · optional (explicit presence)

Whether to create an internet-facing SLB for the Kubernetes API server.
When true, the API server is accessible from the public internet.
When false, only the VPC-internal endpoint is available.
Default: true

- default: `true`

### spec.securityGroupId

`string | valueFrom`

Security group ID for the cluster nodes. If omitted, ACK auto-creates
a security group. Conflicts with is_enterprise_security_group.

- references: AliCloudSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.isEnterpriseSecurityGroup

`bool` · optional (explicit presence)

Whether ACK should auto-create an advanced (enterprise) security group.
Advanced security groups support up to 65,536 rules and 100,000 ENIs.
Conflicts with security_group_id.
Default: false

### spec.enableRrsa

`bool` · optional (explicit presence)

Enable RRSA (RAM Roles for Service Accounts) for pod-level IAM.
When enabled, Kubernetes service accounts can assume RAM roles via OIDC
federation, eliminating the need for static access keys in pods.
Requires Kubernetes 1.22.3+.
WARNING: Once enabled, RRSA cannot be disabled.
Default: false

### spec.deletionProtection

`bool` · optional (explicit presence)

Enable cluster deletion protection. When true, the cluster cannot be
deleted through the API until deletion_protection is explicitly disabled.

### spec.encryptionProviderKey

`string | valueFrom`

KMS key ID for encrypting Kubernetes Secrets at rest.
When set, all Secrets stored in etcd are encrypted with this key.
Immutable after creation.

- references: AliCloudKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.customSan

`string`

Custom Subject Alternative Names for the API server TLS certificate.
Comma-separated list of IP addresses or domain names that the API server
certificate should be valid for, in addition to the default endpoints.
Example: "10.0.0.1,api.example.com"

### spec.addons

`[]AliCloudKubernetesAddon`

Cluster addons to install during creation.

Common addons:
  Network: "flannel", "terway-eniip" (choose one)
  Storage: "csi-plugin", "csi-provisioner"
  Logging: "logtail-ds" (config: {"IngressDashboardEnabled":"true","sls_project_name":"..."})
  Ingress: "nginx-ingress-controller", "alb-ingress-controller"
  Monitoring: "arms-prometheus", "ack-node-problem-detector"
  DNS: "managed-coredns"
  Autoscaling: "metrics-server"

Addons are only configurable at creation time. Post-creation management
requires the alicloud_cs_kubernetes_addon resource.

### spec.addons[].name

`string` · required

Addon name. Must match an ACK addon identifier.
Examples: "flannel", "terway-eniip", "csi-plugin", "logtail-ds",
  "nginx-ingress-controller", "arms-prometheus", "metrics-server".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.addons[].config

`string`

JSON-encoded addon configuration.
Structure is addon-specific. Example for logtail-ds:
  {"IngressDashboardEnabled":"true","sls_project_name":"my-log-project"}

### spec.addons[].version

`string`

Addon version. If omitted, ACK installs the default version for the
cluster's Kubernetes version.

### spec.addons[].disabled

`bool`

Disable automatic installation of this addon.
When true, the addon is registered but not installed.

### spec.logging

`AliCloudKubernetesClusterLogging`

Control plane and audit logging configuration.
When omitted, no control plane logs or audit logs are collected.

### spec.logging.controlPlaneLogProject

`string | valueFrom`

SLS (Log Service) project for control plane component logs.
If omitted, ACK auto-creates a project named "k8s-log-{cluster-id}".

- references: AliCloudLogProject (`status.outputs.project_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudLogProject, name: <that resource's name>, fieldPath: status.outputs.project_name}} -- a bare string does not parse

### spec.logging.controlPlaneLogTtl

`string` · optional (explicit presence)

Retention period for control plane logs in days.
Default: "30"

- default: `30`

### spec.logging.controlPlaneLogComponents

`[]string`

Control plane components to enable logging for.
Valid values: "apiserver", "kcm", "scheduler", "ccm",
  "controlplane-events", "alb", "coredns".

### spec.logging.auditLogEnabled

`bool`

Whether Kubernetes audit logging is enabled.
Audit logs record all API requests for security and compliance.

### spec.logging.auditLogSlsProject

`string`

SLS project for audit logs. If omitted and audit logging is enabled,
audit logs are sent to the control plane log project.

### spec.maintenanceWindow

`AliCloudKubernetesClusterMaintenanceWindow`

Maintenance window for the cluster. When configured, ACK applies
updates and patches only during the specified window.

### spec.maintenanceWindow.enable

`bool`

Whether the maintenance window is enabled.

### spec.maintenanceWindow.maintenanceTime

`string`

Start time of the maintenance window in RFC 3339 format.
Example: "2026-03-01T03:00:00+08:00"

### spec.maintenanceWindow.duration

`string`

Duration of the maintenance window (1-24 hours).
Example: "3h"

### spec.maintenanceWindow.weeklyPeriod

`string`

Days of the week when maintenance is allowed.
Comma-separated, e.g., "Monday,Thursday".
Default: "Thursday"

### spec.autoUpgrade

`AliCloudKubernetesClusterAutoUpgrade`

Automatic cluster version upgrade policy.
Only takes effect when a maintenance_window is configured.

### spec.autoUpgrade.enabled

`bool`

Whether automatic cluster version upgrades are enabled.

### spec.autoUpgrade.channel

`string` · optional (explicit presence)

Upgrade channel that controls how aggressively versions are adopted.
"patch" -- only patch version upgrades (e.g., 1.28.3 -> 1.28.5).
"stable" -- minor version upgrades after stabilization period.
"rapid" -- minor version upgrades as soon as available.

- rule: channel must be one of: patch, stable, rapid

### spec.tags

`map<string, string>`

Tags applied to the cluster and its auto-created resources.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the cluster is placed in the account's default resource group.

### spec.timezone

`string`

IANA timezone for the cluster nodes (e.g., "Asia/Shanghai", "UTC").

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudKubernetesCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | ACK cluster ID assigned by Alibaba Cloud. |
| `status.outputs.cluster_name` | `string` | Cluster name as created. |
| `status.outputs.api_server_internet` | `string` | Public API server endpoint for kubectl and other Kubernetes clients. Only populated when slb_internet_enabled is true. |
| `status.outputs.api_server_intranet` | `string` | Private (VPC-internal) API server endpoint. Always available regardless of slb_internet_enabled. |
| `status.outputs.vpc_id` | `string` | VPC ID where the cluster is deployed. Computed by the cluster resource from the provided VSwitches. |
| `status.outputs.security_group_id` | `string` | Security group ID used by the cluster nodes. Either the user-provided security_group_id or the auto-created one. |
| `status.outputs.nat_gateway_id` | `string` | NAT gateway ID auto-created by the cluster. Only populated when new_nat_gateway is true. |
| `status.outputs.worker_ram_role_name` | `string` | RAM role name attached to worker nodes. Downstream components (e.g., container registries, log services) can grant this role access to their resources. |
| `status.outputs.rrsa_oidc_issuer_url` | `string` | RRSA OIDC issuer URL. Empty when enable_rrsa is false. Used to configure RAM role trust policies for pod-level IAM. |
| `status.outputs.ram_oidc_provider_name` | `string` | RRSA OIDC provider name in RAM. Empty when enable_rrsa is false. Format: "ack-rrsa-c{cluster_id}" |
| `status.outputs.ram_oidc_provider_arn` | `string` | RRSA OIDC provider ARN. Empty when enable_rrsa is false. Format: "acs:ram::{account_id}:oidc-provider/ack-rrsa-c{cluster_id}" |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchIds` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.podVswitchIds` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.securityGroupId` | AliCloudSecurityGroup | `status.outputs.security_group_id` |
| `spec.encryptionProviderKey` | AliCloudKmsKey | `status.outputs.key_id` |
| `spec.logging.controlPlaneLogProject` | AliCloudLogProject | `status.outputs.project_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudKubernetesNodePool | `spec.clusterId` | `status.outputs.cluster_id` |

## See Also

- [Overview](../README.md)
