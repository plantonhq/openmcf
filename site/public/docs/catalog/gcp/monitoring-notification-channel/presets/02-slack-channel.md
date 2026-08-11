---
title: "Slack Channel"
description: "Posts incident notifications into a Slack channel — the team-visibility tier of alerting, where everyone sees opens and closes without being paged individually."
type: "preset"
rank: "02"
presetSlug: "02-slack-channel"
componentSlug: "monitoring-notification-channel"
componentTitle: "Monitoring Notification Channel"
provider: "gcp"
icon: "package"
order: 2
---

# Slack Channel

Posts incident notifications into a Slack channel — the team-visibility
tier of alerting, where everyone sees opens and closes without being
paged individually.

## What it configures

- `type: slack` with the target channel in `channelLabels.channel_name`.
- The Slack app OAuth token in `sensitiveLabels.authToken` — a managed
  secret, never plaintext at rest. The token comes from installing
  Google's Cloud Monitoring app (or your own app with `chat:write`) into
  the workspace.

## Adjust before deploying

- **channel_name** — the target channel, `#`-prefixed. The app must be
  invited to private channels.
- **authToken** — supply through the platform's secret handling (the
  literal in this preset is a placeholder shape, not a working token).

## When to choose something else

Slack messages are easy to miss at 3am. Use the **PagerDuty Service**
preset for the escalation tier and keep Slack as the visibility tier on
the same policies.
