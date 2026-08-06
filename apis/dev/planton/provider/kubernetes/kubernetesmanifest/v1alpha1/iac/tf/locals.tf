# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. Stamped on the namespace
  # this module creates — NEVER injected into the manifest's own documents
  # (the manifest is applied exactly as the user wrote it).
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesManifest"
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
    local.identity_labels_base,
    local.id_label,
    local.org_label,
    local.env_label,
  )

  # The anchor namespace (the generator flattens the spec's value-or-ref to a
  # plain string) and whether this module creates it.
  namespace_name   = var.spec.namespace
  create_namespace = try(var.spec.create_namespace, false)

  # Document split rule (identical in the Pulumi module's inventory parser):
  # separators are lines starting with `---`; the prepended newline makes a
  # manifest that STARTS with `---` yield an empty first chunk instead of a
  # missed document. Blank chunks are dropped here; comment-only chunks
  # decode to null and are dropped below; an INVALID chunk fails the plan
  # through yamldecode — never silently skipped (the Pulumi engine rejects
  # the same document at apply).
  manifest_document_chunks = [
    for doc in split("\n---", "\n${var.spec.manifest_yaml}") : doc
    if trimspace(doc) != ""
  ]

  # Decoded documents in manifest order (nulls from comment-only chunks
  # dropped). List comprehension keeps the original document order — never
  # rebuild this from map keys, which sort lexically.
  manifest_documents = [
    for doc in local.manifest_document_chunks : doc
    if yamldecode(doc) != null
  ]

  # for_each keys carry the document's full identity
  # (apiVersion/Kind/namespace/name) so reordering documents never re-keys
  # Terraform state addresses (index keys would churn every later resource).
  # Two documents with the same identity collide loudly at plan time — a
  # duplicate document is a manifest bug, not something to apply twice.
  manifest_docs_by_identity = {
    for doc in local.manifest_documents :
    format("%s/%s/%s/%s",
      try(yamldecode(doc).apiVersion, ""),
      try(yamldecode(doc).kind, ""),
      try(yamldecode(doc).metadata.namespace, ""),
      try(yamldecode(doc).metadata.name, ""),
    ) => doc
  }

  # The applied-resource inventory ("apiVersion/Kind/name" per document,
  # manifest order) — derived from the input YAML so both engines export an
  # identical list regardless of how each engine tracks child resources.
  applied_resources = [
    for doc in local.manifest_documents :
    format("%s/%s/%s",
      try(yamldecode(doc).apiVersion, ""),
      try(yamldecode(doc).kind, ""),
      try(yamldecode(doc).metadata.name, ""),
    )
  ]
}
