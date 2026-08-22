# CloudflareNotificationWebhook guide

The judgment this guide protects you from: the secret never comes back, the type is not yours to choose, and deleting a destination breaks policies without telling them.

## The secret is write-only, forever

Cloudflare's own schema says it outright: "Secrets are not returned in any API response body." You send the secret on create or update and no read ever shows it again -- so drift on it cannot be detected, an imported destination arrives with an empty secret, and re-applying is what puts your configured value back. Keep it in a managed secret and let the platform resolve it at deploy; treat rotation as "set the new value and apply," never as "check what's there."

## The type is a server-side echo

`type` is an OUTPUT, not an input. Cloudflare inspects the URL and decides whether this is a Slack, Google Chat, Discord, Feishu, Datadog, Opsgenie, Splunk, or generic destination -- there is no way to override that classification, and the classification changes how the payload is shaped. If you need Slack-formatted messages, the URL must be a Slack incoming-webhook URL; pointing a generic proxy at Slack yields generic payloads.

## Where the vendor API key goes

For Datadog, Splunk, and Opsgenie, the credential is not a "shared secret" in the verify-the-sender sense -- it is the vendor's API key, and it travels in the same `secret` field. That is Cloudflare's design, and it is why the field is sensitive regardless of which destination kind you register.

## Deleting a destination silently breaks its policies

Notification policies reference destinations by UUID. Delete the destination and every policy that used it simply stops delivering on that channel -- no error at Cloudflare, no plan diff in the policies, and no symptom until an alert you were counting on does not arrive. Retire the policy references first, then the destination. When several policies share one destination, that ordering matters more than it looks.

## Registration does not prove delivery

Cloudflare accepts a URL that is unreachable, misconfigured, or simply wrong. The destination's `last_success` and `last_failure` timestamps are the only real evidence, and they only populate once an alert actually fires. After registering a destination, cause one alert deliberately (a health check you can fail, for instance) and confirm it landed -- once per channel is enough to prove the pipe.

## Pairs well with

- [CloudflareNotificationPolicy](../cloudflarenotificationpolicy/README.md) -- the alert rules that deliver to this destination.
- [CloudflareLogpushJob](../cloudflarelogpushjob/README.md) -- for log records rather than alert events.
