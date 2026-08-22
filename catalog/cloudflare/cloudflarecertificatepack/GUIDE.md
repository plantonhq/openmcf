# CloudflareCertificatePack guide

Operational judgment for advanced certificate packs. The README covers what each field is; this covers how the pieces interact.

## A pack is an order, not a document you edit

Changing hosts, CA, validation method, or validity days replaces the pack. There is no in-place renew-with-new-hosts path. Plan the hostname set (apex plus wildcards, ≤50) before the first apply, and treat a later edit as a delete-and-recreate with a new validation cycle.

## Import is the wrong tool

The provider will import a pack (`<zone_id>/<certificate_pack_id>`), but the docs advise replacing the certificate instead. Most order fields do not round-trip through import, so a post-import plan is noisy even when nothing operationally changed. Prefer a fresh order over adopting an existing pack.

## Validation is a customer-visible wait

`txt` validation asks the customer to create a DNS record; `http` asks them to serve a well-known URL; `email` sends mail to the CA's contacts. The pack sits in `pending_validation` until that proof lands. Soft-delete statuses (`pending_deletion`, `deleted`) still answer GET — an automation that keys on "the id exists" will lie after a destroy.

## Let's Encrypt vs Google vs SSL.com

Let's Encrypt is the cheapest default and is what most zones want. Google and SSL.com are the other advanced-pack CAs; they do not change the Planton shape, they change the CA's validation and branding. `cloudflare_branding` only applies when the CA supports it — leave it unset unless you know you need the branded SAN.
