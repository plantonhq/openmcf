# AWS App Runner Family: Service Rebuilt with Companions Decomposed to First-Class Kinds

**Date**: July 8, 2026
**Type**: Feature (breaking)
**Components**: API Definitions, Provider Framework, Terraform Modules, Pulumi Modules, E2E Harness

## Summary

The AWS App Runner surface is decomposed to match AWS's own resource model: `AwsAppRunnerService` is rebuilt at the full provider surface with its embedded side resources split into three new first-class kinds — `AwsAppRunnerAutoScalingConfiguration` (enum 368), `AwsAppRunnerVpcConnector` (369), and `AwsAppRunnerObservabilityConfiguration` (370). The service gains folded custom domains (per-domain certificate-validation records as outputs), the WAF `web_acl_arn` association seam, a generator-owned Terraform contract, and the family's first-ever E2E coverage — all 8 live dual-engine lanes green with a zero-orphan sweep.

## What Was Built

### The decomposition

The old service module CREATED a VPC connector (from `subnet_ids`/`security_group_ids`) and an auto scaling configuration (from an embedded `auto_scaling` block) as side resources — with a live cross-engine naming divergence (Terraform suffixed `-vpc`/`-asc`, Pulumi used the bare name) and a Pulumi output contract that dropped two outputs whenever the side resources were absent. AWS designs all three companions as SHARED, versioned resources referenced by ARN across many services; embedding one per service forked what should be tuned in one place.

- **`AwsAppRunnerAutoScalingConfiguration` (368, `awsarasc`)** — the concurrency-based scaling policy: `max_concurrency` (1-200), `max_size`, `min_size` (warm floor, memory-only billing while idle), max>=min as CEL. Every value is create-time immutable by AWS design: changes register a NEW revision whose ARN rolls referencing services through the resource graph.
- **`AwsAppRunnerVpcConnector` (369, `awsarvc`)** — the managed ENI set for outbound VPC access: required subnet refs + required security-group refs (the provider marks both Required — requiredness honesty), total create-time immutability stated plainly, egress-only role stated loudly (the inbound plane is a different resource). Registry `prerequisites: [AwsSubnet, AwsSecurityGroup]`.
- **`AwsAppRunnerObservabilityConfiguration` (370, `awsaroc`)** — the X-Ray tracing policy with the provider's `trace_configuration.vendor` block. On the service, **the reference IS the enable switch**: the provider's separate `observability_enabled` boolean (a drift trap where enabled-without-configuration is representable) is unrepresentable in the spec — both engines derive it from the reference's presence.

### AwsAppRunnerService rebuilt (breaking)

- Embedded `auto_scaling` block and `subnet_ids`/`security_group_ids` removed; `auto_scaling_configuration_arn`, `vpc_connector_arn`, and `observability_configuration_arn` are now refs with default kinds pointing at the new companions. The two side-resource passthrough outputs removed.
- **Custom domains folded** as per-name blocks (`domain_name` + `enable_www_subdomain`): service-keyed, all-ForceNew, many-per-service — the established per-name fold class. Per-domain `dns_target` + certificate-validation CNAMEs export as a two-level repeated-message output composing directly into `AwsRoute53DnsRecord` nodes.
- **WAF seam closed**: optional `web_acl_arn` ref → `AwsWafWebAcl`, with `aws_wafv2_web_acl_association` materialized in both engines (the direction CloudFront models: the protected resource points at the web ACL). Proven live in the composed E2E lane.
- **`auto_deployments_enabled` honesty**: now a plain bool sent explicitly by both engines. AWS's own default is conditional (true for code repos, false for image repos) — inexpressible as one spec default — and **AWS rejects the flag for ECR_PUBLIC images**, so a spec-level CEL blocks the invalid combination at validation time instead of at the API. The old spec defaulted it to true, which made the simplest possible manifest (a public image) undeployable with auto-deploy implied.
- Source honesty: the current managed-runtime set validated as a closed CEL set (the live SDK carries `PYTHON_311`, `NODEJS_18`, `NODEJS_22` — absent from the old spec's documentation), `source_directory` for monorepos, health-check field names reconciled to the provider's `interval`/`timeout` vocabulary.
- `environment_secrets` annotated with `sensitive_exempt_reason` (values are Secrets Manager/SSM ARNs — references, never material) and removed from the secret-coverage accepted-gaps baseline.
- The hand-written Terraform contract replaced by the generator-owned `variables.tf` (all four kinds drift-guard enrolled); identity tags converged on the standard key set (the module carried a stale tag-key convention); provider pins on the family v6 floor.

### E2E (first-ever App Runner coverage)

- apprunner SDK verifiers for all four kinds, state-aware: a deleted service stays describable as DELETED; the versioned companions flip to INACTIVE — both treated as absent.
- Scenarios: companion minimal lanes (named for annotation-driven composition), the connector's two-AZ chain over the shared subnet + security-group fixtures, and two service lanes — a dependency-free minimal service (public sample image straight to HTTPS) and the composed shape (auto scaling + observability by reference + WAF association) with companions declared via the `planton.dev/e2e-prerequisites` annotation.
- **Live dual-engine E2E 8/8 green**: auto scaling configuration 24s/23s; observability configuration ~26s/lane; VPC connector chain 3m03s/3m10s; service minimal 4m30s/4m19s and composed 5m23s/6m09s. Zero-orphan sweep clean.

## Defects Found and Fixed Along the Way

- **App Runner status casing is inconsistent across AWS's own API**: auto scaling configurations return lowercase `"active"` while the SDK enum constant is `"ACTIVE"` (the Terraform provider itself uses lowercase string constants for exactly this resource). A strict enum comparison in the verifier passed compile and failed live; all App Runner verifiers now compare case-insensitively, with the quirk documented at the comparison site.
- **The service's legacy `Pulumi.yaml` carried `runtime.options.binary: main`** — pointing at a binary no clean checkout contains. Every offline gate passes; only a live `pulumi up` fails (`unable to find 'main' executable`). Fixed to the canonical plain `runtime: go`, and the update workflow's entrypoint-completeness gate now checks existing `Pulumi.yaml` files for this residue class rather than only checking existence.

## Deliberately Not Modeled (recorded with reasons)

- `aws_apprunner_connection` — unusable until a human completes the OAuth handshake in the console; composes by literal ARN.
- `aws_apprunner_vpc_ingress_connection` — the inbound private-access plane; composes against the exported `service_arn`.
- `aws_apprunner_default_auto_scaling_configuration_version` — an account-level singleton PUT.
- `aws_apprunner_deployment` — a one-shot operation, not infrastructure.
- Custom domains excluded from the live lane (requires a tenant-owned domain); proven by spec tests + the offline plan gate.

## Validation

Offline: spec tests x4 (go + Bazel), targeted Go + Pulumi entrypoint builds, `make build-go`, kind-map regen, variables.tf drift guard (all four enrolled), outputs conformance (+4 cases incl. the two-level repeated-message custom-domains shape), `tofu init`+`validate`+offline `plan` x4 from the hack manifests, `validate-refs --check`, `secret-coverage --check` (baseline entry removed, gate green), `validate-outputs` x4, every manifest CLI-validated (hack x4, presets x6, scenarios x5), site catalog regenerated (all four kinds listed). Live: all 8 dual-engine lanes green, account swept clean.
