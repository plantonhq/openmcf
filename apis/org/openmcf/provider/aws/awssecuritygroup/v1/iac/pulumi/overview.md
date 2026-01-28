# Overview

The **aws_security_group** Component in OpenMCF streamlines the creation and management of AWS Security Groups—
crucial tools for controlling network traffic at the instance or service level. Security groups act as stateful
firewalls
that filter inbound and outbound traffic based on configurable rules, enhancing the security posture of your AWS
environment.

By embracing the familiar Kubernetes-like resource model (`apiVersion`, `kind`, `metadata`, `spec`, `status`),
**aws_security_group** integrates into OpenMCF’s multi-cloud framework. This allows teams to validate manifests
locally, apply consistent Protobuf-based validations, and deploy via Pulumi or Terraform interchangeably, all with a
single YAML specification. Whether you need simple inbound HTTP access or more advanced scenarios involving multiple
CIDR ranges and self-referencing rules, **aws_security_group** offers a clean, flexible approach to defining network
boundaries in your AWS environment.

---

## Key Features

- **Easy Ingress & Egress Rules**  
  Quickly define inbound (ingress) and outbound (egress) traffic rules for specific ports, protocols, and IP ranges.

- **IPv4 & IPv6 Support**  
  Manage dual-stack environments by specifying both IPv4 and IPv6 CIDR blocks in your security group rules.

- **Self-Referencing**  
  Allow traffic within the same security group using a simple flag, enabling microservice-style communications without
  repetitive resource lookups.

- **Pulumi & Terraform Integration**  
  Seamlessly deploy your AWS Security Group specification using either Pulumi or Terraform, keeping the same high-level
  YAML manifest in OpenMCF.

- **Consistent Resource Model**  
  Follows OpenMCF’s standard resource structure, ensuring a predictable user experience and robust field
  validations across all your multi-cloud deployments.

---

## Next Steps

- Check out the [README.md](./README.md) for detailed resource configuration, field descriptions, and best practices.
- Refer to the [examples.md](./examples.md) for practical usage scenarios and sample YAML manifests.
- Explore the broader OpenMCF documentation for insights on multi-cloud deployments, CLI usage, and advanced
  features.
