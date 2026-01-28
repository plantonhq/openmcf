# Pulumi Examples for KubernetesNats

This document provides Pulumi usage examples for deploying NATS on Kubernetes using the official NATS Helm chart.

## Example 1: Basic NATS Cluster with Default Settings

Deploy a simple NATS cluster with default replicas, resources, and JetStream enabled.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsBasic = new kubernetesnats.NatsKubernetes("nats-basic", {
    metadata: {
        name: "nats-basic",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-basic",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 3,
            resources: {
                limits: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
                requests: {
                    cpu: "100m",
                    memory: "256Mi",
                },
            },
            diskSize: "10Gi",
        },
        disableJetStream: false,
        tlsEnabled: false,
        disableNatsBox: false,
    },
});

// Export the internal client URL
export const internalClientUrl = natsBasic.internalClientUrl;
export const namespace = natsBasic.namespace;
```

## Example 2: NATS Cluster with Bearer Token Authentication

Set up a secure NATS cluster using bearer token authentication.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsSecure = new kubernetesnats.NatsKubernetes("nats-secure", {
    metadata: {
        name: "nats-secure",
        org: "my-organization",
        env: "production",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-secure",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 5,
            resources: {
                limits: {
                    cpu: "2000m",
                    memory: "4Gi",
                },
                requests: {
                    cpu: "500m",
                    memory: "1Gi",
                },
            },
            diskSize: "20Gi",
        },
        auth: {
            enabled: true,
            scheme: "bearer_token",
        },
        tlsEnabled: true,
        disableJetStream: false,
    },
});

// Export connection details
export const namespace = natsSecure.namespace;
export const internalClientUrl = natsSecure.internalClientUrl;
export const authSecretName = "auth-nats";
export const authSecretKey = "token";
```

## Example 3: NATS Cluster with Basic Authentication

Deploy NATS with basic username/password authentication.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsBasicAuth = new kubernetesnats.NatsKubernetes("nats-basic-auth", {
    metadata: {
        name: "nats-basic-auth",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-basic-auth",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 3,
            resources: {
                limits: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
                requests: {
                    cpu: "100m",
                    memory: "256Mi",
                },
            },
            diskSize: "10Gi",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
        },
        tlsEnabled: false,
    },
});

export const namespace = natsBasicAuth.namespace;
export const username = "nats"; // Default admin username
export const authSecretName = "auth-nats";
export const passwordKey = "password";
```

## Example 4: NATS Cluster with External Access via Ingress

Deploy NATS configured with ingress to allow external clients to access messaging services.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsExternal = new kubernetesnats.NatsKubernetes("nats-external", {
    metadata: {
        name: "nats-external",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-external",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 3,
            resources: {
                limits: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
                requests: {
                    cpu: "100m",
                    memory: "256Mi",
                },
            },
            diskSize: "10Gi",
        },
        ingress: {
            enabled: true,
            hostname: "nats.example.com",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
        },
        tlsEnabled: true,
    },
});

// Export connection details
export const namespace = natsExternal.namespace;
export const externalHostname = natsExternal.externalHostname;
export const internalClientUrl = natsExternal.internalClientUrl;
export const externalClientUrl = `nats://nats.example.com:4222`;
```

## Example 5: Lightweight NATS Cluster Without JetStream

Set up a lightweight NATS messaging cluster without JetStream persistence.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsMinimal = new kubernetesnats.NatsKubernetes("nats-minimal", {
    metadata: {
        name: "nats-minimal",
        env: "development",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-minimal",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 1,
            resources: {
                limits: {
                    cpu: "500m",
                    memory: "512Mi",
                },
                requests: {
                    cpu: "100m",
                    memory: "128Mi",
                },
            },
            diskSize: "1Gi",
        },
        disableJetStream: true,
        tlsEnabled: false,
        disableNatsBox: true,
    },
});

export const namespace = natsMinimal.namespace;
export const internalClientUrl = natsMinimal.internalClientUrl;
```

## Example 6: High Availability NATS Cluster with Basic Auth and No-Auth User

Deploy a highly available NATS cluster with both authenticated and unauthenticated access.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsHA = new kubernetesnats.NatsKubernetes("nats-ha-metrics", {
    metadata: {
        name: "nats-ha-metrics",
        org: "acme-corp",
        env: "production",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-ha-metrics",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 7,
            resources: {
                limits: {
                    cpu: "4000m",
                    memory: "8Gi",
                },
                requests: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
            },
            diskSize: "50Gi",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
            noAuthUser: {
                enabled: true,
                publishSubjects: [
                    "telemetry.>",
                    "metrics.>",
                ],
            },
        },
        tlsEnabled: true,
        disableJetStream: false,
        ingress: {
            enabled: true,
            hostname: "nats-ha.example.com",
        },
    },
});

// Export comprehensive connection details
export const namespace = natsHA.namespace;
export const externalHostname = natsHA.externalHostname;
export const internalClientUrl = natsHA.internalClientUrl;
export const metricsEndpoint = `http://nats-prom.${natsHA.namespace}.svc.cluster.local:7777/metrics`;
```

## Connecting to NATS from Pulumi Programs

### From within the Kubernetes cluster:

```typescript
import * as pulumi from "@pulumi/pulumi";

const nats = // ... your NATS deployment

// Use the internal client URL
const clientUrl = nats.internalClientUrl; // nats://<service>.<namespace>.svc.cluster.local:4222
```

### Retrieving Authentication Credentials:

For bearer token authentication:

```typescript
import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";

const nats = // ... your NATS deployment with bearer_token auth

const authSecret = k8s.core.v1.Secret.get("auth-secret", 
    pulumi.interpolate`${nats.namespace}/auth-nats`);

export const bearerToken = authSecret.data.apply(d => 
    Buffer.from(d["token"], "base64").toString());
```

For basic authentication:

```typescript
import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";

const nats = // ... your NATS deployment with basic_auth

const authSecret = k8s.core.v1.Secret.get("auth-secret", 
    pulumi.interpolate`${nats.namespace}/auth-nats`);

export const username = authSecret.data.apply(d => 
    Buffer.from(d["user"], "base64").toString());
export const password = authSecret.data.apply(d => 
    Buffer.from(d["password"], "base64").toString());
```

## Example 7: Using Pre-existing Namespace

Deploy NATS into an existing namespace that's managed separately.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as kubernetesnats from "@openmcf/kubernetesnats";

// Assume namespace is created elsewhere (e.g., GitOps, separate stack)
const existingNamespace = "shared-messaging";

const natsWithExistingNs = new kubernetesnats.NatsKubernetes("nats-existing-ns", {
    metadata: {
        name: "nats-existing-ns",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: existingNamespace,
        },
        createNamespace: false, // Don't create the namespace
        serverContainer: {
            replicas: 3,
            resources: {
                limits: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
                requests: {
                    cpu: "100m",
                    memory: "256Mi",
                },
            },
            diskSize: "10Gi",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
        },
        tlsEnabled: true,
    },
});

export const namespace = existingNamespace;
export const internalClientUrl = natsWithExistingNs.internalClientUrl;
```

**Use Case**: Ideal for GitOps workflows or when namespace has specific ResourceQuotas, NetworkPolicies, or other configurations managed separately.

## Example 8: NATS with NACK Controller and JetStream Streams

Deploy NATS with the NACK (NATS Controllers for Kubernetes) controller for declarative JetStream stream management.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsWithStreams = new kubernetesnats.NatsKubernetes("nats-with-streams", {
    metadata: {
        name: "nats-with-streams",
        org: "my-org",
        env: "development",
    },
    spec: {
        targetCluster: {
            clusterName: "my-gke-cluster",
        },
        namespace: {
            value: "nats-streams",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 3,
            resources: {
                limits: {
                    cpu: "1000m",
                    memory: "2Gi",
                },
                requests: {
                    cpu: "100m",
                    memory: "256Mi",
                },
            },
            diskSize: "10Gi",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
        },
        // Enable NACK controller for declarative stream management
        nackController: {
            enabled: true,
            enableControlLoop: true,  // Required for KeyValue/ObjectStore support
        },
        // Define JetStream streams
        streams: [
            {
                name: "orders",
                subjects: ["orders.*", "orders.>"],
                storage: "file",
                replicas: 3,
                retention: "limits",
                maxAge: "7d",
                maxBytes: 1073741824,  // 1GB
                consumers: [
                    {
                        durableName: "orders-processor",
                        deliverPolicy: "all",
                        ackPolicy: "explicit",
                        maxAckPending: 1000,
                        ackWait: "30s",
                    },
                ],
            },
            {
                name: "events",
                subjects: ["events.>"],
                storage: "memory",
                replicas: 1,
                retention: "interest",
            },
        ],
    },
});

// Export stream information
export const namespace = natsWithStreams.namespace;
export const internalClientUrl = natsWithStreams.internalClientUrl;
export const nackEnabled = true;
export const streamsCreated = ["orders", "events"];
```

## Example 9: Production NATS with Multiple Streams and Consumers

A comprehensive production configuration with multiple streams, consumers, and external access.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetesnats from "@openmcf/kubernetesnats";

const natsProd = new kubernetesnats.NatsKubernetes("nats-prod", {
    metadata: {
        name: "nats-prod",
        org: "acme-corp",
        env: "production",
        labels: {
            team: "platform",
            costCenter: "engineering",
        },
    },
    spec: {
        targetCluster: {
            clusterName: "prod-gke-cluster",
        },
        namespace: {
            value: "nats-prod",
        },
        createNamespace: true,
        serverContainer: {
            replicas: 5,
            resources: {
                limits: {
                    cpu: "2000m",
                    memory: "4Gi",
                },
                requests: {
                    cpu: "500m",
                    memory: "1Gi",
                },
            },
            diskSize: "50Gi",
        },
        auth: {
            enabled: true,
            scheme: "basic_auth",
        },
        tlsEnabled: true,
        ingress: {
            enabled: true,
            hostname: "nats.prod.example.com",
        },
        // NACK controller with control-loop for reliable state enforcement
        nackController: {
            enabled: true,
            enableControlLoop: true,
            helmChartVersion: "0.31.1",
            appVersion: "0.21.1",
        },
        // Production streams configuration
        streams: [
            // High-throughput API events stream
            {
                name: "api-events",
                subjects: ["api.events.>"],
                storage: "file",
                replicas: 3,
                retention: "limits",
                maxAge: "24h",
                maxBytes: 10737418240,  // 10GB
                maxMsgs: 10000000,
                discard: "old",
                description: "API event stream for real-time processing",
                consumers: [
                    {
                        durableName: "api-processor",
                        deliverPolicy: "all",
                        ackPolicy: "explicit",
                        maxAckPending: 5000,
                        maxDeliver: 5,
                        ackWait: "60s",
                        replayPolicy: "instant",
                        description: "Main API event processor",
                    },
                    {
                        durableName: "api-analytics",
                        deliverPolicy: "all",
                        ackPolicy: "explicit",
                        filterSubject: "api.events.orders.*",
                        description: "Analytics consumer for order events",
                    },
                ],
            },
            // Work queue for background jobs
            {
                name: "background-jobs",
                subjects: ["jobs.>"],
                storage: "file",
                replicas: 3,
                retention: "workqueue",
                maxAge: "1h",
                description: "Work queue for background job processing",
                consumers: [
                    {
                        durableName: "job-worker",
                        ackPolicy: "explicit",
                        maxAckPending: 100,
                        maxDeliver: 3,
                        ackWait: "5m",
                    },
                ],
            },
            // Ephemeral notifications stream
            {
                name: "notifications",
                subjects: ["notify.>"],
                storage: "memory",
                replicas: 1,
                retention: "interest",
                maxAge: "5m",
                description: "Real-time notifications (ephemeral)",
            },
        ],
        natsHelmChartVersion: "2.12.3",
    },
});

// Comprehensive exports
export const namespace = natsProd.namespace;
export const internalClientUrl = natsProd.internalClientUrl;
export const externalHostname = natsProd.externalHostname;
```

## Verifying NACK Stream Creation

After deployment, you can verify the streams were created:

```bash
# Check NACK controller is running
kubectl get pods -n <namespace> -l app.kubernetes.io/name=nack

# List Stream custom resources
kubectl get streams -n <namespace>

# List Consumer custom resources
kubectl get consumers -n <namespace>

# Describe a specific stream
kubectl describe stream <stream-name> -n <namespace>

# Check stream status using nats-box
kubectl exec -it <nats-box-pod> -n <namespace> -- nats stream list
kubectl exec -it <nats-box-pod> -n <namespace> -- nats stream info <stream-name>
```

## Notes

- The deployment uses the official NATS Helm chart for production-grade message streaming
- JetStream provides persistent streaming capabilities with file-based storage
- For production deployments, use an odd number of replicas (3, 5, 7) for optimal quorum
- TLS encryption is recommended for production workloads
- The `nats-box` utility pod is deployed by default for debugging and testing
- Metrics are exposed at port 7777 for Prometheus scraping
- **NACK Controller**: Enables declarative stream/consumer management via Kubernetes CRDs
- **Control-Loop Mode**: Required for KeyValue, ObjectStore, and Account support
- **Stream Replication**: Use odd numbers (1, 3, 5) for stream replicas to maintain quorum

