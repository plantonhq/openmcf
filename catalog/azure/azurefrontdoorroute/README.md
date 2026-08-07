# AzureFrontDoorRoute

A route inside an Azure Front Door endpoint: the rule that connects
client traffic to an origin group by URL pattern, with protocol policy
and edge caching. Routes are the traffic-serving edge of the Front Door
graph:

```
endpoint (entry hostname) -> route (match + policy) -> origin group -> origins
```

Routes are many-per-endpoint with independent lifecycles -- one
endpoint commonly serves "/api/*" and "/static/*" from different
backends -- which is why the route is a first-class kind referencing
the endpoint rather than a list folded into it.

## When to Use

Use AzureFrontDoorRoute when you need:

- **To serve traffic at all** -- an endpoint without routes answers
  nothing
- **Path-based backend splits** -- API traffic to one origin group,
  static assets to another, on the same hostname
- **Edge caching and compression** -- per-route cache policy is where
  CDN economics live

## Key Configuration

- `endpoint_id` -- the ARM parent (ForceNew), referenced from an
  AzureFrontDoorEndpoint's output
- `origin_group_id` -- the destination pool; updatable in place
  (repointing routes is how traffic moves between pools)
- `origin_ids` -- deploy-ordering references to the group's origins;
  never sent to Azure, but ARM rejects a route whose group has no
  origins yet
- `patterns_to_match` -- "/*", "/api/*", ...; Front Door picks the most
  specific pattern across the endpoint's routes
- `supported_protocols` + `https_redirect_enabled` -- both protocols
  with the redirect (default) is the standard production posture; the
  redirect requires both (spec-enforced)
- `forwarding_protocol` -- the origin-leg protocol; HTTPS_ONLY keeps
  the origin leg encrypted regardless of the client protocol
- `rule_set_ids` -- the AzureFrontDoorRuleSet delivery policies applied
  to this route's traffic (redirects, header edits, cache overrides)
- `custom_domain_ids` + `link_to_default_domain` -- the hostnames the
  route serves: validated AzureFrontDoorCustomDomain references, the
  generated *.azurefd.net hostname (default), or both; disabling the
  default domain requires at least one custom domain (spec-enforced)
- `cache` -- ABSENT means caching disabled (a real switch); configure
  query-string keying and compression when present

## Composition

```yaml
endpointId:
  valueFrom:
    kind: AzureFrontDoorEndpoint
    name: web-endpoint
    fieldPath: status.outputs.endpoint_id
originGroupId:
  valueFrom:
    kind: AzureFrontDoorOriginGroup
    name: api-backends
    fieldPath: status.outputs.origin_group_id
```

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
