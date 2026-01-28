# AWS Security Group Component Overview

This document provides an at-a-glance description of the **aws_security_group** resource within OpenMCF. It
explains how this Component manages AWS Security Groups (SGs), the core fields in its API, and how it fits into your
multi-cloud deployments.

---

## Resource Description

An AWS Security Group acts as a virtual firewall for EC2 instances, controlling inbound and outbound traffic. In
OpenMCF, the `aws_security_group` resource is used to create or manage one of these Security Groups within a
specific VPC. By defining ingress and egress rules, you can precisely control which ports and protocols are allowed, as
well as restrict or allow traffic from specific CIDR ranges or other security groups.

---

## Behavior in OpenMCF

• **Multi-Cloud Consistency**  
Through our standard Protobuf API definitions, this Component retains the same look and feel as other OpenMCF
resources. Once you define your SG in a YAML manifest, the OpenMCF CLI orchestrates the provisioning steps via
either Pulumi or Terraform—no need to worry about differing AWS or IaC syntax.

• **Lifecycle and Audit**  
The `metadata` and `status` fields in the resource provide Kubernetes-style lifecycle management, including creation
timestamps, version details, and future lifecycle signals (e.g., marking the SG for deletion).

• **VPC Integration**  
An SG must always reside in a specific AWS VPC (`vpc_id` in the spec). If the VPC doesn’t exist yet, you must ensure
it’s created elsewhere in your overall OpenMCF deployment before or alongside this resource.

---

## Argument Reference

The `aws_security_group` resource follows the conventional OpenMCF resource schema:

1. **apiVersion**  
   • Must be `aws.openmcf.org/v1`.

2. **kind**  
   • Must be `aws_security_group` (or a matching variant recognized by the system).

3. **metadata**  
   • Required.  
   • Includes `name` (unique within your org/project), along with optional fields like `labels` and `annotations`.  
   • The `metadata.version.message` field must be set, ensuring that each resource revision is clearly described.

4. **spec**  
   • **vpcId** (string, required)  
   The identifier of the VPC in which to place the Security Group (e.g., `vpc-12345abcde`).  
   • **description** (string, required)  
   A short description of this Security Group’s purpose. AWS requires this field, up to 255 characters.  
   • **ingressRules** (repeated)  
   A list of inbound rules (see [SecurityGroupRule](#securitygrouprule) below). Omitting it means no inbound traffic is
   allowed (the default deny).  
   • **egressRules** (repeated)  
   A list of outbound rules (see [SecurityGroupRule](#securitygrouprule) below). Omitting it typically defaults to allow
   all outbound.

### SecurityGroupRule

Each `SecurityGroupRule` in `ingressRules` or `egressRules` includes:

• **protocol** (string, required)  
Common values: `tcp`, `udp`, `icmp`, or `-1` for all protocols.

• **fromPort** / **toPort** (int32)  
Defines the port range for the rule. For a single port, set both to the same value.

• **ipv4Cidrs** / **ipv6Cidrs** (string list)  
Ranges of allowed or targeted IP addresses. Example: `['0.0.0.0/0']` for all IPv4.

• **sourceSecurityGroupIds** / **destinationSecurityGroupIds** (string list)  
When referencing other Security Groups. Often used to allow internal traffic among multiple SGs.

• **selfReference** (bool)  
A shortcut to reference the Security Group itself, allowing internal traffic within the same SG.

• **description** (string)  
A rule-level description, up to 255 characters.

---

## Attributes Reference

When the deployment completes, the OpenMCF engine populates the `status` field. Notable subfields include:

• **status.lifecycle**  
Indicates if the resource is active, updated, or marked for cleanup.

• **status.audit**  
Creation and last-updated timestamps, plus user info (if available).

• **status.stackUpdateId**  
The ID of the Pulumi/Terraform job that performed the provisioning.

• **status.outputs**  
Although mostly used for more complex networking resources, any relevant SG outputs may appear here (e.g., the final SG
ID). Future enhancements may provide deeper insights (like referencing child or peer resources).

---

## Further Reading

• For hands-on deployment examples, check out **examples.md** in this repository.  
• Refer to the OpenMCF [Comprehensive Guide](../../../../docs/guide.md) for an overview of how to incorporate
multiple resources in a unified workflow.  
• If you’re new to OpenMCF, the restaurant analogy in the main guide helps clarify how each resource and
underlying IaC module fits into our multi-cloud approach.

---

_This overview is part of the standard documentation set for OpenMCF Components. It ensures a consistent look and
feel across providers and resource types, so you can confidently manage AWS Security Groups alongside other building
blocks in your cloud infrastructure._

