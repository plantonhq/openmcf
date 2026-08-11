# Model Scoring Recipe

This preset creates the classic batch inference recipe: a registered model scored on a compute-cluster pool, with the batching dials tuned explicitly and an appended predictions file as output.

## When to Use

- Nightly / scheduled scoring of large datasets with a registered model
- The first deployment behind a new batch endpoint
- Any recipe where "job succeeded" must mean "records actually scored" (the explicit error threshold)

## Key Configuration Choices

- **`model.id.assetId`** -- a registered model version's ARM ID (register with `az ml model create`; registering assets is a data-science workflow step). MLflow models need no code configuration -- the service generates scoring code.
- **`computeId`** -- reference an AzureMachineLearningComputeCluster; a min-0 pool means the recipe costs nothing between jobs. `resources.instanceCount: 4` asks each job for four nodes -- keep it at or under the pool's max.
- **`errorThreshold: 100`** -- an explicit abort threshold instead of the service default (-1, ignore all failures); set it to 0 to abort on the first failure.
- **`retrySettings.timeout: PT1M`** -- the per-mini-batch invocation timeout; the first suspect when large mini-batches "randomly" fail and retry.

## After Deployment

Set the endpoint's `defaultDeploymentName` to `production` (or invoke this recipe by name), then submit jobs against the endpoint's `scoring_uri` with an Entra token.
