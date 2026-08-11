# CPU-Based Scale Set Autoscale

This preset deploys the classic elastic pool: a VM Scale Set that grows on sustained CPU pressure and shrinks conservatively when load drains, inside a 2-10 instance envelope.

## When to Use

- Any stateless VM Scale Set workload whose load tracks CPU (web tiers, API pools, render farms)
- As the starting shape for other metrics -- swap `metricName` for any metric the target emits

## Key Configuration Choices

- **Asymmetric rules** -- scale out at 75% over a 10-minute window with a 5-minute cooldown (react fast); scale in at 25% over 20 minutes with a 15-minute cooldown (avoid flapping). The wide gap between thresholds is deliberate: Azure's anti-flapping guard vetoes scale-ins that would immediately re-trigger a scale-out
- **`default: 3`** -- the metrics-outage posture: when metrics go dark, Azure raises the count to 3 if it is lower (never lowers it)
- **The metric source IS the target** -- both reference the same scale set; point `metricResourceId` elsewhere (e.g. a queue) for load-driven workers

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-region>` | The scale set's region (the setting must match it) | The scale set's overview page |
| `<your-vmss-arm-id>` | The scale set's ARM resource ID | `AzureVirtualMachineScaleSet` status outputs (`scale_set_id`) or the portal's Properties page |
