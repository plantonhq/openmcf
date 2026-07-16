# Custom Role Grant (Fully Composed)

This preset is the full least-privilege composition: a GcpIamCustomRole defines exactly the permissions a workload needs, a GcpServiceAccount is the workload's identity, and this grant is the edge binding them on a project. All three are first-class nodes — permissions, identity, and access each change independently.

## When to Use

- The workload's permission set doesn't match any predefined role (the normal case after a security review)
- You want permission changes (on the role) decoupled from membership changes (on grants)
- Building reusable charts where role definitions and grants ship together

## Key Configuration Choices

- **Both halves referenced** — the role's `name` output and the identity's `member` output are exactly the strings IAM expects; no assembly, no drift
- **Deploy ordering is automatic** — the references make the role and the identity explicit dependencies of the grant
- **One grant per (role, member) pair** — need the same role for three identities? Three grant nodes; each independently removable

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | Project whose IAM policy receives the grant | GCP Console or `GcpProject` outputs |
| `<custom-role-resource-name>` | The Planton resource name of the GcpIamCustomRole | Your GcpIamCustomRole manifest's `metadata.name` |
| `<service-account-resource-name>` | The Planton resource name of the GcpServiceAccount | Your GcpServiceAccount manifest's `metadata.name` |

## Related Presets

- **01-service-account-grant** — The simpler predefined-role variant
- **03-conditional-grant** — Time-bound or resource-scoped access with an IAM condition
