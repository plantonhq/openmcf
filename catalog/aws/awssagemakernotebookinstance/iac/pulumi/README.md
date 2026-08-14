# AwsSagemakerNotebookInstance — Pulumi module (Go)

Deploys an Amazon SageMaker AI notebook instance
(`sagemaker.NotebookInstance`) with its folded lifecycle configuration
(`sagemaker.NotebookInstanceLifecycleConfiguration`).

Module facts worth knowing before editing:

- **Most instance changes ride the provider's stop-update-start
  choreography** — SageMaker requires a Stopped instance for
  UpdateNotebookInstance; budget several minutes per change.
- **Growing `VolumeSize` updates in place, SHRINKING replaces the
  instance** (provider-enforced, mirroring AWS's no-shrink rule).
- **The lifecycle scripts are sent base64-encoded** — the module
  encodes (`base64.StdEncoding`); the spec carries plain shell. They
  run as root with a 5-minute limit.
- **Clearing a script upstream does NOT clear it in AWS** (the
  provider's update omits empty fields) — replacing the text is the
  reliable path, taught on the spec fields.
- **The lifecycle configuration renders only when configured**, under
  a stable derived name (`<name>-lifecycle`); the
  `lifecycle_config_name` output is the empty string otherwise.

Outputs mirror the Terraform module key-for-key:
`notebook_instance_name`, `notebook_instance_arn`, `url`,
`network_interface_id`, `lifecycle_config_name`.
