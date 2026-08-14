# AzureDataFactoryDataFlow

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a mapping data
# flow exercising the full endpoint surface -- dataset, linked service,
# schema linked service, and flowlet bindings on the source, rejected-
# row routing on the sink, an intermediate transformation, script
# lines, annotations, and the Studio folder. References are literal
# values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactoryDataFlow
metadata:
  name: test-data-flow
  id: test-data-flow
  org: test-org
  env: test
spec:
  dataFactoryId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DataFactory/factories/test-df
  name: transform-orders
  description: Cleans and conforms raw orders into the curated zone.
  folder: transformations/daily
  annotations:
    - team:data
    - tier:silver
  sources:
    - name: rawOrders
      description: Raw landing-zone orders.
      dataset:
        name:
          value: raw_orders
        parameters:
          path: raw/orders
      schemaLinkedService:
        name:
          value: schema-store
    - name: enrichment
      description: Shared cleanup embedded as a flowlet.
      flowlet:
        name:
          value: scrub-pii
        parameters:
          mode: strict
        datasetParameters: "path: 'raw/orders'"
  sinks:
    - name: curatedOrders
      description: The curated zone table.
      linkedService:
        name:
          value: lakehouse
        parameters:
          container: curated
      rejectedLinkedService:
        name:
          value: quarantine
  transformations:
    - name: joinReference
      description: Joins the reference data.
      linkedService:
        name:
          value: reference-store
  scriptLines:
    - "source(allowSchemaDrift: true, validateSchema: false) ~> rawOrders"
    - "source(allowSchemaDrift: true) ~> enrichment"
    - "rawOrders, enrichment join(orders.id == enrichment.id, joinType:'inner') ~> joinReference"
    - "joinReference sink(allowSchemaDrift: true, validateSchema: false) ~> curatedOrders"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dataFactoryId` | `string \| valueFrom` | yes |  | AzureDataFactory (`status.outputs.data_factory_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.flowlet` | `bool` |  |  |  |
| `spec.script` | `string` |  |  |  |
| `spec.scriptLines` | `[]string` |  |  |  |
| `spec.sources` | `[]AzureDataFactoryDataFlowSource` |  |  |  |
| `spec.sources[].name` | `string` | yes |  |  |
| `spec.sources[].description` | `string` |  |  |  |
| `spec.sources[].dataset` | `AzureDataFactoryDataFlowDatasetReference` |  |  |  |
| `spec.sources[].dataset.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataset (`status.outputs.dataset_name`) |
| `spec.sources[].dataset.parameters` | `map<string, string>` |  |  |  |
| `spec.sources[].flowlet` | `AzureDataFactoryDataFlowFlowletReference` |  |  |  |
| `spec.sources[].flowlet.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataFlow (`status.outputs.data_flow_name`) |
| `spec.sources[].flowlet.parameters` | `map<string, string>` |  |  |  |
| `spec.sources[].flowlet.datasetParameters` | `string` |  |  |  |
| `spec.sources[].linkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.sources[].linkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sources[].linkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.sources[].schemaLinkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.sources[].schemaLinkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sources[].schemaLinkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.sinks` | `[]AzureDataFactoryDataFlowSink` |  |  |  |
| `spec.sinks[].name` | `string` | yes |  |  |
| `spec.sinks[].description` | `string` |  |  |  |
| `spec.sinks[].dataset` | `AzureDataFactoryDataFlowDatasetReference` |  |  |  |
| `spec.sinks[].dataset.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataset (`status.outputs.dataset_name`) |
| `spec.sinks[].dataset.parameters` | `map<string, string>` |  |  |  |
| `spec.sinks[].flowlet` | `AzureDataFactoryDataFlowFlowletReference` |  |  |  |
| `spec.sinks[].flowlet.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataFlow (`status.outputs.data_flow_name`) |
| `spec.sinks[].flowlet.parameters` | `map<string, string>` |  |  |  |
| `spec.sinks[].flowlet.datasetParameters` | `string` |  |  |  |
| `spec.sinks[].linkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.sinks[].linkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sinks[].linkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.sinks[].schemaLinkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.sinks[].schemaLinkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sinks[].schemaLinkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.sinks[].rejectedLinkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.sinks[].rejectedLinkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sinks[].rejectedLinkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.transformations` | `[]AzureDataFactoryDataFlowTransformation` |  |  |  |
| `spec.transformations[].name` | `string` | yes |  |  |
| `spec.transformations[].description` | `string` |  |  |  |
| `spec.transformations[].dataset` | `AzureDataFactoryDataFlowDatasetReference` |  |  |  |
| `spec.transformations[].dataset.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataset (`status.outputs.dataset_name`) |
| `spec.transformations[].dataset.parameters` | `map<string, string>` |  |  |  |
| `spec.transformations[].flowlet` | `AzureDataFactoryDataFlowFlowletReference` |  |  |  |
| `spec.transformations[].flowlet.name` | `string \| valueFrom` | yes |  | AzureDataFactoryDataFlow (`status.outputs.data_flow_name`) |
| `spec.transformations[].flowlet.parameters` | `map<string, string>` |  |  |  |
| `spec.transformations[].flowlet.datasetParameters` | `string` |  |  |  |
| `spec.transformations[].linkedService` | `AzureDataFactoryDataFlowLinkedServiceReference` |  |  |  |
| `spec.transformations[].linkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.transformations[].linkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.annotations` | `[]string` |  |  |  |
| `spec.folder` | `string` |  |  |  |

## Field Details

### spec.dataFactoryId

`string | valueFrom` · required

- references: AzureDataFactory (`status.outputs.data_factory_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactory, name: <that resource's name>, fieldPath: status.outputs.data_factory_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.flowlet

`bool`

### spec.script

`string`

### spec.scriptLines

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.sources

`[]AzureDataFactoryDataFlowSource`

### spec.sources[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sources[].description

`string`

### spec.sources[].dataset

`AzureDataFactoryDataFlowDatasetReference`

### spec.sources[].dataset.name

`string | valueFrom` · required

- references: AzureDataFactoryDataset (`status.outputs.dataset_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataset, name: <that resource's name>, fieldPath: status.outputs.dataset_name}} -- a bare string does not parse

### spec.sources[].dataset.parameters

`map<string, string>`

### spec.sources[].flowlet

`AzureDataFactoryDataFlowFlowletReference`

### spec.sources[].flowlet.name

`string | valueFrom` · required

- references: AzureDataFactoryDataFlow (`status.outputs.data_flow_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataFlow, name: <that resource's name>, fieldPath: status.outputs.data_flow_name}} -- a bare string does not parse

### spec.sources[].flowlet.parameters

`map<string, string>`

### spec.sources[].flowlet.datasetParameters

`string`

### spec.sources[].linkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.sources[].linkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sources[].linkedService.parameters

`map<string, string>`

### spec.sources[].schemaLinkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.sources[].schemaLinkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sources[].schemaLinkedService.parameters

`map<string, string>`

### spec.sinks

`[]AzureDataFactoryDataFlowSink`

### spec.sinks[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sinks[].description

`string`

### spec.sinks[].dataset

`AzureDataFactoryDataFlowDatasetReference`

### spec.sinks[].dataset.name

`string | valueFrom` · required

- references: AzureDataFactoryDataset (`status.outputs.dataset_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataset, name: <that resource's name>, fieldPath: status.outputs.dataset_name}} -- a bare string does not parse

### spec.sinks[].dataset.parameters

`map<string, string>`

### spec.sinks[].flowlet

`AzureDataFactoryDataFlowFlowletReference`

### spec.sinks[].flowlet.name

`string | valueFrom` · required

- references: AzureDataFactoryDataFlow (`status.outputs.data_flow_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataFlow, name: <that resource's name>, fieldPath: status.outputs.data_flow_name}} -- a bare string does not parse

### spec.sinks[].flowlet.parameters

`map<string, string>`

### spec.sinks[].flowlet.datasetParameters

`string`

### spec.sinks[].linkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.sinks[].linkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sinks[].linkedService.parameters

`map<string, string>`

### spec.sinks[].schemaLinkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.sinks[].schemaLinkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sinks[].schemaLinkedService.parameters

`map<string, string>`

### spec.sinks[].rejectedLinkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.sinks[].rejectedLinkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sinks[].rejectedLinkedService.parameters

`map<string, string>`

### spec.transformations

`[]AzureDataFactoryDataFlowTransformation`

### spec.transformations[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.transformations[].description

`string`

### spec.transformations[].dataset

`AzureDataFactoryDataFlowDatasetReference`

### spec.transformations[].dataset.name

`string | valueFrom` · required

- references: AzureDataFactoryDataset (`status.outputs.dataset_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataset, name: <that resource's name>, fieldPath: status.outputs.dataset_name}} -- a bare string does not parse

### spec.transformations[].dataset.parameters

`map<string, string>`

### spec.transformations[].flowlet

`AzureDataFactoryDataFlowFlowletReference`

### spec.transformations[].flowlet.name

`string | valueFrom` · required

- references: AzureDataFactoryDataFlow (`status.outputs.data_flow_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryDataFlow, name: <that resource's name>, fieldPath: status.outputs.data_flow_name}} -- a bare string does not parse

### spec.transformations[].flowlet.parameters

`map<string, string>`

### spec.transformations[].flowlet.datasetParameters

`string`

### spec.transformations[].linkedService

`AzureDataFactoryDataFlowLinkedServiceReference`

### spec.transformations[].linkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.transformations[].linkedService.parameters

`map<string, string>`

### spec.description

`string`

### spec.annotations

`[]string`

### spec.folder

`string`

## Validation Rules

- `azure_data_factory_data_flow_script_required`: Provide the data flow script -- script, script_lines, or both
- `azure_data_factory_data_flow_mapping_source_sink`: A mapping data flow requires at least one source and one sink -- only flowlets (flowlet: true) may omit them

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactoryDataFlow, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.data_flow_id` | `string` |  |
| `status.outputs.data_flow_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataFactoryId` | AzureDataFactory | `status.outputs.data_factory_id` |
| `spec.sources[].dataset.name` | AzureDataFactoryDataset | `status.outputs.dataset_name` |
| `spec.sources[].flowlet.name` | AzureDataFactoryDataFlow | `status.outputs.data_flow_name` |
| `spec.sources[].linkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sources[].schemaLinkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sinks[].dataset.name` | AzureDataFactoryDataset | `status.outputs.dataset_name` |
| `spec.sinks[].flowlet.name` | AzureDataFactoryDataFlow | `status.outputs.data_flow_name` |
| `spec.sinks[].linkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sinks[].schemaLinkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sinks[].rejectedLinkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.transformations[].dataset.name` | AzureDataFactoryDataset | `status.outputs.dataset_name` |
| `spec.transformations[].flowlet.name` | AzureDataFactoryDataFlow | `status.outputs.data_flow_name` |
| `spec.transformations[].linkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryDataFlow | `spec.sources[].flowlet.name` | `status.outputs.data_flow_name` |
| AzureDataFactoryDataFlow | `spec.sinks[].flowlet.name` | `status.outputs.data_flow_name` |
| AzureDataFactoryDataFlow | `spec.transformations[].flowlet.name` | `status.outputs.data_flow_name` |

## See Also

- [Overview](../README.md)
