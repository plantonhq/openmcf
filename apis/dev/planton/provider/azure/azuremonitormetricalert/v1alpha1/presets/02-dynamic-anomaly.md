# Dynamic Anomaly Alert

This preset creates a machine-learning metric alert: Azure learns the metric's normal band (including daily and weekly seasonality) and fires when the value deviates -- in either direction. No fixed threshold to guess or re-tune as the workload grows.

## When to Use

- Metrics with rhythm (traffic, transactions, queue depth) where any fixed line is wrong half the day
- Detecting drops (an upstream outage silences traffic) as well as spikes (abuse, retry storms)

## Key Configuration Choices

- **GREATER_OR_LESS_THAN** -- both deviation directions alert; narrow to GREATER_THAN or LESS_THAN when only one direction matters
- **MEDIUM sensitivity** -- balanced band width; HIGH pages on small wobbles, LOW only on major departures
- **3 of 4 periods** (`evaluationFailureCount: 3` of `evaluationTotalCount: 4`) -- the flap damper: a single odd window does not page
- **Learning period** -- new rules need history to learn from; expect noisier behavior on brand-new resources, and use `ignoreDataBefore` after a regime change

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the rule | `AzureResourceGroup` status outputs |
| `my-app-storage` | The monitored resource | Any kind's status outputs (`*_id`) |
| `my-platform-oncall` | The action group to notify | `AzureMonitorActionGroup` status outputs |

## Related Presets

- **01-static-threshold** -- Fixed lines for SLO-backed conditions
- **03-webtest-availability** -- Endpoint availability from outside-in
