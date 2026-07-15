// Package envpartition assigns discovered cloud account resources to
// environments, deterministically. It is the engine behind the environment
// partition stage of account import: before any grouping or mapping, every
// discovered resource is placed in an environment (or honestly left
// unassigned), because environment is part of a resource's identity and a
// wrong split poisons every env-scoped reference and the resource graph
// itself.
//
// The engine applies an EnvironmentPartitionRule (iac.planton.dev/v1) with
// a strict precedence ladder -- deterministic user intent always outranks
// inference:
//
//	taught override > authoritative tag > name token > containment > unassigned
//
// Two properties are load-bearing and must survive any change:
//
//   - DETERMINISM. The same rule and the same resources always produce the
//     same result, byte for byte. A re-scan replays a taught rule and must
//     reproduce the confirmed split exactly; review gates diff against
//     prior runs. Nothing here may iterate a map without imposing order.
//   - NO GUESSING. A resource whose signals are absent or contradictory
//     stays unassigned (with the contradiction flagged); the review gate
//     owns the remainder. False confidence is the failure mode this
//     package exists to prevent.
//
// The engine is provider-agnostic: it consumes neutral Resource values.
// The awsscan subpackage adapts AWS Cloud Control scan records (type name
// plus property document) into them.
package envpartition
