# GCP Backend Buckets: Static Content as a First-Class Load-Balancer Node

## Where the Backend Bucket Sits

Google's external HTTP(S) load balancer routes every request through a URL map to exactly one of two backend types: a backend *service* (dynamic content — instance groups, NEGs, serverless) or a backend *bucket* (static content — a Cloud Storage bucket). The backend bucket is the piece that turns "a bucket of files" into "an origin behind a global anycast VIP, with edge caching, compression, custom headers, and edge security policy."

The split between the bucket and the backend bucket is deliberate and worth preserving in composition:

- **One bucket, many serving policies.** The same GCS bucket can sit behind an aggressive-caching backend bucket for `/assets/*` and a no-caching one for `/downloads/*`.
- **Origin swaps are cheap.** `bucket_name` is mutable — pointing a backend bucket at a new bucket (a blue/green static release) is an in-place update that leaves the URL map and everything above it untouched.
- **The reference direction matches GCP's.** URL maps reference backend buckets by self-link; backend buckets reference buckets by name. Each hop is its own composable node.

## CDN Is a Policy, Not a Resource

There is no standalone "Cloud CDN" object in GCP. CDN is `enable_cdn` plus a `cdn_policy` block ON the backend (bucket or service) — flipping the flag on an existing backend bucket is the entire "add CDN" operation. Modeling honesty follows the API: caching behavior lives inside this kind, not in a separate CDN node.

The `cdn_policy` knobs that matter most in practice:

- **`cache_mode`** decides who controls cache lifetimes. `CACHE_ALL_STATIC` (the GCP default) caches recognized static content types and honors origin headers for everything else — right for most origins. `USE_ORIGIN_HEADERS` gives the origin total control (and therefore forbids the TTL fields — the spec enforces what GCP would otherwise silently strip). `FORCE_CACHE_ALL` caches *everything* including responses that say not to — never point it at private or per-user content.
- **`negative_caching`** absorbs error storms: a missing asset that would otherwise hammer the origin with 404s gets cached briefly at the edge (codes are limited by GCP to 300, 301, 308, 404, 405, 410, 421, 451, 501; TTL ≤ 30 minutes).
- **`serve_while_stale`** smooths origin blips by serving expired content while revalidating in the background.
- **`cache_key_policy`** tunes hit rate: fingerprinted immutable assets want the narrowest key (URL only — the default ignores query strings on backend buckets); an image-resizing origin wants exactly its size parameters whitelisted.

Cache *invalidation* is deliberately not modeled: it is an operational verb (`gcloud compute url-maps invalidate-cdn-cache`), not a declarative property, and fingerprinted asset paths sidestep it entirely.

## Signed URLs: the Private-Content Path

Cloud CDN serves private content through signed URLs and signed cookies — expiring, tamper-proof links minted with a named 128-bit key. GCP models each key as a separate API resource (`backend_bucket_signed_url_key`), but this component folds them in as a `signed_url_keys` list rather than a separate kind, on the granularity rule: keys are never referenced by any other resource, GCP caps them at 3 per backend bucket, and their lifecycle is the bucket's. A separate kind would be a glue node.

The folded shape still honors the key facts:

- **`key_value` is a secret** — annotated sensitive in the spec, marked secret in Pulumi state, never exported in outputs. Anyone holding it can mint valid URLs for the content.
- **Keys are immutable in GCP** (add/delete only). Rotation is: add a new key (the 3-slot cap always leaves room), re-sign URLs with it, remove the old one. Changing a key's value in the manifest replaces that key — exactly the rotation semantics.
- **The keys' names are public** (they appear in the signed URL's `KeyName` parameter); only the value is secret.

## The Edge Security Hook

`edge_security_policy` attaches a Cloud Armor policy of type `CLOUD_ARMOR_EDGE` — rate limiting and geo/IP filtering applied *before* the cache, protecting both origin and cached content. It attaches by reference (`StringValueOrRef` → GcpCloudArmorPolicy), never by embedding: the policy is its own first-class node with its own lifecycle. Standard backend-type Cloud Armor policies do not attach here — they belong on backend services.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `bucket_name` | ✅ `bucketName` | `StringValueOrRef` → GcpGcsBucket; mutable (origin swaps) |
| `name` | ✅ `backendBucketName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; falls back to the provider default project |
| `description` | ✅ | |
| `enable_cdn` | ✅ `enableCdn` | With the INTERNAL_MANAGED incompatibility enforced pre-deploy |
| `cdn_policy` (full block) | ✅ `cdnPolicy` | All 11 sub-fields incl. cache-mode/TTL coherence CEL |
| `compression_mode` | ✅ | AUTOMATIC / DISABLED |
| `custom_response_headers` | ✅ | Header-form validated |
| `edge_security_policy` | ✅ | `StringValueOrRef` → GcpCloudArmorPolicy |
| `load_balancing_scheme` | ✅ | INTERNAL_MANAGED (cross-region internal ALB) or unset (external) |
| signed-URL keys (separate resource) | ✅ folded as `signedUrlKeys` | Max 3; `key_value` sensitive; never FK-referenced → folded, not a kind |
| `deletion_policy` | ❌ | A Terraform-provider-level abandon-vs-delete lever, not a property of the backend bucket; Planton's lifecycle management owns this concern |
| `params.resource_manager_tags` | ❌ | Write-only tag bindings, unmodeled across the catalog; adopting them is a catalog-wide decision |
| backend-bucket IAM (policy/binding/member) | ❌ | Beta-only resources; resource-scoped IAM trios are deliberately not modeled as kinds |
| `timeouts` | ❌ | Operation plumbing, not resource configuration |

## Composition

The static half of a global serving path:

1. **GcpGcsBucket** — the origin holding the objects.
2. **GcpBackendBucket** (this component) — references the bucket, owns CDN/compression/edge policy.
3. **GcpUrlMap → GcpTargetHttpsProxy → GcpGlobalForwardingRule** — route `/assets/*` here by self-link; dynamic paths go to backend services.

Pair with **GcpCloudArmorPolicy** (type `CLOUD_ARMOR_EDGE`) for edge filtering of the static path.
