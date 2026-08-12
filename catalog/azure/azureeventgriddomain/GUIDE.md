# Azure Event Grid Domain -- Operational Guide

Judgment calls that matter when you run domains in production.

## Domains are for topic COUNT, not event volume

Choose a domain when the number of STREAMS is the scaling problem (a topic per customer, per device fleet, per team). If you have three well-known streams, three custom topics are simpler -- separate keys, separate firewalls, no per-event topic naming. A domain trades per-stream control at the publish edge for one integration serving thousands of streams.

## Pick the topic lifecycle as a governance decision

Auto-managed topics (the defaults) are frictionless: a subscription materializes its topic, the last unsubscribe removes it. Pinned topics (`auto_create`/`auto_delete` both false, each topic an AzureEventgridDomainTopic resource) make tenant onboarding an auditable, chart-managed act -- and make "who is this stream and why does it exist" answerable from IaC. SaaS products with a real tenant registry should pin; internal platforms can start auto-managed and tighten later (the flags update in place).

## One schema rules every topic

The input schema is domain-wide and create-only: every tenant's events arrive in the same envelope. That uniformity is a feature (one publisher SDK config, one validation path) -- but it means schema choice is a day-one decision for the whole tenant base. CloudEvents 1.0 is the default answer for anything new.

## Publish auth: one key set for every tenant

Domain keys authorize publishing to EVERY topic in the domain -- there are no per-topic keys. That is the right shape when the PUBLISHER is your own service tier. If untrusted parties publish directly, front them with your API (which stamps the topic name server-side); never hand a tenant the domain key. For production, prefer Entra ID (EventGrid Data Sender on the domain) and set `local_auth_enabled: false`.

## Tenant isolation lives on the subscription side

Publishing is shared; CONSUMPTION is isolated per domain topic. A tenant's handlers subscribe to that tenant's topic only -- so treat subscription creation as part of tenant onboarding, next to the domain topic itself. Events published to a topic with no subscription are dropped, not queued (in the pinned posture, create topic and subscription together).

## The domain deletes its topics with it

Destroying a domain removes every domain topic and subscription under it -- extension-resource semantics, no dangling children, and no safety net. In charts, the reverse-dependency destroy order (topics before domain) produces clean teardowns naturally when topics reference the domain by `valueFrom`.
