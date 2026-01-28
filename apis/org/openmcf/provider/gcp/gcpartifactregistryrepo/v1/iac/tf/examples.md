# Terraform Examples for GCP Artifact Registry Repository

This document provides Terraform-specific examples for deploying GCP Artifact Registry repositories using the OpenMCF CLI with Terraform (OpenTofu) backend.

## Prerequisites

- OpenMCF CLI installed
- GCP credentials configured
- Terraform state backend configured (S3 or GCS)

## Example 1: Private Docker Registry

Create a private Docker registry for your organization's container images.

### Manifest: `private-docker.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-docker-prod
spec:
  repoFormat: DOCKER
  projectId:
    value: acme-gcp-project-123
  region: us-central1
  enablePublicAccess: false
```

### Deploy

```bash
# Initialize Terraform state
openmcf tofu init --manifest private-docker.yaml --backend-type s3 \
  --backend-config="bucket=acme-terraform-state" \
  --backend-config="dynamodb_table=acme-terraform-locks" \
  --backend-config="region=us-east-1" \
  --backend-config="key=gcp/artifact-registry/acme-docker-prod.tfstate"

# Plan changes
openmcf tofu plan --manifest private-docker.yaml

# Apply configuration
openmcf tofu apply --manifest private-docker.yaml --auto-approve
```

### Expected Outputs

```
repo_name = "acme-docker-prod-docker"
repo_hostname = "us-central1-docker.pkg.dev"
repo_url = "us-central1-docker.pkg.dev/acme-gcp-project-123/acme-docker-prod-docker"
reader_service_account_email = "acme-docker-prod-abc123-ro@acme-gcp-project-123.iam.gserviceaccount.com"
writer_service_account_email = "acme-docker-prod-abc123-rw@acme-gcp-project-123.iam.gserviceaccount.com"
```

### Use in CI/CD

```bash
# Push image using writer service account
echo $WRITER_SA_KEY_BASE64 | base64 -d > key.json
gcloud auth activate-service-account --key-file=key.json
docker push us-central1-docker.pkg.dev/acme-gcp-project-123/acme-docker-prod-docker/app:latest
```

## Example 2: Public Docker Registry for Open Source

Create a publicly accessible Docker registry for open source projects.

### Manifest: `opensource-docker.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: opensource-project-public
spec:
  repoFormat: DOCKER
  projectId:
    value: opensource-gcp-project
  region: us-west2
  enablePublicAccess: true
```

### Deploy

```bash
# Initialize
openmcf tofu init --manifest opensource-docker.yaml --backend-type s3 \
  --backend-config="bucket=opensource-terraform-state" \
  --backend-config="key=gcp/opensource-docker.tfstate"

# Apply
openmcf tofu apply --manifest opensource-docker.yaml --auto-approve
```

### Public Access

Anyone can pull images without authentication:

```bash
docker pull us-west2-docker.pkg.dev/opensource-gcp-project/opensource-project-public-docker/app:v1.0.0
```

## Example 3: Python Package Repository

Create a private repository for Python packages (PyPI-compatible).

### Manifest: `python-packages.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-python-packages
spec:
  repoFormat: PYTHON
  projectId:
    value: acme-gcp-project-123
  region: us-east1
  enablePublicAccess: false
```

### Deploy

```bash
openmcf tofu init --manifest python-packages.yaml --backend-type gcs \
  --backend-config="bucket=acme-terraform-state-gcs" \
  --backend-config="prefix=gcp/python-packages"

openmcf tofu apply --manifest python-packages.yaml --auto-approve
```

### Upload Python Package

```bash
# Configure pip/twine
export TWINE_USERNAME=oauth2accesstoken
export TWINE_PASSWORD=$(gcloud auth print-access-token)

# Upload package
twine upload --repository-url https://us-east1-python.pkg.dev/acme-gcp-project-123/acme-python-packages/ dist/*
```

### Install from Repository

```bash
# Configure pip
pip install \
  --index-url https://oauth2accesstoken:$(gcloud auth print-access-token)@us-east1-python.pkg.dev/acme-gcp-project-123/acme-python-packages/simple/ \
  your-package-name
```

## Example 4: Maven Repository for Java/JVM

Create a private Maven repository for Java artifacts.

### Manifest: `maven-artifacts.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-maven-libs
spec:
  repoFormat: MAVEN
  projectId:
    value: acme-gcp-project-123
  region: us-central1
  enablePublicAccess: false
```

### Deploy

```bash
openmcf tofu init --manifest maven-artifacts.yaml --backend-type s3 \
  --backend-config="bucket=acme-terraform-state" \
  --backend-config="key=gcp/maven-artifacts.tfstate"

openmcf tofu apply --manifest maven-artifacts.yaml --auto-approve
```

### Maven Configuration

Add to `~/.m2/settings.xml`:

```xml
<settings>
  <servers>
    <server>
      <id>artifact-registry</id>
      <username>_json_key</username>
      <password>${env.MAVEN_REPO_KEY}</password>
    </server>
  </servers>
</settings>
```

Add to `pom.xml`:

```xml
<repositories>
  <repository>
    <id>artifact-registry</id>
    <url>https://us-central1-maven.pkg.dev/acme-gcp-project-123/acme-maven-libs</url>
  </repository>
</repositories>
```

## Example 5: NPM Package Repository

Create a private repository for NPM packages.

### Manifest: `npm-packages.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-npm-packages
spec:
  repoFormat: NPM
  projectId:
    value: acme-gcp-project-123
  region: europe-west1
  enablePublicAccess: false
```

### Deploy

```bash
openmcf tofu init --manifest npm-packages.yaml --backend-type s3 \
  --backend-config="bucket=acme-terraform-state" \
  --backend-config="key=gcp/npm-packages.tfstate"

openmcf tofu apply --manifest npm-packages.yaml --auto-approve
```

### NPM Configuration

Configure `.npmrc`:

```ini
@acme:registry=https://europe-west1-npm.pkg.dev/acme-gcp-project-123/acme-npm-packages/
//europe-west1-npm.pkg.dev/acme-gcp-project-123/acme-npm-packages/:_password=${NPM_REPO_PASSWORD_BASE64}
//europe-west1-npm.pkg.dev/acme-gcp-project-123/acme-npm-packages/:username=_json_key_base64
//europe-west1-npm.pkg.dev/acme-gcp-project-123/acme-npm-packages/:email=not.used@example.com
//europe-west1-npm.pkg.dev/acme-gcp-project-123/acme-npm-packages/:always-auth=true
```

Publish package:

```bash
npm publish
```

## Example 6: Go Module Repository

Create a private repository for Go modules.

### Manifest: `go-modules.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-go-modules
spec:
  repoFormat: GO
  projectId:
    value: acme-gcp-project-123
  region: us-west2
  enablePublicAccess: false
```

### Deploy

```bash
openmcf tofu init --manifest go-modules.yaml --backend-type gcs \
  --backend-config="bucket=acme-terraform-state-gcs" \
  --backend-config="prefix=gcp/go-modules"

openmcf tofu apply --manifest go-modules.yaml --auto-approve
```

### Go Configuration

Configure Go to use the repository:

```bash
# Set GOPROXY
export GOPROXY=https://us-west2-go.pkg.dev/acme-gcp-project-123/acme-go-modules,direct

# Configure authentication
gcloud auth application-default login
```

## Example 7: Multi-Region Deployment

Deploy repositories in multiple regions for high availability and performance.

### Manifest: `multi-region.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-docker-us
spec:
  repoFormat: DOCKER
  projectId:
    value: acme-gcp-project-123
  region: us-central1
  enablePublicAccess: false
---
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-docker-eu
spec:
  repoFormat: DOCKER
  projectId:
    value: acme-gcp-project-123
  region: europe-west1
  enablePublicAccess: false
---
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-docker-asia
spec:
  repoFormat: DOCKER
  projectId:
    value: acme-gcp-project-123
  region: asia-southeast1
  enablePublicAccess: false
```

### Deploy All Regions

```bash
# Deploy to all regions
openmcf tofu init --manifest multi-region.yaml --backend-type s3 \
  --backend-config="bucket=acme-terraform-state" \
  --backend-config="key=gcp/multi-region-registries.tfstate"

openmcf tofu apply --manifest multi-region.yaml --auto-approve
```

## Example 8: Generic Artifact Repository

Create a repository for generic artifacts (tarballs, binaries, etc.).

### Manifest: `generic-artifacts.yaml`

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpArtifactRegistryRepo
metadata:
  name: acme-generic-artifacts
spec:
  repoFormat: GENERIC
  projectId:
    value: acme-gcp-project-123
  region: us-central1
  enablePublicAccess: false
```

### Deploy

```bash
openmcf tofu init --manifest generic-artifacts.yaml --backend-type s3 \
  --backend-config="bucket=acme-terraform-state" \
  --backend-config="key=gcp/generic-artifacts.tfstate"

openmcf tofu apply --manifest generic-artifacts.yaml --auto-approve
```

### Upload Generic Artifacts

```bash
# Upload a file
curl -X POST \
  -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  -F "package=@myfile.tar.gz" \
  "https://us-central1-generic.pkg.dev/acme-gcp-project-123/acme-generic-artifacts/mypackage/1.0.0/myfile.tar.gz"
```

## Best Practices

### 1. State Backend Configuration

Always use remote state backend for team collaboration:

**AWS S3 Backend:**
```bash
--backend-type s3 \
--backend-config="bucket=my-terraform-state" \
--backend-config="dynamodb_table=terraform-locks" \
--backend-config="region=us-east-1" \
--backend-config="key=path/to/state.tfstate"
```

**GCS Backend:**
```bash
--backend-type gcs \
--backend-config="bucket=my-terraform-state-gcs" \
--backend-config="prefix=gcp/artifact-registry"
```

### 2. Co-location with Compute

Deploy repositories in the same region as your workloads:

| Workload Location | Repository Region | Benefit |
|------------------|-------------------|---------|
| us-central1-a (GKE) | us-central1 | Free egress, faster pulls |
| europe-west1-b (GKE) | europe-west1 | Free egress, GKE Image Streaming |
| asia-southeast1 (Cloud Run) | asia-southeast1 | Reduced latency |

### 3. Service Account Key Management

**Store keys securely:**

```bash
# Get outputs (including sensitive keys)
openmcf tofu output --manifest manifest.yaml

# Store in GCP Secret Manager
gcloud secrets create writer-sa-key \
  --data-file=<(openmcf tofu output -raw writer_service_account_key_base64 | base64 -d)

# Use in CI/CD
gcloud secrets versions access latest --secret=writer-sa-key > sa-key.json
```

**Or use Workload Identity (preferred for GKE):**
- No service account keys needed
- Automatic authentication via Kubernetes service accounts

### 4. Cost Optimization

- **Storage:** $0.10/GB/month
- **Egress:** FREE within same region, paid cross-region
- **Strategy:** Always co-locate with primary consumers

### 5. Naming Conventions

Use descriptive names:
- `{company}-{format}-{env}` → `acme-docker-prod`
- `{project}-{format}-{region-code}` → `webapp-npm-us`
- `{purpose}-{format}-{visibility}` → `opensource-docker-public`

## Destroying Resources

### Single Repository

```bash
openmcf tofu destroy --manifest manifest.yaml --auto-approve
```

### All Repositories in a Multi-Region Setup

```bash
openmcf tofu destroy --manifest multi-region.yaml --auto-approve
```

**Warning:** Destroying a repository permanently deletes all artifacts. Ensure you have backups if needed.

## Troubleshooting

### Error: Repository Already Exists

If you encounter "already exists" errors:

1. **Import existing repository:**
   ```bash
   terraform import google_artifact_registry_repository.repo \
     projects/PROJECT_ID/locations/REGION/repositories/REPO_ID
   ```

2. **Or use a different name** in `metadata.name`

### Error: Permission Denied

Ensure your service account has required roles:
- `roles/artifactregistry.admin`
- `roles/iam.serviceAccountAdmin`
- `roles/iam.serviceAccountKeyAdmin`

Grant roles:
```bash
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member=serviceAccount:terraform@PROJECT_ID.iam.gserviceaccount.com \
  --role=roles/artifactregistry.admin
```

### Error: State Lock Failed

If state is locked (from failed previous run):

```bash
# Force unlock (use with caution)
terraform force-unlock LOCK_ID
```

## Additional Resources

- [GCP Artifact Registry Documentation](https://cloud.google.com/artifact-registry/docs)
- [Terraform GCP Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [OpenMCF CLI](https://openmcf.org)
- [OpenTofu Documentation](https://opentofu.org/docs/)

