locals {
  # The spec's service enum arrives as the FULL proto value name; the
  # map carries the complete vocabulary, translated to the API's
  # lowercase wire values.
  service_map = {
    "BLOB" = "blob"
    "FILE" = "file"
  }

  home_directory = var.spec.home_directory == null || var.spec.home_directory == "" ? null : var.spec.home_directory

  # The account name, parsed from the resolved account ARM ID -- used
  # for the storage_account_name output and to compose the full SFTP
  # login ({account}.{user}). The named-group regex fails the plan
  # loudly if the ID is not a storage-account ARM ID.
  storage_account_name = regex("/storageAccounts/(?P<name>[^/]+)$", var.spec.storage_account_id)["name"]

  sftp_username = "${local.storage_account_name}.${var.spec.user_name}"
}
