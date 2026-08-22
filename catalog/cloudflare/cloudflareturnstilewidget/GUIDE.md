# CloudflareTurnstileWidget guide

Operational judgment for Turnstile widgets. The README covers what each field is; this covers how the pieces interact.

## The sitekey is the identity

There is no separate widget id. The public sitekey is what the API, the import ID (`<account_id>/<sitekey>`), and the frontend snippet all use. The secret is a sensitive output for `/siteverify`; it is returned on read, so it does round-trip through import.

## Mode is the product choice

`managed` lets Cloudflare decide the challenge (the default you want). `invisible` never shows a widget. `non-interactive` is the no-interaction challenge. Changing mode is an in-place edit; it is also a UX change for every page embedding the sitekey, so treat it as a product decision, not a tuning knob.

## Domains are an allow-list, not a zone

`domains` are hostnames the widget may be served on. They do not need to be Cloudflare zones. `localhost` is valid for local development. A mismatch between the page's origin and this list is the usual "widget not showing" cause.

## Enterprise flags stay off unless you are on Enterprise

`bot_fight_mode`, `ephemeral_id`, `offlabel`, and some `clearance_level` values are Enterprise-gated. Setting them on a non-Enterprise account fails at the API. Leave them unset; the proto accepts them so an Enterprise account can turn them on without a schema change.
