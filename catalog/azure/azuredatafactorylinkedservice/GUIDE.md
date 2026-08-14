# Azure Data Factory Linked Service -- Operational Guide

Judgment calls that matter when you run Data Factory connections in production.

## Deploy the Key Vault connection first, then stop pasting secrets

The single highest-leverage move: one `key_vault` connection per factory, then every other connection's password / connection string / token lives in the vault and is referenced by name (`keyVaultPassword`, `keyVaultConnectionString`, and friends). Rotation becomes a vault operation -- no manifest changes, no redeploys, and Data Factory picks up the new value on the next run. Grant the factory's managed identity get/list on the vault's secrets; the connection itself carries no credential at all.

## Saving is not connecting

Azure accepts a linked service definition without dialing the target -- a wrong password, a firewalled host, or a missing role saves GREEN and fails at pipeline run time. After deploy, use Studio's Test connection (Manage -> Linked services); for managed-identity connections also verify the factory identity's data-plane role on the target (Storage Blob Data Contributor, Key Vault secrets get/list, database AAD user) -- the most common "worked in dev" failure is a missing grant, not a wrong spec.

## Secrets never come back

Inline secrets (connection strings, passwords, access tokens) are stored as hidden secure strings: ARM reads return them masked or not at all, and an import carries the connection's address but never its credential. Two consequences: keep the SOURCE of truth in the manifest (or better, in Key Vault), and treat a plain-text `connection_string_insecure` as what it is -- readable by anyone who can read the factory; use it only for strings that carry no secret.

## One name, one connection -- switching types replaces it

All 23 types share the factory's linked-service namespace. Changing the variant block redeploys the same ARM address as a different type; every dataset and pipeline referencing the name follows instantly. That is the upgrade lever (swap a paste-the-key blob connection for a managed-identity one with no downstream edits) and the foot-gun (a wrong variant swap breaks every consumer at once) -- diff carefully.

## Private targets ride the integration runtime

The default Azure runtime reaches public endpoints only. On-premises SQL Server, ODBC DSNs, and VNet-isolated systems need a self-hosted integration runtime named in `integrationRuntimeName` -- and the ODBC driver for a DSN must be installed on that runtime's machine, not in Azure. The runtime is part of the connection's failure domain: when it is down, every connection through it is down.

## The custom form trades validation for reach

`custom` carries any connector Data Factory speaks as raw JSON -- with no schema checking until Azure parses it at save time. Keep secrets inside it as Key Vault reference objects (`{"type": "AzureKeyVaultSecretReference", ...}`), never literals, and prefer a first-class variant the moment one exists: typed fields, validation, and secret marking are what the custom form gives up.
