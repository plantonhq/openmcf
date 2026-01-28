# KubernetesExternalSecrets Examples

Complete YAML manifests for deploying External Secrets Operator (ESO) on different cloud providers.

---

## Example 1: Minimal GKE Setup with Google Cloud Secret Manager

Basic deployment for GKE using Workload Identity.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-gke-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 30
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  gke:
    project_id:
      value: my-gcp-project
    gsa_email: external-secrets@my-gcp-project.iam.gserviceaccount.com
```

**How it works:**
- Deploys ESO with Workload Identity for Google Cloud Secret Manager
- Polls for secret changes every 30 seconds
- Uses the specified GCP service account via Workload Identity
- Automatically creates ClusterSecretStore for GCP Secret Manager

**Prerequisites:**
- GKE cluster with Workload Identity enabled
- Google Cloud Secret Manager API enabled
- GCP service account with `secretmanager.secretAccessor` role
- IAM binding: `gcloud iam service-accounts add-iam-policy-binding`

**What gets created:**
- Namespace: `external-secrets` (default)
- Helm release: External Secrets Operator
- ServiceAccount with Workload Identity annotation
- ClusterSecretStore for GCP Secret Manager

---

## Example 2: GKE with Foreign Key References

Reference GCP project and service account from other OpenMCF resources.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-gke-ref
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 60
  
  container:
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
      requests:
        cpu: 25m
        memory: 64Mi
  
  gke:
    project_id:
      ref: gcp-project-platform
    gsa_email: external-secrets@my-project.iam.gserviceaccount.com
```

**Benefits:**
- GCP project ID pulled from existing `GcpProject` resource
- Single source of truth for resource IDs
- Safer than hardcoding values
- Easier to manage across environments

---

## Example 3: Using Existing Namespace

Deploy External Secrets to an existing namespace managed separately.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-existing-ns
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: platform-services  # Must already exist
  create_namespace: false  # Use existing namespace
  
  poll_interval_seconds: 30
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  gke:
    project_id:
      value: my-gcp-project
    gsa_email: external-secrets@my-gcp-project.iam.gserviceaccount.com
```

**Use case:**
- Namespace `platform-services` is managed by platform team
- Multiple platform components share the same namespace
- Namespace has specific labels, quotas, or network policies
- GitOps workflow where namespace is created separately

**Prerequisites:**
- Namespace must exist before deploying this component
- You have permissions to create resources in that namespace

**Important:** If the namespace doesn't exist, deployment will fail with "namespace not found" error.

---

## Example 4: Shared Namespace with Multiple Components

Deploy External Secrets alongside other platform components in a shared namespace.

```yaml
# First, create the namespace (via separate manifest or KubernetesNamespace component)
apiVersion: v1
kind: Namespace
metadata:
  name: platform-infrastructure
  labels:
    managed-by: platform-team
    environment: production

---
# Deploy External Secrets to the shared namespace
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: eso-shared-platform
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: platform-infrastructure
  create_namespace: false  # Namespace already created above
  
  poll_interval_seconds: 60
  
  container:
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
      requests:
        cpu: 25m
        memory: 64Mi
  
  gke:
    project_id:
      value: my-gcp-project
    gsa_email: external-secrets@my-gcp-project.iam.gserviceaccount.com
```

**Benefits:**
- Consolidate multiple platform components in one namespace
- Easier resource management and monitoring
- Shared quotas and RBAC policies
- Reduced namespace sprawl

**Pattern:**
1. Create namespace with appropriate labels and policies
2. Deploy External Secrets with `create_namespace: false`
3. Deploy other platform components to the same namespace
4. All components share the namespace lifecycle

---

## Example 5: EKS with AWS Secrets Manager (Auto-Created IRSA Role)

Deploy ESO on Amazon EKS with automatic IRSA role creation.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-eks-prod
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 15
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  eks:
    region: us-east-1
```

**How it works:**
- Automatically creates an IAM role for IRSA
- Grants SecretsManager permissions (GetSecretValue, DescribeSecret)
- Annotates ESO ServiceAccount with the IAM role ARN
- No static AWS credentials needed

**Prerequisites:**
- EKS cluster with OIDC provider configured
- AWS Secrets Manager in the specified region

**IAM Policy (auto-created):**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": "*"
    }
  ]
}
```

---

## Example 6: EKS with Custom IRSA Role

Provide your own IAM role ARN for fine-grained secret access.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-eks-custom
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 30
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  eks:
    region: us-west-2
    irsa_role_arn_override: arn:aws:iam::123456789012:role/external-secrets-limited
```

**Use case:**
- When you manage IAM roles separately (e.g., via Terraform)
- For compliance requirements around IAM role naming
- To scope secrets to specific prefixes (`prod/*`, `app-name/*`)
- For cross-account secret access

**Example scoped IAM policy:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": [
        "arn:aws:secretsmanager:us-west-2:123456789012:secret:prod/*",
        "arn:aws:secretsmanager:us-west-2:123456789012:secret:shared/*"
      ]
    }
  ]
}
```

---

## Example 7: AKS with Azure Key Vault

Deploy ESO on Azure Kubernetes Service using Azure Workload Identity.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-aks-prod
spec:
  target_cluster:
    cluster_name: prod-aks-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 30
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  aks:
    key_vault_resource_id: /subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-rg/providers/Microsoft.KeyVault/vaults/prod-keyvault
    managed_identity_client_id: 87654321-4321-4321-4321-210987654321
```

**How it works:**
- Uses Azure Workload Identity (federated credentials)
- Managed Identity has "Key Vault Secrets User" role on the Key Vault
- No client secrets or passwords needed
- Annotates ESO ServiceAccount with the Managed Identity client ID

**Prerequisites:**
- AKS cluster with Azure AD Workload Identity enabled
- Azure Key Vault created
- User-assigned Managed Identity with Key Vault permissions
- Federated credential linking Managed Identity to Kubernetes ServiceAccount

**Azure CLI setup:**
```bash
# Create managed identity
az identity create --name external-secrets-identity --resource-group prod-rg

# Assign Key Vault Secrets User role
az role assignment create \
  --assignee <managed-identity-principal-id> \
  --role "Key Vault Secrets User" \
  --scope /subscriptions/.../vaults/prod-keyvault

# Create federated credential
az identity federated-credential create \
  --name eso-fed-cred \
  --identity-name external-secrets-identity \
  --resource-group prod-rg \
  --issuer https://oidc.prod-aks-cluster.azure.com \
  --subject system:serviceaccount:external-secrets:external-secrets
```

---

## Example 8: Cost-Optimized Configuration (Hourly Polling)

Balance secret freshness with cloud API costs using longer poll intervals.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-cost-optimized
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: external-secrets
  
  # Poll every hour to minimize API costs
  poll_interval_seconds: 3600
  
  container:
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
      requests:
        cpu: 25m
        memory: 64Mi
  
  gke:
    project_id:
      value: my-gcp-project
    gsa_email: external-secrets@my-gcp-project.iam.gserviceaccount.com
```

**Cost considerations:**
- **10-second poll** (100 secrets): ~864,000 API calls/day
- **1-hour poll** (100 secrets): ~2,400 API calls/day
- **4-hour poll** (100 secrets): ~600 API calls/day

**AWS Secrets Manager pricing:** $0.05 per 10,000 API calls
- 10-second poll: ~$4.32/day = ~$130/month
- 1-hour poll: ~$0.12/day = ~$3.60/month
- 4-hour poll: ~$0.03/day = ~$0.90/month

**When to use:**
- Configuration secrets that change infrequently (TLS certs, API keys)
- Large number of secrets (100+)
- Cost-sensitive environments

**When NOT to use:**
- Frequently rotating credentials (database passwords with 15-min TTL)
- Critical secrets requiring fast propagation

---

## Example 9: High-Frequency Polling for Dynamic Secrets

Fast polling for short-lived credentials.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-dynamic
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: external-secrets
  
  # Poll every 5 minutes for dynamic secrets
  poll_interval_seconds: 300
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 100m
        memory: 256Mi
  
  eks:
    region: us-east-1
```

**Use case:**
- Database credentials with automatic rotation (e.g., RDS secrets with auto-rotation)
- Short-lived tokens from Vault
- Secrets that change multiple times per day
- Applications that can reload secrets without restart (mounted as volumes)

**Note:** Even with fast polling, Kubernetes Secret updates have kubelet sync delays (default 60s for volume mounts).

---

## Example 10: Resource-Constrained Environment

Minimal resource footprint for dev/staging clusters.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-minimal
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 120
  
  container:
    resources:
      limits:
        cpu: 200m
        memory: 256Mi
      requests:
        cpu: 10m
        memory: 32Mi
  
  gke:
    project_id:
      value: dev-gcp-project
    gsa_email: external-secrets@dev-project.iam.gserviceaccount.com
```

**Resource usage:**
- Minimal CPU: 10m request, 200m limit
- Minimal memory: 32Mi request, 256Mi limit
- Longer poll interval (2 minutes) reduces overhead

**Good for:**
- Development clusters
- CI/CD clusters
- Low-traffic staging environments
- Clusters with limited node resources

---

## Example 11: Production High-Availability Setup

Production-ready configuration with optimal resources.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-ha-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 30
  
  container:
    resources:
      limits:
        cpu: 2000m
        memory: 2Gi
      requests:
        cpu: 100m
        memory: 256Mi
  
  gke:
    project_id:
      value: prod-gcp-project
    gsa_email: external-secrets@prod-project.iam.gserviceaccount.com
```

**What this enables:**
- Multiple controller replicas (configured via Helm chart)
- Leader election for active/standby controllers
- Handles large number of ExternalSecret resources (500+)
- Fast secret synchronization
- Resource headroom for traffic spikes

**Additional production considerations:**
- Configure pod anti-affinity to spread replicas across nodes
- Set up Prometheus monitoring for sync metrics
- Configure alerts for sync failures
- Enable Pod Disruption Budgets (PDBs)

---

## Example 12: Multi-Region AWS Setup

ESO deployment that can access secrets from multiple AWS regions.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesExternalSecrets
metadata:
  name: external-secrets-multi-region
spec:
  target_cluster:
    cluster_name: prod-eks-cluster
  namespace:
    value: external-secrets
  
  poll_interval_seconds: 60
  
  container:
    resources:
      limits:
        cpu: 1000m
        memory: 1Gi
      requests:
        cpu: 50m
        memory: 100Mi
  
  eks:
    # Primary region - controller defaults to this
    region: us-east-1
    irsa_role_arn_override: arn:aws:iam::123456789012:role/eso-multi-region
```

**IAM policy for multi-region access:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": [
        "arn:aws:secretsmanager:us-east-1:123456789012:secret:*",
        "arn:aws:secretsmanager:us-west-2:123456789012:secret:*",
        "arn:aws:secretsmanager:eu-west-1:123456789012:secret:*"
      ]
    }
  ]
}
```

**Then create multiple SecretStores** (one per region):
```yaml
---
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-us-east-1
  namespace: default
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets
---
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-us-west-2
  namespace: default
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-west-2
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets
```

---

## Using External Secrets After Deployment

Once ESO is deployed, create `SecretStore` and `ExternalSecret` resources to sync secrets.

### Step 1: Create a SecretStore (GCP example)

```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: gcpsm-secret-store
  namespace: default
spec:
  provider:
    gcpsm:
      projectID: my-gcp-project
      auth:
        workloadIdentity:
          clusterLocation: us-central1
          clusterName: prod-gke-cluster
          serviceAccountRef:
            name: external-secrets
```

### Step 2: Create an ExternalSecret

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: database-credentials
  namespace: default
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: gcpsm-secret-store
    kind: SecretStore
  target:
    name: db-creds  # Kubernetes Secret name
    creationPolicy: Owner
  data:
  - secretKey: username
    remoteRef:
      key: prod-db-username  # Cloud secret name
  - secretKey: password
    remoteRef:
      key: prod-db-password
```

**Result:** A Kubernetes Secret named `db-creds` with keys `username` and `password` synced from Google Cloud Secret Manager.

### Step 3: Use the Secret in Your Application

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: my-app:latest
        env:
        - name: DB_USERNAME
          valueFrom:
            secretKeyRef:
              name: db-creds
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-creds
              key: password
```

---

## Provider-Specific SecretStore Examples

### AWS Secrets Manager (EKS)

```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secret-store
  namespace: default
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets
```

### Azure Key Vault (AKS)

```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: azure-secret-store
  namespace: default
spec:
  provider:
    azurekv:
      vaultUrl: https://prod-keyvault.vault.azure.net
      authType: WorkloadIdentity
      serviceAccountRef:
        name: external-secrets
```

### ClusterSecretStore (Cluster-Wide)

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ClusterSecretStore
metadata:
  name: gcpsm-cluster-store
spec:
  provider:
    gcpsm:
      projectID: my-gcp-project
      auth:
        workloadIdentity:
          clusterLocation: us-central1
          clusterName: prod-gke-cluster
          serviceAccountRef:
            name: external-secrets
            namespace: external-secrets
```

**Use ClusterSecretStore when:**
- All namespaces should access the same secret backend
- Centralized secret management by platform team
- Reduces duplication of SecretStore definitions

**Use namespaced SecretStore when:**
- Different teams/namespaces use different cloud credentials
- Fine-grained RBAC per namespace
- Multi-tenant clusters

---

## Common Patterns

### Pattern 1: JSON Secret with Multiple Keys

Cloud secret stores often store JSON blobs. Extract specific keys:

**Cloud secret (AWS Secrets Manager):**
```json
{
  "host": "db.example.com",
  "port": "5432",
  "username": "admin",
  "password": "secret123"
}
```

**ExternalSecret:**
```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-connection
spec:
  secretStoreRef:
    name: aws-secret-store
  target:
    name: db-connection
  data:
  - secretKey: DB_HOST
    remoteRef:
      key: prod/database/connection
      property: host
  - secretKey: DB_PORT
    remoteRef:
      key: prod/database/connection
      property: port
  - secretKey: DB_USER
    remoteRef:
      key: prod/database/connection
      property: username
  - secretKey: DB_PASS
    remoteRef:
      key: prod/database/connection
      property: password
```

### Pattern 2: Fetch Entire JSON as One Key

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: app-config
spec:
  secretStoreRef:
    name: gcpsm-secret-store
  target:
    name: app-config
  dataFrom:
  - extract:
      key: prod/app/config  # Entire JSON becomes Secret data
```

### Pattern 3: Secret Templating

Transform secrets or combine multiple sources:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: database-url
spec:
  secretStoreRef:
    name: aws-secret-store
  target:
    name: database-url
    template:
      data:
        DATABASE_URL: "postgresql://{{ .username }}:{{ .password }}@{{ .host }}:{{ .port }}/{{ .database }}"
  data:
  - secretKey: username
    remoteRef:
      key: prod/db/username
  - secretKey: password
    remoteRef:
      key: prod/db/password
  - secretKey: host
    remoteRef:
      key: prod/db/host
  - secretKey: port
    remoteRef:
      key: prod/db/port
  - secretKey: database
    remoteRef:
      key: prod/db/name
```

---

## Troubleshooting

### Verify ESO is Running

```bash
# Check ESO pods
kubectl get pods -n external-secrets

# Should see:
# external-secrets-<hash>         1/1     Running
# external-secrets-webhook-<hash> 1/1     Running
# external-secrets-cert-controller-<hash> 1/1     Running

# Check logs
kubectl logs -n external-secrets deployment/external-secrets
```

### Check SecretStore Status

```bash
# List SecretStores
kubectl get secretstore -A

# Describe to see validation status
kubectl describe secretstore gcpsm-secret-store -n default
```

Look for `Conditions` with `Type: Ready` and `Status: True`.

### Debug ExternalSecret Sync Failures

```bash
# Check ExternalSecret status
kubectl get externalsecret database-credentials -n default
kubectl describe externalsecret database-credentials -n default

# Common status conditions:
# - SecretSynced: True = Success
# - SecretSynced: False = Sync failed (check message)
# - Ready: True = Secret created and healthy
```

### Common Issues

**Authentication Errors:**

**GKE:**
```
Error: permission denied on resource "projects/my-project/secrets/my-secret"
```
**Fix:** Grant GCP service account `secretmanager.secretAccessor` role

**EKS:**
```
Error: User: arn:aws:sts::123456789012:assumed-role/... is not authorized
```
**Fix:** Verify IRSA role has SecretsManager permissions and trust policy

**AKS:**
```
Error: Caller is not authorized to perform action on resource
```
**Fix:** Verify Managed Identity has "Key Vault Secrets User" role

**Secret Not Found:**
```
Error: secret "prod-db-password" not found
```
**Fix:** Verify secret exists in cloud provider and exact name matches

**Polling Not Working:**
```
ExternalSecret shows old value after updating cloud secret
```
**Possible causes:**
- Poll interval hasn't elapsed yet (wait longer or manually trigger)
- Cloud provider propagation delay
- ESO controller restart needed (rare)

---

## Verification

After deployment:

```bash
# 1. Verify ESO is running
kubectl get pods -n external-secrets

# 2. Check ClusterSecretStore (created automatically by the module)
kubectl get clustersecretstore

# 3. Create a test ExternalSecret
kubectl apply -f test-external-secret.yaml

# 4. Wait for sync (based on refreshInterval)
kubectl get externalsecret test-secret -w

# 5. Verify Kubernetes Secret was created
kubectl get secret test-secret-data

# 6. Decode and verify secret content
kubectl get secret test-secret-data -o jsonpath='{.data.my-key}' | base64 -d
```

---

## Next Steps

1. Choose the example that matches your cloud provider and requirements
2. Customize with your cluster ID, project/account IDs, and resource requirements
3. Deploy via `openmcf deploy`
4. Verify ESO pods are running
5. Create SecretStore and ExternalSecret CRs for your secrets
6. Monitor sync status and ESO logs
7. Update applications to consume the synced Kubernetes Secrets

For more details, see:
- [README](README.md) - Overview and feature details
- [Research Documentation](docs/README.md) - Deep dive into ESO architecture and best practices
- [External Secrets Operator Official Docs](https://external-secrets.io/)

---

## Cost Calculator

Estimate your monthly cloud API costs:

```
Monthly API Calls = (Number of Secrets) × (86400 / Poll Interval) × 30

Example 1: 50 secrets, 30-second poll
= 50 × (86400 / 30) × 30
= 50 × 2880 × 30
= 4,320,000 API calls/month

AWS Secrets Manager: 4,320,000 / 10,000 × $0.05 = $21.60/month
GCP Secret Manager: 4,320,000 / 10,000 × $0.06 = $25.92/month

Example 2: 100 secrets, 1-hour poll
= 100 × (86400 / 3600) × 30
= 100 × 24 × 30
= 72,000 API calls/month

AWS Secrets Manager: 72,000 / 10,000 × $0.05 = $0.36/month
GCP Secret Manager: 72,000 / 10,000 × $0.06 = $0.43/month
```

**Recommendation:** Start with 1-hour polling for most use cases, adjust based on secret rotation requirements.


