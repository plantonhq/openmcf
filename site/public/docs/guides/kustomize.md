---
title: "Kustomize Integration"
description: "Using Kustomize with OpenMCF for multi-environment deployments - directory structure, overlays, and workflows"
icon: "layers"
order: 4
---

# Kustomize Integration Guide

Learn how to use Kustomize with OpenMCF for managing multi-environment deployments.

---

## What is Kustomize?

Kustomize is a configuration management tool that lets you create variations of YAML files without duplication. Instead of maintaining separate manifests for dev/staging/prod, you maintain one **base** and environment-specific **overlays** that patch the base.

**The Problem Without Kustomize**:

```
manifests/
├── dev-database.yaml      # Lots of duplication
├── staging-database.yaml  # Same resource, different values
└── prod-database.yaml     # Hard to maintain consistency
```

**The Solution With Kustomize**:

```
manifests/database/
├── base/
│   └── database.yaml      # Shared configuration
└── overlays/
    ├── dev/
    ├── staging/
    └── prod/               # Environment-specific patches
```

### The Clothing Analogy

Think of Kustomize like a clothing store:

- **Base** = The standard shirt design (common to all)
- **Overlays** = Size-specific modifications (small, medium, large)
- **Result** = A shirt that fits each person, derived from the same base design

OpenMCF integrates Kustomize seamlessly, building your manifest at deployment time.

---

## Quick Start

### 1. Create Base Manifest

```bash
mkdir -p services/api/kustomize/base
```

**`services/api/kustomize/base/deployment.yaml`**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: api
spec:
  container:
    image:
      repo: myapp/api
      tag: latest
    replicas: 1
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
```

**`services/api/kustomize/base/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml
```

### 2. Create Environment Overlay

**`services/api/kustomize/overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

patches:
  - path: patch.yaml
```

**`services/api/kustomize/overlays/prod/patch.yaml`**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: api
spec:
  container:
    image:
      tag: v1.0.0
    replicas: 3
    resources:
      limits:
        cpu: 2000m
        memory: 4Gi
```

### 3. Deploy with OpenMCF

```bash
# Deploy to production
openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay prod

# Deploy to dev
openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay dev
```

**What happens**:
1. OpenMCF runs `kustomize build services/api/kustomize/overlays/prod`
2. Merges base + prod overlay into final manifest
3. Validates the result
4. Deploys using Pulumi or OpenTofu

---

## Directory Structure

### Standard Layout

```
<service-name>/kustomize/
├── base/
│   ├── kustomization.yaml          # Base kustomization config
│   └── <resource>.yaml             # Base resource definition
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml      # Dev environment config
    │   └── patch.yaml              # Dev-specific patches
    ├── staging/
    │   ├── kustomization.yaml
    │   └── patch.yaml
    └── prod/
        ├── kustomization.yaml
        └── patch.yaml
```

### Example: Complete Service

```
backend/services/api/kustomize/
├── base/
│   ├── kustomization.yaml
│   ├── deployment.yaml
│   └── database.yaml
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml
    │   ├── deployment-patch.yaml
    │   └── database-patch.yaml
    ├── staging/
    │   ├── kustomization.yaml
    │   ├── deployment-patch.yaml
    │   └── database-patch.yaml
    └── prod/
        ├── kustomization.yaml
        ├── deployment-patch.yaml
        └── database-patch.yaml
```

---

## Creating Patches

### Strategic Merge Patches

The most common approach - specify only the fields you want to change:

**Base**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
spec:
  container:
    replicas: 1
    resources:
      limits:
        cpu: 500m
        memory: 1Gi
    diskSize: 10Gi
```

**Prod Patch** (`overlays/prod/patch.yaml`):

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
spec:
  container:
    replicas: 3              # Override
    resources:
      limits:
        cpu: 2000m           # Override
        memory: 4Gi          # Override
    diskSize: 100Gi          # Override
```

**Result**: Base + patch merged = 3 replicas, 2000m CPU, 4Gi memory, 100Gi disk.

### JSON 6902 Patches

For more complex changes:

**`overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

patches:
  - target:
      kind: KubernetesPostgres
      name: app-database
    patch: |-
      - op: replace
        path: /spec/container/replicas
        value: 3
      - op: add
        path: /metadata/labels/environment
        value: production
```

---

## Common Patterns

### Pattern 1: Environment-Specific Resources

Different instance sizes per environment:

**Dev** (small, cheap):

```yaml
spec:
  container:
    replicas: 1
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
```

**Prod** (large, resilient):

```yaml
spec:
  container:
    replicas: 5
    resources:
      limits:
        cpu: 2000m
        memory: 4Gi
```

### Pattern 2: Environment-Specific Labels

Add labels for cost tracking:

**`overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

commonLabels:
  environment: production
  cost-center: engineering
  team: backend
```

### Pattern 3: Environment-Specific Images

**`overlays/dev/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

images:
  - name: myapp/api
    newTag: latest          # Dev uses latest

patches:
  - path: dev-patch.yaml
```

**`overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

images:
  - name: myapp/api
    newTag: v1.2.3          # Prod uses specific version

patches:
  - path: prod-patch.yaml
```

### Pattern 4: Shared Configuration + Environment Patches

**`base/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: api
  managed-by: openmcf

resources:
  - deployment.yaml
  - database.yaml
  - cache.yaml
```

Each resource in base defines shared configuration, overlays patch for environment needs.

---

## Workflow Examples

### Deploying to Multiple Environments

```bash
# Deploy to dev
openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay dev \
  --yes

# Test in dev...

# Deploy to staging
openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay staging \
  --yes

# Test in staging...

# Deploy to production (with review)
openmcf pulumi preview \
  --kustomize-dir services/api/kustomize \
  --overlay prod

# Review changes...

openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay prod
```

### Combining Kustomize with --set Overrides

```bash
# Kustomize overlay + runtime override
openmcf pulumi up \
  --kustomize-dir services/api/kustomize \
  --overlay prod \
  --set spec.container.image.tag=v1.2.4
```

**Order of precedence**:
1. Base manifest
2. Overlay patches applied
3. `--set` overrides applied last (highest priority)

### Preview Built Manifest

```bash
# See what Kustomize generates (useful for debugging)
cd services/api/kustomize
kustomize build overlays/prod

# Or let OpenMCF build and show it
openmcf load-manifest \
  --kustomize-dir services/api/kustomize \
  --overlay prod
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Deploy API

on:
  push:
    branches:
      - main
      - develop

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Determine Environment
        id: env
        run: |
          if [ "${{ github.ref }}" == "refs/heads/main" ]; then
            echo "overlay=prod" >> $GITHUB_OUTPUT
          else
            echo "overlay=dev" >> $GITHUB_OUTPUT
          fi
      
      - name: Deploy
        run: |
          openmcf pulumi up \
            --kustomize-dir services/api/kustomize \
            --overlay ${{ steps.env.outputs.overlay }} \
            --yes
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

### GitLab CI

```yaml
deploy:
  script:
    - |
      if [ "$CI_COMMIT_BRANCH" == "main" ]; then
        OVERLAY="prod"
      elif [ "$CI_COMMIT_BRANCH" == "staging" ]; then
        OVERLAY="staging"
      else
        OVERLAY="dev"
      fi
    - openmcf pulumi up --kustomize-dir services/api/kustomize --overlay $OVERLAY --yes
  only:
    - main
    - staging
    - develop
```

---

## Advanced Techniques

### Multiple Bases

Useful for shared components:

```
common/
└── base/
    ├── kustomization.yaml
    └── shared-config.yaml

service-a/kustomize/
└── overlays/
    └── prod/
        ├── kustomization.yaml  # References ../../../common/base
        └── patch.yaml

service-b/kustomize/
└── overlays/
    └── prod/
        ├── kustomization.yaml  # Also references ../../../common/base
        └── patch.yaml
```

### Components (Reusable Pieces)

**`components/monitoring/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1alpha1
kind: Component

patches:
  - path: add-monitoring.yaml
```

**Use in overlay**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

components:
  - ../../../components/monitoring
```

### Generating ConfigMaps from Files

**`overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

configMapGenerator:
  - name: app-config
    files:
      - config/prod.yaml
      - config/secrets.encrypted
```

---

## Troubleshooting

### Error: "no such file or directory"

**Problem**: Kustomize can't find referenced files.

**Solution**:
```bash
# Check file paths in kustomization.yaml
# Ensure paths are relative to kustomization.yaml location

# Verify structure
ls -R services/api/kustomize/
```

### Error: "kustomization.yaml not found"

**Problem**: Missing kustomization.yaml in overlay.

**Solution**:
```bash
# Create kustomization.yaml
cat > services/api/kustomize/overlays/prod/kustomization.yaml <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base
EOF
```

### Patch Not Applied

**Problem**: Your patch isn't affecting the final output.

**Solution**:
```bash
# 1. Verify patch is listed in kustomization.yaml
cat overlays/prod/kustomization.yaml

# 2. Check patch targets correct resource
# - apiVersion must match
# - kind must match
# - metadata.name must match

# 3. Test kustomize build directly
cd services/api/kustomize
kustomize build overlays/prod
```

### Wrong Overlay Applied

**Problem**: Deployed dev config to prod (or vice versa).

**Solution**:
```bash
# Always verify overlay before deploying
openmcf pulumi preview \
  --kustomize-dir services/api/kustomize \
  --overlay prod  # Double-check this!

# Use explicit confirmation in CI/CD
if [ "$OVERLAY" != "prod" ]; then
  echo "Deploying to $OVERLAY"
  openmcf pulumi up --kustomize-dir ... --overlay $OVERLAY --yes
else
  echo "Production deployment - manual approval required"
  openmcf pulumi up --kustomize-dir ... --overlay $OVERLAY
fi
```

---

## Best Practices

### 1. **Keep Base Minimal**

```yaml
# ✅ Good: Base has common configuration
base/deployment.yaml:
  name: api
  container:
    image:
      repo: myapp/api
    # No environment-specific values

# ❌ Bad: Base has production values
base/deployment.yaml:
  name: api
  container:
    replicas: 10  # Production-specific, doesn't belong in base
```

### 2. **One Overlay Per Environment**

```
# ✅ Good
overlays/
├── dev/
├── staging/
└── prod/

# ❌ Confusing
overlays/
├── dev-us-west/
├── dev-eu-central/
├── staging-us-west/
└── ... (too many combinations)
```

### 3. **Use Descriptive Patch Names**

```
# ✅ Good
overlays/prod/
├── kustomization.yaml
├── resources-patch.yaml          # Increases resources
├── replicas-patch.yaml            # Scales replicas
└── monitoring-patch.yaml          # Adds monitoring

# ❌ Bad
overlays/prod/
├── kustomization.yaml
├── patch1.yaml
├── patch2.yaml
└── patch3.yaml
```

### 4. **Version Control Everything**

```bash
# ✅ Good: All kustomize files in Git
git add services/api/kustomize/
git commit -m "feat: add production overlay for API"

# ❌ Bad: Generated files, temp files committed
git add services/api/kustomize/overlays/prod/output.yaml  # Generated
```

### 5. **Test Overlays Before Deploying**

```bash
# ✅ Good: Preview before applying
kustomize build overlays/prod | less  # Review output
openmcf pulumi preview --kustomize-dir ... --overlay prod

# ⚠️ Risky: Deploy without review
openmcf pulumi up --kustomize-dir ... --overlay prod --yes
```

---

## Complete Example

Here's a complete real-world example:

### Directory Structure

```
backend/services/api/kustomize/
├── base/
│   ├── kustomization.yaml
│   ├── deployment.yaml
│   ├── database.yaml
│   └── redis.yaml
└── overlays/
    ├── dev/
    │   ├── kustomization.yaml
    │   ├── deployment-patch.yaml
    │   ├── database-patch.yaml
    │   └── redis-patch.yaml
    └── prod/
        ├── kustomization.yaml
        ├── deployment-patch.yaml
        ├── database-patch.yaml
        └── redis-patch.yaml
```

### Files

**`base/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: api
  managed-by: openmcf

resources:
  - deployment.yaml
  - database.yaml
  - redis.yaml
```

**`base/deployment.yaml`**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: api
spec:
  container:
    image:
      repo: mycompany/api
      tag: latest
    replicas: 1
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
```

**`overlays/prod/kustomization.yaml`**:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

commonLabels:
  environment: production

images:
  - name: mycompany/api
    newTag: v1.0.0

patches:
  - path: deployment-patch.yaml
  - path: database-patch.yaml
  - path: redis-patch.yaml
```

**`overlays/prod/deployment-patch.yaml`**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: api
spec:
  container:
    replicas: 5
    resources:
      limits:
        cpu: 2000m
        memory: 4Gi
```

### Deployment

```bash
# Deploy to production
openmcf pulumi up \
  --kustomize-dir backend/services/api/kustomize \
  --overlay prod
```

---

## Related Documentation

- [Manifest Structure](/docs/guides/manifests) - Understanding manifests
- [Pulumi Commands](/docs/cli/pulumi-commands) - Deploying with Pulumi
- [OpenTofu Commands](/docs/cli/tofu-commands) - Deploying with OpenTofu
- [Official Kustomize Docs](https://kustomize.io/) - Learn more about Kustomize

---

## Next Steps

1. **Create Your First Overlay**: Start with a simple dev/prod split
2. **Test Locally**: Use `kustomize build` to preview results
3. **Deploy Gradually**: Test in dev before prod
4. **Iterate**: Add more overlays as needed (staging, regional, etc.)

Kustomize + OpenMCF gives you powerful multi-environment management with minimal duplication. Start simple and grow as needed.

