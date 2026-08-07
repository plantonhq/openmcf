variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsFsxOntapStorageVirtualMachine specification"
  type = object({
    region = string
    file_system_id = string
    name = string
    root_volume_security_style = optional(string)
    svm_admin_password = optional(string, "")
    active_directory_configuration = optional(object({
      netbios_name = optional(string, "")
      domain_name = string
      dns_ips = list(string)
      username = string
      password = string
      file_system_administrators_group = optional(string)
      organizational_unit_distinguished_name = optional(string, "")
    }))
  })
}
