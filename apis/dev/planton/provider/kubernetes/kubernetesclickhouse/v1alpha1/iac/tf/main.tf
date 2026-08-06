# KubernetesClickHouse Terraform module.
#
# Deploys one operator-managed ClickHouse cluster:
#
#   1. the namespace (optional, create_namespace),
#   2. the auth Secret (when users are declared) — one key per user name,
#      referenced from the CHI so no password ever appears in the CR,
#   3. the ClickHouseKeeperInstallation CR (when coordination resolves to
#      a managed Keeper) — the same operator reconciles it,
#   4. the clickhouse.altinity.com/v1 ClickHouseInstallation CR itself —
#      host StatefulSets, Services, generated ConfigMaps and the PDB are
#      all operator-created from it. No ingress resources — exposure
#      composes from first-class kinds referencing the exported handles.
#
# The CRs apply through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a cluster can be PLANNED before the Altinity
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its clusters in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, Keeper quorum, replica sync) that is not part of
# applying the resource — the same never-block-on-a-controller posture as
# the sibling operator-CR modules.
#
# NAMING CONTRACT (operator source at the pinned release): the cluster-wide
# client Service is `clickhouse-<name>`, the per-cluster Service
# `cluster-<name>-<cluster>`, every host's StatefulSet and headless Service
# `chi-<name>-<cluster>-<shard>-<replica>`; the managed Keeper's client
# Service is `keeper-<name>-keeper`.

# The optional namespace. Created before everything else (every other
# resource is namespaced); deleted with the resource. Pre-existing-
# namespace deployments leave create_namespace false.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---- the auth Secret --------------------------------------------------------
# One key per declared user, value = that user's (already-resolved)
# password. The CHI wires each user's password to its key via
# valueFrom.secretKeyRef — rotating the Secret alone does not re-render
# ClickHouse config (upstream-documented); bump any spec field to roll a
# rotation out.
resource "kubernetes_secret_v1" "auth" {
  count = length(try(var.spec.users, [])) > 0 ? 1 : 0

  metadata {
    name      = local.auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = { for u in var.spec.users : u.name => u.password }

  type = "Opaque"

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}

# ---- the ClickHouseKeeperInstallation ----------------------------------------
resource "kubectl_manifest" "keeper" {
  count = local.keeper_managed ? 1 : 0

  yaml_body = yamlencode(local.keeper_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}

# ---- the ClickHouseInstallation ----------------------------------------------
# Applied after the Keeper CR so the coordination reference resolves as
# soon as the operator reconciles (apply order, not readiness — see the
# no-wait posture above).
resource "kubectl_manifest" "clickhouse_installation" {
  yaml_body = yamlencode(local.chi_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_secret_v1.auth,
    kubectl_manifest.keeper,
  ]
}
