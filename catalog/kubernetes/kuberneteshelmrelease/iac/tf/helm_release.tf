# The Helm release itself.
#
# Every argument maps 1:1 onto a helm.v3.Release argument in the Pulumi
# module — keep the two in lockstep.
#
# VALUES-MERGE PARITY: values_yaml rides the `values` list; the three
# override maps ride set/set_sensitive blocks. The provider merges them in
# exactly the documented precedence (values_yaml, then set, then set_string,
# then set_sensitive) because it processes the values list first, then set
# entries in list order, then set_sensitive. Both this provider and the
# Pulumi module's merge run Helm's own strvals parser on override entries,
# so dotted-path syntax and type coercion behave identically on both
# engines. for_each iterates the maps lexically — the same sorted-key order
# the Pulumi module applies.

resource "helm_release" "helm_release" {
  name = local.release_name

  # Chart identity. For oci:// repositories the provider joins repo + chart
  # into the full OCI reference internally (the Pulumi module performs the
  # identical join client-side).
  repository = var.spec.repo
  chart      = var.spec.chart
  version    = var.spec.version

  # Use created namespace if available, otherwise use the namespace name
  # directly; creation stays with the explicit resource in main.tf
  # (module-owned, labeled) rather than helm_release's own create_namespace
  # flag.
  namespace        = local.create_namespace ? kubernetes_namespace_v1.helm_release_namespace[0].metadata[0].name : local.namespace_name
  create_namespace = false

  # Private chart repository / OCI registry credentials.
  repository_username = try(var.spec.repository_username, "") != "" ? var.spec.repository_username : null
  repository_password = try(var.spec.repository_password, "") != "" ? var.spec.repository_password : null

  # ---- values layers (see VALUES-MERGE PARITY above) ----------------------
  # set/set_sensitive are list-of-object ATTRIBUTES in helm provider v3 (the
  # framework rewrite), not legacy blocks. `set` entries come first (Helm
  # --set coercion), then `set_string` entries (literal strings) — list
  # order IS the merge order, and the map comprehensions iterate keys
  # lexically, matching the Pulumi module's sorted-key application.
  # The same three inputs the CRD render consumed (locals.tf), so the
  # derived CRDs can never see different values than the install.
  values        = local.helm_release_values
  set           = local.helm_set
  set_sensitive = local.helm_set_sensitive

  # ---- lifecycle knobs -----------------------------------------------------
  # skip_crds is pinned: the module owns the chart's crds/ directory
  # (helm_crds.tf derives and applies it kept, ahead of this release, when
  # crds.install is true; crds.install false means someone else owns it).
  # Helm never installs its own once-only copy either way.
  atomic                     = try(var.spec.atomic, false)
  cleanup_on_fail            = try(var.spec.cleanup_on_fail, false)
  wait                       = local.wait
  wait_for_jobs              = try(var.spec.wait_for_jobs, false)
  timeout                    = local.timeout_seconds
  skip_crds                  = true
  dependency_update          = try(var.spec.dependency_update, false)
  max_history                = local.max_history
  replace                    = try(var.spec.replace, false)
  force_update               = try(var.spec.force_update, false)
  reuse_values               = try(var.spec.reuse_values, false)
  reset_values               = try(var.spec.reset_values, false)
  disable_webhooks           = try(var.spec.disable_webhooks, false)
  disable_openapi_validation = try(var.spec.disable_openapi_validation, false)
  take_ownership             = try(var.spec.take_ownership, false)
  description                = try(var.spec.description, "") != "" ? var.spec.description : null

  # The namespace exists and the module-owned CRDs are registered before
  # the release installs (a chart's custom resources would otherwise be
  # rejected by an API server that does not serve their kind yet).
  depends_on = [
    kubernetes_namespace_v1.helm_release_namespace,
    kubectl_manifest.helm_crds,
  ]
}
