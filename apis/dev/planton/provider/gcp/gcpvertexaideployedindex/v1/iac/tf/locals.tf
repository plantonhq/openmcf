locals {
  location          = var.spec.location
  deployed_index_id = var.spec.deployed_index_id

  # This resource class carries NO labels and NO project field in the GCP
  # API — the deployment lives inside the index endpoint resource and
  # inherits its project — so there is no label merge here: platform
  # attribution is impossible on a DeployedIndex and none is faked.
}
