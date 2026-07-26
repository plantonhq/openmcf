---
title: "Artifact store preset"
description: "A bucket-centric store for build and release artifacts: three typed buckets created at install — CI artifacts that expire themselves after 30 days (SeaweedFS TTL, no cleanup job), release artifacts..."
type: "preset"
rank: "03"
presetSlug: "03-artifact-store"
componentSlug: "seaweedfs"
componentTitle: "SeaweedFS"
provider: "kubernetes"
icon: "package"
order: 3
---

# Artifact store preset

A bucket-centric store for build and release artifacts: three typed
buckets created at install — CI artifacts that expire themselves
after 30 days (SeaweedFS TTL, no cleanup job), release artifacts with
S3 versioning so an overwrite never destroys a prior release, and a
public-assets bucket with anonymous reads for content served behind
an ingress — on a moderately sized single-node topology with the S3
gateway embedded on the filer.

The bucket policies are the point, so know their edges: TTL expiry is
irreversible for the objects it removes — do not put anything on
`ci-artifacts` that a release depends on. `anonymous_read` opens
READS only (writes still require credentials), and it applies to
everything in that bucket — public means public. And removing a
bucket entry from the spec never deletes the bucket or its data;
buckets are created by IaC but only ever destroyed by a deliberate
data operation. Object Lock (`object_lock: true`) exists for WORM
compliance needs but is deliberately absent here — it is irreversible
and forces versioning.

Change first: the bucket names and the TTL window to match your CI
retention policy, then `volume.data_volume.size` for the artifact
volume you expect — versioned buckets in particular accumulate every
revision. When artifact traffic grows past a single node, carry these
buckets onto the production HA preset's topology unchanged.

See [03-artifact-store.yaml](./03-artifact-store.yaml) for the
manifest.
