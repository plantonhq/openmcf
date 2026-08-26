# CloudflareWorkersKvPair guide

Operational judgment for IaC-seeded KV entries. The README covers what each field is; this covers how the pieces interact.

## This kind is for seeded configuration, not application data

Model an entry here when infrastructure owns its value — feature flags, per-environment endpoints, bootstrap documents. Data the Worker writes at runtime must NOT be modeled: the next apply would overwrite it, and drift detection would flag every runtime change as corruption.

## The key is a third of the identity

Account, namespace, and key together identify the entry; changing any of them replaces it (the old key is deleted, the new one written). Renaming a key in place is therefore a delete-then-create — consumers reading the old key see a gap unless they are updated first.

## Keep keys free of slashes

The terraform import ID is the slash-delimited `{account}/{namespace}/{key}` triple, so a key containing `/` cannot be imported even though the API itself accepts it. Prefer `-` or `:` as separators in key names.

## Values are not secrets

KV values are plaintext storage readable by anything holding the namespace binding or an API token with KV read. Credentials belong in a Worker `secret_text` binding or Cloudflare Secrets Store; put only the *name* of a secret in KV if a Worker needs indirection.

## Adopting an entry rewrites metadata formatting once

Import restores both the value and the metadata, but Cloudflare re-serializes the metadata JSON on read — so the first apply after adopting an entry rewrites the metadata's whitespace and key order to match your declaration, and nothing else (measured 2026-08-26). Semantically identical, safe to apply.
