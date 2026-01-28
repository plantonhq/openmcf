# GCP GKE Workload Identity Binding Examples

This document provides copy-paste ready examples for common Workload Identity binding scenarios.

## Table of Contents

1. [Minimal Example](#minimal-example)
2. [cert-manager with Cloud DNS](#cert-manager-with-cloud-dns)
3. [Application Accessing Cloud Storage](#application-accessing-cloud-storage)
4. [ExternalDNS with Cloud DNS](#externaldns-with-cloud-dns)
5. [Multi-Environment Setup](#multi-environment-setup)
6. [Using Foreign Key References](#using-foreign-key-references)

---

## Minimal Example

The simplest possible binding for testing purposes.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: test-binding
spec:
  projectId: "my-project-123"
  serviceAccountEmail: "test-gsa@my-project-123.iam.gserviceaccount.com"
  ksaNamespace: "default"
  ksaName: "default"
```

**What this does:**

- Binds the `default` ServiceAccount in the `default` namespace
- Allows pods using this KSA to impersonate `test-gsa@my-project-123.iam.gserviceaccount.com`

**Test it:**

```bash
kubectl run test-pod \
  --image=google/cloud-sdk:slim \
  --rm -it \
  -- gcloud auth list
```

You should see the GSA email as the active account.

---

## cert-manager with Cloud DNS

cert-manager needs to create TXT records in Cloud DNS to solve DNS-01 challenges for Let's Encrypt certificates.

### Step 1: Create the Google Service Account

```bash
# Create the GSA
gcloud iam service-accounts create dns01-solver \
  --display-name="cert-manager DNS-01 solver" \
  --project=prod-project

# Grant DNS admin permissions
gcloud projects add-iam-policy-binding prod-project \
  --member="serviceAccount:dns01-solver@prod-project.iam.gserviceaccount.com" \
  --role="roles/dns.admin"
```

### Step 2: Create the Workload Identity Binding

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: cert-manager-dns-binding
spec:
  projectId: "prod-project"
  serviceAccountEmail: "dns01-solver@prod-project.iam.gserviceaccount.com"
  ksaNamespace: "cert-manager"
  ksaName: "cert-manager"
```

### Step 3: Deploy cert-manager

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml
```

### Step 4: Create ClusterIssuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: admin@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        cloudDNS:
          project: prod-project
          # No credentials needed - Workload Identity auto-detected
```

### Step 5: Request a Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - example.com
  - "*.example.com"
```

---

## Application Accessing Cloud Storage

A microservice needs read-only access to a Cloud Storage bucket containing configuration files.

### Step 1: Create the GSA and Grant Permissions

```bash
# Create the GSA
gcloud iam service-accounts create my-app-gsa \
  --display-name="My App GSA" \
  --project=prod-project

# Grant storage viewer permissions on a specific bucket
gsutil iam ch \
  serviceAccount:my-app-gsa@prod-project.iam.gserviceaccount.com:objectViewer \
  gs://my-config-bucket
```

### Step 2: Create the Workload Identity Binding

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: my-app-gcs-binding
spec:
  projectId: "prod-project"
  serviceAccountEmail: "my-app-gsa@prod-project.iam.gserviceaccount.com"
  ksaNamespace: "my-app"
  ksaName: "my-app-ksa"
```

### Step 3: Create the Kubernetes ServiceAccount

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-app-ksa
  namespace: my-app
```

### Step 4: Deploy Your Application

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      serviceAccountName: my-app-ksa  # Use the bound KSA
      containers:
      - name: app
        image: my-app:v1.0.0
        env:
        - name: CONFIG_BUCKET
          value: "gs://my-config-bucket/config.json"
        # No GOOGLE_APPLICATION_CREDENTIALS needed!
        # Google Cloud SDKs automatically detect Workload Identity
```

---

## ExternalDNS with Cloud DNS

ExternalDNS automatically creates and updates DNS records based on Kubernetes Ingress and Service resources.

### Step 1: Create the GSA

```bash
# Create the GSA
gcloud iam service-accounts create external-dns \
  --display-name="ExternalDNS" \
  --project=prod-project

# Grant DNS admin permissions
gcloud projects add-iam-policy-binding prod-project \
  --member="serviceAccount:external-dns@prod-project.iam.gserviceaccount.com" \
  --role="roles/dns.admin"
```

### Step 2: Create the Workload Identity Binding

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: external-dns-binding
spec:
  projectId: "prod-project"
  serviceAccountEmail: "external-dns@prod-project.iam.gserviceaccount.com"
  ksaNamespace: "external-dns"
  ksaName: "external-dns"
```

### Step 3: Deploy ExternalDNS

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: external-dns
  namespace: external-dns
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: external-dns
rules:
- apiGroups: [""]
  resources: ["services", "endpoints", "pods"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["extensions", "networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: external-dns-viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: external-dns
subjects:
- kind: ServiceAccount
  name: external-dns
  namespace: external-dns
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: external-dns
  namespace: external-dns
spec:
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: external-dns
  template:
    metadata:
      labels:
        app: external-dns
    spec:
      serviceAccountName: external-dns
      containers:
      - name: external-dns
        image: registry.k8s.io/external-dns/external-dns:v0.14.0
        args:
        - --source=service
        - --source=ingress
        - --domain-filter=example.com
        - --provider=google
        - --google-project=prod-project
        - --registry=txt
        - --txt-owner-id=my-cluster
        # No credentials needed - Workload Identity handles it
```

---

## Multi-Environment Setup

Separate bindings for development, staging, and production environments using separate GCP projects.

### Development Environment

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: my-app-dev-binding
spec:
  projectId: "dev-project"
  serviceAccountEmail: "my-app-dev@dev-project.iam.gserviceaccount.com"
  ksaNamespace: "my-app-dev"
  ksaName: "my-app-ksa"
```

### Staging Environment

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: my-app-staging-binding
spec:
  projectId: "staging-project"
  serviceAccountEmail: "my-app-staging@staging-project.iam.gserviceaccount.com"
  ksaNamespace: "my-app-staging"
  ksaName: "my-app-ksa"
```

### Production Environment

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: my-app-prod-binding
spec:
  projectId: "prod-project"
  serviceAccountEmail: "my-app-prod@prod-project.iam.gserviceaccount.com"
  ksaNamespace: "my-app-prod"
  ksaName: "my-app-ksa"
```

**Security benefit**: Each environment uses a separate GCP project, creating isolated Workload Identity Pools and preventing cross-environment identity confusion.

---

## Using Foreign Key References

OpenMCF supports referencing other components using foreign keys, eliminating hardcoded values.

### Example: Reference a GcpProject Component

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: app-binding-with-refs
spec:
  projectId:
    fromReference:
      kind: GcpProject
      name: my-prod-project
      fieldPath: status.outputs.project_id
  serviceAccountEmail:
    fromReference:
      kind: GcpServiceAccount
      name: my-app-gsa
      fieldPath: status.outputs.email
  ksaNamespace: "my-app"
  ksaName: "my-app-ksa"
```

**Benefits:**

- No hardcoded project IDs or emails
- Changes to the referenced component automatically propagate
- Clearer dependencies between components

### Example: With Value Directly

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: app-binding-direct
spec:
  projectId:
    value: "prod-project-123"
  serviceAccountEmail:
    value: "my-app-gsa@prod-project-123.iam.gserviceaccount.com"
  ksaNamespace: "my-app"
  ksaName: "my-app-ksa"
```

---

## Complete Example: Web Application with Cloud SQL

A real-world example showing a web app that needs to:
- Access Cloud SQL (using Cloud SQL Proxy sidecar)
- Read from Cloud Storage
- Write logs to Cloud Logging

### Step 1: Create GSA with Multiple Roles

```bash
# Create the GSA
gcloud iam service-accounts create webapp-gsa \
  --display-name="Web App GSA" \
  --project=prod-project

# Grant Cloud SQL Client role
gcloud projects add-iam-policy-binding prod-project \
  --member="serviceAccount:webapp-gsa@prod-project.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# Grant Storage Object Viewer on specific bucket
gsutil iam ch \
  serviceAccount:webapp-gsa@prod-project.iam.gserviceaccount.com:objectViewer \
  gs://webapp-assets

# Grant Logging Writer
gcloud projects add-iam-policy-binding prod-project \
  --member="serviceAccount:webapp-gsa@prod-project.iam.gserviceaccount.com" \
  --role="roles/logging.logWriter"
```

### Step 2: Create Workload Identity Binding

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: webapp-binding
spec:
  projectId: "prod-project"
  serviceAccountEmail: "webapp-gsa@prod-project.iam.gserviceaccount.com"
  ksaNamespace: "webapp"
  ksaName: "webapp-ksa"
```

### Step 3: Deploy Application with Cloud SQL Proxy

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: webapp-ksa
  namespace: webapp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
  namespace: webapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
    spec:
      serviceAccountName: webapp-ksa  # Bound to webapp-gsa
      containers:
      - name: webapp
        image: webapp:v1.0.0
        env:
        - name: DB_HOST
          value: "127.0.0.1"
        - name: DB_PORT
          value: "5432"
        - name: DB_NAME
          value: "webapp_db"
        - name: ASSETS_BUCKET
          value: "gs://webapp-assets"
      - name: cloud-sql-proxy
        image: gcr.io/cloud-sql-connectors/cloud-sql-proxy:2.8.0
        args:
        - "--structured-logs"
        - "--port=5432"
        - "prod-project:us-central1:webapp-db"
        securityContext:
          runAsNonRoot: true
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
```

**No credentials needed!** The webapp container and cloud-sql-proxy sidecar both automatically use the GSA via Workload Identity.

---

## Troubleshooting Examples

### Test Pod to Verify Authentication

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: workload-identity-test
  namespace: my-app
spec:
  serviceAccountName: my-app-ksa
  containers:
  - name: test
    image: google/cloud-sdk:slim
    command: ["sleep", "3600"]
```

**Test commands:**

```bash
# Check which account is active
kubectl exec -it workload-identity-test -n my-app -- gcloud auth list

# Try accessing a GCS bucket
kubectl exec -it workload-identity-test -n my-app -- gsutil ls gs://my-bucket

# Check Cloud Logging permissions
kubectl exec -it workload-identity-test -n my-app -- gcloud logging write test-log "Test message"
```

---

## Additional Resources

- **[Main README](README.md)**: Overview and troubleshooting guide
- **[Research Documentation](docs/README.md)**: Deep dive into Workload Identity architecture and design decisions
- **[Pulumi Deployment](iac/pulumi/README.md)**: Deploy using Pulumi
- **[Terraform Deployment](iac/tf/README.md)**: Deploy using Terraform

## Support

For issues or questions, see the [OpenMCF documentation](https://openmcf.org).


