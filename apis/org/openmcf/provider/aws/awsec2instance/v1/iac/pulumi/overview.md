# Overview

The **AwsEc2Instance** API resource provides a streamlined, production-ready approach to deploying single EC2 virtual machine instances on AWS. Rather than wrestling with the dozens of configuration parameters that `aws_instance` offers, **AwsEc2Instance** focuses on the essential 80/20 of EC2 deployment: secure networking (private subnets, security groups), modern access patterns (AWS Systems Manager Session Manager by default), and the minimum viable set of configuration options needed for typical production workloads.

By defining a single YAML manifest conforming to the familiar Kubernetes-like structure (`apiVersion`, `kind`, `metadata`, `spec`, `status`), teams can validate, provision, and maintain EC2 instances consistently across multiple environments. Whether you're deploying a bastion host, a database server, or an application server in a private subnet, **AwsEc2Instance** handles the heavy lifting—wiring up networking, security, IAM profiles, and SSH keys—while adhering to AWS best practices for high availability and security.

---

## Key Features

- **Secure-by-Default Networking**  
  Designed for deployment in private subnets with custom VPC configuration and security groups. No default VPC, no public IPs by default—network isolation is the starting point, not an afterthought.

- **Modern Access Patterns**  
  Defaults to AWS Systems Manager Session Manager (SSM) for shell access, eliminating the need for SSH keys, bastion hosts, or open SSH ports. Alternatively supports traditional SSH via bastion or EC2 Instance Connect for teams with existing workflows.

- **Minimal, Focused Configuration**  
  Exposes only the most commonly needed parameters: AMI ID, instance type, subnet, security groups, connection method, root volume size, and optional user data. This 80/20 approach keeps manifests readable and maintainable while still covering typical production scenarios.

- **Composability with Foreign Keys**  
  Uses `StringValueOrRef` for VPC subnet IDs, security group IDs, and IAM instance profile ARNs, enabling seamless integration with other OpenMCF resources like `AwsVpc`, `AwsSecurityGroup`, and `AwsIamRole`.

- **Consistent Resource Model**  
  Uses OpenMCF's standard resource layout with Protobuf definitions and built-in validation rules. Advanced CEL validation ensures configuration correctness (e.g., SSM requires IAM instance profile, SSH requires key_name).

- **Pulumi & Terraform Integration**  
  Provision the same EC2 instance specification using either Pulumi (Go) or Terraform (HCL), unified by the OpenMCF CLI's straightforward commands and orchestration.

- **Immutable Infrastructure Philosophy**  
  Designed for the modern pattern of baking configuration into AMIs (via Packer) rather than post-launch configuration management. User data is available for bootstrap scripts, but the expectation is that most setup happens at AMI build time.

---

## Pulumi Module Architecture

The Pulumi implementation follows a clean separation of concerns across multiple files:

### Core Module Structure

- **`module/main.go`**: Entry point that initializes the AWS provider (using either default credentials or explicit access keys from `ProviderConfig`) and orchestrates the resource creation workflow.

- **`module/locals.go`**: Initializes derived values and local variables from the input spec, including metadata extraction, tag generation, and spec field processing.

- **`module/ec2_instance.go`**: Contains the core EC2 instance provisioning logic. This is where the "smart" behavior lives:
  - Validates and enforces connection method requirements (SSM requires IAM profile, SSH/Bastion requires key pair)
  - Auto-generates SSH key pairs if needed (using `tls.PrivateKey` and `ec2.KeyPair`)
  - Assembles the complete `ec2.Instance` resource with all required and optional fields
  - Exports stack outputs (instance ID, private IP, DNS name, availability zone, and optionally generated SSH keys)

- **`module/outputs.go`**: Defines output key constants used when exporting Pulumi stack outputs (e.g., `instance_id`, `private_ip`, `ssh_private_key`).

### Resource Relationships

```
AwsEc2Instance Spec
    │
    ├─► AWS Provider (configured with credentials)
    │
    ├─► VPC Subnet (referenced via StringValueOrRef)
    │
    ├─► Security Groups (1+ referenced via StringValueOrRef)
    │
    ├─► IAM Instance Profile (optional, referenced via StringValueOrRef)
    │       └─► Required if connection_method = SSM
    │
    ├─► SSH Key Pair (optional, can be auto-generated)
    │       └─► Required if connection_method = BASTION or INSTANCE_CONNECT
    │
    └─► EC2 Instance
            ├─► Root EBS Volume (gp3, size configurable)
            ├─► Tags (from metadata.labels + auto-generated tags)
            └─► Outputs → exported to Pulumi stack state
```

### Deployment Flow

1. **Provider Initialization**: The module first configures the AWS provider using either default credentials (from environment) or explicit credentials passed via `ProviderConfig` in the stack input.

2. **Locals Initialization**: Metadata (name, labels), tags, and spec field values are extracted and prepared for resource creation.

3. **Connection Method Validation**: Based on `connection_method`, the module enforces requirements:
   - **SSM**: Validates that `iam_instance_profile_arn` is present
   - **BASTION/INSTANCE_CONNECT**: If `key_name` is not provided, auto-generates a new RSA-4096 key pair and registers it with AWS

4. **Security Group Processing**: Converts `StringValueOrRef` security group IDs to concrete string values for the EC2 instance's VPC security group IDs.

5. **EC2 Instance Creation**: Assembles all required and optional arguments (`ami_id`, `instance_type`, `subnet_id`, `vpc_security_group_ids`, `key_name`, `iam_instance_profile`, `root_block_device`, `user_data`, tags) and provisions the instance.

6. **Output Export**: Exports instance ID, private IP, private DNS name, availability zone, and (if auto-generated) SSH key material to Pulumi stack outputs. These outputs are captured in the `AwsEc2Instance` resource's `status.outputs` field.

---

## State Management

Pulumi manages the state of the EC2 instance in its backend (local file, S3, Pulumi Cloud, etc.). On each `pulumi update`, Pulumi calculates the diff between the desired state (from the manifest) and the actual state (from AWS APIs), then applies the minimal set of changes needed to reconcile.

**Important behaviors:**

- **Immutable Fields**: Changing AMI ID, instance type, or subnet requires instance replacement (destroy + create). Pulumi will show this in the preview.
- **In-Place Updates**: Changing tags, security groups, or IAM instance profile can typically be done in place.
- **User Data**: Changes to `user_data` trigger instance replacement (AWS limitation—user data is applied only at launch time).

---

## Troubleshooting

### Instance Launch Failures

- **AMI not found**: Ensure `ami_id` is valid in the target AWS region. AMI IDs are region-specific.
- **Subnet not accessible**: Ensure the subnet ID references a valid, existing subnet in the same region and VPC.
- **Security group validation failure**: Ensure all security group IDs reference existing groups in the same VPC.

### Connection Method Issues

- **SSM access fails**: Ensure the IAM instance profile includes the `AmazonSSMManagedInstanceCore` managed policy and the instance can reach the SSM endpoint (either via VPC endpoint or internet gateway).
- **SSH key not found**: If using BASTION or INSTANCE_CONNECT with a user-supplied `key_name`, ensure the key pair exists in AWS.
- **Auto-generated key not accessible**: The private key is exported as a Pulumi stack output. Retrieve it with `openmcf pulumi stack output ssh_private_key --show-secrets`.

### Pulumi State Conflicts

- **Concurrent updates**: If multiple users or CI/CD pipelines attempt to update the same Pulumi stack simultaneously, state conflicts can occur. Use Pulumi's backend locking or coordinate deployments.
- **Deleted resources outside Pulumi**: If the EC2 instance is deleted manually (via AWS console or CLI), Pulumi will detect the drift on the next `pulumi refresh` and offer to reconcile.

---

## Next Steps

- Refer to the [README.md](./README.md) for detailed setup instructions, CLI usage patterns, and configuration field reference.
- Review the [examples.md](./examples.md) to explore common EC2 instance use cases, from minimal SSM-based deployments to SSH-enabled bastion hosts.
- Check out the wider OpenMCF documentation for deeper insights into multi-cloud deployments, advanced features, and the CLI usage patterns.
- See `docs/README.md` in the component root for a comprehensive discussion of the infrastructure deployment landscape and the philosophy behind this API.

