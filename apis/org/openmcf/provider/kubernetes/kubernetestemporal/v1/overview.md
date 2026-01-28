The **TemporalKubernetes** component in OpenMCF simplifies the deployment and management of Temporal workflows on
Kubernetes clusters. By encapsulating essential Temporal services and supporting flexible database integrations, it
seamlessly fits into the OpenMCF multi-cloud ecosystem, providing a uniform deployment experience across
different cloud environments.

## Purpose and Functionality

- **Unified Temporal Deployment**: Easily deploy Temporal’s core services—frontend, history, matching, and worker—within
  Kubernetes using minimal configuration.
- **Flexible Database Options**: Supports multiple database backends (Cassandra, PostgreSQL, MySQL) either by
  provisioning in-cluster databases automatically or connecting to existing external databases.
- **Integrated Observability**: Optionally configure external Elasticsearch clusters for enhanced visibility and
  observability of Temporal workflows.
- **Effortless Exposure**: Seamlessly expose Temporal frontend and web UI services through Kubernetes ingress
  controllers or load balancers with simple, declarative configurations.

## Key Benefits

- **Simplified Deployment**: Condenses complex Temporal setup into straightforward YAML manifests validated by
  OpenMCF’s standardized Protobuf schemas.
- **Adaptable Configuration**: Offers sensible defaults with easy overrides, minimizing boilerplate and accelerating
  deployment timelines.
- **Consistent Multi-Cloud Experience**: Leverages OpenMCF’s uniform APIs and CLI workflows, enabling reliable
  and repeatable deployments across AWS, GCP, Azure, and other Kubernetes-supported environments.
- **Optimized Observability**: Integrate external Elasticsearch clusters effortlessly to enhance workflow visibility
  without additional complexity.

Below is a minimal YAML example demonstrating a Temporal deployment with an external PostgreSQL database (note the use
of **camel-case** keys):

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: TemporalKubernetes
metadata:
  name: exampleTemporal
spec:
  database:
    backend: postgresql
    externalDatabase:
      host: postgres.example.com
      port: 5432
      user: temporalUser
      password: temporalPass
  ingress:
    frontend:
      enabled: true
      hostname: temporal-frontend.example.com
    webUi:
      enabled: true
      hostname: temporal-ui.example.com
  externalElasticsearch:
    host: elasticsearch.example.com
    port: 9200
```

Leverage the **TemporalKubernetes** component to effortlessly integrate powerful Temporal workflow orchestration into
your multi-cloud deployments. OpenMCF’s consistent CLI commands and schema-driven validations ensure easy
operations, enhanced reliability, and accelerated productivity.
