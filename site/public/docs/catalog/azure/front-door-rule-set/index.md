---
title: "Front Door Rule Set"
description: "Front Door Rule Set deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorruleset"
---

# Azure Front Door Rule Set

Creates a rule set inside an AzureFrontDoorProfile -- the ordered edge delivery policy (redirects, rewrites, header edits, per-request cache and origin overrides) that AzureFrontDoorRoute resources attach through their `ruleSetIds`. One set is commonly shared by many routes.

## What Gets Created

When you deploy an AzureFrontDoorRuleSet resource, Planton provisions:

- **Front Door Rule Set** -- an `azurerm_cdn_frontdoor_rule_set` on the referenced profile
- **Front Door Rules** -- one `azurerm_cdn_frontdoor_rule` per entry in `rules`, nested under the set and evaluated in ascending `order`

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** to create the rule set in (referenced through `profileId`)

## Quick Start

Create a file `rule-set.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorRuleSet
metadata:
  name: security-headers
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorRuleSet.security-headers
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
      fieldPath: status.outputs.profile_id
  ruleSetName: securityheaders
  rules:
    - name: baseline
      order: 1
      actions:
        responseHeaders:
          - headerAction: OVERWRITE
            headerName: Strict-Transport-Security
            value: max-age=31536000; includeSubDomains
          - headerAction: DELETE
            headerName: X-Powered-By
```

Deploy:

```shell
planton apply -f rule-set.yaml
```

A rule without `conditions` applies to every request -- the shape for set-wide baselines. Conditions AND together (all must match) while the values inside one condition OR together; a rule carries at most 10 conditions and 5 actions, and a redirect never coexists with a rewrite on one rule.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `rule_set_id` | The ARM id -- what AzureFrontDoorRoute references in `ruleSetIds` to attach the policy |
| `rule_set_name` | The set's name inside the profile |

## Related Resources

- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the parent profile
- [Azure Front Door Route](/docs/catalog/azure/front-door-route) -- attaches this policy to traffic
- [Azure Front Door Origin Group](/docs/catalog/azure/front-door-origin-group) -- the target of a rule's origin override
