---
title: "Web Application Firewall Policy"
description: "Web Application Firewall Policy deployment documentation"
icon: "package"
order: 100
componentName: "azurewebapplicationfirewallpolicy"
---

# Azure Web Application Firewall Policy

Creates a regional Web Application Firewall (WAF) policy -- the rule set an Azure Application Gateway enforces on HTTP traffic. Custom rules handle allowlists, geo fencing, rate limiting, and bot challenges; Microsoft's managed rule sets (OWASP core rule set, bot manager) handle attack signatures, tuned with per-rule overrides and scoped exclusions; policy settings govern the enforcement mode and body inspection.

## What Gets Created

When you deploy an AzureWebApplicationFirewallPolicy resource, Planton provisions:

- **WAF Policy** -- an `azurerm_web_application_firewall_policy` in the specified region and resource group, carrying your custom rules, managed-rule configuration, policy settings, and tags

The policy is deliberately standalone: one policy is shared across Application Gateways and attached by reference -- gateway-wide, per listener, or per URL path rule -- so tuning it never touches the gateways.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the policy in (an `AzureResourceGroup` in composed environments)
- **A WAF_v2 Application Gateway** to attach it to (the Standard_v2 and Basic SKUs cannot enforce WAF policies)

## Quick Start

Create a file `waf-policy.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureWebApplicationFirewallPolicy
metadata:
  name: waf-baseline
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureWebApplicationFirewallPolicy.waf-baseline
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  policyName: org-waf-baseline
  managedRules:
    managedRuleSets:
      - version: "3.2"
```

Deploy:

```shell
planton apply -f waf-policy.yaml
```

This creates an OWASP 3.2 policy in Prevention mode. Read `status.outputs.policy_id` for gateway wiring.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region -- must match the gateways that attach the policy. Changing it replaces the policy. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `policyName` | `string` | The policy's name, unique within the resource group. | Required, 1-128 chars |
| `managedRules` | `object` | At least one managed rule set (see below). | Required |

### Managed Rules

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `managedRuleSets[].type` | `enum` | `OWASP` | `OWASP`, `MICROSOFT_BOT_MANAGER_RULE_SET`, or `MICROSOFT_DEFAULT_RULE_SET`. |
| `managedRuleSets[].version` | `string` | -- | OWASP: `3.2`/`3.1`/`3.0`/`2.2.9`; DefaultRuleSet: `2.1`/`2.2`; BotManager: `0.1`/`1.0`/`1.1`. |
| `managedRuleSets[].ruleGroupOverrides` | `list` | `[]` | Per-group tuning: `ruleGroupName` + rules (`id`, `enabled` -- default `false`, so listing a rule disables it -- and an optional `OVERRIDE_*` action). |
| `exclusions` | `list` | `[]` | Request parts the managed rules skip: a collection (`REQUEST_HEADER_NAMES`, `REQUEST_COOKIE_NAMES`, ...), a `SELECTOR_*` operator, the selector key, and an optional narrowed `excludedRuleSet`. |

### Custom Rules

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Letters and digits, starting with a letter. |
| `priority` | `int32` | 1-100, unique; lower runs first. |
| `ruleType` | `enum` | `MATCH_RULE` or `RATE_LIMIT_RULE`. |
| `action` | `enum` | `ALLOW`, `BLOCK`, `LOG`, `JS_CHALLENGE` (rate-limit rules cannot `ALLOW`). |
| `rateLimitDuration` / `rateLimitThreshold` / `groupRateLimitBy` | -- | Required for rate-limit rules: `ONE_MIN`/`FIVE_MINS` window, threshold, and grouping (`CLIENT_ADDR`, `CLIENT_ADDR_XFF_HEADER`, `GEO_LOCATION`, `GEO_LOCATION_XFF_HEADER`, `NONE`). |
| `matchConditions` | `list` | Variables (`REMOTE_ADDR`, `REQUEST_URI`, `REQUEST_HEADERS` + selector, ...), an operator (`IP_MATCH`, `GEO_MATCH`, `CONTAINS`, `REGEX`, ...), values, optional negation and transforms (`LOWERCASE`, `URL_DECODE`, ...). Conditions AND; values OR. |

### Policy Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | `bool` | `true` | Whether the policy is enforced at all. |
| `mode` | `enum` | `PREVENTION` | `PREVENTION` blocks; `DETECTION` only logs (the tuning mode). |
| `requestBodyCheck` / `requestBodyEnforcement` | `bool` | `true` | Body inspection and over-size blocking. |
| `requestBodyInspectLimitInKb` | `int32` | `128` | 0 = unlimited inspection. |
| `maxRequestBodySizeInKb` | `int32` | `128` | 8-2000. |
| `fileUploadEnforcement` / `fileUploadLimitInMb` | -- | `true` / `100` | Upload blocking (OWASP 3.2 only) and the 1-4000 MB limit. |
| `jsChallengeCookieExpirationInMinutes` | `int32` | `30` | 5-1440; how long a solved JS challenge stays valid. |
| `logScrubbing` | `object` | -- | Redact request parts from WAF logs: `SCRUB_REQUEST_HEADER_NAMES`, `SCRUB_REQUEST_COOKIE_NAMES`, `SCRUB_REQUEST_IP_ADDRESS`, ... with `SELECTOR_EQUALS` (one key) or `SELECTOR_EQUALS_ANY` (all keys). |

## Examples

### Rate Limiting and Geo Fencing

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureWebApplicationFirewallPolicy
metadata:
  name: waf-edge
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  policyName: org-waf-edge
  customRules:
    - name: throttleApi
      priority: 20
      ruleType: RATE_LIMIT_RULE
      action: BLOCK
      rateLimitDuration: ONE_MIN
      rateLimitThreshold: 300
      groupRateLimitBy: CLIENT_ADDR
      matchConditions:
        - matchVariables:
            - variableName: REQUEST_URI
          operator: BEGINS_WITH
          matchValues:
            - /api/
  managedRules:
    managedRuleSets:
      - version: "3.2"
```

### Attach to an Application Gateway

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureApplicationGateway
metadata:
  name: web-gateway
spec:
  sku: WAF_V2
  firewallPolicyId:
    valueFrom:
      kind: AzureWebApplicationFirewallPolicy
      name: waf-baseline
      fieldPath: status.outputs.policy_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `policy_id` | `string` | The policy's ARM ID -- referenced by Application Gateways gateway-wide, per listener, and per path rule |
| `policy_name` | `string` | The policy's name |

## Related Components

- [AzureApplicationGateway](/docs/catalog/azure/application-gateway) — the L7 load balancer that enforces the policy
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for policy placement
