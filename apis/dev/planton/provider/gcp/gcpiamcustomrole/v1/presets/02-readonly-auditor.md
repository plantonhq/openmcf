# Read-Only Auditor Role

This preset defines a read-only role for security audits, compliance dashboards, and inventory tooling — visibility into project configuration, IAM policy, service accounts, and storage layout without any `get`-object or mutate permissions.

## When to Use

- A security team or audit tool needs to enumerate what exists and who can access it
- Compliance dashboards that read IAM policies and resource configuration
- You want auditors to see policy without being able to read data (note: no `storage.objects.*` here)

## Key Configuration Choices

- **Policy visibility, not data visibility** — the permission set includes `getIamPolicy` surfaces but deliberately excludes object-read permissions, so audit access never becomes data access
- **List + get pairs** — enumeration (`list`) and inspection (`get`) permissions travel together for each audited surface
- **Extend per audited service** — add the `*.get`/`*.list`/`*.getIamPolicy` permissions for each additional service in scope (e.g. `cloudsql.instances.get`, `container.clusters.get`)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project that owns the role | GCP Console or `GcpProject` outputs |

## Related Presets

- **01-workload-least-privilege** — The standard workload-scoped permission bundle
- **03-ci-cd-deployer** — A deployment role for CI/CD pipelines
