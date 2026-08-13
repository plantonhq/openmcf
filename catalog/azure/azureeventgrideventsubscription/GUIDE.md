# Azure Event Grid Event Subscription -- Operational Guide

Judgment calls that matter when you run event subscriptions in production.

## Dead-letter first, destination second

Event Grid does not queue undeliverable events: after the retry policy gives up, an event without a dead-letter destination is gone. The production-grade order is to create the dead-letter container, wire `dead_letter`, THEN cut traffic over -- not the reverse. Treat a subscription without dead-lettering as fire-and-forget by explicit choice, never by omission.

## Choose the destination by failure mode, not familiarity

A storage queue is the cheapest at-least-once consumer (workers drain at their own pace; poison messages sit visibly in the queue). Service Bus adds ordering-by-session and DLQ semantics of its own. Webhooks couple your delivery success to someone else's uptime -- always pair them with retry tuning and dead-lettering, and remember the create-time handshake means the endpoint must be LIVE before the subscription can even be created (deploy ordering matters in charts).

## The 25-value filter budget is per subscription, not per condition

Azure counts every listed value across every advanced-filter condition against one budget of 25. A filter that grows past it fails -- at plan time on scope-addressed subscriptions (the Terraform engine's check), at deploy time otherwise. When a filter approaches the cap, split consumers into multiple subscriptions (they are free) rather than compressing conditions into unreadable regex-like prefixes.

## Filters are ANDed; values are ORed -- design for it

All configured conditions must match simultaneously; within one condition the values are alternatives. "Type A from region X, OR type B from region Y" is therefore TWO subscriptions, not one -- and that is the cheaper, clearer shape anyway: subscriptions cost nothing at rest and each carries its own metrics.

## Identity delivery has a dependency order

Delivering as a managed identity requires: the identity exists ON the source topic, it has the data-plane role on the destination (Storage Queue Data Message Sender, Azure Service Bus Data Sender, ...), and only then the subscription names it. Getting the order wrong fails at create -- or worse, at first delivery. In charts, sequence the role assignment before the subscription.

## Expiration is a feature, not a bug risk

`expiration_time_utc` kills the subscription silently at the deadline -- ideal for temporary integrations (a partner trial, a migration tap), dangerous when set reflexively. Alert on the topic's unmatched-event count if you use expirations in production paths.
