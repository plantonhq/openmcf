# AWS GuardDuty

The region's threat detection: one detector with protection plans
(S3, EKS, runtime, RDS, Lambda, AI), noise-control filters, trusted
and threat IP lists, findings export to S3, and organization-wide
member administration.

## What Gets Managed

- The detector (one per account per region) and its finding
  re-publish frequency.
- Protection-plan features with agent-management sub-toggles.
- Finding filters, trusted IP lists, and threat intel lists.
- Findings export to S3 under a KMS key.
- Organization delegation, auto-enrollment, org-wide features, and
  member accounts — or the member-side invitation accept.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with GuardDuty permissions.

### AWS Account

- No pre-existing detector in the region (hand-enabled or
  Organizations auto-enabled) — AWS allows exactly one.
- For findings export: a bucket whose policy grants
  `guardduty.amazonaws.com`
  ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket)) and a KMS key whose
  policy grants it `kms:GenerateDataKey`
  ([AWS KMS Key](/cloud-catalog/aws-kms-key)).
- For trusted/threat lists: the list files in S3, readable by
  GuardDuty.

## Deploy

### Console

Create the resource from the AWS catalog, pick the protection plans,
and deploy.

### CLI

```bash
planton apply -f guardduty.yaml
```

## After Deploy

- Findings appear in the GuardDuty console as threats are detected;
  updated findings re-publish at the configured frequency.
- Exported findings land in the S3 destination within minutes.
- Destroy deletes the detector, its findings, and every satellite.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
