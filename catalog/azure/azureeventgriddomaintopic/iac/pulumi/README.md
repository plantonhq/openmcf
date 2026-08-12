# AzureEventgridDomainTopic Pulumi Module

## Overview

Creates one named event stream (domain topic) inside an Azure Event Grid domain -- the per-tenant mailbox of the multi-tenant pattern. Publishers address it by naming the topic in events sent to the domain's endpoint; subscribers attach event subscriptions to the topic's own ARM ID.

## Resources Created

- `eventgrid.DomainTopic` -- the topic entry under the domain

## Outputs

- `domain_topic_id` -- the topic's ARM resource ID (`{domain_id}/topics/{name}`), the scope event subscriptions attach to
- `domain_topic_name` -- the topic's name

## Behavior Notes

- **The spec takes the domain's ARM ID** (the composable reference shape); the SDK addresses domain topics by (resource group, domain name), so the module splits the ID into those segments -- the same ARM object either way.
- **Everything is create-only** -- a domain topic is pure addressing (a name under the domain); changing anything replaces it, briefly interrupting its subscriptions and nothing else.
- **No tags** -- the provider carries no tags argument on domain topics.
- **Works alongside the domain's auto-managed topics**: declaring a topic explicitly pins it regardless of the domain's `auto_create`/`auto_delete` flags; the pinned-topics governance posture sets both flags false on the domain.
- **Billing**: free. Operations are billed on the domain.

## Required Permissions

The deploying principal needs `Microsoft.EventGrid/domains/topics/*` (EventGrid Contributor covers it).
