# Batch Scoring Endpoint

This preset creates the everyday batch endpoint: the minimal routing object batch deployments attach to, with Microsoft Entra authentication applied by default (the only mode the batch service accepts).

## When to Use

- The first endpoint a batch scoring application needs
- The stable address a scheduler or pipeline submits jobs against
- The parent object for one or more batch deployment recipes

## Key Configuration Choices

- **`authMode` unset** -- the platform applies `AADToken`; the batch service rejects every other mode, so there is nothing to decide and no keys to manage. Submitting principals (users or service principals) need rights to create jobs in the workspace.
- **No identity** -- batch jobs run under the submitter's Entra token plus the compute pool's managed identity; the endpoint's own identity is optional and this preset omits it.
- **No default-deployment pointer yet** -- set `defaultDeploymentName` after the first deployment attaches (it updates in place); until then, submissions name their deployment explicitly.

## After Deployment

Attach an **Azure Machine Learning Batch Deployment** wiring `endpointId` by reference, then set `defaultDeploymentName` to its name and submit jobs to the endpoint's `scoring_uri` output with an Entra token.
