# Terraform module for Percona Operator for MongoDB

resource "kubernetes_namespace" "percona_operator" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "percona_operator" {
  # Use computed release name from metadata.name to avoid conflicts when multiple instances share a namespace
  name       = local.helm_release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = var.spec.create_namespace ? kubernetes_namespace.percona_operator[0].metadata[0].name : local.namespace

  set {
    name  = "resources.limits.cpu"
    value = var.spec.container.resources.limits.cpu
  }

  set {
    name  = "resources.limits.memory"
    value = var.spec.container.resources.limits.memory
  }

  set {
    name  = "resources.requests.cpu"
    value = var.spec.container.resources.requests.cpu
  }

  set {
    name  = "resources.requests.memory"
    value = var.spec.container.resources.requests.memory
  }

  set {
    name  = "watchAllNamespaces"
    value = "true"
  }

  timeout         = 300
  atomic          = true
  cleanup_on_fail = true
  wait            = true
}

