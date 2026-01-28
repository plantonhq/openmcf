# Examples

This document provides practical examples for deploying GCP Artifact Registry repositories using the OpenMCF CLI.

## Example 1: Private Docker Registry (Literal Value)

Create a private Docker registry using a literal project ID:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-docker-private
spec:
  repoFormat: DOCKER
  projectId:
    value: my-gcp-project-123
  region: us-central1
  enablePublicAccess: false
```

Deploy:

```bash
openmcf pulumi up --manifest docker-registry.yaml
```

## Example 2: Private Docker Registry (Reference to GcpProject)

Create a private Docker registry referencing a GcpProject resource:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-docker-private
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

Deploy:

```bash
openmcf pulumi up --manifest docker-registry.yaml
```

## Example 3: Public Docker Registry (for Open Source)

Create a publicly accessible Docker registry for open source projects:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: opensource-docker-public
spec:
  repoFormat: DOCKER
  projectId:
    value: my-gcp-project-123
  region: us-west2
  enablePublicAccess: true
```

Deploy:

```bash
openmcf pulumi up --manifest opensource-registry.yaml
```

## Example 4: Python Package Repository

Create a private repository for Python packages (PyPI):

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-python-packages
spec:
  repoFormat: PYTHON
  projectId:
    value: my-gcp-project-123
  region: us-east1
  enablePublicAccess: false
```

Deploy:

```bash
openmcf pulumi up --manifest python-repo.yaml
```

## Example 5: Maven Repository

Create a private repository for Maven artifacts (Java/JVM):

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-maven-artifacts
spec:
  repoFormat: MAVEN
  projectId:
    value: my-gcp-project-123
  region: us-central1
  enablePublicAccess: false
```

Deploy:

```bash
openmcf pulumi up --manifest maven-repo.yaml
```

## Example 6: NPM Package Repository

Create a private repository for NPM packages:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-npm-packages
spec:
  repoFormat: NPM
  projectId:
    value: my-gcp-project-123
  region: europe-west1
  enablePublicAccess: false
```

Deploy:

```bash
openmcf pulumi up --manifest npm-repo.yaml
```

## Example 7: Go Module Repository

Create a private repository for Go modules:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: company-go-modules
spec:
  repoFormat: GO
  projectId:
    value: my-gcp-project-123
  region: us-west2
  enablePublicAccess: false
```

Deploy:

```bash
openmcf pulumi up --manifest go-repo.yaml
```

## Multi-Region Setup

Deploy repositories in multiple regions for high availability:

```yaml
# US Central region
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-docker-us-central
spec:
  repoFormat: DOCKER
  projectId:
    value: my-gcp-project-123
  region: us-central1
  enablePublicAccess: false
---
# Europe West region
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-docker-europe-west
spec:
  repoFormat: DOCKER
  projectId:
    value: my-gcp-project-123
  region: europe-west1
  enablePublicAccess: false
```

## Best Practices

### Co-location with GKE

Always deploy repositories in the same region as your GKE clusters for:
- **Free egress** - No bandwidth charges
- **Faster image pulls** - Reduced latency
- **GKE Image Streaming** - Only works within the same region

```yaml
# GKE cluster in us-west2 → Repository also in us-west2
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: gke-colocated-registry
spec:
  repoFormat: DOCKER
  projectId:
    value: my-gcp-project-123
  region: us-west2  # Same region as GKE cluster
  enablePublicAccess: false
```

### Naming Conventions

Use clear, descriptive names that indicate:
- Purpose (docker, python, maven)
- Visibility (private, public)
- Environment (if applicable)

Good examples:
- `company-docker-prod`
- `opensource-python-public`
- `backend-maven-dev`

### Security Considerations

- Set `enablePublicAccess: false` for internal artifacts
- Use `enablePublicAccess: true` only for genuinely public open-source projects
- Service accounts with appropriate roles are automatically created:
  - **Reader account** - For pulling artifacts
  - **Writer account** - For pushing artifacts

## Outputs and Access

After deployment, the following outputs are available:

- **repo_name** - Repository identifier
- **repo_hostname** - Repository hostname (e.g., `us-central1-docker.pkg.dev`)
- **repo_url** - Complete repository URL
- **reader_service_account** - Service account for read access (email and key)
- **writer_service_account** - Service account for write access (email and key)

Access outputs:

```bash
openmcf pulumi output --manifest artifact-registry.yaml
```

## Managing Lifecycle

### Update Repository

Modify the YAML and apply changes:

```bash
openmcf pulumi up --manifest artifact-registry.yaml
```

### Destroy Repository

Remove the repository (warning: deletes all artifacts):

```bash
openmcf pulumi destroy --manifest artifact-registry.yaml
```

## Additional Resources

- [GCP Artifact Registry Documentation](https://cloud.google.com/artifact-registry/docs)
- [Supported Formats](https://cloud.google.com/artifact-registry/docs/supported-formats)
- [Regional Locations](https://cloud.google.com/artifact-registry/docs/repositories/repo-locations)
