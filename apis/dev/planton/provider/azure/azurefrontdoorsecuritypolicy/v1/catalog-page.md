# Azure Front Door Security Policy

Creates a security policy inside an AzureFrontDoorProfile -- the association that attaches an AzureFrontDoorFirewallPolicy (the Front Door WAF) to the endpoint and custom-domain hostnames the profile serves. The WAF enforces nothing until this association exists.

## What Gets Created

When you deploy an AzureFrontDoorSecurityPolicy resource, Planton provisions:

- **Front Door Security Policy** -- an `azurerm_cdn_frontdoor_security_policy` on the referenced profile, binding the referenced WAF policy to the referenced domains on all paths (`/*`)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** (referenced through `profileId`)
- **An AzureFrontDoorFirewallPolicy** with the SAME sku as the profile (referenced through `firewallPolicyId`)
- **At least one AzureFrontDoorEndpoint or AzureFrontDoorCustomDomain** to protect (referenced through `domainIds`)

## Quick Start

Create a file `security-policy.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorSecurityPolicy
metadata:
  name: waf-attach
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorSecurityPolicy.waf-attach
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
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
        name: my-endpoint
        fieldPath: status.outputs.endpoint_id
```

Deploy it:

```shell
planton apply -f security-policy.yaml
```

## Notes

- `domainIds` accepts endpoint IDs (the generated `*.azurefd.net` hostname) and custom-domain IDs interchangeably; one list can mix them.
- A STANDARD profile allows up to 100 domains per association, PREMIUM up to 500 (checked at deploy time).
- The WAF applies to all paths on every associated domain -- Azure accepts no other pattern; scope enforcement by choosing which domains to associate.
