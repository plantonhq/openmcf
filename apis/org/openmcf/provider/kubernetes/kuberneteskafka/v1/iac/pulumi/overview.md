# Overview

The Pulumi module provided automates the deployment of a Kafka cluster on Kubernetes using the Strimzi operator. Written
in Golang, it leverages Pulumi to define and manage the necessary Kubernetes resources based on user-defined
specifications. The module accepts an API resource input, which includes details such as broker and zookeeper
configurations, Kafka topics, and optional components like Schema Registry and Kafka UI (Kowl).

By encapsulating complex configurations into a reusable module, developers can effortlessly set up Kafka clusters with
custom specifications. The module handles the creation of namespaces, deployment of Kafka and Zookeeper pods, setup of
Kafka topics, and configuration of ingress resources for external access. It also manages TLS certificates and hostnames
for secure communication. This standardized approach simplifies the deployment process, ensuring consistency and
scalability across different environments.
