# GcpGkeWorkloadIdentityBinding Terraform Module

This Terraform module creates one additive IAM grant on a Google Service
Account: `roles/iam.workloadIdentityUser` to the constructed
workload-identity principal
`serviceAccount:<project>.svc.id.goog[<namespace>/<ksa>]`.

## Usage

### With Planton CLI

```bash
planton tofu apply --manifest workload-identity-binding.yaml
```

### Standalone Usage

```hcl
module "workload_identity_binding" {
  source = "./path/to/module"

  metadata = {
    name = "cert-manager-binding"
  }

  spec = {
    # StringValueOrRef fields are flattened to plain strings by the tfvars
    # converter before the module sees them.
    project_id            = "my-gcp-project"
    service_account_email = "cert-manager@my-gcp-project.iam.gserviceaccount.com"
    ksa_namespace         = "cert-manager"
    ksa_name              = "cert-manager"
  }
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| google | ~> 6.0 |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| metadata | Resource metadata including name | object | yes |
| spec.project_id | Project hosting the GKE cluster (names the workload-identity pool) | string | yes |
| spec.service_account_email | Email of the GSA receiving the grant | string | yes |
| spec.ksa_namespace | Kubernetes namespace of the ServiceAccount | string | yes |
| spec.ksa_name | Name of the Kubernetes ServiceAccount | string | yes |
| spec.condition | Optional IAM Condition (title, expression, description) | object | no |

## Outputs

| Name | Description |
|------|-------------|
| member | The constructed workload-identity principal |
| service_account_email | The bound GSA email — the value the KSA annotation needs |

## Required Permissions

The identity running the module needs `roles/iam.serviceAccountAdmin` (or
the narrower `iam.serviceAccounts.setIamPolicy`) on the target service
account.

## Notes

- The provider requires the fully-qualified service-account resource name;
  the module constructs `projects/-/serviceAccounts/<email>` — the IAM
  API's `-` wildcard infers the SA's project from the email, so
  cross-project GSAs work without extra configuration.
- The `iam.gke.io/gcp-service-account` annotation on the KSA is the
  Kubernetes half of the handshake and belongs to the workload's own
  deployment.
