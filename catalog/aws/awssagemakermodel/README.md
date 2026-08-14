<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Model" width="80"/>
</p>

# AWS SageMaker Model

Create and manage [Amazon SageMaker AI models](https://docs.aws.amazon.com/sagemaker/latest/dg/realtime-endpoints-deploy-models.html)
— the immutable serving definition (container image, model artifacts,
execution role, networking) that endpoints, batch transform jobs, and
inference components deploy.

## What Gets Created

- **A model** that is either a single serving container
  (`primary_container`) or an **inference pipeline** of 2–15
  containers (`containers`) invoked in sequence (`Serial`) or
  addressed directly by hostname (`Direct`) — exactly one of the two
  forms.
- Each container serves an ECR image or a registered model package,
  with artifacts from S3 as a compressed `model_data_url` or an
  uncompressed `model_data_source` (with gated-model EULA acceptance),
  plus optional adapter channels, MultiModel mode, and private-registry
  image configuration.
- Optional: VPC attachment (`vpc_config`, 1–16 subnets and 1–5
  security groups) and full network isolation
  (`enable_network_isolation`).

## Everything Is Create-Time Only

SageMaker models are immutable — any spec change replaces the model
(AWS's own contract: roll a new model and repoint the endpoint).
Models themselves are free; only the endpoints deploying them bill.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
