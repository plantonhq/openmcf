# This module deploys through the azapi provider (raw ARM), not azurerm:
# azurerm carries NO resource for Machine Learning batch endpoints (its
# endpoint draft is tracked at hashicorp/terraform-provider-azurerm#32011).
# The resource is written at the pinned ARM api-version 2025-06-01 -- the
# same GA line the ML family's verification pins -- and the kind's spec
# carries the full validation burden (azapi has no provider-side schema).
#
# The pin is EXACT (never pessimistic): a floating raw-API provider would
# move the ARM client under every admitted kind at once. When azurerm
# ships native ML endpoint resources, this module migrates azapi -> native
# in the next minor release (state move / re-import) -- the raw-API
# mechanism is an admission with a mandatory exit, not a new baseline.
terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "2.11.0"
    }
  }
}

provider "azapi" {}
