# Azure Machine Learning Batch Deployment -- Operational Guide

Judgment that saves real time when running batch deployments. The field reference lives in the API Explorer; this is the operational layer above it.

## A deployment is a recipe -- cost lives elsewhere

Creating, updating, or keeping a batch deployment costs nothing: compute provisions when a JOB runs and scales back down after. This inverts the online deployment's economics (a standing VM from create to destroy). The corollary: there is no reason to delete batch deployments aggressively -- keep the last known-good recipe around as an instant rollback target behind the same endpoint.

## Size against the compute pool, not against hope

`resources.instanceCount` is how many nodes each JOB requests, and ARM accepts any positive number -- the real ceiling is the compute cluster's `max_node_count`, checked only at job time. A recipe asking for 8 nodes on a 4-node pool deploys cleanly and then every job queues half its work. When jobs sit in "queued" longer than they run, compare these two numbers first.

## The batching dials multiply -- do the arithmetic once

Total parallelism = instanceCount × maxConcurrencyPerInstance, each invocation receiving miniBatchSize units of input. The defaults (1 node, 1 concurrent, 10 units) are deliberately timid. Tune in this order: raise instanceCount to the pool's ceiling, raise maxConcurrencyPerInstance only if the scoring code is I/O-bound (it shares the node's memory), and set miniBatchSize so one invocation stays comfortably inside the per-invocation timeout (`retrySettings.timeout`, default 30 seconds -- the FIRST suspect when large mini-batches "randomly" fail and retry).

## errorThreshold's default keeps jobs alive -- decide if you want that

The default -1 means IGNORE every scoring failure and finish the job: right for exploratory scoring, wrong for pipelines that treat "job succeeded" as "every record scored." Set an explicit threshold (0 aborts on the first failure) when downstream consumers assume completeness, and read the failure counts from job outputs either way.

## Model updates roll forward in place -- rollback is a second deployment

Everything on this surface updates via full PUT, so pointing the recipe at a new model version is a one-field update -- and the previous version is gone from the recipe. For anything that matters, run model versions as SEPARATE deployments behind the endpoint (`churn-v3`, `churn-v4`), invoke the new one by name until it earns trust, then move the endpoint's default pointer. In-place model updates are for recipes nothing routes to yet.

## Pipeline recipes carry their own compute setting

A `pipelineComponent` recipe runs a registered pipeline per job -- and the pipeline's steps find their compute through the `settings` map (`default_compute`), NOT through the deployment's `computeId`/`resources` (those describe model scoring). A pipeline recipe that omits `default_compute` fails at job time with a compute-resolution error that looks like a permissions problem; it is not.
