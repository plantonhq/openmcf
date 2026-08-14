# Business-Hours Schedule

This preset deploys calendar-shaped capacity for an App Service plan: three instances (elastic to eight) on weekday business hours, one instance otherwise, with the on-call inbox emailed on every scale action (set `customEmails` to your team's address — Azure retired the subscription-administrator email flags in April 2024).

## When to Use

- Internal apps and B2B services whose load follows the work day
- Any App Service plan (Standard tier or above -- Basic plans do not support autoscale) where paying for peak capacity around the clock is waste

## Key Configuration Choices

- **Three profiles, because a schedule marks a START, not a window** -- the business-hours profile activates at 08:00 and would stay in effect forever; the `weekday-evenings` partner profile at 18:00 is what returns capacity to the off-hours shape. The schedule-less `off-hours` profile covers weekends and is the default
- **A CPU rule inside the business-hours profile** -- the schedule sets the envelope, the rule still reacts to surprises within it; the off-hours profiles pin capacity by carrying no rules
- **Timezone from Azure's fixed vocabulary** -- schedules evaluate in "Eastern Standard Time" here (which tracks daylight saving); replace with your region's Azure timezone name

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-region>` | The plan's region (the setting must match it) | The App Service plan's overview page |
| `<your-service-plan-arm-id>` | The App Service plan's ARM resource ID | `AzureServicePlan` status outputs (`service_plan_id`) or the portal's Properties page |
