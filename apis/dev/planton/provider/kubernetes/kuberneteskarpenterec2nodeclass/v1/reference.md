# KubernetesKarpenterEc2NodeClass

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

KubernetesKarpenterEc2NodeClassSpec declares a Karpenter EC2NodeClass —
the AWS-level machine template NodePools launch instances from: which
AMIs to boot, which subnets and security groups to attach, which IAM
identity nodes assume, disk layout, kubelet configuration, and instance
metadata posture. One NodeClass is typically shared by several NodePools
(the pools differ in constraints and taints; the class is the common
"how a node is built").

Requires a Karpenter installation (KubernetesKarpenter) on an EKS or
EKS-compatible cluster. The rendered EC2NodeClass is CLUSTER-SCOPED (no
namespace) and named after metadata.name. 100% fidelity with the
upstream karpenter.k8s.aws/v1 EC2NodeClass CRD at the pinned release
(aws/karpenter-provider-aws pkg/apis/crds/
karpenter.k8s.aws_ec2nodeclasses.yaml); the CRD's own CEL rules are
mirrored below so mistakes surface at validate time, not at apply.

Several field names here carry an explicit json_name: the CRD's keys use
acronym casing (IP, ID, DNS, CFS, GC) that protojson would otherwise
derive differently, and the API server rejects miscased keys as
undeclared fields.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKarpenterEc2NodeClass
metadata:
  name: default-al2023
spec:
  amiSelectorTerms:
    - alias: al2023@v20240807
  amiFamily: AL2023
  subnetSelectorTerms:
    - tags:
        karpenter.sh/discovery: my-eks-cluster
  securityGroupSelectorTerms:
    - tags:
        karpenter.sh/discovery: my-eks-cluster
  role: KarpenterNodeRole-my-eks-cluster
  blockDeviceMappings:
    - deviceName: /dev/xvda
      rootVolume: true
      ebs:
        deleteOnTermination: true
        encrypted: true
        iops: 3000
        throughput: 125
        volumeSize: 100Gi
        volumeType: gp3
  associatePublicIPAddress: false
  connectionTracking:
    tcpEstablishedTimeout: 86400
  detailedMonitoring: true
  instanceStorePolicy: RAID0
  ipPrefixCount: 4
  kubelet:
    clusterDNS:
      - 10.100.0.10
    cpuCFSQuota: true
    evictionHard:
      memory.available: 5%
    evictionMaxPodGracePeriod: 60
    evictionSoft:
      memory.available: 10%
    evictionSoftGracePeriod:
      memory.available: 1m0s
    imageGCHighThresholdPercent: 85
    imageGCLowThresholdPercent: 80
    kubeReserved:
      cpu: 200m
      memory: 512Mi
    maxPods: 110
    podsPerCore: 10
    systemReserved:
      cpu: 100m
      memory: 256Mi
  metadataOptions:
    httpEndpoint: enabled
    httpProtocolIPv6: disabled
    httpPutResponseHopLimit: 1
    httpTokens: required
  tags:
    team: platform
    environment: dev
  userData: |
    #!/bin/bash
    echo "extra bootstrap commands merged by karpenter"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.amiSelectorTerms` | `[]KubernetesKarpenterEc2NodeClassAmiSelectorTerm` | yes |  |  |
| `spec.amiSelectorTerms[].alias` | `string` |  |  |  |
| `spec.amiSelectorTerms[].id` | `string` |  |  |  |
| `spec.amiSelectorTerms[].name` | `string` |  |  |  |
| `spec.amiSelectorTerms[].owner` | `string` |  |  |  |
| `spec.amiSelectorTerms[].ssmParameter` | `string` |  |  |  |
| `spec.amiSelectorTerms[].tags` | `map<string, string>` |  |  |  |
| `spec.amiFamily` | `string` |  |  |  |
| `spec.subnetSelectorTerms` | `[]KubernetesKarpenterEc2NodeClassSubnetSelectorTerm` | yes |  |  |
| `spec.subnetSelectorTerms[].id` | `string` |  |  |  |
| `spec.subnetSelectorTerms[].tags` | `map<string, string>` |  |  |  |
| `spec.securityGroupSelectorTerms` | `[]KubernetesKarpenterEc2NodeClassSecurityGroupSelectorTerm` | yes |  |  |
| `spec.securityGroupSelectorTerms[].id` | `string` |  |  |  |
| `spec.securityGroupSelectorTerms[].name` | `string` |  |  |  |
| `spec.securityGroupSelectorTerms[].tags` | `map<string, string>` |  |  |  |
| `spec.role` | `string` | yes |  |  |
| `spec.instanceProfile` | `string` | yes |  |  |
| `spec.blockDeviceMappings` | `[]KubernetesKarpenterEc2NodeClassBlockDeviceMapping` |  |  |  |
| `spec.blockDeviceMappings[].deviceName` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs` | `KubernetesKarpenterEc2NodeClassEbs` |  |  |  |
| `spec.blockDeviceMappings[].ebs.deleteOnTermination` | `bool` |  |  |  |
| `spec.blockDeviceMappings[].ebs.encrypted` | `bool` |  |  |  |
| `spec.blockDeviceMappings[].ebs.iops` | `int64` |  |  |  |
| `spec.blockDeviceMappings[].ebs.kmsKeyID` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs.snapshotID` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs.throughput` | `int64` |  |  |  |
| `spec.blockDeviceMappings[].ebs.volumeInitializationRate` | `int32` |  |  |  |
| `spec.blockDeviceMappings[].ebs.volumeSize` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs.volumeType` | `string` |  |  |  |
| `spec.blockDeviceMappings[].rootVolume` | `bool` |  |  |  |
| `spec.capacityReservationSelectorTerms` | `[]KubernetesKarpenterEc2NodeClassCapacityReservationSelectorTerm` |  |  |  |
| `spec.capacityReservationSelectorTerms[].id` | `string` |  |  |  |
| `spec.capacityReservationSelectorTerms[].instanceMatchCriteria` | `string` |  |  |  |
| `spec.capacityReservationSelectorTerms[].ownerID` | `string` |  |  |  |
| `spec.capacityReservationSelectorTerms[].tags` | `map<string, string>` |  |  |  |
| `spec.associatePublicIPAddress` | `bool` |  |  |  |
| `spec.connectionTracking` | `KubernetesKarpenterEc2NodeClassConnectionTracking` |  |  |  |
| `spec.connectionTracking.tcpEstablishedTimeout` | `int32` |  |  |  |
| `spec.connectionTracking.udpStreamTimeout` | `int32` |  |  |  |
| `spec.connectionTracking.udpTimeout` | `int32` |  |  |  |
| `spec.context` | `string` |  |  |  |
| `spec.cpuOptions` | `KubernetesKarpenterEc2NodeClassCpuOptions` |  |  |  |
| `spec.cpuOptions.nestedVirtualization` | `string` |  |  |  |
| `spec.detailedMonitoring` | `bool` |  |  |  |
| `spec.instanceStorePolicy` | `string` |  |  |  |
| `spec.ipPrefixCount` | `int32` |  |  |  |
| `spec.kubelet` | `KubernetesKarpenterEc2NodeClassKubelet` |  |  |  |
| `spec.kubelet.clusterDNS` | `[]string` |  |  |  |
| `spec.kubelet.cpuCFSQuota` | `bool` |  |  |  |
| `spec.kubelet.evictionHard` | `map<string, string>` |  |  |  |
| `spec.kubelet.evictionMaxPodGracePeriod` | `int32` |  |  |  |
| `spec.kubelet.evictionSoft` | `map<string, string>` |  |  |  |
| `spec.kubelet.evictionSoftGracePeriod` | `map<string, string>` |  |  |  |
| `spec.kubelet.imageGCHighThresholdPercent` | `int32` |  |  |  |
| `spec.kubelet.imageGCLowThresholdPercent` | `int32` |  |  |  |
| `spec.kubelet.kubeReserved` | `map<string, string>` |  |  |  |
| `spec.kubelet.maxPods` | `int32` |  |  |  |
| `spec.kubelet.podsPerCore` | `int32` |  |  |  |
| `spec.kubelet.systemReserved` | `map<string, string>` |  |  |  |
| `spec.metadataOptions` | `KubernetesKarpenterEc2NodeClassMetadataOptions` |  |  |  |
| `spec.metadataOptions.httpEndpoint` | `string` |  | `enabled` |  |
| `spec.metadataOptions.httpProtocolIPv6` | `string` |  | `disabled` |  |
| `spec.metadataOptions.httpPutResponseHopLimit` | `int64` |  | `1` |  |
| `spec.metadataOptions.httpTokens` | `string` |  | `required` |  |
| `spec.networkInterfaces` | `[]KubernetesKarpenterEc2NodeClassNetworkInterface` |  |  |  |
| `spec.networkInterfaces[].deviceIndex` | `int32` |  |  |  |
| `spec.networkInterfaces[].interfaceType` | `string` | yes |  |  |
| `spec.networkInterfaces[].networkCardIndex` | `int32` |  |  |  |
| `spec.placementGroupSelector` | `KubernetesKarpenterEc2NodeClassPlacementGroupSelector` |  |  |  |
| `spec.placementGroupSelector.name` | `string` | yes |  |  |
| `spec.placementGroupSelector.id` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.userData` | `string` |  |  |  |

## Field Details

### spec.amiSelectorTerms

`[]KubernetesKarpenterEc2NodeClassAmiSelectorTerm` · required

AMI selector terms — how the machine image is chosen. Terms are ORed;
fields within a term are ANDed. The simplest arm is a single term with
an alias ("al2023@v20240807" or "bottlerocket@latest"); tag/name/id/
ssm-parameter selection serves custom-AMI pipelines. Pin versions in
production — "latest" drifts nodes when a new AMI ships.

- rule: each AMI selector term needs at least one of alias, id, name, ssm_parameter or tags
- rule: an AMI selector term with 'id' cannot combine it with alias, tags, name or owner
- rule: an AMI selector term with 'alias' cannot combine it with id, tags, name or owner
- rule: when an alias term is used it must be the ONLY AMI selector term
- rule: {"repeated":{"minItems":"1","maxItems":"30"}}

### spec.amiSelectorTerms[].alias

`string` · optional (explicit presence)

EKS-optimized AMI alias, "family@version" — families al2, al2023,
bottlerocket, windows2019, windows2022, windows2025. Version is an AMI
release ("al2023@v20240807") or "latest" (drifts nodes on new
releases; the only option for Windows families). Mutually exclusive
with every other selector field, and an alias term must be the list's
only term.

- rule: alias must match 'family@version' with family one of al2, al2023, bottlerocket, windows2019, windows2022, windows2025 (Windows families only support 'latest')
- rule: {"string":{"maxLen":"30"}}

### spec.amiSelectorTerms[].id

`string` · optional (explicit presence)

Explicit AMI id (ami-...). Mutually exclusive with the other fields.

- rule: {"string":{"pattern":"^ami-[0-9a-z]+$"}}

### spec.amiSelectorTerms[].name

`string` · optional (explicit presence)

AMI name in EC2 (the name FIELD, not the Name tag). ANDs with owner.

### spec.amiSelectorTerms[].owner

`string` · optional (explicit presence)

AMI owner filter: account IDs, "self", "amazon", "aws-marketplace".

### spec.amiSelectorTerms[].ssmParameter

`string` · optional (explicit presence)

Name (or ARN) of an SSM parameter holding the image id — the
custom-AMI-pipeline arm.

### spec.amiSelectorTerms[].tags

`map<string, string>`

Tag selectors ('*' value matches any). Empty keys/values are invalid.

- rule: AMI selector tags may not have empty keys or values
- rule: {"map":{"maxPairs":"20"}}

### spec.amiFamily

`string` · optional (explicit presence)

UserData format / bootstrap family of the chosen AMIs: AL2, AL2023,
Bottlerocket, Windows2019, Windows2022, Windows2025 or Custom.
Optional when an alias term is used (inferred from the alias — and
then only its own family or Custom may be set); REQUIRED when
selecting AMIs by id/name/tags/ssm-parameter.

- rule: ami_family must be one of 'AL2', 'AL2023', 'Bottlerocket', 'Custom', 'Windows2019', 'Windows2022' or 'Windows2025'

### spec.subnetSelectorTerms

`[]KubernetesKarpenterEc2NodeClassSubnetSelectorTerm` · required

Subnet selector terms — where nodes launch. Terms are ORed; selection
is by tags (the standard karpenter.sh/discovery pattern) or explicit
subnet id. Karpenter spreads across the selected subnets' zones.

- rule: each subnet selector term needs 'id' or 'tags'
- rule: a subnet selector term with 'id' cannot also carry tags
- rule: {"repeated":{"minItems":"1","maxItems":"30"}}

### spec.subnetSelectorTerms[].id

`string` · optional (explicit presence)

Explicit subnet id (subnet-...). Mutually exclusive with tags.

- rule: {"string":{"pattern":"^subnet-[0-9a-z]+$"}}

### spec.subnetSelectorTerms[].tags

`map<string, string>`

Tag selectors ('*' value matches any) — `karpenter.sh/discovery:
<cluster>` is the convention EKS setups tag their private subnets
with.

- rule: subnet selector tags may not have empty keys or values
- rule: {"map":{"maxPairs":"20"}}

### spec.securityGroupSelectorTerms

`[]KubernetesKarpenterEc2NodeClassSecurityGroupSelectorTerm` · required

Security-group selector terms — what nodes' ENIs attach. Terms are
ORed; selection by tags, group name, or group id.

- rule: each security-group selector term needs one of id, name or tags
- rule: a security-group selector term with 'id' cannot combine it with tags or name
- rule: a security-group selector term with 'name' cannot combine it with tags or id
- rule: {"repeated":{"minItems":"1","maxItems":"30"}}

### spec.securityGroupSelectorTerms[].id

`string` · optional (explicit presence)

Explicit security-group id (sg-...). Mutually exclusive with the
other fields.

- rule: {"string":{"pattern":"^sg-[0-9a-z]+$"}}

### spec.securityGroupSelectorTerms[].name

`string` · optional (explicit presence)

Security-group name (the name FIELD, not the Name tag). Mutually
exclusive with the other fields.

### spec.securityGroupSelectorTerms[].tags

`map<string, string>`

Tag selectors ('*' value matches any).

- rule: security-group selector tags may not have empty keys or values
- rule: {"map":{"maxPairs":"20"}}

### spec.role

`string` · required · optional (explicit presence)

IAM role name nodes assume (Karpenter manages the instance profile for
it — needs iam:PassRole on the controller). Exactly one of role /
instance_profile must be set; role is the recommended arm.

- rule: {"string":{"minLen":"1"}}

### spec.instanceProfile

`string` · required · optional (explicit presence)

Pre-existing IAM instance-profile name nodes use, for accounts where
Karpenter must not manage profiles. Exactly one of role /
instance_profile must be set.

- rule: {"string":{"minLen":"1"}}

### spec.blockDeviceMappings

`[]KubernetesKarpenterEc2NodeClassBlockDeviceMapping`

Block device mappings for provisioned nodes. Empty = the AMI family's
defaults. At most one mapping may be the root volume.

- rule: at most one block device mapping may set root_volume
- rule: {"repeated":{"maxItems":"50"}}

### spec.blockDeviceMappings[].deviceName

`string`

Device name (e.g. "/dev/xvda").

### spec.blockDeviceMappings[].ebs

`KubernetesKarpenterEc2NodeClassEbs`

EBS volume parameters for the device.

- rule: an EBS mapping needs snapshot_id or volume_size (or both) — without either there is nothing to create
- rule: volume_initialization_rate only applies when snapshot_id is set — it is the snapshot-download rate

### spec.blockDeviceMappings[].ebs.deleteOnTermination

`bool` · optional (explicit presence)

Delete the volume when the instance terminates.

### spec.blockDeviceMappings[].ebs.encrypted

`bool` · optional (explicit presence)

Encrypt the volume (cannot be set when snapshot_id supplies the
encryption state).

### spec.blockDeviceMappings[].ebs.iops

`int64` · optional (explicit presence)

Provisioned IOPS (gp3: 3000-16000; io1/io2: 100-64000; unsupported
for gp2/st1/sc1/standard).

### spec.blockDeviceMappings[].ebs.kmsKeyID

`string` · optional (explicit presence)

Customer-managed KMS key (id, alias, or ARN) for encryption.

### spec.blockDeviceMappings[].ebs.snapshotID

`string` · optional (explicit presence)

Snapshot to create the volume from. Either snapshot_id or volume_size
must be set — mirrored from the CRD.

### spec.blockDeviceMappings[].ebs.throughput

`int64` · optional (explicit presence)

gp3 throughput in MiB/s (125-1000).

- rule: {"int64":{"lte":"1000","gte":"125"}}

### spec.blockDeviceMappings[].ebs.volumeInitializationRate

`int32` · optional (explicit presence)

EBS Provisioned Rate for Volume Initialization in MiB/s (100-300);
only valid with snapshot_id.

- rule: {"int32":{"lte":300,"gte":100}}

### spec.blockDeviceMappings[].ebs.volumeSize

`string` · optional (explicit presence)

Volume size as a quantity in Gi/G/Ti/T (e.g. "100Gi"). Either
snapshot_id or volume_size must be set.

- rule: {"string":{"pattern":"^((?:[1-9][0-9]{0,3}|[1-4][0-9]{4}|[5][0-8][0-9]{3}|59000)Gi|(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-3][0-9]{3}|64000)G|([1-9]||[1-5][0-7]|58)Ti|([1-9]||[1-5][0-9]|6[0-3]|64)T)$"}}

### spec.blockDeviceMappings[].ebs.volumeType

`string` · optional (explicit presence)

Volume type: standard, io1, io2, gp2, sc1, st1, gp3.

- rule: volume_type must be one of 'standard', 'io1', 'io2', 'gp2', 'sc1', 'st1' or 'gp3'

### spec.blockDeviceMappings[].rootVolume

`bool` · optional (explicit presence)

Marks this device as the kubelet root volume (at most one mapping).

### spec.capacityReservationSelectorTerms

`[]KubernetesKarpenterEc2NodeClassCapacityReservationSelectorTerm`

Capacity-reservation selector terms (requires the reservedCapacity
feature gate, on by default): which On-Demand Capacity Reservations /
capacity blocks this class may launch into.

- rule: each capacity-reservation selector term needs one of id, tags or instance_match_criteria
- rule: a capacity-reservation selector term with 'id' cannot combine it with tags, owner_id or instance_match_criteria
- rule: {"repeated":{"maxItems":"30"}}

### spec.capacityReservationSelectorTerms[].id

`string` · optional (explicit presence)

Explicit capacity-reservation id (cr-...). Mutually exclusive with
the other fields.

- rule: {"string":{"pattern":"^cr-[0-9a-z]+$"}}

### spec.capacityReservationSelectorTerms[].instanceMatchCriteria

`string` · optional (explicit presence)

Match by how the reservation accepts instances: "open" or "targeted".

- rule: instance_match_criteria must be 'open' or 'targeted'

### spec.capacityReservationSelectorTerms[].ownerID

`string` · optional (explicit presence)

Owning AWS account id (12 digits) — for reservations shared across
accounts.

- rule: {"string":{"pattern":"^[0-9]{12}$"}}

### spec.capacityReservationSelectorTerms[].tags

`map<string, string>`

Tag selectors ('*' value matches any).

- rule: capacity-reservation selector tags may not have empty keys or values
- rule: {"map":{"maxPairs":"20"}}

### spec.associatePublicIPAddress

`bool` · optional (explicit presence)

Assign public IPv4 addresses to instances. Unset = AWS defaults by
subnet setting (public subnets assign, private do not) — set only to
override.

### spec.connectionTracking

`KubernetesKarpenterEc2NodeClassConnectionTracking`

ENI idle connection-tracking timeouts in the launch template.

- rule: connection_tracking needs at least one of tcp_established_timeout, udp_stream_timeout or udp_timeout

### spec.connectionTracking.tcpEstablishedTimeout

`int32` · optional (explicit presence)

Idle timeout for established TCP connections, seconds (60-432000).

- rule: {"int32":{"lte":432000,"gte":60}}

### spec.connectionTracking.udpStreamTimeout

`int32` · optional (explicit presence)

Idle timeout for UDP stream flows, seconds (60-180).

- rule: {"int32":{"lte":180,"gte":60}}

### spec.connectionTracking.udpTimeout

`int32` · optional (explicit presence)

Idle timeout for single-transaction UDP flows, seconds (30-60).

- rule: {"int32":{"lte":60,"gte":30}}

### spec.context

`string`

Reserved field in the EC2 CreateFleet API (rarely used; AWS-assigned
contexts).

### spec.cpuOptions

`KubernetesKarpenterEc2NodeClassCpuOptions`

CPU options for instances (nested virtualization).

### spec.cpuOptions.nestedVirtualization

`string` · optional (explicit presence)

Nested virtualization: "enabled" filters to instance types supporting
it; "disabled" turns it off.

- rule: nested_virtualization must be 'enabled' or 'disabled'

### spec.detailedMonitoring

`bool`

Enable detailed (1-minute) CloudWatch monitoring on instances.

### spec.instanceStorePolicy

`string` · optional (explicit presence)

How instance-store (ephemeral NVMe) disks are handled: "RAID0" stripes
them into one volume for pod ephemeral storage. Unset = untouched.

- rule: instance_store_policy only supports 'RAID0'

### spec.ipPrefixCount

`int32` · optional (explicit presence)

Number of IPv4 prefixes to assign to each node's ENI (prefix
delegation — raises pod density beyond per-ENI IP limits).

- rule: {"int32":{"gte":0}}

### spec.kubelet

`KubernetesKarpenterEc2NodeClassKubelet`

Kubelet configuration for provisioned nodes (a supported subset of
upstream kubelet config, merged into bootstrap user data).

- rule: image_gc_high_threshold_percent must be greater than image_gc_low_threshold_percent
- rule: every eviction_soft signal needs a matching eviction_soft_grace_period entry (and vice versa)

### spec.kubelet.clusterDNS

`[]string`

Cluster DNS server IPs.

### spec.kubelet.cpuCFSQuota

`bool` · optional (explicit presence)

Enforce CPU CFS quota for containers with CPU limits.

### spec.kubelet.evictionHard

`map<string, string>`

Hard eviction thresholds by signal (memory.available,
nodefs.available, nodefs.inodesFree, imagefs.available,
imagefs.inodesFree, pid.available) — quantity or percentage values.

- rule: eviction_hard keys must be one of memory.available, nodefs.available, nodefs.inodesFree, imagefs.available, imagefs.inodesFree, pid.available

### spec.kubelet.evictionMaxPodGracePeriod

`int32` · optional (explicit presence)

Maximum grace period (seconds) for pods terminated on soft eviction.

### spec.kubelet.evictionSoft

`map<string, string>`

Soft eviction thresholds by signal (same keys as eviction_hard). Every
signal here must have a matching eviction_soft_grace_period.

- rule: eviction_soft keys must be one of memory.available, nodefs.available, nodefs.inodesFree, imagefs.available, imagefs.inodesFree, pid.available

### spec.kubelet.evictionSoftGracePeriod

`map<string, string>`

Grace periods per soft-eviction signal (same keys as eviction_soft).

- rule: eviction_soft_grace_period keys must be one of memory.available, nodefs.available, nodefs.inodesFree, imagefs.available, imagefs.inodesFree, pid.available

### spec.kubelet.imageGCHighThresholdPercent

`int32` · optional (explicit presence)

Disk usage percentage that always triggers image garbage collection
(0-100; must exceed image_gc_low_threshold_percent).

- rule: {"int32":{"lte":100,"gte":0}}

### spec.kubelet.imageGCLowThresholdPercent

`int32` · optional (explicit presence)

Disk usage percentage below which image GC never runs (0-100).

- rule: {"int32":{"lte":100,"gte":0}}

### spec.kubelet.kubeReserved

`map<string, string>`

Resources reserved for Kubernetes system components (keys cpu, memory,
ephemeral-storage, pid; non-negative quantities).

- rule: kube_reserved keys must be one of cpu, memory, ephemeral-storage, pid, with non-negative quantity values

### spec.kubelet.maxPods

`int32` · optional (explicit presence)

Maximum pods per node (overrides the ENI-based default).

- rule: {"int32":{"gte":0}}

### spec.kubelet.podsPerCore

`int32` · optional (explicit presence)

Pods per CPU core cap (max_pods still wins when lower).

- rule: {"int32":{"gte":0}}

### spec.kubelet.systemReserved

`map<string, string>`

Resources reserved for OS daemons and kernel (same keys as
kube_reserved).

- rule: system_reserved keys must be one of cpu, memory, ephemeral-storage, pid, with non-negative quantity values

### spec.metadataOptions

`KubernetesKarpenterEc2NodeClassMetadataOptions`

Instance Metadata Service exposure. CRD defaults enforce the EKS
security best practice: IMDSv2 required, hop limit 1 (pods cannot
reach the node's credentials), IPv6 endpoint disabled.

### spec.metadataOptions.httpEndpoint

`string` · optional (explicit presence)

Enable/disable the HTTP metadata endpoint. CRD default: "enabled"
("disabled" removes instance metadata entirely — most node bootstraps
need it).

- default: `enabled`
- rule: http_endpoint must be 'enabled' or 'disabled'

### spec.metadataOptions.httpProtocolIPv6

`string` · optional (explicit presence)

Enable the IPv6 metadata endpoint. CRD default: "disabled".

- default: `disabled`
- rule: http_protocol_ipv6 must be 'enabled' or 'disabled'

### spec.metadataOptions.httpPutResponseHopLimit

`int64` · optional (explicit presence)

IMDS PUT response hop limit (1-64). CRD default: 1 — pods cannot reach
the node's instance credentials (the EKS security best practice).

- default: `1`
- rule: {"int64":{"lte":"64","gte":"1"}}

### spec.metadataOptions.httpTokens

`string` · optional (explicit presence)

IMDSv2 token requirement: "required" (CRD default — v2 only) or
"optional" (v1 fallback allowed; weaker).

- default: `required`
- rule: http_tokens must be 'required' or 'optional'

### spec.networkInterfaces

`[]KubernetesKarpenterEc2NodeClassNetworkInterface`

Explicit network-interface layout (multi-card / EFA topologies). Must
include a primary 'interface' at device index 0 on card 0; at most one
efa-only interface per card. Empty = one default interface.

- rule: network_interfaces must include a primary interface with interface_type 'interface' at device_index 0 on network_card_index 0
- rule: network_interfaces must not repeat a (network_card_index, device_index) pair, and each network card may carry at most one efa-only interface
- rule: {"repeated":{"maxItems":"150"}}

### spec.networkInterfaces[].deviceIndex

`int32`

Device index for the attachment (0 = primary).

- rule: {"int32":{"gte":0}}

### spec.networkInterfaces[].interfaceType

`string` · required

Interface type: "interface" (standard ENI) or "efa-only" (Elastic
Fabric Adapter without IP networking).

- rule: interface_type must be 'interface' or 'efa-only'
- rule: {"required":true}

### spec.networkInterfaces[].networkCardIndex

`int32`

Index of the network card to attach to (multi-card instance types).

- rule: {"int32":{"gte":0}}

### spec.placementGroupSelector

`KubernetesKarpenterEc2NodeClassPlacementGroupSelector`

Placement group for provisioned instances (cluster/spread/partition
strategies), by name or id — exactly one when set.

- rule: exactly one of name or id must identify the placement group

### spec.placementGroupSelector.name

`string` · required · optional (explicit presence)

Placement-group name.

- rule: {"string":{"minLen":"1"}}

### spec.placementGroupSelector.id

`string` · optional (explicit presence)

Placement-group id (pg-...).

- rule: {"string":{"pattern":"^pg-[0-9a-z]+$"}}

### spec.tags

`map<string, string>`

Tags applied to EC2 resources Karpenter creates (instances, launch
templates, volumes). The kubernetes.io/cluster/*, karpenter.sh/
nodepool, karpenter.sh/nodeclaim, karpenter.k8s.aws/ec2nodeclass and
eks:eks-cluster-name keys are controller-owned and rejected —
mirrored from the CRD.

- rule: tags may not use controller-owned keys: eks:eks-cluster-name, kubernetes.io/cluster/*, karpenter.sh/nodepool, karpenter.sh/nodeclaim, karpenter.k8s.aws/ec2nodeclass

### spec.userData

`string`

Extra user data merged into the generated bootstrap configuration —
format follows the AMI family (MIME multi-part for AL2/AL2023, TOML
for Bottlerocket, PowerShell for Windows). Karpenter merges its own
required bootstrap fields into this.

## Validation Rules

- `spec.role_xor_instance_profile`: exactly one of role or instance_profile must be set — role lets Karpenter manage the instance profile (recommended); instance_profile brings your own
- `spec.ami_family_required_without_alias`: ami_family is required when no AMI selector term uses an alias (the family cannot be inferred)
- `spec.ami_family_matches_alias`: when an alias term is used, ami_family may only be set to the alias's own family or 'Custom'

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesKarpenterEc2NodeClass, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_class_name` | `string` | Name of the EC2NodeClass object (cluster-scoped; equals metadata.name) — the value NodePools reference through their node_class_ref. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| KubernetesKarpenterNodePool | `spec.template.nodeClassRef.name` | `status.outputs.node_class_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
