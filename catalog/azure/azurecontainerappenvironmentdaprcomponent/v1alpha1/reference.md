# AzureContainerAppEnvironmentDaprComponent

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureContainerAppEnvironmentDaprComponentSpec** defines the configuration
for registering a Dapr component on a Container App Environment
(Microsoft.App/managedEnvironments/daprComponents).

Dapr components are the pluggable backends behind Dapr's building blocks:
a state store (state.azure.blobstorage, state.redis), a pub/sub broker
(pubsub.azure.servicebus, pubsub.kafka), a secret store, an input/output
binding. Components are registered once on the environment and consumed
by any Dapr-enabled app whose `dapr.app_id` appears in `scopes` (an empty
scopes list exposes the component to every Dapr-enabled app in the
environment).

**Configuration model**: `metadata` entries carry the component's
connection settings (host, connection string name, consistency level --
the keys depend on the component type; see the Dapr component reference).
Values that are secrets go into `secrets` and are referenced from
metadata by `secret_name`, never inlined.

**ForceNew fields** (changing these destroys and recreates the
component): `component_name`, `component_type`,
`container_app_environment_id`.

**Referenced by**: None directly -- apps consume components through
Dapr's runtime by component name, scoped via their dapr.app_id.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppEnvironmentDaprComponent
metadata:
  name: test-dapr-component
spec:
  container_app_environment_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/managedEnvironments/test-env
  component_name: statestore
  component_type: state.azure.blobstorage
  version: v1
  init_timeout: 10s
  ignore_errors: false
  secrets:
    - name: account-key
      value: dGVzdC1hY2NvdW50LWtleQ==
  metadata:
    - name: accountName
      value:
        value: teststorageaccount
    - name: containerName
      value:
        value: dapr-state
    - name: accountKey
      secret_name: account-key
  scopes:
    - orders
    - billing
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.containerAppEnvironmentId` | `string \| valueFrom` | yes |  | AzureContainerAppEnvironment (`status.outputs.environment_id`) |
| `spec.componentName` | `string` | yes |  |  |
| `spec.componentType` | `string` | yes |  |  |
| `spec.version` | `string` | yes |  |  |
| `spec.initTimeout` | `string` |  | `5s` |  |
| `spec.ignoreErrors` | `bool` |  | `false` |  |
| `spec.secrets` | `[]AzureContainerAppEnvironmentDaprComponentSecret` |  |  |  |
| `spec.secrets[].name` | `string` | yes |  |  |
| `spec.secrets[].value` | `string` (sensitive) | yes |  |  |
| `spec.metadata` | `[]AzureContainerAppEnvironmentDaprComponentMetadata` |  |  |  |
| `spec.metadata[].name` | `string` | yes |  |  |
| `spec.metadata[].value` | `string \| valueFrom` |  |  |  |
| `spec.metadata[].secretName` | `string` |  |  |  |
| `spec.scopes` | `[]string` |  |  |  |

## Field Details

### spec.containerAppEnvironmentId

`string | valueFrom` · required

The Container App Environment to register the Dapr component on.

**ForceNew**: Changing this destroys and recreates the component.

- references: AzureContainerAppEnvironment (`status.outputs.environment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerAppEnvironment, name: <that resource's name>, fieldPath: status.outputs.environment_id}} -- a bare string does not parse

### spec.componentName

`string` · required

The name of the Dapr component -- what application code passes to the
Dapr API when using the building block (e.g. the store name in a
state.get call).
Lowercase alphanumeric characters and hyphens; must start with a
letter and end with an alphanumeric character; no consecutive
hyphens; at most 60 characters.

**ForceNew**: Changing this destroys and recreates the component.

- rule: component name must be lowercase alphanumeric characters or hyphens, start with a letter, end with an alphanumeric character, and contain no consecutive hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"60"}}

### spec.componentType

`string` · required

The Dapr component type identifying the building block and backend,
in Dapr's own dotted notation.
Examples: "state.azure.blobstorage", "pubsub.azure.servicebus",
"secretstores.azure.keyvault", "bindings.cron".

**ForceNew**: Changing this destroys and recreates the component.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.version

`string` · required

The component version. Dapr components are versioned independently of
the Dapr runtime; "v1" for virtually all stable components.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initTimeout

`string` · optional (explicit presence)

How long the Dapr sidecar waits for the component to initialise.
Whole intervals of seconds, minutes, or hours: "5s", "10m", "1h".

Default: "5s"

- default: `5s`
- rule: init_timeout must be a whole interval of seconds, minutes, or hours, e.g. 5s, 10m, 1h

### spec.ignoreErrors

`bool` · optional (explicit presence)

Whether the Dapr sidecar continues initialising when this component
fails to load. Leave false so a broken state store fails loudly at
startup instead of surfacing as runtime errors on first use.

Default: false

- default: `false`

### spec.secrets

`[]AzureContainerAppEnvironmentDaprComponentSecret`

Secrets available to this component's metadata. Referenced from
metadata entries by `secret_name` -- connection strings and access
keys never appear as plain metadata values.

### spec.secrets[].name

`string` · required

Secret name. Lowercase alphanumeric or hyphens.
Referenced from metadata entries by secret_name.

- rule: secret name must be lowercase alphanumeric or hyphens, start and end with alphanumeric
- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.secrets[].value

`string` · required · sensitive

The secret value (a connection string, access key, or token the
component's backend requires).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.metadata

`[]AzureContainerAppEnvironmentDaprComponentMetadata`

The component's configuration entries. Keys depend on the component
type -- consult the Dapr reference for the chosen backend.
Example for state.azure.blobstorage: accountName, containerName, and
a secret-backed accountKey.

- rule: a metadata entry takes either a value or a secret_name, not both -- move secrets into the component's secrets list

### spec.metadata[].name

`string` · required

The metadata key (component-type specific, e.g. "accountName",
"redisHost", "consumerID").

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.metadata[].value

`string | valueFrom`

The entry's value: a reference or a literal. Reference another
resource's output when the value is minted at deploy time -- the
keyless-auth entries are the canonical case (an "azureClientId"
entry tracking an AzureUserAssignedIdentity's client_id output, so
the component authenticates with a managed identity instead of a
connection-string secret). Pass literals for everything knowable up
front (hosts, namespace names, consumer IDs, feature flags). No kind
dominates Dapr metadata values, so references declare their kind
explicitly. Mutually exclusive with secret_name.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.metadata[].secretName

`string`

Reference to a secret in the component's `secrets` list. Use for
connection strings and keys. Mutually exclusive with value.

### spec.scopes

`[]string`

The `dapr.app_id` values of the Container Apps allowed to use this
component. An empty list exposes the component to every Dapr-enabled
app in the environment -- scope production components deliberately.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerAppEnvironmentDaprComponent, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.dapr_component_id` | `string` | The Azure Resource Manager ID of the Dapr component. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.App/managedEnvironments/{env}/daprComponents/{name} |
| `status.outputs.component_name` | `string` | The name of the Dapr component -- what application code passes to the Dapr API when using the building block. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.containerAppEnvironmentId` | AzureContainerAppEnvironment | `status.outputs.environment_id` |

## See Also

- [Overview](../README.md)
