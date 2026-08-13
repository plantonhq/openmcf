# Overview

The **AzureMachineLearningOnlineDeployment** component creates a MANAGED online deployment on an Azure Machine Learning online endpoint -- a running copy of a model behind the endpoint's address, on Azure-managed VMs with health probes, request limits, and optional model data collection. The endpoint's traffic map routes scoring requests to deployments by name, which is how models roll out blue/green.

## Purpose

- **The model-serving runtime as declarative infrastructure**: the model, environment, instance size and count, probes, and data collection -- reviewed and versioned like everything else.
- **Blue/green by construction**: deployments are the disposable layer behind a stable endpoint; ship a new model as a new deployment and move the endpoint's traffic map.
- **Typed references end-to-end**: the deployment wires its endpoint by reference -- chart-ready.
- **Monitoring-ready**: model data collection captures scoring inputs and outputs to workspace storage for drift detection.

## Key Features

- The full pinned ARM surface (api-version 2025-06-01) for the managed compute type: model / environment / scoring-code references, instance type and count (the scale dial, carried as the ARM SKU capacity), all three health probes with ISO-8601 durations validated at manifest time, request settings, egress control, and per-collection model data capture.
- azurerm carries NO resource for ML deployments -- the modules write the raw ARM shape (azapi on Terraform, azure-native on Pulumi) at one pinned api-version, and the spec's validation rules carry the full contract burden.
- Honest scope, recorded: the Kubernetes and AzureMLCompute variants are a recorded deferral (they need attached compute whose supported story does not exist yet), and managed deployments have no scale-to-zero.

## Use Cases

- **Serve a registered model**: reference a model version and environment from the workspace registry on managed VMs.
- **Blue/green model rollout**: run the new model as `green` beside `blue` and shift the endpoint's traffic map in steps.
- **Custom containers**: bring an image that embeds its model and server, with probes and mount paths under your control.

## Future Enhancements

- Kubernetes-attached deployment variants when a supported compute-attach story lands.
- Migration to azurerm's native ML deployment resources when the provider ships them (tracked upstream), with no manifest change.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
