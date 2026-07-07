# AwsCloudwatchCompositeAlarm — Research Notes

## What AWS models

A composite alarm (`PutCompositeAlarm`) is a first-class CloudWatch alarm type
whose state derives entirely from a boolean **rule expression** over other
alarms' states — `ALARM("name")`, `OK("name")`, `INSUFFICIENT_DATA("name")`
combined with `AND`/`OR`/`NOT`, parentheses, and the `TRUE`/`FALSE` constants
(useful while wiring a new composite). It shares the metric-alarm namespace
(names are unique per account/region across both alarm types) and appears in
`DescribeAlarms` only when `AlarmTypes` includes `CompositeAlarm`.

Key properties:

- **No evaluation plumbing of its own.** No metrics, periods, statistics, or
  thresholds — the composite re-evaluates its rule whenever any referenced
  alarm transitions. This is why the spec is small and honest rather than a
  thin metric-alarm clone.
- **Name-based composition.** The rule addresses alarms by NAME, not ARN.
  Planton alarms name themselves after `metadata.name` and export
  `alarm_name`; the rule text passes through verbatim. Rule references are
  the same composition class as metric-math expression variables — strings
  the service interprets, not foreign keys the platform can resolve.
- **Action restrictions.** Composite alarms support SNS topics and Systems
  Manager OpsItems/incidents as actions — NOT Auto Scaling policies or EC2
  actions (those stay on metric alarms).
- **`actions_enabled` is ForceNew here** (unlike on metric alarms, where it
  updates in place) — the provider replaces the composite when it changes.
- **Actions suppressor.** `actions_suppressor` designates an alarm whose
  ALARM state withholds this composite's actions; `wait_period` bounds how
  long the composite waits for the suppressor to catch up after its own
  transition, and `extension_period` keeps suppression active after the
  suppressor clears. State transitions always record — only actions are
  suppressed, so alarm history stays honest through maintenance windows.
- **Cycles are rejected by AWS** at PutCompositeAlarm time (a composite may
  not reference itself directly or transitively).

## 90/10 decisions

- The full provider surface is modeled: rule, description, three action
  lists, presence-aware `actions_enabled`, and the suppressor block. Nothing
  in `aws_cloudwatch_composite_alarm` is skipped.
- Rule-expression grammar validation is delegated to the AWS API (the spec
  bounds length only) — re-implementing the parser in CEL would chase a
  moving grammar for no user benefit; deploy-time errors are precise.
- The suppressor's `alarm` is a `StringValueOrRef` defaulting to
  `AwsCloudwatchAlarm.status.outputs.alarm_name`, because unlike rule text it
  is a single addressable field the platform CAN resolve.

## Composition patterns

1. **Shared-cause page**: two symptom alarms AND-ed into one page.
2. **Dependency-aware page**: `ALARM("svc") AND NOT ALARM("upstream")` routes
   pages to the team that can act.
3. **Maintenance window**: a deploy pipeline publishes a flag metric; a flag
   alarm watches it; every paging composite lists the flag alarm as its
   suppressor.
4. **Hierarchies**: composites reference other composites by name (exported
   as `alarm_name`), building service → product rollups.
