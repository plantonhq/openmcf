# internal/status

Status helpers for the PlantonPlatform CRD.

## Responsibility

This package owns all mutations to `PlantonPlatformStatus`. It provides pure functions that modify the status struct in memory. The caller (controller) is responsible for writing the updated status to the API server.

## Key Functions

- **Initialize**: Sets the status to Pending, allocates a status slot for every enabled component (and retires the slot of a component toggled off), echoes the declared version and license mode, and sets `Ready=False`. Idempotent -- returns false when nothing changed.
- **RefuseVersion**: Records that the operator will not run the declared platform version: phase `Error`, `VersionSupported=False` with the finer reason, `Ready=False` with the same plain-language message (the one the `MESSAGE` column prints). Component statuses are left alone so a running platform keeps reporting what is running. Returns false when the refusal is already recorded, so it costs no status write.
- **SetCondition**: Sets a condition (`Ready`, `VersionSupported`).
- **SetComponentPhase**: Sets a component's phase and message.
- **UpdateReadyCondition** / **ComputeOverallPhase**: Derive the `Ready` condition and the overall deployment phase from all component phases (Error > Ready > Deploying > Pending).

## Design Decisions

**Pure functions, not methods**: Status helpers take a `*PlantonPlatform` parameter rather than being methods on a status manager. This keeps the functions testable without a Kubernetes client and avoids premature abstraction.

**Single status update per reconcile**: Components modify the in-memory status throughout the reconciliation loop. The controller performs one `Status().Update()` at the end. This avoids multiple API calls and reduces conflict risk.
