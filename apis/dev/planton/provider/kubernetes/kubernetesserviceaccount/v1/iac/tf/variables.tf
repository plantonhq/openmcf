# Input variables for Kubernetes ServiceAccount Terraform module

variable "metadata" {
  description = "Metadata for the service account resource"
  type = object({
    name = string
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes ServiceAccount"
  type = object({
    name        = string
    namespace   = optional(string, "default")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # Names of kubernetes.io/dockerconfigjson secrets in the same namespace,
    # attached as imagePullSecrets so pods running as this identity can pull
    # from private registries without repeating credentials per pod spec.
    image_pull_secrets = optional(list(string), [])

    # Tri-state, mirroring the Kubernetes API field: null (unset) defers to the
    # cluster default, false hardens pods that never talk to the kube-apiserver,
    # true makes the mount explicit. See main.tf for how null is applied.
    automount_service_account_token = optional(bool)

    # Cloud workload-identity binding. At most one arm (gke/eks/aks) should be
    # set; the module translates it into the ServiceAccount annotation the
    # cloud's webhook expects. Omit for ServiceAccounts that never leave the
    # cluster.
    workload_identity = optional(object({
      gke = optional(object({
        service_account_email = string
      }))
      eks = optional(object({
        role_arn = string
      }))
      aks = optional(object({
        client_id = string
        tenant_id = optional(string)
      }))
    }))
  })
}
