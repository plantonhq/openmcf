# AzureApplicationSecurityGroup

## Overview

`AzureApplicationSecurityGroup` provisions an Azure Application Security
Group (ASG): a named, logical grouping of network interfaces that network
security group (NSG) rules can target by name instead of by IP address. A
rule can say "allow the web tier to reach the app tier" rather than "allow
10.0.1.0/24 to reach 10.0.2.0/24" -- the policy follows the workload as
instances scale in and out and change addresses.

## Why a First-Class Resource?

An ASG is real infrastructure with its own lifecycle:

- **Created once, referenced by many** -- network interfaces list the ASGs
  they join, VM scale set IP configurations declare membership, and NSG
  security rules reference source/destination ASGs; nothing lives inside
  the group
- **Independent lifecycle** -- the group is a stable anchor that outlives
  the individual NICs joining and leaving it
- **Intent-carrying composition seam** -- `application_security_group_id`
  is the join key every member and rule references

The inversion -- membership declared from the member side, not the group
side -- is what makes the ASG composable. The group itself holds no
members.

## Key Features

- **Micro-segmentation by role** -- name the group after the workload tier
  ("web", "app-tier", "db") so security rules read as intent
- **Address-independent policy** -- rules target the group, not a CIDR, so
  scaling never rewrites the firewall
- **Composable** -- the resource group is referenced by name, defaulting to
  an `AzureResourceGroup`'s output in composed environments

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region (must match the NICs referencing it) |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | -- | Group name, unique within the resource group (1-80 chars); fixed at creation |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `application_security_group_id` | Full ARM ID -- the join key NICs, scale sets, and NSG rules reference |
| `application_security_group_name` | The group's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureApplicationSecurityGroup
metadata:
  name: web-tier
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: web-tier
```

Target it from an NSG rule:

```yaml
spec:
  securityRules:
    - name: allow-web-to-app
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      priority: 100
      sourceApplicationSecurityGroupIds:
        - valueFrom:
            name: web-tier
      destinationApplicationSecurityGroupIds:
        - valueFrom:
            name: app-tier
```

## Lifecycle Notes

- `name` and `region` are **fixed at creation**; changing either replaces
  the group, and every NSG rule and network interface that referenced it
  must be re-pointed
- `tags` is the only field that updates in place
- An ASG can only be referenced by network interfaces in the **same
  region**

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
