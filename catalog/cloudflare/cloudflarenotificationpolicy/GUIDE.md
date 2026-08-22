# CloudflareNotificationPolicy guide

The judgment this guide protects you from: which filters a given alert type reads is Cloudflare's contract, and an alert with no working destination fails silently in the worst possible moment.

## Filters belong to alert types, and the pairing is API-owned

The `filters` message carries 43 fields because that is the provider's full set -- but a given alert type reads only a handful. `incident_impact` means something to `incident_alert` and nothing to `tunnel_health_event`; `pool_id` matters to the load-balancing families only. Each field's comment carries the provider's own guidance ("Used for configuring radar_notification", and so on) -- start there. A filter the type ignores is accepted and silently does nothing, which is the quiet way to build a policy that never fires the way you expected.

Only two filter fields are enum-walled here, because only two are walled in the provider: `incident_impact` and `traffic_exclusions`. Everything else is free-form by Cloudflare's design, including numeric thresholds and booleans, which all travel as strings.

## At least one destination, and destinations must be real

This spec requires at least one of `emails`, `pagerduty_ids`, or `webhook_ids` -- the provider does not, so a destination-less policy is a shape Cloudflare will happily accept and never deliver from. That is a real footgun on an alerting product, hence the tightening.

Being present is not the same as working. A PagerDuty service UUID must come from an integration already connected in the Cloudflare dashboard, and a webhook UUID must point at a live `CloudflareNotificationWebhook`. Neither is validated at deploy; both fail at delivery time, which is exactly when you are relying on the alert.

## Many alert types need a bigger plan

Advanced DDoS (L4 and L7), bot traffic, several script-monitor and security-insight families, and others are Business or Enterprise features. Cloudflare refuses the create when the plan does not carry the type -- the failure is honest and immediate, but it happens at apply time, not at review time. Check the type against the account's plan before you promote the manifest.

## Deleting a policy is silent

Nothing downstream changes when a policy disappears: no destination is touched and no other policy notices. The only symptom is alerts that stop arriving, which nobody observes until the thing you were watching for happens. Treat alert coverage as reviewable infrastructure -- it belongs in the manifest for exactly that reason.

## Test what you can before you rely on it

Cloudflare's own alerting surface offers no synthetic fire, so the cheapest confidence is a policy on a type you can actually cause (a health check you can fail on purpose, a tunnel you can stop) delivered to a destination you can watch. Do that once per channel -- email, PagerDuty, webhook -- and you have proven the pipe rather than the configuration.

## Pairs well with

- [CloudflareNotificationWebhook](../cloudflarenotificationwebhook/README.md) -- register the webhook destinations this policy delivers to.
- [CloudflareHealthcheck](../cloudflarehealthcheck/README.md) -- the origin probes behind health-check alerts.
- [CloudflareLogpushJob](../cloudflarelogpushjob/README.md) -- pair a policy on `failing_logpush_job_disabled_alert` with every job that matters.
