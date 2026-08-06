---
title: "EC2 Instance"
description: "EC2 Instance deployment documentation"
icon: "package"
order: 100
componentName: "awsec2instance"
---

# AWS EC2 Instance

Deploys a single Amazon EC2 virtual machine through one declarative manifest: the launch source (an AMI + instance type, a launch template, or a template with per-instance overrides), network placement, IAM identity, storage, instance posture (IMDSv2, monitoring, CPU topology), purchase options (On-Demand or Spot), and lifecycle protections.

## What Gets Created

When you deploy an AwsEc2Instance resource, Planton provisions:

- **EC2 Instance** — exactly one virtual machine whose `Name` tag carries `metadata.name`, launched from the inline AMI/type or the referenced launch template
- **Root volume override** — when `rootBlockDevice` is set, the AMI's boot volume reshaped (size, gp3 tuning, encryption with an optional KMS reference)
- **Additional EBS volumes** — from `ebsBlockDevices`, attached at launch by device name
- **Secondary network interfaces** — from `secondaryNetworkInterfaces`, on non-primary network cards of multi-card instance types

Nothing else is created: the subnet, security groups, IAM instance profile, launch template, key pair, and KMS keys are referenced first-class resources. All created resources are tagged with Planton metadata (organization, environment, resource kind, resource ID).

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A subnet** (recommended) — reference an `AwsSubnet`; unset launches into the account's default VPC, which is not a production posture
- **Security groups** (recommended) — reference `AwsSecurityGroup` resources; unset falls back to the VPC's default group
- **An IAM instance profile** (recommended) — reference an `AwsIamInstanceProfile` whose role carries `AmazonSSMManagedInstanceCore` for keyless SSM access
- **A launch template** (optional) — reference an `AwsLaunchTemplate` to inherit an org-wide golden baseline

## Quick Start

Create a file `instance.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEc2Instance
metadata:
  name: my-instance
spec:
  region: us-west-2
  ami: ami-0123456789abcdef0
  instanceType: t4g.small
  subnetId:
    valueFrom:
      kind: AwsSubnet
      name: my-private-subnet
      fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: my-instance-sg
        fieldPath: status.outputs.security_group_id
  instanceProfile:
    valueFrom:
      kind: AwsIamInstanceProfile
      name: my-ssm-profile
      fieldPath: status.outputs.instance_profile_name
  metadataOptions:
    httpTokens: required
```

Deploy:

```shell
planton apply -f instance.yaml
```

This launches one hardened instance reachable through SSM Session Manager -- no SSH key, no inbound ports.

## Configuration Reference

### Launch Source

| Field | Type | Description |
|-------|------|-------------|
| `ami` | `string` | The AMI the instance boots from (`ami-...`, region-specific). Required unless the launch template supplies one. Changing it replaces the instance. |
| `instanceType` | `string` | vCPU/memory/network shape (e.g. `t4g.nano`, `m7g.large`). Required unless the launch template supplies one. |
| `launchTemplate.id` | `ref` | Reference an `AwsLaunchTemplate`'s `launch_template_id` output. Mutually exclusive with `launchTemplate.name`. |
| `launchTemplate.version` | `string` | A version number, `$Latest`, or `$Default` (the default). Both Planton launch-template modules promote each release to default. |

Every inline field set on the instance overrides the template's value -- the template is the baseline, the instance spec carries the deviations.

### IAM and Access

| Field | Type | Description |
|-------|------|-------------|
| `instanceProfile` | `ref` | The IAM instance profile by NAME (the EC2 API takes names; reference `instance_profile_name`). The instance's identity for SSM, ECR, S3, and every AWS API call. |
| `keyName` | `string` | An existing EC2 key pair for SSH. Leave unset for keyless SSM posture. |

### Networking

| Field | Type | Description |
|-------|------|-------------|
| `subnetId` | `ref` | The subnet (reference `AwsSubnet.subnet_id`). Unset = default VPC (not production). Create-time. |
| `securityGroupIds` | `ref[]` | Groups on the primary interface. Unset = VPC default group (not production). Updatable in place. |
| `primaryNetworkInterfaceId` | `string` | Attach an EXISTING ENI as eth0 (static IP/MAC identity). Excludes all inline networking fields. |
| `privateIp` / `secondaryPrivateIps` | `string` / `string[]` | Static primary private IPv4 / additional addresses on eth0. |
| `associatePublicIpAddress` | `bool?` | Tri-state: unset inherits the subnet's setting. Public IPv4 addresses bill hourly. |
| `sourceDestCheck` | `bool?` | Default `true`. Set `false` only for NAT/router/VPN instances that forward traffic. |
| `ipv6AddressCount` / `ipv6Addresses` | `int` / `string[]` | Auto-assigned count XOR explicit addresses from the subnet's IPv6 range. |
| `enablePrimaryIpv6` | `bool?` | Pin the first IPv6 address as the stable primary (required posture for IPv6-only workloads). |
| `privateDnsNameOptions` | `object` | Hostname scheme (`ip-name`/`resource-name`) and A/AAAA record publication. |
| `secondaryNetworkInterfaces` | `object[]` | Interfaces on network cards 1+ of multi-card types: card index, subnet reference, address count, delete-on-termination. |

### Storage

| Field | Type | Description |
|-------|------|-------------|
| `rootBlockDevice` | `object` | Reshape the AMI's boot volume: `volumeSizeGb`, `volumeType` (gp3/gp2/io1/io2/st1/sc1/standard), `iops`, `throughputMibps` (gp3 only), `encrypted`, `kmsKeyId` (ref), `deleteOnTermination`. |
| `ebsBlockDevices` | `object[]` | Additional data volumes keyed by `deviceName`, same shape plus `snapshotId`. Create-time on the instance -- prefer standalone volumes for independent lifecycles. |
| `ephemeralBlockDevices` | `object[]` | Instance-store mappings (`deviceName`, `virtualName` XOR `noDevice`). Data does not survive stop/terminate. |
| `ebsOptimized` | `bool` | Dedicated EBS throughput on types where it is optional. |

### Instance Posture

| Field | Type | Description |
|-------|------|-------------|
| `metadataOptions` | `object` | IMDS posture. `httpTokens: required` enforces IMDSv2 (the recommended hardening); `httpPutResponseHopLimit: 2` keeps containers working; plus endpoint on/off, IPv6 endpoint, tag exposure. |
| `detailedMonitoring` | `bool` | 1-minute CloudWatch metrics instead of the free 5-minute tier. |
| `cpuOptions` | `object` | `coreCount` + `threadsPerCore` (1 disables SMT -- per-core licensing), `amdSevSnp`, `nestedVirtualization`. Fixed at launch. |
| `cpuCredits` | `string` | `standard` or `unlimited` for burstable (T-family) types. |
| `enclaveEnabled` / `hibernationEnabled` | `bool` | Nitro Enclaves parent XOR hibernation support (mutually exclusive per AWS). |
| `autoRecovery` | `string` | `default` (recover onto healthy hardware) or `disabled`. |
| `instanceInitiatedShutdownBehavior` | `string` | `stop` (default) or `terminate` on OS shutdown. |
| `disableApiStop` / `disableApiTermination` | `bool?` | API stop/termination protection for pet instances. |

### Purchase Options and Placement

| Field | Type | Description |
|-------|------|-------------|
| `spotOptions` | `object` | Presence = Spot. `maxPrice` (unset caps at On-Demand price -- the AWS recommendation), `spotInstanceType` (`one-time`/`persistent`), `instanceInterruptionBehavior` (`terminate`/`stop`/`hibernate` -- stop/hibernate need persistent), `validUntil`. |
| `capacityReservation` | `object` | `preference` (`open`/`none`) XOR a specific `capacityReservationId` / `capacityReservationResourceGroupArn`. |
| `placement` | `object` | `availabilityZone`, placement group (`groupName` XOR `groupId`, `partitionNumber`), `tenancy` (`default`/`dedicated`/`host`), `hostId`, `hostResourceGroupArn`. Fixed at launch. |

### User Data

| Field | Type | Description |
|-------|------|-------------|
| `userData` | `string` | Plain-text cloud-init/shell script (16 KiB), run on FIRST boot. `${...}` shell syntax passes through literally. |
| `userDataBase64` | `string` | Base64 for binary payloads. Mutually exclusive with `userData`. |
| `userDataReplaceOnChange` | `bool` | Replace the instance on user-data changes instead of stop-update-start -- the right choice when user data IS the provisioning mechanism. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_id` | The instance ID -- what `AwsLbTargetGroup` instance targets reference |
| `arn` | The instance ARN, for IAM policies and EventBridge rules |
| `instance_state` | Lifecycle state as of the last deploy |
| `availability_zone` | The AZ the instance runs in |
| `private_ip` / `private_dns` | The primary private address and VPC-internal hostname |
| `public_ip` / `public_dns` | The public address pair (empty for private-only instances; changes across stop/start -- compose an `AwsElasticIp` for stability) |
| `primary_network_interface_id` | The eth0 ENI -- the attachment point for Elastic IP associations and flow logs |

## Related Resources

- [AWS Launch Template](/docs/catalog/aws/launch-template) — the golden baseline this instance can launch from
- [AWS Auto Scaling Group](/docs/catalog/aws/auto-scaling-group) — when the pet should become cattle
- [AWS Subnet](/docs/catalog/aws/subnet) and [AWS Security Group](/docs/catalog/aws/security-group) — network placement and reach
- [AWS IAM Instance Profile](/docs/catalog/aws/iam-instance-profile) — the instance's AWS identity
- [AWS Elastic IP](/docs/catalog/aws/elastic-ip) — a stable public address across stop/start cycles
