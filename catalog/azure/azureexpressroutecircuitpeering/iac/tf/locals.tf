locals {
  # NOTE: no tag locals -- an ExpressRoute circuit peering is an ARM
  # child of the circuit and carries no tags of its own (the provider
  # schema has no tags argument; governance tags live on the parent
  # circuit).

  # The spec's enum NAME mapped onto ARM's vocabulary -- the value is
  # also the ARM child's NAME on the circuit, so a circuit carries at
  # most one peering of each type.
  peering_type_wire = {
    "AZURE_PRIVATE_PEERING" = "AzurePrivatePeering"
    "AZURE_PUBLIC_PEERING"  = "AzurePublicPeering"
    "MICROSOFT_PEERING"     = "MicrosoftPeering"
  }
}
