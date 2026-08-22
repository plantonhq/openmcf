# JSON Contract Schema

This preset registers a strict JSON Schema for a topic's message values -- closed to unknown properties, with an enumerated status field -- turning a loose JSON pipeline into a validated contract.

## When to Use

- Teams producing plain JSON (no Avro toolchain) that still want registry-enforced structure
- Integration topics where unknown fields should be rejected, not silently passed through

## Key Configuration Choices

- **`additionalProperties: false`** -- the schema is a real contract: producers cannot smuggle new fields past consumers.
- **`enum` on status** -- value-level validation, not just shape.
- **`subjectName: shipment-events-value`** -- the `<topic>-value` convention registry-aware serializers resolve.

## What You Get

A registered JSON Schema subject on the cluster's registry. Evolving a closed schema is a breaking change by design -- and on this resource any change replaces the subject and drops prior versions, so plan contract changes with consumers first.
