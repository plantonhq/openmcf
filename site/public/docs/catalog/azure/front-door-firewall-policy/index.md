---
title: "Front Door Firewall Policy"
description: "Front Door Firewall Policy deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorfirewallpolicy"
---

# Azure Front Door Firewall Policy

Deploys a Web Application Firewall (WAF) policy for Azure Front Door -- the rule set Microsoft's edge enforces on HTTP traffic before requests ever reach an origin. Custom match and rate-limit rules evaluate first by ascending priority, Microsoft's managed rule sets (PREMIUM) evaluate next with three-scope tuning, and policy settings decide the enforcement mode, body inspection, block-response customization, and access-log scrubbing. The policy is GLOBAL and resource-group-scoped, and it enforces nothing on its own: an AzureFrontDoorSecurityPolicy associates it with a profile's domains, and that association turns the WAF on. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door WAF Policy** -- a global `Microsoft.Network/frontDoorWebApplicationFirewallPolicies` resource
- **Custom rules** -- your IP/geo allowlists, header exceptions, and per-client rate limits, evaluated first
- **Managed rules** (PREMIUM) -- Microsoft_DefaultRuleSet and Bot Manager with set/group/rule-scoped tuning
- **Policy settings** -- DETECTION/PREVENTION mode, body inspection, block-response customization, log scrubbing
- **Azure Tags** -- resource metadata tags merged with user tags

## The Policy in the Front Door Family

- **AzureFrontDoorSecurityPolicy** -- references this policy's `firewall_policy_id` and names the domains to protect; without it the policy sits idle
- **AzureFrontDoorProfile** -- the sku must match: STANDARD policies pair with STANDARD profiles, PREMIUM with PREMIUM
- **This is the FRONT DOOR policy type** -- the regional WAF that Application Gateways attach (AzureWebApplicationFirewallPolicy) is a different resource with a different rule vocabulary

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **A resource group** for the policy (the policy itself is global -- the group scopes RBAC and lifecycle).
- **PREMIUM Front Door profiles** if you plan managed rules or the JS-challenge/CAPTCHA actions -- the sku pairing is enforced at association time.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Firewall Policy**, and click **Deploy**. The wizard walks you through placement, the tier and enforcement mode (DETECTION preselected for the safe rollout), inspection and block-response behavior, the custom-rules builder, the PREMIUM-gated managed-rules tuning tree, and log scrubbing. Start from the **Standard Rate Limit** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorFirewallPolicy
metadata:
  name: edge-waf
  org: acme-corp
  env: prod
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: security-rg
      fieldPath: status.outputs.resource_group_name
  policyName: acmeedgewaf
  mode: DETECTION
  customRules:
    - name: LoginRateLimit
      ruleType: RATE_LIMIT_RULE
      rateLimitThreshold: 20
      action: BLOCK
      matchConditions:
        - matchVariable: REQUEST_URI
          operator: BEGINS_WITH
          matchValues:
            - /login
```

```shell
planton apply -f front-door-firewall-policy.yaml
```

This creates a STANDARD policy in DETECTION mode. Attach it to domains with an AzureFrontDoorSecurityPolicy referencing the `firewall_policy_id` output; flip the mode to PREVENTION once the match log looks right.

### InfraChart

```yaml
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: security-rg
      fieldPath: status.outputs.resource_group_name
  policyName: acmeedgewafpremium
  sku: PREMIUM
  mode: PREVENTION
  managedRules:
    - type: Microsoft_DefaultRuleSet
      version: "2.1"
      action: RULE_SET_BLOCK
    - type: Microsoft_BotManagerRuleSet
      version: "1.1"
      action: RULE_SET_BLOCK
```

## Key Configuration

**Policy name** -- 1-128 characters, letter-first, letters and numbers ONLY (no hyphens -- the Front Door family outlier). ForceNew: renaming detaches the policy from every security policy referencing it.

**SKU** -- unspecified deploys STANDARD. PREMIUM unlocks managed rules and the JS-challenge/CAPTCHA actions, must pair with PREMIUM profiles, and Azure refuses a PREMIUM → STANDARD downgrade outright.

**Mode** -- DETECTION logs matches (the tuning posture for new policies); PREVENTION blocks/redirects/challenges (the production posture). Updates in place.

**Custom rules** -- up to 100, evaluated by ascending priority; conditions AND together, values within a condition OR. RATE_LIMIT rules count per-client matches in a sliding window.

**Managed rules** (PREMIUM) -- Microsoft_DefaultRuleSet 1.1/2.0/2.1 (2.x uses anomaly scoring: overrides may only substitute ANOMALY_SCORING or LOG), Microsoft_BotManagerRuleSet 1.0/1.1 (the only home of the JS-challenge override), legacy DefaultRuleSet 1.0. A rule listed in an override is DISABLED unless `enabled: true` -- the false-positive tuning gesture.

**Log scrubbing** (preview) -- masks auth headers, PII arguments, or client IPs before WAF logs land in storage. IP address and URI accept only the every-key scope.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewall_policy_id` | ARM resource ID of the WAF policy | AzureFrontDoorSecurityPolicy.`firewallPolicyId` |
| `firewall_policy_name` | The policy's name | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Standard Rate Limit | 1 | STANDARD tier with a per-client rate-limit rule |
| Premium Managed Prevention | 2 | PREMIUM with Microsoft's managed rule sets in PREVENTION |
| Detection Rollout | 3 | The tune-first posture -- rules ship logging-only |

## Related Components

- **AzureFrontDoorSecurityPolicy** -- attaches this policy to domains (the enforcement switch)
- **AzureFrontDoorProfile** -- the sku-matched delivery container
- **AzureFrontDoorCustomDomain / AzureFrontDoorEndpoint** -- the domains a security policy scopes this WAF to
