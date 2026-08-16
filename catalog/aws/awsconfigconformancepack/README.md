<p align="center">
  <img src="logo.svg" alt="AWS Config Conformance Pack" width="80"/>
</p>

# AWS Config Conformance Pack

Manage an [AWS Config conformance pack](https://docs.aws.amazon.com/config/latest/developerguide/conformance-packs.html)
— a template bundle that deploys a set of Config rules (and optional
remediations) as one governed unit, at account or organization scope.

## What Gets Managed

- **The pack** (`metadata.name` is the pack name): the rule-bundle
  template, inline (`template_body`) or from S3 (`template_s3_uri`),
  with values for the template's parameters.
- **The scope**: this account (the default) or every account in the
  AWS Organization (`organization_scope: true`, with optional
  `excluded_accounts`).
- **Results delivery**: optionally an S3 bucket + prefix where Config
  stores the pack's evaluation results.

The pack's template CREATES its own rules — it references no
standalone Config rules (that is [AwsConfigRule](../awsconfigrule)'s
surface). Deploying a pack **requires a running Config recorder in the
region** ([AwsConfigRecorder](../awsconfigrecorder)); AWS rejects it
otherwise. Conformance packs carry no tags at the provider.

Destroying this component **deletes the pack and every rule it
created** (organization packs unwind from all member accounts, which
can take minutes).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
