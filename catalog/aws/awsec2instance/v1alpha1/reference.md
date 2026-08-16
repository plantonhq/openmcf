# AwsEc2Instance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsEc2InstanceSpec defines a single EC2 virtual machine: the launch source
(an AMI + instance type, a launch template, or a template with inline
overrides), network placement, IAM identity, storage, instance posture
(IMDSv2, monitoring, CPU topology), purchase options (On-Demand or Spot),
and lifecycle protections.

A standalone instance is the pet of EC2 compute -- a bastion, a license
server, a singleton stateful workload -- as opposed to the cattle an
auto-scaling group manages. For fleets, compose AwsLaunchTemplate +
AwsAutoScalingGroup instead; this kind deliberately shares the launch
template's vocabulary (metadata_options, cpu_options, spot_options, block
devices) so a pet can graduate into a templated fleet without relearning
field names.

Composition over embedding: the instance ATTACHES referenced first-class
nodes -- a subnet, security groups, an IAM instance profile, a launch
template, KMS keys for volume encryption -- and creates none of them.
Access posture is composed the same way: attach an instance profile whose
role carries AmazonSSMManagedInstanceCore for keyless SSM Session Manager
access (the modern default), or set key_name for SSH against a key pair
you manage.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEc2Instance
metadata:
  name: ec2-demo
spec:
  region: us-west-2
  ami: ami-0123456789abcdef0 # replace with a current AMI in your region
  instanceType: t4g.nano
  subnetId:
    value: subnet-aaa111 # replace with your subnet id
  securityGroupIds:
    - value: sg-000111222 # replace with your security group id
  instanceProfile:
    value: ec2-demo-profile # replace with your instance profile NAME (SSM access posture)
  metadataOptions:
    httpTokens: required # enforce IMDSv2 -- the recommended hardening
    httpPutResponseHopLimit: 2
  rootBlockDevice:
    volumeSizeGb: 30
    volumeType: gp3
    encrypted: true
  volumeTags: # uniform at-creation tags on EVERY volume (ABAC-compliant); per-device tags are the post-creation alternative
    cost-center: platform
  disableApiTermination: true # protect the pet from fat-fingered deletion
  forceDestroy: true # ...while keeping teardown declarative: destroy lifts the protection itself
  userData: |
    #!/bin/bash
    echo "hello from ${HOSTNAME}" > /etc/motd
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.ami` | `string` |  |  |  |
| `spec.instanceType` | `string` |  |  |  |
| `spec.launchTemplate` | `AwsEc2InstanceLaunchTemplate` |  |  |  |
| `spec.launchTemplate.id` | `string \| valueFrom` |  |  | AwsLaunchTemplate (`status.outputs.launch_template_id`) |
| `spec.launchTemplate.name` | `string` |  |  |  |
| `spec.launchTemplate.version` | `string` |  |  |  |
| `spec.instanceProfile` | `string \| valueFrom` |  |  | AwsIamInstanceProfile (`status.outputs.instance_profile_name`) |
| `spec.keyName` | `string` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` |  |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.primaryNetworkInterfaceId` | `string` |  |  |  |
| `spec.privateIp` | `string` |  |  |  |
| `spec.secondaryPrivateIps` | `[]string` |  |  |  |
| `spec.associatePublicIpAddress` | `bool` |  |  |  |
| `spec.sourceDestCheck` | `bool` |  |  |  |
| `spec.ipv6AddressCount` | `int32` |  |  |  |
| `spec.ipv6Addresses` | `[]string` |  |  |  |
| `spec.enablePrimaryIpv6` | `bool` |  |  |  |
| `spec.privateDnsNameOptions` | `AwsEc2InstancePrivateDnsNameOptions` |  |  |  |
| `spec.privateDnsNameOptions.hostnameType` | `string` |  |  |  |
| `spec.privateDnsNameOptions.enableResourceNameDnsARecord` | `bool` |  |  |  |
| `spec.privateDnsNameOptions.enableResourceNameDnsAaaaRecord` | `bool` |  |  |  |
| `spec.secondaryNetworkInterfaces` | `[]AwsEc2InstanceSecondaryNetworkInterface` |  |  |  |
| `spec.secondaryNetworkInterfaces[].networkCardIndex` | `int32` |  |  |  |
| `spec.secondaryNetworkInterfaces[].deviceIndex` | `int32` |  |  |  |
| `spec.secondaryNetworkInterfaces[].subnetId` | `string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.secondaryNetworkInterfaces[].privateIpAddressCount` | `int32` |  |  |  |
| `spec.secondaryNetworkInterfaces[].deleteOnTermination` | `bool` |  |  |  |
| `spec.rootBlockDevice` | `AwsEc2InstanceRootBlockDevice` |  |  |  |
| `spec.rootBlockDevice.volumeSizeGb` | `int32` |  |  |  |
| `spec.rootBlockDevice.volumeType` | `string` |  |  |  |
| `spec.rootBlockDevice.iops` | `int32` |  |  |  |
| `spec.rootBlockDevice.throughputMibps` | `int32` |  |  |  |
| `spec.rootBlockDevice.encrypted` | `bool` |  |  |  |
| `spec.rootBlockDevice.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.rootBlockDevice.deleteOnTermination` | `bool` |  |  |  |
| `spec.rootBlockDevice.tags` | `map<string, string>` |  |  |  |
| `spec.ebsBlockDevices` | `[]AwsEc2InstanceEbsBlockDevice` |  |  |  |
| `spec.ebsBlockDevices[].deviceName` | `string` | yes |  |  |
| `spec.ebsBlockDevices[].volumeSizeGb` | `int32` |  |  |  |
| `spec.ebsBlockDevices[].volumeType` | `string` |  |  |  |
| `spec.ebsBlockDevices[].iops` | `int32` |  |  |  |
| `spec.ebsBlockDevices[].throughputMibps` | `int32` |  |  |  |
| `spec.ebsBlockDevices[].encrypted` | `bool` |  |  |  |
| `spec.ebsBlockDevices[].kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.ebsBlockDevices[].snapshotId` | `string` |  |  |  |
| `spec.ebsBlockDevices[].deleteOnTermination` | `bool` |  |  |  |
| `spec.ebsBlockDevices[].tags` | `map<string, string>` |  |  |  |
| `spec.ephemeralBlockDevices` | `[]AwsEc2InstanceEphemeralBlockDevice` |  |  |  |
| `spec.ephemeralBlockDevices[].deviceName` | `string` | yes |  |  |
| `spec.ephemeralBlockDevices[].virtualName` | `string` |  |  |  |
| `spec.ephemeralBlockDevices[].noDevice` | `bool` |  |  |  |
| `spec.volumeTags` | `map<string, string>` |  |  |  |
| `spec.ebsOptimized` | `bool` |  |  |  |
| `spec.metadataOptions` | `AwsEc2InstanceMetadataOptions` |  |  |  |
| `spec.metadataOptions.httpEndpoint` | `string` |  |  |  |
| `spec.metadataOptions.httpTokens` | `string` |  |  |  |
| `spec.metadataOptions.httpPutResponseHopLimit` | `int32` |  |  |  |
| `spec.metadataOptions.httpProtocolIpv6` | `string` |  |  |  |
| `spec.metadataOptions.instanceMetadataTags` | `string` |  |  |  |
| `spec.detailedMonitoring` | `bool` |  |  |  |
| `spec.cpuOptions` | `AwsEc2InstanceCpuOptions` |  |  |  |
| `spec.cpuOptions.coreCount` | `int32` |  |  |  |
| `spec.cpuOptions.threadsPerCore` | `int32` |  |  |  |
| `spec.cpuOptions.amdSevSnp` | `string` |  |  |  |
| `spec.cpuOptions.nestedVirtualization` | `string` |  |  |  |
| `spec.cpuCredits` | `string` |  |  |  |
| `spec.marketType` | `string` |  |  |  |
| `spec.spotOptions` | `AwsEc2InstanceSpotOptions` |  |  |  |
| `spec.spotOptions.maxPrice` | `string` |  |  |  |
| `spec.spotOptions.spotInstanceType` | `string` |  |  |  |
| `spec.spotOptions.instanceInterruptionBehavior` | `string` |  |  |  |
| `spec.spotOptions.validUntil` | `string` |  |  |  |
| `spec.capacityReservation` | `AwsEc2InstanceCapacityReservation` |  |  |  |
| `spec.capacityReservation.preference` | `string` |  |  |  |
| `spec.capacityReservation.capacityReservationId` | `string` |  |  |  |
| `spec.capacityReservation.capacityReservationResourceGroupArn` | `string` |  |  |  |
| `spec.placement` | `AwsEc2InstancePlacement` |  |  |  |
| `spec.placement.availabilityZone` | `string` |  |  |  |
| `spec.placement.groupName` | `string` |  |  |  |
| `spec.placement.groupId` | `string` |  |  |  |
| `spec.placement.partitionNumber` | `int32` |  |  |  |
| `spec.placement.tenancy` | `string` |  |  |  |
| `spec.placement.hostId` | `string` |  |  |  |
| `spec.placement.hostResourceGroupArn` | `string` |  |  |  |
| `spec.enclaveEnabled` | `bool` |  |  |  |
| `spec.hibernationEnabled` | `bool` |  |  |  |
| `spec.autoRecovery` | `string` |  |  |  |
| `spec.instanceInitiatedShutdownBehavior` | `string` |  |  |  |
| `spec.disableApiStop` | `bool` |  |  |  |
| `spec.disableApiTermination` | `bool` |  |  |  |
| `spec.forceDestroy` | `bool` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.userDataBase64` | `string` |  |  |  |
| `spec.userDataReplaceOnChange` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.ami

`string`

The Amazon Machine Image the instance boots from (e.g.
"ami-0abcdef1234567890"). AMI IDs are region-specific. Required unless
launch_template supplies one. ForceNew: changing the AMI replaces the
instance (see user_data_replace_on_change for the related user-data
semantics).

### spec.instanceType

`string`

The EC2 instance type (e.g. "t4g.nano", "m7g.large", "c5.xlarge")
determining vCPU count, memory, and network/EBS bandwidth. Required
unless launch_template supplies one. Changing the type stops and
restarts the instance in place (EBS-backed instances only) -- UNLESS
the old and new types share no CPU architecture (x86 -> ARM), which
replaces the instance.

### spec.launchTemplate

`AwsEc2InstanceLaunchTemplate`

Launch from an AwsLaunchTemplate instead of (or in addition to) the
inline fields. Every inline field set here OVERRIDES the template's
value for this instance -- the template is the org's golden baseline,
the inline fields are this pet's deviations.

- rule: identify the launch template by id or by name, not both
- rule: identify the launch template by id or by name

### spec.launchTemplate.id

`string | valueFrom`

The launch template ID. Reference an AwsLaunchTemplate's
launch_template_id output or pass a literal ID (e.g.
"lt-0123456789abcdef0"). Mutually exclusive with name.

- references: AwsLaunchTemplate (`status.outputs.launch_template_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLaunchTemplate, name: <that resource's name>, fieldPath: status.outputs.launch_template_id}} -- a bare string does not parse

### spec.launchTemplate.name

`string`

The launch template name, for templates managed outside the resource
graph. Mutually exclusive with id.

### spec.launchTemplate.version

`string`

The template version to launch from: a version number ("12"),
"$Latest" (track every new version -- each publish restarts the
instance), or "$Default" (AWS default when unset; both Planton launch
template modules promote each new version to default, so this tracks
the template's releases). Changing this to a version the instance is
not already running REPLACES the instance -- the provider verifies
against the live instance's actual template version.

### spec.instanceProfile

`string | valueFrom`

The IAM instance profile attached to the instance -- its identity for
SSM Session Manager, ECR pulls, S3 access, and every other AWS API call
the workload makes. Reference an AwsIamInstanceProfile's
instance_profile_name output or pass a literal profile NAME (the EC2
instance API takes the profile by name, unlike launch templates which
accept an ARN). Attachable and replaceable in place on a running
instance.

- references: AwsIamInstanceProfile (`status.outputs.instance_profile_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamInstanceProfile, name: <that resource's name>, fieldPath: status.outputs.instance_profile_name}} -- a bare string does not parse

### spec.keyName

`string`

The name of an existing EC2 key pair injected for SSH access. Leave
unset for keyless instances (SSM Session Manager via the instance
profile is the modern posture). ForceNew: changing the key pair
replaces the instance.

### spec.subnetId

`string | valueFrom`

The subnet the instance lives in. Reference an AwsSubnet's subnet_id
output or pass a literal subnet ID. Unset launches into the account's
default VPC (its default subnet in the chosen AZ) -- acceptable for
experiments, not a production posture; org accounts frequently have no
default VPC at all. ForceNew: changing the subnet replaces the
instance.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

Security groups attached to the instance's primary network interface --
what can reach the instance and what it can reach. Reference
AwsSecurityGroup security_group_id outputs or pass literal group IDs.
Unset falls back to the VPC's default security group -- not a
production posture. Updatable in place.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.primaryNetworkInterfaceId

`string`

Attach an EXISTING network interface (ENI) as the primary interface
(eth0) instead of creating one -- how an instance inherits a
pre-provisioned static identity (fixed private IP, fixed MAC, existing
security groups). The ENI carries the network configuration, so the
subnet/security-group/addressing fields here must stay unset.
ForceNew: changing the ENI replaces the instance.

### spec.privateIp

`string`

A specific primary private IPv4 address from the subnet's range.
Unset lets AWS pick one. ForceNew: changing it replaces the instance.

### spec.secondaryPrivateIps

`[]string`

Additional private IPv4 addresses on the primary interface -- for
hosting multiple TLS endpoints or failover addresses on one instance.
Updatable in place.

### spec.associatePublicIpAddress

`bool` · optional (explicit presence)

Associate a public IPv4 address at launch. Optional tri-state: unset
inherits the subnet's map-public-IP-on-launch setting; an explicit
value overrides it either way. ForceNew: changing it replaces the
instance. Note AWS now bills every public IPv4 address.

### spec.sourceDestCheck

`bool` · optional (explicit presence)

Source/destination checking on the primary interface. Optional
tri-state: unset keeps AWS's default (true -- checking on); an
explicit false is the forwarding posture for instances that carry
traffic they neither originated nor terminate (NAT instances,
software routers, VPN appliances) -- without it the network silently
drops their forwarded packets. Leave unset when attaching a
pre-provisioned primary ENI: the provider rejects the combination
(the ENI carries its own source/dest-check setting), and setting the
field on the ENI's kind is the right home for that posture.

### spec.ipv6AddressCount

`int32`

Number of IPv6 addresses AWS auto-assigns from the subnet's IPv6
range. Mutually exclusive with ipv6_addresses.

### spec.ipv6Addresses

`[]string`

Specific IPv6 addresses from the subnet's range. Mutually exclusive
with ipv6_address_count.

### spec.enablePrimaryIpv6

`bool` · optional (explicit presence)

Designate the first IPv6 address as the instance's stable primary
IPv6 address (kept until the instance or ENI is deleted) -- required
posture for IPv6-only workloads that must keep one consistent
address. One-way: disabling it after enablement forces replacement.

### spec.privateDnsNameOptions

`AwsEc2InstancePrivateDnsNameOptions`

How the guest's private DNS hostname is formed and which DNS records
resolve to it. Defaults inherit from the subnet's settings.

- rule: hostname_type must be 'ip-name' or 'resource-name' when set

### spec.privateDnsNameOptions.hostnameType

`string`

Hostname scheme: "ip-name" (ip-10-0-1-5.ec2.internal, the classic
default) or "resource-name" (i-0123....ec2.internal -- stable across
IP changes, required for IPv6-only subnets).

### spec.privateDnsNameOptions.enableResourceNameDnsARecord

`bool`

Publish an A record (IPv4) for the hostname.

### spec.privateDnsNameOptions.enableResourceNameDnsAaaaRecord

`bool`

Publish an AAAA record (IPv6) for the hostname.

### spec.secondaryNetworkInterfaces

`[]AwsEc2InstanceSecondaryNetworkInterface`

Additional network interfaces created and attached at launch on OTHER
network cards -- for high-bandwidth instance types with multiple
network cards. For ordinary multi-homing on card 0, create standalone
ENIs instead; for the primary interface, see
primary_network_interface_id.

### spec.secondaryNetworkInterfaces[].networkCardIndex

`int32`

The physical network card the interface binds to. Required; card 0
carries the primary interface, so secondary interfaces target 1+.

- rule: {"int32":{"gte":1}}

### spec.secondaryNetworkInterfaces[].deviceIndex

`int32`

Position of the interface in the attachment order on its card.
AWS default: 0.

### spec.secondaryNetworkInterfaces[].subnetId

`string | valueFrom` · required

The subnet the interface lives in. Reference an AwsSubnet's subnet_id
output or pass a literal subnet ID; must be in the same AZ as the
instance. Security groups are not configurable on secondary
interfaces at launch (the EC2 API applies the VPC default group);
manage them post-launch on the ENI itself.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.secondaryNetworkInterfaces[].privateIpAddressCount

`int32`

Number of private IPv4 addresses to assign. AWS default: 1.

### spec.secondaryNetworkInterfaces[].deleteOnTermination

`bool` · optional (explicit presence)

Delete the interface when the instance terminates. AWS default: true.
Optional so an explicit false ("keep the ENI for reuse") is
distinguishable from unset.

### spec.rootBlockDevice

`AwsEc2InstanceRootBlockDevice`

Reshape the root (boot) volume the AMI defines: grow it, switch it to
gp3, encrypt it with a customer-managed key. Unset fields inherit
from the AMI's block device mapping. Size, type, IOPS, throughput,
and delete-on-termination are updatable in place; flipping encryption
requires replacement.

- rule: volume_type must be one of: gp2, gp3, io1, io2, st1, sc1, standard
- rule: throughput_mibps must be between 125 and 2000 when set
- rule: throughput_mibps only applies to gp3 volumes
- rule: iops only applies to gp3, io1, and io2 volumes

### spec.rootBlockDevice.volumeSizeGb

`int32`

Volume size in GiB. Must be at least the AMI snapshot's size.
Growable in place; shrinking requires replacement.

### spec.rootBlockDevice.volumeType

`string`

Volume type: "gp3" (the current general-purpose default choice),
"gp2", "io1", "io2" (provisioned IOPS), "st1", "sc1"
(throughput/cold HDD), "standard" (legacy magnetic). Unset inherits
from the AMI mapping.

### spec.rootBlockDevice.iops

`int32`

Provisioned IOPS. Required for "io1"/"io2"; optional for "gp3"
(baseline 3000 without it); not valid for other types.

### spec.rootBlockDevice.throughputMibps

`int32`

Throughput in MiB/s, 125-2000. "gp3" only (baseline 125 without it;
above 1000 requires a matching iops floor per the gp3 ratio rules AWS
enforces at the API).

### spec.rootBlockDevice.encrypted

`bool`

Encrypt the root volume at rest. ForceNew on the root device: flipping
encryption replaces the instance. When the account enforces
EBS-encryption-by-default this is already true regardless.

### spec.rootBlockDevice.kmsKeyId

`string | valueFrom`

The KMS key for encryption. Reference an AwsKmsKey's key_arn output
or pass a literal key ARN. Unset with encrypted = true uses the AWS
managed aws/ebs key; a customer-managed key adds revocation and
cross-account control.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.rootBlockDevice.deleteOnTermination

`bool` · optional (explicit presence)

Delete the root volume when the instance terminates. AWS default:
true. Optional so an explicit false ("keep the boot disk for
forensics/reuse") is distinguishable from unset.

### spec.rootBlockDevice.tags

`map<string, string>`

Tags on the root volume itself. Applied AFTER instance creation by a
separate tagging call -- incompatible with ABAC/SCP policies that
require tags at creation time (use the spec-level volume_tags for
those). Mutually exclusive with volume_tags. Updatable in place.

### spec.ebsBlockDevices

`[]AwsEc2InstanceEbsBlockDevice`

Additional EBS data volumes attached at launch, keyed by device name
(e.g. "/dev/sdf"). Each volume's shape is create-time: changing a
mapping replaces the instance -- attach post-launch volumes as
separate resources when independent lifecycles matter.

- rule: volume_type must be one of: gp2, gp3, io1, io2, st1, sc1, standard
- rule: throughput_mibps must be between 125 and 2000 when set
- rule: throughput_mibps only applies to gp3 volumes
- rule: iops only applies to gp3, io1, and io2 volumes
- rule: provide volume_size_gb, or a snapshot_id that defines the size

### spec.ebsBlockDevices[].deviceName

`string` · required

The device name exposed to the instance (e.g. "/dev/sdf",
"/dev/xvdb"). Required; must not collide with the root device.

- rule: {"required":true}

### spec.ebsBlockDevices[].volumeSizeGb

`int32`

Volume size in GiB. Required unless snapshot_id supplies the size.

### spec.ebsBlockDevices[].volumeType

`string`

Volume type: "gp3", "gp2", "io1", "io2", "st1", "sc1", "standard".
AWS default: gp3 for new volumes.

### spec.ebsBlockDevices[].iops

`int32`

Provisioned IOPS. Required for "io1"/"io2"; optional for "gp3".

### spec.ebsBlockDevices[].throughputMibps

`int32`

Throughput in MiB/s, 125-2000. "gp3" only.

### spec.ebsBlockDevices[].encrypted

`bool`

Encrypt the volume at rest.

### spec.ebsBlockDevices[].kmsKeyId

`string | valueFrom`

The KMS key for encryption. Reference an AwsKmsKey's key_arn output
or pass a literal key ARN.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.ebsBlockDevices[].snapshotId

`string`

Create the volume from this EBS snapshot (e.g. a pre-warmed data
volume) instead of empty.

### spec.ebsBlockDevices[].deleteOnTermination

`bool` · optional (explicit presence)

Delete the volume when the instance terminates. AWS default: true.
Set an explicit false to keep the data volume after termination.

### spec.ebsBlockDevices[].tags

`map<string, string>`

Tags on this volume. Applied AFTER instance creation by a separate
tagging call -- incompatible with ABAC/SCP policies that require tags
at creation time (use the spec-level volume_tags for those).
Mutually exclusive with volume_tags. Updatable in place.

### spec.ephemeralBlockDevices

`[]AwsEc2InstanceEphemeralBlockDevice`

Instance-store (ephemeral local disk) mappings for instance types
with local disks. Data on instance store does not survive stop or
termination. ForceNew: changing mappings replaces the instance.

- rule: virtual_name and no_device are mutually exclusive on one mapping

### spec.ephemeralBlockDevices[].deviceName

`string` · required

The device name exposed to the instance (e.g. "/dev/sdc"). Required.

- rule: {"required":true}

### spec.ephemeralBlockDevices[].virtualName

`string`

An instance-store virtual device name ("ephemeral0", "ephemeral1",
...). Mutually exclusive with no_device.

### spec.ephemeralBlockDevices[].noDevice

`bool`

Suppress a device the AMI would otherwise attach -- the way to DROP
an AMI-baked mapping. Mutually exclusive with virtual_name.

### spec.volumeTags

`map<string, string>`

Tags applied uniformly to EVERY EBS volume at instance creation --
including volumes the AMI's block-device mapping creates that are not
declared here. Because they ride the launch call itself, these satisfy
ABAC/SCP policies that require tags at creation time. Mutually
exclusive with per-device tags (root_block_device.tags /
ebs_block_devices[].tags), which allow per-volume values but are
applied AFTER creation by a separate tagging call. Updatable in place.

### spec.ebsOptimized

`bool`

Dedicated EBS throughput between the instance and its volumes. Only
meaningful for instance types where EBS optimization is optional
(most current-generation types have it always-on at no charge);
enabling it on an unsupported type fails the launch. ForceNew.

### spec.metadataOptions

`AwsEc2InstanceMetadataOptions`

Instance Metadata Service posture. Set http_tokens = "required" to
enforce IMDSv2 -- the single most effective hardening against
credential-stealing SSRF attacks, and the recommended default for
every new instance. Updatable in place.

- rule: http_endpoint must be 'enabled' or 'disabled' when set
- rule: http_tokens must be 'required' or 'optional' when set
- rule: http_put_response_hop_limit must be between 1 and 64 when set
- rule: http_protocol_ipv6 must be 'enabled' or 'disabled' when set
- rule: instance_metadata_tags must be 'enabled' or 'disabled' when set

### spec.metadataOptions.httpEndpoint

`string`

Whether the metadata service is reachable at all: "enabled" (AWS
default) or "disabled" (nothing on the instance can fetch credentials
-- rare, for fully static workloads).

### spec.metadataOptions.httpTokens

`string`

IMDS version enforcement: "required" (IMDSv2 session tokens only --
the recommended hardening) or "optional" (v1 and v2 both answer; the
AWS default for backward compatibility).

### spec.metadataOptions.httpPutResponseHopLimit

`int32`

TTL of the IMDSv2 token-fetch packet, 1-64. 1 (AWS default) confines
metadata to the instance itself; 2 lets containerized workloads (one
extra network hop) reach it.

### spec.metadataOptions.httpProtocolIpv6

`string`

Serve metadata over the interface's IPv6 endpoint as well: "enabled"
or "disabled" (AWS default).

### spec.metadataOptions.instanceMetadataTags

`string`

Expose the instance's tags through the metadata service: "enabled" or
"disabled" (AWS default). Lets on-instance agents read tags without
ec2:DescribeTags permission.

### spec.detailedMonitoring

`bool`

Detailed CloudWatch monitoring: metrics at 1-minute granularity
instead of the free 5-minute default. Costs per instance-metric;
alarms and dashboards react meaningfully faster with it enabled.

### spec.cpuOptions

`AwsEc2InstanceCpuOptions`

CPU topology and features: trim vCPUs on license-bound workloads
(threads_per_core = 1), enable AMD SEV-SNP memory encryption, or
enable nested virtualization on supported types. ForceNew: CPU
options are fixed at launch.

- rule: threads_per_core must be 1 or 2 when set
- rule: amd_sev_snp must be 'enabled' or 'disabled' when set
- rule: nested_virtualization must be 'enabled' or 'disabled' when set

### spec.cpuOptions.coreCount

`int32`

Number of physical cores. Combined with threads_per_core = 1 this
trims the vCPU count -- the standard move for per-core licensed
software (databases) on large instance types.

### spec.cpuOptions.threadsPerCore

`int32`

Threads per core: 2 (hyper-threading, the default on x86) or 1
(disable SMT -- per-core licensing, or HPC codes that fight over
shared core resources).

### spec.cpuOptions.amdSevSnp

`string`

AMD SEV-SNP memory encryption on supported AMD types: "enabled" or
"disabled". Confidential-computing hardening.

### spec.cpuOptions.nestedVirtualization

`string`

Nested virtualization on supported 8th-generation Intel types
("enabled" or "disabled") -- run hypervisors inside the instance.
Enabling it automatically disables Virtual Secure Mode.

### spec.cpuCredits

`string`

Credit option for burstable (T-family) instance types: "standard"
(throttle when credits run out) or "unlimited" (keep bursting, pay
for the excess). AWS default: "unlimited" for recent T families.
Ignored for non-burstable types. Updatable in place.

### spec.marketType

`string`

The instance's purchase market. Unset = On-Demand (with spot_options
present implying "spot" for that classic shape). "spot" pairs with
spot_options; "capacity-block" launches into a pre-purchased ML
Capacity Block (target the block's reservation via
capacity_reservation -- required); "interruptible-capacity-reservation"
launches into an interruptible Capacity Reservation (target required).
The AwsLaunchTemplate sibling carries only spot/capacity-block -- the
interruptible market is instance-level surface at the pinned provider.
ForceNew: the purchase option is fixed at launch. Note: the provider
keeps capacity-block in state after launch (a re-plan shows no diff --
the 6.53.0 perpetual-diff fix).

### spec.spotOptions

`AwsEc2InstanceSpotOptions`

Request Spot capacity instead of On-Demand -- for interruption-
tolerant standalone workloads (a build agent, a batch box). The
instance can be reclaimed by AWS with two minutes' notice; pair with
instance_interruption_behavior "stop"/"hibernate" (persistent
requests) to survive reclaims. Presence implies market_type "spot"
when that field is unset. ForceNew: the purchase option is fixed at
launch.

- rule: spot_instance_type must be 'one-time' or 'persistent' when set
- rule: instance_interruption_behavior must be 'terminate', 'stop', or 'hibernate' when set
- rule: instance_interruption_behavior 'stop' or 'hibernate' requires spot_instance_type 'persistent'
- rule: valid_until only applies when spot_instance_type is 'persistent'

### spec.spotOptions.maxPrice

`string`

Maximum price per instance-hour, as a decimal string (e.g. "0.05").
AWS default: the On-Demand price -- and leaving it unset is the AWS
recommendation, since Spot's discount comes from interruption risk,
not bidding.

### spec.spotOptions.spotInstanceType

`string`

Spot request type: "one-time" (the AWS default -- the request ends
when the instance is interrupted) or "persistent" (EC2 re-requests
capacity after interruption; required for "stop"/"hibernate"
interruption behavior).

### spec.spotOptions.instanceInterruptionBehavior

`string`

What happens to the instance on interruption: "terminate" (default),
"stop", or "hibernate". "stop"/"hibernate" require a persistent
request.

### spec.spotOptions.validUntil

`string`

Expiry of a persistent request, RFC3339 (e.g.
"2027-01-01T00:00:00Z"). Only valid when spot_instance_type is
"persistent".

### spec.capacityReservation

`AwsEc2InstanceCapacityReservation`

Target an EC2 Capacity Reservation: "open" (use a matching
reservation if one exists -- the AWS default behavior), "none" (never
consume a reservation), or a specific reservation / reservation
group. Capacity reservations are how latency-critical or
failover-critical instances guarantee capacity in an AZ.

- rule: preference must be 'open', 'none', or 'capacity-reservations-only' when set
- rule: set preference, or target a specific reservation/group, not both
- rule: capacity_reservation_id and capacity_reservation_resource_group_arn are mutually exclusive
- rule: set preference, or target a specific reservation/group

### spec.capacityReservation.preference

`string`

Reservation preference: "open" (consume a matching reservation when
one exists -- AWS's default behavior), "none" (never consume a
reservation, even when one matches), or "capacity-reservations-only"
(launch ONLY into a matching reservation -- the launch fails when
none matches, guaranteeing reserved capacity is what runs). Mutually
exclusive with the specific-target fields.

### spec.capacityReservation.capacityReservationId

`string`

Target one specific Capacity Reservation by ID (e.g. "cr-0123...").
Mutually exclusive with preference and
capacity_reservation_resource_group_arn.

### spec.capacityReservation.capacityReservationResourceGroupArn

`string`

Target a Capacity Reservation resource group by ARN -- a pool of
reservations managed together. Mutually exclusive with preference and
capacity_reservation_id.

### spec.placement

`AwsEc2InstancePlacement`

Where the instance lands: availability zone pinning, placement group
membership, tenancy, and Dedicated Host targeting.

- rule: tenancy must be 'default', 'dedicated', or 'host' when set
- rule: group_name and group_id are mutually exclusive
- rule: host_resource_group_arn is mutually exclusive with group_name and group_id
- rule: partition_number requires a placement group (group_name or group_id)

### spec.placement.availabilityZone

`string`

Pin the instance to one availability zone (e.g. "us-west-2a"). Unset
derives from the subnet (or lets AWS choose in the default VPC).

### spec.placement.groupName

`string`

The placement group to launch into, by name: "cluster" groups pack
instances for lowest latency, "spread" and "partition" groups
separate them for fault isolation. Mutually exclusive with group_id
and host_resource_group_arn.

### spec.placement.groupId

`string`

The placement group by ID instead of name. Mutually exclusive with
group_name and host_resource_group_arn.

### spec.placement.partitionNumber

`int32`

Partition number within a partition placement group.

### spec.placement.tenancy

`string`

Instance tenancy: "default" (shared hardware), "dedicated"
(single-tenant hardware), or "host" (a Dedicated Host -- licensing
scenarios).

### spec.placement.hostId

`string`

Launch onto a specific Dedicated Host by ID (e.g. "h-0123...").

### spec.placement.hostResourceGroupArn

`string`

Launch into a host resource group (License Manager-managed Dedicated
Hosts) by ARN. Mutually exclusive with placement groups; omit
tenancy or set it to "host" alongside this.

### spec.enclaveEnabled

`bool`

Launch as an AWS Nitro Enclaves parent, enabling isolated enclave VMs
for secret-processing workloads. Not supported on every instance
type; incompatible with hibernation. ForceNew.

### spec.hibernationEnabled

`bool`

Pre-provision the instance for hibernation (requires an encrypted
root volume large enough to hold RAM contents). Incompatible with
Nitro Enclaves. ForceNew.

### spec.autoRecovery

`string`

Simplified automatic recovery on instance impairment: "default"
(recover onto healthy hardware) or "disabled" (leave the instance
impaired -- for workloads with instance-store state that recovery
would silently discard). Updatable in place.

### spec.instanceInitiatedShutdownBehavior

`string`

What an OS-initiated shutdown does to the instance: "stop" (AWS
default; the instance can be started again) or "terminate" (the
instance is gone). Updatable in place; cannot be set on
instance-store-backed instances.

### spec.disableApiStop

`bool` · optional (explicit presence)

Protect the instance from being stopped via the API -- a guard for
long-lived pets. Updatable in place.

### spec.disableApiTermination

`bool` · optional (explicit presence)

Protect the instance from being terminated via the API (termination
protection) -- the classic guard against fat-fingered deletion of a
stateful pet. Updatable in place; the module's destroy flips it off
first only if you remove the protection from the spec.

### spec.forceDestroy

`bool`

Allow destroy to proceed even while disable_api_termination /
disable_api_stop are true: the engine lifts the protections itself
before terminating, instead of failing the destroy. The declarative
escape hatch for tearing down a protected pet without first editing
its spec. Imported instances always start with this false regardless
of the live value (the provider cannot read it back).

### spec.userData

`string`

Instance user data: a cloud-init config or shell script executed on
first boot. Provide PLAIN TEXT here (16 KiB limit before encoding);
shell variable syntax like ${HOME} passes through literally. By
default, changing user data stops and restarts the instance without
replacing it -- and the NEW script does not re-run on an
already-initialized instance (cloud-init runs per-instance); set
user_data_replace_on_change to get a fresh boot. Mutually exclusive
with user_data_base64.

- rule: {"string":{"maxBytes":"16384"}}

### spec.userDataBase64

`string`

Base64-encoded user data for binary payloads (e.g. gzip-compressed
cloud-init) that cannot survive plain-text handling. Mutually
exclusive with user_data.

### spec.userDataReplaceOnChange

`bool`

Replace the instance (destroy and recreate) whenever user data
changes, instead of the default stop-update-start. The right choice
when user data IS the provisioning mechanism and a stale instance is
worse than a new one.

## Validation Rules

- `ami_or_launch_template`: provide ami, or reference a launch_template that supplies the image
- `instance_type_or_launch_template`: provide instance_type, or reference a launch_template that supplies the type
- `ami_format`: ami must be an AMI ID beginning with 'ami-'
- `primary_eni_excludes_inline_networking`: when primary_network_interface_id is set, the ENI defines networking -- clear subnet_id, security_group_ids, private_ip, secondary_private_ips, the IPv6 fields, associate_public_ip_address, and source_dest_check
- `ipv6_count_xor_addresses`: ipv6_address_count and ipv6_addresses are mutually exclusive
- `private_ip_format`: private_ip must be an IPv4 address (e.g. 10.0.1.5)
- `secondary_private_ips_format`: every secondary_private_ips entry must be an IPv4 address
- `ipv6_addresses_format`: every ipv6_addresses entry must be an IPv6 address
- `cpu_credits_valid`: cpu_credits must be 'standard' or 'unlimited' when set
- `enclave_incompatible_with_hibernation`: enclave_enabled and hibernation_enabled cannot both be true
- `auto_recovery_valid`: auto_recovery must be 'default' or 'disabled' when set
- `shutdown_behavior_valid`: instance_initiated_shutdown_behavior must be 'stop' or 'terminate' when set
- `user_data_xor_base64`: user_data and user_data_base64 are mutually exclusive
- `market_type_valid`: market_type must be 'spot', 'capacity-block', or 'interruptible-capacity-reservation' when set
- `spot_options_require_spot_market`: spot_options only applies to the 'spot' market (leave market_type unset or 'spot')
- `reservation_market_requires_target`: market_type 'capacity-block' / 'interruptible-capacity-reservation' requires capacity_reservation with capacity_reservation_id or capacity_reservation_resource_group_arn
- `volume_tags_xor_per_device_tags`: volume_tags and per-device tags (root_block_device.tags / ebs_block_devices[].tags) are mutually exclusive

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEc2Instance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The instance ID (e.g. "i-0123456789abcdef0"). The primary handle: what load-balancer target groups register (an AwsLbTargetGroup instance target references this output), what the CLI and APIs address. |
| `status.outputs.arn` | `string` | The ARN of the instance, for IAM policies and EventBridge rules scoped to this instance. |
| `status.outputs.instance_state` | `string` | The instance lifecycle state as of the last deploy ("running", "stopped", ...). |
| `status.outputs.availability_zone` | `string` | The availability zone the instance runs in (e.g. "us-west-2a"). |
| `status.outputs.private_ip` | `string` | The primary private IPv4 address. |
| `status.outputs.private_dns` | `string` | The private DNS hostname within the VPC (e.g. "ip-10-0-1-5.us-west-2.compute.internal"). |
| `status.outputs.public_ip` | `string` | The public IPv4 address, when one is associated. Empty for private-only instances. Note this address changes across stop/start cycles -- compose an AwsElasticIp for a stable public address. |
| `status.outputs.public_dns` | `string` | The public DNS hostname, when a public address is associated. Empty for private-only instances. |
| `status.outputs.primary_network_interface_id` | `string` | The ID of the primary network interface (eth0) -- the attachment point for Elastic IP associations and flow-log scoping. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.launchTemplate.id` | AwsLaunchTemplate | `status.outputs.launch_template_id` |
| `spec.instanceProfile` | AwsIamInstanceProfile | `status.outputs.instance_profile_name` |
| `spec.subnetId` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.secondaryNetworkInterfaces[].subnetId` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.rootBlockDevice.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.ebsBlockDevices[].kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsBudget | `spec.actions[].ssmActionDefinition.instanceIds` | `status.outputs.instance_id` |
| AwsElasticIp | `spec.instance` | `status.outputs.instance_id` |
| AwsElasticIp | `spec.networkInterface` | `status.outputs.primary_network_interface_id` |
| AwsLbTargetGroup | `spec.targets[].targetId` | `status.outputs.instance_id` |

## See Also

- [Overview](../README.md)
