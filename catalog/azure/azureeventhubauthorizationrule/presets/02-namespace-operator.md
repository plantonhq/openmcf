# Namespace-Scoped Operator Credential

This preset mints a full-rights (listen+send+manage) SAS credential over
a whole namespace -- for operational tooling that creates and inspects
hubs and consumer groups, distinct from the built-in root rule.

## When to Use

- Streaming-platform tooling that provisions or audits entities across
  the namespace
- A dedicated, revocable operator credential -- unlike the built-in
  `RootManageSharedAccessKey`, deleting this rule cuts off exactly one
  consumer of manage rights

## Key Configuration Choices

- **`namespaceId` scope** -- rights span every hub in the namespace;
  prefer hub-scoped rules for applications
- **The manage trio** -- Azure's contract: manage requires listen AND
  send; the spec enforces it up front
- **Never reuse the root rule** -- `RootManageSharedAccessKey` is
  reserved (its keys already surface as namespace outputs); minting a
  named rule keeps operator access individually revocable

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-telemetry-hubs` | The AzureEventHubNamespace's Planton resource name | Your streaming composition |
| `streaming-operator` | The rule name | Your credential taxonomy |
