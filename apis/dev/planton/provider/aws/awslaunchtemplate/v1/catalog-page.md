# AWS Launch Template

Deploys an EC2 launch template: the versioned blueprint auto-scaling
groups, EKS managed node groups, and AWS Batch compute environments launch
instances from -- AMI, instance type or attribute-based requirements,
storage, networking, IAM identity, metadata posture, and purchase options
defined once and referenced everywhere.

## What Gets Created

When you deploy an AwsLaunchTemplate resource, Planton provisions:

- **Launch template** — an `aws_launch_template` / `ec2.LaunchTemplate`
  named from `metadata.name` (truncated to AWS's 125-character limit when
  necessary), with every configured launch attribute
- **A new template version on every change** — versions are immutable in
  AWS; each update creates the next version and promotes it to the
  template default, so consumers following `$Default` roll forward on
  their next launch or instance refresh

Identity tags are applied to the template AND propagated to the instances
and volumes each launch creates (template tags do not propagate on their
own).

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An IAM instance profile** (`AwsIamInstanceProfile`) if instances need AWS API access (SSM, ECR, S3).
- **Security groups** (`AwsSecurityGroup`) for network access control.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLaunchTemplate
metadata:
  name: web
spec:
  region: us-west-2
  imageId: ami-0123456789abcdef0
  instanceType: t3.small
  instanceProfile:
    valueFrom:
      kind: AwsIamInstanceProfile
      name: web-profile
      fieldPath: status.outputs.instance_profile_arn
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: web-sg
        fieldPath: status.outputs.security_group_id
  metadataOptions:
    httpTokens: required
```

```shell
planton apply -f launch-template.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region the template is created in. A template is regional: consumers must live in the same region. | Required; non-empty |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `description` | `string` | — | Version description recorded on each template version — use it as a change-log line. Max 255 chars. |
| `imageId` | `string` | — | AMI to boot from. Optional at the template level, but auto-scaling groups require their template to carry one. |
| `instanceType` | `string` | — | Exact EC2 instance type. Mutually exclusive with `instanceRequirements`. |
| `instanceRequirements` | `object` | — | Attribute-based selection: `memoryMib` + `vcpuCount` ranges (required), type allow/deny lists, generations, CPU manufacturers, accelerators, local storage, price protection. Mutually exclusive with `instanceType`. |
| `keyName` | `string` | — | EC2 key pair for SSH. Leave unset for keyless (SSM) fleets. |
| `userData` | `string` | — | Plain-text cloud-init / shell script; both modules base64-encode it. Max 16 KiB. |
| `instanceProfile` | `string \| valueFrom` | — | IAM instance profile ARN — the instance's AWS identity. Defaults to referencing an `AwsIamInstanceProfile`'s `instance_profile_arn`. |
| `securityGroupIds` | `string[] \| valueFrom` | `[]` | Security groups for the primary interface. Mutually exclusive with per-interface groups in `networkInterfaces`. |
| `ebsOptimized` | `bool` | type default | Dedicated EBS throughput on types where it is optional. |
| `blockDeviceMappings` | `object[]` | AMI mapping | Root-volume overrides and data volumes: size, type (`gp3` etc.), IOPS/throughput, encryption + `AwsKmsKey` reference, snapshot, delete-on-termination tri-state, `noDevice` suppression. |
| `networkInterfaces` | `object[]` | `[]` | Explicit interfaces: device/card index, `efa` types, public-IP tri-state, subnet + security-group references, static IPs, IPv6, IPv4/IPv6 prefix delegation. |
| `metadataOptions` | `object` | AWS defaults | IMDS posture: `httpTokens: required` enforces IMDSv2; hop limit 2 lets containers reach metadata. |
| `detailedMonitoring` | `bool` | `false` | 1-minute CloudWatch metrics (billed) — scaling policies react faster. |
| `placement` | `object` | — | AZ pinning, placement group, partition, tenancy. |
| `cpuOptions` | `object` | — | Core count / threads-per-core (license trimming), AMD SEV-SNP. |
| `cpuCredits` | `string` | `unlimited` (recent T families) | Burstable credit mode: `standard` or `unlimited`. |
| `spotOptions` | `object` | On-Demand | Makes every launch a Spot request: max price, request type, interruption behavior. |
| `enclaveEnabled` | `bool` | `false` | Nitro Enclaves parent. Incompatible with hibernation. |
| `hibernationEnabled` | `bool` | `false` | Pre-provision for hibernation. Incompatible with enclaves. |
| `autoRecovery` | `string` | `default` | Simplified automatic recovery: `default` or `disabled`. |
| `privateDnsNameOptions` | `object` | — | Hostname scheme (`ip-name` / `resource-name`) and DNS records. |
| `disableApiStop` | `bool` | `false` | Stop protection for launched instances. |
| `disableApiTermination` | `bool` | `false` | Termination protection — note an ASG cannot scale in protected instances. |
| `instanceInitiatedShutdownBehavior` | `string` | `stop` | What an OS shutdown does: `stop` or `terminate`. |

## Examples

### Attribute-based Spot template (no instance type named)

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLaunchTemplate
metadata:
  name: spot-workers
spec:
  region: us-west-2
  imageId: ami-0123456789abcdef0
  instanceRequirements:
    memoryMib:
      min: 4096
      max: 16384
    vcpuCount:
      min: 2
      max: 8
    cpuManufacturers: [intel, amd]
    instanceGenerations: [current]
    bareMetal: excluded
  spotOptions:
    spotInstanceType: one-time
    instanceInterruptionBehavior: terminate
```

### Hardened template with encrypted storage

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLaunchTemplate
metadata:
  name: hardened-base
spec:
  region: us-west-2
  imageId: ami-0123456789abcdef0
  instanceType: m6i.large
  metadataOptions:
    httpTokens: required
    httpPutResponseHopLimit: 1
  blockDeviceMappings:
    - deviceName: /dev/xvda
      ebs:
        volumeSizeGb: 50
        volumeType: gp3
        encrypted: true
        kmsKeyId:
          valueFrom:
            kind: AwsKmsKey
            name: fleet-key
            fieldPath: status.outputs.key_arn
  disableApiStop: true
  disableApiTermination: true
```

### Kubernetes-CNI-friendly interface with prefix delegation

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLaunchTemplate
metadata:
  name: k8s-workers
spec:
  region: us-west-2
  networkInterfaces:
    - deviceIndex: 0
      ipv4PrefixCount: 4
      securityGroupIds:
        - valueFrom:
            kind: AwsSecurityGroup
            name: node-sg
            fieldPath: status.outputs.security_group_id
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `launch_template_id` | Template ID — what `AwsAutoScalingGroup`, EKS node groups, and Batch compute environments reference |
| `launch_template_arn` | ARN, for IAM policies scoping `ec2:RunInstances` to approved templates |
| `latest_version` | Newest version number (every change creates one) |
| `default_version` | Version consumers referencing `$Default` launch from — the modules promote each new version |

## Related Components

- [AwsAutoScalingGroup](/docs/catalog/aws/awsautoscalinggroup) — launches and manages fleets from this template
- [AwsBatchComputeEnvironment](/docs/catalog/aws/awsbatchcomputeenvironment) — launches Batch compute from this template
- [AwsIamInstanceProfile](/docs/catalog/aws/awsiaminstanceprofile) — the instance identity the template attaches
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — network access control for launched instances
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed EBS encryption keys
- [AwsSubnet](/docs/catalog/aws/awssubnet) — subnet placement for explicit network interfaces
