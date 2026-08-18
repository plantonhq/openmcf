# AzureMachineLearningDatastore Terraform Module

## Overview

Registers a datastore on an Azure Machine Learning workspace. The spec's variant block (blob / data-lake / file-share) selects which of the three provider resources is created -- exactly one exists per deployment; all three write the same ARM child collection.

## Resources Created

Exactly one of:

- `azurerm_machine_learning_datastore_blobstorage` -- when `spec.blob_storage` is set
- `azurerm_machine_learning_datastore_datalake_gen2` -- when `spec.data_lake_gen2` is set
- `azurerm_machine_learning_datastore_fileshare` -- when `spec.file_share` is set

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningDatastoreSpec fields; the workspace and storage-target references arrive as resolved literal ARM IDs

## Outputs

- `datastore_id`, `datastore_name` -- resolved from whichever variant resource was created
- `is_default` -- settable on the blob variant, read back on the others

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **The variant switch**: each provider resource carries `count = <variant> != null ? 1 : 0`; the spec's exactly-one-variant CEL guarantees a single resource per deployment, and the outputs `try()` across the three.
- **The identity-argument rename**: the provider calls the same argument `service_data_auth_identity` on the blob resource and `service_data_identity` on the other two -- ONE spec field (`service_data_identity`) feeds both, recorded in the parity manifest.
- **Write-only credentials**: `account_key`, `shared_access_signature`, and `client_secret` are never returned by ARM -- the provider echoes them from configuration (write-normalized in the import catalog).
- **ForceNew breadth**: name, storage target, description, and TAGS are all ForceNew on the provider -- only credentials and the identity mode update in place.
