# AzureContainerInstance Terraform Module

## Overview

Creates an Azure Container Instance container group -- serverless containers billed per second: one or more containers sharing a lifecycle, network, and volumes, plus optional one-shot init containers.

## Resources Created

- `azurerm_container_group` -- the group (containers with probes and volumes, init containers, registry credentials, identity, Log Analytics diagnostics, custom DNS, network posture, CMK encryption, tags)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureContainerInstanceSpec fields; the resource-group, subnet, registry, identity, workspace, storage, and Key Vault references arrive as resolved literals

## Outputs

- `container_group_id` -- the group's ARM resource ID
- `container_group_name` -- the group's name
- `ip_address` -- the group's IP (public or private per posture; empty for "None")
- `fqdn` -- the DNS-label FQDN (empty unless a public group sets dns_name_label)
- `identity_principal_id` / `identity_tenant_id` -- the system-assigned identity (empty unless enabled)

## Behavior Notes

- **Near-total ForceNew**: Azure applies only identity and tag changes in place -- any other change replaces the group. Additionally, `cpu_limit`, `memory_limit`, and `key_vault_user_assigned_identity_id` are accepted by the provider on updates but SILENTLY NEVER APPLIED (its update path covers only identity and tags) -- treat all three as create-only in practice.
- **The volume union flattens here**: the spec validates exactly one of azure_file / empty_dir / git_repo / secret per volume before apply; this module maps the chosen form onto the provider's flat volume block. An empty_dir with the same name in several containers is ONE shared scratch volume.
- **Write-only secrets**: Azure never returns `secure_environment_variables`, volume `storage_account_key` / `secret`, registry `password`, or the Log Analytics `workspace_key` on reads -- the provider echoes them from configuration/state, so they must stay present in the manifest.
- **Probe http_get is singular on the wire**: the provider schema accepts a list but its code keeps only the last entry; the spec models one and this module renders a one-element list.
- **exposed_port omitted = expose all container ports** (the provider derives the group ports). When set, every entry must match a container's port+protocol -- validated in the spec, enforced by the provider at apply.
- **Unset optionals ride the provider defaults**: ip_address_type "Public", restart_policy "Always", sku "Standard", dns_name_label_reuse_policy "Unsecure", port protocol "TCP", probe timings Azure-side.
- **Subnet serialization**: the provider locks the subnet during create -- container groups sharing a subnet deploy one at a time.
