# Terraform Module to Deploy AWSLambda

This module deploys an `AWSLambda` resource using Terraform via the OpenMCF CLI (tofu).

## CLI

```bash
openmcf tofu init --manifest hack/manifest.yaml
openmcf tofu plan --manifest hack/manifest.yaml
openmcf tofu apply --manifest hack/manifest.yaml --auto-approve
openmcf tofu destroy --manifest hack/manifest.yaml --auto-approve
```

- Credentials are provided via the CLI stack input, not stored in the manifest `spec`.
- Example manifest: see `apis/project/planton/provider/aws/awslambda/v1/iac/hack/manifest.yaml`.
