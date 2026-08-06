# OciContainerEngineCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciContainerEngineClusterSpec defines the specification for an Oracle Cloud
Infrastructure Container Engine for Kubernetes (OKE) cluster.

An OKE cluster is a managed Kubernetes control plane that runs on OCI
infrastructure. The cluster manages the API server, etcd, scheduler, and
controller manager. Worker nodes are managed separately via
OciContainerEngineNodePool components.

OKE supports two cluster types: basic (standard features) and enhanced
(virtual node pools, cluster add-on management, workload identity). Pod
networking supports OCI VCN-native IP allocation (pods get VCN IPs) or
flannel overlay (simpler, legacy).

Deprecated provider features intentionally omitted:
  - add_ons (Kubernetes Dashboard, Tiller): removed in modern Kubernetes
  - admission_controller_options (Pod Security Policy): removed in K8s 1.25

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.vcnId` | `string \| valueFrom` | yes |  | OciVcn (`status.outputs.vcn_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.kubernetesVersion` | `string` | yes |  |  |
| `spec.type` | `enum` |  |  |  |
| `spec.cniType` | `enum` |  |  |  |
| `spec.endpointConfig` | `EndpointConfig` |  |  |  |
| `spec.endpointConfig.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.endpointConfig.isPublicIpEnabled` | `bool` |  |  |  |
| `spec.endpointConfig.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.options` | `ClusterOptions` |  |  |  |
| `spec.options.kubernetesNetworkConfig` | `KubernetesNetworkConfig` |  |  |  |
| `spec.options.kubernetesNetworkConfig.podsCidr` | `string` |  |  |  |
| `spec.options.kubernetesNetworkConfig.servicesCidr` | `string` |  |  |  |
| `spec.options.serviceLbSubnetIds` | `[]string \| valueFrom` |  |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.options.ipFamilies` | `[]enum` |  |  |  |
| `spec.options.serviceLbConfig` | `ServiceLbConfig` |  |  |  |
| `spec.options.serviceLbConfig.backendNsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.options.serviceLbConfig.freeformTags` | `map<string, string>` |  |  |  |
| `spec.options.serviceLbConfig.definedTags` | `map<string, string>` |  |  |  |
| `spec.options.persistentVolumeConfig` | `PersistentVolumeConfig` |  |  |  |
| `spec.options.persistentVolumeConfig.freeformTags` | `map<string, string>` |  |  |  |
| `spec.options.persistentVolumeConfig.definedTags` | `map<string, string>` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig` | `OpenIdConnectTokenAuthenticationConfig` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.isOpenIdConnectAuthEnabled` | `bool` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.configurationFile` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.issuerUrl` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.clientId` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.caCertificate` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.usernameClaim` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.usernamePrefix` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.groupsClaim` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.groupsPrefix` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.signingAlgorithms` | `[]string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims` | `[]RequiredClaim` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims[].key` | `string` |  |  |  |
| `spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims[].value` | `string` |  |  |  |
| `spec.options.isOpenIdConnectDiscoveryEnabled` | `bool` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.imagePolicyConfig` | `ImagePolicyConfig` |  |  |  |
| `spec.imagePolicyConfig.isPolicyEnabled` | `bool` |  |  |  |
| `spec.imagePolicyConfig.keyDetails` | `[]ImagePolicyKeyDetail` |  |  |  |
| `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | `string \| valueFrom` |  |  | OciKmsKey (`status.outputs.key_id`) |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the cluster will be created.
Changing this after creation forces cluster recreation.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.vcnId

`string | valueFrom` · required

OCID of the VCN where the cluster will be deployed.
Changing this after creation forces cluster recreation.

- references: OciVcn (`status.outputs.vcn_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciVcn, name: <that resource's name>, fieldPath: status.outputs.vcn_id}} -- a bare string does not parse

### spec.name

`string`

Human-readable name for the cluster shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.kubernetesVersion

`string` · required

Kubernetes version to install on the cluster control plane.
Example: "v1.28.2". Use `oci ce cluster-options list` to see available versions.

- rule: {"string":{"minLen":"1"}}

### spec.type

`enum`

Cluster type. Enhanced clusters provide virtual node pools, workload
identity, and cluster add-on lifecycle management.

Allowed values (use exactly as shown):

- `unspecified`
- `basic_cluster`
- `enhanced_cluster`

### spec.cniType

`enum`

Container Network Interface (CNI) plugin for pod networking.
oci_vcn_ip_native assigns VCN IP addresses directly to pods (recommended
for production; enables network policies and security groups on pods).
flannel_overlay uses an overlay network (simpler, legacy).
Changing this after creation forces cluster recreation.

Allowed values (use exactly as shown):

- `cni_unspecified`
- `flannel_overlay`
- `oci_vcn_ip_native`

### spec.endpointConfig

`EndpointConfig`

Network configuration for the Kubernetes API server endpoint.
Controls whether the API server is reachable via a public or private IP,
which subnet hosts the endpoint, and which NSGs protect it.

### spec.endpointConfig.subnetId

`string | valueFrom` · required

OCID of the regional subnet hosting the cluster API server endpoint.
Changing this after creation forces cluster recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.endpointConfig.isPublicIpEnabled

`bool` · optional (explicit presence)

Whether to assign a public IP to the API server endpoint.
Set to false for private clusters. Must be false when the subnet is
private (no public IP assignment possible).

### spec.endpointConfig.nsgIds

`[]string | valueFrom`

OCIDs of network security groups applied to the API server endpoint.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.options

`ClusterOptions`

Optional cluster configuration for networking, service load balancers,
persistent volumes, and OIDC authentication.

### spec.options.kubernetesNetworkConfig

`KubernetesNetworkConfig`

Kubernetes pod and service CIDR configuration.
Changing this after creation forces cluster recreation.

### spec.options.kubernetesNetworkConfig.podsCidr

`string`

CIDR block for Kubernetes pods.
Default when omitted: 10.244.0.0/16 (IPv4), fd00:eeee:eeee:0000::/96 (IPv6).

### spec.options.kubernetesNetworkConfig.servicesCidr

`string`

CIDR block for Kubernetes services (ClusterIP range).
Default when omitted: 10.96.0.0/16 (IPv4), fd00:eeee:eeee:0001::/108 (IPv6).

### spec.options.serviceLbSubnetIds

`[]string | valueFrom`

OCIDs of subnets where Kubernetes Service load balancers will be placed.
Typically one or two public subnets in different availability domains.
Changing this after creation forces cluster recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.options.ipFamilies

`[]enum`

IP address family for the cluster. Use [ipv4] for IPv4-only (default)
or [ipv4, ipv6] for dual-stack.
Changing this after creation forces cluster recreation.

Allowed values (use exactly as shown):

- `ip_family_unspecified`
- `ipv4`
- `ipv6`

### spec.options.serviceLbConfig

`ServiceLbConfig`

Default configuration applied to load balancers created by Kubernetes
Service resources of type LoadBalancer.

### spec.options.serviceLbConfig.backendNsgIds

`[]string | valueFrom`

OCIDs of NSGs applied to load balancer backends.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.options.serviceLbConfig.freeformTags

`map<string, string>`

Freeform tags applied to service load balancers.

### spec.options.serviceLbConfig.definedTags

`map<string, string>`

Defined tags applied to service load balancers.
Keys use the format "namespace.key" (e.g. "Operations.CostCenter").

### spec.options.persistentVolumeConfig

`PersistentVolumeConfig`

Default configuration applied to block volumes created by Kubernetes
PersistentVolumeClaim resources.

### spec.options.persistentVolumeConfig.freeformTags

`map<string, string>`

Freeform tags applied to persistent volumes.

### spec.options.persistentVolumeConfig.definedTags

`map<string, string>`

Defined tags applied to persistent volumes.
Keys use the format "namespace.key" (e.g. "Operations.CostCenter").

### spec.options.openIdConnectTokenAuthenticationConfig

`OpenIdConnectTokenAuthenticationConfig`

OpenID Connect token authentication for the Kubernetes API server.
Enables external identity provider integration for kubectl and API access.

### spec.options.openIdConnectTokenAuthenticationConfig.isOpenIdConnectAuthEnabled

`bool`

Whether OIDC token authentication is enabled.

### spec.options.openIdConnectTokenAuthenticationConfig.configurationFile

`string`

Base64-encoded Kubernetes OIDC Auth Config file.
Mutually exclusive with the inline fields below.

### spec.options.openIdConnectTokenAuthenticationConfig.issuerUrl

`string`

URL of the OIDC identity provider. Must use https://.
The API server uses this to discover signing keys.

### spec.options.openIdConnectTokenAuthenticationConfig.clientId

`string`

Client ID that all tokens must be issued for.

### spec.options.openIdConnectTokenAuthenticationConfig.caCertificate

`string`

Base64-encoded public RSA or ECDSA certificate of the identity provider.

### spec.options.openIdConnectTokenAuthenticationConfig.usernameClaim

`string`

JWT claim to use as the Kubernetes username. Default: "sub".

### spec.options.openIdConnectTokenAuthenticationConfig.usernamePrefix

`string`

Prefix prepended to username claims to prevent collisions with
existing names (e.g. "oidc:").

### spec.options.openIdConnectTokenAuthenticationConfig.groupsClaim

`string`

JWT claim to use as the user's group. Must be an array of strings
in the token.

### spec.options.openIdConnectTokenAuthenticationConfig.groupsPrefix

`string`

Prefix prepended to group claims.

### spec.options.openIdConnectTokenAuthenticationConfig.signingAlgorithms

`[]string`

Accepted signing algorithms for tokens. Default: ["RS256"].

### spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims

`[]RequiredClaim`

Key-value pairs of required claims in the ID token. If any required
claim is missing or has a different value, authentication is rejected.

### spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims[].key

`string`

### spec.options.openIdConnectTokenAuthenticationConfig.requiredClaims[].value

`string`

### spec.options.isOpenIdConnectDiscoveryEnabled

`bool`

When true, enables the cluster-specific OIDC Discovery endpoint,
allowing external systems to discover the cluster's public signing keys.

### spec.kmsKeyId

`string | valueFrom`

OCID of the KMS key to encrypt Kubernetes secrets at rest.
Requires kubernetes_version >= v1.13.0.
Changing this after creation forces cluster recreation.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.imagePolicyConfig

`ImagePolicyConfig`

Image verification policy for the cluster. When enabled, all container
images deployed to the cluster must be signed with one of the specified
KMS keys.

### spec.imagePolicyConfig.isPolicyEnabled

`bool`

Whether the image verification policy is enabled.

### spec.imagePolicyConfig.keyDetails

`[]ImagePolicyKeyDetail`

KMS keys used for image signature verification.

### spec.imagePolicyConfig.keyDetails[].kmsKeyId

`string | valueFrom`

OCID of the KMS key used to verify image signatures.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciContainerEngineCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | OCID of the OKE cluster. |
| `status.outputs.kubernetes_version` | `string` | Kubernetes version running on the cluster control plane. |
| `status.outputs.kubernetes_endpoint` | `string` | Kubernetes API server endpoint URL (non-native networking). |
| `status.outputs.private_endpoint` | `string` | Private native networking Kubernetes API server endpoint URL. |
| `status.outputs.public_endpoint` | `string` | Public native networking Kubernetes API server endpoint URL. Empty when the cluster endpoint is private-only. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.vcnId` | OciVcn | `status.outputs.vcn_id` |
| `spec.endpointConfig.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.endpointConfig.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.options.serviceLbSubnetIds` | OciSubnet | `status.outputs.subnet_id` |
| `spec.options.serviceLbConfig.backendNsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.kmsKeyId` | OciKmsKey | `status.outputs.key_id` |
| `spec.imagePolicyConfig.keyDetails[].kmsKeyId` | OciKmsKey | `status.outputs.key_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciContainerEngineNodePool | `spec.clusterId` | `status.outputs.cluster_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
