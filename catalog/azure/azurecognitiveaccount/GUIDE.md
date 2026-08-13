# Azure Cognitive Account -- Operational Guide

Judgment that saves real time when running AI services accounts. The field reference lives in the API Explorer; this is the operational layer above it.

## Deleting an account does not free its name

Cognitive accounts soft-delete: the deleted account becomes a purgeable ghost that keeps holding the account name (the Key Vault recycle-bin pattern). A teardown followed by a recreate under the same name fails with a name conflict until the ghost is purged -- the module purges on destroy by default, but accounts deleted outside Planton (portal, scripts) linger. `az cognitiveservices account list-deleted` shows the ghosts; `az cognitiveservices account purge` clears one.

## Set the custom subdomain on day one

Network ACLs and Entra ID (token) authentication both require a custom subdomain, and while one can be ADDED to an account that never had one, CHANGING it replaces the account. Naming it at creation costs nothing and keeps every hardening path open; retrofitting a rename later costs a re-deploy of the account and everything on it.

## OpenAI capacity is quota, not money

A Standard/GlobalStandard account and its deployments carry no idle cost -- but every model deployment draws from a per-subscription, per-region quota of tokens-per-minute. A rejected deployment ("InsufficientQuota") is a quota request away, not a pricing tier away. The expensive class is different: ProvisionedManaged deployments bill their PTU capacity continuously from the moment they exist.

## Kind changes: one pair upgrades, everything else replaces

`OpenAI` <-> `AIServices` changes apply in place -- that is the designed growth path from "we deploy models" to "we run AI Foundry projects and agents". Any other kind change replaces the account (new endpoint, new keys, dependents re-wire). Treat single-service accounts as fixed-purpose.

## Keys off is the hardened posture -- flip it before dependents exist

`localAuthEnabled: false` empties both key outputs and forces Entra ID tokens. Flipping it later invalidates whatever consumed the keys, so decide the auth posture before applications wire themselves to `primary_access_key`. The system-assigned identity output is what Key Vault and storage grants bind to.

## Responsible-AI policies are account-level, selection is per-deployment

Filters live on the account (`raiPolicies`), deployments opt in by name (`raiPolicyName`). Renaming a policy replaces the ARM child AND breaks every deployment that selected the old name -- treat policy names as an interface, not a label. A policy referencing a custom blocklist needs the blocklist on the same account; define both in one spec and the module orders them.
