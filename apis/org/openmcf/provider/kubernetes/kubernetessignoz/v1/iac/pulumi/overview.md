# Overview

The **SignozKubernetes** API resource streamlines and standardizes how developers deploy SigNoz observability platform onto Kubernetes. This Pulumi module interprets a `SignozKubernetesStackInput`, which includes core details like Kubernetes credentials and the SigNoz configuration. From this input, the module automatically creates and manages:

- **Kubernetes Namespaces**
- **Helm Chart Deployments** for SigNoz using the official SigNoz chart
- **SigNoz Binary** (UI, API Server, Ruler, Alertmanager)
- **OpenTelemetry Collector** for data ingestion
- **ClickHouse Database** (self-managed or external connection)
- **ZooKeeper** (optional) for cluster coordination in distributed ClickHouse setups
- **Services** for exposing SigNoz and OpenTelemetry Collector
- **LoadBalancer Services** (optional) for external access

Developers simply provide a declarative resource specification – focusing on resource allocations, database mode, component scaling, and ingress preferences – while the module handles the underlying Kubernetes constructs.

### Key Features

1. **Unified Observability Platform**  
   Deploys a complete observability stack with logs, metrics, and traces in a single application. Built on OpenTelemetry standards for vendor-neutral instrumentation.

2. **Flexible Database Options**  
   Choose between self-managed ClickHouse within the cluster or connect to an external ClickHouse instance. Self-managed mode supports both simple single-node and distributed cluster configurations.

3. **High Availability Support**  
   Deploy distributed ClickHouse clusters with configurable sharding and replication. Zookeeper coordination ensures consistency and failover capabilities.

4. **Scalable Data Ingestion**  
   Deploy multiple OpenTelemetry Collector replicas to handle high-volume telemetry data. Scale horizontally based on ingestion load.

5. **Component Independence**  
   Scale SigNoz UI/API, OpenTelemetry Collector, and ClickHouse independently based on workload requirements. Each component can be tuned for specific performance needs.

6. **Persistence Management**  
   Toggle data persistence with customizable disk sizes for ClickHouse. When enabled, telemetry data survives pod restarts and node failures.

7. **Security First**  
   When using self-managed ClickHouse, the module automatically generates secure random passwords and stores them in Kubernetes Secrets. Credentials never appear in version control or container images.

8. **Resource Optimization**  
   Specify CPU and memory requests/limits for all components. This ensures your observability platform runs efficiently without jeopardizing the cluster's stability.

9. **Optional External Access**  
   Enable ingress to expose SigNoz UI and OpenTelemetry Collector via LoadBalancer services with external DNS annotations. Perfect for accessing dashboards and sending telemetry from external applications.

10. **Helm Chart Integration**  
    Leverages the official SigNoz Helm chart with extensive customization options through the `helmValues` field.

### Architecture

**Standalone Mode (Single-Node ClickHouse):**
```
┌─────────────────────────────────────────────────┐
│  Namespace                                      │
│  ┌───────────────────────────────────────────┐  │
│  │  SigNoz Binary                            │  │
│  │  - UI (ReactJS)                           │  │
│  │  - API Server                             │  │
│  │  - Ruler (Alerting)                       │  │
│  │  - Alertmanager                           │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  OpenTelemetry Collector                  │  │
│  │  - Data Ingestion (gRPC, HTTP)            │  │
│  │  - Processing & Batching                  │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  ClickHouse                               │  │
│  │  - Persistent Volume (Telemetry Data)     │  │
│  │  - Auto-generated Password                │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  Services (ClusterIP)                     │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  LoadBalancers (Optional)                 │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

**Distributed Mode (HA with Clustering):**
```
┌────────────────────────────────────────────────────────┐
│  Namespace                                             │
│  ┌──────────────────────────────────────────────────┐  │
│  │  SigNoz Binary (Multiple Replicas)              │  │
│  │  - UI (ReactJS)                                  │  │
│  │  - API Server (Load Balanced)                    │  │
│  │  - Ruler (Alerting)                              │  │
│  │  - Alertmanager                                  │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  OpenTelemetry Collector (Scaled)                │  │
│  │  - Multiple replicas for high throughput         │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  ZooKeeper (3 or 5 nodes for quorum)            │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  ClickHouse Shard 1                              │  │
│  │  ├─ Replica 1 (Persistent Volume)               │  │
│  │  └─ Replica 2 (Persistent Volume)               │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  ClickHouse Shard 2                              │  │
│  │  ├─ Replica 1 (Persistent Volume)               │  │
│  │  └─ Replica 2 (Persistent Volume)               │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Services & LoadBalancers                        │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

**External ClickHouse Mode:**
```
┌─────────────────────────────────────────────────┐
│  SigNoz Namespace                               │
│  ┌───────────────────────────────────────────┐  │
│  │  SigNoz Binary                            │  │
│  │  - UI, API, Ruler, Alertmanager           │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  OpenTelemetry Collector                  │  │
│  └───────────────────────────────────────────┘  │
│         │                                        │
│         └─────────────┐                          │
└───────────────────────┼──────────────────────────┘
                        │
                        ▼
         ┌──────────────────────────────┐
         │  External ClickHouse Cluster │
         │  (Managed Separately)        │
         └──────────────────────────────┘
```

### Use Cases

- **Microservices Observability**: Monitor distributed systems with unified logs, metrics, and traces
- **Application Performance Monitoring (APM)**: Track service latencies, error rates, and throughput
- **Infrastructure Monitoring**: Collect and analyze metrics from Kubernetes clusters and applications
- **Log Analytics**: Centralized log aggregation and analysis with powerful query capabilities
- **Distributed Tracing**: Visualize request flows across microservices
- **Custom Dashboards**: Create tailored dashboards for specific services and teams
- **Alerting and Incident Response**: Configure alerts based on logs, metrics, or trace data
- **Cost Optimization**: Identify performance bottlenecks and optimize resource utilization

### Data Flow

1. **Instrumentation**: Applications instrumented with OpenTelemetry SDK export telemetry data
2. **Ingestion**: OpenTelemetry Collector receives data via gRPC or HTTP protocols
3. **Processing**: Collector batches, enriches (adds Kubernetes metadata), and processes data
4. **Storage**: Processed telemetry is written to ClickHouse database
5. **Visualization**: Users access SigNoz UI to query and visualize telemetry data
6. **Analysis**: SigNoz API queries ClickHouse and returns results to the frontend

Overall, **SignozKubernetes** helps you focus on observability and application performance. By delegating the Kubernetes resource orchestration to this module, you gain a cleaner, more consistent deployment experience across development, staging, and production environments.

