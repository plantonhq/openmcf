locals {
  # Schema groups carry no Azure tags: ARM does not support tags on
  # Event Hubs entities, so the platform's identity tags live on the
  # parent namespace.

  # Enum wire maps. The tfvars wire format carries the FULL proto enum
  # value name; both enums are required in the spec (unspecified is
  # rejected at validation), so no unspecified fallback row is needed --
  # an unmapped value would fail the plan loudly, which is the right
  # outcome.
  schema_compatibility_map = {
    "NONE"     = "None"
    "BACKWARD" = "Backward"
    "FORWARD"  = "Forward"
  }
  schema_type_map = {
    "AVRO" = "Avro"
    "JSON" = "Json"
  }

  schema_compatibility = local.schema_compatibility_map[var.spec.schema_compatibility]
  schema_type          = local.schema_type_map[var.spec.schema_type]
}
