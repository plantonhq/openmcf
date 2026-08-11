# GcpEventarcTrigger Guide

Operational judgment for running Eventarc triggers as code — the things
the spec reference cannot tell you.

## The first trigger in a project is slow — by design

Creating the first trigger provisions Eventarc's service agent (P4SA)
and, for Pub/Sub-transported types, the transport topic and
subscription. Deliveries in the first few minutes after that first
apply can silently drop while IAM propagates. Do not debug your
destination yet: wait, then publish a fresh test event. Subsequent
triggers in the same project are immediate.

## The trigger's service account is TWO grants, not one

`roles/eventarc.eventReceiver` lets the trigger RECEIVE events;
authenticated Cloud Run destinations additionally need
`roles/run.invoker` ON THE SERVICE so the trigger's identity tokens are
accepted. Missing the second grant looks like "trigger fires, service
403s" in the platform logs. Audit-log triggers refuse to create without
a service account at all — and so do WORKFLOW destinations: the API
rejects the create with "trigger.service_account is empty" (the account
mints the workflow executions; pair `eventReceiver` with
`roles/workflows.invoker`). The spec enforces both requirements at
manifest time so you never meet the raw 400.

## Cloud Run destinations always carry a region on the wire

The API never infers the Cloud Run service's region — a create without
one is rejected ("cloud_run.region is empty"). Leave
`destination.cloudRunService.region` empty and the module sends the
trigger's own location, which is what you want in the common
same-region case. Set it explicitly when the service lives elsewhere,
and always on a `global` trigger (audit-log triggers are usually
global; "global" is not a Cloud Run region — the spec enforces this).

## Destination validation happens at the API, not at plan

The provider sends whatever destination combination you wrote and lets
the API reject it at apply time. The spec's exactly-one wall catches the
structural mistakes at manifest time, but SEMANTIC mismatches (a GKE
destination in a project without gke-destinations init, a Cloud Run
service in another project) surface only live — read the create error,
it names the real constraint.

## Pub/Sub triggers own neither payload format nor content type

`eventDataContentType` is rejected outright on
`messagePublished` triggers ("Pub/Sub triggers do not accept any value
for event_data_content_type") — a Pub/Sub event's payload format is
whatever the publisher sent. The field belongs to the OTHER event types
(Storage, audit-log), where Eventarc serializes the event itself. The
spec enforces this at manifest time.

## Event types decide the filterable attributes

`matchingCriteria` beyond `type` must use attributes the EVENT TYPE
declares filterable — `bucket` exists for Storage events, `serviceName`
+ `methodName` for audit logs. An attribute the type does not support is
rejected at create. The audit-log type
(google.cloud.audit.log.v1.written) is the broadest: it can watch
almost any GCP API call, at the cost of enabling audit logs for the
watched service first.

## Transport topics outlive the trigger deliberately

With `transportPubsubTopic` set, Eventarc uses your existing topic and
NEVER deletes it — destroy removes the trigger and its subscription
only. This is the arm to use when the topic is shared infrastructure (a
GcpPubSubTopic managed in its own right). Without it, Eventarc mints and
owns a hidden topic per trigger.

## Partner channels are a handshake, not a resource

Creating a partner channel yields a one-time `activation_token` output;
the channel sits PENDING and delivers NOTHING until the partner redeems
it in their console. Recreating the channel (name or provider change —
both immutable) mints a NEW token and the handshake must be redone.
Treat the token like a credential: it is marked sensitive in both
engines.

## The google channel config is a singleton — one owner

`googleChannelCryptoKey` manages a per-project-per-location SINGLETON
shared by every non-partner trigger in that project+location. Two
triggers both setting it fight over the same object on every apply. Pick
ONE trigger (or a dedicated config-owner instance) per project+location
to carry it; the provider's delete is a state-only no-op, so destroying
the owner reverts nothing in GCP.

## Teardown discipline

`DELETE` removes the trigger; in-flight events in the transport
subscription are lost with it. `ABANDON` keeps delivery running
unmanaged — the escape hatch when routing must survive an IaC
migration. `PREVENT` protects production event routes from accidental
teardown; the partner channel inherits the same policy, protecting the
handshake investment.
