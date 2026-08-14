# Overview

The **AzureMonitorAutoscaleSetting** component deploys an Azure Monitor autoscale setting -- the rule book that automatically adds and removes instances of one scalable target (a Virtual Machine Scale Set, an App Service plan, and other capacity-bearing resources) based on metric rules and schedules.

## Purpose

- **Capacity that follows load**: metric rules ("CPU average over 10 minutes above 75% -> add 2 instances, cool down 5 minutes") move the instance count inside a min/max envelope, so the service absorbs peaks without paying for them around the clock.
- **Capacity that follows the calendar**: recurrence profiles set known shapes ("10 instances on weekday business hours, 2 overnight"), and fixed-date profiles handle one-off events (launch day, sale weekend) -- with metric rules still active inside each profile's envelope.
- **One object, whole story**: the envelope, the rules, the schedules, predictive autoscale, and scale-event notifications (email/webhooks) live in one resource with one lifecycle.

## Key Features

- Full azurerm v5 surface: up to 20 profiles with capacity envelopes, up to 10 metric rules each (all statistics, aggregations, and comparison operators; dimension filters; per-instance metric division), fixed-date and weekly-recurrence schedules with Azure's full timezone vocabulary, predictive autoscale (ForecastOnly/Enabled with look-ahead), email and webhook notifications.
- The provider's contracts front-loaded as validation: a profile carries at most one schedule (fixed-date XOR recurrence); a notification block must configure at least one channel; capacity bounds, schedule fields, and every vocabulary validate at manifest time.
- Chart-ready: the scale target and metric sources are references -- point them at any kind's `*_id` output (a scale set, a service plan, or a queue whose depth drives a worker pool).

## Use Cases

- **VM Scale Set on CPU**: the classic elastic pool -- scale out on sustained CPU, scale in conservatively with a longer cooldown.
- **App Service plan by schedule**: business-hours capacity weekdays 8-18, minimal capacity otherwise, with a CPU rule guarding surprises inside each window.
- **Queue-driven workers**: scale a worker scale set on a Service Bus queue's message count -- the metric source and the scale target are different resources by design.

## Future Enhancements

- Predictive autoscale currently applies to VM Scale Set targets only (an Azure service boundary, documented on the field).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
