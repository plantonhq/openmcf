# GCP Vertex AI Workbench Instance -- Design Notes

## Deployment Landscape

### What Is Vertex AI Workbench?

Vertex AI Workbench is Google Cloud's managed notebook service. It provides JupyterLab instances running on Compute Engine VMs, pre-configured with popular ML frameworks, GPU drivers, and access to GCP data services. Workbench is the successor to AI Platform Notebooks (now deprecated).

### Resource Lineage

Google has iterated through three generations of managed notebooks:

1. **AI Platform Notebooks** (`google_notebooks_instance`) -- v1 API, deprecated
2. **Vertex AI Workbench Managed Notebooks** -- intermediate generation, also deprecated
3. **Vertex AI Workbench Instances** (`google_workbench_instance`) -- v2 API, current

Planton targets the current v2 API via `google_workbench_instance` (Terraform) and `workbench.Instance` (Pulumi).

## Methods of Deployment

### 1. Google Cloud Console

The Workbench UI in the Cloud Console provides a wizard for creating instances. It's the most accessible method but doesn't support infrastructure-as-code workflows.

### 2. gcloud CLI

```bash
gcloud workbench instances create my-notebook \
  --location=us-central1-a \
  --machine-type=e2-standard-4
```

Good for ad-hoc creation but lacks state management.

### 3. Terraform / Pulumi

```hcl
resource "google_workbench_instance" "notebook" {
  name     = "my-notebook"
  location = "us-central1-a"

  gce_setup {
    machine_type = "e2-standard-4"
  }
}
```

The standard IaC approach. Both of this component's modules follow this pattern.

### 4. Planton (This Component)

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVertexAiNotebook
metadata:
  name: my-notebook
spec:
  location: us-central1-a
  machineType: e2-standard-4
```

Planton provides a declarative YAML interface with cross-resource composition, merged user + platform labels, and dual IaC engine support. `project_id` is optional: when omitted, the instance lands in the provider's default project.

## Feature Coverage

The spec models the full released provider surface for Workbench instances:

| Feature | Coverage |
|---------|----------|
| Machine type selection | Full |
| GPU accelerators | Full (10 types) |
| Boot disk (type, size, CMEK) | Full |
| Data disk (type, size, CMEK) | Full |
| VPC networking (network, subnet, NIC type) | Full |
| Static external IP (access config) | Full -- references GcpAddress |
| Disable public IP | Full |
| Service account | Full |
| VM images (project, family, name) | Full |
| Container images (repo, tag) | Full |
| Shielded VM (Secure Boot, vTPM, integrity) | Full |
| Confidential Computing (AMD SEV) | Full |
| Reservation affinity (ANY / SPECIFIC / NONE) | Full |
| Managed end-user credentials (EUC) | Full |
| Third-party identity access | Full |
| Desired state (ACTIVE/STOPPED) | Full |
| Instance owners | Full |
| Metadata | Full |
| Network tags | Full |
| User labels | Full -- merged beneath platform labels |

### Deliberate Exclusions (with reasons)

| Provider surface | Reason |
|------------------|--------|
| `instance_id` | Vestigial in the provider: the create call derives the instance ID from `name`, so a second identity field would only invite drift. The spec's `instance_name` (falling back to `metadata.name`) is the single identity. |

### Immutable Fields (ForceNew)

These fields cannot be changed after creation without destroying and recreating the instance:

- `location`, `instance_name`, `disable_proxy_access`
- `network_interface` (network, subnet, nic_type, external_ip)
- `disable_public_ip`, `enable_ip_forwarding`
- `service_account`, `tags`
- `vm_image`, `container_image`
- `confidential_instance_config`, `reservation_affinity`
- `boot_disk.disk_type`, `boot_disk.kms_key`
- `data_disk.disk_type`, `data_disk.kms_key`

Disk sizes (boot and data), labels, desired_state, metadata, managed EUC, and third-party identity CAN change without recreation.

## Design Decisions

### 1. Flattened gce_setup

The Terraform/Pulumi providers nest all VM configuration under a `gce_setup` block. We flatten these to the top level of the spec because:

- The component IS a workbench instance -- the `gce_setup` wrapper adds no semantic value
- Matches the GcpComputeInstance pattern (boot_disk, network_interfaces at top level)
- Simpler YAML for users

### 2. Singular Sub-Messages

The providers use repeated fields for accelerator_configs, data_disks, network_interfaces, and service_accounts -- all capped at one entry by the API ("currently supports only one..."). We use singular messages for honesty and clarity:

- `accelerator_config` (not `accelerator_configs`)
- `data_disk` (not `data_disks`)
- `network_interface` (not `network_interfaces`)
- `service_account` (not `service_accounts`)

### 3. Int32 for Disk Sizes and Core Count

The providers use strings for `disk_size_gb` and `core_count`. We use `int32` because:

- Enables proto-level range validation (10-64000 for disks)
- Better developer experience (no quotes around numbers)
- IaC modules convert to strings internally

### 4. Derived disk_encryption

Instead of exposing a `disk_encryption` field, we derive it from the presence of `kms_key`:

- If `kms_key` is set → CMEK
- If `kms_key` is not set → GMEK (default)

This eliminates a redundant field and prevents inconsistent configurations.

### 5. Service Account as StringValueOrRef

Rather than a sub-message with `email` and `scopes` fields, we use a flat `StringValueOrRef` because:

- `scopes` is always `["cloud-platform"]` (computed, not configurable)
- The only user-facing field is the email address
- Flat StringValueOrRef enables direct infra-chart composition with GcpServiceAccount

### 6. External IP as a Reference

The provider's `access_configs` block takes a literal static IP string. The spec models it as a `StringValueOrRef` defaulting to GcpAddress's `address` output, so reserving the IP and pinning the notebook to it compose as two first-class nodes.

## Best Practices

### Cost Management

- Use `desired_state: STOPPED` to suspend notebooks when not in use
- Compute charges stop; storage charges continue
- GPU instances are expensive -- always stop when not training
- Point `reservation_affinity` at pre-purchased capacity to guarantee GPU availability without over-provisioning on demand

### Security

- Set `disable_public_ip: true` for production notebooks
- Configure a dedicated service account (don't use the default compute SA)
- Use CMEK encryption for regulated workloads
- Enable Shielded VM for defense-in-depth
- Use Confidential Computing (SEV) when notebook memory must stay encrypted in use -- pick an n2d machine type
- Enable managed EUC when notebook code should act as the signed-in user rather than the VM's service account

### Networking

- Deploy in the same VPC as your data sources (BigQuery, GCS, etc.)
- Use Private Google Access for accessing GCP APIs without public IP
- Apply network tags for firewall rule targeting
- Pin a static external IP (via GcpAddress) only when firewall allowlists or DNS records depend on the instance address
