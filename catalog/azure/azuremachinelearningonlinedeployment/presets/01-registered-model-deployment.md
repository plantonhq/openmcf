# Registered-Model Deployment

This preset serves a registered model version from the workspace on one managed instance -- the everyday deployment an MLflow model needs, since the service infers scoring code and environment for MLflow flavors.

## When to Use

- The first deployment behind a new online endpoint
- MLflow-registered models (no scoring script or environment needed)
- Any rollout where a new model version ships as a new deployment

## Key Configuration Choices

- **`name: blue`** -- the endpoint's traffic map routes by this key; ship the next model version as `green` and shift the map in steps.
- **`instanceType` / `instanceCount`** -- the type is fixed at creation, the count updates in place (the one change the service applies without rolling containers). There is no scale-to-zero; one instance is the honest floor.
- **`model`** -- the ARM ID of a registered model version. Non-MLflow models add `codeConfiguration` (scoring script) and `environmentId`.

## After Deployment

Point the endpoint's `traffic` map at `blue: 100`, then POST to the endpoint's `scoring_uri`. Watch the deployment's logs (`az ml online-deployment get-logs`) during the first requests.
