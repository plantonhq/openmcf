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
  description = "Azure Cosmos DB MongoDB database specification"
  type = object({
    # The Cosmos DB account the database lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs.
    cosmosdb_account_id = string

    # The database's name -- unique within the account.
    database_name = string

    # Fixed provisioned throughput in RU/s shared by the database's
    # collections. Mutually exclusive with autoscale_max_throughput
    # (enforced by the spec); leave both unset for serverless accounts
    # or collection-dedicated throughput.
    throughput = optional(number)

    # Autoscale ceiling in RU/s (Azure scales between 10% and 100% of
    # this value). Mutually exclusive with throughput.
    autoscale_max_throughput = optional(number)
  })
}
