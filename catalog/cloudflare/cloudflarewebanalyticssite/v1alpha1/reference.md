# CloudflareWebAnalyticsSite

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

CloudflareWebAnalyticsSiteSpec creates a Web Analytics (RUM) site:
Cloudflare's privacy-first real-user monitoring for one website, measured
by a JavaScript beacon. The site is identified by either a hostname (any
site, Cloudflare-proxied or not -- you embed the snippet yourself) or a
Cloudflare zone (auto_install can then inject the snippet at the edge) --
exactly one must be set. Free on every plan.

The optional `rules[]` narrow WHAT gets measured: include or exclude
traffic by host and path. Cloudflare manages rules as separate objects
under the site's ruleset; this kind folds them in, and the modules manage
one rule object per row.

One provider truth this spec teaches: the provider never reads rules back
after writing them (its refresh is deliberately blind), so rule edits
made in the Cloudflare dashboard are invisible to IaC until the next
apply re-asserts the declared rows. Manage rules from here or from the
dashboard -- not both.

## Example

```yaml
# Complete example manifest for CloudflareWebAnalyticsSite.
# Measures a hostname with the full beacon, excluding checkout pages from
# measurement (three resources: the site and two folded rules).
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWebAnalyticsSite
metadata:
  name: www-example-rum
spec:
  account_id: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  host: www.example.com
  auto_install: false
  rules:
    - host: www.example.com
      paths:
        - "/checkout/*"
      inclusive: false
    - paths:
        - "/*"
      inclusive: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.accountId` | `string` | yes |  |  |
| `spec.host` | `string` |  |  |  |
| `spec.zoneTag` | `string \| valueFrom` |  |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.autoInstall` | `bool` |  |  |  |
| `spec.enabled` | `bool` |  |  |  |
| `spec.lite` | `bool` |  |  |  |
| `spec.rules` | `[]CloudflareWebAnalyticsSiteRule` |  |  |  |
| `spec.rules[].host` | `string` |  |  |  |
| `spec.rules[].paths` | `[]string` |  |  |  |
| `spec.rules[].inclusive` | `bool` |  |  |  |
| `spec.rules[].isPaused` | `bool` |  |  |  |

## Field Details

### spec.accountId

`string` · required

The Cloudflare account the site belongs to.

- rule: account_id must be a 32-character hex string
- rule: {"required":true}

### spec.host

`string`

The hostname to measure (e.g. www.example.com). Use for sites not on
Cloudflare, or when you embed the snippet yourself. Set this OR
zone_tag, never both.

### spec.zoneTag

`string | valueFrom`

The Cloudflare zone to measure. Set this OR host, never both.
When using value_from, defaults to CloudflareDnsZone kind and status.outputs.zone_id field path.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.autoInstall

`bool` · optional (explicit presence)

Whether Cloudflare injects the measurement snippet automatically at the
edge (orange-clouded zones only -- pair with zone_tag).

### spec.enabled

`bool` · optional (explicit presence)

Whether measurement is active. Cloudflare's default is enabled; the
modules send the flag only when set.

### spec.lite

`bool` · optional (explicit presence)

Serve the lightweight beacon variant (smaller script, reduced metric
set).

### spec.rules

`[]CloudflareWebAnalyticsSiteRule`

Include/exclude rules narrowing what gets measured. Cloudflare manages
each as its own object under the site's ruleset; the modules keep one
object per row, in order.

### spec.rules[].host

`string`

The hostname the rule applies to (e.g. shop.example.com; empty applies
to every host of the site).

### spec.rules[].paths

`[]string`

The paths the rule applies to (e.g. /checkout/*; * wildcards allowed).

### spec.rules[].inclusive

`bool` · optional (explicit presence)

Whether matching traffic is measured (true, an include rule) or
dropped from measurement (false, an exclude rule).

### spec.rules[].isPaused

`bool` · optional (explicit presence)

Whether the rule is currently paused (kept but not applied).

## Validation Rules

- `spec.site_identity_exactly_one`: set exactly one of host (measure any site by hostname) or zone_tag (measure a Cloudflare zone)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareWebAnalyticsSite, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.site_tag` | `string` | The Cloudflare-assigned site tag (the site's identity in every RUM API path). |
| `status.outputs.site_token` | `string` | The site's measurement token, embedded by the JavaScript beacon. Treated as a credential in these outputs (marked secret by the modules) even though it ships inside public pages once deployed. |
| `status.outputs.snippet` | `string` | The ready-to-embed JavaScript snippet (carries the site token). Secret-marked by the modules for the same reason as site_token. |
| `status.outputs.ruleset_id` | `string` | The site's ruleset ID -- the parent object Cloudflare stores the include/exclude rules under. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneTag` | CloudflareDnsZone | `status.outputs.zone_id` |

## See Also

- [Overview](../README.md)
