# Wiring Resources Together

Charts rarely stand alone — a VPC feeds subnets, subnets feed a cluster, a
cluster feeds add-ons. Planton resolves these links at deploy time using a
dependency graph built from your references.

## valueFrom — the primary wiring mechanism

When resource B needs an ID or ARN that resource A creates after deployment,
never paste a literal. Reference A's stack output:

```yaml
spec:
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: "{{ values.env }}-vpc"
      fieldPath: status.outputs.vpc_id
```

Rules:

- **kind** — PascalCase cloud resource kind of the producer.
- **name** — must match the producer's `metadata.name` exactly (including
  template expressions — both sides must render to the same string).
- **fieldPath** — `status.outputs.<outputName>`. Both snake_case
  (`status.outputs.vpc_id`) and camelCase (`status.outputs.vpcId`) are
  accepted: the validator tries the exact proto field name first, then
  converts camelCase to snake_case. The schema report's `outputs[]` lists
  the camelCase names; the fleet mostly writes snake_case. Either resolves
  to the same output — pick one style per chart and stay consistent.

Fields typed `string | valueFrom` in the schema report accept either a plain
string or the `valueFrom` block above.

### Nested valueFrom

Some fields nest the reference one level deeper (e.g. route targets):

```yaml
routes:
  - destinationCidrBlock: 0.0.0.0/0
    targetType: internet_gateway
    targetId:
      valueFrom:
        kind: AwsInternetGateway
        name: "{{ values.env }}-igw"
        fieldPath: status.outputs.internet_gateway_id
```

## Finding output field paths

1. Run `planton explain <ProducerKind> -o json`.
2. Read the `outputs` array — each entry's `name` is the leaf under
   `status.outputs.`.
3. Cross-check against the chart fleet: grep for `valueFrom` with that kind
   in `github.com/plantonhq/planton/tree/main/charts`.

Build errors for bad references are explicit: `Invalid valueFrom references:
Field 'no_such_output' not found in …StackOutputs for kind: …` — fix the
`fieldPath` leaf to match the schema report, not the provider's API docs.

## DAG semantics

The platform builds a directed acyclic graph from two edge sources — every
`valueFrom` link AND every `metadata.relationships` entry — and deploys
resources in topological order.

Implications for chart authors:

- Every producer must appear in the chart (or already exist in the target
  environment — charts assume greenfield unless documented otherwise).
- Circular references fail at build time.
- Resources with no incoming references deploy in parallel in the earliest
  layer.

## metadata.relationships — real ordering edges

Relationships create REAL dependency edges: a resource with a `runs_on` or
`depends_on` relationship does not start until its target has deployed
successfully.

```yaml
metadata:
  name: "{{ values.env }}-istio"
  relationships:
    - kind: AwsEksCluster            # PascalCase kind of the target
      name: "{{ values.env }}-cluster"  # target's exact metadata.name
      type: runs_on                  # runs_on | depends_on | uses | managed_by
```

The `name` must match the target's `metadata.name` exactly (same template
expression on both sides). Choosing between the two mechanisms:

- **`valueFrom`** when a spec field needs a value the producer creates — it
  carries the data AND the edge.
- **`relationships`** when the dependency is real but no spec field carries a
  value — the canonical case is Kubernetes workloads that must wait for their
  cluster (see `kubernetes-on-cluster.md`), or an operator that must install
  before the instances it serves.

A relationship never substitutes for `valueFrom` when a spec field needs the
actual value.

## metadata.group — layout hint

`metadata.group` (e.g. `network`, `compute`, `kubernetes`) groups resources
in the UI and in build reports. It does not affect deploy order. The fleet
uses consistent group names per concern — follow the same convention in new
charts.

## Common wiring mistakes

| Mistake | Fix |
|---------|-----|
| Hardcoded AWS account ID or VPC ID | valueFrom to the creating resource |
| Mismatched name between producer and consumer | Same template expression on both sides |
| Wrong output leaf (`vpc_id` vs `vpcId`) | Schema report + fleet grep |
| Referencing a conditionally omitted resource | Wrap the consumer in the same `{% if %}` or ensure the producer always renders |
