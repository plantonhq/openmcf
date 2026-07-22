# AzureFrontDoorEndpoint

An endpoint inside an Azure Front Door profile: the public entry point
client traffic arrives at. Each endpoint receives a generated, globally
unique hostname (`{name}-{hash}.z01.azurefd.net`); routes attach to the
endpoint to define which URL patterns it serves, and custom-domain DNS
records CNAME onto its hostname.

Endpoints are many-per-profile with independent lifecycles -- one
profile commonly fronts several applications, each behind its own
endpoint -- which is why the endpoint is a first-class kind referencing
the profile rather than a list folded into the profile's spec.

## When to Use

Use AzureFrontDoorEndpoint when you need:

- **A public entry hostname** for an application behind Front Door
- **Per-app isolation on a shared profile** -- each app gets its own
  endpoint, routes, and (later) custom domains
- **A launch/maintenance switch** -- disabling stops traffic at the
  edge without deleting configuration

## Key Configuration

- `profile_id` -- the parent profile, referenced from an
  AzureFrontDoorProfile's output; fixed at creation
- `endpoint_name` -- 2-46 characters; the prefix of the generated
  hostname, so renaming breaks DNS records pointing at the old one;
  the hash suffix means the name only needs per-profile uniqueness
- `enabled` -- default true; false stops traffic at the edge
- `tags` -- merged over the Planton-derived resource tags

## Composition

```yaml
profileId:
  valueFrom:
    kind: AzureFrontDoorProfile
    name: my-front-door
    fieldPath: status.outputs.profile_id
```

Routes reference this endpoint through its `endpoint_id` output; DNS
records (AzureDnsRecord CNAME) reference its `host_name` output.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
