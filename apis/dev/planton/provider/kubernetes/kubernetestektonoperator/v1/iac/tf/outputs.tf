# Stack outputs — identical names and derivations in the Pulumi module
# (KubernetesTektonOperatorStackOutputs).

output "namespace" {
  description = "Namespace the operator is installed into (always tekton-operator — the release manifest's fixed namespace)"
  value       = local.namespace
}
