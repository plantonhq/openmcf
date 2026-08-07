# Azure Front Door Rule Set

Deploys a rule set inside an Azure Front Door (Standard/Premium) profile -- the delivery-rules engine that runs at Microsoft's edge. A rule set is an ordered collection of rules; each rule pairs match conditions (what traffic it applies to -- up to 10 per rule, ANDed together) with actions (what happens -- header changes, a redirect or rewrite, a route override; 1 to 5 per rule). The set runs against no traffic on its own: AzureFrontDoorRoute resources attach it via `ruleSetIds`, and its rules then evaluate on every request those routes match. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Rule Set** -- a named child of the profile
- **Delivery rules** -- each with conditions (19 request parts: path, headers, cookies, device type, scheme, geo/IP…), actions (request/response header changes, redirect XOR rewrite, route-configuration override), an evaluation order, and CONTINUE/STOP behavior

## The Rule Set in the Front Door Family

- **AzureFrontDoorProfile** -- the parent container, referenced by `profileId`
- **AzureFrontDoorRoute** -- attaches this set via `ruleSetIds`; the route side owns the attachment
- **AzureFrontDoorOriginGroup** -- a rule's route override can steer matches to a different origin group (the canary gesture), referenced by `originGroupId`

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Profile** the rule set nests under.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Rule Set**, and click **Deploy**. The wizard walks you through the profile and name, then the rules builder: pick a condition type (each offers only its legal operators), pair it with actions, and let the live counters hold Azure's 10-condition and 5-action caps. Quick-add templates carry the classic stories. Start from the **Security Headers** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorRuleSet
metadata:
  name: security-headers
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  ruleSetName: securityHeaders
  rules:
    - name: addSecurityHeaders
      order: 1
      actions:
        responseHeaders:
          - headerAction: OVERWRITE
            headerName: Strict-Transport-Security
            value: "max-age=31536000; includeSubDomains"
          - headerAction: OVERWRITE
            headerName: X-Content-Type-Options
            value: nosniff
```

```shell
planton apply -f front-door-rule-set.yaml
```

A rule with no conditions matches every request its routes serve -- the deliberate shape for policy-for-everything rules like security headers. Attach the set by adding its `rule_set_id` output to a route's `ruleSetIds`.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  ruleSetName: canarySteering
  rules:
    - name: toCanary
      order: 5
      behaviorOnMatch: STOP
      conditions:
        requestHeader:
          - headerName: X-Canary
            operator: EQUAL
            matchValues:
              - "true"
      actions:
        routeConfigurationOverride:
          originGroupId:
            valueFrom:
              kind: AzureFrontDoorOriginGroup
              name: canary-origins
              fieldPath: status.outputs.origin_group_id
          forwardingProtocol: MATCH_REQUEST
          cacheBehavior: DISABLED
```

## Key Configuration

**Rule set name** -- 1-60 characters, letter-first, letters and numbers only (no hyphens). ForceNew: renaming replaces the set, stranding every attached route until it repoints.

**Rules** -- unique names (1-260 chars, no hyphens), ascending evaluation order (leave gaps: 10, 20, 30), CONTINUE (later rules still run) or STOP on match. Rules update IN PLACE -- tuning never replaces the set.

**Conditions** -- up to 10 per rule across all 19 types; conditions AND together, the values within one condition OR. Each type accepts only its own operators (URL path adds WILDCARD; remote address takes GEO_MATCH/IP_MATCH; body and file extension always need values).

**Actions** -- 1 to 5 per rule (each header action counts). One redirect XOR rewrite. The route-configuration override steers matches to a different origin group and/or overrides caching -- duration format `[days.]HH:MM:SS`, and the ignore/include-specified query-string behaviors need their parameter list.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |
| **AzureFrontDoorOriginGroup** (route overrides) | `rules[].actions.routeConfigurationOverride.originGroupId` | `status.outputs.origin_group_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `rule_set_id` | ARM resource ID of the rule set | AzureFrontDoorRoute.`ruleSetIds` |
| `rule_set_name` | The set's name within its profile | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Security Headers | 1 | HSTS + nosniff on every request |
| HTTPS and Caching | 2 | Scheme enforcement with cache overrides |
| Path Rewrite and Canary | 3 | Prefix rewrites plus header-steered canary routing |

## Related Components

- **AzureFrontDoorProfile** -- the parent container
- **AzureFrontDoorRoute** -- attaches this set to traffic
- **AzureFrontDoorOriginGroup** -- the canary steering target
