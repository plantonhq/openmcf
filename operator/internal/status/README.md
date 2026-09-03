# internal/status

Status helpers for the PlantonPlatform CRD.

## Responsibility

This package owns all mutations to `PlantonPlatformStatus`. It provides pure functions that modify the status struct in memory. The caller (controller) is responsible for writing the updated status to the API server.

## Key Functions

- **Initialize**: Sets the status to Pending with all 8 component statuses and 3 aggregate conditions. Idempotent -- returns false if already initialized.
- **SetCondition**: Sets an aggregate condition (DataLayerReady, SupportingServicesReady, ApplicationReady).
- **SetComponentPhase**: Sets a component's phase and message.
- **ComputeOverallPhase**: Derives the overall deployment phase from all component phases (Error > Ready > Deploying > Pending).

## Design Decisions

**Pure functions, not methods**: Status helpers take a `*PlantonPlatform` parameter rather than being methods on a status manager. This keeps the functions testable without a Kubernetes client and avoids premature abstraction.

**Single status update per reconcile**: Phases modify the in-memory status throughout the reconciliation loop. The controller performs one `Status().Update()` at the end. This avoids multiple API calls and reduces conflict risk.
