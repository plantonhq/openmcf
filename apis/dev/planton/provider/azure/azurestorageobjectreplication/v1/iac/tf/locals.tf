locals {
  # The shared policy GUID, parsed from the destination-side ARM ID (the
  # authoritative copy) -- what `az storage account or-policy` and the
  # monitoring surfaces key on. Both engines derive it the same way, so
  # the output is byte-identical.
  policy_id = element(split("/", azurerm_storage_object_replication.main.destination_object_replication_id), length(split("/", azurerm_storage_object_replication.main.destination_object_replication_id)) - 1)
}
