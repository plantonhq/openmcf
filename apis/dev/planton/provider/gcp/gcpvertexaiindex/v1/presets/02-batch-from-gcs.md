# Batch Index from Cloud Storage

A bulk-built index seeded from a Cloud Storage directory of embedding
files — the economical shape for corpora refreshed on a schedule.

## What this preset creates

A `BATCH_UPDATE` index named `Catalog Embeddings` in `us-central1`, built
from the files under `gs://ml-embeddings/catalog/` with explicit tree-AH
tuning and medium (20 GB) shards. The initial build is a long-running
operation — expect tens of minutes to hours depending on corpus size.

## When to use

- Product catalogs, knowledge bases, or archives refreshed nightly/weekly
- Large corpora where streaming's per-GB premium isn't justified
- Pipelines that already materialize embeddings to Cloud Storage

## Remix ideas

- Set `isCompleteOverwrite: true` on a subsequent apply to replace the
  whole corpus instead of applying the files as a delta.
- Switch `distanceMeasureType` to `COSINE_DISTANCE` if your vectors are
  not unit-normalized and the model was trained on cosine similarity.
- Drop `contentsDeltaUri` to create the index empty and load data later —
  the field can be set in a follow-up apply.
