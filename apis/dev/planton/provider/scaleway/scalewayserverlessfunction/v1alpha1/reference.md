# ScalewayServerlessFunction

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayServerlessFunctionSpec defines the specification for a
Scaleway serverless function.

This is a composite resource that creates a function namespace, the
function itself, and optional cron triggers. The namespace is an
internal implementation detail -- users interact with the function
as a single resource.

**Scaleway serverless function model:**
Scaleway organizes functions into namespaces. Each namespace can
contain multiple functions, but in the Planton model we create one
namespace per function for clean lifecycle management and isolation.
Environment variables and secrets are set on the function (not the
namespace) for simplicity and clarity.

**Code deployment:**
This resource can optionally deploy function code via `zip_file` and
`zip_hash`. When these fields are provided, the IaC module uploads
the zip and triggers deployment. Alternatively, code can be deployed
separately via the Scaleway CLI (`scw function deploy`) or CI/CD.

**Composition pattern:** Mid-tier resource (DAG Layer 2).
Upstream: `private_network_id` references ScalewayPrivateNetwork.
Downstream: `domain_name` output for ScalewayDnsRecord.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.runtime` | `string` | yes |  |  |
| `spec.handler` | `string` | yes |  |  |
| `spec.privacy` | `enum` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.memoryLimitMb` | `uint32` |  | `256` |  |
| `spec.minScale` | `uint32` |  |  |  |
| `spec.maxScale` | `uint32` |  | `20` |  |
| `spec.timeoutSeconds` | `uint32` |  | `300` |  |
| `spec.httpOption` | `enum` |  |  |  |
| `spec.env` | `ScalewayServerlessFunctionEnv` |  |  |  |
| `spec.env.variables` | `[]ScalewayServerlessFunctionEnvVar` |  |  |  |
| `spec.env.variables[].name` | `string` | yes |  |  |
| `spec.env.variables[].value` | `string` | yes |  |  |
| `spec.env.secrets` | `[]ScalewayServerlessFunctionEnvVar` |  |  |  |
| `spec.env.secrets[].name` | `string` | yes |  |  |
| `spec.env.secrets[].value` | `string` | yes |  |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.sandbox` | `string` |  |  |  |
| `spec.zipFile` | `string` |  |  |  |
| `spec.zipHash` | `string` |  |  |  |
| `spec.cronTriggers` | `[]ScalewayServerlessFunctionCronTrigger` |  |  |  |
| `spec.cronTriggers[].name` | `string` |  |  |  |
| `spec.cronTriggers[].schedule` | `string` | yes |  |  |
| `spec.cronTriggers[].args` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

region where the function namespace and function are deployed.

Scaleway serverless functions are regional resources.
Valid regions: "fr-par", "nl-ams", "pl-waw".

- rule: {"string":{"minLen":"1"}}

### spec.runtime

`string` · required

runtime specifies the language runtime for the function.

This is a plain string (not an enum) because Scaleway adds new
runtimes frequently. Using a string avoids proto staleness.

Known runtimes (as of 2026):
  Node.js: "node20", "node22"
  Python:  "python39", "python310", "python311", "python312", "python313"
  Go:      "go122", "go123", "go124"
  Rust:    "rust165"
  PHP:     "php82"

See https://www.scaleway.com/en/docs/serverless-functions/reference-content/functions-runtimes/

- rule: {"string":{"minLen":"1"}}

### spec.handler

`string` · required

handler is the function entrypoint, runtime-dependent.

Examples:
  Python: "handler.handle"
  Node.js: "handler.handler"
  Go: "Handle"

Refer to Scaleway documentation for handler conventions per runtime.

- rule: {"string":{"minLen":"1"}}

### spec.privacy

`enum` · required

privacy controls how the function is authenticated.

- rule: privacy must be specified
- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `scaleway_serverless_function_privacy_unspecified` -- Unspecified (invalid).
- `privacy_public` -- Public: the function endpoint is publicly accessible without authentication.
- `privacy_private` -- Private: the function endpoint requires a valid authentication token. Tokens are managed separately via the Scaleway console or API (not bundled in this resource).

### spec.description

`string`

description of the function. Optional, human-readable.

### spec.memoryLimitMb

`uint32`

memory_limit_mb is the memory allocated to each function instance
in megabytes.

Higher memory also increases CPU allocation proportionally.
Common values: 128, 256, 512, 1024, 2048.
Defaults to 256 MB if not specified.

- default: `256`

### spec.minScale

`uint32`

min_scale is the minimum number of always-running function instances.

Set to 0 (default) for scale-to-zero behavior -- the function
spins down when idle and incurs no compute charges.
Set to 1+ for always-warm instances (eliminates cold starts but
incurs continuous billing).

### spec.maxScale

`uint32`

max_scale is the maximum number of concurrent function instances.

The function auto-scales based on incoming workload but never
exceeds this limit. Defaults to 20 if not specified.

- default: `20`

### spec.timeoutSeconds

`uint32`

timeout_seconds is the maximum execution time for a single
function invocation in seconds.

If the function exceeds this duration, Scaleway terminates it.
Defaults to 300 seconds (5 minutes) if not specified.

- default: `300`

### spec.httpOption

`enum`

http_option controls HTTP/HTTPS behavior for the function endpoint.

Defaults to "enabled" (both HTTP and HTTPS allowed) if not specified.

Allowed values (use exactly as shown):

- `scaleway_serverless_function_http_option_unspecified` -- Unspecified -- defaults to "enabled" behavior.
- `enabled` -- Enabled: both HTTP and HTTPS requests are accepted.
- `redirected` -- Redirected: HTTP requests are automatically redirected to HTTPS.

### spec.env

`ScalewayServerlessFunctionEnv`

env groups environment variables and secrets for the function.

Variables are non-secret and visible in logs/dashboards.
Secrets are encrypted at rest and masked in the Scaleway console.

Modeled as repeated name-value messages (Kubernetes-style) rather
than maps to preserve sort order and enable future `valueFrom`
extension.

### spec.env.variables

`[]ScalewayServerlessFunctionEnvVar`

variables are non-secret environment variables.

These are visible in the Scaleway console and may appear in
function logs. Do not put sensitive values here -- use `secrets`.

### spec.env.variables[].name

`string` · required

name is the environment variable name (e.g., "DATABASE_URL").

- rule: {"string":{"minLen":"1"}}

### spec.env.variables[].value

`string` · required

value is the environment variable value.

- rule: {"string":{"minLen":"1"}}

### spec.env.secrets

`[]ScalewayServerlessFunctionEnvVar`

secrets are encrypted environment variables.

Scaleway stores these encrypted at rest and masks them in the
console. Use for database URLs, API keys, tokens, and other
sensitive configuration.

### spec.env.secrets[].name

`string` · required

name is the environment variable name (e.g., "DATABASE_URL").

- rule: {"string":{"minLen":"1"}}

### spec.env.secrets[].value

`string` · required

value is the environment variable value.

- rule: {"string":{"minLen":"1"}}

### spec.privateNetworkId

`string | valueFrom`

private_network_id optionally connects the function to a Scaleway
Private Network for VPC-internal communication.

When connected, the function can reach resources on the Private
Network (databases, Redis clusters, other services) without
traversing the public internet.

Leave unset for functions that only need public internet access.

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.sandbox

`string`

sandbox selects the execution environment for the function.

Available sandboxes depend on the Scaleway platform version.
Common values: "v1" (standard), "v2" (enhanced security).
Leave empty to use the platform default.

### spec.zipFile

`string`

zip_file is the local path to a zip archive containing the
function source code.

When provided, the IaC module uploads this archive and triggers
deployment. Use together with `zip_hash` for change detection.

Leave empty when deploying code separately via the Scaleway CLI
or CI/CD pipeline.

### spec.zipHash

`string`

zip_hash is a hash of the zip archive for change detection.

When the hash changes, the IaC module re-uploads and redeploys
the function. Can be any string (e.g., SHA256 of the zip file).

Only meaningful when `zip_file` is also set.

### spec.cronTriggers

`[]ScalewayServerlessFunctionCronTrigger`

cron_triggers defines optional scheduled triggers for the function.

Each trigger creates a `scaleway_function_cron` resource that
invokes the function on the specified schedule with the given
JSON arguments.

Common patterns:
  Hourly cleanup:    schedule = "0 * * * *"
  Nightly backup:    schedule = "0 2 * * *"
  Every 5 minutes:   schedule = "*/5 * * * *"

### spec.cronTriggers[].name

`string`

name is an optional human-readable identifier for the trigger.

If not provided, Scaleway auto-generates a name. Providing a name
makes triggers easier to identify in the console and logs.

### spec.cronTriggers[].schedule

`string` · required

schedule is a UNIX CRON expression defining when the function
is invoked.

Examples:
  "0 * * * *"    -- every hour at minute 0
  "0 2 * * *"    -- daily at 2:00 AM
  "*/5 * * * *"  -- every 5 minutes
  "0 0 * * 0"    -- weekly on Sunday at midnight

See https://www.scaleway.com/en/docs/serverless/functions/reference-content/cron-schedules/

- rule: {"string":{"minLen":"1"}}

### spec.cronTriggers[].args

`string` · required

args is a JSON string passed to the function's event object on
each scheduled invocation.

Must be valid JSON. Use "{}" for no arguments.

Example: '{"cleanup_type": "stale_sessions", "batch_size": 100}'

- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayServerlessFunction, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.function_id` | `string` | The unique identifier of the deployed serverless function. Format: Scaleway-assigned UUID. Used for API operations, Scaleway CLI commands, and Terraform import. |
| `status.outputs.namespace_id` | `string` | The unique identifier of the function namespace. Useful for managing additional resources that reference the namespace (external cron triggers, tokens, or additional functions not managed by this resource). |
| `status.outputs.domain_name` | `string` | The native Scaleway domain name for invoking the function. This is the HTTPS endpoint automatically assigned by Scaleway. Downstream ScalewayDnsRecord resources can create CNAME records pointing to this domain for custom domain routing. Example: "myfunc-abc123.functions.fnc.fr-par.scw.cloud" |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](./README.md)
