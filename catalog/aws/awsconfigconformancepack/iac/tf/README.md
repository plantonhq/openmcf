# AwsConfigConformancePack — Terraform/OpenTofu module

Manages a conformance pack at one of two scopes: the account pack
(`aws_config_conformance_pack`) or the organization pack
(`aws_config_organization_conformance_pack`) — exactly one renders,
count-gated on `spec.organization_scope`.

Module facts worth knowing before editing:

- **The recorder is the consumer's contract.** AWS rejects pack
  deployment without an active Config recorder in the region — the
  recorder lives on AwsConfigRecorder, never here.
- **The template asymmetry is upstream.** The account resource
  accepts both template forms at once; the organization resource
  accepts exactly one — the spec CELs guarantee only legal
  combinations arrive.
- **Template drift is undetectable** (neither resource reads it
  back); imports re-assert the template on the first apply.
- **No tags render** — neither provider resource carries a tags
  argument.

Outputs mirror the Pulumi module key-for-key: `pack_name`, `pack_arn`
(the deployed scope's ARN).
