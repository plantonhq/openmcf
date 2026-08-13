# Overview

The **AzureMachineLearningBatchDeployment** component creates a batch deployment on an Azure Machine Learning batch endpoint -- the job RECIPE behind the endpoint's address: which model to run, on which compute pool, and how to split, retry, and collect the work. Nothing runs or bills at create time; each endpoint invocation materializes a job from the recipe.

## Purpose

- **Scoring recipes as declarative infrastructure**: model reference, compute target, mini-batch sizing, retry policy, and output shape -- reviewed and versioned like everything else.
- **Typed references end-to-end**: the endpoint wires by reference, and the compute pool wires by reference to the compute-cluster kind -- chart-ready.
- **Model OR pipeline**: the same kind expresses a model scoring recipe and a pipeline-component recipe (a whole batch pipeline behind the endpoint's stable address).
- **Everything updates in place**: ARM flags nothing immutable on this surface -- tune mini-batch size, retries, or the model reference without replacing the deployment.

## Key Features

- The full pinned ARM surface (api-version 2025-06-01): the three-arm model reference union (registered asset ID, datastore path, job output), code configuration, environment, per-job compute sizing, mini-batch/concurrency/error-threshold dials with the service's defaults documented, retry settings on ISO-8601 durations, output action/file, logging level, and the pipeline-component deployment type.
- azurerm carries NO resource for ML deployments -- the modules write the raw ARM shape (azapi on Terraform, azure-native on Pulumi) at one pinned api-version, and the spec's validation rules carry the full contract burden.
- Honest boundaries recorded where they bite: ARM's untyped `resources.properties` bag is deliberately not modeled (a recorded exclusion), and pipeline-component recipes carry their own steps (the model/code/environment fields describe the Model type).

## Use Cases

- **Nightly scoring**: a registered model + a scale-to-zero compute cluster + `AppendRow` output -- the classic batch inference recipe.
- **Blue/green recipe rollout**: two deployments behind one endpoint; invoke the new one by name until it earns the default pointer.
- **Pipeline-as-a-service**: a pipeline-component deployment turns a registered pipeline into a stable, invocable address.

## Future Enhancements

- Migration to azurerm's native ML deployment resources when the provider ships them (tracked upstream), with no manifest change.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
