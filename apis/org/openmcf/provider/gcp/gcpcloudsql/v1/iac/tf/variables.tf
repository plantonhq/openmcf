variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for GCP Cloud SQL instance"
  type = object({

    # GCP project ID where the Cloud SQL instance will be created.
    # Supports StringValueOrRef: either { value = "..." } or { value_from = { kind = "...", name = "...", field_path = "..." } }
    project_id = object({
      value      = optional(string)
      value_from = optional(object({
        kind       = optional(string)
        env        = optional(string)
        name       = string
        field_path = optional(string)
      }))
    })

    # Region where the instance is deployed, for example "us-central1".
    region = string

    # Database engine type (MYSQL or POSTGRESQL).
    database_engine = string

    # Database version string, engine-specific (e.g., "MYSQL_8_0" or "POSTGRES_15").
    database_version = string

    # Machine tier/type for the instance, for example "db-n1-standard-1".
    tier = string

    # Storage size in gigabytes for the database instance.
    storage_gb = number

    # Network configuration for the Cloud SQL instance.
    network = optional(object({
      # VPC network ID for private IP connectivity.
      # Supports StringValueOrRef: either { value = "..." } or { value_from = { kind = "...", name = "...", field_path = "..." } }
      vpc_id = optional(object({
        value      = optional(string)
        value_from = optional(object({
          kind       = optional(string)
          env        = optional(string)
          name       = string
          field_path = optional(string)
        }))
      }))

      # Whether to enable private IP for the instance.
      private_ip_enabled = optional(bool, false)

      # List of authorized networks (CIDR blocks) allowed to connect via public IP.
      authorized_networks = optional(list(string), [])
    }), {})

    # High availability configuration.
    high_availability = optional(object({
      # Whether to enable high availability (regional failover).
      enabled = optional(bool, false)

      # Secondary zone for high availability failover.
      zone = optional(string)
    }), {})

    # Backup configuration settings.
    backup = optional(object({
      # Whether automated backups are enabled.
      enabled = optional(bool, false)

      # Start time for daily backup window in HH:MM format (UTC).
      start_time = optional(string)

      # Number of days to retain automated backups.
      retention_days = optional(number, 7)
    }), {})

    # Database-specific configuration flags as key-value pairs.
    database_flags = optional(map(string), {})

    # Initial root password for the database instance.
    root_password = string
  })
}

