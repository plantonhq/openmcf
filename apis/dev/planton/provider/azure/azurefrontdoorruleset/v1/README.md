# AzureFrontDoorRuleSet

A rule set inside an Azure Front Door profile: the ordered edge delivery
policy routes attach by ARM id. Each rule pairs match conditions (which
requests it applies to -- 19 condition types over paths, headers,
cookies, methods, addresses, TLS) with actions (redirects, rewrites,
header edits, or a per-request override of the route's caching and
forwarding).

The rules live inside the set rather than as standalone resources: they
form one ordered policy document, nothing references an individual
rule, and a rule is meaningless outside its set's evaluation order. One
set is commonly attached to many routes -- which is exactly why the SET
is its own composable resource.

## When to Use

Use AzureFrontDoorRuleSet when you need:

- **A security-header baseline** stamped on every response at the edge,
  once, for every route that attaches the set
- **HTTPS upgrades, apex-to-www redirects, or vanity-URL rewrites**
  answered at the edge without touching a backend
- **Per-request cache and origin overrides** -- disable caching for
  `/api/*`, force long TTLs for fingerprinted assets, or steer
  cookie-flagged traffic to a canary origin group

## Key Configuration

- `profile_id` -- the parent profile, referenced from an
  AzureFrontDoorProfile's output; fixed at creation
- `rule_set_name` -- 1-60 letters/digits (no hyphens), unique within the
  profile; ForceNew (replaces the set and every rule)
- `rules[]` -- each with a unique `name` (its ARM identity), an `order`
  (ascending evaluation), optional `conditions` (all must match; at most
  10), and required `actions` (at least one, at most 5; a redirect and a
  rewrite never coexist)
- `behavior_on_match` -- CONTINUE (default) keeps evaluating later
  rules; STOP is right for terminal actions like redirects

## Composition

```yaml
profileId:
  valueFrom:
    kind: AzureFrontDoorProfile
    name: my-front-door
    fieldPath: status.outputs.profile_id
```

Routes attach the whole set through its `rule_set_id` output (the
route's `rule_set_ids` list).

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
