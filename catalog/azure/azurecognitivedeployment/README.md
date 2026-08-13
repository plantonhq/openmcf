# Overview

The **AzureCognitiveDeployment** component deploys a model onto an Azure AI services account -- which actual model an application calls (gpt-4o on deployment "chat", text-embedding-3-large on deployment "embeddings"), at which throughput class and capacity. The account (AzureCognitiveAccount, kind `OpenAI` or `AIServices`) owns the endpoint and keys; the deployment decides the model.

## Purpose

- **Models as declarative infrastructure**: which model, which version, which throughput class -- reviewed and versioned like everything else.
- **The scale knob in one field**: `sku.capacity` is thousands of tokens-per-minute on the pay-per-token SKUs (a rate limit, not idle cost) and PTUs on the provisioned SKUs; it updates in place.
- **Version policy made explicit**: track Azure's default version, upgrade on expiry, or pin with `NO_AUTO_UPGRADE` for compliance.
- **Responsible AI per deployment**: select an account-level content-filter policy by name.

## Key Features

- Full azurerm v5 surface: the 8-value SKU vocabulary (Standard, GlobalStandard, DataZone and Batch variants, ProvisionedManaged), model format/name/version, dynamic throttling, version upgrade options.
- Typed reference to the parent account -- chart-ready wiring.
- `rai_policy_name` honesty: omitted-when-unset so ARM's default policy applies without read drift.

## Use Cases

- **Chat and completion workloads**: gpt-4o / gpt-4o-mini on GlobalStandard.
- **Embedding pipelines**: text-embedding-3-large behind a rate-limited Standard deployment.
- **Reserved-capacity production**: ProvisionedManaged PTUs for latency-guaranteed serving.

## Future Enhancements

- Batch-SKU-specific presets as batch inference patterns settle.
