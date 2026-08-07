# Container Engine Cluster on OCI

Deploys an OKE (Oracle Container Engine for Kubernetes) cluster -- a managed Kubernetes control plane on OCI with configurable cluster type (basic vs enhanced), VCN-native or flannel pod networking, public or private API endpoint, optional OIDC authentication, and optional KMS secrets encryption. Worker nodes are managed separately via OciContainerEngineNodePool components. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for compartment, VCN, subnet, and security group wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **OKE Cluster** -- a managed Kubernetes control plane (API server, etcd, scheduler, controller manager) in the specified compartment and VCN, with the configured cluster type, CNI plugin, Kubernetes version, and endpoint access settings
- **API Server Endpoint** -- a network endpoint in the specified subnet with optional public IP and NSG association. Private endpoints restrict API access to the VCN; public endpoints allow kubectl from the internet
- **Pod Network Configuration** -- configured only when `cniType` is set; determines whether pods get VCN IP addresses (oci_vcn_ip_native) or use a flannel overlay network
- **Image Verification Policy** -- configured only when `imagePolicyConfig.isPolicyEnabled` is `true`; requires container images to be signed with specified KMS keys
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the cluster

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the cluster in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A VCN with subnets for the API endpoint and service load balancers. The VCN must have sufficient CIDR space that does not overlap with the pod and service CIDRs. Provide the VCN OCID directly or reference an OciVcn Cloud Resource via ValueFromRef.
- A regional subnet for the Kubernetes API server endpoint. For production, use a dedicated private subnet with restrictive NSGs.
- One or more subnets for Kubernetes Service load balancers (typically public subnets in different availability domains).
- A supported Kubernetes version. Use `oci ce cluster-options list` to see available versions.

## Deploy

### Console

Open the deployment store, find **Container Engine Cluster on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Production** preset in the [Presets](#presets) tab to pre-populate an enhanced cluster with VCN-native networking and a public API endpoint.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciContainerEngineCluster
metadata:
  name: platform-cluster
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vcnId:
    value: "ocid1.vcn.oc1..example"
  kubernetesVersion: "v1.30.1"
  type: enhanced_cluster
  cniType: oci_vcn_ip_native
  endpointConfig:
    subnetId:
      value: "ocid1.subnet.oc1..example"
    isPublicIpEnabled: true
```

```shell
planton apply -f oke-cluster.yaml
```

This creates an enhanced OKE cluster with VCN-native pod networking and a publicly accessible API endpoint. No service load balancer subnets, KMS encryption, or OIDC authentication are configured. Worker nodes must be added separately via an OciContainerEngineNodePool.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a compartment, VCN, subnets, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
  vcnId:
    valueFrom:
      kind: OciVcn
      name: platform-vcn
      fieldPath: status.outputs.vcnId
  endpointConfig:
    subnetId:
      valueFrom:
        kind: OciSubnet
        name: api-endpoint-subnet
        fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: api-endpoint-nsg
          fieldPath: status.outputs.networkSecurityGroupId
  options:
    serviceLbSubnetIds:
      - valueFrom:
          kind: OciSubnet
          name: service-lb-subnet
          fieldPath: status.outputs.subnetId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, VCN, subnets, and security groups first, then provisions the OKE cluster with the resolved values.

## Key Configuration

These are the most important decisions when configuring an OKE cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster type** -- The `type` field selects between `basic_cluster` (standard features, lower cost) and `enhanced_cluster` (virtual node pools, cluster add-on management, workload identity). Enhanced is recommended for production; basic is sufficient for development. Changing the type after creation is not supported.

**CNI plugin** -- The `cniType` selects pod networking: `oci_vcn_ip_native` assigns VCN IP addresses directly to pods, enabling NSGs on pods and native network policies. `flannel_overlay` uses an overlay network with simpler subnet requirements. VCN-native is recommended for production. Changing CNI after creation forces cluster recreation.

**Public vs private API endpoint** -- Set `endpointConfig.isPublicIpEnabled: true` for kubectl access from the internet; `false` for VCN-only access requiring a VPN, bastion, or peered network. Apply `endpointConfig.nsgIds` to restrict traffic to the API server regardless of public/private setting.

**Pod and service CIDRs** -- The `options.kubernetesNetworkConfig` sets `podsCidr` (default `10.244.0.0/16`) and `servicesCidr` (default `10.96.0.0/16`). These must not overlap with the VCN CIDR or each other. Plan these ranges before creation -- they force cluster recreation if changed.

**OIDC authentication** -- The `options.openIdConnectTokenAuthenticationConfig` enables external identity provider integration for kubectl and API access. Supports both inline configuration (issuer URL, client ID, claims) and a base64-encoded configuration file. Set `isOpenIdConnectDiscoveryEnabled: true` to expose the cluster's OIDC Discovery endpoint for external workload identity federation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciVcn** | `vcnId` | `status.outputs.vcnId` |
| **OciSubnet** | `endpointConfig.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `endpointConfig.nsgIds` | `status.outputs.networkSecurityGroupId` |
| **OciSubnet** (optional) | `options.serviceLbSubnetIds` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `options.serviceLbConfig.backendNsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | OCID of the OKE cluster | OciContainerEngineNodePool attachment, monitoring dashboards |
| `kubernetes_version` | Kubernetes version running on the control plane | Node pool version pinning, upgrade coordination |
| `kubernetes_endpoint` | Kubernetes API server endpoint URL (non-native networking) | kubectl configuration, CI/CD pipeline access |
| `private_endpoint` | Private VCN-native API server endpoint URL | kubectl from within the VCN, bastion port forwarding |
| `public_endpoint` | Public VCN-native API server endpoint URL (empty when private-only) | kubectl from the internet, CI/CD pipeline access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard production** -- An enhanced cluster with VCN-native pod networking, a public API endpoint with NSG protection, and service load balancer subnet configuration. The starting point for most production OKE deployments. Start from the **Standard Production** preset.

**Private cluster** -- An enhanced cluster with VCN-native networking, a private-only API endpoint (no public IP), KMS secrets encryption, and image verification policy. Required for compliance-sensitive environments where all cluster access goes through a VPN or bastion. Start from the **Private Cluster** preset.

**Development** -- A basic cluster with flannel overlay networking, a public API endpoint, and minimal configuration. Lower cost, simpler subnet requirements, suitable for development and experimentation. Start from the **Development** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this cluster
- [**Virtual Cloud Network on OCI**](/cloud-catalog/oci-vcn) -- provides the VCN for cluster networking
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides subnets for the API endpoint and service load balancers
- [**Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides NSGs for API endpoint and service load balancer traffic control