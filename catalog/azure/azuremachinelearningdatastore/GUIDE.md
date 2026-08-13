# Azure Machine Learning Datastore -- Operational Guide

Judgment that saves real time when running datastores. The field reference lives in the API Explorer; this is the operational layer above it.

## Prefer identity auth; it removes secrets entirely

`serviceDataIdentity: WORKSPACE_SYSTEM_ASSIGNED_IDENTITY` with no embedded credentials is the posture to aim for: no key rotation, no secret in the manifest, one role assignment (Storage Blob Data Contributor / Reader on the target) instead. The FILE SHARE variant is the exception -- the provider demands a key or SAS there regardless of identity mode; treat file-share datastores as inherently secret-bearing.

## Credentials are write-only -- rotation is a deliberate act

ARM never returns keys, SAS tokens, or client secrets; the state simply echoes what was configured. Rotating a storage key does NOT break the datastore visibly at plan time -- jobs fail at data access instead. When you rotate the account key, update the secret the manifest references and re-apply in the same change window.

## Almost everything is ForceNew -- including tags and description

The provider replaces the datastore when the name, storage target, DESCRIPTION, or TAGS change -- unusual, and worth knowing before wiring tags to anything dynamic. Replacement is cheap (the datastore is a pointer; the DATA is untouched) but running jobs referencing it mid-replace will fail. Update credentials and identity mode freely; change anything else in a maintenance window.

## The built-in datastores are not yours to manage

Every workspace ships with `workspaceblobstore` and `workspacefilestore` pointing at the default storage account. This kind manages ADDITIONAL datastores; do not recreate the built-ins under management (their names are reserved by the service) and do not point a managed datastore at the workspace's own default container -- confusing double-registration follows.

## Claiming is_default moves job outputs

`isDefault: true` (blob variant only) re-points where job outputs land workspace-wide the moment it applies. Flip it deliberately, and remember exactly one datastore holds the flag -- claiming it silently demotes `workspaceblobstore`.
