# Queue Driven

This preset sizes a worker fleet by its backlog instead of its CPU: roughly one pod per 30 ready messages in one queue. Queue depth is the honest signal for workers — a worker grinding through a backlog can look CPU-calm while the queue grows without bound. The `external` metric family reads signals with no in-cluster object at all: broker queue depths, cloud load balancer QPS, anything an external-metrics adapter serves.

## When to Use

- Message-queue consumers (RabbitMQ, SQS, Pub/Sub, Kafka consumer groups) whose correct size is a function of backlog
- Batch-ish workloads where "how much work is waiting" beats "how hot are the pods" as a scaling signal
- Any workload scaled by a metric that lives outside the cluster

## Key Configuration Choices

- **Requires an external-metrics adapter** — the HPA reads metrics, it does not produce them. Something must serve the `external.metrics.k8s.io` API: KEDA (which can also manage HPAs itself — never both on one target) or prometheus-adapter are the usual choices. Without an adapter, the metric reads unavailable and the autoscaler holds the current count
- **`metric.name` + `match_labels`** — the identity as the ADAPTER exposes it, with the selector narrowing which series is read (one queue's depth, not the sum of all queues). Both are adapter-side names, not Kubernetes object names
- **`average_value: "30"`** — the per-pod form: the controller targets `depth ÷ pods = 30`, i.e. one pod per 30 ready messages. This is the usual form for external metrics; `raw_value` would instead compare the whole queue depth against a fixed number regardless of fleet size
- **`min_replicas: 1`** — the floor keeps one consumer alive through quiet periods. Scale-to-zero is feature-gated upstream and not modeled; for event-driven scale-to-zero semantics, use KEDA end-to-end instead of a hand-written HPA
- **`max_replicas: 30`** — with backlog-proportional scaling, the ceiling is what stops a poison-message flood from buying the whole cluster

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The scale target's own namespace — an HPA cannot scale across namespaces | Your namespace management |
| `<your-workload-name>` | The worker Deployment's name | The workload's manifest |
| `<your-queue-name>` | The queue label value as the external-metrics adapter exposes it | Your metrics adapter configuration |

The metric name `queue_messages_ready` and the 30-messages-per-pod target are working examples — replace them with the adapter's actual metric name and your throughput math.

## Related Presets

- **01-cpu-autoscale** — the CPU workhorse for request-serving workloads
- **04-behavior-tuned** — pair with conservative scale-down so the fleet does not cliff-drop the moment the queue empties
