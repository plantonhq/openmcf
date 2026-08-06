# RAG workload preset

A single-node Qdrant sized for a typical RAG corpus: 8Gi of memory
(the real capacity bound — Qdrant serves from RAM-resident segments
and HNSW indexes; 8Gi holds an embedding set in the
several-million-vector range at common dimensions), a generated API
key because ingestion and retrieval cross service boundaries even in
a single-team setup, a snapshots volume so the index restores in
minutes instead of re-embedding for hours, and one engine-tuning
entry through `helm_values` raising the request-size ceiling for
bulk-ingestion batches.

One node is a deliberate choice here, not an oversight: most RAG
stores are derived data — rebuildable from source documents — so
quorum availability is usually not worth tripling the memory bill.
The snapshot volume is the cheaper insurance. If retrieval downtime
is genuinely unacceptable, take the production-cluster preset
instead; growing THIS preset later is only a `replicas` change (the
new pods join the existing consensus), plus per-collection
replication factors set through the Qdrant API.

Change first: the memory sizing, from your actual corpus — vectors ×
dimensions × 4 bytes, plus index overhead, is the back-of-envelope;
quantization (declared per collection through the API) cuts it
severalfold. Then `storage.size` and `snapshots.size` together (keep
them equal, per upstream guidance). The `helm_values` block is the
template for further engine tuning — collection defaults and
optimizer settings ride the same `config:` document — and must never
carry secrets; key material belongs in the `api_key` arms.

See [03-rag-workload.yaml](./03-rag-workload.yaml) for the manifest.
