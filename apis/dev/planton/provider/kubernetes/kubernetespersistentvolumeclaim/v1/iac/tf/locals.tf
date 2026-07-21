# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesPersistentVolumeClaim"
  }

  id_label = (
    var.metadata.id != null && try(var.metadata.id, "") != ""
  ) ? { "planton.ai/id" = var.metadata.id } : {}

  org_label = (
    var.metadata.org != null && try(var.metadata.org, "") != ""
  ) ? { "planton.ai/organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && try(var.metadata.env, "") != ""
  ) ? { "planton.ai/environment" = var.metadata.env } : {}

  labels = merge(
    try(var.spec.labels, {}),
    local.identity_labels_base,
    local.id_label,
    local.org_label,
    local.env_label,
  )

  annotations = try(var.spec.annotations, {})

  # Fall back to the cluster's "default" namespace when the field arrives
  # null or empty — the same behavior as kubectl without a namespace flag.
  namespace = (
    try(var.spec.namespace, null) == null || try(var.spec.namespace, "") == ""
    ? "default"
    : var.spec.namespace
  )

  # Kubernetes-default access mode applied module-side (the API REQUIRES
  # accessModes; there is no server default), identical to the Pulumi module.
  access_modes = length(try(var.spec.access_modes, [])) > 0 ? var.spec.access_modes : ["ReadWriteOnce"]

  # The API string for volume_mode with the server default (Filesystem).
  volume_mode_map = {
    "filesystem" = "Filesystem"
    "block"      = "Block"
  }
  volume_mode = lookup(local.volume_mode_map, try(var.spec.volume_mode, "filesystem"), "Filesystem")

  # The storageClassName wire value: the class name, "" when dynamic
  # provisioning is explicitly disabled, or null (unset — cluster default
  # applies). The empty-vs-absent distinction is load-bearing upstream.
  storage_class_name = (
    var.spec.disable_dynamic_provisioning
    ? ""
    : (try(var.spec.storage_class_name, null) == null || try(var.spec.storage_class_name, "") == "" ? null : var.spec.storage_class_name)
  )

  # PVC clones live in the core group (apiGroup omitted); VolumeSnapshot
  # restores name the snapshot.storage.k8s.io group.
  data_source_api_group = (
    try(var.spec.data_source, null) == null
    ? null
    : (var.spec.data_source.kind == "volume_snapshot" ? "snapshot.storage.k8s.io" : null)
  )
  data_source_kind = (
    try(var.spec.data_source, null) == null
    ? null
    : (var.spec.data_source.kind == "volume_snapshot" ? "VolumeSnapshot" : "PersistentVolumeClaim")
  )
}
