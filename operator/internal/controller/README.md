# internal/controller -- the reconcile loop

This package implements the reconciliation loop for the `PlantonPlatform` custom resource. It is a thin orchestrator: it decides the order and the gates, and delegates every deployment decision to the components registered in `internal/component/` and every status mutation to `internal/status/`.

## How It Works

Each reconciliation follows this sequence:

1. Fetch the `PlantonPlatform` (a deleted resource ends the reconcile quietly).
2. Initialize status on a fresh resource or when the spec has changed shape (component slots, the version echo, the toggled slots); write it and requeue immediately.
3. Judge the declared platform version against the operator's floor (`internal/platformversion`). A version this operator cannot run is refused whole: phase `Error`, the `VersionSupported` and `Ready` conditions carry the explanation, no component is touched, and the reconcile ends without a requeue -- a spec change re-enqueues on its own.
4. Walk every registered component. A component whose dependencies are not yet Ready is marked Pending with the dependency named; the rest reconcile their own resources and report Ready, Deploying, or Error with a message.
5. Compute the overall phase from the component phases, set the `Ready` condition, write the status once, and requeue after 30 seconds.

`PlantonIdentityProvider` deliberately has no controller of its own: a change to one re-enqueues the platforms in its namespace, and the identity component resolves the binding inside the same loop -- one loop, one cadence, no second writer.

## Design Decisions

### One status write per reconcile

Components mutate the in-memory status as they run; the controller performs a single `Status().Update()` at the end. Fewer API calls, fewer conflicts.

### The version floor runs here, not only at admission

The API server enforces the shape of `spec.version` (a release as `vMAJOR.MINOR.PATCH`) through the CRD. Whether this operator RUNS that release is a property of the operator binary, and an operator upgrade can outgrow a platform that is already running -- no admission rule sees that moment. So the reconciler judges the version on every pass and explains a refusal where a person reads first: the `MESSAGE` column of `kubectl get plantonplatform`.

### Requeue Strategy

- **After initialization**: immediate requeue (`Requeue: true`) so components start without waiting.
- **Version refused**: no requeue. Nothing changes until the spec does, and the spec change is itself the trigger.
- **Otherwise**: requeue after 30 seconds, whether deploying or healthy, for ongoing convergence and health monitoring.

## Testing

Tests use envtest (a real API server and etcd, the actual CRD applied) through the Kubebuilder harness. They verify status initialization on first reconcile, component and condition setup, requeue behavior, graceful handling of deleted resources, the CRD's validation rules (the messages the API server prints), and both branches of the version floor.
