# AzureEventgridDomainTopic Terraform Module

## Overview

Creates one named event stream (domain topic) inside an Azure Event Grid domain -- the per-tenant mailbox of the multi-tenant pattern. Publishers address it by naming the topic in events sent to the domain's endpoint; subscribers attach event subscriptions to the topic's own ARM ID.

## Resources Created

- `azurerm_eventgrid_domain_topic` -- the topic entry under the domain

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridDomainTopicSpec fields; the domain reference arrives as a resolved literal ARM ID

## Outputs

- `domain_topic_id` -- the topic's ARM resource ID (`{domain_id}/topics/{name}`), the scope event subscriptions attach to
- `domain_topic_name` -- the topic's name

## Behavior Notes

- **The spec takes the domain's ARM ID** (the composable reference shape); the provider addresses domain topics by (resource group, domain name), so the module splits the ID into those segments -- the same ARM object either way.
- **Everything is create-only** -- a domain topic is pure addressing (a name under the domain); changing anything replaces it, briefly interrupting its subscriptions and nothing else.
- **No tags** -- the provider carries no tags argument on domain topics.
- **Works alongside the domain's auto-managed topics**: declaring a topic explicitly pins it regardless of the domain's `auto_create`/`auto_delete` flags; the pinned-topics governance posture sets both flags false on the domain.
- **Billing**: free. Operations are billed on the domain.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege actions the deploying principal needs.
