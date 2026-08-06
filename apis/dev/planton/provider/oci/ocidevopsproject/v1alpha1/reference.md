# OciDevopsProject

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciDevopsProjectSpec defines the specification for an OCI DevOps
Project -- the organizational container for CI/CD pipelines, code
repositories, deployment environments, artifacts, and triggers in
OCI's managed DevOps service.

A project provides a shared namespace and a notification topic that
receives pipeline events (build started, deployment succeeded, etc.).
All downstream DevOps resources (build pipelines, deploy pipelines,
repositories, connections) reference the project by its OCID.

Key behaviors:
  - name (from metadata.name) is ForceNew (immutable after creation)
  - compartment_id is updatable (supports compartment moves)
  - notification_topic_id is updatable

Provider nesting flattened:
  - notification_config.topic_id -> notification_topic_id
    (single-field block, no benefit to nesting)

Excluded from v1:
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.notificationTopicId` | `string \| valueFrom` | yes |  |  |
| `spec.description` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the project will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.notificationTopicId

`string | valueFrom` · required

OCID of the ONS (Oracle Notification Service) topic that receives
DevOps pipeline events such as build completions, deployment
successes, and failures. Required by the OCI API.

Flattened from the provider's notification_config.topic_id block.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.description

`string`

Human-readable description of the project's purpose.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciDevopsProject, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.project_id` | `string` | OCID of the DevOps project. |
| `status.outputs.namespace` | `string` | Namespace associated with the project, used in container registry paths and artifact references. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
