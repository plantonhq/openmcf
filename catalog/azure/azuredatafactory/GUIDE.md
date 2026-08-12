# Azure Data Factory -- Operational Guide

Judgment calls that matter when you run Data Factory workspaces in production.

## Decide the managed-VNet question before you create

Enabling the managed virtual network is an in-place update, but DISABLING it replaces the factory -- and a replaced factory drops every pipeline, dataset, and linked service inside it. If there is ANY chance the factory will need private egress (managed private endpoints require the managed VNet), enable it at creation; the network itself is free and changes nothing until endpoints use it.

## The factory is cheap; what runs inside it is not

A factory at rest costs almost nothing -- billing follows pipeline ACTIVITY (integration runtime hours, data movement volume). That shapes the deployment pattern: one factory per environment, many pipelines inside it, rather than one factory per team. Workspace-level settings (identity, credentials, network posture) are shared by everything inside, so treat the factory manifest as platform-owned and let teams own their pipelines.

## Private endpoints complete on the other side

A managed private endpoint reaches Succeeded while its connection on the TARGET resource is still Pending -- the deploy going green does not mean traffic flows. Approval lives on the target (storage account, SQL server) under Networking -> Private endpoint connections, usually owned by a different team. Wire the approval into your onboarding runbook, or pipelines fail at runtime with connection errors long after the factory deployed cleanly.

## Removing the git block does not detach the repo

The repository binding is applied through a side-channel call after the factory exists, and the provider never calls a repo-clear API -- deleting `github_configuration` or `vsts_configuration` from the manifest leaves the factory bound to the repository. Detach deliberately in the Data Factory Studio (Manage -> Git configuration -> Disconnect). Also decide publishing early: with a repo bound, the collaboration branch is the source of truth and "Publish" is how changes reach the live factory.

## CMK is a one-way door

Customer-managed-key encryption can be enabled at create but never removed -- Azure has no decrypt path back to service-managed keys; only a replacement changes it. Prefer the versionless key reference (the default when wiring an AzureKeyVaultKey) so key rotation propagates without touching the factory; and remember the unwrap identity must hold vault permissions BEFORE create -- Azure validates at deploy time, not at plan time.
