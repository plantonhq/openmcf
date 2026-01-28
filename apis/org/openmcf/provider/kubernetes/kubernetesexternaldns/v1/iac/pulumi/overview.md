# Overview

The KubernetesExternalDns Pulumi module automates the deployment and configuration of ExternalDNS on Kubernetes clusters across multiple cloud providers. ExternalDNS synchronizes Kubernetes Service and Ingress resources with external DNS providers, eliminating manual DNS record management.

This module integrates seamlessly with cloud-native authentication mechanisms (IRSA for AWS, Workload Identity for GCP, Managed Identity for Azure, and API tokens for Cloudflare) by leveraging the Kubernetes cluster credentials provided in the `stack-input` to create a Kubernetes provider. It handles the complete setup: namespace creation, service account configuration with proper cloud provider annotations, secret management (for Cloudflare), and Helm chart deployment with provider-specific configuration.

The module supports four provider configurations:

- **GKE**: Uses Google Cloud DNS with Workload Identity for secure, keyless authentication
- **EKS**: Uses AWS Route53 with IRSA for IAM-based permissions without static credentials  
- **AKS**: Uses Azure DNS with Managed Identity for Azure AD-based authentication
- **Cloudflare**: Uses Cloudflare DNS with API tokens stored as Kubernetes secrets

Key features include:

- **Zone scoping**: Automatically limits ExternalDNS to manage only specified DNS zones for safety
- **Cloud-native security**: No long-lived credentials; uses IRSA, Workload Identity, or Managed Identity  
- **Multi-instance support**: Can deploy multiple ExternalDNS instances in the same cluster for different zones
- **Production defaults**: Upsert-only policy, proper RBAC, and namespace isolation

The module captures and exposes important output details, such as the namespace where ExternalDNS is deployed and the port-forward command for debugging, making it easier for operators to manage and monitor DNS automation. By taking in a pre-defined YAML specification as input, the module eliminates manual effort in setting up ExternalDNS across multiple clusters and cloud providers, while following OpenMCF's standardized approach for Kubernetes add-on deployment.
