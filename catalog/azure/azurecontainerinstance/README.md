# Overview

The **AzureContainerInstance** component deploys an Azure Container Instance container group -- serverless containers billed per second: hand Azure an image plus CPU and memory, and it runs, with no cluster or VM to manage.

## Purpose

- **The fastest path from image to running container**: no AKS cluster to operate, no VM to patch -- one manifest, one group, running in seconds.
- **Run-once jobs and burst workloads**: `restartPolicy: Never` turns the group into a job; per-second billing means a group that runs ten minutes costs ten minutes.
- **Sidecar groups without a cluster**: several containers share one lifecycle, network namespace, and volumes -- the pod pattern, serverless.

## Key Features

- Full azurerm v5 surface: multi-container groups with liveness/readiness probes, one-shot init containers, all four volume forms (Azure File share, empty scratch dir, git checkout, inline secret files), private-registry credentials (identity or username/password), managed identity, Log Analytics diagnostics, custom DNS, public/private/IP-less postures, Spot priority, zones, and customer-managed-key encryption.
- The volume forms are a first-class choice: each volume takes exactly ONE of azure_file / empty_dir / git_repo / secret, validated before anything reaches Azure -- the provider itself only rejects a wrong mix at apply time.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, `subnet_id` to AzureSubnet, registry credentials to AzureContainerRegistry and AzureUserAssignedIdentity, diagnostics to AzureLogAnalyticsWorkspace (customer GUID + shared key), volume storage to AzureStorageAccount and AzureStorageShare, and encryption to AzureKeyVaultKey; the `ip_address` and `fqdn` outputs are what DNS records and upstream proxies consume.
- Secure by default: every secret-bearing field (secure environment variables, volume storage key and secret files, registry password, workspace key) is marked sensitive, and the docs record that Azure never returns any of them on reads.

## Use Cases

- **Scheduled and one-off jobs**: data loads, migrations, report generation -- `restartPolicy: Never` with per-second billing.
- **Simple always-on services**: a small API or webhook receiver with a public IP and DNS label, no cluster overhead.
- **VNet-internal workers**: a private group inside a delegated subnet processing internal queues.
- **Build and CI burst capacity**: spin up, run, destroy -- the group is the unit.

## Future Enhancements

- GPU containers were retired from the provider's surface at v5 (Azure retired the underlying SKUs) -- if azurerm regains a GPU surface, it lands here.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
