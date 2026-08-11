# GcpVertexAiIndex Guide

The judgment this guide protects: an index's geometry is chosen ONCE.
`dimensions`, the algorithm arm, `shardSize`, `distanceMeasureType`,
`featureNormType`, and `indexUpdateMethod` are all immutable — every one
of them is a replace-the-index event, and replacing an index means
rebuilding the corpus (tens of minutes to hours for batch, or replaying
months of streaming upserts).

## Pick the update regime for the pipeline you actually run

`BATCH_UPDATE` (the default) suits corpora rebuilt from a pipeline:
cheaper per GB, loaded from Cloud Storage, rebuilt wholesale.
`STREAM_UPDATE` suits corpora that mutate continuously: upserts land in
seconds, at a higher per-GB price, and still compact in the background.
Migrating between regimes is a rebuild — if any part of the roadmap
says "real-time updates later", start streaming now.

## Geometry follows the embedding model, not preference

`dimensions` IS the encoder's output size. `distanceMeasureType` must
match how the model was trained — a mismatched measure returns
plausible-looking neighbors that are silently worse, the kind of
quality bug nobody attributes to infrastructure. `DOT_PRODUCT_DISTANCE`
with `UNIT_L2_NORM` ranks identically to cosine similarity; that pair
is the safe default for most text encoders.

## Tree-AH tuning is a recall/latency dial, not a checkbox

`approximateNeighborsCount` (required with tree-AH) is how many
candidates approximate search fetches before exact reordering — raise
it for recall, pay in latency. `leafNodesToSearchPercent` is the same
trade at the tree level. Brute force exists for two honest uses: small
corpora, and producing ground truth to measure a tree-AH index's recall
against. Do not serve production traffic from brute force.

## Data loading has two provider quirks worth expecting

`contentsDeltaUri` is write-only upstream: GCP never reports it back,
so an out-of-band load shows as a one-field diff on the next plan —
expected, not drift. A change to it also travels in its OWN single-field
update; no other field can ride along in the same apply.

## CMEK is a create-time decision with an IAM prerequisite

`kmsKeyName` (a `GcpKmsKey` reference) must name a key in the index's
region, and the Vertex AI service agent needs
`roles/cloudkms.cryptoKeyEncrypterDecrypter` on it BEFORE the create —
a missing grant fails the long-running create late, not fast.
Immutable: encrypting an existing index means rebuilding it.

## deletionPolicy protects the hours you spent building

Empty/`DELETE` deletes the corpus with the resource. `PREVENT` makes
destroy fail — set it on any index whose build time is measured in
hours or whose streaming history is not reproducible. `ABANDON` hands
the running index to another management plane (it keeps billing for
stored vectors).

## What is deliberately absent

Serving configuration lives elsewhere by design: the endpoint
(`GcpVertexAiIndexEndpoint`) owns connectivity and the deployment
(`GcpVertexAiDeployedIndex`) owns compute. An index alone costs storage
and serves nothing.
