# Managed Widget

The recommended default: a `managed` Turnstile widget where Cloudflare chooses the
challenge dynamically. Best general-purpose protection for forms.

## When to use

- Default choice for protecting a login/signup/contact form.

## Key choices

- `mode: managed` — Cloudflare decides whether to show an interactive challenge.
- `domains` — list every domain (and `localhost` for local development).
- The `secret` output wires into a Worker that calls `/siteverify`.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<widget-name>` | Human-readable widget name |
| `<domain>` | A domain the widget runs on |
