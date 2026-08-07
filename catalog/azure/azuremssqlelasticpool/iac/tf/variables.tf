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
  description = "Azure SQL elastic pool specification"
  type = object({
    # The logical server the pool lives on, by resolved ARM ID. The
    # server's name and resource group are derived from it.
    server_id = string

    # The Azure region the pool lives in -- must match the parent
    # server's region (ARM rejects a mismatch).
    region = string

    # The pool's name, unique within the server.
    pool_name = string

    # The pool SKU name (BasicPool/StandardPool/PremiumPool or
    # {GP|BC|HS}_{family}); the service tier and hardware family are
    # derived from it.
    sku_name = string

    # The pool's capacity: eDTUs for DTU pools, vCores for vCore pools.
    capacity = number

    # Per-database consumption bounds inside the pool, in the pool's
    # capacity unit. min_capacity defaults to 0 (reserve nothing) because
    # the tfvars wire format drops zero-valued proto fields entirely.
    per_database_settings = object({
      min_capacity = optional(number, 0)
      max_capacity = number
    })

    # The pool's total storage cap: gigabytes XOR bytes (spec-validated).
    max_size_gb    = optional(number)
    max_size_bytes = optional(number)

    # Spread the pool's replicas across availability zones.
    zone_redundant = optional(bool, false)

    # The confidential-computing enclave for every database in the pool,
    # as the spec enum's name string (VBS/DEFAULT_ENCLAVE).
    enclave_type = optional(string)

    # Azure Hybrid Benefit for vCore pools, as the spec enum's name
    # string (BASE_PRICE/LICENSE_INCLUDED).
    license_type = optional(string)

    # Hyperscale pools only: readable HA replicas per database (0-4).
    high_availability_replica_count = optional(number)

    # The maintenance window the pool (and its databases) patch in.
    maintenance_configuration_name = optional(string, "SQL_Default")

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
