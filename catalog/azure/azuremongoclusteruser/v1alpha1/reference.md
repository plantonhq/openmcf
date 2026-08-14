# AzureMongoClusterUser

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Example for docs and offline validation: a managed identity's
# service principal granted root on the admin database (cluster-wide
# access). References are literal values so the manifest validates
# standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMongoClusterUser
metadata:
  name: test-mongo-cluster-user
  id: test-mongo-cluster-user
  org: test-org
  env: test
spec:
  mongoClusterId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DocumentDB/mongoClusters/acme-orders-db
  objectId:
    value: 11111111-2222-3333-4444-555555555555
  principalType: servicePrincipal
  roles:
    - database: admin
      role: root
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.mongoClusterId` | `string \| valueFrom` | yes |  | AzureMongoCluster (`status.outputs.mongo_cluster_id`) |
| `spec.objectId` | `string \| valueFrom` | yes |  |  |
| `spec.principalType` | `string` | yes |  |  |
| `spec.roles` | `[]AzureMongoClusterUserRole` | yes |  |  |
| `spec.roles[].database` | `string` | yes |  |  |
| `spec.roles[].role` | `string` | yes |  |  |

## Field Details

### spec.mongoClusterId

`string | valueFrom` · required

- references: AzureMongoCluster (`status.outputs.mongo_cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMongoCluster, name: <that resource's name>, fieldPath: status.outputs.mongo_cluster_id}} -- a bare string does not parse

### spec.objectId

`string | valueFrom` · required

- rule: object_id must be the principal's object id -- a UUID like 00000000-0000-0000-0000-000000000000
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.principalType

`string` · required

- rule: {"required":true,"string":{"in":["user","servicePrincipal"]}}

### spec.roles

`[]AzureMongoClusterUserRole` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.roles[].database

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.roles[].role

`string` · required

- rule: {"required":true,"string":{"in":["root"]}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMongoClusterUser, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.mongo_cluster_user_id` | `string` |  |
| `status.outputs.mongo_cluster_user_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.mongoClusterId` | AzureMongoCluster | `status.outputs.mongo_cluster_id` |

## See Also

- [Overview](../README.md)
