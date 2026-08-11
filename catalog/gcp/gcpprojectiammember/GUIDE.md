# GcpProjectIamMember Guide

The judgment this guide protects: this kind is the SAFE way to grant —
one (role, member) pair, additive, composing with grants made anywhere
else. Its value is precisely what it refuses to do: it never rewrites
the project's policy.

## Additive means removal is exact, too

Destroying this resource subtracts exactly its own (role, member,
condition) tuple — grants made by other tools, charts, or humans stay
untouched. That is why fleets of these compose safely where the
authoritative `iam_binding`/`iam_policy` forms would clobber each
other. One resource per grant; resist the urge to "batch" by reaching
for authoritative forms.

## Everything is ForceNew — a grant edit is revoke-then-grant

The provider has no update path: changing role, member, or condition
replaces the resource. The window between revoke and re-grant is real
(seconds) — for a live workload's only grant, add the new grant first
as its own resource, then remove the old one.

## Conditions are part of the grant's identity

An IAM condition (title + CEL expression) makes this a DIFFERENT grant
from its unconditioned twin — both can coexist, and the unconditioned
one wins in effect. When narrowing an existing grant with a condition,
remember to remove the unconditioned original or the condition is
decorative.

## The ambient-project fallback is a convenience, not a pattern

An empty `projectId` resolves to the provider's default project at plan
time (the resource requires an explicit project argument, so the module
reads the provider's own configuration). Fine for single-project
setups; in anything multi-project, name the project — a grant landing
in whichever project the credentials default to is how permissions end
up somewhere surprising.
