---
title: "Front Door Security Policy"
description: "Front Door Security Policy deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorsecuritypolicy"
---

# Azure Front Door Security Policy

Deploys a security policy inside an Azure Front Door (Standard/Premium) profile -- the association that attaches a Front Door WAF policy (AzureFrontDoorFirewallPolicy) to the hostnames the profile serves. A WAF policy enforces nothing until a security policy associates it; this kind is the enforcement seam. The `domainIds` list names the protected hostnames, and each entry is either an endpoint's ARM ID (the generated `*.azurefd.net` hostname) or a custom domain's ARM ID -- Azure accepts both interchangeably, and one list can mix them. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Security Policy** -- a named child of the profile binding one WAF policy to a domain list
- **Live enforcement** -- from the moment it deploys, the WAF's rules run on every request to the associated domains

## The Association in the Front Door Family

- **AzureFrontDoorFirewallPolicy** -- the WAF rule set this association enforces, referenced by `firewallPolicyId`; the WAF's sku must match the profile's
- **AzureFrontDoorProfile** -- the parent container whose sku also caps the domain list (100 on STANDARD, 500 on PREMIUM)
- **AzureFrontDoorEndpoint / AzureFrontDoorCustomDomain** -- the two domain types a row may reference

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door profile** with at least one endpoint or validated custom domain.
- **A Front Door WAF policy of the SAME sku** as the profile -- Azure rejects a mismatched pairing at deploy time.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Security Policy**, and click **Deploy**. The wizard walks you through the profile and name, then the enforcement scope: the WAF policy and the protected-domains list, whose per-row toggle switches between endpoint and custom-domain targets. Start from the **Endpoint Association** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorSecurityPolicy
metadata:
  name: waf-attach
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  securityPolicyName: waf-attach
  firewallPolicyId:
    valueFrom:
      kind: AzureFrontDoorFirewallPolicy
      name: edge-waf
      fieldPath: status.outputs.firewall_policy_id
  domainIds:
    - valueFrom:
        kind: AzureFrontDoorEndpoint
        name: web-endpoint
        fieldPath: status.outputs.endpoint_id
```

```shell
planton apply -f front-door-security-policy.yaml
```

This turns the WAF on for the endpoint's generated hostname. Add custom domains to the same list as they validate -- the list updates in place.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  securityPolicyName: waf-attach
  firewallPolicyId:
    valueFrom:
      kind: AzureFrontDoorFirewallPolicy
      name: edge-waf
      fieldPath: status.outputs.firewall_policy_id
  domainIds:
    - valueFrom:
        kind: AzureFrontDoorEndpoint
        name: web-endpoint
        fieldPath: status.outputs.endpoint_id
    - valueFrom:
        kind: AzureFrontDoorCustomDomain
        name: www-example-com
        fieldPath: status.outputs.custom_domain_id
```

## Key Configuration

**Security policy name** -- unique within the profile; letters, digits, and hyphens. ForceNew: a rename replaces the association with a brief enforcement gap.

**Firewall policy** -- the WAF to enforce. ForceNew: swapping it replaces the association (fast, metadata-only). The skus must match.

**Domain list** -- at least 1, at most 500 (100 on a STANDARD profile -- the cap rides the profile's sku, checked at deploy time). Updatable in place. The WAF applies to ALL paths (`/*`) on every associated domain -- Azure accepts no other pattern; scope enforcement by choosing WHICH domains to associate.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |
| **AzureFrontDoorFirewallPolicy** | `firewallPolicyId` | `status.outputs.firewall_policy_id` |
| **AzureFrontDoorEndpoint** | `domainIds[]` | `status.outputs.endpoint_id` |
| **AzureFrontDoorCustomDomain** | `domainIds[]` | `status.outputs.custom_domain_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `security_policy_id` | ARM resource ID of the association | Operator tooling |
| `security_policy_name` | The association's name | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Endpoint Association | 1 | Protect the endpoint's generated hostname |
| Custom Domain Association | 2 | Protect a branded custom hostname |

## Related Components

- **AzureFrontDoorFirewallPolicy** -- the WAF rule set this association enforces
- **AzureFrontDoorProfile** -- the parent container
- **AzureFrontDoorEndpoint / AzureFrontDoorCustomDomain** -- the protected hostnames
