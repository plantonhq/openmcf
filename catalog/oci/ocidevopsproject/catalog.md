# DevOps Project on OCI

Deploys an Oracle Cloud Infrastructure DevOps project -- the organizational container for CI/CD pipelines, code repositories, deployment environments, artifacts, and triggers. The project provides a shared namespace and an ONS notification topic for pipeline events. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DevOps Project** -- a `devops.Project` in the specified compartment with a notification topic for build and deployment event delivery
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the project

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the project in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- An ONS (Oracle Notification Service) topic OCID for receiving pipeline events (build completions, deployment successes/failures). Create the topic beforehand via the OCI Console or CLI.

## Deploy

### Console

Open the deployment store, find **DevOps Project on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Project** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDevopsProject
metadata:
  name: app-pipelines
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  notificationTopicId:
    value: "ocid1.onstopic.oc1..example"
  description: CI/CD project for application pipelines
```

```shell
planton apply -f devops-project.yaml
```

This creates a DevOps project with the specified notification topic. The project name is immutable after creation. All downstream DevOps resources (build pipelines, deploy pipelines, repositories) reference this project by its OCID.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the project to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: devops-compartment
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the DevOps project with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a DevOps project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Project name** -- Derived from `metadata.name`. This is a ForceNew field -- changing it destroys and recreates the project, losing all associated pipelines, repositories, and artifacts. Choose a stable, descriptive name (e.g., `platform-services`, `data-pipelines`).

**Notification topic** -- The `notificationTopicId` is required and receives all project-level events (build started, deployment succeeded, approval required). Use a dedicated topic per project to isolate event streams. The topic can be updated after creation without affecting existing subscriptions.

**Compartment** -- The `compartmentId` is updatable, supporting compartment moves for organizational restructuring. Moving a project moves its visibility scope but does not affect running pipelines.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `project_id` | OCID of the DevOps project | Build pipelines, deploy pipelines, code repositories, triggers |
| `namespace` | Namespace associated with the project | Container registry paths, artifact references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard project** -- A DevOps project with a notification topic and description. The baseline for all CI/CD workflows on OCI. Start from the **Standard Project** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DevOps project