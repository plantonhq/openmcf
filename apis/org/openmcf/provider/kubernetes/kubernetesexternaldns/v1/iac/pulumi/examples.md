# Pulumi Examples for KubernetesExternalDns

Examples of deploying ExternalDNS using the Pulumi module directly.

---

## Example 1: GKE with Cloud DNS (Programmatic)

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	kubernetesexternaldnsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	module "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1/iac/pulumi/module"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-gke-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-prod",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "prod-gke-cluster",
						},
					},
					CreateNamespace: true,  // Module will create the namespace
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Gke{
						Gke: &kubernetesexternaldnsv1.KubernetesExternalDnsGkeConfig{
							ProjectId: &v1.StringValueOrRef{
								Value: "my-gcp-project",
							},
							DnsZoneId: &v1.StringValueOrRef{
								Value: "my-dns-zone-id",
							},
						},
					},
				},
			},
		}

		return module.Resources(ctx, stackInput)
	})
}
```

---

## Example 2: EKS with Route53

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	kubernetesexternaldnsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	module "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1/iac/pulumi/module"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-eks-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-eks",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "prod-eks-cluster",
						},
					},
					CreateNamespace: true,  // Module will create the namespace
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Eks{
						Eks: &kubernetesexternaldnsv1.KubernetesExternalDnsEksConfig{
							Route53ZoneId: &v1.StringValueOrRef{
								Value: "Z1234567890ABC",
							},
							IrsaRoleArnOverride: "arn:aws:iam::123456789012:role/external-dns-role",
						},
					},
				},
			},
		}

		return module.Resources(ctx, stackInput)
	})
}
```

---

## Example 3: Cloudflare DNS

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	kubernetesexternaldnsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	module "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1/iac/pulumi/module"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-k8s-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-cloudflare",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "my-cluster",
						},
					},
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Cloudflare{
						Cloudflare: &kubernetesexternaldnsv1.KubernetesExternalDnsCloudflareConfig{
							ApiToken: "your-cloudflare-api-token",
							DnsZoneId: &v1.StringValueOrRef{
								Value: "1234567890abcdef1234567890abcdef",
							},
							IsProxied: true,
						},
					},
				},
			},
		}

		return module.Resources(ctx, stackInput)
	})
}
```

---

## Example 4: Using Pulumi Config

Store sensitive values in Pulumi config instead of hardcoding:

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	kubernetesexternaldnsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	module "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1/iac/pulumi/module"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		
		// Read config values
		cloudflareToken := cfg.RequireSecret("cloudflare-token")
		dnsZoneId := cfg.Require("dns-zone-id")
		
		stackInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-k8s-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-cloudflare",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "my-cluster",
						},
					},
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Cloudflare{
						Cloudflare: &kubernetesexternaldnsv1.KubernetesExternalDnsCloudflareConfig{
							ApiToken: string(cloudflareToken),
							DnsZoneId: &v1.StringValueOrRef{
								Value: dnsZoneId,
							},
						},
					},
				},
			},
		}

		return module.Resources(ctx, stackInput)
	})
}
```

Set config values:

```bash
pulumi config set cloudflare-token --secret your-token-here
pulumi config set dns-zone-id 1234567890abcdef1234567890abcdef
```

---

## Example 5: Multiple Instances

Deploy multiple ExternalDNS instances for different zones:

```go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	kubernetesexternaldnsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	module "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesexternaldns/v1/iac/pulumi/module"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Production domain
		prodInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-gke-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-prod",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "shared-cluster",
						},
					},
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Gke{
						Gke: &kubernetesexternaldnsv1.KubernetesExternalDnsGkeConfig{
							ProjectId: &v1.StringValueOrRef{
								Value: "my-project",
							},
							DnsZoneId: &v1.StringValueOrRef{
								Value: "prod-zone-id",
							},
						},
					},
				},
			},
		}

		if err := module.Resources(ctx, prodInput); err != nil {
			return err
		}

		// Staging domain
		stagingInput := &kubernetesexternaldnsv1.KubernetesExternalDnsStackInput{
			ProviderConfig: &kubernetesexternaldnsv1.KubernetesProviderConfig{
				KubernetesProviderConfigId: "my-gke-credentials",
			},
			Target: &kubernetesexternaldnsv1.KubernetesExternalDns{
				Metadata: &shared.CloudResourceMetadata{
					Name: "external-dns-staging",
				},
				Spec: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec{
					TargetCluster: &kubernetes.KubernetesAddonTargetCluster{
						KubernetesClusterId: &v1.StringValueOrRef{
							Value: "shared-cluster",
						},
					},
					ProviderConfig: &kubernetesexternaldnsv1.KubernetesExternalDnsSpec_Gke{
						Gke: &kubernetesexternaldnsv1.KubernetesExternalDnsGkeConfig{
							ProjectId: &v1.StringValueOrRef{
								Value: "my-project",
							},
							DnsZoneId: &v1.StringValueOrRef{
								Value: "staging-zone-id",
							},
						},
					},
				},
			},
		}

		return module.Resources(ctx, stagingInput)
	})
}
```

---

## Running Examples

### Option 1: Direct Execution

```bash
# Initialize
cd iac/pulumi
pulumi stack init dev

# Set credentials
pulumi config set kubernetes-credentials-id <your-cred-id>

# Deploy
pulumi up
```

### Option 2: Via OpenMCF CLI

Create a YAML manifest (see `examples.md` in parent directory) and deploy:

```bash
openmcf deploy external-dns.yaml
```

---

## Best Practices

### Security

1. **Never hardcode secrets** - Use Pulumi config secrets or external secret managers
2. **Use cloud-native auth** - Prefer IRSA/Workload Identity over static credentials
3. **Scope to zones** - Always configure zone filtering to prevent managing wrong DNS

### Organization

1. **One instance per zone** - Deploy separate ExternalDNS instances for different DNS zones
2. **Unique release names** - Use descriptive names that indicate the zone/environment
3. **Consistent naming** - Follow a naming convention like `external-dns-<env>-<domain>`

### Operations

1. **Pin versions** - Specify exact Helm chart and ExternalDNS versions
2. **Test in staging** - Always test DNS automation in non-production first
3. **Monitor logs** - Set up log aggregation for ExternalDNS pods
4. **Use GitOps** - Store Pulumi code in Git and use CI/CD for deployment

---

## Troubleshooting

**Import errors?**

Make sure dependencies are installed:

```bash
go mod download
go mod tidy
```

**Authentication errors?**

Verify Kubernetes credentials:

```bash
pulumi config set kubernetes-credentials-id <correct-id>
```

**Deployment fails?**

Check Pulumi logs:

```bash
pulumi up --debug
```

For more examples, see the [main examples file](../examples.md).

