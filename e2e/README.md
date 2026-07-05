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
preference, the dependency's `v1/e2e/prerequisite.yaml` (its published install
profile) or its `v1/e2e/scenarios/minimal.yaml`. Declaring `prerequisites: [X]`
is all that is needed -- no per-component wiring.
*Example:* every Gateway API kind declares `KubernetesGatewayApiCrds`, so the
harness installs the Gateway API CRDs (experimental channel, version-pinned)
before applying a GatewayClass / Gateway / route / ReferenceGrant. The Tier 3
operator-dependent components (Postgres, Kafka, ...) likewise declare their
operator kind, which installs from the operator's `scenarios/minimal.yaml`.

### Scenario-declared extra fixtures (optional composition seams)

Registry `prerequisites` carry a strict meaning -- the parents a resource cannot
exist without -- and they double as deploy-ordering metadata, so **optional**
composition seams must never be encoded there: adding an optional kind to a
registry prerequisite list would force every downstream kind's fixture chain to
deploy it forever. When a scenario needs to live-prove an optional edge (a
subnet attaching a route table, a NAT gateway associating a public IP, a
network peering to a second network), it declares the extra fixtures itself via
the `planton.dev/e2e-extra-prerequisites` annotation on the scenario manifest:

```yaml
metadata:
  annotations:
    # Kind names install through the kind's standard install profile;
    # repo-relative manifest paths deploy an EXTRA INSTANCE of their declared
    # kind (for scenarios needing more instances than the profiles provide).
    planton.dev/e2e-extra-prerequisites: "AzureRouteTable, AzureNetworkSecurityGroup"
```

Entries deploy in listed order after the registry prerequisites, each preceded
by any of its own transitive registry prerequisites not already deployed
(shared parents are deduplicated -- one resource group serves the whole chain).
Kind-name entries already present in the registry chain are skipped; path
entries always deploy, because they exist precisely to add another instance.
All fixtures join the same transitive `value_from` reference resolution the
registry chain uses, and teardown runs in reverse across the merged chain.
Every kind that appears in the annotation needs a verifier and an install
profile, exactly like a registry prerequisite.

A dependency whose `pulumi up` FAILS is still tracked for teardown: a failed
update may have created any number of resources before erroring, and skipping
its destroy would orphan them -- and, because Azure-style parents refuse to
delete while children exist, a single orphaned fixture (say, a load balancer
holding a frontend in the fixture subnet) blocks the entire reverse teardown
chain behind it. Destroying a stack whose update failed is safe; it removes
whatever was actually created.

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

### Long-running Azure components (AKS)

AKS clusters take roughly 5–10 minutes to create and a similar time to delete.
Component E2E profiles for `azureakscluster` and `azureaksnodepool` set
`timeout_minutes: 60–75`. When invoking tests directly, size `-timeout` beyond
the default 30m:

```bash
go test -tags=e2e -timeout=90m -v -count=1 \
  -run 'TestAzureAksCluster_Pulumi/minimal' ./e2e/azure/...
```

Burstable VM sizes in `eastus` may support only availability zone `1` — multi-zone
lists fail with `AvailabilityZoneNotSupported`. AKS E2E scenarios in this repo
use `zones: ["1"]` for the test subscription.

### Offer-restricted services on free/PAYG subscriptions (probe, don't roulette)

Some Azure database services restrict provisioning PER REGION on free-tier
and new pay-as-you-go subscriptions: the ARM API returns
`LocationIsOfferRestricted` (PostgreSQL/MySQL Flexible Server) or
`ProvisioningDisabled` (Microsoft.Sql logical servers) in the restricted
regions. Verified on the test subscription: `eastus`, `eastus2`, and
`westus2` are blocked for PostgreSQL while `westus3`, `centralus`,
`canadacentral`, `northeurope`, and `uksouth` are clean — the restriction
is REGIONAL, not subscription-wide, so do not record a deferral after two
failures; find a clean region instead.

For PostgreSQL the per-region flag is queryable UP FRONT — probe before
picking (or moving) a scenario region:

```bash
az rest --method get \
  --url "https://management.azure.com/subscriptions/{sub}/providers/Microsoft.DBforPostgreSQL/locations/{region}/capabilities?api-version=2024-08-01" \
  --query "value[0].restricted"   # "Disabled" == provisioning allowed
```

Each RDBMS carries its OWN restriction footprint — do not assume they
match. Verified live on the test subscription: `westus3` is clean for
PostgreSQL but blocked for MySQL (`ProvisionNotSupportedForRegion`),
while `westus2` is blocked for PostgreSQL but clean for MySQL. MySQL's
capabilities endpoint is itself the probe, just with a different signal:
a restricted region answers 500 `InternalServerError`, a usable region
returns its edition ladder:

```bash
az rest --method get \
  --url "https://management.azure.com/subscriptions/{sub}/providers/Microsoft.DBforMySQL/locations/{region}/capabilities?api-version=2023-12-30" \
  --query "value[0].supportedFlexibleServerEditions[].name"
```

Microsoft.Sql has no public equivalent — fall back to one live attempt.
A server may live in a different region than its resource group, so the
shared fixture RG serves unchanged when a scenario pins a different
region. True subscription-wide gates do exist (the quota-increase link
in the ARM error is the fix) — but conclude that only after a
probe-verified clean region also fails.

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

### Real-cloud harness Setup: validate credentials and export preconditions

For a real-cloud provider (`test_substrate: real_cloud`), the framework builds
every stack input with a **nil provider config**, so the IaC modules resolve
credentials from the SDK's ambient chain (a keyless CLI/SSO login locally, OIDC
federation in CI) rather than from a stored secret. The harness `Setup` therefore
owns two responsibilities beyond wiring verifiers:

- **Fail fast on credentials.** Probe the ambient chain with a zero-permission,
  side-effect-free identity call (e.g. AWS `sts:GetCallerIdentity`; Azure: acquire
  a management token) so a missing/expired login is reported before any deploy.
- **Export the environment preconditions both engines need.** The Terraform path's
  provider block is empty and the Pulumi builders leave unset args to env-var
  fallbacks, so any value the provider cannot infer must be exported into the
  process environment (the Pulumi runner inherits `os.Environ()`; the Terraform
  path layers extracted vars on top of it). Two categories recur:
  - *Identity/scope the provider cannot infer* — e.g. Azure `azurerm` v4 no longer
    infers the subscription, so the harness must export `ARM_SUBSCRIPTION_ID`.
  - *Test-environment behavior that differs from production defaults* — e.g. the
    Azure providers auto-register a broad set of resource providers at init and
    fire the registrations concurrently, which returns HTTP 409 on a subscription
    whose providers are not yet registered. Resource-provider registration is a
    one-time subscription bootstrap, orthogonal to whether a module creates its
    resource correctly (the contract E2E validates), so the harness opts out for
    the ephemeral run via `ARM_SKIP_PROVIDER_REGISTRATION=true`. Keep such opt-outs
    scoped to the harness and documented, so the test stays honest about what it
    proves.

### Authoring verifiers and prerequisite fixtures

- **Every kind that appears in another kind's `kind_meta.prerequisites` needs its
  own verifier registered in the provider harness** — the dependency deployer
  verifies each fixture right after installing it, so a missing verifier fails the
  composed scenario at DEPENDENCIES-UP with "no verifier registered", not at the
  component under test. Wiring a new composed component therefore means wiring
  verifiers for its whole prerequisite chain (plus a `prerequisite.yaml` or minimal
  scenario for each prerequisite).
- **Prefer a GET-by-ID existence probe over a HEAD unless the service is known to
  support HEAD.** ARM's generic `CheckExistenceByID` (HEAD) is not implemented by
  every resource provider — e.g. `Microsoft.ManagedIdentity` answers HEAD with
  405 Method Not Allowed while GET works fine. A GET with the typed 404 as the
  absence signal works everywhere; treat every non-404 failure as a real error so
  auth/network problems never masquerade as "absent".
- **Data-plane resources need a data-plane grant AND sometimes a data-plane
  verifier.** Some resources live behind a service's own endpoint rather than the
  control plane (e.g. Azure Key Vault keys/certificates behind
  `{vault}.vault.azure.net`): the deploying credential needs an explicit
  data-plane authorization even when it owns the subscription. Grant it once at
  the subscription scope as an idempotent harness-Setup bootstrap rather than
  per-ephemeral-resource — data-plane RBAC takes minutes to propagate, so a
  per-run grant makes every scenario race its own authorization. For
  verification, check whether the control plane exposes a read proxy first
  (Azure vault KEYS have one; CERTIFICATES do not) and fall back to the service's
  data-plane SDK only when it does not.
- **Soft-delete/retention services add an orphan class the resource list does not
  show.** A destroyed Azure Key Vault lingers soft-deleted (its globally unique
  name stays reserved) unless purged; both IaC engines purge on destroy by
  default, but an interrupted run can strand one. The zero-orphan sweep for such
  services must check the recycle bin too (`az keyvault list-deleted`), and
  scenarios should keep purge protection OFF so teardown can actually purge.
  Expect destroys to be slow (a vault purge runs ~10 minutes) and size test
  timeouts accordingly.
- **Some ARM resources have no true delete — destroy flips a state field and the
  object stays GETtable, so verify-absent must be STATE-AWARE, not 404-based.**
  An Azure Storage encryption scope is the canonical case: "delete" PATCHes
  `properties.state` to `Disabled`, a GET keeps answering 200, and the name stays
  reserved inside its parent (recreating the same name re-enables it). A plain
  GET-with-typed-404 probe reports such a resource as still-existing after a
  clean destroy and fails the run spuriously. Check the provider's own Read/Delete
  source for state-flip semantics when writing a verifier: if delete is a
  soft-disable, the verifier must read the state field and treat the disabled
  value as absent (mirroring the provider's own removed-from-state behavior),
  and verify-exists should require the ENABLED state, not mere presence. No
  recycle-bin sweep exists for this class — the object intentionally persists;
  parent teardown is what actually removes it.
- **Never let sequential scenarios destroy and recreate the same globally unique
  parent name.** When a kind's registry prerequisite chain deploys a fixture
  whose name is globally unique (an Azure SQL logical server, a Key Vault),
  every scenario of a multi-scenario component tears the fixture down and the
  next scenario recreates it — and Azure can hold the just-deleted name long
  enough that the recreate hangs indefinitely (a `Microsoft.Sql/servers` create
  stuck 20+ minutes with no write in the activity log). Give each scenario its
  own uniquely named parent instead: declare the parent through the
  `e2e-extra-prerequisites` annotation with a scenario-local manifest (kept
  OUTSIDE `e2e/scenarios/`, which the discoverer treats as test cases), and
  drop the registry prerequisite if it would force the shared fixture chain in
  anyway. Registry prerequisites are for parents a kind cannot exist without
  AND that are safe to recreate per scenario.
- **The tfvars wire format drops zero-valued proto fields — TF object attributes
  where zero is meaningful must be `optional()` with the zero default.**
  `ProtoToTFVars()` serializes from protojson, which omits scalar zeros
  (proto3 implicit presence). A required attribute like
  `per_database_settings.min_capacity = number` then fails
  `terraform apply` with "attribute is required" whenever the manifest
  legitimately sets `0`. Declare such attributes
  `optional(number, 0)` so the dropped zero round-trips.

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
