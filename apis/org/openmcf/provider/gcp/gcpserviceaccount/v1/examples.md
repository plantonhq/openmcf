# GcpServiceAccount Examples

The `GcpServiceAccount` component lets you declaratively create Google Cloud service accounts, optionally generate a
key, and grant IAM roles at the project or organization level—all through the familiar OpenMCF YAML workflow. The
examples below illustrate typical scenarios.

---

## 1. Minimal Service Account (no key, default project)

Creates a service account in the provider's default project without generating a key or granting any IAM roles.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: analytics-sa
spec:
  serviceAccountId: analytics-sa
```

---

## 2. Service Account with Project-Level Roles and Key (Literal Value)

Creates the service account in a specific project using a literal project ID, generates a JSON key, and binds two project-level roles.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: logging-writer-sa
spec:
  serviceAccountId: logging-writer
  projectId:
    value: my-app-prod-1234
  createKey: true
  projectIamRoles:
    - roles/logging.logWriter
    - roles/monitoring.metricWriter
```

---

## 3. Service Account with Organization-Level Roles

Creates the account, skips key creation, assigns one project role and one organization-level role.  
Useful when the service account must operate across multiple projects in the same organization.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: org-auditor-sa
spec:
  serviceAccountId: org-auditor
  projectId:
    value: shared-infra-5678
  orgId: "123456789012"
  createKey: false
  projectIamRoles:
    - roles/viewer
  orgIamRoles:
    - roles/resourcemanager.organizationViewer
```

---

## 4. Service Account with Foreign Key Reference

References a `GcpProject` resource instead of hardcoding the project ID. This enables clean dependency management
and ensures the service account is created only after the project exists.

### Project Resource (defined separately)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: myapp-project
spec:
  projectId: myapp-prod-123
  billingAccountId: 012345-ABCDEF-678910
  folderId: "123456789012"
```

### Service Account Resource (referencing the project)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: myapp-service-account
spec:
  serviceAccountId: myapp-worker
  projectId:
    ref:
      kind: GcpProject
      name: myapp-project  # References the project defined above
  projectIamRoles:
    - roles/logging.logWriter
    - roles/storage.objectViewer
```

### How It Works

1. OpenMCF creates the `GcpProject` resource first
2. Once the project exists, it extracts `status.outputs.project_id`
3. The `GcpServiceAccount` resource uses that project ID automatically
4. Ensures proper dependency ordering (service account created after project)

---

## 5. Service Account for GKE Workload Identity

Creates a service account intended to be used with GKE Workload Identity. No key is created—the account
is attached to Kubernetes pods via Workload Identity Federation.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: gke-workload-sa
  org: platform-team
  env: production
spec:
  serviceAccountId: gke-workload
  projectId:
    value: gke-prod-project
  createKey: false  # Keyless auth via Workload Identity
  projectIamRoles:
    - roles/logging.logWriter
    - roles/monitoring.metricWriter
    - roles/cloudtrace.agent
    - roles/secretmanager.secretAccessor
```

### Next Steps

After creating this service account, bind it to a Kubernetes service account using `GcpGkeWorkloadIdentityBinding`:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeWorkloadIdentityBinding
metadata:
  name: gke-workload-binding
spec:
  projectId:
    ref:
      kind: GcpServiceAccount
      name: gke-workload-sa
      fieldPath: spec.project_id
  serviceAccountEmail:
    ref:
      kind: GcpServiceAccount
      name: gke-workload-sa
  kubernetesServiceAccountName: my-k8s-service-account
  kubernetesNamespace: my-namespace
```

---

### Deploying the Examples

```bash
# Validate the manifest
openmcf validate --manifest example.yaml

# Deploy with Pulumi
openmcf pulumi up --manifest example.yaml --stack org/project/stack

# …or with Terraform
openmcf terraform apply --manifest example.yaml --stack org/project/stack
```

These examples demonstrate how the `GcpServiceAccount` component streamlines service-account management—eliminating
manual steps for key generation and IAM binding while fitting seamlessly into OpenMCF's multi-cloud workflow.
