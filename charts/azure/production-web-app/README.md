# Production Web App

A customer-facing web application from edge to database. One deploy
produces a VNet-integrated Linux App Service on PremiumV3, a private
PostgreSQL Flexible Server, an RBAC-mode Key Vault wired to the app's
system identity, WAF-fronted Azure Front Door with origin lockdown, and
the observability pack (diagnostics, outside-in availability test, paging
alert). For platform teams shipping a real production web tier without
composing twenty resources by hand.

## The architecture

User traffic enters through Azure Front Door, not the App Service's
default hostname:

- **The network** is a dedicated VNet with two delegated subnets: one for
  App Service integration (`Microsoft.Web/serverFarms`) and one for
  PostgreSQL VNet injection (`Microsoft.DBforPostgreSQL/flexibleServers`).
  A private DNS zone (`privatelink.postgres.database.azure.com`) resolves
  the database internally. Implicit outbound is off on the app subnet --
  user traffic exits through Front Door; the database never needs the
  public internet once injected.
- **PostgreSQL** is VNet-injected with `publicNetworkAccessEnabled:
  false`. Credentials are password auth at defaults (Entra-only is the
  documented hardening path). The administrator password is an honest
  must-change parameter -- nothing in the graph outputs it.
- **The web app** runs on a PremiumV3 plan, VNet-integrated, HTTPS-only,
  with system-assigned identity and Application Insights wired by
  reference. **Origin lockdown** is the critical security seam: only
  Front Door backends filtered to *this* profile's GUID (`xAzureFdid`
  referencing the profile's `resource_guid`) and Application Insights
  availability probes may reach the app. Everything else is denied --
  without this, WAF-fronted Front Door is bypassable through
  `{app}.azurewebsites.net`.
- **Front Door** terminates TLS at the edge, runs a WAF policy in
  PREVENTION mode, and forwards to the web app's `default_hostname`
  output. The origin hostname is referenced, never hand-copied.
- **Key Vault** is RBAC-mode. A Secrets User grant binds to the web app's
  `identity_principal_id` output (not a generic UAI). Secret *values*
  are loaded through your deployment workflow -- the chart provisions the
  vault and the read grant only.
- **Observability** routes web-app and Postgres diagnostics into a Log
  Analytics workspace, runs an outside-in availability probe against
  `https://{web_app_name}.azurewebsites.net`, and pages on multi-
  location failures.

Optional **custom domain** wiring (toggle off at defaults) creates a
public DNS zone, a Front Door custom domain with managed TLS, and TXT/
CNAME records that self-complete validation through reference paths.

## What is on by default

- **Private database** (VNet injection, public access off, zone-redundant
  HA on `postgres_ha_enabled`).
- **WAF in PREVENTION** (`front_door_waf_mode`).
- **Origin lockdown** (Front Door FDID filter + probe allowance).
- **HTTPS-only** app and Front Door forwarding.
- **Custom domain off** (`custom_domain_enabled`) -- the generated
  `*.azurefd.net` endpoint works immediately; enable custom domain when
  you have a registrar-delegated zone.
- **Password database auth** with a documented must-change password.

## Parameters worth understanding

- **`web_app_name` / `key_vault_name`**: globally unique across Azure.
  Change both before deploying.
- **`postgres_admin_password`**: MUST be changed. Rotate through your
  operational workflow; consider Key Vault references in app settings for
  the connection string once secrets are loaded.
- **`postgres_ha_enabled`**: zone-redundant HA doubles database cost;
  disable for non-production estates.
- **`custom_domain_enabled`**: when true, set `dns_zone_name` to a zone
  you control and `custom_domain_label` to the subdomain to serve (e.g.
  `www`), delegate NS records at your registrar to Azure DNS, then wait
  for Front Door validation to complete.
- **`ops_email` / `critical_email`**: Azure sends confirmation mail on
  first deploy. Critical receives availability pages.

## After deployment

Provisioning takes roughly 15-25 minutes (PostgreSQL VNet injection and
Front Door custom-domain validation dominate).

- **Upload application code** to the web app (Git deploy, container, or
  zip deploy). The chart ships Python 3.12 as a neutral runtime default
  -- change `applicationStack` in the template for your stack.
- **Load Key Vault secrets** and wire app settings (e.g.
  `@Microsoft.KeyVault(SecretUri=...)` references).
- **Verify origin lockdown**: browsing `{web_app_name}.azurewebsites.net`
  directly should be denied; traffic through Front Door should succeed.
- **Confirm the availability probe** is green in Application Insights.
- **Natural next steps**: enable `custom_domain_enabled`, add Entra-only
  database auth, or peer this VNet into a hub-spoke foundation.
