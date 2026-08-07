# Service Account Grant (Predefined Role)

This preset grants one predefined role to a service account — the most common IAM grant in infrastructure code. The member references a GcpServiceAccount resource, so the access relationship is a visible edge in the resource graph rather than a string buried in configuration.

## When to Use

- A workload identity (Cloud Run service, GKE pod, Cloud Function, CI job) needs one well-known predefined role
- You want grants owned and versioned independently of the identity they attach to
- Composing charts where the identity and the grant come from different manifests

## Key Configuration Choices

- **`valueFrom` member** — referencing the GcpServiceAccount's `member` output (the ready-made `serviceAccount:<email>` string) instead of a literal keeps the graph honest and survives identity recreation
- **Predefined role as a literal** — predefined roles (`roles/...`) are stable global names; use a literal. Custom roles should be referenced (see preset 02)
- **Additive semantics** — this grant never touches other members' bindings on the same role

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | Project whose IAM policy receives the grant | GCP Console or `GcpProject` outputs |
| `<service-account-resource-name>` | The Planton resource name of the GcpServiceAccount | Your GcpServiceAccount manifest's `metadata.name` |

## Related Presets

- **02-custom-role-grant** — Grant a custom role defined as a first-class node
- **03-conditional-grant** — Time-bound or resource-scoped access with an IAM condition
