# Pub/Sub to Cloud Run

The workhorse event route: a message published to a Pub/Sub topic
invokes a Cloud Run service — with your own topic as transport so the
topic remains shared infrastructure.

## What it configures

- A messagePublished trigger consuming an existing GcpPubSubTopic as
  transport (never deleted with the trigger).
- An authenticated Cloud Run destination POSTing CloudEvents to
  `/events`.
- A dedicated invoker identity.

## Adjust before deploying

- **destination.cloudRunService.service** — reference your GcpCloudRun's
  `service_name` output (same project as the trigger).
- **transportPubsubTopic** — reference your GcpPubSubTopic's `topic_id`
  output, or REMOVE the field to let Eventarc mint a hidden per-trigger
  topic.
- **serviceAccount** — needs `roles/eventarc.eventReceiver` and
  `roles/run.invoker` on the destination service (two grants, not one).

## After deploying

Publish to the topic; the service receives a CloudEvent within seconds.
FIRST trigger in a fresh project: allow a few minutes for Eventarc's
service agent to propagate before judging deliveries.

## When to choose something else

To react to CONTROL-PLANE changes (a bucket created, an IAM grant)
rather than your own messages, start from the **Audit Log to Workflow**
preset.
