# Overview

The **AzureMachineLearningOnlineEndpoint** component creates a managed online endpoint on an Azure Machine Learning workspace -- the stable HTTPS address applications call to score against deployed models. The endpoint owns authentication and a traffic dial that splits requests across its deployments, so models roll out blue/green and shadow-test against production traffic without callers noticing.

## Purpose

- **Model serving as declarative infrastructure**: the scoring address, its authentication mode, its identity, and its traffic split -- reviewed and versioned like everything else.
- **Blue/green rollouts by design**: the `traffic` map routes whole-percent shares of requests to deployments by name, and `mirrorTraffic` shadow-tests a candidate against live traffic without returning its answers.
- **Typed references end-to-end**: the workspace wires by reference, and deployments wire back to the endpoint by reference -- chart-ready.
- **Identity-first artifact access**: the endpoint's managed identity is how its deployments pull container images and model artifacts without embedded credentials.

## Key Features

- The full pinned ARM surface (api-version 2025-06-01): all three auth modes (static keys, Azure ML tokens, Entra tokens), traffic and mirror-traffic maps with the service's percentage bounds validated at manifest time, public-network toggle, bring-your-own initial auth keys from a secret store, and the ARM property dictionary.
- azurerm carries NO resource for ML endpoints -- the modules write the raw ARM shape (azapi on Terraform, azure-native on Pulumi) at one pinned api-version, and the spec's validation rules carry the full contract burden.
- The service's realities recorded where they bite: endpoint names are reserved region-wide per subscription, and ARM never returns key values on reads (retrieval is the separate listKeys action).

## Use Cases

- **Real-time inference**: one endpoint per model-serving application, with deployments behind it carrying model versions.
- **Blue/green model rollout**: run old and new deployments side by side and move the traffic map in steps.
- **Shadow evaluation**: mirror a fraction of production traffic to a candidate deployment before it ever answers a caller.

## Future Enhancements

- Batch endpoint kinds for asynchronous, job-based scoring as their contracts land in the catalog.
- Migration to azurerm's native ML endpoint resources when the provider ships them (tracked upstream), with no manifest change.
