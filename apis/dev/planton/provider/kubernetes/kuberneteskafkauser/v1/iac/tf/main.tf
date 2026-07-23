# KubernetesKafkaUser Terraform module.
#
# Declares one Kafka client identity on the Strimzi KafkaUser custom
# resource. The target Kafka cluster's USER OPERATOR reconciles it: it
# generates the credentials into a Secret named after the user
# (scram-sha-512 password / tls client certificate) and applies the
# declared ACLs when the cluster runs simple authorization. The CR must
# live in the cluster's own namespace with the strimzi.io/cluster label
# (both rendered from the spec; the placement contract is on the spec's
# field comments).
#
# The CR applies through kubectl_manifest (alekc/kubectl): no cluster
# connection at plan time, so users plan before the Strimzi CRDs exist
# (single-run infra charts, offline plan proofs).
#
# No wait_for block, deliberately: reconciliation belongs to the user
# operator, not to applying the resource.

resource "kubectl_manifest" "user" {
  yaml_body = yamlencode(local.user_manifest)

  server_side_apply = true
}
