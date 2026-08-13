# Azure Event Grid Namespace -- Operational Guide

Judgment calls that matter when you run Event Grid namespaces in production.

## Decide the MQTT question before you create

The topic-spaces block is the namespace's one irreversible choice: it cannot be added, removed, or changed later -- only a replacement changes it, and a replaced namespace drops every namespace topic inside it. If there is ANY chance the namespace will serve MQTT clients, create it with the block (an empty block is legal and costs nothing extra); a pure-CloudEvents namespace that later needs MQTT is a migration, not an update.

## A namespace is not free the way a topic is

Classic Event Grid topics cost nothing at rest. A namespace bills per throughput unit per hour from the moment it exists -- one TU is the floor. That changes the deployment pattern: share one namespace across services in an environment (streams are cheap namespace topics) instead of stamping one namespace per service the way classic topics are stamped.

## Capacity is shared, and it moves in place

Throughput units cap ingress and egress for ALL topics and MQTT traffic in the namespace together -- one noisy publisher can starve its neighbors' streams. Watch the throttled-requests metric, and raise `capacity` in place when it climbs; no downtime, no replacement.

## The classic and namespace worlds do not mix

Namespace topics speak CloudEvents v1.0 only and are NOT valid targets for the classic kinds' machinery: an AzureEventgridEventSubscription's `scope` arm addresses classic topics/domains, not namespace topics (Azure models namespace-topic subscriptions as a different resource the Terraform provider does not ship yet -- until it does, manage those subscriptions outside or route MQTT into a classic topic and subscribe there). The one supported bridge today is the MQTT route: `topic_spaces_configuration.route_topic_id` forwards every MQTT message into a classic custom topic, where the full delivery machinery applies.

## Client certificates are the MQTT identity model

The broker authenticates MQTT clients against certificates; `alternative_authentication_name_sources` decides which certificate field carries the client's identity (subject, DNS SAN, URI SAN, IP SAN, email). Pick ONE convention fleet-wide before onboarding devices -- changing the source later means re-issuing or re-mapping every client. The session dials (max sessions per name, session expiry hours) are per-authentication-name; devices sharing a name share that budget.
