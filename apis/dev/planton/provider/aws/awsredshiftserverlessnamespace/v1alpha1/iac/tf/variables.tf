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
  description = "AwsRedshiftServerlessNamespace specification"
  type = object({
    region = string
    db_name = optional(string, "")
    admin_username = optional(string, "")
    manage_admin_password = optional(bool, false)
    admin_user_password = optional(string, "")
    admin_password_secret_kms_key_id = optional(string, "")
    kms_key_id = optional(string, "")
    iam_roles = optional(list(string), [])
    default_iam_role_arn = optional(string, "")
    log_exports = optional(list(string), [])
  })
}
