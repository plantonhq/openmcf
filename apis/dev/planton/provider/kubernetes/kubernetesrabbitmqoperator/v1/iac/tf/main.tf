# KubernetesRabbitMqOperator Terraform module.
#
# Installs the RabbitMQ Cluster Operator from its released single-file
# manifest — the operator's OFFICIAL distribution (it has no Helm chart).
# The manifest applies per document, with the spec's typed overrides
# patched onto the operator Deployment (see locals.tf) and every other
# document applied verbatim: the namespace `rabbitmq-system`, the
# RabbitmqCluster CRD, RBAC, the webhook + metrics Services, the
# cert-manager Issuer and Certificates, and the mutating/validating
# webhook configurations.
#
# CERT-MANAGER IS A HARD PREREQUISITE (a registry prerequisite of this
# kind): the webhook serving certificate is a cert-manager Certificate
# with CA injection — without a running cert-manager the certificate never
# issues and every RabbitmqCluster admission fails (the webhooks are
# failurePolicy: Fail).
#
# DESTROY SEMANTICS: every document deletes with the resource, INCLUDING
# the CRD — which cascade-deletes every RabbitmqCluster on the cluster.
# The spec's CRD-lifecycle note carries the warning.

# Fetch the released manifest from the pinned release tag.
data "http" "cluster_operator_manifest" {
  url = local.manifest_url

  request_headers = {
    Accept = "application/yaml"
  }
}

# Apply every document. Keyed by each document's COMPOSED IDENTITY
# (apiVersion//kind//name[//namespace]) — stable across manifest
# reorderings, and exactly the ID form the kubectl importer takes so the
# import map derives IDs blind from the address keys.
#
# server_side_apply is REQUIRED, not stylistic: the RabbitmqCluster CRD
# document (~342 KB) exceeds the client-side last-applied-configuration
# annotation cap (256 KB). The Pulumi twin rides its provider's
# server-side apply for the same reason.
resource "kubectl_manifest" "cluster_operator" {
  for_each = local.applied_documents

  yaml_body = yamlencode(each.value)

  server_side_apply = true
  force_conflicts   = true
}
