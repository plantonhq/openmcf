# AzureFrontDoorSecurityPolicy

A security policy inside an Azure Front Door profile: the association
that attaches a Front Door WAF policy (AzureFrontDoorFirewallPolicy) to
the hostnames the profile serves. A WAF policy enforces nothing until a
security policy associates it -- this kind is the enforcement seam.

The `domain_ids` list names the protected hostnames: endpoint ARM IDs
(the generated `*.azurefd.net` hostname) and/or custom-domain ARM IDs
-- Azure accepts both interchangeably, and one list can mix them. The
WAF applies to all paths (`/*`) on every associated domain; Azure's
security policies accept no other pattern, so scope enforcement by
choosing WHICH domains to associate.

## When to Use

Use AzureFrontDoorSecurityPolicy when you need:

- **To turn a WAF on** -- every AzureFrontDoorFirewallPolicy needs one
  of these before it inspects any traffic
- **Different protection per hostname** -- one profile can carry
  multiple security policies, each pairing a different WAF policy with
  a different domain set (e.g. strict rules for the API domain, lighter
  rules for marketing pages)
- **One org-standard WAF across profiles** -- each profile gets its own
  association referencing the same shared policy

## Key Configuration

- `profile_id` -- the parent profile, referenced from an
  AzureFrontDoorProfile's output; ForceNew
- `security_policy_name` -- unique within the profile;
  letters/digits/hyphens, begins and ends alphanumeric; ForceNew
- `firewall_policy_id` -- the WAF policy to enforce, referenced from an
  AzureFrontDoorFirewallPolicy's output; ForceNew (swapping the policy
  replaces the association -- a fast, metadata-only operation). The WAF
  policy's sku must MATCH the profile's sku.
- `domain_ids` -- 1-500 endpoint/custom-domain references (a STANDARD
  profile caps the list at 100); updates in place

## Composition

```yaml
firewallPolicyId:
  valueFrom:
    kind: AzureFrontDoorFirewallPolicy
    name: my-edge-waf
    fieldPath: status.outputs.firewall_policy_id
domainIds:
  - valueFrom:
      kind: AzureFrontDoorEndpoint
      name: my-endpoint
      fieldPath: status.outputs.endpoint_id
```

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
