# Overview

The **Azure NAT Gateway API Resource** provides a consistent and standardized interface for deploying and managing Azure NAT Gateway resources within our infrastructure. This resource solves the critical problem of SNAT port exhaustion in private Azure workloads by providing dynamic, scalable outbound internet connectivity.

## Purpose

We developed this API resource to streamline the deployment of Azure NAT Gateway for production workloads. NAT Gateway eliminates the common failure mode of SNAT port exhaustion that plagues default Azure outbound connectivity, enabling:

- **Dynamic SNAT**: On-demand port allocation from a shared pool, eliminating pre-allocation inefficiency
- **Massive Scale**: Up to 1 million+ SNAT ports per subnet using public IP prefixes
- **Predictable Egress**: Static public IP addresses for firewall allow-listing
- **Zero Configuration for Workloads**: Automatic outbound routing for all subnet resources

## Key Features

- **Consistent Interface**: Aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composition by Reference**: The gateway SNATs through first-class AzurePublicIp and AzurePublicIpPrefix resources (`publicIpIds` / `publicIpPrefixIds`), keeping egress addresses visible in the resource graph, allowlistable, and reusable
- **Subnet-Side Attachment**: Subnets attach themselves via AzureSubnet's `natGatewayId`, matching Azure's model where one gateway serves many subnets
- **SKU Choice**: `STANDARD` for zonal deployments (optionally pinned to an availability zone) or `STANDARD_V2` for automatic zone redundancy
- **AKS Integration**: Perfect for providing predictable outbound connectivity for private AKS clusters

## Use Cases

- **Private AKS Clusters**: Eliminate SNAT port exhaustion for high-connection workloads running in Kubernetes
- **VM Scale Sets**: Provide reliable outbound connectivity for autoscaling VM fleets
- **Microservices**: Support high-volume API calls from private application subnets
- **CI/CD Pipelines**: Predictable egress IPs for build agents that need to be allow-listed
- **Compliance Requirements**: Static egress IPs for auditable outbound traffic patterns
- **Multi-Region Deployments**: Consistent NAT Gateway configuration across Azure regions

## Future Enhancements

Future updates will include:

- **Monitoring Integration**: Built-in Azure Monitor metrics and alerts for port utilization
- **Cost Optimization**: Automatic configuration of Private Link bypass patterns
- **Multi-IP Management**: Dynamic scaling of public IP count based on load
- **Comprehensive Documentation**: Expanded troubleshooting guides and performance tuning recommendations
