# AzureVpnServerConfiguration Terraform Module

## Overview

This Terraform module provisions a VPN Server Configuration using the
`azurerm` provider. It creates one
`azurerm_vpn_server_configuration` -- the reusable point-to-site
authentication policy (Entra ID / certificate / RADIUS, trusted
certificates, tunnel protocols) -- plus one
`azurerm_vpn_server_configuration_policy_group` per `policy_groups`
entry.

## Resources Created

- `azurerm_vpn_server_configuration.main` -- the authentication policy
  object
- `azurerm_vpn_server_configuration_policy_group.policy_groups` -- one
  ARM child per spec entry, keyed by the group's name (`for_each`)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | VPN Server Configuration specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | The configuration object's region (ForceNew) |
| `resource_group` | yes | Resource group name (ForceNew) |
| `name` | yes | The configuration's name (ForceNew) |
| `vpn_authentication_types` | yes | "AAD" / "Certificate" / "Radius" -- each enabled type requires its block |
| `aad_authentication` | when AAD | Entra ID audience/issuer/tenant trio |
| `client_root_certificates` | when Certificate | Trusted roots (at least one) |
| `radius` | when Radius | RADIUS servers and trust anchors |
| `ipsec_policy` | no | Pinned client IPsec proposal (all fields required when set) |
| `vpn_protocols` | no | "IkeV2" / "OpenVPN"; empty applies ARM's default |
| `policy_groups` | no | Named member-matching rules (ARM children) |

## Outputs

| Output | Description |
|--------|-------------|
| `vpn_server_configuration_id` | Full ARM ID -- what a point-to-site gateway's `vpn_server_configuration_id` references |
| `vpn_server_configuration_name` | The configuration's name |
| `policy_group_ids` | Each policy group's ARM ID keyed by group name |

## Usage

```hcl
module "remote_workforce" {
  source = "./iac/tf"

  metadata = { name = "remote-workforce", org = "mycompany", env = "production" }

  spec = {
    name                     = "remote-workforce"
    region                   = "eastus"
    resource_group           = "network-rg"
    vpn_authentication_types = ["AAD"]
    aad_authentication = {
      audience = "41b23e61-6c1e-4545-b367-cd054e0ed4b4"
      issuer   = "https://sts.windows.net/00000000-0000-0000-0000-000000000000/"
      tenant   = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000"
    }
    vpn_protocols = ["OpenVPN"]
  }
}
```

## Behavior Notes

- The spec's CEL contracts guarantee each enabled authentication type
  brings its block (AAD → aad_authentication, Certificate →
  client_root_certificates, Radius → radius) -- the provider enforces
  the same three rules at apply time.
- `vpn_protocols` is Optional+Computed on the provider: the module
  emits null when the spec leaves it empty so ARM's default selection
  applies without read drift.
- The RADIUS server `secret` is Sensitive and ARM never returns it on
  reads -- the import map declares the matching tolerance.
- Policy groups deploy as standalone ARM children keyed by name;
  `name` and `is_default` are ForceNew on each group.
