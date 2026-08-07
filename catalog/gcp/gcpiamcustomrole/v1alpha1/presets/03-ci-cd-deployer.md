# CI/CD Deployer Role

This preset defines a deployment role for CI/CD pipelines that roll out new Cloud Run revisions — update access to the services being deployed and pull access to the artifact repository, without `roles/run.admin`'s ability to delete services or rewrite their IAM.

## When to Use

- A pipeline identity (often a Workload-Identity-federated CI job) deploys application revisions
- You want deploy access separated from service administration (create/delete/IAM stays with humans or a different role)
- The pipeline pulls container images from Artifact Registry in the same project

## Key Configuration Choices

- **`run.services.update` but not `run.services.create`/`delete`** — the pipeline updates existing services; creating and deleting services is an infrastructure change that belongs to IaC, not CI
- **`artifactregistry.repositories.downloadArtifacts`** — image pull for the deploy; push permissions belong to the build identity, not the deploy identity
- **Adapt per deploy target** — for GKE swap in `container.*` permissions, for Cloud Functions the `cloudfunctions.*` set

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project that owns the role | GCP Console or `GcpProject` outputs |

## Related Presets

- **01-workload-least-privilege** — The standard workload-scoped permission bundle
- **02-readonly-auditor** — A read-only role for auditors and dashboards
