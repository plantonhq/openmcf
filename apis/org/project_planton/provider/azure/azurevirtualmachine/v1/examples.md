# Azure Virtual Machine Examples

Copy-paste ready examples for deploying Azure Virtual Machines with Project Planton.

## Minimal Linux VM with SSH

A minimal Ubuntu VM with SSH key authentication. Uses cross-reference to AzureVpc for subnet.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: dev-vm
spec:
  region: eastus
  resource_group: dev-rg
  subnet_id:
    value_from:
      kind: AzureVpc
      name: dev-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
```

## Production Linux VM with Availability Zone

A production-ready Ubuntu VM with zone redundancy and managed identity.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: prod-web-server
spec:
  region: eastus
  resource_group: prod-rg
  vm_size: Standard_D4s_v5
  availability_zone: "1"
  subnet_id:
    value_from:
      kind: AzureVpc
      name: prod-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
  enable_system_assigned_identity: true
  enable_boot_diagnostics: true
  tags:
    environment: production
    team: platform
```

## Windows Server VM

A Windows Server 2022 VM with password authentication and public IP.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: windows-server
spec:
  region: westus2
  resource_group: windows-rg
  vm_size: Standard_D4s_v5
  subnet_id:
    value: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet/subnets/default"
  image:
    publisher: MicrosoftWindowsServer
    offer: WindowsServer
    sku: 2022-datacenter-g2
  admin_username: adminuser
  admin_password:
    value: "SuperSecurePassword123!"
  network:
    enable_public_ip: true
    public_ip_sku: standard
    public_ip_allocation: static
```

## VM with Password from Key Vault

Reference password from an AzureKeyVault resource.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: secure-vm
spec:
  region: eastus
  resource_group: secure-rg
  vm_size: Standard_D2s_v3
  subnet_id:
    value_from:
      kind: AzureVpc
      name: secure-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  admin_password:
    value_from:
      kind: AzureKeyVault
      name: my-keyvault
      field_path: status.outputs.vault_uri
```

## VM with Data Disks

VM with additional data disks for storage.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: database-vm
spec:
  region: eastus
  resource_group: db-rg
  vm_size: Standard_D8s_v5
  subnet_id:
    value_from:
      kind: AzureVpc
      name: db-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
  os_disk:
    size_gb: 128
    storage_type: premium_lrs
    caching: read_write
  data_disks:
    - name: data-disk-1
      size_gb: 512
      lun: 0
      storage_type: premium_lrs
      caching: read_only
    - name: log-disk
      size_gb: 256
      lun: 1
      storage_type: premium_lrs
      caching: none
```

## Spot Instance for Cost Savings

Use spot pricing for dev/test or fault-tolerant workloads.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: spot-worker
spec:
  region: eastus
  resource_group: workers-rg
  vm_size: Standard_D8s_v5
  subnet_id:
    value_from:
      kind: AzureVpc
      name: workers-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
  is_spot_instance: true
  spot_max_price: 0.5  # Maximum hourly price in USD
  tags:
    workload: batch-processing
```

## Custom Image VM

Deploy a VM using a custom or shared image.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: custom-image-vm
spec:
  region: eastus
  resource_group: custom-rg
  vm_size: Standard_D4s_v5
  subnet_id:
    value_from:
      kind: AzureVpc
      name: custom-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    custom_image_id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/images-rg/providers/Microsoft.Compute/images/my-golden-image"
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
```

## VM with Accelerated Networking and Static Private IP

Configure advanced networking options.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureVirtualMachine
metadata:
  name: network-optimized-vm
spec:
  region: eastus
  resource_group: network-rg
  vm_size: Standard_D8s_v5
  subnet_id:
    value_from:
      kind: AzureVpc
      name: network-vpc
      field_path: status.outputs.nodes_subnet_id
  image:
    publisher: Canonical
    offer: 0001-com-ubuntu-server-jammy
    sku: 22_04-lts-gen2
  ssh_public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host"
  network:
    enable_accelerated_networking: true
    private_ip_allocation: private_static
    private_ip_address: "10.0.1.100"
```

## Deploy Commands

```bash
# Deploy with Pulumi
project-planton pulumi up --manifest vm.yaml

# Deploy with Terraform
project-planton terraform apply --manifest vm.yaml

# Preview changes before deploying
project-planton pulumi preview --manifest vm.yaml
```
