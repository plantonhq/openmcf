# Azure Cognitive Deployment -- Operational Guide

Judgment that saves real time when running model deployments. The field reference lives in the API Explorer; this is the operational layer above it.

## Know which SKUs cost money while idle

`Standard`, `GlobalStandard`, `DataZoneStandard` and the Batch variants bill per token -- an idle deployment costs nothing and `capacity` is only a rate limit. The `ProvisionedManaged` variants are the opposite: reserved PTU capacity billed continuously from the moment the deployment exists, whether or not anything calls it. A forgotten provisioned deployment is a standing bill; a forgotten standard one is free.

## Quota errors are requests, not defects

Every deployment draws tokens-per-minute quota from a per-subscription, per-region pool. When ARM rejects a create or a capacity bump with an insufficient-quota error, the fix is a quota increase request (portal -> the account -> Quotas) or a smaller capacity -- nothing in the manifest is wrong. GlobalStandard usually has the most headroom.

## The deployment name is your API contract

Applications pass the DEPLOYMENT name -- not the model name -- when calling the endpoint. Renaming a deployment replaces the ARM object and breaks every caller. Name deployments after their role ("chat", "embeddings"), keep the role stable, and change the model behind it freely (model version updates in place; model name is a new deployment by design).

## Decide the version policy explicitly

Unset `model.version` tracks Azure's current default and `versionUpgradeOption` defaults to upgrading when a new default ships -- right for staying current, wrong for reproducibility. For compliance-pinned workloads set BOTH: an explicit version and `NO_AUTO_UPGRADE`. Watch Azure's model retirement calendar either way: a retired version upgrades or stops regardless of preference.

## The content-filter policy is a name-coupled contract

`raiPolicyName` selects an account-level policy by NAME. Deleting or renaming that policy on the account strands the deployments that selected it. Define policies and deployments in specs that change together, and treat policy names as an interface.
