---
title: "Unified Commands"
description: "kubectl-style unified commands that automatically detect your provisioner - simplify your workflow with apply, destroy, init, plan, and refresh"
icon: "rocket"
order: 2
---

# Unified Commands

The unified commands provide a kubectl-like experience by automatically detecting the IaC provisioner from your manifest and routing to the appropriate tool (Pulumi, Tofu, or Terraform).

**Available unified commands**: `apply`, `destroy`, `init`, `plan`/`preview`, `refresh`

---

## Why Unified Commands?

### The Problem

Previously, you had to remember different commands for different provisioners:

```bash
# Pulumi
project-planton pulumi init -f app.yaml
project-planton pulumi preview -f app.yaml
project-planton pulumi up -f app.yaml --stack org/project/env
project-planton pulumi refresh -f app.yaml
project-planton pulumi destroy -f app.yaml

# OpenTofu
project-planton tofu init -f app.yaml
project-planton tofu plan -f app.yaml
project-planton tofu apply -f app.yaml
project-planton tofu refresh -f app.yaml
project-planton tofu destroy -f app.yaml

# Different commands, different flags, cognitive overhead!
```

### The Solution

Now, use the same commands regardless of provisioner:

```bash
# Works for Pulumi, Tofu, or Terraform - complete lifecycle
project-planton init -f app.yaml
project-planton plan -f app.yaml       # or 'preview'
project-planton apply -f app.yaml
project-planton refresh -f app.yaml
project-planton destroy -f app.yaml
```

The CLI automatically:
1. Reads the `project-planton.org/provisioner` label from your manifest
2. Routes to the appropriate provisioner
3. Passes through all relevant flags

---

## The Provisioner Label

Add this label to your manifest metadata:

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesPostgres
metadata:
  name: my-database
  labels:
    project-planton.org/provisioner: pulumi  # or tofu or terraform
spec:
  # ... your spec
```

**Supported values** (case-insensitive):
- `pulumi` - Use Pulumi for deployment
- `tofu` - Use OpenTofu for deployment
- `terraform` - Use Terraform for deployment

---

## Commands

### apply

Deploy infrastructure changes.

**Usage**:

```bash
project-planton apply -f <file> [flags]
project-planton apply -f <file> [flags]
```

**Examples**:

```bash
# Basic usage
project-planton apply -f database.yaml

# With kustomize
project-planton apply --kustomize-dir services/api --overlay prod

# With field overrides
project-planton apply -f app.yaml --set spec.replicas=5

# Auto-approve (Pulumi)
project-planton apply -f app.yaml --yes

# Auto-approve (Tofu/Terraform)
project-planton apply -f app.yaml --auto-approve
```

**What it does**:
1. Loads and validates your manifest
2. Detects provisioner from `project-planton.org/provisioner` label
3. If label missing, prompts interactively (defaults to Pulumi)
4. Routes to appropriate provisioner:
   - **Pulumi**: Runs `pulumi update`
   - **Tofu**: Runs `tofu apply`
   - **Terraform**: Runs `terraform apply`

---

### destroy

Teardown infrastructure.

**Usage**:

```bash
project-planton destroy -f <file> [flags]
project-planton delete -f <file> [flags]  # kubectl-style alias
```

**Examples**:

```bash
# Basic usage
project-planton destroy -f database.yaml

# kubectl-style delete
project-planton delete -f database.yaml

# With kustomize
project-planton destroy --kustomize-dir services/api --overlay staging

# Auto-approve
project-planton destroy -f app.yaml --auto-approve
```

**What it does**:
1. Loads and validates your manifest
2. Detects provisioner from label
3. Routes to appropriate destroy operation:
   - **Pulumi**: Runs `pulumi destroy`
   - **Tofu**: Runs `tofu destroy`
   - **Terraform**: Runs `terraform destroy`

---

### init

Initialize infrastructure backend or stack.

**Usage**:

```bash
project-planton init -f <file> [flags]
```

**Examples**:

```bash
# Basic usage
project-planton init -f database.yaml

# With kustomize
project-planton init --kustomize-dir services/api --overlay prod

# With OpenTofu backend config
project-planton init -f app.yaml \
    --backend-type s3 \
    --backend-config bucket=my-terraform-state \
    --backend-config key=app/terraform.tfstate \
    --backend-config region=us-west-2
```

**What it does**:
1. Loads and validates your manifest
2. Detects provisioner from label
3. Routes to appropriate initialization:
   - **Pulumi**: Creates stack if it doesn't exist (idempotent)
   - **Tofu**: Initializes backend, downloads providers
   - **Terraform**: Not yet implemented

**When to use**:
- First time deploying a resource
- After cleaning local state/cache
- After changing backend configuration

---

### plan / preview

Preview infrastructure changes without applying them.

**Usage**:

```bash
project-planton plan -f <file> [flags]
project-planton preview -f <file> [flags]  # Pulumi-style alias
```

**Examples**:

```bash
# Basic preview
project-planton plan -f database.yaml

# Using Pulumi-style alias
project-planton preview -f database.yaml

# With kustomize
project-planton plan --kustomize-dir services/api --overlay staging

# Preview destroy plan (OpenTofu)
project-planton plan -f app.yaml --destroy

# Show detailed diffs (Pulumi)
project-planton plan -f app.yaml --diff
```

**What it does**:
1. Loads and validates your manifest
2. Detects provisioner from label
3. Routes to appropriate preview operation:
   - **Pulumi**: Runs `pulumi preview` (dry-run of update)
   - **Tofu**: Runs `tofu plan`
   - **Terraform**: Not yet implemented

**Output**: Shows what resources would be created, modified, or deleted without making any changes.

**When to use**:
- Before applying changes to production
- In pull request CI checks
- To understand impact of manifest changes
- For review and approval workflows

---

### refresh

Sync state with cloud reality.

**Usage**:

```bash
project-planton refresh -f <file> [flags]
```

**Examples**:

```bash
# Basic refresh
project-planton refresh -f database.yaml

# With kustomize
project-planton refresh --kustomize-dir services/api --overlay prod

# Show detailed diffs (Pulumi)
project-planton refresh -f app.yaml --diff
```

**What it does**:
1. Loads and validates your manifest
2. Detects provisioner from label
3. Routes to appropriate refresh operation:
   - **Pulumi**: Runs `pulumi refresh`
   - **Tofu**: Runs `tofu refresh`
   - **Terraform**: Not yet implemented
4. Queries cloud provider for current resource state
5. Updates state file to match reality
6. **Does NOT modify any cloud resources** (read-only)

**When to use**:
- After manual changes made outside IaC (console, CLI, other tools)
- Before applying updates to ensure state accuracy
- After failed deployments to resynchronize state
- When troubleshooting drift between desired and actual state

---

## Interactive Provisioner Selection

If your manifest doesn't have the `project-planton.org/provisioner` label, the CLI prompts you:

```bash
$ project-planton apply -f database.yaml

✓ Manifest loaded
✓ Manifest validated
• Detecting provisioner...
ℹ Provisioner not specified in manifest
Select provisioner [Pulumi]/tofu/terraform: 
```

Simply press **Enter** to use Pulumi (the default), or type your choice.

**Tips**:
- Input is case-insensitive (`Pulumi`, `pulumi`, `PULUMI` all work)
- The prompt only appears when the label is missing
- Add the label to your manifests for a fully automated workflow

---

## Supported Flags

All unified commands support flags from their respective provisioners.

### Common Flags (All Commands)

| Flag | Description |
|------|-------------|
| `-f, -f <path>` | Path to manifest file (kubectl-style `-f` shorthand) |
| `--kustomize-dir <dir>` | Kustomize base directory |
| `--overlay <name>` | Kustomize overlay (prod, dev, staging, etc.) |
| `--set key=value` | Override manifest fields (repeatable) |
| `--module-dir <path>` | Override IaC module directory |

### Pulumi-Specific Flags

| Flag | Commands | Description |
|------|----------|-------------|
| `--stack <org>/<project>/<stack>` | All | Override stack FQDN (or use manifest label) |
| `--yes` | apply, destroy | Auto-approve without confirmation |
| `--diff` | apply, destroy, plan, refresh | Show detailed resource diffs |

### Tofu/Terraform-Specific Flags

| Flag | Commands | Description |
|------|----------|-------------|
| `--auto-approve` | apply, destroy | Skip interactive approval |
| `--destroy` | plan | Create destroy plan instead of apply plan |
| `--reconfigure` | init, apply, destroy, plan, refresh | Reconfigure backend, ignoring saved configuration |
| `--backend-type <type>` | init | Backend type (s3, gcs, local, etc.) |
| `--backend-config <key=value>` | init | Backend configuration (repeatable) |

### Provider Credentials

By default, the CLI reads credentials from **environment variables** - the same ones used by the cloud provider CLIs (AWS CLI, gcloud, az, etc.). No additional configuration is needed if your environment is already set up.

**Default Behavior (Environment Variables)**:

```bash
# AWS - uses AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_DEFAULT_REGION
project-planton apply -f aws-vpc.yaml

# GCP - uses GOOGLE_APPLICATION_CREDENTIALS
project-planton apply -f gcp-cluster.yaml

# Azure - uses ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID, ARM_SUBSCRIPTION_ID
project-planton apply -f azure-aks.yaml
```

**Explicit Credentials (Override with `-p` flag)**:

| Flag | Description |
|------|-------------|
| `-p, --provider-config <file>` | Provider credentials file (overrides environment variables) |

Use the `-p` flag when you need to explicitly specify credentials, such as for multi-account scenarios:

```bash
# Use explicit AWS credentials file
project-planton apply -f aws-vpc.yaml -p ~/.config/aws-prod-creds.yaml

# Use explicit GCP credentials file
project-planton apply -f gcp-cluster.yaml -p ~/.config/gcp-prod-creds.yaml
```

The CLI auto-detects which provider is needed from your manifest's `apiVersion`. See the [Credentials Guide](/docs/guides/credentials) for environment variable details per provider.

---

## Complete Examples

### Example 1: Pulumi PostgreSQL Database

**Manifest** (`postgres.yaml`):

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
  labels:
    project-planton.org/provisioner: pulumi
    pulumi.project-planton.org/stack.name: prod.PostgresKubernetes.app-database
spec:
  container:
    replicas: 1
    resources:
      limits:
        cpu: 1000m
        memory: 2Gi
```

**Commands**:

```bash
# Deploy
project-planton apply -f postgres.yaml

# Update with more replicas
project-planton apply -f postgres.yaml --set spec.container.replicas=3

# Destroy
project-planton destroy -f postgres.yaml
```

---

### Example 2: Tofu AWS VPC

**Manifest** (`vpc.yaml`):

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsVpc
metadata:
  name: production-vpc
  labels:
    project-planton.org/provisioner: tofu
    tofu.project-planton.org/backend.type: s3
    tofu.project-planton.org/backend.bucket: terraform-state
    tofu.project-planton.org/backend.key: vpc/prod.tfstate
    tofu.project-planton.org/backend.region: us-west-2
spec:
  cidrBlock: 10.0.0.0/16
  region: us-west-2
```

**Commands**:

```bash
# Deploy
project-planton apply -f vpc.yaml --auto-approve

# If backend config changed, use --reconfigure
project-planton apply -f vpc.yaml --auto-approve --reconfigure

# Destroy
project-planton delete -f vpc.yaml --auto-approve
```

---

### Example 3: Multi-Environment with Kustomize

**Directory structure**:

```
services/api/
├── base/
│   └── kustomization.yaml
└── overlays/
    ├── dev/
    │   └── kustomization.yaml
    ├── staging/
    │   └── kustomization.yaml
    └── prod/
        └── kustomization.yaml
```

**Deploy to all environments**:

```bash
for env in dev staging prod; do
    echo "Deploying to $env..."
    project-planton apply \
        --kustomize-dir services/api \
        --overlay $env \
        --yes
done
```

---

### Example 4: Complete Infrastructure Lifecycle

```bash
# Full workflow from init to destroy
cd my-infrastructure/

# 1. Initialize backend/stack
project-planton init -f database.yaml

# 2. Preview changes before applying
project-planton plan -f database.yaml

# 3. Apply infrastructure
project-planton apply -f database.yaml --yes

# 4. Refresh state after manual changes
project-planton refresh -f database.yaml

# 5. Destroy when done
project-planton destroy -f database.yaml --yes
```

---

### Example 5: CI/CD Pipeline

```bash
#!/bin/bash
# deploy.sh - CI/CD deployment script

set -e

# Variables from CI environment
IMAGE_TAG="${CI_COMMIT_SHA}"
ENVIRONMENT="${CI_ENVIRONMENT_NAME}"

# Initialize (idempotent)
project-planton init -f deployment.yaml

# Preview changes in pull requests
if [ "$CI_PIPELINE_SOURCE" = "merge_request_event" ]; then
    project-planton plan -f deployment.yaml \
        --set spec.container.image.tag="$IMAGE_TAG"
    exit 0
fi

# Deploy to environment
project-planton apply \
    -f deployment.yaml \
    --set spec.container.image.tag="$IMAGE_TAG" \
    --set metadata.labels.environment="$ENVIRONMENT" \
    --yes

echo "Deployed version $IMAGE_TAG to $ENVIRONMENT"
```

---

## Migration from Provisioner-Specific Commands

The unified commands are **fully backward compatible**. You can migrate gradually:

### Before (Provisioner-Specific)

```bash
# Pulumi - complete workflow
project-planton pulumi init -f app.yaml --stack org/project/dev
project-planton pulumi preview -f app.yaml --stack org/project/dev
project-planton pulumi up -f app.yaml --stack org/project/dev
project-planton pulumi refresh -f app.yaml --stack org/project/dev
project-planton pulumi destroy -f app.yaml --stack org/project/dev

# OpenTofu - complete workflow
project-planton tofu init -f app.yaml
project-planton tofu plan -f app.yaml
project-planton tofu apply -f app.yaml
project-planton tofu refresh -f app.yaml
project-planton tofu destroy -f app.yaml
```

### After (Unified)

```bash
# Works for both Pulumi and OpenTofu!
project-planton init -f app.yaml
project-planton plan -f app.yaml        # or 'preview'
project-planton apply -f app.yaml
project-planton refresh -f app.yaml
project-planton destroy -f app.yaml     # or 'delete'
```

### Migration Steps

1. **Add provisioner label** to your manifests:
   ```yaml
   metadata:
     labels:
       project-planton.org/provisioner: pulumi  # or tofu
   ```

2. **Replace commands** in your scripts:
   - `pulumi init` → `init`
   - `pulumi preview` → `plan` or `preview`
   - `pulumi up` → `apply`
   - `pulumi refresh` → `refresh`
   - `pulumi destroy` → `destroy`
   - `tofu init` → `init`
   - `tofu plan` → `plan`
   - `tofu apply` → `apply`
   - `tofu refresh` → `refresh`
   - `tofu destroy` → `destroy`

3. **Update flags**:
   - `-f` → `-f` (both work, `-f` is shorter)
   - `--stack` → can be in manifest labels now
   - All other flags work the same

---

## Benefits

### 1. Complete Lifecycle Coverage

Unified commands for the entire infrastructure lifecycle:

```bash
project-planton init -f <manifest>      # Initialize
project-planton plan -f <manifest>      # Preview
project-planton apply -f <manifest>     # Deploy
project-planton refresh -f <manifest>   # Sync
project-planton destroy -f <manifest>   # Teardown
```

### 2. Simplified Mental Model

One command pattern for all provisioners, all operations.

### 3. kubectl-like Experience

Familiar kubectl patterns throughout:

```bash
kubectl apply -f deployment.yaml
project-planton apply -f deployment.yaml

kubectl delete -f deployment.yaml
project-planton delete -f deployment.yaml
```

### 4. Easier Automation

Write scripts that work regardless of provisioner:

```bash
for manifest in manifests/*.yaml; do
    project-planton init -f "$manifest"
    project-planton plan -f "$manifest"
    project-planton apply -f "$manifest" --yes
done
```

### 5. Better CI/CD Integration

Use the same commands across all pipelines:

```bash
# Pull Request - preview only
project-planton plan -f app.yaml

# Main branch - deploy
project-planton apply -f app.yaml --yes
```

### 6. Lower Barrier to Entry

New team members learn one command set, not multiple provisioner-specific patterns.

### 7. Gradual Migration

All existing commands still work - migrate at your own pace.

---

## Troubleshooting

### "Invalid provisioner value"

**Error**:
```
Invalid provisioner in manifest: invalid provisioner value 'pulum': must be one of 'pulumi', 'tofu', or 'terraform'
```

**Solution**: Check your provisioner label for typos. Valid values are: `pulumi`, `tofu`, or `terraform` (case-insensitive).

---

### "Provisioner not specified in manifest"

**Behavior**: CLI prompts you to select a provisioner.

**Solutions**:
1. **Add the label** to your manifest (recommended):
   ```yaml
   metadata:
     labels:
       project-planton.org/provisioner: pulumi
   ```

2. **Select interactively**: Type your choice when prompted

3. **Use provisioner-specific commands** if you prefer:
   ```bash
   project-planton pulumi up -f app.yaml
   ```

---

### Pulumi Backend Not Configured

**Error** (when using Pulumi provisioner):
```
error: no Pulumi backend configured
```

**Solution**: Set up Pulumi backend:

```bash
# Local backend (for testing)
pulumi login --local

# Or cloud backend
pulumi login
```

---

### Tofu/Terraform Backend Not Configured

**Behavior**: Uses local backend by default.

**Solution**: Configure backend via labels using provisioner-specific prefixes:

**For Terraform (S3):**
```yaml
metadata:
  labels:
    project-planton.org/provisioner: terraform
    terraform.project-planton.org/backend.type: s3
    terraform.project-planton.org/backend.bucket: my-terraform-state
    terraform.project-planton.org/backend.key: path/to/state.tfstate
    terraform.project-planton.org/backend.region: us-west-2
```

**For OpenTofu (GCS):**
```yaml
metadata:
  labels:
    project-planton.org/provisioner: tofu
    tofu.project-planton.org/backend.type: gcs
    tofu.project-planton.org/backend.bucket: my-tofu-state
    tofu.project-planton.org/backend.key: path/to/state
```

For complete backend configuration options, see the [State Backends Guide](/docs/guides/state-backends).

---

### Backend Configuration Changed

**Error**: Terraform/Tofu prompts for backend migration or reconfiguration.

**Solution**: Use the `--reconfigure` flag:

```bash
project-planton init -f manifest.yaml --reconfigure

# Also works with apply, destroy, plan, refresh
project-planton apply -f manifest.yaml --reconfigure
```

---

## Best Practices

### 1. Always Add Provisioner Label

Include the provisioner label in all manifests:

```yaml
metadata:
  labels:
    project-planton.org/provisioner: pulumi
```

This enables fully automated workflows.

---

### 2. Follow the Complete Lifecycle

Use all commands for a robust workflow:

```bash
# 1. Initialize (first time or after cleaning cache)
project-planton init -f app.yaml

# 2. Always preview before applying
project-planton plan -f app.yaml

# 3. Apply changes
project-planton apply -f app.yaml --yes

# 4. Refresh after manual changes
project-planton refresh -f app.yaml

# 5. Clean up when done
project-planton destroy -f app.yaml
```

---

### 3. Use -f Flag for Consistency

Prefer `-f` over `-f` to match kubectl style:

```bash
# Good (kubectl-style)
project-planton plan -f app.yaml
project-planton apply -f app.yaml

# Also works
project-planton apply -f app.yaml
```

---

### 4. Always Preview in CI/CD

Add preview step to pull request pipelines:

```bash
# .gitlab-ci.yml or similar
preview:
  script:
    - project-planton plan -f deployment.yaml
  only:
    - merge_requests

deploy:
  script:
    - project-planton apply -f deployment.yaml --yes
  only:
    - main
```

---

### 5. Combine with Validation

Always validate before applying in production:

```bash
project-planton validate -f app.yaml && \
project-planton plan -f app.yaml && \
project-planton apply -f app.yaml --yes
```

---

### 6. Use Refresh After Manual Changes

If you make manual changes via console or CLI, refresh state:

```bash
# Made manual changes to database in cloud console
project-planton refresh -f database.yaml

# Now apply can work with accurate state
project-planton apply -f database.yaml
```

---

### 7. Use Kustomize for Multi-Environment

Organize environments with kustomize:

```bash
project-planton init --kustomize-dir services/api --overlay prod
project-planton plan --kustomize-dir services/api --overlay prod
project-planton apply --kustomize-dir services/api --overlay prod
```

---

### 8. Override Values for CI/CD

Use `--set` for dynamic values in pipelines:

```bash
project-planton apply -f app.yaml \
    --set spec.image.tag="$CI_COMMIT_SHA" \
    --set metadata.labels.build="$CI_BUILD_ID"
```

---

### 9. Use Aliases for Familiarity

Use command aliases that match your background:

```bash
# Pulumi users might prefer
project-planton preview -f app.yaml

# Terraform/Tofu users might prefer
project-planton plan -f app.yaml

# kubectl users might prefer
project-planton delete -f app.yaml
```

---

## Related Documentation

- [CLI Reference](/docs/cli/cli-reference) - Complete command reference
- [Pulumi Commands](/docs/cli/pulumi-commands) - Pulumi-specific details
- [OpenTofu Commands](/docs/cli/tofu-commands) - OpenTofu-specific details
- [Manifest Structure](/docs/guides/manifests) - Writing manifests
- [Kustomize Integration](/docs/guides/kustomize) - Multi-environment setup
- [State Backends](/docs/guides/state-backends) - Configure state storage

---

## Feedback

Found an issue or have a suggestion? [Open an issue](https://github.com/plantonhq/project-planton/issues) on GitHub.

