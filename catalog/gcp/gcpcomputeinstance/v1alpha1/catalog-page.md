# GCP Compute Instance

Creates a Google Compute Engine virtual machine — the general-purpose compute node behind databases-on-VM, stateful servers, GPU workers, and bastion hosts — with the full single-instance surface: boot from image/snapshot/existing disk, first-class data-disk attachments, Spot provisioning, Shielded and Confidential VM, multi-NIC networking with static IPs by reference, and an in-place start/suspend/stop lever.

## What Gets Created

- The Compute Engine API is enabled on the project (never disabled on destroy)
- A `google_compute_instance` carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set)
- The instance name falls back to `metadata.name` when `instanceName` is not set; SSH keys fold into the instance metadata `ssh-keys` key

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **IAM permissions** — `roles/compute.instanceAdmin.v1` on the target project; `roles/iam.serviceAccountUser` on any attached service account
- For CMEK boot disks: a `GcpKmsKey` the Compute Engine service agent can use

## Quick Start

Create a file `vm.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpComputeInstance
metadata:
  name: dev-vm
spec:
  zone: us-central1-a
  machineType: e2-medium
  bootDisk:
    image: debian-cloud/debian-12
    type: pd-balanced
  networkInterfaces:
    - network:
        value: default
      accessConfigs:
        - {}
  scheduling:
    provisioningModel: SPOT
    instanceTerminationAction: DELETE
```

Deploy:

```shell
planton apply -f vm.yaml
```

This creates a Spot Debian 12 VM on the default network with an ephemeral external IP — deeply discounted capacity that GCP may reclaim, ideal for dev and interruption-tolerant work.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Owning project (reference a `GcpProject`); empty uses the provider default. Immutable |
| `instanceName` | string | No | Falls back to `metadata.name`. Immutable |
| `zone` | string | Yes | e.g. `us-central1-a`. Immutable |
| `machineType` | string | Yes | e.g. `e2-medium`, `n2-standard-4`. Changing it stops/restarts the VM |
| `hostname` | string | No | Custom FQDN. Immutable |
| `bootDisk` | object | Yes | Exactly one source: `image`, `sourceSnapshot`, or `sourceDisk` (`GcpComputeDisk` reference); plus size/type/`autoDelete`/CMEK |
| `attachedDisks[]` | list | No | First-class `GcpComputeDisk` references — disks keep their own lifecycle |
| `scratchDisks[]` | list | No | Ephemeral local SSDs (contents lost on stop). Create-time only |
| `networkInterfaces[]` | list | Yes | `network` (`GcpVpcNetwork`) or `subnetwork` (`GcpSubnetwork`); static internal `networkIp` and external `accessConfigs[].natIp` reference `GcpAddress`; IPv6, gVNIC, alias ranges |
| `serviceAccount` | object | No | `email` (`GcpServiceAccount` reference) + `scopes`. Omitted = Compute Engine default account |
| `scheduling` | object | No | Spot (`provisioningModel: SPOT`), maintenance (`MIGRATE`/`TERMINATE`), Spot termination action (`STOP`/`DELETE`), run-duration limits, sole tenancy |
| `shieldedInstanceConfig` | object | No | Secure boot, vTPM, integrity monitoring |
| `confidentialInstanceConfig` | object | No | SEV / SEV_SNP / TDX memory encryption; requires `onHostMaintenance: TERMINATE` |
| `guestAccelerators[]` | list | No | GPUs; require `onHostMaintenance: TERMINATE` |
| `advancedMachineFeatures` | object | No | Nested virtualization, SMT control, PMU, turbo mode |
| `reservationAffinity` | object | No | Reserved-capacity consumption mode |
| `totalEgressBandwidthTier` | string | No | `DEFAULT` or `TIER_1` |
| `metadata` / `startupScript` / `sshKeys[]` | — | No | Guest metadata, per-boot startup script, SSH keys (folded into `ssh-keys`) |
| `labels` / `tags[]` / `resourceManagerTags` | — | No | Labels (merged beneath platform labels), firewall network tags, org-policy tag bindings |
| `deletionProtection` | bool | No | Guards the VM object only (default false — data levers live on disks) |
| `desiredStatus` | string | No | `RUNNING` / `SUSPENDED` / `TERMINATED` — start/suspend/stop in place |
| `allowStoppingForUpdate` | bool | No | Allow stop/restart for updates that need it (recommended true) |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_name` | Name of the instance |
| `instance_id` | Server-assigned unique numeric identifier |
| `self_link` | API self link — what consumers reference |
| `internal_ip` | Primary internal IP (first interface) |
| `external_ip` | External IP of the first interface; `""` when private |
| `status` | Current status (RUNNING, TERMINATED, ...) |
| `zone` | Zone |
| `machine_type` | Machine type |
| `cpu_platform` | CPU platform the instance landed on |

## Related Resources

- **GcpComputeDisk** — boot-from-existing-disk and attached data disks with their own lifecycles
- **GcpAddress** — static internal (`networkIp`) and external (`accessConfigs[].natIp`) IPs
- **GcpVpcNetwork / GcpSubnetwork** — the networks each vNIC attaches to
- **GcpServiceAccount** — the VM's least-privilege workload identity
- **GcpKmsKey** — CMEK protection for boot and attached disks
