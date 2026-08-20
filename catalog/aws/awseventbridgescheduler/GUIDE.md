# AwsEventBridgeScheduler — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The first deploy waits on IAM

CreateSchedule validates the execution role is assumable by `scheduler.amazonaws.com`, and a freshly created role takes up to two minutes to propagate. The provider retries exactly this error — a first deploy that sits for a minute is IAM propagation, not a hang.

## DELETE-after-completion deletes out from under you

`action_after_completion: DELETE` makes AWS remove a completed one-time (`at(...)`) schedule itself — the next deploy finds it missing and recreates it. Use DELETE only for fire-and-forget schedules nothing re-deploys; use NONE (the default) when IaC owns the lifecycle.

## The group is a filing cabinet, not a feature

Groups exist to tag and bulk-delete schedules; they have no behavior. One owned group per team or chart (it carries the identity tags — the schedule itself is untaggable), joined by name from the team's other schedules. Moving a schedule between groups replaces it — cheap for stateless schedules, but plan for the momentary gap.

## Timezones belong on the schedule, not in your head

`cron()` expressions evaluate in `schedule_expression_timezone` — set it to the business timezone ("America/New_York") and daylight-saving shifts are AWS's problem. Unset means UTC.

## Retries need somewhere to go

The default retry policy (24h, 185 attempts) can hammer a broken target for a day. Bound it (`maximum_retry_attempts`) and give exhausted events a `dead_letter_queue_arn` — an SQS queue you actually monitor.
