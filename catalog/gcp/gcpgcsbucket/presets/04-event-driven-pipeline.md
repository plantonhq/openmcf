# Event-Driven Pipeline Bucket

An intake bucket that announces its own changes: every new object under
`uploads/` (and every delete) publishes an event to a Pub/Sub topic, which
is where the processing pipeline — Cloud Run services, Cloud Functions,
Eventarc triggers — picks up.

## What this preset creates

A private, IAM-only bucket with one notification config: `OBJECT_FINALIZE`
and `OBJECT_DELETE` events scoped to the `uploads/` prefix, carrying the
full object metadata (`JSON_API_V1`) so consumers can act without an extra
read. A lifecycle rule ages processed uploads out after 30 days.

## Prerequisites

- A `GcpPubSubTopic` named `upload-events` (replace with your topic). Its
  `topic_id` output is exactly the fully-qualified form the notification
  requires.
- The project's GCS service agent must hold `roles/pubsub.publisher` on
  the topic BEFORE this deploys — the Storage API refuses the config
  otherwise. The agent's email is
  `service-{projectNumber}@gs-project-accounts.iam.gserviceaccount.com`
  (`projectNumber` is this bucket kind's own output); compose the grant
  with a `GcpProjectIamMember`, or scope it to the topic's own IAM
  surface.

## Composing the pipeline

Subscribers attach to the topic, never to the bucket: a
`GcpPubSubSubscription` pushing to a Cloud Run service, or an Eventarc
trigger routing into Workflows. Delivery is at-least-once and unordered —
design consumers idempotent.

## Remix ideas

- Narrow to `OBJECT_FINALIZE` only when deletes carry no meaning for the
  pipeline (fewer messages, cheaper).
- Add `customAttributes` (e.g. `env: prod`) so one shared topic can route
  by subscription filter.
- Notification configs are IMMUTABLE — every change replaces the config
  with a brief un-replayed gap. Change them in maintenance windows if the
  pipeline cannot tolerate gaps.
