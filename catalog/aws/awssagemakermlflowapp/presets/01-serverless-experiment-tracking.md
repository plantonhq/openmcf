# Serverless Experiment Tracking

This preset is MLflow with no meter running: a serverless MLflow 3.x
app storing artifacts in your S3 bucket, billed per use and $0 when
idle — the successor to the hourly-billed tracking server.

## When to Use

- New MLflow deployments — per-use billing beats an always-on server
  for all but sustained heavy tracking
- Teams that track intermittently and shouldn't pay for idle capacity

## What You Get

- A serverless MLflow 3.x app with `app_arn` as its identity in the
  outputs
- Model files and run outputs stored in your S3 bucket through the
  app's IAM role
- Nothing to size, nothing to resize — capacity is AWS's problem

## Customize

- Get `roleArn` right the first time — it is the ONE field whose
  change replaces the app; everything else updates in place
- Add `defaultDomainIds` so Studio users in those SageMaker domains
  track to this app automatically
- Point `roleArn` at a composed `AwsIamRole` with `valueFrom` instead
  of a literal ARN
