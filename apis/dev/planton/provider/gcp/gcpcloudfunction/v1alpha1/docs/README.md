# GcpCloudFunction — Research and Design Documentation

## 1. Introduction

### What Is a Cloud Functions (Gen 2) Function?

Gen 2 Cloud Functions is not a standalone FaaS platform — it is a source-based abstraction over Cloud Run and Eventarc. You ship a source archive; Cloud Build containerizes it with buildpacks; Cloud Run serves it; Eventarc routes events to it. Every deployed function is backed by a real Cloud Run service (exported as `cloud_run_service_id`), which is why Gen 2 inherits Cloud Run's serving characteristics: up to 1000 concurrent requests per instance, 60-minute HTTP timeouts, up to 32 GB / 8 vCPU shapes.

Gen 1 is legacy: lower limits, a separate serving stack, and no reason to choose it for new work. This kind models Gen 2 only (`google_cloudfunctions2_function`).

### The Composition Boundary

- **GcpCloudFunction** owns the function: build shape, serving shape, trigger.
- **GcpGcsBucket** holds the source archive (`buildConfig.source.storageSource.bucket` is a reference).
- **GcpServerlessVpcConnector** bridges egress into a VPC (`serviceConfig.vpcConnector` is a reference) — the only VPC path for Cloud Functions.
- **GcpPubSubTopic** feeds Pub/Sub triggers (`eventTrigger.pubsubTopic` is a reference).
- **GcpServiceAccount** provides three distinct identities: the build SA (fully-qualified name), the runtime SA (bare email), and the Eventarc invoker SA (bare email) — all reference-shaped.
- **GcpKmsKey** provides CMEK for the built image and source artifacts.
- **GcpRegionNetworkEndpointGroup** bridges the function into the external Application Load Balancer family (the NEG's `cloudFunction.function` consumes this kind's bare `name` output).

## 2. The Two-Config Split (mirrored from the API)

- **build_config** — HOW source becomes a container: runtime (a free string, because GCP adds runtimes faster than any allowlist can track — deprecated ones are rejected at deploy time), entry point, source (GCS archive XOR Cloud Source Repositories revision), build env, build identity, worker pool, docker repository, and the base-image update policy (AUTOMATIC continuous patching vs ON_DEPLOY pinning).
- **service_config** — HOW it runs: quantity-string memory ("256M"/"1Gi") and explicit CPU, timeout, per-instance concurrency (GCP defaults to 1), plain env, Secret Manager references (env vars and volume-projected files — the material never enters the spec), VPC connector egress, ingress, scaling bounds, traffic pinning (`all_traffic_on_latest_revision: false` is the manual canary/rollback lever), and Binary Authorization.

Secret references are the API's only secret mechanism — Gen 2 accepts no literal secret material, which is why the spec models `{key, secret, version, project_id}` reference messages (with `sensitive_exempt_reason` on the secret NAME) rather than a value map.

## 3. Trigger Model

HTTP (the default) or a CloudEvent through Eventarc: an `event_type`, optional attribute filters (`match-path-pattern` for Firestore-style paths), a Pub/Sub topic for messagePublished, a trigger region (multi-region Storage sources use "us"/"eu"), a retry policy (RETRY = at-least-once, handlers must be idempotent), and the invoking service account (needs `run.invoker`).

"Allow unauthenticated" is not a Cloud Functions IAM concept in Gen 2 — it is `run.invoker` for `allUsers` on the UNDERLYING Cloud Run service, which is exactly what both modules create.

## 4. Terraform Provider Floor

Designed from `google_cloudfunctions2_function` on the released Terraform Google provider 6.x line (`~> 6.0`, schema-probed at 6.50.0). Both engines enable the five APIs a Gen 2 deploy exercises (cloudfunctions, cloudbuild, run, artifactregistry, eventarc).

### Deliberately Not Modeled (verified against the released schema)

- **Direct VPC egress** (`direct_vpc_network_interface` / `direct_vpc_egress`) — absent from the released 6.x line for this resource (present only on the unreleased provider main and the bridged Pulumi SDK; modeling it would create a one-engine field). Cloud Run and Cloud Run jobs have it; for functions, the connector is the modeled VPC path. Revisit when the provider line carries it.
- **deletion_policy** — absent from the released 6.x line for this resource (same one-engine class).
- **Function IAM binding/policy resources** — authoritative IAM clobbers composed grants; the additive public-invoker member (opt-in) is the only IAM write this kind performs.

## 5. Registry

- **Enum:** 602
- **ID prefix:** `cldfunc`
- **Prerequisites:** `GcpServerlessVpcConnector` (the commonly-composed private-egress reference)
