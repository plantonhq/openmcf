# Pulumi Module Examples - GCP Cloud Storage Bucket

This document provides Pulumi-specific deployment examples and patterns for the GcpGcsBucket module. These examples demonstrate how to use the Pulumi CLI and Go-based Pulumi programs to deploy GCS buckets.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Basic Deployment](#basic-deployment)
- [Using the CLI](#using-the-cli)
- [Development Workflow](#development-workflow)
- [Testing and Debugging](#testing-and-debugging)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

```bash
# Install Pulumi CLI
curl -fsSL https://get.pulumi.com | sh

# Install Go 1.21+
# Visit: https://go.dev/doc/install

# Authenticate with GCP
gcloud auth application-default login
```

### GCP Setup

```bash
# Set your GCP project
export GCP_PROJECT=my-gcp-project-123

# Enable required APIs
gcloud services enable storage-api.googleapis.com \
  --project=${GCP_PROJECT}

# Create a service account for Pulumi (optional, for CI/CD)
gcloud iam service-accounts create pulumi-deployer \
  --display-name="Pulumi Deployment Service Account" \
  --project=${GCP_PROJECT}

# Grant necessary permissions
gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
  --member="serviceAccount:pulumi-deployer@${GCP_PROJECT}.iam.gserviceaccount.com" \
  --role="roles/storage.admin"
```

---

## Basic Deployment

### Directory Structure

```
your-project/
├── stack-input.yaml          # GcpGcsBucket resource definition
├── Pulumi.yaml               # Pulumi project configuration (from module)
├── main.go                   # Pulumi program (from module)
└── go.mod                    # Go dependencies
```

### Step 1: Create Stack Input

Create `stack-input.yaml` with your bucket configuration:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: my-app-storage-prod
spec:
  gcp_project_id:
    value: my-gcp-project-123
  location: us-east1
  bucket_name: my-app-storage-prod
  uniform_bucket_level_access_enabled: true
  versioning_enabled: true
  storage_class: STANDARD
  lifecycle_rules:
    - action:
        type: "Delete"
      condition:
        num_newer_versions: 5
  iam_bindings:
    - role: "roles/storage.objectAdmin"
      members:
        - "serviceAccount:app-backend@my-gcp-project-123.iam.gserviceaccount.com"
  gcp_labels:
    environment: production
    team: platform
```

**Using valueFrom for Cross-Resource Reference:**

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: my-app-storage-prod
spec:
  gcp_project_id:
    value_from:
      kind: GcpProject
      name: main-project
      field_path: status.outputs.project_id
  location: us-east1
  bucket_name: my-app-storage-prod
  uniform_bucket_level_access_enabled: true
```

### Step 2: Initialize Pulumi Stack

```bash
# Initialize a new Pulumi stack
pulumi stack init prod

# Set the GCP project (if not using gcloud default)
pulumi config set gcp:project my-gcp-project-123

# Set the region (optional, can be different from bucket location)
pulumi config set gcp:region us-east1
```

### Step 3: Deploy

```bash
# Preview changes
pulumi preview

# Deploy the bucket
pulumi up
```

### Step 4: View Outputs

```bash
# Show stack outputs
pulumi stack output

# Example output:
# bucket_id: projects/my-gcp-project-123/buckets/my-app-storage-prod
```

---

## Using the CLI

### Preview Changes

```bash
# Show what would be created/updated/deleted
pulumi preview

# Save preview to a file
pulumi preview --json > preview.json

# Show detailed diff
pulumi preview --diff
```

### Deploy Stack

```bash
# Deploy with default options
pulumi up

# Auto-approve (for CI/CD)
pulumi up --yes

# Deploy with specific target
pulumi up --target="urn:pulumi:prod::my-project::gcp:storage/bucket:Bucket::my-app-storage-prod"
```

### Update Stack

```bash
# Refresh state from actual infrastructure
pulumi refresh

# Update with new configuration
# (after editing stack-input.yaml)
pulumi up
```

### Destroy Stack

```bash
# Destroy all resources
pulumi destroy

# Auto-approve
pulumi destroy --yes

# Destroy specific resource
pulumi destroy --target="urn:pulumi:prod::my-project::gcp:storage/bucket:Bucket::my-app-storage-prod"
```

### Stack Management

```bash
# List all stacks
pulumi stack ls

# Select a different stack
pulumi stack select dev

# Show stack configuration
pulumi stack output

# Export stack state
pulumi stack export > stack-backup.json

# Import stack state
pulumi stack import < stack-backup.json
```

---

## Development Workflow

### Local Development

Use the provided `debug.sh` script for local testing:

```bash
# Navigate to Pulumi module directory
cd iac/pulumi

# Run debug script
./debug.sh
```

The debug script:
1. Looks for `stack-input.yaml` in the module directory
2. Sets up environment variables
3. Runs Pulumi preview
4. Optionally deploys if confirmed

### Testing with Different Configurations

```bash
# Test with minimal configuration
cat > stack-input-minimal.yaml <<EOF
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: test-bucket-minimal
spec:
  gcp_project_id:
    value: my-gcp-project-123
  location: us-east1
  bucket_name: test-bucket-minimal
  uniform_bucket_level_access_enabled: true
EOF

# Deploy test stack
export STACK_INPUT_FILE=stack-input-minimal.yaml
pulumi up --stack=test-minimal
```

### Iterative Development

```bash
# Watch mode (re-run on file changes)
pulumi watch

# This will automatically preview changes when stack-input.yaml is modified
```

---

## Testing and Debugging

### Enable Debug Logging

```bash
# Set debug log level
export PULUMI_LOG_LEVEL=debug

# Enable GCP API logging
export TF_LOG=DEBUG

# Run Pulumi with verbose output
pulumi up --logtostderr -v=9
```

### Validate Configuration

```bash
# Validate stack input YAML
# (using yq or similar tool)
yq eval stack-input.yaml

# Validate against protobuf schema
# (if buf CLI is installed)
buf validate stack-input.yaml
```

### Dry Run Testing

```bash
# Preview without state changes
pulumi preview --expect-no-changes

# Preview with detailed plan
pulumi preview --show-config --show-replacement-steps
```

### Integration Testing

Create a test script `test-deployment.sh`:

```bash
#!/bin/bash
set -e

STACK_NAME="test-$(date +%s)"
BUCKET_NAME="test-bucket-${STACK_NAME}"

# Create test stack input
cat > stack-input-test.yaml <<EOF
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: ${BUCKET_NAME}
spec:
  gcp_project_id:
    value: my-gcp-project-123
  location: us-east1
  bucket_name: ${BUCKET_NAME}
  uniform_bucket_level_access_enabled: true
EOF

# Deploy
export STACK_INPUT_FILE=stack-input-test.yaml
pulumi stack init ${STACK_NAME}
pulumi up --yes

# Test bucket accessibility
gsutil ls gs://${BUCKET_NAME}

# Cleanup
pulumi destroy --yes
pulumi stack rm ${STACK_NAME} --yes
rm stack-input-test.yaml
```

---

## Common Patterns

### Multi-Environment Deployment

Use separate stacks for different environments:

```bash
# Create stack inputs for each environment
cat > stack-input-dev.yaml <<EOF
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-storage-dev
spec:
  gcp_project_id:
    value: my-dev-project
  location: us-east1
  bucket_name: app-storage-dev
  uniform_bucket_level_access_enabled: true
  storage_class: STANDARD
  lifecycle_rules:
    - action:
        type: "Delete"
      condition:
        age_days: 30  # Aggressive cleanup for dev
  gcp_labels:
    environment: development
EOF

cat > stack-input-prod.yaml <<EOF
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-storage-prod
spec:
  gcp_project_id:
    value: my-prod-project
  location: us-east1
  bucket_name: app-storage-prod
  uniform_bucket_level_access_enabled: true
  storage_class: STANDARD
  versioning_enabled: true
  lifecycle_rules:
    - action:
        type: "Delete"
      condition:
        num_newer_versions: 10
  gcp_labels:
    environment: production
EOF

# Deploy to dev
export STACK_INPUT_FILE=stack-input-dev.yaml
pulumi stack select dev
pulumi up --yes

# Deploy to prod
export STACK_INPUT_FILE=stack-input-prod.yaml
pulumi stack select prod
pulumi up --yes
```

### CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Deploy GCS Bucket

on:
  push:
    branches: [main]
    paths:
      - 'stack-input-*.yaml'
      - '.github/workflows/deploy-bucket.yaml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Pulumi
        uses: pulumi/actions@v4
        
      - name: Authenticate to GCP
        uses: google-github-actions/auth@v1
        with:
          credentials_json: ${{ secrets.GCP_CREDENTIALS }}
          
      - name: Deploy Stack
        env:
          STACK_INPUT_FILE: stack-input-prod.yaml
          PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}
        run: |
          cd iac/pulumi
          pulumi stack select prod --create
          pulumi up --yes
```

### Blue-Green Deployment

Deploy a new bucket alongside existing one, then switch traffic:

```yaml
# stack-input-blue.yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-storage-blue
spec:
  gcp_project_id: my-gcp-project-123
  location: us-east1
  uniform_bucket_level_access_enabled: true
  # ... configuration ...
```

```yaml
# stack-input-green.yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-storage-green
spec:
  gcp_project_id: my-gcp-project-123
  location: us-east1
  uniform_bucket_level_access_enabled: true
  # ... configuration ...
```

```bash
# Deploy green alongside blue
export STACK_INPUT_FILE=stack-input-green.yaml
pulumi stack init green
pulumi up --yes

# Test green deployment
gsutil ls gs://app-storage-green

# Switch application to use green bucket
# Update application configuration or DNS

# Destroy blue after validation
pulumi stack select blue
pulumi destroy --yes
```

### Programmatic Stack Input

Generate stack input programmatically:

```bash
# generate-stack-input.sh
#!/bin/bash

ENVIRONMENT=$1
BUCKET_NAME="app-storage-${ENVIRONMENT}"
PROJECT_ID="my-${ENVIRONMENT}-project"

cat > stack-input-${ENVIRONMENT}.yaml <<EOF
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: ${BUCKET_NAME}
spec:
  gcp_project_id:
    value: ${PROJECT_ID}
  location: us-east1
  bucket_name: ${BUCKET_NAME}
  uniform_bucket_level_access_enabled: true
  gcp_labels:
    environment: ${ENVIRONMENT}
    generated: "true"
    timestamp: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

echo "Generated stack-input-${ENVIRONMENT}.yaml"
```

Usage:

```bash
./generate-stack-input.sh dev
./generate-stack-input.sh staging
./generate-stack-input.sh prod

# Deploy each environment
for env in dev staging prod; do
  export STACK_INPUT_FILE=stack-input-${env}.yaml
  pulumi stack select ${env} --create
  pulumi up --yes
done
```

---

## Troubleshooting

### Common Errors

#### Error: Bucket Already Exists

```
error: resource already exists (HTTP 409)
```

**Solution:** Bucket names are globally unique. Choose a different name or check if bucket exists in a different project.

```bash
# Check if bucket exists
gsutil ls gs://my-bucket-name

# If it's your bucket in different project, import it
pulumi import gcp:storage/bucket:Bucket my-bucket gs://my-bucket-name
```

#### Error: Permission Denied

```
error: googleapi: Error 403: Caller does not have storage.buckets.create permission
```

**Solution:** Grant the deploying service account appropriate permissions:

```bash
gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/storage.admin"
```

#### Error: CMEK Key Not Found

```
error: encryption key not found or inaccessible
```

**Solution:** Ensure the KMS key exists and GCS service account has access:

```bash
# Get GCS service account
PROJECT_NUMBER=$(gcloud projects describe ${GCP_PROJECT} --format="value(projectNumber)")
GCS_SA="service-${PROJECT_NUMBER}@gs-project-accounts.iam.gserviceaccount.com"

# Grant KMS permissions
gcloud kms keys add-iam-policy-binding ${KEY_NAME} \
  --location=${LOCATION} \
  --keyring=${KEY_RING} \
  --member="serviceAccount:${GCS_SA}" \
  --role="roles/cloudkms.cryptoKeyEncrypterDecrypter"
```

### Debug Workflow

1. **Enable verbose logging:**

```bash
export PULUMI_LOG_LEVEL=debug
export TF_LOG=DEBUG
```

2. **Check stack state:**

```bash
pulumi stack export > state.json
cat state.json | jq '.deployment.resources'
```

3. **Validate GCP credentials:**

```bash
gcloud auth application-default print-access-token
gcloud config get-value project
```

4. **Test GCS access manually:**

```bash
gsutil mb gs://test-bucket-${RANDOM}
gsutil ls
gsutil rb gs://test-bucket-*
```

5. **Check for drift:**

```bash
pulumi refresh
pulumi preview --expect-no-changes
```

### Performance Optimization

#### Reduce Deployment Time

```bash
# Use parallel operations
pulumi up --parallel 10

# Skip preview
pulumi up --yes --skip-preview
```

#### State Management

```bash
# Use remote state backend for team collaboration
pulumi login s3://my-pulumi-state-bucket

# Or use Pulumi Cloud
pulumi login
```

---

## Further Reading

- [Pulumi Documentation](https://www.pulumi.com/docs/)
- [Pulumi GCP Provider](https://www.pulumi.com/registry/packages/gcp/)
- [Module README](README.md)
- [Module Overview](overview.md)
- [Component Examples](../../examples.md)

---

## Need Help?

For issues specific to Pulumi deployment:
1. Check [Pulumi Community](https://pulumi.com/community/)
2. Review [GCP Provider Issues](https://github.com/pulumi/pulumi-gcp/issues)
3. Consult the [module overview](overview.md) for architecture details


