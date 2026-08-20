# Compliance Audit Trail

This preset creates the audit posture security reviews ask for: a
multi-region trail with tamper-evident log-file validation, all
management events, and both Insights anomaly engines.

## When to Use

- The account's primary audit trail (the first copy of management
  events is free)
- Compliance regimes that require multi-region, validated audit logs

## What You Get

- Every management API call in every region, delivered to S3 under
  `audit/AWSLogs/<account-id>/`
- Hourly SHA-256 digest files so log tampering is detectable
- Call-rate and error-rate anomaly detection (Insights bills per
  event)

## Customize

- Point `s3BucketName` at a bucket carrying the
  `cloudtrail.amazonaws.com` bucket policy (GetBucketAcl on the
  bucket, PutObject under it) — AWS rejects the trail without it
- Add a `cloudwatchLogs` block (group + role) to query events live
- Drop `insightTypes` to skip the per-event Insights bill

## Composing

```yaml
# The delivery bucket this preset expects:
apiVersion: aws.planton.dev/v1alpha1
kind: AwsS3Bucket
metadata:
  name: audit-logs-bucket
spec:
  region: <aws-region>
  policy:
    Version: "2012-10-17"
    Statement:
      - Sid: AWSCloudTrailAclCheck
        Effect: Allow
        Principal: { Service: cloudtrail.amazonaws.com }
        Action: s3:GetBucketAcl
        Resource: arn:aws:s3:::audit-logs-bucket
      - Sid: AWSCloudTrailWrite
        Effect: Allow
        Principal: { Service: cloudtrail.amazonaws.com }
        Action: s3:PutObject
        Resource: arn:aws:s3:::audit-logs-bucket/*
        Condition:
          StringEquals: { s3:x-amz-acl: bucket-owner-full-control }
```
