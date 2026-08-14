# AwsSagemakerNotebookInstance — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI notebook instance
(`aws_sagemaker_notebook_instance`) with its folded lifecycle
configuration
(`aws_sagemaker_notebook_instance_lifecycle_configuration`).

Module facts worth knowing before editing:

- **Most instance changes ride the provider's stop-update-start
  choreography** — SageMaker requires a Stopped instance for
  UpdateNotebookInstance; budget several minutes per change.
- **Growing `volume_size` updates in place, SHRINKING replaces the
  instance** (provider CustomizeDiff mirrors AWS's no-shrink rule).
- **The lifecycle scripts are sent base64-encoded** — the module
  encodes; the spec carries plain shell. They run as root with a
  5-minute limit.
- **Clearing a script upstream does NOT clear it in AWS** (the
  provider's update omits empty fields) — replacing the text is the
  reliable path, taught on the spec fields.
- **The lifecycle configuration renders only when configured**
  (`count` on `has_lifecycle`) under a stable derived name
  (`<name>-lifecycle`).

Outputs mirror the Pulumi module key-for-key:
`notebook_instance_name`, `notebook_instance_arn`, `url`,
`network_interface_id`, `lifecycle_config_name`.
