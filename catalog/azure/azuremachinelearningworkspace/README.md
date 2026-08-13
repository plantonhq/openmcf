# Overview

The **AzureMachineLearningWorkspace** component deploys an Azure Machine Learning workspace -- the central home a data-science team works in: experiments, models, endpoints, datastores, and compute all live on a workspace. The workspace is a thin coordination object over three required companion services (a storage account, a key vault, and an application-insights component) plus an optional container registry, all fixed at creation.

## Purpose

- **The ML platform's foundation block**: datastores, compute clusters and instances, and model endpoints are all children of a workspace.
- **Companion wiring as typed references**: storage, key vault, application insights, and container registry wire by reference -- chart-ready, no copy-pasted ARM IDs.
- **Managed network isolation made declarative**: the isolation mode plus named outbound rules (FQDN, private-endpoint, service-tag) in one spec.
- **Feature stores as a flavor**: `kind: FEATURE_STORE` with the feature-store block, validated both ways at manifest time.

## Key Features

- Full azurerm v5 surface: identity (system/user-assigned), customer-managed-key encryption with service-side CMK, managed virtual network with all three outbound-rule types composed as children, serverless compute placement, high-business-impact mode, v1 legacy mode.
- THIRTEEN provider code contracts front-loaded as CEL validations -- kind/feature-store pairing, encryption pairings, the serverless no-public-IP rule, cross-type outbound-rule name uniqueness (the three rule lists share one ARM collection).
- Name-keyed outbound-rule ID maps as outputs -- charts reference individual rules.

## Use Cases

- **Classic ML platform**: workspace + compute cluster + datastores for training pipelines.
- **Locked-down enterprise ML**: private workspace with managed-network isolation and approved-outbound rules only.
- **Feature store**: the FEATURE_STORE flavor backing online/offline feature serving.

## Future Enhancements

- Managed online/batch endpoint kinds (the model-serving layer) join the family next.
