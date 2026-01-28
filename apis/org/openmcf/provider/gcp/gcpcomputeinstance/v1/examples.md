# GCP Compute Instance API - Examples

Here are examples of how to create and configure a **GcpComputeInstance** API resource using the OpenMCF CLI. The examples cover various instance configurations from basic to production-grade setups.

## Create using CLI

First, create a YAML file using the examples provided below. After the YAML file is created, you can apply the configuration using the following command:

```shell
openmcf pulumi up --manifest <yaml-path> --stack <org>/<stack-name>/<environment>
```

Or using Terraform:

```shell
openmcf tofu apply --manifest <yaml-path> --auto-approve
```

## Basic Example (Literal Value)

This example demonstrates how to create a basic Compute Engine instance with a hardcoded project ID.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: my-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: e2-medium
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 20
  networkInterfaces:
    - network:
        value: default
      accessConfigs:
        - networkTier: PREMIUM
```

## Basic Example (Value From Reference)

This example demonstrates referencing a GcpProject resource for the project ID.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: my-vm
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  zone: us-central1-a
  machineType: e2-medium
  bootDisk:
    image: debian-cloud/debian-11
  networkInterfaces:
    - network:
        value: default
```

## Ubuntu Instance with SSH Access

This example creates an Ubuntu instance with SSH key configuration.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: ubuntu-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: e2-standard-2
  bootDisk:
    image: ubuntu-os-cloud/ubuntu-2204-lts
    sizeGb: 50
    type: pd-ssd
  networkInterfaces:
    - network:
        value: default
      accessConfigs:
        - networkTier: PREMIUM
  sshKeys:
    - "myuser:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDExample..."
  metadata:
    enable-oslogin: "TRUE"
```

## Instance with Custom VPC (Literal Values)

This example creates an instance connected to a custom VPC network using literal values.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: vpc-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: n1-standard-1
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 20
  networkInterfaces:
    - network:
        value: projects/my-gcp-project/global/networks/my-vpc
      subnetwork:
        value: projects/my-gcp-project/regions/us-central1/subnetworks/my-subnet
      accessConfigs:
        - networkTier: PREMIUM
```

## Instance with Custom VPC (Value From References)

This example creates an instance by referencing existing GcpVpc and GcpSubnetwork resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: vpc-vm
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  zone: us-central1-a
  machineType: n1-standard-1
  bootDisk:
    image: debian-cloud/debian-11
  networkInterfaces:
    - network:
        valueFrom:
          kind: GcpVpc
          name: main-vpc
          fieldPath: status.outputs.network_self_link
      subnetwork:
        valueFrom:
          kind: GcpSubnetwork
          name: main-subnet
          fieldPath: status.outputs.subnetwork_self_link
```

## Spot VM for Cost Optimization

This example creates a Spot VM for fault-tolerant workloads at reduced cost.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: spot-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: e2-standard-4
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 20
  networkInterfaces:
    - network:
        value: default
  spot: true
  scheduling:
    preemptible: true
    automaticRestart: false
    onHostMaintenance: TERMINATE
    provisioningModel: SPOT
    instanceTerminationAction: STOP
```

## Instance with Service Account (Literal Value)

This example creates an instance with a custom service account using a literal email value.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: sa-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: e2-medium
  bootDisk:
    image: debian-cloud/debian-11
  networkInterfaces:
    - network:
        value: default
  serviceAccount:
    email:
      value: my-sa@my-gcp-project.iam.gserviceaccount.com
    scopes:
      - https://www.googleapis.com/auth/cloud-platform
```

## Instance with Service Account (Value From Reference)

This example creates an instance by referencing a GcpServiceAccount resource.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: sa-vm
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  zone: us-central1-a
  machineType: e2-medium
  bootDisk:
    image: debian-cloud/debian-11
  networkInterfaces:
    - network:
        value: default
  serviceAccount:
    email:
      valueFrom:
        kind: GcpServiceAccount
        name: my-service-account
        fieldPath: status.outputs.email
    scopes:
      - https://www.googleapis.com/auth/cloud-platform
```

## Web Server with Startup Script

This example creates a web server with a startup script to install and configure nginx.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: web-server
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: e2-small
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 20
  networkInterfaces:
    - network:
        value: default
      accessConfigs:
        - networkTier: PREMIUM
  tags:
    - http-server
    - https-server
  startupScript: |
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
```

## Instance with Attached Data Disk

This example creates an instance with an additional data disk attached.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: data-vm
spec:
  projectId:
    value: my-gcp-project
  zone: us-central1-a
  machineType: n1-standard-2
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 20
  networkInterfaces:
    - network:
        value: default
  attachedDisks:
    - source: projects/my-gcp-project/zones/us-central1-a/disks/data-disk
      deviceName: data-disk
      mode: READ_WRITE
```

## Production-Grade Instance (Literal Values)

This comprehensive example creates a production-ready instance with all security and performance options configured.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: prod-server
spec:
  projectId:
    value: prod-project
  zone: us-central1-a
  machineType: n2-standard-4
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 100
    type: pd-ssd
    autoDelete: true
  networkInterfaces:
    - network:
        value: projects/prod-project/global/networks/prod-vpc
      subnetwork:
        value: projects/prod-project/regions/us-central1/subnetworks/prod-subnet
      accessConfigs:
        - networkTier: PREMIUM
  serviceAccount:
    email:
      value: prod-sa@prod-project.iam.gserviceaccount.com
    scopes:
      - https://www.googleapis.com/auth/cloud-platform
  deletionProtection: true
  allowStoppingForUpdate: true
  labels:
    env: production
    app: webserver
    team: platform
  tags:
    - web-server
    - https-server
  metadata:
    enable-oslogin: "TRUE"
  scheduling:
    automaticRestart: true
    onHostMaintenance: MIGRATE
    provisioningModel: STANDARD
  startupScript: |
    #!/bin/bash
    echo "Starting production server"
```

## Production-Grade Instance (Value From References)

This comprehensive example uses references to existing resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: prod-server
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: prod-project
      fieldPath: status.outputs.project_id
  zone: us-central1-a
  machineType: n2-standard-4
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 100
    type: pd-ssd
  networkInterfaces:
    - network:
        valueFrom:
          kind: GcpVpc
          name: prod-vpc
          fieldPath: status.outputs.network_self_link
      subnetwork:
        valueFrom:
          kind: GcpSubnetwork
          name: prod-subnet
          fieldPath: status.outputs.subnetwork_self_link
  serviceAccount:
    email:
      valueFrom:
        kind: GcpServiceAccount
        name: prod-service-account
        fieldPath: status.outputs.email
    scopes:
      - https://www.googleapis.com/auth/cloud-platform
  deletionProtection: true
  labels:
    env: production
```

## Development Instance (Cost Optimized)

This example creates a small, cost-effective instance for development.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpComputeInstance
metadata:
  name: dev-vm
spec:
  projectId:
    value: dev-project
  zone: us-central1-a
  machineType: e2-micro
  bootDisk:
    image: debian-cloud/debian-11
    sizeGb: 10
    type: pd-standard
  networkInterfaces:
    - network:
        value: default
  labels:
    env: development
```

## Notes

- **Zone**: Choose a zone close to your users or other resources for lower latency.
- **Machine Type**: Select based on CPU and memory requirements. E2 series is cost-effective for general workloads.
- **Boot Disk**: Use SSD (`pd-ssd`) for production, standard (`pd-standard`) for development.
- **Spot VMs**: Save up to 60-91% but can be preempted with 30-second notice.
- **Service Account**: Always use least-privilege service accounts in production.
- **Network Tags**: Used for firewall rule targeting - create corresponding firewall rules.
- **Labels**: Use for organization, cost allocation, and filtering.
- **Deletion Protection**: Enable for production instances to prevent accidental deletion.

## StringValueOrRef Fields

The following fields support both literal values and references to other resources:

### projectId
Use a literal value or reference a GcpProject resource:
```yaml
# Literal value
projectId:
  value: my-gcp-project

# Reference to GcpProject
projectId:
  valueFrom:
    kind: GcpProject
    name: my-project
    fieldPath: status.outputs.project_id
```

### networkInterfaces[].network
Use a literal value or reference a GcpVpc resource:
```yaml
# Literal value
network:
  value: projects/my-project/global/networks/my-vpc

# Reference to GcpVpc
network:
  valueFrom:
    kind: GcpVpc
    name: my-vpc
    fieldPath: status.outputs.network_self_link
```

### networkInterfaces[].subnetwork
Use a literal value or reference a GcpSubnetwork resource:
```yaml
# Literal value
subnetwork:
  value: projects/my-project/regions/us-central1/subnetworks/my-subnet

# Reference to GcpSubnetwork
subnetwork:
  valueFrom:
    kind: GcpSubnetwork
    name: my-subnet
    fieldPath: status.outputs.subnetwork_self_link
```

### serviceAccount.email
Use a literal value or reference a GcpServiceAccount resource:
```yaml
# Literal value
serviceAccount:
  email:
    value: my-sa@my-project.iam.gserviceaccount.com
  scopes:
    - https://www.googleapis.com/auth/cloud-platform

# Reference to GcpServiceAccount
serviceAccount:
  email:
    valueFrom:
      kind: GcpServiceAccount
      name: my-service-account
      fieldPath: status.outputs.email
  scopes:
    - https://www.googleapis.com/auth/cloud-platform
```
