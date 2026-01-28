# Overview

The provided Pulumi module automates the deployment of Elasticsearch and Kibana on Kubernetes clusters based on a
unified API resource specification. It interprets the `ElasticsearchKubernetes` resource defined in the proto files to
create and configure necessary Kubernetes resources such as namespaces, StatefulSets, Services, and Ingress
configurations. The module utilizes custom resource definitions (CRDs) for Elasticsearch and Kibana to ensure that
deployments are consistent, scalable, and adhere to best practices.

Key features of this module include customizable replicas, resource allocation (CPU and memory), and persistence options
for Elasticsearch and Kibana containers. It supports optional persistence by attaching persistent volumes when enabled.
The module also configures ingress resources using Istio and the Gateway API, enabling both internal and external access
with secure TLS termination. By handling complex configurations and integrations, this module allows developers to
deploy robust Elasticsearch and Kibana instances using a simple YAML specification, streamlining the infrastructure
deployment process.
