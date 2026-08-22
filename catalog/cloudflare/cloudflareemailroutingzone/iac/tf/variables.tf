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
  description = "CloudflareEmailRoutingZoneSpec enables Email Routing on a zone"
  type = object({
    # (Required) The zone ID. StringValueOrRef is flattened to a plain string.
    zone_id = optional(string)

    # (Optional) The single per-zone catch-all rule. Carries a LIST of typed
    # actions (matching the Cloudflare API): each is forward/worker/drop.
    catch_all = optional(object({
      enabled = optional(bool, false)
      name    = optional(string, "")
      actions = list(object({
        type       = string
        forward_to = optional(list(string), [])
        worker     = optional(string, "")
      }))
    }))

    # (Optional) Lock the Email Routing DNS records.
    lock_dns_records = optional(bool, false)

    # (Optional) The domain the managed Email Routing DNS records are created
    # for; empty means the zone apex. Only effective with lock_dns_records.
    dns_name = optional(string, "")
  })
}
