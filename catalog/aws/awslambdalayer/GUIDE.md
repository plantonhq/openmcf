# AwsLambdaLayer — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## A "change" is always a new version — plan rollout, not update

Every argument is fixed for life at AWS; applying any change publishes a new version with a new ARN. Consumers do NOT follow automatically — each function keeps the version it was configured with until its own next deploy points at the new ARN. Treat a layer change like a library release: publish, then roll consumers at their own pace.

## source_code_hash is your change detector, not AWS's

AWS never reports the hash back — it exists to make content updates declarative. Set it from your build pipeline (`base64(sha256(zip))`): a new hash publishes a new version even when the S3 key stays the same; without it, rewriting the object in place is invisible to the modules.

## skip_destroy trades cleanup for continuity

With `skip_destroy`, replaced versions stay available in AWS (functions pinning them keep working) and bill nothing — but nothing deletes them either. Sweep dormant versions with `aws lambda delete-layer-version` when consumers have moved on.

## The runtime lists are advisory

`compatible_runtimes`/`compatible_architectures` drive console filtering and warnings only — the API attaches any layer to any function. Do not rely on them as a compatibility wall; they are documentation for humans.

## Multi-statement grants import lossily upstream

The provider's permission importer reads only the first policy statement, so a version carrying several grants round-trips lossily on import (each imported grant reflects statement one). Keep grants few and coarse (one org grant beats many account grants) when import fidelity matters.
