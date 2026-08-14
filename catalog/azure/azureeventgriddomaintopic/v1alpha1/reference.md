# AzureEventgridDomainTopic

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: one tenant's
# stream inside a domain. References are literal ARM ids so the
# manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventgridDomainTopic
metadata:
  name: test-eventgrid-domain-topic
  id: test-eventgrid-domain-topic
  org: test-org
  env: test
spec:
  domainId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventGrid/domains/test-org-tenant-events
  # The name publishers stamp into every event's topic field.
  name: customer-fabrikam
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainId` | `string \| valueFrom` | yes |  | AzureEventgridDomain (`status.outputs.domain_id`) |
| `spec.name` | `string` | yes |  |  |

## Field Details

### spec.domainId

`string | valueFrom` · required

- references: AzureEventgridDomain (`status.outputs.domain_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventgridDomain, name: <that resource's name>, fieldPath: status.outputs.domain_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Domain topic names must be 3-128 characters of letters, numbers, and hyphens
- rule: {"required":true}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventgridDomainTopic, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_topic_id` | `string` |  |
| `status.outputs.domain_topic_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.domainId` | AzureEventgridDomain | `status.outputs.domain_id` |

## See Also

- [Overview](../README.md)
