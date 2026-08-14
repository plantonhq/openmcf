# AWS Config Recorder

The region's AWS Config recording posture — what changed, when, and
what it looked like before: the configuration recorder, its S3
delivery channel, and the retention window, managed as one regional
singleton.

## What Gets Managed

- The recorder: service role, recording scope (all / inclusion /
  exclusion), frequency (continuous or daily with overrides), and the
  running state.
- The delivery channel: history bucket, optional KMS encryption and
  SNS notices, snapshot frequency.
- Configuration-item retention (30–2557 days).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Config permissions.

### AWS Account

- A service role trusting `config.amazonaws.com`
  ([AWS IAM Role](/cloud-catalog/aws-iam-role) with the managed
  `AWS_ConfigRole` policy plus write access to the history bucket).
- A history bucket carrying the `config.amazonaws.com` bucket policy
  ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the role and bucket,
scope the recording group, and deploy.

### CLI

```bash
planton apply -f config-recorder.yaml
```

## After Deploy

- Configuration items land in the bucket and back every
  [AWS Config Rule](/cloud-catalog/aws-config-rule) evaluation in the
  region.
- Recording bills per configuration item — the recording group's
  scope is the cost lever.
- Destroy stops the recorder, then removes recorder, channel, and
  retention (recorded history ages out on its own).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
