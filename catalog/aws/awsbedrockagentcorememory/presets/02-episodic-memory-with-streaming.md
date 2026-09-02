# Episodic Memory with Streaming

This preset captures experience episodes with EPISODIC reflection —
"what happened and what worked" — indexed by customer and streamed to
Kinesis as records are written, for agents that should learn from past
runs.

## When to Use

- Task agents whose past attempts should inform future ones
- Teams building analytics or replication on memory records

## What You Get

- Episodes under `/episodes/{actorId}` plus reflection records
  consolidated under `/episodes` (a reflection namespace must equal, or
  be a whole-segment prefix of, an episode namespace — AWS rejects a
  disjoint pair)
- A `customer_id` index for filtered retrieval
- Every record's metadata on your Kinesis stream in near-real-time

## Customize

- Set `contentLevel: FULL_CONTENT` to stream record bodies, not just
  metadata — mind the payload sizes
- The execution role must be able to write the stream — the reference
  wires the role; the grant is yours
