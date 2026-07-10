# AWS SageMaker Domain at the Full Provider Surface

**Date**: July 10, 2026
**Type**: Enhancement | Breaking Change
**Components**: AWS Provider, API Definitions, Terraform Modules, Pulumi Modules, Testing Framework

## Summary

`AwsSagemakerDomain` was rebuilt from roughly a third of the
`aws_sagemaker_domain` surface to the complete provider contract with
dual-engine parity. The hand-written Terraform contract that had never
been deployable against its own proto was replaced by the generator under
the drift guard, the last `planton.org/*` tag straggler retired, naming
converged to `metadata.name`, and the kind gained its first end-to-end
suite — live dual-engine runs green on both scenarios with zero-orphan
teardown via `home_efs_retention_policy: Delete`.

## Problem Statement / Motivation

The spec covered only the baseline JupyterLab/JupyterServer/KernelGateway
arms. Canvas, Code Editor, RStudio, TensorBoard, Studio web-portal
governance, default space inheritance, trusted identity propagation, tag
propagation, and the orphaned-EFS retention policy were all unrepresentable.
Beyond breadth, the module pair carried latent defects that blocked every
deploy path:

### Pain Points

- **The TF module was never deployable** (the Cognito/DynamoDB class): the
  hand-written `variables.tf` renamed proto fields to provider names
  (`execution_role` vs `execution_role_arn`, `security_groups` vs
  `security_group_ids`, …). The proto→tfvars converter emits proto field
  names, so every manifest — including the minimal hack manifest — failed
  HCL object-type conversion at plan time.
- **Cross-engine parity defect**: Terraform dropped
  `sagemaker_image_version_alias` and `sagemaker_image_version_arn` in
  every resource-spec block while Pulumi sent them.
- **Tag straggler**: both engines still emitted bare `planton.org/*` keys
  instead of the settled `planton.ai/*` identity-tag set.
- **Naming basis**: both engines keyed `domain_name` off `metadata.id`
  instead of `metadata.name` (family convention).
- **`Pulumi.yaml` carried `options: binary: main` residue** (fails only at
  preview/deploy) and the entrypoint anatomy was incomplete (no `Makefile`).
- **Provider pin `~> 5.0`** — below the floor for
  `trusted_identity_propagation_settings` (provider 6.33.0).
- **Zero E2E artifacts** of its own (verifier existed; no profile,
  scenarios, conformance case, or registry prerequisites).

## Solution / What's New

### Full domain surface

Top-level additions: `tag_propagation`, `home_efs_retention_policy`
(Retain/Delete — the orphaned-EFS fix), `app_security_group_management`,
flattened domain-settings arms (`execution_role_identity_config`,
`trusted_identity_propagation_status`, RStudio domain settings), and the
complete `default_space_settings` inheritance plane. User-settings depth
now covers Canvas (incl. OAuth `secret_arn` as a sensitive-exempt
reference), Code Editor, classic Jupyter Server, R Session, RStudio Pro,
TensorBoard, Studio web-portal hiding, auto-mount-home-EFS, custom EFS
mounts (FK → `AwsElasticFileSystem`), and POSIX identity bounds. New
output: `single_sign_on_managed_application_instance_id`.

Provider-sourced couplings encoded as CELs: trusted identity propagation
requires SSO auth; `app_security_group_management` only when RStudio is
configured; RStudio `user_group` only when `access_status` is ENABLED.

### Cross-engine convergence

- Generator-owned `variables.tf` under the drift guard; TF module
  rewritten with presence-driven dynamic blocks; floor `>= 6.33.0`.
- Pulumi module extended to full parity (`user_settings.go`,
  `space_settings.go`); `Pulumi.yaml` residue removed; `Makefile` added.
- Identity tags converged to `planton.ai/*` via `awstagkeys`; naming
  converged to `metadata.name` (breaking; ml-workbench chart shape
  preserved — it composes only region/auth/vpc/subnets/execution role).
- `idle_settings` kept folded directly under each app-settings block
  (the Firehose/DRA single-purpose-wrapper class — both engines
  reconstruct the SDK's `app_lifecycle_management` shape).

### SageMaker satellites deferred

The whole satellite family (user profile, space, lifecycle config, image,
app, MLflow, model/endpoint planes) is recorded in DD-005. The join path
stays open via `domain_id`.

## Implementation Details

- **Guards**: `AwsSagemakerDomain` enrolled in `migratedKinds`; outputs
  conformance case; 545+ lines of spec tests (a case per new rule/arm);
  `validate-refs --check` and `secret-coverage --check` green.
- **E2E**: profile, minimal + full-surface scenarios (both with
  `homeEfsRetentionPolicy: Delete`), dual-engine entrypoints, SageMaker
  execution-role fixture document (13th IAM doc), registry prerequisites
  `[AwsSubnet, AwsIamRole]`. Live exclusions recorded: SSO + trusted
  identity propagation, RStudio (Posit license), Canvas OAuth, custom
  images/EFS mounts.
- **Docs**: README, docs/README, catalog page rewritten timeless (v2
  roadmap narration retired); presets 01–03 refreshed; new preset
  `04-governed-canvas-workspace`; hack manifest extended to full surface;
  site catalog mirror regenerated.
- **Workflow uplift**: forge rule gains secret-coverage exemption
  discipline (annotate only what the gate flags) and E2E prerequisite
  annotation audit guidance (unresolved refs fail Terraform live-only).

## Benefits

- Manifest authors can express the full SageMaker Studio workspace
  contract without out-of-band console steps for governance arms.
- The TF module is deployable for the first time — the proto and
  variables contract finally agree.
- Dual-engine parity on image-version fields, tags, and naming.
- Zero-orphan teardown: `Delete` retention destroys the auto-created EFS
  home file system with the domain.

## Impact

- **Breaking**: `domain_name` now derives from `metadata.name`; identity
  tags move from `planton.org/*` to `planton.ai/*`. ml-workbench chart
  fields unchanged.
- **Users/operators**: richer presets including governed Canvas; E2E
  profile documents ~9–12 min per engine lane.
- **Maintainers**: generator-owned contract prevents the hand-written
  rename class from recurring.

## Related Work

- DD-004 fold entries (session 041): flattened provider wrappers
  documented.
- DD-005 SageMaker satellite demand-check ledger (session 041).
- Closes the breadth-wave tail — Phase 1 charts wave is next.

---

**Status**: ✅ Production Ready
**Timeline**: Session 043 (breadth-wave tail; concurrent with sessions 041–042)
