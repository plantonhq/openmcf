# Global Static Website

A globally-cached static website from storage to TLS. One deploy produces
a StorageV2 account with static-website hosting, Azure Front Door with
edge caching and compression, and a custom domain that self-completes Front
Door validation through DNS records. For any team shipping a SPA,
documentation site, or marketing page without touching a VM or container
platform.

## The architecture

Traffic flows from the user's browser through Front Door to blob storage:

- **Storage** is a locally-redundant StorageV2 account with the static-
  website feature enabled (`index_document` / `error_document`). Site
  assets live in the auto-created `$web` container. The Front Door origin
  references the account's `primary_web_host` output.
- **Front Door** provides the global edge: HTTPS termination, caching
  (query strings ignored -- assets version by path), and compression for
  text-based MIME types. The default `*.azurefd.net` hostname works
  immediately for smoke testing.
- **Custom domain** is on by default -- this chart's point. It creates an
  Azure DNS zone, a Front Door custom domain with managed TLS, a
  `_dnsauth.<label>` TXT record (validation token by reference), and a
  CNAME from your hostname to the endpoint's `host_name` output. The
  domain deploys in pending-validation state until registrar NS delegation
  is live and Azure confirms the TXT record.
- **Optional apex alias** (`apex_alias_enabled`): an alias A record at `@`
  tracking the Front Door endpoint by ARM ID -- the DNS-correct way to
  serve a naked domain (CNAME is forbidden at the zone apex).
- **Optional WAF** (`front_door_waf_enabled`, off at defaults): attaches
  a Front Door WAF policy for sites that need edge protection without the
  full production-web-app estate.

## What is on by default

- **Custom domain wiring on** -- honest about pending validation until NS
  delegation completes.
- **Edge caching and compression** on both the default and custom-domain
  routes.
- **WAF off** -- low cost, high delight; enable when the site faces
  untrusted input.
- **LRS storage** -- sufficient for static content; upgrade replication
  in the template if compliance requires geo-redundant blobs.

## Parameters worth understanding

- **`storage_account_name`**: globally unique (3-24 lowercase alphanumerics).
  Change before deploying.
- **`dns_zone_name` / `custom_domain_label`**: the zone MUST match a
  domain you can delegate; the label is the subdomain Front Door serves
  (e.g. `www` serves `www.<dns_zone_name>`). After deploy, point your
  registrar's NS records at the Azure DNS zone's name servers (visible in
  the portal or zone outputs).
- **`index_document` / `error_document`**: set both to `index.html` for
  client-routed SPAs.
- **`apex_alias_enabled`**: enable only when you want `example.com` (not
  just `www.example.com`) to resolve through Front Door.
- **`front_door_waf_enabled`**: enable for public forms or APIs co-hosted
  as static assets; start in DETECTION if tuning is needed.

## After deployment

Provisioning takes roughly 10-15 minutes; custom-domain TLS may take
longer once NS delegation propagates.

- **Upload site files** to the `$web` container in the storage account
  (Azure CLI, portal, or CI pipeline).
- **Delegate DNS** at your registrar to the chart-created zone's name
  servers, then wait for Front Door domain validation to flip to approved.
- **Smoke test** the `*.azurefd.net` endpoint before DNS propagates.
- **Natural next steps**: enable WAF, turn on `apex_alias_enabled`, or
  add a CDN purge workflow for cache-busted deployments.
