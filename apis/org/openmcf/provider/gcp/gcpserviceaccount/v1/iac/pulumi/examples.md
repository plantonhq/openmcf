# GcpServiceAccount Pulumi Examples

Below are several YAML examples demonstrating how to create and configure Google Cloud service accounts using OpenMCF's
`GcpServiceAccount` resource. After creating a manifest, you can apply it with Pulumi or Terraform via the OpenMCF CLI,
just like any other resource in the OpenMCF ecosystem.

```shell
# Pulumi
openmcf pulumi up --manifest <yaml-path> --stack <stack-name>

# Terraform
openmcf terraform apply --manifest <yaml-path> --stack <stack-name>
```

---

## Example 1: Minimal Service Account (No Key)

This example creates a service account with just the required fields. No IAM roles or key generation.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpServiceAccount
metadata:
  name: minimal-service-account
spec:
  serviceAccountId: minimal-sa
  projectId:
    value: my-gcp-project-123
```

---

## Example 2: Service Account with Project-Level Roles and Key

Creates a service account with IAM roles and generates a JSON key for authentication.

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

## Example 3: Service Account with Organization-Level Roles

Creates a service account with both project and organization-level IAM roles. Useful for cross-project operations.

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

## Example 4: Service Account with Foreign Key Reference

Reference a `GcpProject` resource instead of hardcoding the project ID. This ensures proper dependency ordering.

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
      name: myapp-project
  projectIamRoles:
    - roles/logging.logWriter
    - roles/storage.objectViewer
```

---

## Example 5: Service Account for GKE Workload Identity

Creates a service account intended for GKE Workload Identity. No key is created—pods authenticate via Workload Identity Federation.

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
  createKey: false
  projectIamRoles:
    - roles/logging.logWriter
    - roles/monitoring.metricWriter
    - roles/cloudtrace.agent
    - roles/secretmanager.secretAccessor
```

---

After deploying any of these manifests, you can confirm the newly created service account in your GCP account:

```shell
gcloud iam service-accounts list --project=<your-project-id>
gcloud iam service-accounts describe <sa-email> --project=<your-project-id>
```

You should see the service account with the specified roles and configuration. From there, you can attach it to GCP resources (Cloud Run, GKE, Compute Engine) or use it with Workload Identity Federation for external systems.
