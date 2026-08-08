# AzureApplicationInsightsStandardWebTest - Pulumi Module

Pulumi implementation for the AzureApplicationInsightsStandardWebTest
deployment component.

## Architecture

```
appinsights.StandardWebTest (one synthetic availability test)
```

## Key Design Decisions

- **Optional request and validation fields are presence-guarded** --
  unspecified specs omit the argument entirely so Azure applies its own
  defaults (GET, follow redirects, 5-minute frequency, status 200) and
  both engines send identical request bodies.
- **`geo_locations` are the run-from locations**, passed through as
  Azure location IDs; the test's reporting target is the referenced
  Application Insights component.
- **Validation rules are built only when the spec carries them**,
  including the nested content-match block.
- **PARITY-EXCEPTION on tag shape** (documented in both engines): the
  Terraform module writes the snake-case `resource_kind` literal and
  falls back to `metadata.name` for `resource_id`, while this module
  lowers the kind enum string and omits an empty id -- output-neutral,
  tags only.

## Provider

Built via the shared `pulumiazureprovider.Get` builder -- static client
secret, keyless web identity, or ambient chain, resolved from the stack
input. Never construct the provider inline.
