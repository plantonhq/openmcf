# AWS EC2 Instance — Architecture and Design

## Overview

`AwsEc2Instance` models one Amazon EC2 virtual machine at the full
provider surface: launch sources (inline AMI/type, a referenced launch
template, or both with per-instance overrides), networking down to
static ENI identity and multi-card secondary interfaces, storage with
KMS-encrypted volumes, the instance-posture surface (IMDSv2, CPU
topology, monitoring, protections), Spot and Capacity Reservation
purchase options, and placement (groups, tenancy, Dedicated Hosts).

The instance composes onto first-class neighbors and creates none of
them: `AwsSubnet`, `AwsSecurityGroup`, `AwsIamInstanceProfile`,
`AwsLaunchTemplate`, and `AwsKmsKey` all attach by reference. Access
posture is composed the same way -- attach a profile whose role carries
`AmazonSSMManagedInstanceCore` for keyless SSM access (the modern
default), or set `key_name` for SSH against a key pair you manage.

## Design Decisions

- **No synthetic access abstraction.** The spec models what EC2 models:
  an optional instance profile and an optional key pair. How you reach
  the instance (SSM, SSH, EC2 Instance Connect) is a property of the
  IAM profile and key material you compose, not an instance field --
  AWS has no "connection method" argument, so neither does this spec.
- **No module-side key generation.** The module never creates side
  resources it was not asked for -- private key material generated
  inside an IaC module would live in state and stack outputs, visible
  to everything that can read either. Bring a key pair by name, or go
  keyless with SSM.
- **Instance profile referenced by NAME.** The EC2 API attaches
  profiles by name (launch templates accept ARNs; instances do not),
  so the ref resolves `AwsIamInstanceProfile`'s `instance_profile_name`
  output. The asymmetry with `AwsLaunchTemplate.instance_profile`
  (ARN-based) is the providers' own.
- **Launch-template composition with override semantics.** A referenced
  template supplies the baseline; every inline field set on the
  instance overrides it for this instance only. AWS's launch validation
  (an image and a type must come from somewhere) is enforced at
  manifest validation, not at deploy.
- **The Name tag is the identity.** EC2 instances have no name
  argument; both engines carry `metadata.name` in the `Name` tag inside
  the shared identity tag set, so a manifest deploys identically on
  either engine.
- **Presence-honest optionals.** Subnet and security groups are
  optional exactly as the API treats them (unset falls back to the
  default VPC and default security group -- documented as a
  non-production posture); tri-states like `associate_public_ip_address`
  and `enable_primary_ipv6` stay unset unless the manifest speaks.
- **Enclaves XOR hibernation, IMDSv2 by choice, honest ForceNew.**
  AWS's own conflict rules (enclave/hibernation, placement group
  name/ID/host-resource-group, capacity-reservation preference/target,
  user-data plain/base64) are CEL rules that fail in seconds, and every
  create-time field says so in its comment.

## Deliberately Skipped Provider Surface

| Provider surface | Verdict | Reason |
|---|---|---|
| `get_password_data` / `password_data` | SKIP | Retrieves the Windows Administrator password into state and outputs -- credential material where it does not belong. Fetch it out-of-band (`aws ec2 get-password-data`). |
| `volume_tags` / per-device `tags` | SKIP | Identity tags derive from metadata platform-wide; custom user tags are a single platform-wide follow-up, not per-kind surface. |
| `network_interface` block | SKIP | Deprecated by the provider in favor of `primary_network_interface` (modeled) and standalone ENI attachments. |
| `security_groups` (by-name) | SKIP | EC2-Classic-era by-name attachment; `vpc_security_group_ids` (modeled as refs) is the VPC-era shape. |
| Secondary-interface security groups | N/A | Not configurable at launch in the EC2 API (the VPC default group applies); manage them post-launch on the ENI. |
| `instance_lifecycle` / `spot_instance_request_id` | SKIP (outputs) | Spot-internals attributes; the purchase option is already explicit in the spec. |

## Billing Note

The instance bills from `running` state by the second. The E2E scenario
runs one t3.micro for the create-verify-destroy window (minutes), an
effectively zero cost. Termination protection in presets guards pets;
the E2E scenario leaves it off so destroy is unattended.

## References

- [Amazon EC2 User Guide](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/concepts.html)
- [IMDSv2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html)
- [Spot Instances](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-spot-instances.html)
- [Placement groups](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html)
- [Optimize CPU options](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-optimize-cpu.html)
