# AzureMachineLearningBatchDeployment Terraform Module

## Overview

Creates a batch deployment on an Azure Machine Learning batch endpoint -- the job recipe (model, compute, batching behavior) behind the endpoint's address. Nothing runs or bills at create time; each endpoint invocation materializes a job from the recipe.

## Resources Created

- `azapi_resource` -- the deployment, written at the pinned raw-ARM shape `Microsoft.MachineLearningServices/workspaces/batchEndpoints/deployments@2025-06-01`, an ARM child of the endpoint

**Why azapi, not azurerm:** azurerm carries NO resource for ML batch deployments (its endpoint draft is tracked at hashicorp/terraform-provider-azurerm#32011). The azapi provider is pinned EXACT (2.11.0) and the api-version is pinned in the resource type -- never `latest`. When azurerm ships native resources, this module migrates azapi → native in the next minor release (state move / re-import). The kind's spec carries the full validation burden: azapi has no provider-side schema, so every ARM contract is a manifest-time rule or a documented apply-time boundary.

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningBatchDeploymentSpec fields; the endpoint and compute references arrive as resolved literal ARM IDs

## Outputs

- `batch_deployment_id` -- the deployment's full ARM ID
- `batch_deployment_name` -- the key the endpoint's default-deployment pointer routes by

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Everything updates in place** via full PUT: ARM flags nothing immutable on this surface (unlike the online deployment's ForceNew instance type); name, region, and endpoint replace the deployment.
- **The model reference is a discriminated union**: the spec's variant block (id / data_path / output_path) IS the ARM `referenceType` -- the locals build the discriminator from which block is set, guaranteed single by the spec's exactly-one rule.
- **The PipelineComponent deployment type** is written only when the spec's `pipeline_component` block is present; absent means ARM's default Model type (the enum's Model value has no concrete ARM shape of its own).
- **No envelope `sku`**: batch scale lives on `resources.instanceCount` per job, not on an autoscaling SKU (the online deployment's dial) -- the envelope's sku is deliberately absent.
- **`resources.properties` is deliberately not modeled**: an untyped ARM bag (arbitrary JSON objects, no documented keys) a string-keyed contract cannot carry faithfully -- a recorded exclusion on the spec message.
- **Unset optionals are OMITTED from the body** so ARM's own defaults apply: miniBatchSize 10, maxConcurrencyPerInstance 1, errorThreshold -1, loggingLevel Info, outputAction AppendRow, outputFileName predictions.csv, retries 3/PT30S.
