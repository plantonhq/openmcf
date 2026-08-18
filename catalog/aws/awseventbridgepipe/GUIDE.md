# AwsEventBridgePipe — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The source is a marriage, the target is a job

Changing the source (or its family's position fields — `starting_position`, `topic_name`, `queue_name`) replaces the whole pipe and resets consumer positions. Changing the target updates in place. Design pipes around their source; move targets freely.

## STOPPED is the money lever

Pipes bill per event processed. `desired_state: STOPPED` halts consumption without losing the pipe or its stream positions — pause through incident triage or cost freezes instead of deleting. State flips ride the same provisioning state machine as creates (minutes, not seconds).

## Match the family block to the ARN

The source/target family blocks must match the ARN's service — a `kinesis` block on an SQS source fails at AWS, not at validation (the ARN may be an unresolved reference at validation time). The at-most-one CELs stop double-blocks; the family-to-service match is yours.

## Credentials never enter the manifest

Kafka and MQ sources authenticate through Secrets Manager secret ARNs — the spec's patterns reject anything that is not a `secretsmanager` ARN. Rotate the secret in Secrets Manager; the pipe follows without a deploy.

## Filters are the cheapest optimization

Only events matching `filter_criteria` reach the enrichment/target — and only matching events bill. A tight filter pattern is both a correctness and a cost control.

## TRACE logging can leak payloads

`log_configuration.level: TRACE` with `include_execution_data: true` writes event payloads into the log destination. Enable it deliberately, and point it at a log group with matching access controls.
