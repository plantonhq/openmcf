# AwsAppRunnerAutoScalingConfiguration -- Research Notes

## Provider Surface

Modeled 1:1 on `aws_apprunner_auto_scaling_configuration_version`:

- `auto_scaling_configuration_name` (Required, ForceNew) -- derived from `metadata.name`, never a spec field.
- `max_concurrency` (Optional, ForceNew, default 100, 1-200) -- spec `max_concurrency` with the provider's own range as buf.validate bounds.
- `max_size` (Optional, ForceNew, default 25, >=1) -- spec `max_size`. The provider validates only the lower bound; the 25 ceiling is an adjustable service quota, so the spec deliberately does NOT freeze it (documented instead).
- `min_size` (Optional, ForceNew, default 1, >=1) -- spec `min_size`.
- Computed: `arn`, `auto_scaling_configuration_revision`, `latest`, `is_default`, `has_associated_service`, `status`.

Every settable attribute is ForceNew because AWS versions the resource: an apply that changes a value destroys the state entry and registers the next revision under the same name. The kind embraces this as the composition mechanism -- the exported ARN carries the revision, so referencing services roll on their next deployment.

## Design Decisions

- **max >= min as CEL** -- the one cross-field invariant, caught at validation instead of at the API.
- **`latest` exported, `is_default`/`has_associated_service` not** -- `latest` tells a reader whether this resource still owns the name's newest revision; the other two are account-state flags that describe things outside this resource's control (the account default setter is deliberately unmodeled; association is visible from the service side).
- **`status` not exported** -- unlike long-lifecycle resources, a configuration is `active` from the moment create returns; the E2E verifier checks it via the API instead.

## Deferral Ledger

- `aws_apprunner_default_auto_scaling_configuration_version` -- SKIP: an account-level regional singleton PUT (Create/Update both call `UpdateDefaultAutoScalingConfiguration`, Delete is a no-op). The account-settings class, not a graph resource.

## Verification

- Spec tests cover the range bounds and the max>=min CEL (happy + failure paths).
- E2E: register -> DescribeAutoScalingConfiguration (active) -> destroy -> verify the revision flipped to inactive (deletion never hard-removes revisions).
