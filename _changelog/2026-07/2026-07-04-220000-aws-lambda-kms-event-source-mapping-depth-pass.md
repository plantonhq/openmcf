# AWS Lambda Rebuild, Event Source Mapping Forge, and KMS Key Depth Pass

**Date:** 2026-07-04  
**Scope:** Components #21–#23 — `AwsLambda`, `AwsLambdaEventSourceMapping` (enum 352), `AwsKmsKey`

## Summary

Three related serverless and encryption kinds brought to the 90/10 bar in one
session: Lambda rebuilt from a never-deployable legacy module to the full
`aws_lambda_function` surface with folded satellites; event source mappings
split out as a first-class composable kind; KMS keys expanded from a six-field
stub to the complete customer-managed key surface with multi-alias support.

## Product impact

- **Lambda manifests change shape** — `function_name` and `codeSourceType` are
  gone; naming is `metadata.name`; code is exactly one of `s3` or `image_uri`.
  Embedded IAM roles, log groups, and invoke permissions are no longer auto-
  created — compose `AwsIamRole`, `AwsCloudwatchLogGroup`, and explicit
  `invoke_permissions` instead.
- **Event-driven pipelines gain a composable edge** — SQS/Kinesis/DynamoDB/
  MSK/MQ/DocumentDB → Lambda wiring is now `AwsLambdaEventSourceMapping`, not
  buried inside the function spec.
- **KMS keys expose the real provider surface** — `key_spec` is a provider
  string (13 values), `policy` is a JSON document, rotation is honest, aliases
  are a repeated list. Thirty-nine downstream consumers of `key_arn` are unchanged.

## Technical highlights

- Registry-driven E2E test-name discovery in `pkg/e2e/profile/discover.go`
  eliminates per-kind PascalCase table maintenance.
- Complete E2E artifacts for all three kinds (verifiers, profiles, scenarios,
  dual-engine test entrypoints); live lanes deferred pending credential refresh.
- `function_arn` output kept stable — 13 FK references across HTTP API Gateway,
  Cognito triggers, and Firehose verified safe.

## Breaking changes

| Kind | Change |
|------|--------|
| AwsLambda | Drop `function_name`; retire embedded role/log-group/permission creation |
| AwsKmsKey | `key_spec` enum → strings; `alias_name` → `aliases[]`; drop `rotation_enabled` output |
| charts/aws/serverless-api | Breaks on new Lambda shape — charts wave, not this session |

## Validation

Offline gate green: spec tests, drift guard, outputs conformance, tofu validate
(all three TF modules), Pulumi bazel builds. Live E2E not run (`InvalidClientTokenId`).
