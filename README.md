# OpenMCF

> **Deploy anywhere with one workflow.** Write declarative YAML once, deploy to AWS, GCP, Azure, or Kubernetes with the same CLI and consistent experience.

## What is OpenMCF?

An open-source multi-cloud deployment framework that brings Kubernetes-style consistency to infrastructure deployments everywhere. No vendor lock-in, no artificial abstractions—just provider-specific configurations with a unified structure and workflow.

**[Documentation](https://openmcf.org)** · **[Component Catalog](https://openmcf.org/docs/catalog)** · **[Website](https://openmcf.org)**

---

## Why OpenMCF?

- **One structure, any cloud** - Kubernetes Resource Model (apiVersion/kind/metadata/spec) for all deployments
- **Validate before deploy** - Protocol Buffer validations catch errors in seconds, not minutes
- **Zero abstraction** - Provider-specific configs preserve cloud capabilities, consistent experience across all
- **Choose your IaC** - Built-in Pulumi and Terraform/OpenTofu modules with feature parity
- **Build on top** - Auto-generated SDKs in Go, Python, TypeScript, Java from Protocol Buffer definitions

---

## Quick Start

### 1. Install the CLI

```bash
brew install plantonhq/tap/openmcf
```

### 2. Create a YAML Manifest

Example: Deploy Redis to Kubernetes using the [redis-kubernetes](https://buf.build/openmcf/apis/file/main:project/planton/provider/kubernetes/workload/rediskubernetes/v1/spec.proto) deployment component.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: RedisKubernetes
metadata:
  name: payments
  id: payments-namespace
spec:
  container:
    replicas: 1
    resources:
      limits:
        cpu: 50m
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 100Mi
    isPersistenceEnabled: true
    diskSize: 1Gi
```

You can create similar manifests for [AWS VPC](https://github.com/plantonhq/openmcf/tree/main/apis/project/planton/provider/aws/awsvpc/v1), [GKE Cluster](https://github.com/plantonhq/openmcf/tree/main/apis/project/planton/provider/gcp/gkecluster/v1), [Kafka on Kubernetes](https://github.com/plantonhq/openmcf/tree/main/apis/project/planton/provider/kubernetes/workload/kafkakubernetes/v1), and [many more](https://github.com/plantonhq/openmcf/tree/main/apis/project/planton/provider).

### 3. Deploy

```bash
openmcf pulumi up --manifest redis.yaml
```

---

## Learn More

- **[Getting Started Guide](https://openmcf.org/docs/getting-started)** - Your first deployment in 5 minutes
- **[Component Catalog](https://openmcf.org/docs/catalog)** - Browse 118+ deployment components across 10 providers
- **[Architecture](https://openmcf.org/docs/concepts/architecture)** - How Protocol Buffers, IaC modules, and CLI work together
- **[Planton Cloud](https://planton.cloud)** - Commercial SaaS platform with UI, CI/CD, and team features

---

## Contributing

Visit [CONTRIBUTING.md](CONTRIBUTING.md) for information on building OpenMCF from source or contributing improvements.

Also, refer to this [Contributor Guide](https://openmcf.org/docs/guide/contributor-guide) for detailed information about becoming a contributor to Project-Planton.

## License

OpenMCF is released under the [Apache 2.0 license](LICENSE). You are free to use, modify, and distribute this software in accordance with the license terms.

## Acknowledgments

- **Brian Grant & Kubernetes API team** for their foundational work on the Kubernetes Resource Model.
- The **[Protobuf Team](https://protobuf.dev/)** for laying the foundation for a powerful language neutral contract definition language.
- The **[Buf](https://github.com/bufbuild/buf) Team** for their Protobuf tooling—including BSR Docs, BSR SDKs, and ProtoValidate — which collectively democratized protobuf adoption and made this project possible.
- The **[Pulumi](https://github.com/pulumi/pulumi)** team for providing a powerful infrastructure as code platform that enables multi-language support.
- The **[spf13/cobra](https://github.com/spf13/cobra)** team for making building command line tools a bliss.
