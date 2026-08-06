---
title: "GKE Cluster"
description: "GKE Cluster deployment documentation"
icon: "package"
order: 100
componentName: "gcpgkecluster"
---

# GCP GKE Cluster

Deploys a Google Kubernetes Engine cluster — the managed Kubernetes control plane plus every cluster-wide configuration surface: private topology, upgrades and maintenance, node auto-provisioning, security, observability, addons, and Autopilot mode.

## What Gets Created

When you deploy a GcpGkeCluster resource, Planton provisions:

- **GKE cluster** — a `google_container_cluster` in the specified project and location, tagged with organization, environment, and resource labels; the Kubernetes Engine API is enabled automatically so a fresh project works on the first deploy
- **VPC-native networking** — pods and services draw from secondary ranges on your subnetwork (named ranges, or GKE-managed ranges when none are named); the default node pool is removed at create time on Standard clusters so all compute comes from `GcpGkeNodePool` resources
- **Private topology** (when configured) — private nodes, a peering- or PSC-based private control plane, master authorized networks, and the DNS endpoint
- **Autopilot cluster** (when `enableAutopilot` is set) — GKE provisions and manages nodes, bills per pod, and enforces a hardened posture; no node pool resources attach

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** (or the provider's default project) where the cluster is created
- **A VPC network and subnetwork** — referenced via `network` and `subnetwork`
- **A Cloud NAT** (`GcpRouterNat`) on the network if nodes are private and need outbound internet
- **A Cloud KMS key** if using CMEK etcd encryption

## Quick Start

Create a file `cluster.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeCluster
metadata:
  name: dev-cluster
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpGkeCluster.dev-cluster
spec:
  projectId:
    value: my-gcp-project
  location: us-central1-a
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: my-vpc
      fieldPath: status.outputs.network_self_link
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: my-subnet
      fieldPath: status.outputs.subnetwork_self_link
  deletionProtection: false
```

Deploy:

```shell
planton apply -f cluster.yaml
```

This creates a zonal cluster with GKE-managed pod/service ranges, Workload Identity on, the REGULAR release channel, and no node pools (add a `GcpGkeNodePool` to run workloads).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `location` | `string` | Region (`us-central1`) for a regional cluster or zone (`us-central1-a`) for a zonal one. Immutable. | GCP region/zone format |
| `network` | `StringValueOrRef` | The VPC network. Can reference a GcpVpcNetwork resource via `valueFrom`. Immutable. | Required |
| `subnetwork` | `StringValueOrRef` | The subnetwork nodes attach to; must be in the cluster's region. Can reference a GcpSubnetwork resource. Immutable. | Required |

### Core Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference a GcpProject resource. |
| `clusterName` | `string` | `metadata.name` | Cluster name in GCP (1–40 chars, lowercase). Immutable. |
| `description` | `string` | — | Human-readable description. Immutable. |
| `nodeLocations` | `string[]` | GKE defaults | Zones nodes run in (narrows a regional cluster / widens a zonal one). Mutable. |
| `resourceLabels` | `map<string,string>` | `{}` | GCE labels merged with the platform labels. Mutable. |
| `deletionProtection` | `bool` | `true` | While true, a destroy plan fails before touching the cluster. |
| `enableAutopilot` | `bool` | `false` | Autopilot mode: GKE manages nodes; node-management fields are rejected pre-deploy. Immutable. |
| `allowNetAdmin` | `bool` | `false` | Autopilot only: permit NET_ADMIN workloads. |
| `fleetProject` | `string` | — | Registers the cluster with a fleet in the given project. |

### Networking

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `ipAllocation.clusterSecondaryRangeName` | `StringValueOrRef` | GKE-managed | Existing subnetwork range for pod IPs. Immutable. Exclusive with `clusterIpv4CidrBlock`. |
| `ipAllocation.servicesSecondaryRangeName` | `StringValueOrRef` | GKE-managed | Existing subnetwork range for service IPs. Immutable. Exclusive with `servicesIpv4CidrBlock`. |
| `ipAllocation.clusterIpv4CidrBlock` / `servicesIpv4CidrBlock` | `string` | GKE-picked | CIDR (or netmask size like `/14`) for GKE-created ranges. Immutable. |
| `ipAllocation.stackType` | `string` | `IPV4` | `IPV4` or `IPV4_IPV6` (dual-stack; needs a dual-stack subnetwork). Immutable. |
| `ipAllocation.additionalPodRangeNames` | `string[]` | — | Additional pod ranges node pools may use — the post-creation pod-space growth lever. Mutable. |
| `ipAllocation.podCidrOverprovisionDisabled` | `bool` | `false` | Disables the 2× per-node pod-CIDR overprovisioning. Immutable. |
| `datapathProvider` | `string` | GKE default | `ADVANCED_DATAPATH` (Dataplane V2, recommended) or `LEGACY_DATAPATH`. Immutable. |
| `defaultMaxPodsPerNode` | `int` | `110` | Cluster default pods-per-node (8–256). Immutable; pools can override. |
| `enableIntranodeVisibility` | `bool` | `false` | Mirrors same-node pod traffic to the VPC dataplane (flow logs visibility). |
| `enableL4IlbSubsetting` | `bool` | `false` | L4 ILB subsetting for >250-node clusters. One-way (cannot be disabled). |
| `enableFqdnNetworkPolicy` / `enableCiliumClusterwideNetworkPolicy` | `bool` | `false` | FQDN / cluster-wide policy objects. Require Dataplane V2. |
| `enableMultiNetworking` | `bool` | `false` | Additional pod network interfaces. Immutable; requires Dataplane V2. |
| `enableNetworkPolicy` | `bool` | `false` | Calico NetworkPolicy enforcement (legacy path — Dataplane V2 enforces natively). |
| `privateIpv6GoogleAccess` | `string` | GCP default | `PRIVATE_IPV6_GOOGLE_ACCESS_TO_GOOGLE`, `_BIDIRECTIONAL`, or `_DISABLED`. |
| `inTransitEncryption` | `string` | GCP default | `IN_TRANSIT_ENCRYPTION_INTER_NODE_TRANSPARENT` encrypts inter-node pod traffic (Dataplane V2). |
| `disableDefaultSnat` | `bool` | `false` | Preserves pod IPs end to end for routable-pod designs. |
| `dnsConfig` | object | GKE default | Cloud DNS vs kube-dns, scope (`CLUSTER_SCOPE`/`VPC_SCOPE`), custom domain, additive VPC domain. |
| `gatewayApiChannel` | `string` | — | `CHANNEL_STANDARD` installs the Gateway API controller; `CHANNEL_DISABLED` turns it off. |
| `enableServiceExternalIps` | `bool` | `false` | Allows Services to use external IPs (CVE-2020-8554 mitigation keeps it off). |
| `totalEgressBandwidthTier` | `string` | — | `TIER_1` unlocks high-bandwidth egress on supported machine families. |
| `disableL4LbFirewallReconciliation` | `bool` | `false` | Stops GKE reconciling its L4 LB VPC firewall rules. |

### Private Topology & Control-Plane Access

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `privateCluster.enablePrivateNodes` | `bool` | `false` | Nodes get only internal IPs (compose a GcpRouterNat for egress). |
| `privateCluster.enablePrivateEndpoint` | `bool` | `false` | Removes the public control-plane endpoint entirely. Requires private nodes. |
| `privateCluster.masterIpv4CidrBlock` | `string` | — | RFC1918 `/28` for peering-based control planes. Immutable. Exclusive with `privateEndpointSubnetwork`. |
| `privateCluster.privateEndpointSubnetwork` | `StringValueOrRef` | — | Subnetwork for PSC-based control-plane endpoint placement. Immutable. |
| `privateCluster.enableMasterGlobalAccess` | `bool` | `false` | Private endpoint reachable from other regions / on-prem. |
| `masterAuthorizedNetworks.cidrBlocks[]` | list | — | CIDRs allowed to reach the API server (`cidrBlock` + `displayName`). |
| `masterAuthorizedNetworks.gcpPublicCidrsAccessEnabled` | `bool` | GCP default | Whether Google Cloud public IPs may reach the control plane. |
| `masterAuthorizedNetworks.privateEndpointEnforcementEnabled` | `bool` | GCP default | Enforce the allowlist on the private endpoint too. |
| `controlPlaneEndpoints.dnsEndpointAllowExternalTraffic` | `bool` | `false` | The `*.gke.goog` DNS endpoint accepts IAM-authenticated traffic from outside the VPC. |
| `controlPlaneEndpoints.ipEndpointsEnabled` | `bool` | `true` | Set false for DNS-endpoint-only clusters. |

### Upgrades, Maintenance & Autoscaling

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `releaseChannel` | `enum` | `REGULAR` | `RAPID`, `REGULAR`, `STABLE`, `EXTENDED`, or `NONE` (opt out of channel upgrades). |
| `minMasterVersion` | `string` | channel-driven | Minimum control-plane Kubernetes version. |
| `maintenancePolicy.dailyWindow.startTime` | `string` | — | Daily 4-hour window start, `HH:MM` UTC. Exactly one of daily/recurring. |
| `maintenancePolicy.recurringWindow` | object | — | RFC3339 start/end + RFC5545 `recurrence` (e.g. weekends only). |
| `maintenancePolicy.exclusions[]` | list (≤20) | — | Change freezes: name, start/end, scope (`NO_UPGRADES`, `NO_MINOR_UPGRADES`, `NO_MINOR_OR_NODE_UPGRADES`). |
| `clusterAutoscaling.enabled` | `bool` | `false` | Node auto-provisioning: GKE creates/deletes node pools within limits. |
| `clusterAutoscaling.resourceLimits[]` | list | — | Per-resource bounds (`cpu`, `memory`, accelerators). Required when NAP is enabled. |
| `clusterAutoscaling.autoscalingProfile` | `string` | `BALANCED` | `BALANCED` or `OPTIMIZE_UTILIZATION`. |
| `clusterAutoscaling.autoProvisioningDefaults` | object | — | SA ref, OAuth scopes, disk size/type, image, min CPU platform, boot-disk CMEK, shielding, auto-repair/upgrade. |
| `enableVerticalPodAutoscaling` | `bool` | `false` | VPA recommendations/updates for pod requests. |
| `hpaProfile` | `string` | — | `PERFORMANCE` for faster HPA reactions. |

### Security

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `workloadIdentityEnabled` | `bool` | `true` | Workload Identity Federation for GKE (`PROJECT_ID.svc.id.goog`). |
| `enableShieldedNodes` | `bool` | GCP default (`true`) | Secure boot + integrity monitoring on nodes. Leave unset on Autopilot. |
| `databaseEncryption.state` | `string` | `DECRYPTED` | `ENCRYPTED` wraps Kubernetes secrets with `keyName` (CMEK). |
| `databaseEncryption.keyName` | `StringValueOrRef` | — | Cloud KMS key (GcpKmsKey reference). Required when state is `ENCRYPTED`. |
| `binaryAuthorizationEvaluationMode` | `string` | — | `PROJECT_SINGLETON_POLICY_ENFORCE` admits only policy-satisfying images. |
| `securityPosture.mode` / `.vulnerabilityMode` | `string` | — | Security Posture dashboard auditing and vulnerability scanning tiers. |
| `authenticatorSecurityGroup` | `string` | — | `gke-security-groups@YOURDOMAIN` Google Group for RBAC. |
| `confidentialNodes` | object | — | Hardware memory encryption (`SEV`, `SEV_SNP`, `TDX`). Immutable. |
| `anonymousAuthenticationMode` | `string` | GCP default | `LIMITED` restricts anonymous access to health endpoints. |
| `enableIdentityService` | `bool` | `false` | External OIDC identity providers for the Kubernetes API. |
| `enableMeshCertificates` | `bool` | `false` | Workload mTLS certificates via the mesh CA. |
| `enableSecretManagerCsi` | `bool` | `false` | Mount Secret Manager secrets via the built-in CSI driver. |
| `enableLegacyAbac` | `bool` | `false` | Legacy ABAC (leave off). |

### Observability & Cost

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `logging.components` | `string[]` | GKE default | Components shipping logs (`SYSTEM_COMPONENTS`, `WORKLOADS`, `APISERVER`, …). Empty list disables Cloud Logging. |
| `monitoring.components` | `string[]` | GKE default | Components shipping metrics (incl. `POD`, `DEPLOYMENT`, `DCGM`, …). |
| `monitoring.managedPrometheusEnabled` | `bool` | `true` | Managed Service for Prometheus collection. |
| `monitoring.autoMonitoringScope` | `string` | — | `ALL` deploys packaged PodMonitorings automatically. |
| `monitoring.advancedDatapathMetricsEnabled` / `RelayEnabled` | `bool` | `false` | Dataplane V2 observability metrics / Hubble relay. |
| `notificationPubsub` | object | — | Lifecycle events to a Pub/Sub topic (GcpPubSubTopic ref), optionally filtered by event type. |
| `enableCostManagement` | `bool` | `false` | Per-namespace cost allocation in the billing export. |
| `resourceUsageExport` | object | — | Continuous usage metering into a BigQuery dataset (GcpBigQueryDataset ref). |

### Addons

| Field | Default | Description |
|-------|---------|-------------|
| `addons.httpLoadBalancingEnabled` | `true` | GCE ingress controller. |
| `addons.horizontalPodAutoscalingEnabled` | `true` | HPA controller. |
| `addons.gcePersistentDiskCsiDriverEnabled` | `true` | PD CSI driver. |
| `addons.gcpFilestoreCsiDriverEnabled` | `false` | Filestore (NFS) CSI driver. |
| `addons.gcsFuseCsiDriverEnabled` | `false` | GCS FUSE CSI driver. |
| `addons.gkeBackupAgentEnabled` | `false` | Backup for GKE agent. |
| `addons.dnsCacheEnabled` | `false` | NodeLocal DNSCache. |
| `addons.configConnectorEnabled` | `false` | Config Connector. |
| `addons.statefulHaEnabled` | `false` | Stateful HA operator. |
| `addons.rayOperatorEnabled` | `false` | Ray (KubeRay) operator. |

## Examples

### Private Production Cluster (Standard)

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeCluster
metadata:
  name: prod-primary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeCluster.prod-primary
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: prod-gke-subnet
      fieldPath: status.outputs.subnetwork_self_link
  ipAllocation:
    clusterSecondaryRangeName:
      value: pods
    servicesSecondaryRangeName:
      value: services
  datapathProvider: ADVANCED_DATAPATH
  privateCluster:
    enablePrivateNodes: true
    masterIpv4CidrBlock: "172.16.0.16/28"
  masterAuthorizedNetworks:
    cidrBlocks:
      - cidrBlock: 203.0.113.0/24
        displayName: corp-vpn
  maintenancePolicy:
    dailyWindow:
      startTime: "03:00"
  securityPosture:
    mode: BASIC
    vulnerabilityMode: VULNERABILITY_BASIC
  enableCostManagement: true
```

### Autopilot Cluster

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeCluster
metadata:
  name: autopilot-primary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeCluster.autopilot-primary
spec:
  projectId:
    value: my-gcp-project
  location: us-central1
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: prod-gke-subnet
      fieldPath: status.outputs.subnetwork_self_link
  enableAutopilot: true
  privateCluster:
    enablePrivateNodes: true
    masterIpv4CidrBlock: "172.16.0.32/28"
```

### CMEK + Hardened Control Plane

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeCluster
metadata:
  name: secure-cluster
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeCluster.secure-cluster
spec:
  projectId:
    value: my-gcp-project
  location: us-east1
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: secure-vpc
      fieldPath: status.outputs.network_self_link
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: secure-subnet
      fieldPath: status.outputs.subnetwork_self_link
  privateCluster:
    enablePrivateNodes: true
    enablePrivateEndpoint: true
    masterIpv4CidrBlock: "172.16.0.48/28"
  databaseEncryption:
    state: ENCRYPTED
    keyName:
      valueFrom:
        kind: GcpKmsKey
        name: etcd-secrets-key
        fieldPath: status.outputs.key_id
  binaryAuthorizationEvaluationMode: PROJECT_SINGLETON_POLICY_ENFORCE
  anonymousAuthenticationMode: LIMITED
  confidentialNodes:
    enabled: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `endpoint` | `string` | Kubernetes API server IP (the private endpoint on private-only control planes) |
| `cluster_ca_certificate` | `string` | Base64 CA certificate for validating the API server's TLS certificate (public trust material) |
| `workload_identity_pool` | `string` | `PROJECT_ID.svc.id.goog`; empty when Workload Identity is disabled on a Standard cluster |
| `cluster_id` | `string` | `projects/{project}/locations/{location}/clusters/{name}` |
| `name` | `string` | The cluster name in GCP — the handle node pools and gcloud use |
| `location` | `string` | Region or zone, exactly as provided in the spec |
| `self_link` | `string` | Server-defined URL of the cluster resource |
| `master_version` | `string` | Kubernetes version currently running on the control plane |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — the network the cluster lives in
- [GcpSubnetwork](/docs/catalog/gcp/subnetwork) — carries the pod/service secondary ranges
- [GcpGkeNodePool](/docs/catalog/gcp/gke-node-pool) — compute for Standard clusters
- [GcpRouterNat](/docs/catalog/gcp/router-nat) — outbound internet for private nodes
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — CMEK key for etcd secrets encryption
- [GcpGkeWorkloadIdentityBinding](/docs/catalog/gcp/gke-workload-identity-binding) — binds Kubernetes service accounts to IAM service accounts
- [GcpPubSubTopic](/docs/catalog/gcp/pubsub-topic) — receives cluster lifecycle notifications
- [GcpBigQueryDataset](/docs/catalog/gcp/bigquery-dataset) — receives resource-usage export records
