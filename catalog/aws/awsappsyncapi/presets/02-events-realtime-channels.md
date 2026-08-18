# Real-time Events API with channel namespaces

WebSocket pub/sub with nothing to operate: browsers subscribe on channels, backend services publish, and AWS fans events out. The shape for chat, notifications, presence, and live dashboards.

## What this shape gives you

- **Per-phase authorization**: anyone with the API key can connect and subscribe; publishing defaults to IAM so only your backend writes. The presence namespace overrides both directions to the key — clients broadcast their own presence.
- **Namespaces as product surfaces**: `chat/*`, `notifications/*`, and `presence/*` channels are grouped, authorized, and handled per namespace.
- **Inline event handling**: the chat namespace's APPSYNC_JS handler drops events without a message payload before fan-out — validation without a Lambda.

## Adapt it

- Point clients at the `events_http_endpoint` (publish) and `events_realtime_endpoint` (subscribe) outputs.
- Add a Cognito or Lambda auth provider for per-user authorization; namespace handlers can also route events into a Lambda data source with a DIRECT integration.
- Rotate the API key by overlap before its expiry (maximum life 365 days).
