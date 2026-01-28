This AWS Security Group Component extends the OpenMCF ecosystem by simplifying how you manage inbound and
outbound traffic rules for AWS EC2 instances and other resources. Through a single OpenMCF manifest, you can
define precise network boundaries, control allowed protocols and ports, and reference other security groups for more
complex configurations—all without leaving the standard OpenMCF workflow.

## Purpose and Functionality

• **Unified Network Access Management** – Centralize your security-group definitions in OpenMCF, specifying all
ingress and egress rules in a single YAML manifest.  
• **Consistent and Validated** – Leverages the same Protobuf-based validation standards as all OpenMCF
components, ensuring rule consistency and preventing misconfigurations before they ever reach AWS.  
• **Multi-Cloud Friendly** – By using OpenMCF’s CLI and manifest structure, you can continue working seamlessly
across different clouds. Even if you only need AWS Security Groups right now, the same approach easily extends to other
providers when needed.

## Key Benefits

1. **Fine-Grained Control**  
   Define rules for any protocol or port range, reference IPv4 or IPv6 CIDRs, and link to other security groups. This
   flexibility helps you tailor network policies to your exact requirements.

2. **Reuse and Modular Design**  
   Repeated rules or commonly used security profiles can be codified and reused as part of your manifest library,
   reducing duplication across environments or teams.

3. **Consistent Deployment Experience**  
   Just like other OpenMCF components, you can opt for Pulumi or Terraform behind the scenes—yet still manage
   everything through a single manifest and CLI command.

4. **Built-in Validation**  
   With built-in Protobuf validations (for example, required descriptions and restricted input lengths), you catch
   configuration errors early, streamlining your deployment process.

## Example Manifest

Below is a sample AWS Security Group manifest that shows how you can define an inbound rule for SSH (port 22) and allow
all outbound traffic. Notice how all keys use camel-case:

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsSecurityGroup
metadata:
  name: my-ssh-access-sg
  version:
    major: 1
    minor: 0
    message: "Initial version for AWS Security Group"
spec:
  vpcId: vpc-12345abcde
  description: "Allows inbound SSH traffic from a known CIDR block"
  ingress:
    - protocol: "tcp"
      fromPort: 22
      toPort: 22
      ipv4Cidrs:
        - "10.0.0.0/16"
  egress:
    - protocol: "-1"
      fromPort: 0
      toPort: 0
      ipv4Cidrs:
        - "0.0.0.0/0"
```

Once validated, deploying this manifest with OpenMCF automatically translates the specification into Pulumi or
Terraform code—depending on your chosen provisioner—and creates the security group with your defined rules.

---

This AWS Security Group Component adds powerful network security features to OpenMCF’s portfolio of cloud-native
capabilities. By wrapping complex AWS configurations in a straightforward, validated manifest, you gain a repeatable
pattern to apply and evolve across environments, speeding up your security operations in multi-cloud or hybrid-cloud
strategies.
