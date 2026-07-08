---
title: "Secure Private Notebook"
description: "The hardened posture for regulated data: private networking, CMEK, Shielded VM, Confidential Computing, and per-user credentials — every security lever the platform models, composed from first-class..."
type: "preset"
rank: "03"
presetSlug: "03-secure-private-notebook"
componentSlug: "vertex-ai-notebook"
componentTitle: "Vertex AI Notebook"
provider: "gcp"
icon: "package"
order: 3
---

# Secure Private Notebook

The hardened posture for regulated data: private networking, CMEK,
Shielded VM, Confidential Computing, and per-user credentials — every
security lever the platform models, composed from first-class nodes.

## What this preset creates

A Workbench instance named `regulated-research` on an `n2d-standard-4`
VM with no public IP, attached to a referenced VPC network and
subnetwork, running as a referenced dedicated service account. Both
disks are CMEK-encrypted under a referenced KMS key. Shielded VM
(Secure Boot, vTPM, integrity monitoring) and AMD SEV memory encryption
are enabled, and managed end-user credentials make notebook code act as
the signed-in user.

## When to use

- Notebooks handling PII or sensitive financial data
- Healthcare / life-sciences workloads under HIPAA controls
- Government workloads requiring FedRAMP-style hardening
- Any environment where per-user auditability of data access matters

## Composition

The referenced resources deploy as their own nodes: a `GcpVpcNetwork`
(`research-vpc`), a `GcpSubnetwork` (`research-subnet`), a
`GcpServiceAccount` (`research-notebook`), and a `GcpKmsKey`
(`research-data-key`). Rename the references to match your environment.

## Remix ideas

- Drop `confidentialInstanceConfig` (and switch back to an e2/n1
  machine type) when in-use memory encryption is not required.
- Add a static external IP via a referenced `GcpAddress` instead of
  `disablePublicIp` when an allowlisted egress address is the
  requirement rather than full network isolation.
- Add `instanceOwners` to pin the notebook to a single user identity.
