# Terraform module for Altinity ClickHouse Operator
# This is a placeholder - the operator is primarily deployed via Helm/Pulumi

resource "kubernetes_namespace" "kubernetes_altinity_operator" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name = local.namespace
  }
}

resource "helm_release" "kubernetes_altinity_operator" {
  name       = local.helm_release_name
  repository = "https://docs.altinity.com/clickhouse-operator/"
  chart      = "altinity-clickhouse-operator"
  version    = "0.25.4"
  namespace  = local.namespace

  set {
    name  = "operator.createCRD"
    value = "true"
  }

  set {
    name  = "watchNamespaces"
    value = "{}"
  }

  set {
    name  = "operator.resources.limits.cpu"
    value = var.spec.container.resources.limits.cpu
  }

  set {
    name  = "operator.resources.limits.memory"
    value = var.spec.container.resources.limits.memory
  }

  set {
    name  = "operator.resources.requests.cpu"
    value = var.spec.container.resources.requests.cpu
  }

  set {
    name  = "operator.resources.requests.memory"
    value = var.spec.container.resources.requests.memory
  }

  timeout         = 300
  atomic          = true
  cleanup_on_fail = true
  wait            = true
}

