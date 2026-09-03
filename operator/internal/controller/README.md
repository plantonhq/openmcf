# internal/controller -- Phase Executor

This package implements the main reconciliation loop for the `PlantonPlatform` custom resource. It acts as a thin orchestrator that delegates work to phases defined in `internal/phases/`.

## How It Works

Each reconciliation follows this sequence:

1. Fetch the PlantonPlatform CR (handle deletion gracefully)
2. Initialize status if this is a fresh resource (delegate to `internal/status/`)
3. Run phases in order: Prerequisites, DataLayer, (future: SupportingServices, Application)
4. If a phase is not ready, update the overall phase and requeue after 30 seconds
5. If all phases are ready, compute and set the overall phase to Ready

The controller does not contain deployment logic itself. Each phase owns its resources and reports readiness through the Phase interface.

## Design Decisions

### Phase Gating

Phases execute sequentially. The DataLayer phase does not run until Prerequisites reports Ready. This ensures sub-operators are healthy before we create CRDs that depend on them. The overhead is negligible since phases check status (fast) rather than block on readiness.

### Status Update Strategy

Status is updated at two points:
- After initialization (immediate requeue to start phases)
- When the overall phase changes (computed from component phases via `status.ComputeOverallPhase`)

Phases modify the in-memory status struct. The controller performs the API server write. This minimizes API calls per reconciliation.

### Requeue Strategy

- **After initialization**: Immediate requeue (`Requeue: true`) to start phases without waiting
- **Phase not ready**: Requeue after 30 seconds (`RequeueAfter: 30s`)
- **All phases ready**: Requeue after 30 seconds for ongoing health monitoring

As the operator matures, the interval will adapt -- shorter when actively deploying, longer when everything is healthy.

## Testing

Tests use envtest (lightweight Kubernetes API server) via the Kubebuilder test harness. They verify:
- Status initialization on first reconcile
- Component and condition setup
- Phase execution and requeue behavior
- Graceful handling of deleted resources
