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

**IDEMPOTENCY phase (per-provider opt-in):** when a provider's E2E profile
sets `assert_apply_idempotency: true`, every lifecycle gains an IDEMPOTENCY
phase right after DEPLOY: the runner re-plans the just-applied configuration
(`terraform plan -detailed-exitcode` must exit 0 / `pulumi preview
--expect-no-changes`) and fails the lane on any pending change. A dirty
second plan means the module and the provider disagree about applied state —
the send-omitted-value and Optional+Computed echo defect classes, which
users otherwise meet as a perpetual diff on every re-apply (first live
catch: the Identity Platform config's server-materialized
`sign_in.phone_number` block). The gate covers the component under test
only; prerequisite fixtures belong to other kinds' contracts. Arm it per
provider after the catalog's known no-op re-plan classes are burned down.

## Directory Layout

### Component E2E Structure

Test scenarios, profiles, and fixtures live **next to their components** at the
component root's `e2e/` level:

```
catalog/{provider}/{component}/
  e2e/
    manifest.yaml          <-- the canonical validated example manifest
    profile.yaml           <-- E2E profile (tier, status, provisioners, timeout)
    scenarios/             <-- test scenario manifests
      minimal.yaml
      with-probes.yaml
      with-hpa.yaml
    prerequisite.yaml      <-- optional: this kind's install profile, used when it
                               is itself a prerequisite of another component
  iac/
    pulumi/                <-- Pulumi module
    tf/                    <-- Terraform module
  v1alpha1/
    spec.proto
```

### Provider Harness Structure

Each cloud provider has a harness that manages test infrastructure and
verification, plus a provider-level E2E profile:

```
catalog/{provider}/aa_e2e/
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
`e2e/prerequisites/<dep>.yaml` (for when the same prerequisite kind needs a
different install shape per consumer — e.g. GcpGlobalAddress as an EXTERNAL
VIP for a forwarding rule vs an INTERNAL VPC_PEERING range for a service
networking connection), then the dependency's `e2e/prerequisite.yaml` (its
published install profile), then its `e2e/scenarios/minimal.yaml`. Declaring
`prerequisites: [X]` is all that is needed -- no per-component wiring.

**Transitive prerequisites resolve against the component under test, not against
intermediate dependencies.** If kind A depends on B and B depends on C, the install
manifest for C is looked up under A's `e2e/prerequisites/c.yaml` (then C's published
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

**Undeletable resource classes redefine "zero orphans".** Some GCP resources
have NO delete API — KMS key rings and crypto keys are the canonical class:
destroy removes a ring from state only, and destroys a key's *versions*
(disabling rotation) while the key object persists forever. For these
classes: (1) every cloud-side name MUST carry `${E2E_RUN_ID}` — the leftover
objects are permanent, so a fixed name can never deploy twice; (2) the
verifier's absent-check asserts the honest destroyed posture (all key
versions `DESTROYED`/`DESTROY_SCHEDULED`, rotation off) instead of absence;
(3) the post-run sweep expectation is "no ACTIVE material", not "no objects"
— run-scoped rings/keys accumulate in the test project as inert, zero-cost
residue by GCP design. Do not hand-sweep them; there is nothing to sweep.
(4) A component whose PREREQUISITE is undeletable gets ONE live scenario:
prerequisites redeploy per scenario with the same engine-scoped run id, so
a second scenario re-creates the just-"destroyed" (state-only) prerequisite
and 409s on its own leftover. Fold the arms into one scenario and record
the rest as offline-proven.

**Deletes that finalize asynchronously make FIXED cloud-side names 409 across
engines.** Some control planes acknowledge a delete before the name is
actually reusable (first hit: Cloud Scheduler jobs — the destroy returns, the
verifier's GET already 404s, yet a create of the same name minutes later can
still 409 `already exists`; Cloud Tasks goes further and documents a queue-ID
reservation of up to 7 days after deletion). The dual-engine runner recreates
every scenario back-to-back under both engines, so any scenario whose
cloud-side name is FIXED — including one that deliberately omits the explicit
name to prove the metadata.name fallback — sits exactly in this window. The
rule: the fallback-proof scenario carries `${E2E_RUN_ID}` in `metadata.name`
itself (the fallback then makes it the cloud-side name, proving the contract
without a fixed identifier); everything else carries the token in the spec's
explicit name field.

**Identifier classes that forbid hyphens take the `${E2E_RUN_ID_UNDERSCORE}`
token.** The engine-scoped run id carries a hyphenated engine suffix
(`-p`/`-t`), so the plain token cannot be embedded in fields whose character
class is letters/numbers/underscores only (first hit: a Vertex AI
`deployed_index_id`, which must start with a letter and forbids hyphens).
The runner expands `${E2E_RUN_ID_UNDERSCORE}` with every hyphen replaced by
an underscore ([tokens.go](framework/runner/tokens.go)); use it whenever a
run-scoped identifier lives in an underscore-only field. Do NOT fall back to
a fixed identifier to dodge the character class — see the next paragraph for
why that 400s.

**Some holds are keyed by a user-chosen SUB-resource ID and survive deleting
the parent.** A Vertex AI DeployedIndex that is in a failed state or still
undeploying holds its user-chosen `deployed_index_id` with a 400 ("It will
be deleted automatically, please retry again later or use a different ID")
— and the hold survives DELETING THE PARENT index endpoint, so neither the
dual-engine back-to-back recreate nor a fresh prerequisite chain escapes a
fixed ID. This is a distinct class from the fixed-NAME 409s above: the held
identifier is a spec field, not a resource name, so the run-scoped token
must go into the spec field itself (`deployedIndexId:
planton_oss_e2e_${E2E_RUN_ID_UNDERSCORE}` — the live-proven shape).

**Retries cannot cover releases the cloud documents in HOURS.** The retry
budget is bounded (~6 minutes) on purpose: some cloud-side holds outlive any
reasonable wait — a direct-VPC Cloud Run service leaves a serverless address
reservation (`serverless-ipv4-*`) in its subnetwork for **1-2 hours** after
the service is destroyed, and until GCP garbage-collects it the subnetwork
destroy fails with `resourceInUseByAnotherResource` (the reservation itself
cannot be deleted — it is held by the serverless service agent). A scenario
whose prerequisite teardown depends on such a release must NOT run live:
record it as an E2E exclusion in the component's `e2e/profile.yaml` with the
specific reason, and prove the surface offline instead. Fixed prerequisite
names make this worse (a stranded subnet collides with the next run), so any
scenario in an async-release blast radius should carry `${E2E_RUN_ID}` in its
prerequisites' cloud-side names.

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

### Scenario-declared extra fixtures (optional composition seams)

Registry `prerequisites` carry a strict meaning -- the parents a resource cannot
exist without -- and they double as deploy-ordering metadata, so **optional**
composition seams must never be encoded there: adding an optional kind to a
registry prerequisite list would force every downstream kind's fixture chain to
deploy it forever. When a scenario needs to live-prove an optional edge (a
subnet attaching a route table, a NAT gateway associating a public IP, a
network peering to a second network), it declares the extra fixtures itself via
the `planton.dev/e2e-prerequisites` annotation on the scenario manifest:

```yaml
metadata:
  annotations:
    # Kind names install through the kind's standard install profile;
    # repo-relative manifest paths deploy an EXTRA INSTANCE of their declared
    # kind (for scenarios needing more instances than the profiles provide).
    planton.dev/e2e-prerequisites: "AzureRouteTable, AzureNetworkSecurityGroup"
```

Kind-name entries join the registry prerequisite graph and deploy in
topological order, expanded through their own prerequisite edges and
deduplicated against the chain (one resource group serves everything);
entries the chain already deploys are skipped. Manifest-path entries deploy
in listed order after the kind-driven chain, each preceded by any of its own
transitive prerequisites not already deployed -- and always deploy, because
they exist precisely to add another instance; a path entry never substitutes
for the kind's install profile. The resolver also honors the same annotation
(kind names only) on each prerequisite's OWN install manifest, ordering the
fixtures an install profile's `value_from` references compose BEFORE the
declaring kind -- recursively and cycle-checked. All fixtures join the same
transitive `value_from` reference resolution the registry chain uses, and
teardown runs in reverse across the merged chain. Every kind that appears in
the annotation needs a verifier and an install profile, exactly like a
registry prerequisite.

A dependency whose `pulumi up` FAILS is still tracked for teardown: a failed
update may have created any number of resources before erroring, and skipping
its destroy would orphan them -- and, because Azure-style parents refuse to
delete while children exist, a single orphaned fixture (say, a load balancer
holding a frontend in the fixture subnet) blocks the entire reverse teardown
chain behind it. Destroying a stack whose update failed is safe; it removes
whatever was actually created.

### Bare polymorphic references need an explicit `kind:` in scenario valueFrom

Reference resolution determines the referenced kind from the `valueFrom.kind`
field, falling back to the spec field's `default_kind` annotation. A **bare
polymorphic reference** (a `StringValueOrRef` deliberately carrying NO
`default_kind` because no kind dominates -- alias-record targets, diagnostic
setting targets, metric-alert scopes) has no fallback, and the resolver leaves
a kind-less reference on such a field UNTOUCHED rather than erroring. The
module then receives an unresolved ref, sends nothing (or an empty string),
and the failure surfaces at DEPLOY as a provider validation error that reads
like a module defect ("one of either X or Y must be specified"). In scenario
manifests, every reference on a bare polymorphic field must therefore carry
the explicit `kind:` alongside `name:` and `fieldPath:` -- no offline gate
catches the omission, because the manifest validates and the plan renders
without it.

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

### Component Profile (`e2e/profile.yaml`)

Declares a component's E2E readiness:

```yaml
apiVersion: qa.planton.dev/v1
kind: ComponentE2EProfile
metadata:
  name: kubernetesvalkey
spec:
  tier: 1
  status: green
  validated_provisioners: [pulumi, terraform]
  timeout_minutes: 25
```

Status values:
- **green** -- passes on CI, included in scheduled runs
- **deferred** -- known failure with documented reason, skipped in CI
- **skip** -- intentionally excluded (needs cloud credentials, etc.)
- **stub** -- module is a stub with no real deployment logic

The status is enforced at two layers: CI matrices are built from `planton
e2e discover` filters, and the provider test runners load the component's
profile and `t.Skip` any non-green component (with the profile's
`deferred_reason` in the skip message) -- so a full-provider suite run
never fails on a documented deferral.

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

### Pulumi CLI minimum version: the destroy path passes `--run-program`

The runner's `pulumi destroy` invocations pass `--run-program` (it keeps the
program available during destroy so BeforeDelete/AfterDelete resource hooks
fire). Older Pulumi CLIs do not know the flag, and the failure mode is nasty:
every phase up to and including VERIFY-RES passes, then DESTROY fails
instantly with `unknown flag: --run-program` -- so the lane fails AFTER
creating real cloud resources, whose stack state lives in the run's temp
backend and is discarded when the process exits. The resources must then be
swept by hand (`az group list` / the provider's own list commands) before a
re-run. Verified live: v3.137.0 fails exactly this way; v3.256.0 works.
Check `pulumi destroy --help | grep run-program` before the first lane on
any machine, and upgrade the CLI rather than editing the runner -- the flag
is load-bearing for delete-hook correctness.

### Sensitive stack outputs: the Pulumi output reader passes `--show-secrets`

The runner reads Pulumi outputs with `pulumi stack output --json
--show-secrets`. The flag is load-bearing: without it Pulumi masks every
secret output as the literal string `[secret]`, which corrupts a
sensitive SCALAR output silently (the sentinel populates the proto field
and counts as verified) and fails a sensitive MAP output's JSON parse
with `invalid character 's' looking for beginning of value` (first live
hit: a name-keyed `authorization_keys` map — the transform skipped the
field and VERIFY-OUT under-reported N-1/N populated). The Terraform path
reads real values (`terraform output -json` includes sensitive outputs),
so the flag keeps both engines' verification bar identical. Dependency
output capture rides the same helper, so without the flag a scenario's
`value_from` reference to a FIXTURE's sensitive output resolves to the
sentinel and deploys silently wrong. Runner code never prints raw output
values (counts and names only), so unmasked values stay out of logs —
keep it that way when touching output-path logging.

### Terraform binary selection

Terraform E2E defaults to `tofu` (OpenTofu), matching the Planton CLI.
To use HashiCorp Terraform instead:

```bash
PLANTON_E2E_TF_BINARY=terraform make e2e-test-kubernetes-terraform-tier1
```

### Background suite lanes: smoke-check the launch, never trust nohup blindly

Long real-cloud suites are best launched as background processes writing to a
log file (a crashed observer then loses only observation, never the run). One
failure mode recurs: a `nohup`-spawned child of a short-lived shell can be
reaped with its parent BEFORE the suite starts, leaving a zero-byte log and no
error -- the launch "succeeds" and nothing runs. After launching any background
lane, smoke-check it within a minute: the PID must still be alive AND the log
must be non-empty (the harness prints its authentication line first). If the
process is gone with an empty log, relaunch under a supervised/managed shell
rather than retrying `nohup` from another transient shell.

### Long-running Azure components (AKS)

AKS clusters take roughly 5–10 minutes to create and a similar time to delete.
Component E2E profiles for `azureakscluster` and `azureaksnodepool` set
`timeout_minutes: 60–75`. When invoking tests directly, size `-timeout` beyond
the default 30m:

```bash
go test -tags=e2e -timeout=90m -v -count=1 \
  -run 'TestAzureAksCluster_Pulumi/minimal' ./e2e/azure/...
```

### Long-running Azure components (VPN gateways)

Classic virtual network gateways are the suite's slowest single
resource: measured live (VpnGw1AZ, eastus), the create ran **36m** and
the delete **15m** — a full single-engine lane (fixture chain up, deploy,
verify, destroy, verify-gone, chain down) totals **~60m**. Budget
`-timeout=120m` per engine, and budget the SAME cycle inside any lane
whose prerequisite chain deploys the fixture VPN gateway (the gateway
connection's does). Lanes for gateways in the same VNet must run
SEQUENTIALLY — ARM allows one VPN-type gateway per VNet, so the scenario
gateway and the fixture gateway can never coexist.

Burstable VM sizes in `eastus` may support only availability zone `1` — multi-zone
lists fail with `AvailabilityZoneNotSupported`. AKS E2E scenarios in this repo
use `zones: ["1"]` for the test subscription.

### Long-running Azure components (ExpressRoute circuits)

An ExpressRoute circuit CREATE is a slow ARM long-running operation:
measured live at **~17-19 minutes** per create (Equinix / "Washington
DC", 50 Mbps metered, Standard SKU -- consistent across four
consecutive creates on both engines), while the delete completes in a
minute or two. Budget it in every lane whose prerequisite chain deploys
the fixture circuit (the circuit-peering lane pays it inside
DEPENDENCIES-UP), and note the peering's own DELETE runs ~7-8 minutes.
The circuit-family component profiles carry `timeout_minutes: 45` for
exactly this class.

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
Verified on the test subscription: `eastus` is blocked for Microsoft.Sql
logical servers (`ProvisioningDisabled`) while `westus3` and `centralus`
are clean.

Microsoft.Sql adds its own TRANSIENT failure shape on top of the
restriction class: a logical-server create in a CLEAN region can return
`OperationTimedOut` ("The operation timed out and automatically rolled
back. Please retry the operation."). The rollback is asynchronous — for
several minutes afterward the half-created server is invisible to
`az sql server show`/`az resource list` yet still blocks its resource
group's deletion (`ResourceDeletionFailed`), and an explicit
`az sql server delete` answers `DropLogicalServerAlreadyInProgress`. A
`pulumi destroy` of the resource-group stack issued during that window
can hang far past the RG-delete norm. Treat it as transient: kill the
hung destroy if needed, wait for the in-progress drop to finish (the RG
delete then succeeds), and retry the lane — do not debug the module.
A server may live in a different region than its resource group, so the
shared fixture RG serves unchanged when a scenario pins a different
region. True subscription-wide gates do exist (the quota-increase link
in the ARM error is the fix) — but conclude that only after a
probe-verified clean region also fails.

App Service adds its own shape: new PAYG subscriptions carry ZERO
Basic-tier VM quota in some regions, and the rejection is a misleading
`401 Unauthorized` whose body says "Operation cannot be completed
without additional quota ... Current Limit (B1 VMs): 0" — a quota
signal wearing an auth status code, not a credential problem. The
restriction is per-region (verified on the test subscription: `eastus`
blocked, `westus3` clean) and the cheap probe is a one-off
`az appservice plan create --sku B1` in a throwaway resource group. An
App Service PLAN pins the region for every app on it — apps must live
in their plan's region — so re-regioning an app scenario means
re-regioning its plan fixture (the fixture resource group serves
unchanged; a plan may live in a different region than its group).

Cosmos DB adds a THIRD failure shape to this class: transient CAPACITY
rejection rather than offer restriction. `eastus` answered
`ServiceUnavailable` ("high demand in East US region for the zonal
redundant (Availability Zones) accounts") for a plain single-region
account that never asked for zones — Azure routes new accounts through
zonal placement internally, so the error hits non-zonal manifests too.
There is no up-front probe (the provider `locations` list still
advertises the region); treat `ServiceUnavailable` on Cosmos account
creation as a region-capacity signal and move the scenario-local
accounts to a quieter region (`westus3` verified clean) instead of
recording a deferral or filing quota requests.

Cosmos DB SQL data-plane RBAC (`cosmosdb_sql_role_assignment`) adds a
Pulumi-module constraint azurerm does not surface: pulumi-azure enforces
a 24-character logical resource name on `SqlRoleAssignment`, and when
the Azure `name` (a GUID) is unset it autogenerates from that logical
name — producing invalid non-UUID values like `mainaa5aa87`. Forge the
Pulumi module to (1) use a short logical name (`"main"`, mirroring the
Terraform module's single resource) and (2) generate an explicit UUID
for the Azure `name` when the spec does not pin one, matching azurerm's
create-time UUID generation. The sibling `SqlRoleDefinition` needs the
same explicit-GUID treatment for `role_definition_id` when unset; its
logical name can stay on `metadata.name`.

### Front Door: fast creates, ~18-minute profile deletes

Azure Front Door (Standard/Premium) inverts the usual timing profile:
every resource in the family CREATES in seconds-to-minutes (profile
~1-2 min; endpoints, origin groups, origins, and routes each well under
a minute), but PROFILE DELETION runs ~18 minutes wall time (verified
live on both engines; the service tears the global edge deployment down
before ARM confirms). Consequences for scenario budgeting:

- Any scenario whose chain includes a Front Door profile fixture costs
  ~20-27 minutes per engine, dominated entirely by the fixture teardown.
  Child-resource deletes inside a live profile are fast; it is only the
  profile itself that is slow.
- Deleting the fixture RESOURCE GROUP does not dodge this: the RG delete
  waits on the profile delete inside it.
- Estimate a full family suite accordingly (10 dual-engine runs measured
  ~3.7 h total) rather than assuming CDN-family resources are cheap
  because they create quickly.

### Vendor-catalog IDs in fixtures: look them up, never infer them

When a scenario (or preset, or hack manifest) carries an ID from a
provider-curated catalog -- a managed WAF rule ID, a curated rule-set
name/version, any value whose legal vocabulary lives server-side -- look
the real IDs up before writing the fixture; a plausible-looking ID
passes every offline gate (the schema types it as a free string) and
fails only at live apply. The lookup sources, in order: the service's
own catalog API (e.g. Azure's
`GET .../providers/Microsoft.Network/frontDoorWebApplicationFirewallManagedRuleSets`
lists every rule group and rule ID per set version) and the provider
clone's acceptance-test fixtures. The trap that motivated this: Azure
bot-manager rule IDs carry a `Bot` prefix (`Bot300700`) while the
default-rule-set IDs are bare numerics (`942100`) -- an inferred bare
`300700` was rejected by ARM on both engines with "managed rule IDs are
not supported".

### Angle-bracket placeholders in canonical manifests fail the offline plan gate

A kind's `e2e/manifest.yaml` is the canonical validated example — refgen's
Example source AND an offline-plan fixture — so every reference-shaped
value in it must be syntactically REAL, never a `<placeholder>`. Provider
schemas validate resource-identifier syntax at plan time (the AWS provider
runs `ValidARN` on every ARN-typed argument), so `value: "<sns-topic-arn>"`
passes proto validation (a free string) and fails the offline `tofu plan`
with "invalid ARN: arn: invalid prefix". Use realistic well-formed values
(`arn:aws:sns:us-west-2:123456789012:order-events`,
`subnet-0a1b2c3d4e5f60001`) — they document the expected shape for users
and agents, and they keep the manifest a working plan fixture. The same
rule already governs presets; the trap here is that a manifest predating
the offline gate can hide placeholders for years.

### Placeholder domains in receiver/endpoint URLs are server-validated

When a fixture (or preset, or hack manifest) carries a URL the SERVICE
will call back to -- a notification webhook, an alert receiver, an
integration endpoint -- do not use documentation placeholder domains.
Some services validate the URI server-side at create and reject blocked
domains outright: Azure Monitor action groups return 400
`WebhookServiceUriBlocked` for webhook receivers on `example.com`, on
both engines, while every offline gate passes (the schema only checks
the http/https scheme). Use a real domain the fixture plausibly owns
(the project's own domain works; the service does not probe the URL at
create, it only screens the domain). Presets and hack manifests should
carry a domain-you-own shape (e.g. `hooks.yourcompany.com`) so users
never copy a blocked placeholder.

### Vendor-constant casing: an SDK constant's Go identifier is not its wire value

When a module maps a spec enum to a provider's string vocabulary, read
the SDK constant's STRING VALUE, never its Go identifier -- the two can
differ in casing, and provider schemas validate case-sensitively (e.g.
azurerm's `StringInSlice(..., false)`). The trap that motivated this:
the frontdoor SDK declares `TransformTypeURLDecode TransformType =
"UrlDecode"` -- an identifier-derived `"URLDecode"` passes every offline
gate that does not render a plan and fails at plan time on both
engines. Two defenses: (1) verify every enum map row against the SDK
constants file (or the provider schema's validator list), and (2) make
sure at least one scenario or hack manifest exercises every mapping row
whose casing is irregular, so the offline plan gate and the live suite
keep it covered -- a mapping row no fixture uses is dead code that
validation cannot see.

### Service-created default sub-resources cannot be adopted declaratively

Some Azure services auto-create a well-known child resource the moment a
parent exists -- Service Bus subscriptions get a catch-all filter rule
named `$Default`, and similar service-created defaults exist elsewhere.
It is tempting to model "replace the default" as declaring a resource
with the same name, expecting ARM's CreateOrUpdate to upsert it. That
never reaches ARM: both engines' azurerm create paths run an
import-existence check first and fail with "a resource with the ID ...
already exists -- needs to be imported" (live-confirmed on both engines;
the check exists precisely so Terraform-model tools never silently adopt
state they did not create). The honest modeling: reserve the
service-created name in the spec (a CEL rejecting it with a message that
teaches the alternative), document declared siblings as ADDITIVE
alongside the default, and teach the one-time out-of-band removal (an
`az` CLI delete) where restrictive behavior is wanted. No offline gate
catches this class -- the plan renders fine and only the live create
collides -- so treat every "declare over a service-created default"
design as suspect until a live run proves it.

### Provider-accepted values can be SERVER-RETIRED: SKU consolidations reject at create

The pinned provider's static vocabulary lags Azure's retirement schedule,
so a value every offline gate accepts can be rejected by ARM at create
with a retirement error. Live-confirmed classes on the VPN gateway
family: `NonAzSkusNotAllowedForVPNGateway` (new non-AZ VpnGw1-5 creates
blocked since 2025-11-01 -- only the AZ tiers and BASIC remain
creatable) and, chained behind it,
`VmssVpnGatewayPublicIpsMustHaveZonesConfigured` (an AZ-SKU gateway
demands ZONES on its Standard public IP -- a no-zone fixture address
that served the non-AZ world fails the AZ world). When a create rejects
with a retirement-class error: fix the SPEC (retire the values with
reserved numbers/names per the catalog's removal-hygiene precedent),
move scenarios/fixtures/presets to the successor values, and land the
retirement on the kind's user surfaces -- never band-aid just the
scenario. Retirements also arrive in CHAINS: re-run the lane after each
fix expecting the next constraint in the family to surface.

### Some ARM contracts exist ONLY server-side: budget one live probe for the flagship combination

The azurerm provider is the completeness floor, but it is NOT a complete
map of ARM's rejection surface: some cross-field contracts appear in no
schema validator, no CustomizeDiff, and no create-path check -- they live
only in ARM's regional service. Example (live-confirmed): a firewall
attached to a firewall policy must not carry firewall-level DNS
parameters -- ARM rejects the create with
`AzureFirewallDNSConfigNotAllowedForVhubOrVnetWithPolicy` ("DNS
configuration should be managed by policy"), yet the provider happily
plans and sends the combination. The consequence for scenario design:
make the FIRST live scenario for a kind exercise the flagship FIELD
COMBINATION users will actually deploy (here: policy-attached firewall
plus every side-channel knob the policy also owns), because a
source-diff of the provider cannot prove combinations the provider never
validates. When such a rejection surfaces, front-load it as a spec CEL
(with the ARM error code in the message trail), fix the scenario, and
record the contract in the component docs -- the next component in the
same service family should check for sibling "managed by the parent"
exclusions up front.

### "no stack named ..." for a fixture that just deployed: backend state loss, not a module defect

When a scenario fails at DEPENDENCIES-UP with `failed to read outputs for
dependency ...: no stack named '<stack>' found` -- for a fixture whose
`pulumi up` just returned success (or one that "deployed and verified"
moments earlier, then fails teardown the same way) -- the suite's LOCAL
PULUMI FILE BACKEND has lost its stack state mid-run. Nothing is wrong
with the module, the manifest, or the cloud: the resources were created,
but the state tracking them vanished. The backend lives in a per-run temp
directory (`TestMain` creates it with `os.MkdirTemp`), so it is exposed to
anything that disturbs temp storage on the host. Diagnose by the
signature: SEVERAL stacks of the same scenario reporting "no stack named"
at once -- including ones that already deployed and verified cleanly (a
real per-stack failure never spreads to unrelated stacks) -- and an
identical re-run passing end to end.

Recovery: the stackless fixtures' destroys were skipped, so sweep the
cloud for the failed run's fixtures before re-running. A same-named
fixture RESOURCE GROUP in a later scenario can mask the orphans -- ARM's
resource-group PUT is an upsert, so the later run's fixture "creates" the
existing group and its teardown then deletes it, orphans included -- but
never rely on that; sweep explicitly (`az group list`, plus the service's
own list for account-level orphans).

### A second "no stack named" cause: uniquifying suffixes truncated off stack names (fixed in the runner)

The same "no stack named" signature at DEPENDENCY destroy -- for stacks
whose deploys all "deployed and verified" -- once had a different cause
than backend state loss: stack names were blindly truncated to a length
cap AFTER their uniquifying suffixes (an install profile's per-document
index, the run id) were appended, so every document of a long-named
multi-document profile collapsed onto ONE shared stack name. Each
document's `pulumi up` then silently REPLACED the previous document's
resources (each "deployed and verified" against a sibling's grave --
the gateway subnet was gone by the time the gateway deployed), the
first destroy removed the shared stack, and every later destroy failed
"no stack named". `GenerateStackName` now enforces the cap by replacing
the tail with a short hash of the full composed name (deterministic,
collision-free -- see `stackname_test.go`); never reintroduce a blind
`name[:N]` truncation on stack names, and treat "several same-kind
dependencies destroy-failing on ONE shared name" as this class, not
state loss.

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

**ADC preflight must assert a NON-EMPTY token, not just exit code 0:** with a
stale-but-present ADC file, `gcloud auth application-default print-access-token`
can exit 0 while printing an EMPTY token (observed live on a credential file
whose refresh token had aged out) — so `print-access-token >/dev/null && echo OK`
false-passes and the failure surfaces minutes later inside a lane. A machine
preflight should capture the token, assert it is non-empty, and make one live
API call with it (e.g. GET the test project via Cloud Resource Manager). The
harness `Setup` itself fail-fasts correctly — this trap only bites pre-session
checks that trust the exit code.

**IAM service-account read-after-delete is eventually consistent:** a GET
issued within seconds of a successful `serviceAccounts.delete` can still
answer 200 from a stale replica (live-caught on the Terraform SA minimal
scenario: destroy succeeded, VERIFY-CLN 1.5s later still saw the account;
an identical cycle minutes earlier read 404 immediately). The SA verifier
polls for up to 90s and treats 403 as a recently-deleted signal alongside
404 — do not tighten VERIFY-CLN to a single GET for this kind.

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
GKE cluster scenarios across both engines. Dataproc clusters boot real VMs:
budget **≥90m** when batching multiple cluster scenarios across both engines.
Cloud Composer environments are the longest single resource: a create runs
25-45 minutes and a destroy 10-15, per scenario per engine — budget **≥240m**
for a both-engines batch, and the same for any kind whose chain installs an
environment as a prerequisite (the user-workloads Secret/ConfigMap).

**Managed services that VALIDATE runtime identity at create:** some control
planes reject creation unless the workload identity already holds the right
role — Dataproc requires the VM service account to hold `roles/dataproc.worker`
(hardened projects grant the default compute service account nothing), and
Composer 3 requires an explicitly specified workloads service account holding
`roles/composer.worker`. Model the identity as a registry prerequisite with a
consumer-scoped `GcpServiceAccount` profile carrying the additive grant in
`project_iam_roles` — a plain SA prerequisite without the grant fails the
scenario at create with a permissions error, not at verify.

**Test-side data staging (kinds whose deploy needs pre-existing object bytes):**
Some resources cannot apply until a blob already exists in cloud storage — Gen
2 Cloud Functions read their source from GCS at deploy time, and object bytes
cannot be expressed as IaC. The runner has no hook between prerequisite deploy
and the scenario apply, and `valueFrom` resolution only touches
`StringValueOrRef` fields (not plain-string bucket/object paths). The pattern is:
check in a minimal source tree under `e2e/fixtures/`, and in the test
entrypoint (before calling the runner) zip it and upload to a run-scoped bucket
whose name uses the same engine-scoped `${E2E_RUN_ID}` expansion the runner
applies (`runner.EngineScopedRunID`). Register cleanup on the test handle — the
harness `Teardown` is a no-op by design. See `e2e/gcp/gcp_test.go`
`stageCloudFunctionSource` and `gcpcloudfunction/e2e/scenarios/*.yaml`.
Serverless VPC Access connectors are slow: budget **≥15m per connector scenario**
(CREATING → READY is typically 5-10 minutes). Cloud Functions Gen 2 deploys
(including Cloud Build) need **≥20m per scenario**.

**One manifest per prerequisite kind — scenarios needing TWO instances of the
same kind cannot express the second live.** The dependency resolver installs
each prerequisite kind from exactly one manifest (consumer-scoped override,
else the published profile, else the minimal scenario), so a scenario that
composes two instances of one kind — e.g. a Pub/Sub subscription whose parent
topic AND dead-letter topic are distinct topics — can only resolve the first
by reference. Prove the second-instance arm with the offline converter plan
(the reference-resolution mechanism is identical to the live-proven first
instance) and record the exclusion, rather than pointing both references at
one instance (semantically wrong) or hand-rolling a second install path.

**Reservation-window resources need consumer-unique AND scenario-unique
prerequisite names:** `${E2E_RUN_ID}` makes cloud-side names unique per run
— but two scenario chains in the SAME run that install the same
prerequisite profile still produce the same name, and for resource classes
that reserve a deleted ID for minutes after destroy (Firestore database
IDs are held ~3-5 minutes), the second chain's create collides with the
first chain's just-deleted ghost ("not available ... retry in N seconds").
Two grains of the same trap:
- *Multiple consumers* chaining one prerequisite kind: give each consumer a
  consumer-scoped override whose cloud-side name embeds the consumer
  (e.g. `e2e-fsidx-db-...` vs `e2e-fsbs-db-...`).
- *Multiple scenarios of ONE kind*: the kind's own override redeploys per
  scenario, so even a consumer-unique run-scoped name is deleted and
  recreated at every scenario boundary. Embed `${E2E_SCENARIO}` (expanded
  to the running scenario's slug) alongside the run id
  (`e2e-fsidx-db-${E2E_SCENARIO}-${E2E_RUN_ID}`) so no name is ever
  recreated within a run.

**Same-kind dependency stacks and the 50-character bound:** dependency
stack names carry the manifest name precisely so two instances of one kind
get two stacks — and the name bound is enforced by a uniqueness-preserving
truncation (`boundStackName`: readable head + stable hash of the full
name). This is load-bearing: a plain head-truncate once collapsed a
two-database Firestore chain onto ONE stack (the kind slug alone is 20
characters), so the second fixture's `up` silently REPLACED the first and
teardown destroyed the shared stack once, then failed `no stack named` on
the other. A `no stack named` error during DEPENDENCIES-DOWN therefore has
TWO known causes: backend state loss mid-run (see the Azure section) and —
if two same-kind dependencies' stack names agree — a truncation collision
in a modified stack-naming path.

**Identity-derived cloud IDs need run-scoped metadata names:** when a module
derives the cloud-side resource ID deterministically from the resource
identity (org/env/metadata.name — Vertex AI endpoints were first) AND the
service reserves deleted IDs (a destroyed endpoint's numeric ID 409s on
recreate while its GET returns 404), a fixed `metadata.name` makes the
scenario single-shot: the second engine — deriving the byte-identical ID —
collides with the first engine's ghost. Put `${E2E_RUN_ID}` in
`metadata.name` itself, an explicit exception to the fixed-name rule that is
only legal for leaf kinds nothing FK-references (prerequisite resolution
keys off metadata names). Record the exception in the scenario comment.
Corollary: such a cross-engine collision is itself the strongest live proof
that both engines derive the same ID.

**Validate every scenario and prerequisite manifest offline before the first
live run:** manifests load through `protojson.Unmarshal`, which hard-rejects
unknown fields — a mistyped field name (`subnetwork` for `subnet`,
`network` for `vpcSelfLink`) fails at manifest load after minutes of
prerequisite deploys. `planton validate <manifest>` (with `${E2E_RUN_ID}`
text-substituted) catches this in seconds; run it on every new scenario,
prerequisite, and preset YAML as part of the offline gate.

**Offline fixture-integrity gate (FK resolvability):** a `valueFrom` that
names a prerequisite the chain never deploys — or an ambiguous name among
several instances of the same kind — fails only deep into a live run (or
worse, silently: an unresolvable ref is left untouched, and DEPLOY fails as
a provider validation error that reads like a module defect). The gate at
`e2e/framework/runner/fixtureintegrity.go` (`TestCatalogFixtureIntegrity`)
replays `ResolveDependencies` + `forEachRefField` statically against every
committed scenario, with no cloud and no credentials. It mirrors
`lookupRefValue` one for one — including the sole-instance fallback (a
topology-named reference still resolves when exactly one instance of the
kind exists; that is design, not a defect). Pre-existing findings live in
`fixture_integrity_baseline.yaml` and only ever burn down; a new finding is
a manifest fix, never a new baseline entry. Run it with
`go test ./e2e/framework/runner/ -run TestCatalogFixtureIntegrity`.

**Inherited fixture orphans 409 the whole chain — sweep the fixture
families, not just the kinds' own objects:** prerequisite fixtures use
FIXED names by design (FK resolution keys off `metadata.name`), so a
leftover fixture from ANY earlier era — a crashed run, a pre-harness
experiment — collides with a fresh chain as `Error 409: already exists`
at DEPENDENCIES-UP, from a run-scoped backend that has no state for it.
Live-caught with a month-old `planton-oss-e2e-gcphc-prereq` health check
plus its backend service and a serverless NEG: three sequential 409s,
each surfacing only after the previous one was cleared. The fix is
always to DELETE the leftover (a pre-existing fixture object is by
definition unmanaged — no lane is running when a session starts), and
the honest end-of-session sweep must enumerate every fixture FAMILY the
session's chains deployed (for LB chains: health checks, backend
services, NEGs, URL maps, proxies, addresses), not only the proven
kinds' own resource classes.

**Named image families are a staleness trap in scenarios:** GCP retires
image families (the deep-learning-VM notebook families like
`common-cpu-notebooks` no longer resolve), and a scenario pinned to one
fails live years after it was written. Prefer the service's default image
(omit the image block) in E2E scenarios and presets; pin a family only when
the arm under test requires it, and prefer the service's maintained family
(e.g. `cloud-notebooks-managed/workbench-instances`) over frozen ones.

**First-ever API activation outruns in-module enablement:** every module
enables its service APIs with a dependency edge before creating resources,
and that is sufficient for a project where the API has been active before.
But the FIRST-ever activation of a service on a project propagates slowly —
the create call can land minutes before the activation is visible and fail
with a 403 "API has not been used in project ... before or it is disabled",
even though the enablement resource completed. When a new service family
first joins the harness, pre-enable its APIs on the test project once
(`gcloud services enable <api> --project <test-project>`) before the first
live run; from then on the in-module enablement carries every future run.

**Some GCP APIs require a quota project on user-credential calls:** the
Identity Toolkit API (Identity Platform) rejects calls made with plain
user ADC — 403 "requires a quota project, which is not set by default" —
and the Terraform provider sends the `X-Goog-User-Project` header ONLY
when `user_project_override` is armed in the provider config (the ADC
file's own `quota_project_id` is NOT honored by the provider transport;
live-verified). Kinds whose APIs carry this requirement wire
`user_project_override = true` in their TF provider block and build their
Pulumi provider via `pulumigoogleprovider.GetWithUserProjectOverride` —
a module-level, customer-facing fix, never a harness export. The typed
verifier clients are unaffected (google.golang.org/api transports honor
the ADC quota project).

**Once-only initialization singletons need an adopt arm for every lane
after the first:** some project singletons initialize exactly once EVER —
Identity Platform's `initializeAuth` hard-fails a second call with 400
"Identity Platform has already been enabled for this project"
(live-verified), and destroy only abandons. The first live lane spends
the project's single create; every later lane (the same engine's second
scenario, the other engine, every prerequisite chain) must ADOPT instead.
The catalog pattern is a spec-level `adopt_existing` switch (TF: a
for_each-gated `import` block on the deterministic singleton ID; Pulumi:
a conditional `pulumi.Import` resource option) with the kind's e2e
fixtures arming it permanently once the test project is initialized —
Pulumi's import matches because the modules only specify spec-driven
inputs, and Terraform's config-driven import applies spec drift as an
update in the same apply.

**Typed-client gaps in the pinned Google API line:** a brand-new GCP service
may have no typed client in the repo's pinned `google.golang.org/api`
version (Memorystore for Valkey was first). Do not bump the shared
dependency mid-component for one verifier — the harness carries an
ADC-authenticated plain HTTP client (`Services.RestClient`) for exactly
this: probe the service's documented REST GET path and decode the few
fields the posture assertions need. Swap to the typed client whenever the
dependency is next upgraded for its own reasons.

**Run GCP e2e batches SEQUENTIALLY — never two `go test` processes at once:**
the GCP dependency deploys write each prerequisite's stack-input into the
SHARED module source directories under `apis/.../iac/pulumi`, so a second
concurrent `go test` process picks up the first's on-disk manifests and
collides on the other run's names (409 already-exists / cross-contaminated
stacks). Run one component's batch to completion (or one scenario at a time)
before starting the next; parallelism within a single process is bounded by
the framework, but two processes are not isolated from each other.

**Intra-cluster firewall is part of a multi-node cluster's minimum
composition on a custom VPC:** custom-mode VPCs have no default firewall
rules, and managed multi-node clusters (Dataproc among them) refuse to reach
RUNNING because master↔worker traffic is dropped — the create hangs, then
times out, rather than failing fast. Model an intra-cluster `GcpFirewallRule`
(allow ingress within the VPC/subnet) as a registry prerequisite for such
kinds, exactly like the identity and networking prerequisites.

**Dataproc's regional staging/temp buckets persist and carry per-creator
ACLs:** Dataproc auto-creates `dataproc-staging-<region>-<project#>-<hash>`
and `dataproc-temp-...` buckets (a deterministic name per project+region) and
reuses them across runs; it never deletes them. If an earlier run created them
under a since-deleted VM service account, a later run's identity can be locked
out ("Permissions are missing ... on the staging_bucket") even with correct
project-level IAM, because the bucket's fine-grained ACLs still name the old
principal. When a Dataproc live batch fails on bucket permissions after an
identity change, delete the leftover `dataproc-staging-*`/`dataproc-temp-*`
buckets so Dataproc recreates them fresh under the current identity, then
rerun. (The `roles/dataproc.worker` role already carries the full
`storage.buckets.get` + `storage.objects.*` set a custom VM identity needs.)

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

1. Create a YAML manifest in `{component}/e2e/scenarios/` with a descriptive
   filename
2. Use a unique `metadata.name` (and unique namespace if the component creates
   one) to avoid collisions with other scenarios
3. Run `make e2e-test-component component={ComponentName}` to verify it works
4. That's it -- the framework discovers and runs it automatically

## Adding a New Component

1. Create the IaC modules (`iac/pulumi/`, `iac/tf/`)
2. Create `e2e/profile.yaml` with the component's E2E profile
3. Create at least `e2e/scenarios/minimal.yaml` with a minimal test manifest
4. If the component needs other resources installed first, declare them as
   `prerequisites` on the kind in `cloud_resource_kind.proto` (the harness
   installs them automatically -- see "Component Dependencies").
5. Add a `Test{ComponentName}_{Provisioner}` function in the appropriate test
   file (e.g., `kubernetes_test.go`), and -- if the component name does not
   PascalCase trivially -- a `toPascalCase` entry in
   `pkg/e2e/profile/discover.go` so the CI matrix regex matches it
6. The CI workflow picks up the new component automatically from the profile

## Adding a New Provider

1. Create `catalog/{provider}/aa_e2e/` with harness, verify
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
  - *Because of that opt-out, a kind whose ARM namespace has never been used on
    the subscription fails its FIRST live run with a 409
    `MissingSubscriptionRegistration`* (e.g. `Microsoft.App` before the first
    Container Apps run). The fix is the same one-time subscription bootstrap the
    opt-out defers: `az provider register --namespace <ns>` (registration takes
    a minute or two; poll `az provider show -n <ns>` until `Registered`), then
    re-run. Expect this once per new resource-provider family, never per kind.

### Authoring verifiers and prerequisite fixtures

- **Never call `ctx.Export` inside an `ApplyT` callback in a Pulumi module.**
  Apply callbacks run on output-resolution goroutines, and the exports map
  they write is the same map the SDK's end-of-program stack-output
  marshaling reads — a data race that kills the whole program with
  `fatal error: concurrent map read and map write`, timing-dependent and
  therefore FLAKY (first hit: the subnetwork module's per-index
  secondary-range exports crashed one chain deploy and passed the next,
  identical code). A module that needs per-element export keys derives each
  key's index space from the spec (known before the program returns) and
  registers every export synchronously with a value DERIVED via `ApplyT`
  (`ctx.Export(key, out.ApplyT(...))` is safe — it is the export
  registration itself that must stay on the program goroutine). When a
  fixture deploy dies this way mid-chain, the created cloud resource
  usually IS recorded in the run-scoped stack state, so the next scenario's
  `up --refresh` adopts it — but the failed scenario's own teardown can
  strand it against the chain's parent (a VPC refusing deletion while the
  crashed subnet exists) until a later scenario's teardown sweeps both.
- **Never name a verifier (or any production `.go` file) with a `_test.go`
  suffix.** Go treats every file ending in `_test.go` as a test file and
  excludes it from the normal package build, so a verifier in, say,
  `foo_web_test.go` compiles fine under `go test` but is INVISIBLE to
  `go build` — the registration in `verifier.go` then fails with a
  bewildering "undefined: fooVerifier" even though the type is plainly
  defined. Name such files to avoid the suffix (e.g.
  `application_insights_standard_webtest.go`, not `..._web_test.go`).
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
- **Some services access customer resources with the PROVIDER'S OWN service
  principal, which needs its own bootstrap.** E.g. Azure Front Door reads
  customer Key Vaults for bring-your-own TLS certificates as the
  `Microsoft.AzureFrontDoor-Cdn` enterprise application (a well-known,
  tenant-invariant app id) — NOT as the deploying credential. Two one-time
  steps, both idempotent harness-Setup bootstraps in the same class as the
  test-principal grant above: (1) the principal may not exist in a fresh
  tenant until instantiated (`az ad sp create --id <well-known-app-id>`;
  resolve with `az ad sp show` first — create fails when it already exists),
  and (2) it needs the read role granted at the subscription scope
  (`--assignee-principal-type ServicePrincipal`). Without both, the dependent
  resource's create fails with an access-denied error naming the vault, which
  looks like a module defect but is tenant bootstrap.
- **Azure auto-creates `NetworkWatcherRG` the first time a virtual network
  exists in a region** -- it appears mid-run without any manifest creating it
  and survives every teardown (it is subscription furniture, not a test
  orphan). The zero-orphan sweep should recognize it, and deleting it at
  session end is safe and keeps the zero-resource-group baseline exact:
  Azure recreates it on demand the next time a VNet appears.
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
- **A sibling of the state-flip class: some association resources have NO ARM
  object at all — their state lives in a property of the resources they
  associate.** An Azure Managed Redis geo-replication group is the canonical
  case: creating it links existing databases (no new ARM object; the resource ID
  is just the managing cluster's ID) and destroying it unlinks them (nothing is
  deleted, so a 404 probe can NEVER pass verify-absent — and verify-exists
  against the resource ID would pass even when the association never took
  effect). The verifier must GET the carrying resource and read the association
  property (here `properties.geoReplication.linkedDatabases` on the managing
  default database): verify-exists requires the association to be genuinely
  established, verify-absent requires it collapsed. Read the provider's Read
  source to find which resource carries the property.
- **Never let sequential scenarios destroy and recreate the same globally unique
  parent name.** When a kind's registry prerequisite chain deploys a fixture
  whose name is globally unique (an Azure SQL logical server, a Key Vault),
  every scenario of a multi-scenario component tears the fixture down and the
  next scenario recreates it — and Azure can hold the just-deleted name long
  enough that the recreate hangs indefinitely (a `Microsoft.Sql/servers` create
  stuck 20+ minutes with no write in the activity log). Give each scenario its
  own uniquely named parent instead: declare the parent through the
  `e2e-prerequisites` annotation with a scenario-local manifest (kept
  OUTSIDE `e2e/scenarios/`, which the discoverer treats as test cases), and
  drop the registry prerequisite if it would force the shared fixture chain in
  anyway. Registry prerequisites are for parents a kind cannot exist without
  AND that are safe to recreate per scenario. The post-delete name hold is
  SERVICE-SPECIFIC: SQL logical servers hold the name after the delete
  returns, while Azure Cache for Redis frees the name the moment its (slow,
  several-minute) delete completes — same-name recreates across sequential
  engine runs verified live. Scenario-local parents are still the right shape
  for very expensive fixtures regardless of name-hold behavior: a Redis cache
  runs 15-40 minutes per creation, and a scenario that owns its parent never
  serializes against another scenario recreating a shared one.
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

qa/      -- Proto schema for E2E profiles (KRM-style)
  shared/                 -- Shared enums (CostClass)
  providere2eprofile/v1/  -- ProviderE2EProfile KRM API
  componente2eprofile/v1/ -- ComponentE2EProfile KRM API
```

The framework is engine-agnostic. The runner supports both Pulumi and Terraform
execution paths. Each component test runs through the same lifecycle regardless
of which engine is used.
