# DigitalOceanFunction

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanFunctionSpec defines the configuration for deploying a serverless function on DigitalOcean.
Functions are deployed via DigitalOcean App Platform for production-ready VPC integration, monitoring, and IaC support.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.functionName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.runtime` | `enum` | yes |  |  |
| `spec.githubSource` | `DigitalOceanFunctionGithubSource` |  |  |  |
| `spec.githubSource.repo` | `string` | yes |  |  |
| `spec.githubSource.branch` | `string` | yes |  |  |
| `spec.githubSource.deployOnPush` | `bool` |  | `true` |  |
| `spec.sourceDirectory` | `string` | yes |  |  |
| `spec.memoryMb` | `uint32` |  | `256` |  |
| `spec.timeoutMs` | `uint32` |  | `3000` |  |
| `spec.environmentVariables` | `map<string, string>` |  |  |  |
| `spec.secretEnvironmentVariables` | `map<string, string>` |  |  |  |
| `spec.entrypoint` | `string` |  |  |  |
| `spec.cronSchedule` | `string` |  |  |  |
| `spec.isWeb` | `bool` |  | `true` |  |

## Field Details

### spec.functionName

`string` · required

function_name is the name of the function. Must be unique within the project.

- rule: {"string":{"minLen":"1","maxLen":"64"}}

### spec.region

`enum` · required

region specifies the DigitalOcean region to deploy the function.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.runtime

`enum` · required

runtime specifies the runtime environment for the function (e.g., nodejs, python, go).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_function_runtime_unspecified`
- `nodejs_18` -- Node.js 18 (LTS)
- `nodejs_20` -- Node.js 20 (Current)
- `python_39` -- Python 3.9
- `python_310` -- Python 3.10
- `python_311` -- Python 3.11 (Latest)
- `go_120` -- Go 1.20
- `go_121` -- Go 1.21
- `php_82` -- PHP 8.2

### spec.githubSource

`DigitalOceanFunctionGithubSource`

GitHub repository configuration for function source code

### spec.githubSource.repo

`string` · required

repo is the GitHub repository in the format "owner/repo"

- rule: {"string":{"minLen":"1","pattern":"^[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+$"}}

### spec.githubSource.branch

`string` · required

branch is the Git branch to deploy from (e.g., "main", "production")

- rule: {"string":{"minLen":"1","maxLen":"255"}}

### spec.githubSource.deployOnPush

`bool`

deploy_on_push enables automatic redeployment when changes are pushed to the branch

- default: `true`

### spec.sourceDirectory

`string` · required

source_directory is the path within the repository containing the function code and project.yml
Example: "/functions/api-handler"

- rule: {"string":{"minLen":"1"}}

### spec.memoryMb

`uint32`

memory_mb is the memory allocated to the function (in megabytes). Defaults to 256 if not specified.
Valid values: 128, 256, 512, 1024, 2048

- default: `256`
- rule: {"uint32":{"in":[128,256,512,1024,2048]}}

### spec.timeoutMs

`uint32`

timeout_ms is the maximum execution time for the function in milliseconds.
Defaults to 3000ms (3 seconds) if not specified. Max: 300000ms (5 minutes)

- default: `3000`
- rule: {"uint32":{"lte":300000}}

### spec.environmentVariables

`map<string, string>`

environment_variables are non-secret environment variables for the function

### spec.secretEnvironmentVariables

`map<string, string>`

secret_environment_variables are encrypted environment variables (e.g., database URLs, API keys)
These are stored securely in App Platform's secret store

### spec.entrypoint

`string`

entrypoint is an optional function or script entrypoint name within the code
Example: "main" for Go, "handler" for Node.js

### spec.cronSchedule

`string`

cron_schedule is an optional cron expression for scheduled function execution
Example: "0 * * * *" for hourly execution
If set, the function will not be exposed as an HTTP endpoint

### spec.isWeb

`bool`

is_web indicates if the function should be exposed as an HTTP endpoint
Defaults to true. Set to false for background/scheduled functions.

- default: `true`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanFunction, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.function_id` | `string` | function_id is the unique identifier of the deployed function. |
| `status.outputs.https_endpoint` | `string` | https_endpoint is the public HTTPS URL endpoint for invoking the function. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
