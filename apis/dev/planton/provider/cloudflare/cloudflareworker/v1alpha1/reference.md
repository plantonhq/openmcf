# CloudflareWorker

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1alpha1`

CloudflareWorkerSpec deploys a Cloudflare Worker script and everything that
hangs off it: its resource bindings, routing (workers.dev, custom domains,
routes), scheduled (cron) invocations, and runtime settings. Bindings are
grouped by type (the wrangler.toml authoring grain) and each cross-resource
binding accepts a literal value or a reference to the producing resource, so a
Worker composes as a real node in the resource graph.

## Example

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWorker
metadata:
  name: test-worker
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  workerName: test-worker
  content: |
    export default {
      async fetch(request, env, ctx) {
        return new Response("ok");
      }
    };
  compatibilityDate: "2025-01-15"
  compatibilityFlags:
    - nodejs_compat
  vars:
    LOG_LEVEL: info
  workersDev:
    enabled: true
  observability:
    enabled: true
    headSamplingRate: 1
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.accountId` | `string` | yes |  |  |
| `spec.workerName` | `string` | yes |  |  |
| `spec.compatibilityDate` | `string` |  |  |  |
| `spec.content` | `string` |  |  |  |
| `spec.r2Bundle` | `CloudflareWorkerScriptBundle` |  |  |  |
| `spec.r2Bundle.bucket` | `string` | yes |  |  |
| `spec.r2Bundle.path` | `string` | yes |  |  |
| `spec.mainModule` | `string` |  | `index.js` |  |
| `spec.compatibilityFlags` | `[]string` |  |  |  |
| `spec.vars` | `map<string, string>` |  |  |  |
| `spec.secrets` | `[]CloudflareWorkerSecretBinding` |  |  |  |
| `spec.secrets[].name` | `string` | yes |  |  |
| `spec.secrets[].value` | `string \| valueFrom` (sensitive) | yes |  |  |
| `spec.kvNamespaces` | `[]CloudflareWorkerKvBinding` |  |  |  |
| `spec.kvNamespaces[].name` | `string` | yes |  |  |
| `spec.kvNamespaces[].namespaceId` | `string \| valueFrom` | yes |  | CloudflareKvNamespace (`status.outputs.namespace_id`) |
| `spec.r2Buckets` | `[]CloudflareWorkerR2Binding` |  |  |  |
| `spec.r2Buckets[].name` | `string` | yes |  |  |
| `spec.r2Buckets[].bucketName` | `string \| valueFrom` | yes |  | CloudflareR2Bucket (`status.outputs.bucket_name`) |
| `spec.r2Buckets[].jurisdiction` | `string` |  |  |  |
| `spec.d1Databases` | `[]CloudflareWorkerD1Binding` |  |  |  |
| `spec.d1Databases[].name` | `string` | yes |  |  |
| `spec.d1Databases[].databaseId` | `string \| valueFrom` | yes |  | CloudflareD1Database (`status.outputs.database_id`) |
| `spec.hyperdriveConfigs` | `[]CloudflareWorkerHyperdriveBinding` |  |  |  |
| `spec.hyperdriveConfigs[].name` | `string` | yes |  |  |
| `spec.hyperdriveConfigs[].configId` | `string \| valueFrom` | yes |  | CloudflareHyperdriveConfig (`status.outputs.hyperdrive_id`) |
| `spec.services` | `[]CloudflareWorkerServiceBinding` |  |  |  |
| `spec.services[].name` | `string` | yes |  |  |
| `spec.services[].service` | `string \| valueFrom` | yes |  | CloudflareWorker (`status.outputs.script_name`) |
| `spec.services[].environment` | `string` |  |  |  |
| `spec.services[].entrypoint` | `string` |  |  |  |
| `spec.queues` | `[]CloudflareWorkerQueueBinding` |  |  |  |
| `spec.queues[].name` | `string` | yes |  |  |
| `spec.queues[].queueName` | `string \| valueFrom` | yes |  | CloudflareQueue (`status.outputs.queue_name`) |
| `spec.durableObjects` | `[]CloudflareWorkerDurableObjectBinding` |  |  |  |
| `spec.durableObjects[].name` | `string` | yes |  |  |
| `spec.durableObjects[].className` | `string` | yes |  |  |
| `spec.durableObjects[].scriptName` | `string` |  |  |  |
| `spec.durableObjects[].environment` | `string` |  |  |  |
| `spec.analyticsEngineDatasets` | `[]CloudflareWorkerAnalyticsEngineBinding` |  |  |  |
| `spec.analyticsEngineDatasets[].name` | `string` | yes |  |  |
| `spec.analyticsEngineDatasets[].dataset` | `string` | yes |  |  |
| `spec.vectorizeIndexes` | `[]CloudflareWorkerVectorizeBinding` |  |  |  |
| `spec.vectorizeIndexes[].name` | `string` | yes |  |  |
| `spec.vectorizeIndexes[].indexName` | `string` | yes |  |  |
| `spec.ai` | `[]CloudflareWorkerAiBinding` |  |  |  |
| `spec.ai[].name` | `string` | yes |  |  |
| `spec.versionMetadata` | `[]CloudflareWorkerVersionMetadataBinding` |  |  |  |
| `spec.versionMetadata[].name` | `string` | yes |  |  |
| `spec.workersDev` | `CloudflareWorkerWorkersDev` |  |  |  |
| `spec.workersDev.enabled` | `bool` |  |  |  |
| `spec.workersDev.previewsEnabled` | `bool` |  |  |  |
| `spec.customDomains` | `[]CloudflareWorkerCustomDomain` |  |  |  |
| `spec.customDomains[].hostname` | `string` | yes |  |  |
| `spec.customDomains[].zoneId` | `string \| valueFrom` |  |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.routes` | `[]CloudflareWorkerRoute` |  |  |  |
| `spec.routes[].zoneId` | `string \| valueFrom` | yes |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.routes[].pattern` | `string` | yes |  |  |
| `spec.schedules` | `[]string` |  |  |  |
| `spec.observability` | `CloudflareWorkerObservability` |  |  |  |
| `spec.observability.enabled` | `bool` |  |  |  |
| `spec.observability.headSamplingRate` | `double` |  |  |  |
| `spec.placement` | `CloudflareWorkerPlacement` |  |  |  |
| `spec.placement.mode` | `string` |  |  |  |
| `spec.limits` | `CloudflareWorkerLimits` |  |  |  |
| `spec.limits.cpuMs` | `int64` |  |  |  |
| `spec.limits.subrequests` | `int64` |  |  |  |
| `spec.logpush` | `bool` |  |  |  |
| `spec.tailConsumers` | `[]CloudflareWorkerTailConsumer` |  |  |  |
| `spec.tailConsumers[].service` | `string` | yes |  |  |
| `spec.tailConsumers[].environment` | `string` |  |  |  |
| `spec.tailConsumers[].namespace` | `string` |  |  |  |
| `spec.assets` | `CloudflareWorkerAssets` |  |  |  |
| `spec.assets.directory` | `string` | yes |  |  |
| `spec.assets.config` | `CloudflareWorkerAssetsConfig` |  |  |  |
| `spec.assets.config.htmlHandling` | `string` |  |  |  |
| `spec.assets.config.notFoundHandling` | `string` |  |  |  |
| `spec.assets.config.headers` | `string` |  |  |  |
| `spec.assets.config.redirects` | `string` |  |  |  |
| `spec.assets.config.runWorkerFirst` | `bool` |  |  |  |
| `spec.assets.config.runWorkerFirstRules` | `[]string` |  |  |  |
| `spec.assets.bindingName` | `string` |  |  |  |

## Field Details

### spec.accountId

`string` · required

The Cloudflare account ID in which to create the worker.

- rule: {"required":true,"string":{"len":"32","pattern":"^[0-9a-fA-F]{32}$"}}

### spec.workerName

`string` · required

The Worker script name, visible in the Cloudflare dashboard and used as the
target of service bindings and routes.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63"}}

### spec.compatibilityDate

`string`

Compatibility date (YYYY-MM-DD) pinning the Workers runtime behavior. When
unset, the module pins today's date at deploy.

- rule: {"string":{"pattern":"^([0-9]{4}-[0-9]{2}-[0-9]{2})?$"}}

### spec.content

`string`

Inline ES-module source. Convenient for small or generated scripts.

### spec.r2Bundle

`CloudflareWorkerScriptBundle`

A pre-built script bundle stored in an R2 bucket (the CI build-artifact
flow). The module fetches the object and deploys it.

### spec.r2Bundle.bucket

`string` · required

The R2 bucket name where the script bundle is stored.

- rule: {"required":true}

### spec.r2Bundle.path

`string` · required

The object key (path) of the script bundle within the bucket.

- rule: {"required":true}

### spec.mainModule

`string`

The entrypoint module filename within the bundle (module-format workers).

- default: `index.js`

### spec.compatibilityFlags

`[]string`

Runtime compatibility flags (e.g. "nodejs_compat"). See the Cloudflare
compatibility-flags docs for the set valid on a given compatibility_date.

### spec.vars

`map<string, string>`

Plain-text variable bindings: non-secret configuration exposed to the Worker
as env vars. Map key is the JS binding name, value is the literal string.

### spec.secrets

`[]CloudflareWorkerSecretBinding`

Secret bindings: sensitive values exposed to the Worker as env vars. Provide
each as a managed-secret reference; resolved just-in-time at deploy.

### spec.secrets[].name

`string` · required

The JS binding (variable) name the Worker reads the secret through.

- rule: {"required":true}

### spec.secrets[].value

`string | valueFrom` · required · sensitive

The secret value. Provide a managed-secret reference; the platform resolves
it just-in-time at deploy and never stores plaintext.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.kvNamespaces

`[]CloudflareWorkerKvBinding`

KV namespace bindings. Each binds a CloudflareKvNamespace (by id or
reference) to a JS variable the Worker uses for edge key-value access.

### spec.kvNamespaces[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.kvNamespaces[].namespaceId

`string | valueFrom` · required

The KV namespace ID, or a reference to a CloudflareKvNamespace resource.

- references: CloudflareKvNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareKvNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.r2Buckets

`[]CloudflareWorkerR2Binding`

R2 bucket bindings. Each binds a CloudflareR2Bucket (by name or reference)
for object storage access from the Worker.

### spec.r2Buckets[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.r2Buckets[].bucketName

`string | valueFrom` · required

The R2 bucket name, or a reference to a CloudflareR2Bucket resource.

- references: CloudflareR2Bucket (`status.outputs.bucket_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareR2Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_name}} -- a bare string does not parse

### spec.r2Buckets[].jurisdiction

`string`

Optional data-residency jurisdiction of the bucket. One of "eu", "fedramp",
or "fedramp-high"; leave empty for the default jurisdiction.

- rule: jurisdiction must be one of "eu", "fedramp", "fedramp-high"

### spec.d1Databases

`[]CloudflareWorkerD1Binding`

D1 database bindings. Each binds a CloudflareD1Database (by id or reference)
for serverless SQL access from the Worker.

### spec.d1Databases[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.d1Databases[].databaseId

`string | valueFrom` · required

The D1 database ID, or a reference to a CloudflareD1Database resource.

- references: CloudflareD1Database (`status.outputs.database_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareD1Database, name: <that resource's name>, fieldPath: status.outputs.database_id}} -- a bare string does not parse

### spec.hyperdriveConfigs

`[]CloudflareWorkerHyperdriveBinding`

Hyperdrive bindings. Each binds a CloudflareHyperdriveConfig (by id or
reference) for pooled, cached access to a regional SQL database.

### spec.hyperdriveConfigs[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.hyperdriveConfigs[].configId

`string | valueFrom` · required

The Hyperdrive config ID, or a reference to a CloudflareHyperdriveConfig resource.

- references: CloudflareHyperdriveConfig (`status.outputs.hyperdrive_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareHyperdriveConfig, name: <that resource's name>, fieldPath: status.outputs.hyperdrive_id}} -- a bare string does not parse

### spec.services

`[]CloudflareWorkerServiceBinding`

Service bindings: bind another Worker (by name or reference) for direct
worker-to-worker calls without going over the public internet.

### spec.services[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.services[].service

`string | valueFrom` · required

The target Worker's script name, or a reference to a CloudflareWorker resource.

- references: CloudflareWorker (`status.outputs.script_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareWorker, name: <that resource's name>, fieldPath: status.outputs.script_name}} -- a bare string does not parse

### spec.services[].environment

`string`

Optional environment of the target Worker to bind to.

### spec.services[].entrypoint

`string`

Optional named entrypoint (WorkerEntrypoint) to invoke on the target Worker.

### spec.queues

`[]CloudflareWorkerQueueBinding`

Queue producer bindings: let the Worker enqueue messages to a Cloudflare
Queue by name.

### spec.queues[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.queues[].queueName

`string | valueFrom` · required

The Cloudflare Queue name to produce to, or a reference to a CloudflareQueue resource.

- references: CloudflareQueue (`status.outputs.queue_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareQueue, name: <that resource's name>, fieldPath: status.outputs.queue_name}} -- a bare string does not parse

### spec.durableObjects

`[]CloudflareWorkerDurableObjectBinding`

Durable Object bindings: bind a Durable Object namespace (a class) for
strongly-consistent, stateful coordination.

### spec.durableObjects[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.durableObjects[].className

`string` · required

The Durable Object class name that implements the namespace.

- rule: {"required":true}

### spec.durableObjects[].scriptName

`string`

Script that defines the class, when it lives in a different Worker. Leave
empty when the class is defined in this Worker.

### spec.durableObjects[].environment

`string`

Optional environment of the defining script.

### spec.analyticsEngineDatasets

`[]CloudflareWorkerAnalyticsEngineBinding`

Analytics Engine bindings: bind a dataset the Worker writes analytics data points to.

### spec.analyticsEngineDatasets[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.analyticsEngineDatasets[].dataset

`string` · required

The Analytics Engine dataset the Worker writes data points to.

- rule: {"required":true}

### spec.vectorizeIndexes

`[]CloudflareWorkerVectorizeBinding`

Vectorize bindings: bind a vector-search index for AI/embedding workloads.

### spec.vectorizeIndexes[].name

`string` · required

The JS binding (variable) name.

- rule: {"required":true}

### spec.vectorizeIndexes[].indexName

`string` · required

The Vectorize index name.

- rule: {"required":true}

### spec.ai

`[]CloudflareWorkerAiBinding`

Workers AI bindings: bind the account's AI inference gateway. Each entry is
just the JS variable name the Worker calls models through.

### spec.ai[].name

`string` · required

The JS binding (variable) name the Worker calls AI models through.

- rule: {"required":true}

### spec.versionMetadata

`[]CloudflareWorkerVersionMetadataBinding`

Version-metadata bindings: expose the deployed version id/tag to the Worker
at runtime via the named binding.

### spec.versionMetadata[].name

`string` · required

The JS binding (variable) name exposing the version id/tag.

- rule: {"required":true}

### spec.workersDev

`CloudflareWorkerWorkersDev`

workers.dev subdomain exposure (e.g. <name>.<subdomain>.workers.dev).

### spec.workersDev.enabled

`bool`

Expose the Worker at <name>.<account-subdomain>.workers.dev.

### spec.workersDev.previewsEnabled

`bool`

Also enable per-version preview URLs (<version>-<name>.<subdomain>.workers.dev).

### spec.customDomains

`[]CloudflareWorkerCustomDomain`

Custom domains that route directly to this Worker (a managed hostname with
automatic TLS), distinct from pattern-based routes.

### spec.customDomains[].hostname

`string` · required

The fully qualified hostname to route to the Worker (e.g. "api.example.com").
Cloudflare provisions and manages TLS for it automatically.

- rule: {"required":true}

### spec.customDomains[].zoneId

`string | valueFrom`

The Cloudflare Zone ID hosting the hostname, or a reference to a
CloudflareDnsZone. Optional — Cloudflare can infer the zone from the hostname.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.routes

`[]CloudflareWorkerRoute`

Route patterns within a zone that map matching requests to this Worker.

### spec.routes[].zoneId

`string | valueFrom` · required

The Cloudflare Zone ID the route belongs to, or a reference to a CloudflareDnsZone.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.routes[].pattern

`string` · required

The route pattern (e.g. "api.example.com/webhooks/*").

- rule: {"required":true}

### spec.schedules

`[]string`

Cron schedules that invoke the Worker's scheduled handler. Each entry is a
cron expression (e.g. "0 * * * *").

### spec.observability

`CloudflareWorkerObservability`

Observability (Workers Logs) configuration.

### spec.observability.enabled

`bool`

Enable Workers Logs (structured logs visible in the dashboard / queryable).

### spec.observability.headSamplingRate

`double`

Fraction of requests sampled for logging, 0.0–1.0. Leave 0 for the default
(1.0 = sample all) when enabled.

- rule: head_sampling_rate must be between 0 and 1

### spec.placement

`CloudflareWorkerPlacement`

Smart Placement: let Cloudflare run the Worker closer to backend services it
calls, instead of closest to the user.

### spec.placement.mode

`string`

Placement mode. "smart" lets Cloudflare run the Worker near the backends it
calls. Leave empty to keep the default (run near the user).

- rule: placement mode must be "smart"

### spec.limits

`CloudflareWorkerLimits`

Per-invocation resource limits.

### spec.limits.cpuMs

`int64`

CPU time limit in milliseconds per invocation. Leave 0 for the plan default.

- rule: cpu_ms must be 0 (default) or positive

### spec.limits.subrequests

`int64`

Maximum number of subrequests the Worker may make per invocation. Leave 0
for the plan default.

- rule: subrequests must be 0 (default) or positive

### spec.logpush

`bool`

Enable Logpush for the Worker's logs (push to a configured Logpush job).

### spec.tailConsumers

`[]CloudflareWorkerTailConsumer`

Tail consumers: other Workers that receive this Worker's tail (trace) events.

### spec.tailConsumers[].service

`string` · required

The consuming Worker's service (script) name.

- rule: {"required":true}

### spec.tailConsumers[].environment

`string`

Optional environment of the consuming Worker.

### spec.tailConsumers[].namespace

`string`

Optional dispatch namespace of the consuming Worker.

### spec.assets

`CloudflareWorkerAssets`

Static assets served by this Worker (Cloudflare Workers Static Assets). Point
it at a built site directory to host a static site or single-page app at the
edge; combine it with a script source above to run a full-stack app whose
Functions handle dynamic routes while everything else is served as a static
asset. When assets are set without a script source, the Worker is a pure
static site. This is Cloudflare's converged successor to Pages for the
"build-and-upload" hosting model.

### spec.assets.directory

`string` · required

Filesystem path to the directory of built assets to upload. Interpreted on
the deploy runner — point it at your build output (e.g. "dist" or "build").

- rule: {"required":true}

### spec.assets.config

`CloudflareWorkerAssetsConfig`

Behavior for serving the assets (header injection, redirects, HTML routing,
SPA handling). Optional; sensible Cloudflare defaults apply when omitted.

- rule: set run_worker_first (apply to all paths) or run_worker_first_rules (per-path), not both

### spec.assets.config.htmlHandling

`string`

Trailing-slash and rewrite behavior for HTML requests. One of
"auto-trailing-slash", "force-trailing-slash", "drop-trailing-slash", or
"none". Empty uses Cloudflare's default.

- rule: html_handling must be one of "auto-trailing-slash", "force-trailing-slash", "drop-trailing-slash", "none"

### spec.assets.config.notFoundHandling

`string`

Response when a request matches no static asset and no Worker script handles
it. One of "none", "404-page" (serve /404.html), or "single-page-application"
(serve /index.html — the SPA fallback). Empty uses Cloudflare's default.

- rule: not_found_handling must be one of "none", "404-page", "single-page-application"

### spec.assets.config.headers

`string`

The contents of a `_headers` file: rules that attach custom headers to asset
responses. See Cloudflare's _headers syntax.

### spec.assets.config.redirects

`string`

The contents of a `_redirects` file: redirect/rewrite rules applied ahead of
asset serving. See Cloudflare's _redirects syntax.

### spec.assets.config.runWorkerFirst

`bool`

Invoke the Worker script before attempting to serve any asset (applies to all
paths). Use this when the script must run on every request. Mutually
exclusive with run_worker_first_rules. When neither is set, Cloudflare serves
a matching asset first and falls back to the script.

### spec.assets.config.runWorkerFirstRules

`[]string`

Per-path control over whether the Worker runs before asset serving. Each
entry is a path rule; "/api/*" routes through the script, "!/api/docs/*"
(negative, higher precedence) excludes a subtree. Mutually exclusive with
run_worker_first.

### spec.assets.bindingName

`string`

When set, exposes the asset namespace to the Worker script as a binding of
this JS variable name (so code can call e.g. env.ASSETS.fetch(request)).
Leave empty for a pure static site that has no script. Only meaningful when
a script source is also provided.

## Validation Rules

- `spec.code_or_assets_required`: provide a script source (content or r2_bundle) and/or a static-asset directory (assets)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareWorker, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.script_id` | `string` | The Cloudflare-assigned identifier of the deployed Worker script (equals the script name in v5). |
| `status.outputs.script_name` | `string` | The Worker script name — the target a service binding references. |
| `status.outputs.custom_domain_hostnames` | `[]string` | The custom-domain hostnames attached to this Worker. |
| `status.outputs.route_patterns` | `[]string` | The route patterns mapped to this Worker. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.kvNamespaces[].namespaceId` | CloudflareKvNamespace | `status.outputs.namespace_id` |
| `spec.r2Buckets[].bucketName` | CloudflareR2Bucket | `status.outputs.bucket_name` |
| `spec.d1Databases[].databaseId` | CloudflareD1Database | `status.outputs.database_id` |
| `spec.hyperdriveConfigs[].configId` | CloudflareHyperdriveConfig | `status.outputs.hyperdrive_id` |
| `spec.services[].service` | CloudflareWorker | `status.outputs.script_name` |
| `spec.queues[].queueName` | CloudflareQueue | `status.outputs.queue_name` |
| `spec.customDomains[].zoneId` | CloudflareDnsZone | `status.outputs.zone_id` |
| `spec.routes[].zoneId` | CloudflareDnsZone | `status.outputs.zone_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CloudflareEmailRoutingRule | `spec.action.worker` | `status.outputs.script_name` |
| CloudflareEmailRoutingZone | `spec.catchAll.worker` | `status.outputs.script_name` |
| CloudflarePagesProject | `spec.deploymentConfigs.preview.services[].service` | `status.outputs.script_name` |
| CloudflarePagesProject | `spec.deploymentConfigs.production.services[].service` | `status.outputs.script_name` |
| CloudflareQueue | `spec.consumer.scriptName` | `status.outputs.script_name` |
| CloudflareWorker | `spec.services[].service` | `status.outputs.script_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
