# Overview

The **AzureMachineLearningBatchEndpoint** component creates a batch endpoint on an Azure Machine Learning workspace -- the stable address batch scoring jobs are submitted to. Invoking the endpoint creates a JOB that runs a deployment's recipe over a large input (score a million records overnight) on pooled compute; nothing runs or bills while no job is active.

## Purpose

- **Batch inference as declarative infrastructure**: the job-submission address, its authentication, and its default-deployment routing -- reviewed and versioned like everything else.
- **A stable address over changing recipes**: deployments (model, compute, batching behavior) come and go behind the endpoint; the `defaultDeploymentName` pointer decides which one answers unrouted submissions.
- **Typed references end-to-end**: the workspace wires by reference, and deployments wire back to the endpoint by reference -- chart-ready.
- **Zero idle cost by design**: a batch endpoint is a routing object; compute provisions per job and scales back to zero.

## Key Features

- The full pinned ARM surface (api-version 2025-06-01): Microsoft Entra token authentication (the only mode the batch service accepts -- the spec validates it so misconfigurations cannot reach Azure), the default-deployment pointer, and the ARM property dictionary.
- azurerm carries NO resource for ML endpoints -- the modules write the raw ARM shape (azapi on Terraform, azure-native on Pulumi) at one pinned api-version, and the spec's validation rules carry the full contract burden.
- The service's realities recorded where they bite: batch auth is Entra-only (ARM's shared enum advertises key modes the service rejects), and the surface carries no public-network toggle (the workspace's network settings govern reachability).

## Use Cases

- **Scheduled batch scoring**: one endpoint per scoring application; a scheduler (or pipeline) submits jobs against the stable address.
- **Safe recipe rollout**: register a new deployment behind the same endpoint, test it by name, then move the default-deployment pointer.
- **Pipeline invocation**: with a pipeline-component deployment behind it, the endpoint becomes a stable trigger for a whole batch pipeline.

## Future Enhancements

- Migration to azurerm's native ML endpoint resources when the provider ships them (tracked upstream), with no manifest change.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
