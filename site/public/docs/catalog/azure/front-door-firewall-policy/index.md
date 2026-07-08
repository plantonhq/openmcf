---
title: "Front Door Firewall Policy"
description: "Front Door Firewall Policy deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorfirewallpolicy"
---

# Azure Front Door Firewall Policy

Creates a Web Application Firewall (WAF) policy for Azure Front Door -- the edge rule set (custom match/rate-limit rules, Microsoft's managed rule sets on Premium, log scrubbing) that an AzureFrontDoorSecurityPolicy attaches to the hostnames a profile serves. Global and resource-group-scoped; one policy commonly protects many profiles.

## What Gets Created

When you deploy an AzureFrontDoorFirewallPolicy resource, Planton provisions:

- **Front Door WAF Policy** -- an `azurerm_cdn_frontdoor_firewall_policy` in the referenced resource group, carrying the policy settings, custom rules, managed rule sets, and log-scrubbing rules

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the policy in (referenced through `resourceGroup`)

## Quick Start

Create a file `waf-policy.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorFirewallPolicy
metadata:
  name: edge-waf
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorFirewallPolicy.edge-waf
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-resource-group
      fieldPath: status.outputs.resource_group_name
  policyName: edgewaf
  sku: PREMIUM
  mode: PREVENTION
  managedRules:
    - type: Microsoft_DefaultRuleSet
      version: "2.1"
      action: RULE_SET_BLOCK
```

Deploy it:

```shell
planton apply -f waf-policy.yaml
```

## Notes

- The policy enforces nothing until an **AzureFrontDoorSecurityPolicy** associates it with a profile's endpoints or custom domains.
- The policy's `sku` must **match the profile's sku** at association time, and Azure refuses a PREMIUM-to-STANDARD change outright.
- Start new policies in `DETECTION` mode against real traffic, tune exclusions, then flip to `PREVENTION` (an in-place update).
