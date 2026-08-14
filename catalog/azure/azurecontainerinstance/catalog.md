# Azure Container Instance

Deploys an Azure Container Instance container group -- serverless containers billed per second, with no cluster or VM to manage. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Container group** -- one or more containers sharing a lifecycle, network, and volumes: images, CPU/memory, probes, init containers, registry credentials, managed identity, Log Analytics diagnostics, custom DNS, network posture, and tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Private images** -- a registry credential entry per private registry (prefer the user-assigned-identity form over username/password).
- **Private posture** -- a subnet delegated to `Microsoft.ContainerInstance/containerGroups` (the AzureSubnet kind's delegations field).

### Azure Subscription

- **Nearly everything is fixed at creation** -- Azure applies only identity and tag changes in place; any other change replaces the group. Design the shape first.
- **Pick the restart policy for the workload** -- "Always" for services, "Never" for run-once jobs (the group shows Terminated when done), "OnFailure" for retried batch work.
- **Per-second billing** -- the group bills while it runs; a destroyed group costs nothing at rest.

## Deploy

### Console

Open the deployment store, find **Azure Container Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public Web Container** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f container-instance.yaml
```

## After Deploy

Use the `ip_address` or `fqdn` output where the group is consumed: a DNS record, an upstream proxy, or a health monitor. For jobs, the group's containers report their exit state in the portal and CLI; delete the group when the job's artifacts are collected (it bills only while it exists and runs).
