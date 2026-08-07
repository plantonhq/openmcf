# KubernetesKafkaTopic Terraform module.
#
# Declares one Kafka topic on the Strimzi KafkaTopic custom resource. The
# target Kafka cluster's TOPIC OPERATOR reconciles it into a real topic —
# the CR must live in the cluster's own namespace with the
# strimzi.io/cluster label (both rendered from the spec; the placement
# contract is on the spec's field comments).
#
# The CR applies through kubectl_manifest (alekc/kubectl): no cluster
# connection at plan time, so topics plan before the Strimzi CRDs exist
# (single-run infra charts, offline plan proofs).
#
# No wait_for block, deliberately: reconciliation belongs to the topic
# operator, not to applying the resource. Deleting the resource deletes
# the TOPIC AND ITS DATA (the topic operator propagates deletion).

resource "kubectl_manifest" "topic" {
  yaml_body = yamlencode(local.topic_manifest)

  server_side_apply = true
}
