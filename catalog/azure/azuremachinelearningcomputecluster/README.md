# Overview

The **AzureMachineLearningComputeCluster** component creates an auto-scaling compute cluster on an Azure Machine Learning workspace -- the pool of VMs that training jobs and pipelines run on. The cluster grows and shrinks between its configured node bounds, and with a minimum of zero nodes it costs nothing while idle.

## Purpose

- **Training compute as declarative infrastructure**: VM size, priority, scale bounds, networking, and identity -- reviewed and versioned like everything else.
- **Scale-to-zero economics**: `minNodeCount: 0` makes the cluster free between jobs; the idle-duration knob balances warm-node latency against cost.
- **Typed references end-to-end**: the workspace, subnet, and user-assigned identities all wire by reference -- chart-ready.
- **Identity-first data access**: the cluster's managed identity is how jobs reach storage, Key Vault, and the container registry without embedded credentials.

## Key Features

- Full azurerm v5 surface: VM size and priority (dedicated or evictable low-priority at a deep discount), required scale settings with ISO-8601 idle duration, system/user-assigned identity, per-node SSH admin account with the provider's at-least-one-credential contract validated at manifest time, VNet placement, node public-IP and local-auth toggles.
- The provider's update contract recorded where it bites: only identity, scale settings, and tags update in place -- everything else replaces the cluster.
- Cross-region capability modeled honestly: the cluster's `region` is where NODES run and may differ from the workspace's region (the only ML compute with this ability); ARM reports the cluster envelope at the workspace's region.

## Use Cases

- **Shared training pool**: one scale-to-zero CPU cluster the whole team submits jobs to.
- **GPU training**: a dedicated GPU-family cluster sized to the team's regional quota.
- **Cost-optimized batch**: a low-priority (spot-class) cluster for checkpointed, fault-tolerant workloads.

## Future Enhancements

- Managed online/batch endpoint kinds for model serving as their contracts land in the catalog.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
