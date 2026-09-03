# Planton Operator

Kubernetes operator that deploys and manages a complete Planton platform on any Kubernetes cluster through a single `PlantonPlatform` custom resource.

## What It Does

The operator encodes all Planton deployment knowledge -- component sequencing, health gates, credential wiring, and post-deployment bootstrapping -- so users only need to `kubectl apply` one YAML file. The operator continuously reconciles the desired state, self-healing when components drift.

## Quick Start

```bash
# Create a Kind cluster for local development
make kind-create

# Install CRDs into the cluster
make install

# Run the operator locally (connects to Kind's API server)
make run

# In another terminal, apply a sample resource
kubectl apply -f config/samples/minimal.yaml

# Watch the status converge
kubectl get plantonplatform -w
```

## Architecture

The operator uses a **phased reconciliation** pattern:

- **Phase 0 (Prerequisites)**: Namespace creation, sub-operator deployment (e.g., CloudNativePG)
- **Phase 1 (Data Layer)**: PostgreSQL, Valkey (redis-protocol cache)
- **Phase 2 (Supporting Services)**: OpenFGA, Temporal
- **Phase 3 (Application Layer)**: Control plane monolith, web console

Phases are gated: Phase 2 does not start until all Phase 1 components report Ready.

## CRD

```yaml
apiVersion: planton.ai/v1
kind: PlantonPlatform
metadata:
  name: my-planton
spec:
  version: v1.0.0
```

See `config/samples/` for minimal and full examples.

## Development

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `make run`         | Run operator locally against current kubeconfig |
| `make test`        | Run unit tests with envtest                     |
| `make manifests`   | Regenerate CRD and RBAC manifests               |
| `make generate`    | Regenerate DeepCopy methods                     |
| `make kind-create` | Create a Kind cluster                           |
| `make kind-reset`  | Delete + recreate Kind cluster with CRDs        |
| `make kind-deploy` | Build image, load into Kind, deploy in-cluster  |

## Package Structure

Every Go package has its own README.md with design context.

```
api/v1/             CRD type definitions (PlantonPlatform spec/status)
internal/controller/ Reconciliation loop
```

Additional packages are added as the operator grows through its phased implementation (phases/, resources/, status/).

## Related

- [Operator Plan](_projects/20260224.01.sp.planton-operator/tasks/T01_0_plan.md)
- [Infrastructure Architecture Overview](../../../wiki/_legacy/architecture/infrastrcuture/README.md)
- [Architecture Decisions](docs/architecture/infrastrcuture/architecture-decisions.md)
