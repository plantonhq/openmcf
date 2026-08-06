# Brute-Force Evaluation Index

An exact-search index over a sample corpus — the ground truth for
measuring an approximate index's recall.

## What this preset creates

A `BATCH_UPDATE` (default) index named `Eval Ground Truth` in
`us-central1` using brute-force (exhaustive) search over the embeddings
under `gs://ml-embeddings/eval-sample/`. Every query scans every vector:
perfect recall, linear cost.

## When to use

- Measuring recall of a tree-AH index: run the same queries against both
  and compare the neighbor sets
- Small corpora (thousands of vectors) where exact search is affordable
- Correctness-critical lookups where a missed neighbor is unacceptable

## Remix ideas

- Keep the sample small — brute-force cost grows linearly with corpus
  size, and a few thousand vectors usually suffice for recall estimates.
- Match `distanceMeasureType` and `featureNormType` exactly to the
  approximate index under evaluation; a mismatch invalidates the
  comparison.
