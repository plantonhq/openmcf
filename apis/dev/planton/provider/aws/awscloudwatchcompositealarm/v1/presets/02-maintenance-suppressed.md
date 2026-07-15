# Preset: Maintenance-Suppressed Paging

**Use case:** Stop paging during planned maintenance without deleting or
disabling alarms.

This composite pages when the API is unhealthy — unless a designated
"maintenance flag" alarm is in ALARM, in which case actions are withheld.
Deploy pipelines raise the flag (e.g. by publishing a metric the flag alarm
watches) before a rollout and clear it after.

## What You Get

- A composite alarm over the API's symptom alarms, combined with `OR`
- An actions suppressor referencing the maintenance-flag alarm by its
  exported `alarm_name` output
- State transitions still evaluate and record during maintenance — only the
  paging is silenced, so the alarm history stays honest
- Outputs: `alarm_arn`, `alarm_name`

## When to Use

- Planned deploys that predictably trip latency/error alarms
- Scheduled batch windows that saturate shared infrastructure
- Any "we know, don't page" window that should never require editing alarms

## Key Configuration Choices

- **`waitPeriod: 60`** — after this composite transitions, it waits up to 60
  seconds for the suppressor to enter ALARM before concluding the window is
  not active and acting
- **`extensionPeriod: 120`** — actions stay suppressed for 2 minutes after
  the flag clears, absorbing transitions from the wind-down itself
- **Reference, not string** — the suppressor alarm composes via `valueFrom`
  against the flag alarm's `alarm_name` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region of the composite and its constituent alarms | Your deployment region |
| `<api-5xx-alarm-name>` | Name of the 5xx symptom alarm | The alarm's `alarm_name` output |
| `<api-latency-alarm-name>` | Name of the latency symptom alarm | The alarm's `alarm_name` output |
| `<maintenance-flag-alarm>` | Planton resource name of the maintenance-flag alarm | Your maintenance-flag AwsCloudwatchAlarm |
