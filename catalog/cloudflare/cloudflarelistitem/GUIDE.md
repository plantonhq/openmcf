# CloudflareListItem guide

Operational judgment for list entries. The README covers what each field is; this covers how the pieces interact.

## One writer per list

If the parent CloudflareList also declared inline `items`, this kind and that block will overwrite each other. The parent must be an empty container. Grow and trim the set by adding and deleting CloudflareListItem resources.

## The value is the identity

Item values are immutable in the provider: changing an IP, ASN, hostname, or redirect replaces the item (new id). A comment-only edit may in-place; a value edit will not. Treat the value as the key you meant to write, not as something you will mutate later.

## Shape must match the parent kind

An ip list rejects a redirect item at the API, not at YAML parse time. The spec's oneof lets you write any shape; the list's `kind` is the real constraint. Look at the parent before picking `ip` / `asn` / `hostname` / `redirect`.

## Deletes can lag

A just-deleted item may still answer GET for a few seconds. An automation that checks absence immediately after destroy can false-fail; a short retry is the expected workaround, not a sign the delete was ignored.
