# AliCloudContainerRegistry

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudContainerRegistrySpec defines the configuration for an Alibaba Cloud
Container Registry (ACR) Enterprise Edition instance bundled with namespaces.

ACR Enterprise Edition provides a managed container image registry with
enterprise-grade security, multi-region replication, and image scanning.
Three tiers are available -- Basic, Standard, and Advanced -- selected via
the instance_type field, similar to how Azure Container Registry uses
Basic/Standard/Premium SKUs.

Namespaces are bundled into this component (per DD07 composite bundling)
because a registry instance without namespaces is an empty shell that cannot
store images. Namespaces are the organizational units within the registry.

Note: ACR Enterprise Edition instances do not support tags in the provider.
This is a provider limitation for BSS-provisioned resources.

Provider resources:
  Terraform: alicloud_cr_ee_instance + alicloud_cr_ee_namespace
  Pulumi:    cr.RegistryEnterpriseInstance + cs.RegistryEnterpriseNamespace

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudContainerRegistry
metadata:
  name: test-registry
spec:
  region: cn-hangzhou
  instanceName: test-acr-registry
  instanceType: Basic
  paymentType: Subscription
  period: 1
  namespaces:
    - name: platform
      autoCreate: true
      defaultVisibility: PRIVATE
    - name: frontend
      autoCreate: false
      defaultVisibility: PRIVATE
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.instanceName` | `string` | yes |  |  |
| `spec.instanceType` | `string` | yes |  |  |
| `spec.paymentType` | `string` |  | `Subscription` |  |
| `spec.period` | `int32` |  |  |  |
| `spec.password` | `string` (sensitive) |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.namespaces` | `[]AliCloudContainerRegistryNamespace` |  |  |  |
| `spec.namespaces[].name` | `string` | yes |  |  |
| `spec.namespaces[].autoCreate` | `bool` |  | `false` |  |
| `spec.namespaces[].defaultVisibility` | `string` |  | `PRIVATE` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the registry instance will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceName

`string` · required

Registry instance name. This is the human-readable identifier for the
ACR Enterprise Edition instance.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceType

`string` · required

Tier of the ACR Enterprise Edition instance.
Basic: suitable for individual developers or small teams.
Standard: suitable for small and medium enterprises.
Advanced: suitable for large enterprises with geo-replication needs.

The instance type is immutable after creation (ForceNew in the provider).

- rule: instance_type must be one of: Basic, Standard, Advanced
- rule: {"required":true}

### spec.paymentType

`string` · optional (explicit presence)

Payment model for the registry instance.
Subscription: pre-paid monthly/yearly billing (lower cost for long-term use).
PayAsYouGo: post-paid usage-based billing (flexible for dev/test).

Immutable after creation (ForceNew in the provider).
Default: "Subscription"

- default: `Subscription`
- rule: payment_type must be one of: Subscription, PayAsYouGo

### spec.period

`int32`

Subscription period in months. Only applies when payment_type is
"Subscription". Common values: 1, 3, 6, 12.

### spec.password

`string` · sensitive

Login password for the registry instance. Used for authenticating
docker login / image push and pull operations.
8-32 characters.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the instance is placed in the account's default resource group.

### spec.namespaces

`[]AliCloudContainerRegistryNamespace`

Namespaces to create within the registry instance.
Namespaces are the top-level organizational units for container images.
Images are addressed as: {registry-endpoint}/{namespace}/{repo}:{tag}

### spec.namespaces[].name

`string` · required

Namespace name. 2-120 characters; can contain lowercase letters, digits,
underscores, hyphens, and periods. Cannot start or end with a delimiter.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"120"}}

### spec.namespaces[].autoCreate

`bool` · optional (explicit presence)

When true, image repositories are automatically created within this
namespace when an image is pushed to a repository that does not yet exist.
This is convenient for CI/CD pipelines that push to dynamic repository names.
Default: false

- default: `false`

### spec.namespaces[].defaultVisibility

`string` · optional (explicit presence)

Default visibility for repositories auto-created within this namespace.
PUBLIC: anyone can pull images (no authentication required).
PRIVATE: only authenticated users with access can pull images.
Default: "PRIVATE"

- default: `PRIVATE`
- rule: default_visibility must be one of: PUBLIC, PRIVATE

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudContainerRegistry, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The ACR Enterprise Edition instance ID. Referenced by downstream resources that need to interact with this registry. |
| `status.outputs.instance_name` | `string` | The registry instance name as created. |
| `status.outputs.public_endpoint` | `string` | The internet-facing registry endpoint domain. Used for docker login, push, and pull from outside VPC. Example: "myregistry-registry.cn-hangzhou.cr.aliyuncs.com" |
| `status.outputs.vpc_endpoint` | `string` | The VPC-internal registry endpoint domain. Used for pulling images from within the same VPC (e.g., from ACK nodes) for faster transfer and no internet egress cost. |
| `status.outputs.namespace_ids` | `map<string, string>` | Map of created namespace names to their IDs. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
