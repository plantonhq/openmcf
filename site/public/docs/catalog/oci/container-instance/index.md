---
title: "Container Instance"
description: "Container Instance deployment documentation"
icon: "package"
order: 100
componentName: "ocicontainerinstance"
---

# Container Instance on OCI

Deploys an OCI Container Instance -- a serverless container runtime that runs one or more containers in a pod-like construct sharing networking and volumes, without managing underlying compute infrastructure. Supports multi-container sidecar patterns, HTTP/TCP health checks, security contexts (non-root enforcement, read-only rootfs, Linux capabilities), and two volume types (emptydir and configfile). Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for compartment, subnet, and security group wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Container Instance** -- a serverless container runtime in the specified compartment and availability domain, with the configured shape, CPU/memory allocation, containers, VNICs, volumes, and restart policy
- **Containers** -- one or more containers within the instance sharing the same network namespace. Each container has its own image, environment variables, health checks, security context, and volume mounts
- **VNICs** -- virtual network interface cards providing network connectivity. All containers share the instance's VNICs and can communicate over localhost
- **Volumes** -- created only when `volumes` is populated; emptydir (ephemeral or tmpfs-backed) and configfile (inline file injection) volumes mounted into containers
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the container instance

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the container instance in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for at least one VNIC. Public subnets allow public IP assignment for direct access; private subnets require a load balancer or bastion for inbound traffic. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef.
- Container images accessible from OCI (public registries like Docker Hub/GHCR, or private registries with image pull secrets configured).

## Deploy

### Console

Open the deployment store, find **Container Instance on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Service** preset in the [Presets](#presets) tab to pre-populate a single-container instance with an HTTP health check.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciContainerInstance
metadata:
  name: api-service
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Ixxj:US-ASHBURN-AD-1"
  shape: CI.Standard.E4.Flex
  shapeConfig:
    ocpus: 1
    memoryInGbs: 2
  containers:
    - imageUrl: "ghcr.io/acme/api:v1.0"
      displayName: api
  vnics:
    - subnetId:
        value: "ocid1.subnet.oc1..example"
```

```shell
planton apply -f container-instance.yaml
```

This creates a single-container instance with 1 OCPU and 2 GiB memory, using the default restart policy (ALWAYS). No health checks, security context, or volumes are configured. The container inherits DNS from the subnet's DHCP options.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the container instance to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
  vnics:
    - subnetId:
        valueFrom:
          kind: OciSubnet
          name: private-app-subnet
          fieldPath: status.outputs.subnetId
      nsgIds:
        - valueFrom:
            kind: OciSecurityGroup
            name: app-nsg
            fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the container instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring a container instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Shape and resource allocation** -- Container Instance shapes are always flex (`CI.Standard.E4.Flex`, `CI.Standard.E3.Flex`). The `shapeConfig` sets the CPU and memory envelope for the entire instance. Individual containers can set `resourceConfig` limits within this envelope, or leave them unset to share all available resources.

**Multi-container patterns** -- Multiple containers in the `containers` list share the same network namespace (communicate over localhost) and can share volumes via volume mounts. Use this for sidecar patterns: log forwarders, metrics exporters, reverse proxies, or config reloaders alongside the main application container.

**Health checks** -- Each container supports HTTP and TCP health checks with configurable thresholds, intervals, and failure actions. Set `failureAction: kill` to restart unhealthy containers automatically (when `containerRestartPolicy` is `always` or `on_failure`). Health checks are the primary mechanism for automated container lifecycle management.

**Security context** -- The `securityContext` block enables non-root enforcement (`isNonRootUserCheckEnabled`), read-only root filesystem (`isRootFileSystemReadonly`), explicit UID/GID (`runAsUser`/`runAsGroup`), and Linux capability management (`capabilities.dropCapabilities: ["ALL"]`). Apply these for production workloads to minimize the container's attack surface.

**Restart policy** -- The `containerRestartPolicy` controls whether containers restart after exit: `always` (default, suitable for long-running services), `on_failure` (retry on non-zero exit), or `never` (one-shot jobs). Combined with `gracefulShutdownTimeoutInSeconds`, this controls the complete container lifecycle.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `vnics.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `vnics.nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `container_instance_id` | OCID of the container instance | Monitoring dashboards, OCI CLI operations, log queries |
| `container_ids` | Comma-separated OCIDs of individual containers | Per-container log retrieval, exec operations |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web service** -- A single-container instance with a public IP, HTTP health check, and automatic restart. The standard pattern for stateless web servers and API backends that need a public endpoint without a load balancer. Start from the **Web Service** preset.

**Private hardened** -- A security-hardened container in a private subnet with NSG association, non-root enforcement, read-only root filesystem, all capabilities dropped, and graceful shutdown timeout. The production pattern for application containers behind a load balancer. Start from the **Private Hardened** preset.

**Multi-container sidecar** -- Two containers sharing an emptydir volume and a configfile volume: a main application writing logs and a log-forwarder sidecar reading them. Demonstrates the sidecar pattern with per-container resource limits and shared volumes. Start from the **Multi-Container Sidecar** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this container instance
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet for VNIC network connectivity
- [**Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security groups for fine-grained traffic control on VNICs