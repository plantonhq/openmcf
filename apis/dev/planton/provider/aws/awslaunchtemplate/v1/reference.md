# AwsLaunchTemplate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsLaunchTemplateSpec defines an EC2 launch template: the reusable
blueprint that describes how to launch an instance -- AMI, instance type
(or attribute-based requirements), storage, networking, IAM identity,
metadata-service posture, and purchase options.

A launch template is the composition anchor of EC2 fleet compute. It has
its own lifecycle and is referenced from many places: auto-scaling groups
(directly or through mixed-instances overrides), EKS managed node groups,
AWS Batch compute environments, EC2 Fleet, and one-off RunInstances calls.
That is why it is a first-class resource rather than a detail of the
group that launches from it.

Launch templates are versioned and versions are immutable in AWS: every
change to this spec produces a NEW template version rather than mutating
the old one, and both IaC modules promote that new version to the
template's default -- so consumers that follow "$Default" (the common ASG
and node-group setup) pick up the change on their next launch or instance
refresh, while consumers pinned to a numeric version keep exactly what
they tested. Only the template name is create-only; the name comes from
metadata.name (AWS limit: 125 characters; both modules truncate longer
names deterministically).

Everything here is a launch-time DEFAULT, not a constraint: a consumer
(an ASG override, a RunInstances call) may override any of it. Leaving a
field unset omits it from the template so the consumer or the AWS
account-level default decides -- which is exactly how an org-wide
"golden" template is built: set the opinionated parts (IMDSv2, encrypted
volumes, monitoring), leave the workload-specific parts open.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLaunchTemplate
metadata:
  name: web-demo
spec:
  region: us-west-2
  description: hardened web fleet blueprint
  imageId: ami-0123456789abcdef0
  instanceType: t3.small
  instanceProfile:
    value: arn:aws:iam::123456789012:instance-profile/web
  securityGroupIds:
    - value: sg-0123456789abcdef0
  # The hardened-fleet shape: IMDSv2 enforced, encrypted gp3 root volume,
  # detailed monitoring, and a user-data bootstrap. Exercises the nested
  # metadata-options and block-device objects so the fixture proves the full
  # variable contract, not just the scalars.
  detailedMonitoring: true
  ebsOptimized: true
  userData: |
    #!/bin/bash
    echo "bootstrapping" > /var/log/bootstrap.log
  metadataOptions:
    httpEndpoint: enabled
    httpTokens: required
    httpPutResponseHopLimit: 2
    instanceMetadataTags: enabled
  blockDeviceMappings:
    - deviceName: /dev/xvda
      ebs:
        volumeSizeGb: 50
        volumeType: gp3
        iops: 4000
        throughputMibps: 250
        encrypted: true
        deleteOnTermination: true
  instanceInitiatedShutdownBehavior: terminate
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.imageId` | `string` |  |  |  |
| `spec.instanceType` | `string` |  |  |  |
| `spec.instanceRequirements` | `AwsLaunchTemplateInstanceRequirements` |  |  |  |
| `spec.instanceRequirements.memoryMib` | `AwsLaunchTemplateIntRange` | yes |  |  |
| `spec.instanceRequirements.memoryMib.min` | `int32` |  |  |  |
| `spec.instanceRequirements.memoryMib.max` | `int32` |  |  |  |
| `spec.instanceRequirements.vcpuCount` | `AwsLaunchTemplateIntRange` | yes |  |  |
| `spec.instanceRequirements.vcpuCount.min` | `int32` |  |  |  |
| `spec.instanceRequirements.vcpuCount.max` | `int32` |  |  |  |
| `spec.instanceRequirements.allowedInstanceTypes` | `[]string` |  |  |  |
| `spec.instanceRequirements.excludedInstanceTypes` | `[]string` |  |  |  |
| `spec.instanceRequirements.instanceGenerations` | `[]string` |  |  |  |
| `spec.instanceRequirements.cpuManufacturers` | `[]string` |  |  |  |
| `spec.instanceRequirements.bareMetal` | `string` |  |  |  |
| `spec.instanceRequirements.burstablePerformance` | `string` |  |  |  |
| `spec.instanceRequirements.requireHibernateSupport` | `bool` |  |  |  |
| `spec.instanceRequirements.spotMaxPricePercentageOverLowestPrice` | `int32` |  |  |  |
| `spec.instanceRequirements.maxSpotPriceAsPercentageOfOptimalOnDemandPrice` | `int32` |  |  |  |
| `spec.instanceRequirements.onDemandMaxPricePercentageOverLowestPrice` | `int32` |  |  |  |
| `spec.instanceRequirements.localStorage` | `string` |  |  |  |
| `spec.instanceRequirements.localStorageTypes` | `[]string` |  |  |  |
| `spec.instanceRequirements.totalLocalStorageGb` | `AwsLaunchTemplateDoubleRange` |  |  |  |
| `spec.instanceRequirements.totalLocalStorageGb.min` | `double` |  |  |  |
| `spec.instanceRequirements.totalLocalStorageGb.max` | `double` |  |  |  |
| `spec.instanceRequirements.memoryGibPerVcpu` | `AwsLaunchTemplateDoubleRange` |  |  |  |
| `spec.instanceRequirements.memoryGibPerVcpu.min` | `double` |  |  |  |
| `spec.instanceRequirements.memoryGibPerVcpu.max` | `double` |  |  |  |
| `spec.instanceRequirements.networkInterfaceCount` | `AwsLaunchTemplateIntRange` |  |  |  |
| `spec.instanceRequirements.networkInterfaceCount.min` | `int32` |  |  |  |
| `spec.instanceRequirements.networkInterfaceCount.max` | `int32` |  |  |  |
| `spec.instanceRequirements.networkBandwidthGbps` | `AwsLaunchTemplateDoubleRange` |  |  |  |
| `spec.instanceRequirements.networkBandwidthGbps.min` | `double` |  |  |  |
| `spec.instanceRequirements.networkBandwidthGbps.max` | `double` |  |  |  |
| `spec.instanceRequirements.baselineEbsBandwidthMbps` | `AwsLaunchTemplateIntRange` |  |  |  |
| `spec.instanceRequirements.baselineEbsBandwidthMbps.min` | `int32` |  |  |  |
| `spec.instanceRequirements.baselineEbsBandwidthMbps.max` | `int32` |  |  |  |
| `spec.instanceRequirements.acceleratorCount` | `AwsLaunchTemplateIntRange` |  |  |  |
| `spec.instanceRequirements.acceleratorCount.min` | `int32` |  |  |  |
| `spec.instanceRequirements.acceleratorCount.max` | `int32` |  |  |  |
| `spec.instanceRequirements.acceleratorManufacturers` | `[]string` |  |  |  |
| `spec.instanceRequirements.acceleratorNames` | `[]string` |  |  |  |
| `spec.instanceRequirements.acceleratorTypes` | `[]string` |  |  |  |
| `spec.instanceRequirements.acceleratorTotalMemoryMib` | `AwsLaunchTemplateIntRange` |  |  |  |
| `spec.instanceRequirements.acceleratorTotalMemoryMib.min` | `int32` |  |  |  |
| `spec.instanceRequirements.acceleratorTotalMemoryMib.max` | `int32` |  |  |  |
| `spec.keyName` | `string` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.instanceProfile` | `string \| valueFrom` |  |  | AwsIamInstanceProfile (`status.outputs.instance_profile_arn`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.ebsOptimized` | `bool` |  |  |  |
| `spec.blockDeviceMappings` | `[]AwsLaunchTemplateBlockDeviceMapping` |  |  |  |
| `spec.blockDeviceMappings[].deviceName` | `string` | yes |  |  |
| `spec.blockDeviceMappings[].virtualName` | `string` |  |  |  |
| `spec.blockDeviceMappings[].noDevice` | `bool` |  |  |  |
| `spec.blockDeviceMappings[].ebs` | `AwsLaunchTemplateEbs` |  |  |  |
| `spec.blockDeviceMappings[].ebs.volumeSizeGb` | `int32` |  |  |  |
| `spec.blockDeviceMappings[].ebs.volumeType` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs.iops` | `int32` |  |  |  |
| `spec.blockDeviceMappings[].ebs.throughputMibps` | `int32` |  |  |  |
| `spec.blockDeviceMappings[].ebs.encrypted` | `bool` |  |  |  |
| `spec.blockDeviceMappings[].ebs.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.blockDeviceMappings[].ebs.snapshotId` | `string` |  |  |  |
| `spec.blockDeviceMappings[].ebs.deleteOnTermination` | `bool` |  |  |  |
| `spec.networkInterfaces` | `[]AwsLaunchTemplateNetworkInterface` |  |  |  |
| `spec.networkInterfaces[].deviceIndex` | `int32` |  |  |  |
| `spec.networkInterfaces[].networkCardIndex` | `int32` |  |  |  |
| `spec.networkInterfaces[].description` | `string` |  |  |  |
| `spec.networkInterfaces[].interfaceType` | `string` |  |  |  |
| `spec.networkInterfaces[].networkInterfaceId` | `string` |  |  |  |
| `spec.networkInterfaces[].associatePublicIpAddress` | `bool` |  |  |  |
| `spec.networkInterfaces[].deleteOnTermination` | `bool` |  |  |  |
| `spec.networkInterfaces[].subnetId` | `string \| valueFrom` |  |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.networkInterfaces[].securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.networkInterfaces[].privateIpAddress` | `string` |  |  |  |
| `spec.networkInterfaces[].ipv4AddressCount` | `int32` |  |  |  |
| `spec.networkInterfaces[].ipv4Addresses` | `[]string` |  |  |  |
| `spec.networkInterfaces[].ipv6AddressCount` | `int32` |  |  |  |
| `spec.networkInterfaces[].ipv6Addresses` | `[]string` |  |  |  |
| `spec.networkInterfaces[].ipv4PrefixCount` | `int32` |  |  |  |
| `spec.networkInterfaces[].ipv4Prefixes` | `[]string` |  |  |  |
| `spec.networkInterfaces[].ipv6PrefixCount` | `int32` |  |  |  |
| `spec.networkInterfaces[].ipv6Prefixes` | `[]string` |  |  |  |
| `spec.metadataOptions` | `AwsLaunchTemplateMetadataOptions` |  |  |  |
| `spec.metadataOptions.httpEndpoint` | `string` |  |  |  |
| `spec.metadataOptions.httpTokens` | `string` |  |  |  |
| `spec.metadataOptions.httpPutResponseHopLimit` | `int32` |  |  |  |
| `spec.metadataOptions.httpProtocolIpv6` | `string` |  |  |  |
| `spec.metadataOptions.instanceMetadataTags` | `string` |  |  |  |
| `spec.detailedMonitoring` | `bool` |  |  |  |
| `spec.placement` | `AwsLaunchTemplatePlacement` |  |  |  |
| `spec.placement.availabilityZone` | `string` |  |  |  |
| `spec.placement.groupName` | `string` |  |  |  |
| `spec.placement.partitionNumber` | `int32` |  |  |  |
| `spec.placement.tenancy` | `string` |  |  |  |
| `spec.cpuOptions` | `AwsLaunchTemplateCpuOptions` |  |  |  |
| `spec.cpuOptions.coreCount` | `int32` |  |  |  |
| `spec.cpuOptions.threadsPerCore` | `int32` |  |  |  |
| `spec.cpuOptions.amdSevSnp` | `string` |  |  |  |
| `spec.cpuCredits` | `string` |  |  |  |
| `spec.spotOptions` | `AwsLaunchTemplateSpotOptions` |  |  |  |
| `spec.spotOptions.maxPrice` | `string` |  |  |  |
| `spec.spotOptions.spotInstanceType` | `string` |  |  |  |
| `spec.spotOptions.instanceInterruptionBehavior` | `string` |  |  |  |
| `spec.spotOptions.validUntil` | `string` |  |  |  |
| `spec.enclaveEnabled` | `bool` |  |  |  |
| `spec.hibernationEnabled` | `bool` |  |  |  |
| `spec.autoRecovery` | `string` |  |  |  |
| `spec.privateDnsNameOptions` | `AwsLaunchTemplatePrivateDnsNameOptions` |  |  |  |
| `spec.privateDnsNameOptions.hostnameType` | `string` |  |  |  |
| `spec.privateDnsNameOptions.enableResourceNameDnsARecord` | `bool` |  |  |  |
| `spec.privateDnsNameOptions.enableResourceNameDnsAaaaRecord` | `bool` |  |  |  |
| `spec.disableApiStop` | `bool` |  |  |  |
| `spec.disableApiTermination` | `bool` |  |  |  |
| `spec.instanceInitiatedShutdownBehavior` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the launch template is created in. A template is a
regional object: an auto-scaling group or node group can only launch
from a template in its own region.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

Human-readable description recorded on each template VERSION (AWS calls
this the version description, up to 255 characters). Use it as a
change-log line: "amzn2023 + IMDSv2 + gp3", "rotate AMI 2026-07".

- rule: {"string":{"maxLen":"255"}}

### spec.imageId

`string`

The Amazon Machine Image the instance boots from (e.g.
"ami-0abcdef1234567890"). Optional at the template level: a template
without an AMI is a partial blueprint whose consumer must supply one --
but note that auto-scaling groups require the template they reference
to carry an AMI, so leave this unset only for consumers that inject
their own image (EKS managed node groups, EC2 Fleet overrides).

### spec.instanceType

`string`

The EC2 instance type launched by default (e.g. "t3.small",
"m7g.large"). Mutually exclusive with instance_requirements -- name an
exact type here, or describe what you need there and let AWS pick
matching types. Leave both unset when the consumer supplies the type
(an ASG mixed-instances override, an EKS node group).

### spec.instanceRequirements

`AwsLaunchTemplateInstanceRequirements`

Attribute-based instance selection: instead of naming a type, describe
the compute you need (vCPU and memory ranges, CPU generations,
accelerators, price protection) and AWS resolves the matching set of
instance types at launch. The foundation of Spot diversification -- a
fleet that can draw from dozens of pools rides out capacity events that
would starve a single-type group. Mutually exclusive with
instance_type.

- rule: allowed_instance_types and excluded_instance_types are mutually exclusive
- rule: spot_max_price_percentage_over_lowest_price and max_spot_price_as_percentage_of_optimal_on_demand_price are mutually exclusive
- rule: bare_metal must be 'included', 'excluded', or 'required' when set
- rule: burstable_performance must be 'included', 'excluded', or 'required' when set
- rule: local_storage must be 'included', 'excluded', or 'required' when set

### spec.instanceRequirements.memoryMib

`AwsLaunchTemplateIntRange` · required

Required. Memory per instance, in MiB. min is required; leave max unset
(0) for no upper bound.

- rule: {"required":true}
- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.memoryMib.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.memoryMib.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.vcpuCount

`AwsLaunchTemplateIntRange` · required

Required. vCPUs per instance. min is required; leave max unset (0) for
no upper bound.

- rule: {"required":true}
- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.vcpuCount.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.vcpuCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.allowedInstanceTypes

`[]string`

Allow-list of instance types or families, with wildcards ("m5.large",
"m5.*", "c*"). At most 400 entries. Mutually exclusive with
excluded_instance_types.

- rule: {"repeated":{"maxItems":"400"}}

### spec.instanceRequirements.excludedInstanceTypes

`[]string`

Deny-list of instance types or families, with wildcards. At most 400
entries. Mutually exclusive with allowed_instance_types.

- rule: {"repeated":{"maxItems":"400"}}

### spec.instanceRequirements.instanceGenerations

`[]string`

Instance generations to include: "current" and/or "previous". AWS
default: any generation matching the other requirements.

### spec.instanceRequirements.cpuManufacturers

`[]string`

CPU manufacturers to include: "intel", "amd", "amazon-web-services"
(Graviton), "apple". AWS default: any. Selecting only
"amazon-web-services" is how an arm64 fleet is expressed -- pair it
with an arm64 AMI.

### spec.instanceRequirements.bareMetal

`string`

Bare-metal eligibility: "included", "excluded" (AWS default), or
"required".

### spec.instanceRequirements.burstablePerformance

`string`

Burstable (T-family) eligibility: "included", "excluded" (AWS
default), or "required".

### spec.instanceRequirements.requireHibernateSupport

`bool`

Only instance types that support hibernation.

### spec.instanceRequirements.spotMaxPricePercentageOverLowestPrice

`int32`

Spot price protection: exclude types whose Spot price exceeds the
identified lowest-priced type's Spot price by more than this
percentage. Mutually exclusive with
max_spot_price_as_percentage_of_optimal_on_demand_price.

### spec.instanceRequirements.maxSpotPriceAsPercentageOfOptimalOnDemandPrice

`int32`

Spot price protection anchored to On-Demand: exclude types whose Spot
price exceeds this percentage of the optimal type's On-Demand price --
steadier than the lowest-Spot anchor because On-Demand prices do not
fluctuate. Mutually exclusive with
spot_max_price_percentage_over_lowest_price.

### spec.instanceRequirements.onDemandMaxPricePercentageOverLowestPrice

`int32`

On-Demand price protection: exclude types whose On-Demand price exceeds
the identified lowest-priced type's by more than this percentage. AWS
default: 20.

### spec.instanceRequirements.localStorage

`string`

Instance-store (local disk) eligibility: "included" (AWS default),
"excluded", or "required".

### spec.instanceRequirements.localStorageTypes

`[]string`

Local storage technologies when instance-store is in play: "hdd"
and/or "ssd".

### spec.instanceRequirements.totalLocalStorageGb

`AwsLaunchTemplateDoubleRange`

Total local (instance-store) storage, in GB.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.totalLocalStorageGb.min

`double`

Lower bound, inclusive.

### spec.instanceRequirements.totalLocalStorageGb.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.memoryGibPerVcpu

`AwsLaunchTemplateDoubleRange`

Memory-to-vCPU ratio, in GiB per vCPU -- a compact way to say "memory
optimized" (min 8) or "compute optimized" (max 2) without naming
families.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.memoryGibPerVcpu.min

`double`

Lower bound, inclusive.

### spec.instanceRequirements.memoryGibPerVcpu.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.networkInterfaceCount

`AwsLaunchTemplateIntRange`

Number of network interfaces the type must support.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.networkInterfaceCount.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.networkInterfaceCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.networkBandwidthGbps

`AwsLaunchTemplateDoubleRange`

Network bandwidth, in Gbps.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.networkBandwidthGbps.min

`double`

Lower bound, inclusive.

### spec.instanceRequirements.networkBandwidthGbps.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.baselineEbsBandwidthMbps

`AwsLaunchTemplateIntRange`

Baseline (non-burst) EBS bandwidth, in Mbps.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.baselineEbsBandwidthMbps.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.baselineEbsBandwidthMbps.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.acceleratorCount

`AwsLaunchTemplateIntRange`

Number of accelerators (GPUs, FPGAs, inference chips). Set min 1 to
require accelerated types; set max 0 explicitly via {min:0, max:0} is
not expressible -- to EXCLUDE accelerators, leave this unset and rely
on accelerator_types being empty.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.acceleratorCount.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.acceleratorCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.instanceRequirements.acceleratorManufacturers

`[]string`

Accelerator manufacturers: "nvidia", "amd", "amazon-web-services",
"xilinx", "habana".

### spec.instanceRequirements.acceleratorNames

`[]string`

Specific accelerator models (e.g. "a100", "v100", "t4",
"inferentia", "radeon-pro-v520").

### spec.instanceRequirements.acceleratorTypes

`[]string`

Accelerator categories: "gpu", "fpga", "inference".

### spec.instanceRequirements.acceleratorTotalMemoryMib

`AwsLaunchTemplateIntRange`

Total accelerator memory, in MiB.

- rule: max must be greater than or equal to min when both are set

### spec.instanceRequirements.acceleratorTotalMemoryMib.min

`int32`

Lower bound, inclusive.

### spec.instanceRequirements.acceleratorTotalMemoryMib.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.keyName

`string`

The name of an existing EC2 key pair injected for SSH access. Leave
unset for keyless fleets (SSM Session Manager via the instance
profile is the modern posture).

### spec.userData

`string`

Instance user data: a cloud-init config or shell script executed on
first boot. Provide PLAIN TEXT here -- both IaC modules base64-encode
it for the EC2 API, so the manifest stays readable. AWS limit: 16 KiB
before encoding.

- rule: {"string":{"maxBytes":"16384"}}

### spec.instanceProfile

`string | valueFrom`

The IAM instance profile attached to launched instances -- the
instance's identity for SSM access, ECR pulls, S3 access, and every
other AWS API call the workload makes. Reference an
AwsIamInstanceProfile's instance_profile_arn output or pass a literal
profile ARN.

- references: AwsIamInstanceProfile (`status.outputs.instance_profile_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamInstanceProfile, name: <that resource's name>, fieldPath: status.outputs.instance_profile_arn}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

Security groups attached to the instance's primary network interface.
Mutually exclusive with per-interface security groups inside
network_interfaces -- when the template declares explicit interfaces,
attach security groups on the interface instead.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.ebsOptimized

`bool`

Dedicated EBS throughput between the instance and its volumes. Only
meaningful for instance types where EBS optimization is optional (most
current-generation types have it always-on at no charge); enabling it
on an unsupported type fails the launch.

### spec.blockDeviceMappings

`[]AwsLaunchTemplateBlockDeviceMapping`

Block device mappings: the volumes attached at launch, keyed by device
name. Override the AMI's root volume (grow it, switch to gp3, encrypt
with a CMK) or attach additional data volumes.

- rule: ebs and virtual_name are mutually exclusive on one mapping
- rule: ebs and no_device are mutually exclusive on one mapping

### spec.blockDeviceMappings[].deviceName

`string` · required

The device name exposed to the instance (e.g. "/dev/xvda",
"/dev/sdf"). Required.

- rule: {"required":true}

### spec.blockDeviceMappings[].virtualName

`string`

An instance-store virtual device name ("ephemeral0", "ephemeral1", ...)
for instance types with local disks. Mutually exclusive with ebs.

### spec.blockDeviceMappings[].noDevice

`bool`

Suppress a device the AMI would otherwise attach -- the way to DROP an
AMI-baked data volume. Mutually exclusive with ebs.

### spec.blockDeviceMappings[].ebs

`AwsLaunchTemplateEbs`

EBS volume configuration for this device.

- rule: volume_type must be one of: gp2, gp3, io1, io2, st1, sc1, standard
- rule: throughput_mibps must be between 125 and 1000 when set
- rule: throughput_mibps only applies to gp3 volumes
- rule: iops only applies to gp3, io1, and io2 volumes

### spec.blockDeviceMappings[].ebs.volumeSizeGb

`int32`

Volume size in GiB. Must be at least the AMI snapshot's size when
overriding the root device.

### spec.blockDeviceMappings[].ebs.volumeType

`string`

Volume type: "gp3" (the current general-purpose default choice),
"gp2", "io1", "io2" (provisioned IOPS), "st1", "sc1" (throughput/cold
HDD), "standard" (legacy magnetic). Unset inherits from the AMI
mapping.

### spec.blockDeviceMappings[].ebs.iops

`int32`

Provisioned IOPS. Required for "io1"/"io2"; optional for "gp3"
(baseline 3000 without it); not valid for other types.

### spec.blockDeviceMappings[].ebs.throughputMibps

`int32`

Throughput in MiB/s, 125-1000. "gp3" only (baseline 125 without it).

### spec.blockDeviceMappings[].ebs.encrypted

`bool`

Encrypt the volume at rest. Snapshots and volumes created from them
stay encrypted. When the account enforces EBS-encryption-by-default
this is already true regardless.

### spec.blockDeviceMappings[].ebs.kmsKeyId

`string | valueFrom`

The KMS key for encryption. Reference an AwsKmsKey's key_arn output or
pass a literal key ARN. Unset with encrypted = true uses the AWS
managed aws/ebs key; a customer-managed key adds revocation and
cross-account control.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.blockDeviceMappings[].ebs.snapshotId

`string`

Create the volume from this EBS snapshot (e.g. a pre-warmed data
volume).

### spec.blockDeviceMappings[].ebs.deleteOnTermination

`bool` · optional (explicit presence)

Delete the volume when the instance terminates. AWS default: true for
the root volume, false for additional volumes (inherited from the
AMI). Optional rather than plain bool so an explicit false ("keep the
data volume") is distinguishable from unset ("keep the AMI default").

### spec.networkInterfaces

`[]AwsLaunchTemplateNetworkInterface`

Explicit network interfaces. Most templates leave this empty and let
the consumer place the instance (an ASG spreads across its subnets);
declare interfaces to control public-IP association, static private
IPs/prefixes, multiple NICs, or EFA for HPC/ML workloads. When any
interface is declared, put security groups on the interface, not on
security_group_ids.

- rule: interface_type must be 'interface', 'efa', or 'efa-only' when set
- rule: ipv4_address_count and ipv4_addresses are mutually exclusive
- rule: ipv6_address_count and ipv6_addresses are mutually exclusive
- rule: ipv4_prefix_count and ipv4_prefixes are mutually exclusive
- rule: ipv6_prefix_count and ipv6_prefixes are mutually exclusive

### spec.networkInterfaces[].deviceIndex

`int32`

Position of the interface in the attachment order. The primary
interface is 0.

### spec.networkInterfaces[].networkCardIndex

`int32`

The physical network card the interface binds to, for instance types
with multiple cards (high-bandwidth/EFA types). Default 0.

### spec.networkInterfaces[].description

`string`

Free-text description of the interface.

### spec.networkInterfaces[].interfaceType

`string`

Interface type: "interface" (standard, the default), "efa" (Elastic
Fabric Adapter with OS-bypass for tightly coupled HPC/ML), or
"efa-only" (EFA without an IP -- secondary cards on multi-card
types).

### spec.networkInterfaces[].networkInterfaceId

`string`

Attach an existing ENI by ID instead of creating one -- for a
pre-provisioned static identity (fixed IP/MAC). Mutually exclusive
with subnet placement and addressing fields.

### spec.networkInterfaces[].associatePublicIpAddress

`bool` · optional (explicit presence)

Associate a public IPv4 address. Optional tri-state: unset inherits
the subnet's map-public-IP setting; an explicit value overrides it
either way.

### spec.networkInterfaces[].deleteOnTermination

`bool` · optional (explicit presence)

Delete the interface when the instance terminates. AWS default: true
for interfaces the launch creates. Optional so an explicit false
("keep the ENI for reuse") is distinguishable from unset.

### spec.networkInterfaces[].subnetId

`string | valueFrom`

The subnet the interface lives in. Reference an AwsSubnet's subnet_id
output or pass a literal subnet ID. Setting it pins every launch to
this subnet; auto-scaling templates normally leave it unset.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.networkInterfaces[].securityGroupIds

`[]string | valueFrom`

Security groups on this interface. Reference AwsSecurityGroup
security_group_id outputs or pass literal group IDs.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.networkInterfaces[].privateIpAddress

`string`

A specific primary private IPv4 address from the subnet's range.

### spec.networkInterfaces[].ipv4AddressCount

`int32`

Number of additional private IPv4 addresses AWS auto-assigns.
Mutually exclusive with ipv4_addresses.

### spec.networkInterfaces[].ipv4Addresses

`[]string`

Specific secondary private IPv4 addresses. Mutually exclusive with
ipv4_address_count.

### spec.networkInterfaces[].ipv6AddressCount

`int32`

Number of IPv6 addresses AWS auto-assigns from the subnet's IPv6
range. Mutually exclusive with ipv6_addresses.

### spec.networkInterfaces[].ipv6Addresses

`[]string`

Specific IPv6 addresses. Mutually exclusive with ipv6_address_count.

### spec.networkInterfaces[].ipv4PrefixCount

`int32`

Number of /28 IPv4 prefixes AWS auto-assigns (prefix delegation --
how Kubernetes CNIs scale pod IPs per node). Mutually exclusive with
ipv4_prefixes.

### spec.networkInterfaces[].ipv4Prefixes

`[]string`

Specific /28 IPv4 prefixes. Mutually exclusive with
ipv4_prefix_count.

### spec.networkInterfaces[].ipv6PrefixCount

`int32`

Number of /80 IPv6 prefixes AWS auto-assigns. Mutually exclusive with
ipv6_prefixes.

### spec.networkInterfaces[].ipv6Prefixes

`[]string`

Specific /80 IPv6 prefixes. Mutually exclusive with
ipv6_prefix_count.

### spec.metadataOptions

`AwsLaunchTemplateMetadataOptions`

Instance Metadata Service posture. Set http_tokens = "required" to
enforce IMDSv2 -- the single most effective hardening against
credential-stealing SSRF attacks, and the recommended default for every
new template.

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

Detailed CloudWatch monitoring: metrics at 1-minute granularity instead
of the free 5-minute default. Costs per instance-metric; auto-scaling
policies react meaningfully faster with it enabled.

### spec.placement

`AwsLaunchTemplatePlacement`

Placement of launched instances: availability zone pinning, placement
group membership, and tenancy.

- rule: tenancy must be 'default', 'dedicated', or 'host' when set

### spec.placement.availabilityZone

`string`

Pin launches to one availability zone (e.g. "us-west-2a"). For ASG
templates, prefer the group's subnets over AZ pinning here.

### spec.placement.groupName

`string`

The placement group to launch into: "cluster" groups pack instances
for lowest latency, "spread" and "partition" groups separate them for
fault isolation.

### spec.placement.partitionNumber

`int32`

Partition number within a partition placement group.

### spec.placement.tenancy

`string`

Instance tenancy: "default" (shared hardware), "dedicated"
(single-tenant hardware), or "host" (a specific Dedicated Host --
licensing scenarios).

### spec.cpuOptions

`AwsLaunchTemplateCpuOptions`

CPU topology and features: trim vCPUs on license-bound workloads
(threads_per_core = 1), or enable AMD SEV-SNP memory encryption.

- rule: threads_per_core must be 1 or 2 when set
- rule: amd_sev_snp must be 'enabled' or 'disabled' when set

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

### spec.cpuCredits

`string`

Credit option for burstable (T-family) instance types: "standard"
(throttle when credits run out) or "unlimited" (keep bursting, pay for
the excess). AWS default: "unlimited" for recent T families. Ignored
for non-burstable types.

### spec.spotOptions

`AwsLaunchTemplateSpotOptions`

Request Spot capacity instead of On-Demand. Configuring this block
makes every launch from the template a Spot request -- appropriate for
interruption-tolerant workloads. For a mixed On-Demand/Spot fleet,
leave this unset and blend purchase options in the auto-scaling group's
mixed-instances policy instead.

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

Spot request type: "one-time" (the default for templates consumed by
auto-scaling -- the ASG replaces interrupted capacity itself) or
"persistent" (EC2 re-requests the instance after interruption;
standalone instances only).

### spec.spotOptions.instanceInterruptionBehavior

`string`

What happens to the instance on interruption: "terminate" (default),
"stop", or "hibernate". "stop"/"hibernate" require a persistent
request.

### spec.spotOptions.validUntil

`string`

Expiry of a persistent request, RFC3339 (e.g. "2027-01-01T00:00:00Z").
Only valid when spot_instance_type is "persistent".

### spec.enclaveEnabled

`bool`

Launch instances as AWS Nitro Enclaves parents, enabling isolated
enclave VMs for secret-processing workloads. Not supported on every
instance type; incompatible with hibernation.

### spec.hibernationEnabled

`bool`

Pre-provision instances for hibernation (encrypted root volume large
enough to hold RAM contents required). Incompatible with Nitro
Enclaves.

### spec.autoRecovery

`string`

Simplified automatic recovery on instance impairment: "default"
(recover onto healthy hardware) or "disabled" (leave the instance
impaired -- for workloads with instance-store state that recovery would
silently discard).

### spec.privateDnsNameOptions

`AwsLaunchTemplatePrivateDnsNameOptions`

How the guest's private DNS hostname is formed and which DNS records
resolve to it.

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

### spec.disableApiStop

`bool`

Protect launched instances from being stopped via the API. A guard for
pet-like fleet members; leave false for disposable fleet instances.

### spec.disableApiTermination

`bool`

Protect launched instances from being terminated via the API. NOTE: an
auto-scaling group cannot scale in instances with termination
protection -- use the ASG's protect_from_scale_in instead for fleet
instances.

### spec.instanceInitiatedShutdownBehavior

`string`

What an OS-initiated shutdown does to the instance: "stop" (AWS
default; the instance can be started again) or "terminate" (the
instance is gone -- the right choice for immutable fleet members that
should never be resurrected by hand).

## Validation Rules

- `instance_type_xor_requirements`: instance_type and instance_requirements are mutually exclusive -- name an exact type, or describe requirements, not both
- `image_id_format`: image_id must be an AMI ID beginning with 'ami-'
- `cpu_credits_valid`: cpu_credits must be 'standard' or 'unlimited' when set
- `auto_recovery_valid`: auto_recovery must be 'default' or 'disabled' when set
- `shutdown_behavior_valid`: instance_initiated_shutdown_behavior must be 'stop' or 'terminate' when set
- `enclave_incompatible_with_hibernation`: enclave_enabled and hibernation_enabled cannot both be true
- `security_groups_conflict_with_interfaces`: when network_interfaces are declared, attach security groups on each interface instead of security_group_ids

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsLaunchTemplate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.launchTemplateId` | `string` | The launch template ID (e.g. "lt-0123456789abcdef0"). The primary handle other resources reference via status.outputs.launch_template_id -- auto-scaling groups, EKS managed node groups, and Batch compute environments all take this value. |
| `status.outputs.launchTemplateArn` | `string` | The ARN of the launch template, for IAM policies that scope ec2:RunInstances to approved templates. |
| `status.outputs.latestVersion` | `int64` | The latest version number of the template. Every spec change creates a new immutable version; this tracks the newest one. |
| `status.outputs.defaultVersion` | `int64` | The default version number -- what consumers referencing "$Default" launch from. The modules promote each new version to default, so this normally equals latest_version. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.instanceProfile` | AwsIamInstanceProfile | `status.outputs.instance_profile_arn` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.blockDeviceMappings[].ebs.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.networkInterfaces[].subnetId` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.networkInterfaces[].securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
