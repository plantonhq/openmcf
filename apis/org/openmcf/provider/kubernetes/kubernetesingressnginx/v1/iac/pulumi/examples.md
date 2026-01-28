# Kubernetes Ingress NGINX - Pulumi Module Examples

This document provides Pulumi-specific examples for deploying the NGINX Ingress Controller.

## Prerequisites

- Pulumi CLI installed
- Go 1.21 or later
- Access to Kubernetes cluster
- kubectl configured

## Setup

Create a new Pulumi Go project:

```bash
mkdir my-ingress-deployment
cd my-ingress-deployment
pulumi new go
```

Add dependencies to `go.mod`:

```go
require (
    github.com/plantonhq/openmcf v1.0.0
    github.com/pulumi/pulumi-kubernetes/sdk/v4 v4.0.0
    github.com/pulumi/pulumi/sdk/v3 v3.0.0
)
```

## Example 1: Basic External Ingress Controller

Deploy NGINX Ingress Controller with external load balancer on any Kubernetes cluster.

```go
package main

import (
    kubernetesingressnginxv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        ingressSpec := &kubernetesingressnginxv1.KubernetesIngressNginxSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "ingress-nginx",
                },
            },
            CreateNamespace: true,
            ChartVersion:    "4.11.1",
            Internal:        false,
        }

        ingress := &kubernetesingressnginxv1.KubernetesIngressNginx{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesIngressNginx",
            Metadata: &shared.CloudResourceMetadata{
                Name: "basic-ingress",
            },
            Spec: ingressSpec,
        }

        stackInput := &kubernetesingressnginxv1.KubernetesIngressNginxStackInput{
            Target: ingress,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 1b: Using Existing Namespace

Deploy to a namespace created by your platform team.

```go
package main

import (
    kubernetesingressnginxv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        ingressSpec := &kubernetesingressnginxv1.KubernetesIngressNginxSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "platform-ingress",
                },
            },
            CreateNamespace: false,  // Use existing namespace
            ChartVersion:    "4.11.1",
            Internal:        false,
        }

        ingress := &kubernetesingressnginxv1.KubernetesIngressNginx{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesIngressNginx",
            Metadata: &shared.CloudResourceMetadata{
                Name: "existing-ns-ingress",
            },
            Spec: ingressSpec,
        }

        stackInput := &kubernetesingressnginxv1.KubernetesIngressNginxStackInput{
            Target: ingress,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 2: GKE with Static IP

Deploy on Google Kubernetes Engine with a reserved static IP address.

```go
package main

import (
    kubernetesingressnginxv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        ingressSpec := &kubernetesingressnginxv1.KubernetesIngressNginxSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "gke-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "ingress-nginx",
                },
            },
            CreateNamespace: true,
            ChartVersion:    "4.11.1",
            Internal:        false,
            ProviderConfig: &kubernetesingressnginxv1.KubernetesIngressNginxSpec_Gke{
                Gke: &kubernetesingressnginxv1.KubernetesIngressNginxGkeConfig{
                    StaticIpName: "prod-ingress-static-ip",
                },
            },
        }

        ingress := &kubernetesingressnginxv1.KubernetesIngressNginx{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesIngressNginx",
            Metadata: &shared.CloudResourceMetadata{
                Name: "gke-ingress",
                Org:  "my-company",
                Env:  "production",
            },
            Spec: ingressSpec,
        }

        stackInput := &kubernetesingressnginxv1.KubernetesIngressNginxStackInput{
            Target: ingress,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

**Prerequisites:**

```bash
# Reserve static IP in GCP
gcloud compute addresses create prod-ingress-static-ip --global
```

## Example 3: EKS with NLB and Security Groups

Deploy on Amazon EKS with Network Load Balancer and custom security groups.

```go
package main

import (
    kubernetesingressnginxv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        ingressSpec := &kubernetesingressnginxv1.KubernetesIngressNginxSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "eks-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "ingress-nginx",
                },
            },
            CreateNamespace: true,
            ChartVersion:    "4.11.1",
            Internal:        false,
            ProviderConfig: &kubernetesingressnginxv1.KubernetesIngressNginxSpec_Eks{
                Eks: &kubernetesingressnginxv1.KubernetesIngressNginxEksConfig{
                    AdditionalSecurityGroupIds: []*foreignkeyv1.StringValueOrRef{
                        {Value: &foreignkeyv1.StringValueOrRef_Value{Value: "sg-0123456789abcdef0"}},
                        {Value: &foreignkeyv1.StringValueOrRef_Value{Value: "sg-fedcba9876543210"}},
                    },
                    SubnetIds: []*foreignkeyv1.StringValueOrRef{
                        {Value: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-abc123"}},
                        {Value: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-def456"}},
                    },
                    IrsaRoleArnOverride: "arn:aws:iam::123456789012:role/ingress-nginx-role",
                },
            },
        }

        ingress := &kubernetesingressnginxv1.KubernetesIngressNginx{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesIngressNginx",
            Metadata: &shared.CloudResourceMetadata{
                Name: "eks-ingress",
                Org:  "my-company",
                Env:  "production",
            },
            Spec: ingressSpec,
        }

        stackInput := &kubernetesingressnginxv1.KubernetesIngressNginxStackInput{
            Target: ingress,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 4: AKS with Workload Identity

Deploy on Azure Kubernetes Service with Workload Identity integration.

```go
package main

import (
    kubernetesingressnginxv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesingressnginx/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        ingressSpec := &kubernetesingressnginxv1.KubernetesIngressNginxSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "aks-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "ingress-nginx",
                },
            },
            CreateNamespace: true,
            ChartVersion:    "4.11.1",
            Internal:        false,
            ProviderConfig: &kubernetesingressnginxv1.KubernetesIngressNginxSpec_Aks{
                Aks: &kubernetesingressnginxv1.KubernetesIngressNginxAksConfig{
                    ManagedIdentityClientId: "12345678-1234-1234-1234-123456789012",
                    PublicIpName:            "prod-ingress-public-ip",
                },
            },
        }

        ingress := &kubernetesingressnginxv1.KubernetesIngressNginx{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesIngressNginx",
            Metadata: &shared.CloudResourceMetadata{
                Name: "aks-ingress",
                Org:  "my-company",
                Env:  "production",
            },
            Spec: ingressSpec,
        }

        stackInput := &kubernetesingressnginxv1.KubernetesIngressNginxStackInput{
            Target: ingress,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Stack Outputs

All examples export the following outputs:

```bash
# View all outputs
pulumi stack output

# View specific output
pulumi stack output namespace
pulumi stack output service-name
pulumi stack output release-name
pulumi stack output service-type
```

## Using with Pulumi Stacks

### Development Stack

```bash
pulumi stack init dev
pulumi config set kubernetes-cluster dev-cluster-credential
pulumi up
```

### Production Stack

```bash
pulumi stack init prod
pulumi config set kubernetes-cluster prod-cluster-credential
pulumi config set chart-version 4.11.1
pulumi up
```

## Best Practices

1. **Use Pulumi Secrets**: Store credentials securely
   ```bash
   pulumi config set --secret cluster-credential <value>
   ```

2. **Separate Stacks**: Use different stacks for environments
   - `dev` - Development deployments
   - `staging` - Pre-production testing
   - `prod` - Production deployments

3. **Version Pinning**: Specify exact chart versions in production

4. **Resource Tagging**: Use metadata labels for cost tracking
   ```go
   Metadata: &shared.CloudResourceMetadata{
       Name: "prod-ingress",
       Org:  "platform-team",
       Env:  "production",
       Labels: map[string]string{
           "cost-center": "infrastructure",
           "team": "platform",
       },
   }
   ```

## Troubleshooting

### Check Deployment Status

```bash
pulumi stack
pulumi stack output
kubectl get all -n kubernetes-ingress-nginx
```

### View Logs

```bash
kubectl logs -n kubernetes-ingress-nginx -l app.kubernetes.io/component=controller
```

### Debug Helm Release

```bash
helm list -n kubernetes-ingress-nginx
helm get values kubernetes-ingress-nginx -n kubernetes-ingress-nginx
```

## Additional Resources

- [Pulumi Documentation](https://www.pulumi.com/docs/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/)

