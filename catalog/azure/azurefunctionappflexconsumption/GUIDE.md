# Azure Function App Flex Consumption -- Operational Guide

Judgment calls that matter when you run Flex Consumption apps in production.

## The plan tier is a hard gate, checked before anything else

Azure rejects app creation on any plan whose SKU is not FC1 ("the sku name is ... which is not valid for a flex consumption function app") -- and the check runs against the LIVE plan, so a plan re-tiered after the reference was wired still fails. Keep flex apps on dedicated FC1 plans; one plan hosts many apps and costs nothing at rest. Flex Consumption is also not offered in every region -- verify availability before committing a region, because region is ForceNew. Live-proven: FC1 is available in centralus (the fixture resource group's region); both smoke lanes create there in about a minute.

## Deployment storage: pick the auth mode deliberately

The app's code package lives in the blob container behind `storageContainerEndpoint`, and the auth mode decides your credential posture:

- **Connection-string** (`STORAGE_ACCOUNT_CONNECTION_STRING`): the simplest mode. Azure derives the connection string from the endpoint's account name and your `storageAccessKey`, and manages it as the `DEPLOYMENT_STORAGE_CONNECTION_STRING` app setting (filtered from reads -- you will never see it in `appSettings`). A rotated storage key must be updated in the manifest before the next apply.
- **System-assigned identity** (`SYSTEM_ASSIGNED_IDENTITY`): credential-free, but the grant is day-2 by construction -- the identity exists only after the app does. Create the app, then grant its `identity_principal_id` "Storage Blob Data Contributor" on the storage account; package deployments fail until the grant lands. Live-proven: ARM does not validate storage access at create -- the site object succeeds without the grant. Do not treat a green create as proof the package path works.
- **User-assigned identity** (`USER_ASSIGNED_IDENTITY`): pre-grant the identity before the app exists (the ordering system-assigned cannot offer). The same identity must be attached via `identity.identityIds` AND named in `storageUserAssignedIdentityId`.

## The write-only class: what imports cannot recover

Azure never returns these on reads -- re-supplying them in the manifest is expected, not drift:

- `storageAccessKey` -- in connection-string mode it comes back only embedded in the `AzureWebJobsStorage` app setting; in the identity modes it never comes back at all. Live-proven: a blind import of the connection-string smoke still re-plans empty (the key is re-supplied from the manifest; Azure's omission is not drift).
- `zipDeployFile` -- consumed by the publish endpoint, never stored.
- `siteConfig.appServiceLogs` -- applied on UPDATE operations only and never read back; expect the portal, not the manifest, to reflect its drift. On a fresh create the block lands on the first subsequent update.
- `siteConfig.elasticInstanceMinimum` -- accepted on the wire but never returned on this hosting model; `alwaysReady` is the flex-native warm-instance mechanism, prefer it.

## Always-ready pools are the only idle cost -- and the only cold-start cure

Everything else in a flex app bills per execution. Each `alwaysReady` entry keeps N instances warm for a scope ("http", "durable", "blob", or "function:{name}") and bills for their uptime. The counts' sum must stay within `maximumInstanceCount` (Azure enforces at apply time -- a sum rule manifest validation cannot express). Azure lower-cases pool names on save; treat them case-insensitively.

## Easy Auth secrets live in app settings, never in the auth block

Every identity provider references its client secret by APP SETTING NAME (`clientSecretSettingName`); put the actual value in `appSettings` -- ideally as a Key Vault reference (`@Microsoft.KeyVault(SecretUri=...)`) -- and pin it with `stickySettings` if you use deployment slots. The health-check eviction time also travels as an app setting (`WEBSITE_HEALTHCHECK_MAXPINGFAILURES`) -- managed for you, filtered from reads.

## Hardening the deployment surface

`webdeployPublishBasicAuthenticationEnabled: false` closes the classic username/password publishing path and forces identity-based deployment -- pair it with identity-based storage auth for a credential-free posture. While basic auth IS enabled, the `site_credential_password` output is a working deploy credential; treat it like an admin password.
