# Locked-Down GPU Notebook

This preset puts a GPU under Jupyter with the security posture
tightened — no root for users, IMDSv2 only, the current Amazon Linux
2023 platform, and the team's repository cloned as the working
directory.

## When to Use

- Deep-learning experimentation that needs a GPU (`ml.g4dn.xlarge`)
- Shared or governed environments where notebook users should not
  hold root

## What You Get

- A GPU instance on the AL2023 platform with a 100 GB volume for
  datasets and checkpoints
- `rootAccess: Disabled` (lifecycle scripts, if added, still run as
  root) and IMDSv2-only instance metadata — the hardened choice
- The default code repository cloned into the working directory on
  start

## Customize

- Add up to three `additionalCodeRepositories` cloned alongside the
  default one
- Confine the notebook to your VPC with `subnetId`,
  `securityGroupIds`, and `directInternetAccess: Disabled` — plan a
  NAT or endpoint path for SageMaker API calls, and decide up front:
  changing network posture or platform replaces the instance
- Add `kmsKeyArn` to encrypt the storage volume with your own key
