# Azure Monitor Autoscale Setting

Deploys an Azure Monitor autoscale setting -- the rule book that automatically adds and removes instances of one scalable target (a VM Scale Set, an App Service plan, ...) based on metric rules and schedules. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Autoscale setting** -- the scaling rule book bound to one target resource: capacity envelopes, metric rules, schedules, predictive autoscale, and scale-event notifications

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **A resource group** -- the setting lives in a referenced resource group.
- **A scalable target** -- the resource the setting controls, referenced by its ARM id output (e.g. an AzureVirtualMachineScaleSet's `scale_set_id` or an AzureServicePlan's `service_plan_id`).

### Azure Subscription

- **One autoscale setting per target** -- Azure rejects a second setting on the same resource at apply time.
- **The setting must live in the target's region** -- autoscale is evaluated regionally.
- **The target's tier must support autoscale** -- e.g. App Service plans need Standard tier or above (Basic plans are rejected at apply time).
- **The maximum instance count is also bounded by the subscription's core quota** -- the 0-1000 spec bounds mirror the provider's static validation only.
- **The autoscale setting object itself is free** -- you pay only for the instances it creates.

## Deploy

### Console

Open the deployment store, find **Azure Monitor Autoscale Setting**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CPU-Based Scale Set Autoscale** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f autoscale-setting.yaml
```

## After Deploy

Watch the portal's **Run history** tab on the autoscale setting: every evaluation and every scale action lands there, including the reason a rule did or did not fire. If the instance count never moves, check that the metric rule's `metric_resource_id` points at the resource actually emitting the metric, and that the capacity envelope leaves room to move.
