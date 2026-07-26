# KubernetesKafkaConnector Terraform module.
#
# Declares one connector instance on the Strimzi KafkaConnector custom
# resource. The target Connect cluster's operator-managed workers run it —
# the CR must live in the Connect cluster's own namespace with the
# strimzi.io/cluster label (both rendered from the spec; the placement
# contract is on the spec's field comments).
#
# The CR applies through kubectl_manifest (alekc/kubectl): no cluster
# connection at plan time, so connectors plan before the Strimzi CRDs
# exist (single-run infra charts, offline plan proofs).
#
# No wait_for block, deliberately: reconciliation belongs to the cluster
# operator (which drives the connector through the Connect REST API), not
# to applying the resource.

resource "kubectl_manifest" "connector" {
  yaml_body = yamlencode(local.connector_manifest)

  server_side_apply = true
}
