# Create the Azure Firewall -- the managed, stateful firewall data plane
# that enforces an attached firewall policy. The firewall carries WHERE
# enforcement runs (subnet, public IPs, zones, deployment model); the
# policy carries WHAT is enforced.
#
# Classic inline rule collections are deliberately not modeled --
# policy-based management is Azure's direction, and ARM rejects mixing the
# two on one firewall.
#
# Azure Firewall provisions and deletes SLOWLY (10-20+ minutes each way).
# The ForceNew surface (name, region, sku_name, zones, the management
# block, each ip configuration's subnet) is therefore expensive -- design
# changes to avoid replacement.
resource "azurerm_firewall" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Both ARM-required; the tier must match the attached policy's tier
  # (ARM validates the pairing at apply time).
  sku_name = local.sku_name
  sku_tier = local.sku_tier

  # Rules, threat intelligence, TLS inspection, and IDPS live on the
  # policy; the firewall is the enforcement point.
  firewall_policy_id = (
    var.spec.firewall_policy_id != null && var.spec.firewall_policy_id != ""
  ) ? var.spec.firewall_policy_id : null

  # Only meaningful without a policy; null lets Azure own its default
  # (Alert) instead of the module guessing it.
  threat_intel_mode = local.threat_intel_mode

  # DNS: setting servers implicitly forces the DNS proxy ON in Azure's
  # wire encoding (the provider couples them); dns_proxy_enabled alone
  # enables proxying without custom upstreams. Passed through verbatim so
  # the coupling stays Azure's, not the module's.
  dns_servers       = length(var.spec.dns_servers) > 0 ? var.spec.dns_servers : null
  dns_proxy_enabled = var.spec.dns_proxy_enabled ? true : null

  # SNAT ranges: CIDRs or the literal "IANAPrivateRanges" token; sent
  # only when the user overrides Azure's IANA-private default.
  private_ip_ranges = (
    length(var.spec.private_ip_ranges) > 0
  ) ? var.spec.private_ip_ranges : null

  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  # The data-path IP configurations: exactly one carries the
  # "AzureFirewallSubnet" subnet (spec-validated); extra blocks add
  # public IPs (each adds SNAT ports and a DNAT frontend). The subnet
  # name/size contract (/26+) is ARM's -- the referenced AzureSubnet must
  # be created to it.
  dynamic "ip_configuration" {
    for_each = var.spec.ip_configurations
    content {
      name = ip_configuration.value.name
      subnet_id = (
        ip_configuration.value.subnet_id != null && ip_configuration.value.subnet_id != ""
      ) ? ip_configuration.value.subnet_id : null
      public_ip_address_id = (
        ip_configuration.value.public_ip_address_id != null && ip_configuration.value.public_ip_address_id != ""
      ) ? ip_configuration.value.public_ip_address_id : null
    }
  }

  # The management path (forced tunneling / BASIC tier): a dedicated
  # "AzureFirewallManagementSubnet" (/26+) with its REQUIRED public IP.
  # ForceNew -- adding or removing this block replaces the firewall.
  dynamic "management_ip_configuration" {
    for_each = var.spec.management_ip_configuration != null ? [var.spec.management_ip_configuration] : []
    content {
      name                 = management_ip_configuration.value.name
      subnet_id            = management_ip_configuration.value.subnet_id
      public_ip_address_id = management_ip_configuration.value.public_ip_address_id
    }
  }

  # The Virtual WAN hub target (AZFW_HUB model only -- spec-validated
  # pairing). Azure manages the hub firewall's addressing.
  dynamic "virtual_hub" {
    for_each = var.spec.virtual_hub != null ? [var.spec.virtual_hub] : []
    content {
      virtual_hub_id  = virtual_hub.value.virtual_hub_id
      public_ip_count = virtual_hub.value.public_ip_count
    }
  }

  tags = local.final_tags
}
