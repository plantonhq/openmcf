<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Notebook Instance" width="80"/>
</p>

# AWS SageMaker Notebook Instance

Create and manage [Amazon SageMaker AI notebook instances](https://docs.aws.amazon.com/sagemaker/latest/dg/nbi.html)
— managed EC2 instances running Jupyter — together with their folded
lifecycle configuration (bootstrap scripts).

## What Gets Created

- **A notebook instance** on an `ml.*` instance type with an ML
  storage volume (`volume_size_gb`, 5–16384 GB) and an IAM role for
  the AWS calls made from it.
- **A lifecycle configuration** (folded) carrying plain-shell
  bootstrap scripts — `on_create` runs once at creation, `on_start` on
  every start; the modules base64-encode them, and they run as root
  with AWS's 5-minute limit.
- Optional: VPC placement (`subnet_id` + `security_group_ids`) with
  direct internet access disabled, KMS volume encryption, root-access
  and IMDSv2 lockdown, platform selection, and Git code repositories
  cloned as the working directory.

## Updates Ride a Stop-Update-Start Cycle

SageMaker stops the instance to apply most changes and restarts it
afterwards — budget several minutes per change. Growing
`volume_size_gb` updates in place; shrinking it REPLACES the instance
(AWS cannot shrink a volume).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
