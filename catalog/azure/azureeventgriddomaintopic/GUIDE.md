# Azure Event Grid Domain Topic -- Operational Guide

Judgment calls that matter when you run pinned domain topics in production.

## The topic name is the publisher contract

Publishers stamp this exact name into every event's topic field -- it is API surface, not decoration. Name topics after the stable identity they carry (`customer-fabrikam`, `orders`), never after infrastructure or environments (the domain already scopes those). The name is create-only: renaming replaces the topic and drops its subscriptions, so treat a rename like an API version change.

## Pin topics where tenant identity is real

Declared topics turn "which streams exist" into reviewable IaC -- the right posture when topics map to customers, billing entities, or data-isolation boundaries. Pair this kind with the domain's `auto_create`/`auto_delete` flags set false; otherwise Azure happily materializes undeclared topics the moment someone creates a subscription, and your inventory lies.

## Onboard topic and subscriptions together

A topic with no subscriptions drops every event silently (Event Grid stores nothing at the topic). Tenant onboarding is therefore one chart move: the domain topic AND its dead-letter-backed subscription(s) in the same set, references flowing from `domain_topic_id`. Offboarding is the reverse teardown -- and destroying the topic removes its subscriptions with it.

## Deletion is the isolation boundary working

Destroying a tenant's topic instantly stops delivery to that tenant's handlers while every sibling stream continues untouched -- that independence is exactly why the topic is its own resource. Verify offboarding by the DOMAIN's metrics (delivery failures against a deleted topic mean a publisher is still stamping its name).
