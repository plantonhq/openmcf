# CloudflareSnippet guide

The judgment this guide protects you from: Snippets need a paid-plan zone, the snippet name is the identity, and create is an upsert. Deploying a name that already exists in the zone silently overwrites it -- there is no "already exists" error to catch the collision.

## Snippets are a paid-plan product

A free-plan zone refuses every snippet upload with "snippets are not allowed" (measured live 2026-08-27, on pending and active free zones alike). Pro and above include Snippets at no extra cost, with per-plan count limits.

## Adopting an existing snippet: apply, never import

Do NOT `tofu import` a snippet: the provider's importer crashes at v5.23.0 (a nil-pointer panic in its own read path -- "Plugin did not respond"; unchanged through v5.24.0 and master, checked 2026-08-27). You do not need it: because create is an upsert by name, applying a manifest whose `snippet_name` matches the live snippet adopts and overwrites it in place. Make the manifest's content the truth you want, apply, done.

## Name is identity: create is an upsert

`snippet_name` is how Cloudflare addresses the snippet. The create call adopts whatever already lives under that name and overwrites its files. Two teams, two manifests, one name -- the last apply wins, and the first team's code is gone.

Letters, numbers, and underscores only. Hyphens are rejected. Changing `snippet_name` replaces the resource (new identity). Rules that referenced the old name keep pointing at it -- they do not follow the rename.

Pick names the way you pick Worker script names: unique in the zone, stable, and not recycled.

## main_module must name a file

`main_module` is the entry file. It must equal one of `files[].name` -- validation rejects a dangling entry point. The provider argument is `metadata` (a nested `{main_module}` object); the spec lifts the leaf so the manifest does not wrap a one-field object.

Most snippets are a single `main.js`. Multi-file snippets import siblings by name via ES module imports.

## Content must be byte-stable

The provider refetches stored content on refresh, rebuilt from the API's multipart response. Server-side normalization -- or a trailing newline you added in the editor -- reads back as configuration drift. Keep file content byte-stable: consistent line endings, no trailing-whitespace churn, no "pretty-print on save" that the API will not echo.

Cloudflare caps snippet code size (32 KB per snippet at the time of writing); oversized code fails at the API. Plan limits also cap snippet COUNT per zone.

## Destroy is a real delete

Destroy removes the snippet. Rules that referenced `snippet_name` start invoking nothing. Update or delete those rules first, or accept a window where matching requests hit a missing snippet.

## Pairs well with

- [CloudflareSnippetRules](../cloudflaresnippetrules/README.md) -- the zone's routing table; one manifest owns the whole table and invokes this snippet by name.
- [CloudflareWorker](../cloudflareworker/README.md) -- the full Worker when you need bindings, cron, or a custom domain.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
