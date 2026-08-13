# AWS SageMaker Notebook Instance

A managed Jupyter workstation — an EC2 instance SageMaker runs for
you, bootstrapped by declarative lifecycle scripts, with an ML storage
volume, an IAM role for the AWS calls made from it, and optional VPC
confinement and security lockdown.

## What Gets Created

- A notebook instance on an `ml.*` type (ml.t3.medium is the cheapest
  current-generation choice, ~$0.05/hour) with a 5–16384 GB storage
  volume.
- A folded lifecycle configuration: plain-shell `on_create` (runs
  once) and `on_start` (runs every start) scripts, base64-encoded by
  the modules, run as root under AWS's 5-minute limit.
- Optional: VPC placement with direct internet access disabled, KMS
  volume encryption, root-access and IMDSv2 lockdown, platform
  selection, and up to four Git repositories cloned into the working
  directory.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateNotebookInstance` and its siblings).

### AWS Account

- An IAM role trusting `sagemaker.amazonaws.com` (`role_arn`) — the
  notebook assumes it for every AWS call made from Jupyter.
- For VPC confinement (`direct_internet_access: Disabled`): a subnet
  and security groups, plus a NAT or endpoint path so the notebook can
  still reach SageMaker APIs.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, instance
type, and role, add bootstrap scripts if needed, and deploy.

### CLI

```bash
planton apply -f notebook.yaml
```

## After Deploy

- `notebook_instance_name` / `notebook_instance_arn` identify the
  instance; `url` opens Jupyter; `lifecycle_config_name` echoes the
  folded configuration.
- Most changes stop the instance, apply, and restart it — budget
  several minutes and expect the notebook to be unavailable during the
  cycle.
- Growing the volume updates in place; shrinking it replaces the
  instance — size generously up front.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
