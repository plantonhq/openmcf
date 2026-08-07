##############################################
# main.tf
#
# Main orchestration file for installing the
# Istio base CRDs.
#
# This module fetches and applies the official
# istio/base crd-all.gen.yaml bundle from the
# istio/istio GitHub repository (CRDs only --
# no istiod, no controller).
#
# Resources Created:
#  1. Istio CRDs (cluster-scoped)
##############################################

##############################################
# 1. Fetch Istio base CRD bundle
#
# Downloads the crd-all.gen.yaml bundle from the
# pinned istio/istio release ref.
##############################################
data "http" "istio_base_crds" {
  url = local.manifest_url

  request_headers = {
    Accept = "application/yaml"
  }
}

##############################################
# 2. Apply Istio CRDs
#
# The Istio CRDs are cluster-scoped resources
# that enable DestinationRule, ServiceEntry,
# PeerAuthentication, RequestAuthentication,
# AuthorizationPolicy, Telemetry, EnvoyFilter,
# and the rest of the Istio API.
#
# Applied with kubectl_manifest (server-side
# apply) which handles multi-document YAML and
# the large CRD schemas Istio ships.
##############################################
resource "kubectl_manifest" "istio_base_crds" {
  # Keyed by each CRD's OWN NAME (never the split index): the name is the
  # document's identity, so state addresses stay stable across bundle
  # reorderings AND the address key feeds the composed import ID blind
  # (from_address_key in the import map).
  for_each = {
    for doc in split("---", data.http.istio_base_crds.response_body) :
    yamldecode(doc).metadata.name => doc
    if trimspace(doc) != "" && can(yamldecode(doc).metadata.name)
  }

  yaml_body = each.value

  server_side_apply = true
  force_conflicts   = true
}
