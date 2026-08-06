# AwsAppRunnerService -- Research Notes

## Provider Surface

Modeled on `aws_apprunner_service` plus two folded service-scoped satellites:

- `aws_apprunner_custom_domain_association` -- FOLDED as per-name `custom_domains` blocks (service-keyed, all-ForceNew, never FK-referenced; per-domain validation records exported).
- `aws_wafv2_web_acl_association` -- FOLDED as the optional `web_acl_arn` reference (pure glue; the protected resource points at the web ACL, matching how CloudFront models the same association as a distribution argument).

### `aws_apprunner_service` mapping

- `service_name` (Required, ForceNew) -- `metadata.name`, never a spec field.
- `source_configuration` (Required) -- the `image_source` XOR `code_source` arms; `auto_deployments_enabled`; `authentication_configuration` derived from whichever arm is active (access role for private ECR, connection for code repos -- the arms are mutually exclusive so the two never coexist).
  - The runtime settings (`port`, `start_command`, env var/secret maps) live at the spec top level and are routed into the active arm's `image_configuration` / `code_configuration_values` -- they configure the container either way, and exactly-one-of makes the routing unambiguous.
  - `source_code_version` supports only `type: BRANCH`, so the spec models `branch` directly instead of a one-value enum wrapper.
  - The runtime closed set mirrors the live SDK's `Runtime` enum: PYTHON_3, PYTHON_311, NODEJS_12/14/16/18/22, CORRETTO_8/11, GO_1, DOTNET_6, PHP_81, RUBY_31.
- `instance_configuration` -- `cpu`/`memory` (the provider's regex sets kept as closed CEL sets, both numeric and human-readable forms) + `instance_role_arn` ref.
- `health_check_configuration` -- provider-authentic field names (`interval`, `timeout`); path meaningful only for HTTP.
- `network_configuration` -- egress (DEFAULT vs VPC via the connector ref), ingress (`is_publicly_accessible`), `ip_address_type`. Both engines send the block explicitly so the deployed shape never depends on AWS-side defaults.
- `encryption_configuration.kms_key` (ForceNew) -- the KMS ref; the one replace-on-change field, stated in the spec comment.
- `observability_configuration` -- the spec models ONLY the configuration reference; presence sets `observability_enabled = true` in both engines. The provider's separate boolean is a drift trap (enabled-without-configuration is representable there); here it is unrepresentable.
- `auto_scaling_configuration_arn` (Optional+Computed) -- the configuration ref; unset falls back to the account default configuration.

## Design Decisions

- **The embedded side resources are gone.** The previous shape created a VPC connector from `subnet_ids`/`security_group_ids` and an auto scaling configuration from an `auto_scaling` block inside the service module -- with cross-engine naming divergence (TF suffixed `-vpc`/`-asc`, Pulumi used the bare name) and a Pulumi output contract that dropped two outputs when the side resources were absent. AWS designs both resources to be SHARED across services; embedding one per service forked what should be tuned in one place. Both are now first-class kinds referenced by ARN, and the service exports only what it owns.
- **`auto_deployments_enabled` is a plain bool (default false), sent explicitly by both engines.** AWS's own default is conditional (true for code repos, false for image repos), which a single spec default cannot honestly express -- and ECR_PUBLIC images REJECT the flag at create time. A spec-level CEL blocks the invalid combination at validation instead of at the API; explicit-send keeps the deployed value deterministic.
- **`image_identifier` stays a literal string** -- it carries a repository-plus-tag coordinate no single upstream output represents (`AwsEcsTaskDefinition.container.image` is the settled precedent).
- **`connection_arn` is a ref without a default kind** -- App Runner connections require an out-of-band OAuth handshake in the console to become usable, so there is no kind to point at (see the deferral ledger).
- **Custom domains fold, excluded from the live lane** -- association requires a domain the test tenant owns (the ACM validated-certificate class); the fold is proven by spec tests and the offline plan gate. The association waiter deliberately returns at `pending_certificate_dns_validation` -- full activation requires the DNS records the module does not manage, so the exported records ARE the deliverable.

## Deferral Ledger

- `aws_apprunner_connection` -- DEFER: created in seconds but unusable until a human completes the provider OAuth handshake in the console (the resource does not even wait past PENDING_HANDSHAKE). A kind whose deploy cannot reach a usable state without a manual console step fails the self-service test; composes by literal ARN.
- `aws_apprunner_vpc_ingress_connection` -- DEFER: the INBOUND private-access plane (PrivateLink into a private service via an `apprunner.requests` VPC endpoint). A separate product surface composing against the exported `service_arn` with zero rework.
- `aws_apprunner_default_auto_scaling_configuration_version` -- SKIP: an account-level regional singleton PUT (the account-settings class).
- `aws_apprunner_deployment` -- SKIP: a one-shot StartDeployment operation with a no-op delete -- an action, not infrastructure.

## Verification

- Spec tests cover both source arms, the runtime closed set, the ECR_PUBLIC auto-deploy rejection, health-check ranges, cpu/memory sets, custom domains, and the envelope.
- E2E: a dependency-free minimal lane (public sample image; App Runner provisions and health-checks in ~3-5 minutes, waiter ceiling 20) and a composed lane (auto scaling + observability configurations by reference, WAF association) whose companions install through the e2e-prerequisites annotation. Verifier: DescribeService with DELETED treated as absent.
