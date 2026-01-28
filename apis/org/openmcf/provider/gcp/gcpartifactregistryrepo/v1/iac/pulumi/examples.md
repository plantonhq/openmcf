## Usage

### Sample YAML Configuration (Literal Value)

Create a YAML file (`artifact-registry.yaml`) with a literal project ID:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: my-artifact-registry
spec:
  repoFormat: DOCKER
  projectId:
    value: your-gcp-project-id
  region: us-central1
  enablePublicAccess: false
```

### Sample YAML Configuration (Reference to GcpProject)

Reference a GcpProject resource to dynamically get the project ID:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: my-artifact-registry
spec:
  repoFormat: DOCKER
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  region: us-central1
  enablePublicAccess: false
```

> **Note:** Reference resolution (`valueFrom`) is not yet fully implemented. Currently, only literal `value` is supported. References will be resolved in a future version.

### Deploying with CLI

Use the provided CLI tool to deploy the Artifact Registry repositories:

```bash
openmcf pulumi up --manifest artifact-registry.yaml
```

If no Pulumi module is specified, the CLI uses the default module corresponding to the API resource.
