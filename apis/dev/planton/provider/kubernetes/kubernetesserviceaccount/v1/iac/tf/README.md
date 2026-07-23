# Kubernetes ServiceAccount - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes ServiceAccount — the in-cluster identity pods run as. It covers image-pull secrets, the token-automount flag, and cloud workload-identity binding (GKE Workload Identity, EKS IRSA, Azure AD Workload Identity) expressed as ServiceAccount annotations. Behavior is kept in exact parity with the sibling Pulumi module.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Label merging, workload-identity annotation translation, RBAC subject
├── main.tf         # Creates kubernetes_service_account_v1 resource
├── outputs.tf      # Exports service_account_name, namespace, rbac_subject, workload_identity_handle
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the proto spec; the workload-identity oneof arrives flattened as an object with at most one of `gke`/`eks`/`aks` set
2. **Annotation Translation**: `locals.tf` translates the selected workload-identity arm into the cloud webhook's annotation and merges it over user annotations (workload-identity wins on collision)
3. **Resource Creation**: `main.tf` creates a single `kubernetes_service_account_v1` resource with dynamic `image_pull_secret` blocks
4. **Output Export**: Name, namespace, RBAC subject, and workload-identity handle are exported

## Workload-Identity Annotation Mapping

| Variable Arm | Annotation | Value |
|--------------|-----------|-------|
| `gke` | `iam.gke.io/gcp-service-account` | `service_account_email` |
| `eks` | `eks.amazonaws.com/role-arn` | `role_arn` |
| `aks` | `azure.workload.identity/client-id` | `client_id` |
| `aks` (optional) | `azure.workload.identity/tenant-id` | `tenant_id`, only when set |

## Token Automount

`spec.automount_service_account_token` is tri-state (`optional(bool)`, defaults to null). The Terraform provider cannot omit this attribute (its schema defaults it to true), so the module applies `true` when the variable is null — behaviorally identical to leaving the field unset, because the Kubernetes cluster default for an absent field is "mount the token". An explicit `false` or `true` is passed through as-is.

## Usage

```hcl
module "service_account" {
  source = "./iac/tf"

  metadata = {
    name = "dns-manager"
  }

  spec = {
    name      = "dns-manager"
    namespace = "dns-system"

    image_pull_secrets = ["registry-credentials"]

    workload_identity = {
      gke = {
        service_account_email = "dns-manager@my-project.iam.gserviceaccount.com"
      }
    }
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | ServiceAccount specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `service_account_name` | Name of the created ServiceAccount |
| `namespace` | Namespace of the created ServiceAccount |
| `rbac_subject` | `system:serviceaccount:<namespace>:<name>` |
| `workload_identity_handle` | Bound cloud identity handle, or empty |
