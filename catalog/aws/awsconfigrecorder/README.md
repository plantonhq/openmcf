<p align="center">
  <img src="logo.svg" alt="AWS Config Recorder" width="80"/>
</p>

# AWS Config Recorder

Manage the region's [AWS Config recording posture](https://docs.aws.amazon.com/config/latest/developerguide/WhatIsConfig.html):
the configuration recorder, its S3/SNS delivery channel, and the
retention window for recorded configuration items.

This is a **region singleton**: AWS allows exactly one recorder and
one delivery channel per region, both named `default` by AWS
convention — the names are not an identity you choose, so
`metadata.name` never reaches AWS. Deploy at most one instance per
region.

## What Gets Managed

- **The recorder**: its service role (trusting
  `config.amazonaws.com`), what it records (everything, an inclusion
  list, or everything-except an exclusion list), how often
  (continuous or daily, with per-type overrides), and whether it is
  RUNNING (the folded start/stop toggle).
- **The delivery channel**: the history bucket (which must carry the
  `config.amazonaws.com` bucket policy — see
  [AwsS3Bucket](../awss3bucket) `spec.policy`), optional KMS
  encryption, SNS notification, and snapshot frequency.
- **Retention**: how many days recorded configuration items stay
  queryable (30–2557).

Destroying this component is a **real delete with ordering**: the
recorder stops first, then recorder + channel + retention are removed
(already-recorded items age out on their own). Recording bills per
configuration item — scope the recording group deliberately.

Compliance rules that evaluate what this recorder captures are their
own component: [AwsConfigRule](../awsconfigrule).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
