# KubernetesHelmRelease Terraform module.
#
# Installs an upstream chart as a real Helm release via helm_release — the
# semantic twin of the Pulumi module's helm.v3.Release. The release itself
# lives in helm_release.tf; this file owns the optional namespace.

# The optional release namespace. Created before the release installs;
# module-owned (labeled with the identity set) rather than delegated to
# helm_release's create_namespace flag, so both engines stamp identical
# governance labels on it.
resource "kubernetes_namespace_v1" "helm_release_namespace" {
  count = local.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace_name
    labels = local.labels
  }
}
