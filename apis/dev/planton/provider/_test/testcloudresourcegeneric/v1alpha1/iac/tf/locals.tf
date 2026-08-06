# Local values and computed configuration

locals {
  # The synthetic resource id mirrors the kind's registry id_prefix so output
  # shapes look exactly like a real kind's.
  resource_id = "tcrg-${var.metadata.name}"

  # A deterministic endpoint derived from inputs: proves output values flow
  # from resolved inputs through the engine into stack outputs.
  endpoint = "test://${var.metadata.name}"
}
