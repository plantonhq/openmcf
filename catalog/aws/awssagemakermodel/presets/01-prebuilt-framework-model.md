# Prebuilt Framework Model

This preset serves a trained scikit-learn artifact on AWS's own
prebuilt framework image — the fastest path from a `model.tar.gz` in
S3 to a deployable model, with no container of your own to build.

## When to Use

- The first model for a classic ML workload (scikit-learn, XGBoost,
  and friends)
- Teams that train anywhere and only need AWS to serve

## What You Get

- A single-container model on the AWS-owned scikit-learn registry
  image for the region
- Compressed artifacts loaded from `modelDataUrl`, with the inference
  entry point named through the container environment

## Customize

- Swap the image for another prebuilt framework image (the registry
  path is per-region — keep it aligned with `region`)
- Point `executionRoleArn` at a role trusting `sagemaker.amazonaws.com`
  with S3 read on the artifact
- Any change replaces the model (models are immutable) — roll a new
  one and repoint the endpoint; keeping old versions costs nothing
