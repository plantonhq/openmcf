# Kubernetes ConfigMap - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes ConfigMap. It supports the complete ConfigMap surface: UTF-8 `data` entries, base64-encoded `binary_data` entries, and the `immutable` flag.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, namespace default, data maps
├── main.tf         # Creates kubernetes_config_map_v1 resource
├── outputs.tf      # Exports configmap_name, namespace
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable accepts `data` (UTF-8 strings), `binary_data` (base64 strings), and `immutable`
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels
4. **Resource Creation**: `main.tf` creates a single `kubernetes_config_map_v1` resource
5. **Output Export**: ConfigMap name and namespace are exported

## Data Handling

| Variable Field | ConfigMap Field | Encoding |
|---------------|----------------|----------|
| `data` | `data` | Plain UTF-8 strings, passed through |
| `binary_data` | `binaryData` | Base64 strings, passed through (Kubernetes stores binaryData as base64) |

## Usage

```hcl
module "configmap" {
  source = "./iac/tf"

  metadata = {
    name = "app-config"
  }

  spec = {
    name      = "app-config"
    namespace = "production"

    data = {
      "application.properties" = "log.level=INFO"
    }

    binary_data = {
      "logo.png" = "iVBORw0KGgo..." # already base64-encoded
    }

    immutable = false
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | ConfigMap specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `configmap_name` | Name of the created configmap |
| `namespace` | Namespace of the created configmap |
