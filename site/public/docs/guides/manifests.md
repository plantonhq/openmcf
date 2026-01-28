---
title: "Manifest Structure Guide"
description: "Understanding OpenMCF manifests - KRM structure, validation, defaults, and best practices for writing infrastructure definitions"
icon: "document"
order: 1
---

# Manifest Structure Guide

Your complete guide to writing and understanding OpenMCF manifests.

---

## What is a Manifest?

A manifest is a YAML file that describes a piece of infrastructure you want to deploy. Think of it as a recipe card: it lists what you want to create (a database, a Kubernetes cluster, a storage bucket) and how you want it configured (size, region, settings).

**Simple example**:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareR2Bucket
metadata:
  name: my-app-assets
spec:
  accountId: abc123
  location: WNAM
```

This manifest says: "Create a Cloudflare R2 bucket named `my-app-assets` in the Western North America location."

### The Restaurant Menu Analogy

Think of manifests like ordering from a restaurant:

- **apiVersion** = Which menu you're ordering from (Cloudflare menu, AWS menu, GCP menu)
- **kind** = What dish you're ordering (R2 Bucket, S3 Bucket, GCS Bucket)
- **metadata** = Your order details (name, table number, special instructions)
- **spec** = How you want it prepared (size, toppings, customizations)
- **status** = The kitchen's feedback (order ready, here's your table number, etc.)

---

## Anatomy of a Manifest

Every OpenMCF manifest follows the **Kubernetes Resource Model (KRM)** structure. This isn't an accident—it's the same pattern used by Kubernetes, which millions of developers already know.

### The Five Sections

```yaml
apiVersion: <provider>.openmcf.org/<version>
kind: <ResourceType>
metadata:
  name: <resource-name>
  # ... more metadata
spec:
  # ... resource-specific configuration
status:
  # ... read-only status (populated after deployment)
```

Let's break down each section.

---

## apiVersion: The Menu Selection

**Format**: `<provider>.openmcf.org/<version>`

**Purpose**: Identifies which cloud provider and API version you're using.

**Examples**:
- `aws.openmcf.org/v1` - AWS resources
- `gcp.openmcf.org/v1` - GCP resources
- `azure.openmcf.org/v1` - Azure resources
- `cloudflare.openmcf.org/v1` - Cloudflare resources
- `kubernetes.openmcf.org/v1` - Kubernetes resources

**Why it matters**: The `apiVersion` tells OpenMCF which API definitions to use for validation. As APIs evolve, version numbers allow you to opt into changes gradually.

**Versioning**:
- `v1` = First stable version
- `v1alpha1`, `v1beta1` = Pre-stable versions (use with caution in production)
- `v2` = Second major version (indicates breaking changes from v1)

---

## kind: The Dish You're Ordering

**Format**: PascalCase string (e.g., `AwsS3Bucket`, `GcpGkeCluster`, `RedisKubernetes`)

**Purpose**: Specifies exactly which type of infrastructure resource you want to deploy.

**Naming convention**: Usually combines provider prefix with resource type:
- AWS: `AwsS3Bucket`, `AwsEksCluster`, `AwsRdsInstance`
- GCP: `GcpGcsBucket`, `GcpGkeCluster`, `GcpCloudSql`
- Kubernetes: `PostgresKubernetes`, `RedisKubernetes`, `KafkaKubernetes`

**Finding available kinds**:
- Browse the [Catalog](/docs/catalog) - all 118 deployment components
- Check [Buf Schema Registry](https://buf.build/openmcf/apis) - API documentation
- Search the repository: `provider/<provider>/`

**Example kinds**:

```yaml
# Deploy Redis to Kubernetes
kind: RedisKubernetes

# Deploy Postgres to AWS RDS
kind: AwsRdsInstance

# Deploy an application to GCP Cloud Run
kind: GcpCloudRun
```

---

## metadata: The Shipping Label

The `metadata` section contains identifying information and administrative details about your resource.

### Required Fields

**name** (string): Unique identifier for this resource.

```yaml
metadata:
  name: production-database
```

**Naming rules**:
- Lowercase alphanumeric with hyphens
- Start and end with alphanumeric
- No underscores or special characters
- Max 63 characters
- Must be unique within your organization/environment

### Optional Fields

**labels** (key-value pairs): Arbitrary tags for organization and querying.

```yaml
metadata:
  name: api-deployment
  labels:
    environment: production
    team: backend
    cost-center: engineering
    pulumi.openmcf.org/stack.name: "acme/platform/prod.ApiDeployment.api-v2"
```

**Common labels**:
- `environment`: dev, staging, prod
- `team`: owning team name
- `project`: project identifier
- `version`: version tag
- `pulumi.openmcf.org/stack.name`: Pulumi stack FQDN (when using Pulumi)

**Why labels matter**: They help you:
- Organize resources
- Track costs by team/project
- Query resources programmatically
- Configure backend state management

---

## spec: The Customization Options

The `spec` (specification) section contains the resource-specific configuration. **This is the most important part of your manifest** — it defines exactly how your infrastructure should be configured.

### Structure is Resource-Specific

Every `kind` has its own `spec` structure defined in Protocol Buffers. There's no universal spec—what goes in an `AwsS3Bucket` spec is completely different from a `GcpGkeCluster` spec.

### Example: Cloudflare R2 Bucket

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareR2Bucket
metadata:
  name: pipeline-logs
spec:
  accountId: abc123xyz
  location: WNAM
```

### Example: PostgreSQL on Kubernetes

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
spec:
  container:
    replicas: 3
    resources:
      limits:
        cpu: 2000m
        memory: 4Gi
      requests:
        cpu: 1000m
        memory: 2Gi
    diskSize: 100Gi
    isPersistenceEnabled: true
```

### Finding Spec Fields

**Method 1**: Browse component documentation in the [Catalog](/docs/catalog)

**Method 2**: Check the protobuf definition in the repository:
- Location: `provider/<provider>/<component>/v1/spec.proto`
- Shows all available fields, types, and validation rules

**Method 3**: Use `buf.build` Schema Registry (coming soon)

### Required vs Optional Fields

Fields can be:
- **Required**: Must be specified (validation fails if missing)
- **Optional**: Can be omitted
- **Optional with defaults**: Can be omitted (gets a default value automatically)

Check the proto definition or documentation to see which fields are required.

---

## status: The Read-Only Feedback

The `status` section contains outputs and state information populated **after deployment**. You never write this section manually—it's filled in by the deployment process.

### What Goes in Status

**Deployment outputs**: Information you need after resources are created:

```yaml
status:
  outputs:
    connectionString: "postgres://prod-db.us-west-2.rds.amazonaws.com:5432"
    databaseName: "myapp"
    endpoint: "prod-db.us-west-2.rds.amazonaws.com"
```

**Kubernetes status** (for Kubernetes resources):

```yaml
status:
  kubernetes:
    podStatus: Running
    replicas: 3
    readyReplicas: 3
```

**Common output types**:
- Connection strings / endpoints
- Resource IDs
- DNS names
- API keys (encrypted or references)
- Status indicators

### Why Status is Separate

Keeping status separate from spec follows the Kubernetes philosophy:
- **spec** = Desired state (what you want)
- **status** = Observed state (what actually exists)
- **Separation** = Makes it easy to detect drift

---

## Validation: Early Error Detection

One of OpenMCF's superpowers is **validation before deployment**. Instead of waiting 5 minutes for an AWS deployment to fail because you typo'd a field name, you catch errors in seconds.

### Three Validation Layers

**1. Schema Validation** (via Protocol Buffers)

```bash
# This fails immediately if YAML doesn't match proto structure
openmcf validate -f my-resource.yaml
```

**2. Field-Level Validation** (via buf-validate)

Protocol buffer definitions include validation rules:

```protobuf
message PostgresKubernetesSpec {
  int32 replicas = 1 [(buf.validate.field).int32 = {gte: 1, lte: 10}];
  string cpu = 2 [(buf.validate.field).string.pattern = "^[0-9]+m$"];
}
```

This catches errors like:
- Invalid replica count (must be 1-10)
- Malformed CPU value (must be like "500m")
- Missing required fields

**3. Provider Validation** (during deployment)

Final validation by the actual cloud provider APIs. This catches provider-specific constraints.

### Validating Your Manifest

```bash
# Validate before deploying
openmcf validate -f my-app.yaml

# If validation fails, you'll see exactly what's wrong
❌  MANIFEST VALIDATION FAILED

⚠️  Validation Errors:

spec.replicas: value must be between 1 and 10 (got: 15)
spec.container.cpu: value must match pattern "^[0-9]+m$" (got: "invalid")
```

### Why Validation Matters

**Without validation**:
1. Edit manifest
2. Run `pulumi up` (or `tofu apply`)
3. Wait 5-10 minutes
4. Deployment fails
5. Fix manifest
6. Repeat

**With validation**:
1. Edit manifest
2. Run `validate` (2 seconds)
3. See errors immediately
4. Fix manifest
5. Deploy with confidence

---

## Default Values: Keeping Manifests Concise

Many fields have sensible defaults so you don't have to specify everything.

### How Defaults Work

Defaults are defined in the protobuf with the `(org.openmcf.shared.options.default)` extension:

```protobuf
message KubernetesExternalDnsKubernetesSpec {
  optional string namespace = 1 [(org.openmcf.shared.options.default) = "kubernetes-external-dns"];
  optional string version = 2 [(org.openmcf.shared.options.default) = "v0.19.0"];
}
```

### Minimal vs Explicit Manifests

**Minimal (using defaults)**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDnsKubernetes
metadata:
  name: kubernetes-external-dns
spec:
  targetCluster:
    kubernetesProviderConfigId: my-cluster
  # namespace and version get defaults automatically
```

**Explicit (overriding defaults)**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalDnsKubernetes
metadata:
  name: kubernetes-external-dns
spec:
  namespace: custom-dns-namespace
  version: v0.20.0
  targetCluster:
    kubernetesProviderConfigId: my-cluster
```

### Viewing Applied Defaults

Use `load-manifest` to see the effective configuration with defaults:

```bash
openmcf load-manifest kubernetes-external-dns.yaml

# Output shows defaults filled in:
# spec:
#   namespace: kubernetes-external-dns        # ← Default applied
#   version: v0.19.0               # ← Default applied
#   targetCluster:
#     kubernetesProviderConfigId: my-cluster
```

---

## Complete Example: AWS S3 Bucket

Let's walk through a complete manifest with annotations:

```yaml
# Which API (AWS resources, version 1)
apiVersion: aws.openmcf.org/v1

# What resource (S3 Bucket)
kind: AwsS3Bucket

# Identifying information
metadata:
  name: user-uploads-prod
  labels:
    environment: production
    team: backend
    purpose: user-uploads
    pulumi.openmcf.org/stack.name: "acme/storage/prod.AwsS3Bucket.user-uploads"

# Configuration
spec:
  # AWS account region
  region: us-west-2
  
  # Bucket configuration
  versioning: true
  
  # Encryption settings
  encryption:
    enabled: true
    kmsKeyId: arn:aws:kms:us-west-2:123456789:key/abc-123
  
  # Lifecycle rules
  lifecycleRules:
    - id: delete-old-versions
      enabled: true
      expirationDays: 90
  
  # Access control
  blockPublicAccess: true
  
  # Tags
  tags:
    CostCenter: "engineering"
    DataClassification: "internal"

# Status will be populated after deployment
# (never write this manually)
```

---

## Multi-Resource Example: PostgreSQL with Backups

Sometimes you need multiple manifests for related resources:

**postgres.yaml**:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
spec:
  container:
    replicas: 3
    resources:
      limits:
        cpu: 2000m
        memory: 4Gi
    diskSize: 100Gi
```

**postgres-backup-bucket.yaml**:

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsS3Bucket
metadata:
  name: postgres-backups
spec:
  region: us-west-2
  versioning: true
  lifecycleRules:
    - id: delete-old-backups
      enabled: true
      expirationDays: 30
```

Deploy separately but manage together:

```bash
# Deploy database
openmcf pulumi up -f postgres.yaml

# Deploy backup bucket
openmcf pulumi up -f postgres-backup-bucket.yaml
```

---

## Loading Manifests from URLs

Manifests don't have to be local files—you can load them from URLs:

```bash
# Load from GitHub raw URL
openmcf pulumi up \
  -f https://raw.githubusercontent.com/my-org/manifests/main/prod/database.yaml

# Load from any HTTPS URL
openmcf pulumi up \
  -f https://config-server.example.com/manifests/vpc.yaml
```

**How it works**:
1. CLI downloads the manifest to a temporary file
2. Validates it
3. Deploys it
4. Cleans up temporary file

**Use cases**:
- Centralized manifest repository
- Generated manifests from CI/CD
- Shared manifests across teams
- Version-controlled remote manifests

---

## Best Practices

### 1. **Use Version Control**

```bash
# ✅ Good: Track manifests in Git
git add ops/manifests/prod-database.yaml
git commit -m "feat: increase database resources"
git push

# ❌ Bad: Ad-hoc manifests
vim /tmp/database.yaml
```

**Why**: Version control gives you history, rollback, code review, and auditability.

### 2. **Validate Before Deploying**

```bash
# ✅ Good: Catch errors early
openmcf validate -f resource.yaml
openmcf pulumi up -f resource.yaml

# ⚠️ Risky: Deploy without validation
openmcf pulumi up -f resource.yaml
```

**Why**: Validation catches 90% of errors before making cloud API calls.

### 3. **Use Meaningful Names**

```yaml
# ✅ Good: Clear, descriptive names
metadata:
  name: prod-api-postgres-primary

# ❌ Bad: Generic names
metadata:
  name: db1
```

**Why**: Good names make it obvious what the resource does.

### 4. **Organize by Environment**

```
manifests/
├── dev/
│   ├── database.yaml
│   └── cache.yaml
├── staging/
│   ├── database.yaml
│   └── cache.yaml
└── prod/
    ├── database.yaml
    └── cache.yaml
```

**Why**: Clear separation prevents accidents (deploying to wrong environment).

### 5. **Document Non-Obvious Choices**

```yaml
spec:
  # Using m5.2xlarge because m5.xlarge caused OOM errors under load
  # See: https://jira.company.com/PROJECT-123
  instanceType: m5.2xlarge
```

**Why**: Comments explain "why" when it's not obvious from the "what".

### 6. **Use Labels for Organization**

```yaml
metadata:
  labels:
    team: payments
    cost-center: "12345"
    environment: production
    data-classification: pii
```

**Why**: Labels enable querying, cost tracking, and organization.

### 7. **Keep Secrets Out of Manifests**

```yaml
# ❌ Bad: Secrets in manifest
spec:
  apiKey: "sk_live_abc123xyz"

# ✅ Good: Reference to secret
spec:
  apiKeySecretRef:
    name: payment-processor-key
    key: api-key
```

**Why**: Manifests are often committed to Git. Secrets should be in secret managers.

### 8. **One Resource Per Manifest (Usually)**

```bash
# ✅ Good: Separate files for separate resources
ops/
├── database.yaml
├── cache.yaml
└── storage-bucket.yaml

# ⚠️ Sometimes OK: Related resources in one file
ops/
└── database-with-backup-bucket.yaml
```

**Why**: Separate files make it easier to deploy, version, and manage resources independently.

---

## Troubleshooting

### "kind not supported" Error

**Problem**: CLI doesn't recognize your `kind`.

**Solution**:
- Check spelling (case-sensitive, PascalCase)
- Verify kind exists in catalog: `/docs/catalog`
- Ensure you're using the right `apiVersion`

### Validation Fails with Cryptic Error

**Problem**: Validation error message isn't clear.

**Solution**:
```bash
# 1. Check YAML syntax
cat manifest.yaml | yq .

# 2. Verify required fields
# - Check proto definition or docs for required fields

# 3. Look for typos in field names
# - Proto field names use snake_case: disk_size, not diskSize
```

### Defaults Not Applied

**Problem**: Expected defaults aren't showing up.

**Solution**:
```bash
# Defaults are only applied when field is omitted
# Check if you accidentally set field to empty string or 0

# View effective manifest with defaults:
openmcf load-manifest resource.yaml
```

### Manifest from URL Fails

**Problem**: Can't load manifest from URL.

**Solutions**:
- Check URL is publicly accessible
- Verify HTTPS (not HTTP)
- Ensure URL returns raw YAML (not HTML page)
- For GitHub: Use "Raw" button to get correct URL

---

## Related Documentation

- [Pulumi Commands](/docs/cli/pulumi-commands) - Deploying with Pulumi
- [OpenTofu Commands](/docs/cli/tofu-commands) - Deploying with OpenTofu
- [Credentials Guide](/docs/guides/credentials) - Setting up provider credentials
- [Advanced Usage](/docs/guides/advanced-usage) - Using --set, URL manifests, and more
- [Deployment Component Catalog](/docs/catalog) - Browse all available kinds

---

## Next Steps

Now that you understand manifests:

1. **Browse the Catalog**: See what resources you can deploy at `/docs/catalog`
2. **Write Your First Manifest**: Start with something simple like a storage bucket
3. **Validate It**: Use `openmcf validate -f your-resource.yaml`
4. **Deploy It**: Follow the Pulumi or OpenTofu command guides

**Remember**: Manifests are declarative—you describe what you want, not how to create it. OpenMCF handles the "how."

