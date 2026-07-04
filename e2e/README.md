# Planton E2E Test Framework

End-to-end tests that deploy real infrastructure using Planton IaC modules and
verify the results against real providers.

## What This Framework Does

Every Planton component ships with Pulumi and Terraform modules that create
cloud infrastructure. These E2E tests prove that those modules actually work by
executing the full lifecycle against real providers:

1. **VALIDATE** -- load the manifest and build the stack input
2. **DEPLOY** -- run the IaC module (Pulumi up or Terraform apply)
3. **VERIFY-OUT** -- check that stack outputs are populated
4. **VERIFY-RES** -- confirm resources exist using provider-native tools
5. **DESTROY** -- tear down all created resources
6. **VERIFY-CLN** -- confirm resources are gone

If any phase fails, the framework still attempts DESTROY to avoid leaking
resources.

When a component has dependencies (see "Component Dependencies" below), the
framework wraps this lifecycle with a **DEPENDENCIES-UP** phase before VALIDATE
and a **DEPENDENCIES-DOWN** phase after VERIFY-CLN (teardown in reverse order).

## Directory Layout

### Component E2E Structure

Test scenarios, profiles, and fixtures live **next to their components** at the
`v1/e2e/` level:

```
apis/dev/planton/provider/{provider}/{component}/v1/
  e2e/
    profile.yaml           <-- E2E profile (tier, status, provisioners, timeout)
    scenarios/             <-- test scenario manifests
      minimal.yaml
      with-probes.yaml
      with-hpa.yaml
    prerequisite.yaml      <-- optional: this kind's install profile, used when it
                               is itself a prerequisite of another component
  iac/
    hack/manifest.yaml     <-- the canonical example manifest
    pulumi/                <-- Pulumi module
    tf/                    <-- Terraform module
  spec.proto
```

### Provider Harness Structure

Each cloud provider has a harness that manages test infrastructure and
verification, plus a provider-level E2E profile:

```
apis/dev/planton/provider/{provider}/aa_e2e/
  profile.yaml             <-- Provider E2E profile (credentials, substrate, tools)
  harness.go               <-- Provider lifecycle (setup/teardown)
  verify/                  <-- Resource verification logic
```

For Kubernetes, the harness creates a `kind` cluster and uses `kubectl` for
verification.

## Component Dependencies

Some components need other resources installed before they can be applied -- an
operator that owns their CRD, or the CRDs themselves. The harness deploys these
dependencies (via Pulumi) before the component under test and tears them down in
reverse order afterward, resolved by `ResolveDependencies`
([dependencies.go](framework/runner/dependencies.go)) from the proto registry:

Each kind declares its prerequisites in the proto registry
(`CloudResourceKindMeta.prerequisites` in `cloud_resource_kind.proto`). The
harness resolves them transitively and installs each one using, in order of
preference, a consumer-scoped override at the consuming component's
`v1/e2e/prerequisites/<dep>.yaml` (for when the same prerequisite kind needs a
different install shape per consumer — e.g. GcpGlobalAddress as an EXTERNAL
VIP for a forwarding rule vs an INTERNAL VPC_PEERING range for a service
networking connection), then the dependency's `v1/e2e/prerequisite.yaml` (its
published install profile), then its `v1/e2e/scenarios/minimal.yaml`. Declaring
`prerequisites: [X]` is all that is needed -- no per-component wiring.

**Transitive prerequisites resolve against the component under test, not against
intermediate dependencies.** If kind A depends on B and B depends on C, the install
manifest for C is looked up under A's `v1/e2e/prerequisites/c.yaml` (then C's published
profile), NOT under B's consumer-scoped overrides. So when B's install profile
references a specifically shaped C (e.g. a GKE cluster prerequisite referencing a
subnetwork with named secondary ranges), the consuming kind A must ship its own
consumer-scoped copy of that transitive prerequisite with matching `metadata.name`s —
B's consumer-scoped shape does not carry over.

**Cloud SQL chain note:** `GcpCloudSql` declares
`prerequisites: [GcpServiceNetworkingConnection]`, so every instance/database/user
scenario transitively installs VPC → VPC_PEERING range → connection before the
instance. The Cloud SQL kinds ship consumer-scoped overrides for those three
prerequisites with `${E2E_RUN_ID}` in the cloud-side VPC and address names (fixed
`metadata.name` for FK resolution) so parallel runs and stale orphans from a
half-finished teardown do not 409 on recreate.

*Example:* every Gateway API kind declares `KubernetesGatewayApiCrds`, so the
harness installs the Gateway API CRDs (experimental channel, version-pinned)
before applying a GatewayClass / Gateway / route / ReferenceGrant. The Tier 3
operator-dependent components (Postgres, Kafka, ...) likewise declare their
operator kind, which installs from the operator's `scenarios/minimal.yaml`.

### Dependency lifecycle robustness (asynchronous producer cleanup)

Two behaviors in the runner exist because cloud-side cleanup is not always
synchronous with the IaC engine's view of it (first hit: the Cloud SQL →
service networking connection chain):

- **Dependency destroys retry with backoff** (6 attempts, 60s apart, in
  [dependencies.go](framework/runner/dependencies.go)). Deleting a Cloud SQL
  instance returns minutes before the service producer releases its hold on
  the service networking connection; until then the connection destroy fails
  with "Producer services are still using this connection". Skipping past
  that failure is worse than an orphan: the framework would then force-delete
  the VPC, stranding a producer-side connection record that poisons the next
  scenario's same-named prerequisite chain (its connection create silently
  attaches to the stale record and no peering ever materializes).
- **Dependency deploys run `pulumi up --refresh`** ([pulumi.go](framework/runner/pulumi.go)).
  Dependency stacks are keyed by run id, so every scenario in a run reuses the
  same stack name; if an earlier scenario's teardown half-completed, stale
  state would otherwise make a later `up` a silent no-op while the actual
  cloud resource is gone.

Verifiers that read a resource through a *different* API than the one that
created it should poll briefly before declaring it absent (see the
service-networking-connection verifier: the peering is created via the
Service Networking API but read back through the Compute API, and that
cross-API view is eventually consistent).

### Resolving `valueFrom` references in composed scenarios

Before a scenario (or a prerequisite manifest) is applied, the runner resolves
every `valueFrom` reference in the spec against already-deployed prerequisite
outputs ([refresolve.go](framework/runner/refresolve.go)). Resolution walks the
full spec tree — nested messages, repeated message elements such as
`backends[0].group`, and repeated ref fields themselves (a `repeated
StringValueOrRef` like a target HTTPS proxy's `ssl_certificates` resolves each
list element in place) — and honors an explicit `valueFrom.kind` when the
field has no `default_kind` (e.g. a URL map's `default_service` pointing at
either a backend service or a backend bucket). Top-level refs with a
field-level `default_kind` behave exactly as before; existing manifests are
unchanged.

## E2E Profiles

Profiles are KRM-style YAML files (`apiVersion: qa.planton.dev/v1`) that
declare how E2E tests are executed. The CI workflow reads these profiles to
dynamically generate the test matrix -- no hardcoded component lists.

### Provider Profile (`aa_e2e/profile.yaml`)

Configures provider-wide E2E behavior:

```yaml
apiVersion: qa.planton.dev/v1
kind: ProviderE2EProfile
metadata:
  name: kubernetes
spec:
  credential_approach: none
  test_substrate: kind
  default_cost_class: free_local
  default_schedule_lane: weekly
  required_tools: [kind, kubectl, pulumi, tofu]
  github_environment: e2e-kubernetes
  max_concurrent_tests: 8
```

### Component Profile (`v1/e2e/profile.yaml`)

Declares a component's E2E readiness:

```yaml
apiVersion: qa.planton.dev/v1
kind: ComponentE2EProfile
metadata:
  name: kubernetesredis
spec:
  tier: 2
  status: green
  validated_provisioners: [pulumi, terraform]
  timeout_minutes: 15
```

Status values:
- **green** -- passes on CI, included in scheduled runs
- **deferred** -- known failure with documented reason, skipped in CI
- **skip** -- intentionally excluded (needs cloud credentials, etc.)
- **stub** -- module is a stub with no real deployment logic

## Discovering Components

The `planton e2e discover` CLI command scans profiles and displays component
readiness:

```bash
# Interactive TUI (default in terminal)
planton e2e discover --provider kubernetes

# Plain table (default when piped)
planton e2e discover --provider kubernetes --output table

# GitHub Actions matrix JSON (for CI consumption)
planton e2e discover --provider kubernetes --output github-matrix

# Filter to GREEN Pulumi Tier 1 only
planton e2e discover --provider kubernetes --status green --tier 1 --provisioner pulumi
```

## Running Tests

Prerequisites: Docker running, `kind`, `kubectl` installed, plus at least one
of `pulumi` (for Pulumi E2E) or `tofu`/`terraform` (for Terraform E2E).

```bash
# All Kubernetes E2E tests (Pulumi + Terraform, all tiers)
make e2e-test-kubernetes

# Pulumi-only, single component
make e2e-test-component component=KubernetesNamespace_Pulumi

# Terraform-only, Tier 1
make e2e-test-kubernetes-terraform-tier1

# Terraform-only, single component
go test -tags=e2e -timeout=30m -v -count=1 \
  -run "TestKubernetesNamespace_Terraform/minimal$" ./e2e/...
```

### Terraform binary selection

Terraform E2E defaults to `tofu` (OpenTofu), matching the Planton CLI.
To use HashiCorp Terraform instead:

```bash
PLANTON_E2E_TF_BINARY=terraform make e2e-test-kubernetes-terraform-tier1
```

### How the Terraform path works

The Terraform runner uses [Terratest](https://github.com/gruntwork-io/terratest)
as its execution layer. For each test scenario:

1. The TF module (`iac/tf/`) is copied to a temp directory
2. `terraform.tfvars` is generated from the manifest proto via `ProtoToTFVars()`
3. `backend.tf` is written with a local backend
4. Provider env vars (KUBECONFIG, etc.) are extracted from the stack-input YAML
5. Terratest runs `tofu init` + `tofu apply` with built-in transient error retry
6. The same kubectl verifiers validate the deployed infrastructure
7. Terratest runs `tofu destroy`
8. The same kubectl verifiers confirm cleanup
9. The temp directory is removed

## CI Workflow

The `e2e-kubernetes.yaml` GitHub Actions workflow automates E2E on a weekly
schedule:

1. **build-check** -- compiles E2E code + go vet (runs on every PR too)
2. **discover** -- runs `planton e2e discover --output github-matrix` to
   generate the test matrix from profiles
3. **e2e** -- dynamic matrix of (tier x provisioner) cells, each with its own
   kind cluster, using `gotestsum` for JUnit output
4. **summary** -- aggregates JUnit XML into GitHub Step Summary

To trigger manually: Actions > e2e-kubernetes > Run workflow > select branch.

## GCP E2E

GCP tests live under `e2e/gcp/` and use the shared `aa_e2e` harness with
real GCP API verification. Set `E2E_GCP_PROJECT` to a dedicated test project
with Application Default Credentials.

```bash
E2E_GCP_PROJECT=planton-e2e go test -tags=e2e -timeout=120m -v ./e2e/gcp/...
```

**Backend contract:** every scenario's test context must carry the run-scoped
Pulumi file backend URL — even Terraform scenarios, because dependency
prerequisites always deploy via Pulumi. An empty backend URL silently falls
back to the machine's ambient `pulumi login`, coupling runs to developer
state that can vanish mid-run (e.g. a stale `/tmp` backend).

**Timeout guidance:** scenarios with private services access prerequisite chains
(VPC → global address → service networking connection) or Memorystore
create/destroy need substantial headroom. Budget **≥35m per PSA-backed scenario**
and **≥90m** when batching multiple PSA or Memorystore scenarios in one `go test`
run — Redis instance destroy alone can exceed 15 minutes. GKE clusters are the
slowest resources in the harness: a control-plane create runs 10-25 minutes and
a destroy 5-10, per scenario per engine — budget **≥150m** when batching the
GKE cluster scenarios across both engines.

## Build Tag Isolation

All E2E test files use `//go:build e2e`. This means:

- `go test ./...` and `make test` never trigger E2E tests
- `go build ./...` never compiles E2E test binaries
- You must pass `-tags=e2e` explicitly to run them

The framework packages under `e2e/framework/` have **no** build tag -- they are
ordinary Go libraries that get compiled normally. Only the test files that
create real infrastructure are gated.

## How Verification Works

The framework does not hardcode resource names. Instead, it parses each test
manifest at runtime to extract the resource name, namespace, and kind, then
builds the appropriate verification dynamically. This means adding a new test
scenario is as simple as dropping a YAML file into the component's
`e2e/scenarios/` folder -- no Go code changes needed.

## Adding a New Test Scenario

1. Create a YAML manifest in `{component}/v1/e2e/scenarios/` with a descriptive
   filename
2. Use a unique `metadata.name` (and unique namespace if the component creates
   one) to avoid collisions with other scenarios
3. Run `make e2e-test-component component={ComponentName}` to verify it works
4. That's it -- the framework discovers and runs it automatically

## Adding a New Component

1. Create the IaC modules (`iac/pulumi/`, `iac/tf/`)
2. Create `v1/e2e/profile.yaml` with the component's E2E profile
3. Create at least `v1/e2e/scenarios/minimal.yaml` with a minimal test manifest
4. If the component needs other resources installed first, declare them as
   `prerequisites` on the kind in `cloud_resource_kind.proto` (the harness
   installs them automatically -- see "Component Dependencies").
5. Add a `Test{ComponentName}_{Provisioner}` function in the appropriate test
   file (e.g., `kubernetes_test.go`), and -- if the component name does not
   PascalCase trivially -- a `toPascalCase` entry in
   `pkg/e2e/profile/discover.go` so the CI matrix regex matches it
6. The CI workflow picks up the new component automatically from the profile

## Adding a New Provider

1. Create `apis/dev/planton/provider/{provider}/aa_e2e/` with harness, verify
   files, and `profile.yaml`
2. Implement the `provider.Harness` interface (Setup, Teardown, VerifyDeployed,
   VerifyDestroyed)
3. Add a test entry point in `e2e/` that creates the harness and discovers
   scenarios for that provider
4. Add Makefile targets
5. Create `.github/workflows/e2e-{provider}.yaml` with the appropriate trigger
   schedule and credential configuration

## Architecture

```
e2e/
  e2e_test.go             -- TestMain: shared infrastructure lifecycle
  kubernetes_test.go      -- Kubernetes test entry points (per-component)
  framework/
    runner/               -- 6-phase lifecycle engine, Pulumi/Terraform execution
    provider/             -- Harness interface definition
    discovery/            -- Filesystem scanner for components and scenarios
    reporter/             -- JSON + Markdown report generation

pkg/e2e/profile/          -- E2E profile loader and discovery
  loader.go               -- YAML→proto loading for provider and component profiles
  discover.go             -- Profile scanning, filtering, GitHub matrix generation
  paths.go                -- Well-known filesystem paths

apis/dev/planton/qa/      -- Proto schema for E2E profiles (KRM-style)
  shared/                 -- Shared enums (CostClass)
  providere2eprofile/v1/  -- ProviderE2EProfile KRM API
  componente2eprofile/v1/ -- ComponentE2EProfile KRM API
```

The framework is engine-agnostic. The runner supports both Pulumi and Terraform
execution paths. Each component test runs through the same lifecycle regardless
of which engine is used.
