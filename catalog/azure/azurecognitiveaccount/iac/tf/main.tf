# Create the Azure AI services account -- the container Azure AI
# capabilities are provisioned and billed through (Azure OpenAI, the
# multi-service AIServices account, the single-service accounts). The
# account owns the endpoint, keys, network perimeter, and
# responsible-AI policy; model deployments and AI Foundry projects are
# separate kinds that reference it.
#
# The spec's CEL contracts already enforce every kind-gated rule the
# provider checks at apply time (project management / network
# injection only on AIServices, bypass only on the AI kinds, the
# QnAMaker / TextAnalytics / MetricsAdvisor field gates) -- by the
# time this module runs, the shape is legal.
#
# Deletion is a SOFT delete: the account becomes a purgeable ghost
# that keeps holding its name. The provider purges it when the
# features flag `cognitive_account.purge_soft_delete_on_destroy` is
# enabled (the provider's own default).
resource "azurerm_cognitive_account" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Both vocabularies are already wire values in the spec.
  kind     = var.spec.kind
  sku_name = var.spec.sku_name

  # Only legal on kind "AIServices" (spec CEL); lets
  # AzureCognitiveAccountProject resources be created on the account.
  project_management_enabled = var.spec.project_management_enabled

  # Set-once: the provider replaces the account only when CHANGING an
  # existing subdomain, not when adding one to an account without it.
  custom_subdomain_name = var.spec.custom_subdomain_name != "" ? var.spec.custom_subdomain_name : null

  dynamic "customer_managed_key" {
    for_each = var.spec.customer_managed_key != null ? [var.spec.customer_managed_key] : []
    content {
      # A Key Vault key data-plane URL; versionless follows rotation.
      key_vault_key_id   = customer_managed_key.value.key_vault_key_id
      identity_client_id = customer_managed_key.value.identity_client_id != "" ? customer_managed_key.value.identity_client_id : null
    }
  }

  dynamic_throttling_enabled = var.spec.dynamic_throttling_enabled

  # The outbound FQDN allowlist (with outbound_network_access_restricted).
  fqdns = length(var.spec.fqdns) > 0 ? var.spec.fqdns : null

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # Optional-with-default-true on the provider: emit null when the spec
  # leaves it unset so the provider default applies.
  local_auth_enabled = var.spec.local_auth_enabled

  # MetricsAdvisor-kind only (spec CEL); all four are ForceNew.
  metrics_advisor_aad_client_id   = var.spec.metrics_advisor_aad_client_id != "" ? var.spec.metrics_advisor_aad_client_id : null
  metrics_advisor_aad_tenant_id   = var.spec.metrics_advisor_aad_tenant_id != "" ? var.spec.metrics_advisor_aad_tenant_id : null
  metrics_advisor_super_user_name = var.spec.metrics_advisor_super_user_name != "" ? var.spec.metrics_advisor_super_user_name : null
  metrics_advisor_website_name    = var.spec.metrics_advisor_website_name != "" ? var.spec.metrics_advisor_website_name : null

  dynamic "network_acls" {
    for_each = var.spec.network_acls != null ? [var.spec.network_acls] : []
    content {
      # Already a wire value ("Allow"/"Deny") in the spec.
      default_action = network_acls.value.default_action
      ip_rules       = length(network_acls.value.ip_rules) > 0 ? network_acls.value.ip_rules : null

      dynamic "virtual_network_rules" {
        for_each = network_acls.value.virtual_network_rules
        content {
          subnet_id                            = virtual_network_rules.value.subnet_id
          ignore_missing_vnet_service_endpoint = virtual_network_rules.value.ignore_missing_vnet_service_endpoint
        }
      }

      # Enum name -> wire value; unspecified ("") omits the property so
      # ARM applies its default. Only legal on the AI kinds (spec CEL).
      bypass = lookup(local.network_acls_bypass_wire, network_acls.value.bypass, null)
    }
  }

  # AIServices-kind only (spec CEL): inject agent workloads into the
  # given delegated subnet. NOTE: after the account deletes, ARM
  # removes the subnet's service association link asynchronously --
  # the provider waits for that before finishing the destroy.
  dynamic "network_injection" {
    for_each = var.spec.network_injection != null ? [var.spec.network_injection] : []
    content {
      scenario  = network_injection.value.scenario
      subnet_id = network_injection.value.subnet_id
    }
  }

  outbound_network_access_restricted = var.spec.outbound_network_access_restricted

  # Optional-with-default-true on the provider: emit null when unset.
  public_network_access_enabled = var.spec.public_network_access_enabled

  # QnAMaker-kind only (spec CEL).
  qna_runtime_endpoint = var.spec.qna_runtime_endpoint != "" ? var.spec.qna_runtime_endpoint : null

  # TextAnalytics-kind only (spec CEL). The id stays a plain string
  # until AzureSearchService registers (recorded in-place upgrade);
  # the key is sensitive -- resolved from a secret reference, masked
  # in plan output.
  custom_question_answering_search_service_id  = var.spec.custom_question_answering_search_service_id != "" ? var.spec.custom_question_answering_search_service_id : null
  custom_question_answering_search_service_key = var.spec.custom_question_answering_search_service_key != "" ? var.spec.custom_question_answering_search_service_key : null

  dynamic "storage" {
    for_each = var.spec.storage
    content {
      storage_account_id = storage.value.storage_account_id
      identity_client_id = storage.value.identity_client_id != "" ? storage.value.identity_client_id : null
    }
  }

  tags = local.final_tags
}

# The composed responsible-AI blocklists: standalone ARM children of
# the account, one per spec entry, keyed by name (named containers for
# custom blocked content; their ITEMS are data-plane). The
# rai_blocklist_ids output publishes each blocklist's ARM id.
resource "azurerm_cognitive_account_rai_blocklist" "rai_blocklists" {
  for_each = { for blocklist in var.spec.rai_blocklists : blocklist.name => blocklist }

  name                 = each.value.name
  cognitive_account_id = azurerm_cognitive_account.main.id
  description          = each.value.description != "" ? each.value.description : null
  tags                 = length(each.value.tags) > 0 ? each.value.tags : null
}

# The composed responsible-AI (content-filter) policies: standalone
# ARM children of the account, one per spec entry, keyed by name.
# Model deployments select a policy by NAME via their rai_policy_name.
# A policy's content filter may reference a blocklist defined in the
# same spec by its name -- hence the explicit dependency.
resource "azurerm_cognitive_account_rai_policy" "rai_policies" {
  for_each = { for policy in var.spec.rai_policies : policy.name => policy }

  name                 = each.value.name
  cognitive_account_id = azurerm_cognitive_account.main.id
  base_policy_name     = each.value.base_policy_name

  dynamic "content_filter" {
    for_each = each.value.content_filters
    content {
      name           = content_filter.value.name
      filter_enabled = content_filter.value.filter_enabled
      block_enabled  = content_filter.value.block_enabled
      # Already a wire value in the spec.
      source = content_filter.value.source
      # Enum name -> wire value; unspecified ("") omits the property
      # (the binary filters carry no severity -- spec CEL enforces it).
      # PARITY-EXCEPTION: omitting the severity is legal on THIS
      # engine only -- the classic Pulumi SDK bridges the pre-v5
      # provider where the property is required on every filter, so
      # severity-less filters (all binary filters included) deploy
      # through Terraform only until a v5-bridged pulumi-azure major.
      severity_threshold = lookup(local.rai_content_level_wire, content_filter.value.severity_threshold, null)
    }
  }

  # Enum name -> wire value; unspecified ("") omits the property so
  # ARM applies its default.
  mode = lookup(local.rai_policy_mode_wire, each.value.mode, null)

  tags = length(each.value.tags) > 0 ? each.value.tags : null

  depends_on = [azurerm_cognitive_account_rai_blocklist.rai_blocklists]
}
