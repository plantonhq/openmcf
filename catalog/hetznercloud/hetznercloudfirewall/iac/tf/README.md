# HetznerCloudFirewall Terraform Module

Terraform IaC module for creating firewalls with inline rules in Hetzner Cloud.

## Structure

```
.
├── main.tf           # Firewall resource with dynamic rule blocks
├── outputs.tf        # Stack output definitions
├── variables.tf      # Input variable definitions (metadata, spec, hcloud_token)
├── locals.tf         # Firewall name and standard label computation
├── provider.tf       # HetznerCloud provider configuration
└── BUILD.bazel       # Bazel build definition
```

## Resources Created

- `hcloud_firewall` — Firewall with dynamic `rule` blocks from `var.spec.rules`

## Outputs

| Name | Description |
|------|-------------|
| `firewall_id` | Hetzner Cloud numeric ID of the created firewall |

## Usage

```bash
# Build
bazel build //catalog/hetznercloud/hetznercloudfirewall/iac/tf:terraform_module

# Initialize
planton tofu init --manifest ../../e2e/manifest.yaml --module-dir .

# Plan
planton tofu plan --manifest ../../e2e/manifest.yaml --module-dir .

# Apply
planton tofu apply --manifest ../../e2e/manifest.yaml --module-dir . --auto-approve
```
