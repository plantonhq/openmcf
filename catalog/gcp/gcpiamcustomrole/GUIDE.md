# GcpIamCustomRole Guide

The judgment this guide protects: a custom role is a security contract
with every binding that references it. Editing the role edits every
grant simultaneously — that leverage is the point, and the risk.

## Permissions edits propagate instantly, both directions

Adding a permission silently widens every existing grant of the role;
removing one can break a workload that was quietly depending on it.
Treat the `permissions` list like a public API: additive changes are
cheap, removals deserve a search for who holds the role first. Build
roles from the narrowest permission set that works — start from what
the workload's audit logs show it uses, not from a predefined role's
superset.

## DISABLED is the kill switch

`stage: DISABLED` keeps the role defined while every binding of it stops
granting — the fastest way to neutralize a role you suspect is
over-broad without destroying the definition or the bindings. The other
stage values are labels, not behavior.

## Deletion is soft, and names have memory

Destroying a custom role soft-deletes it: bindings stop granting
immediately, the role is recoverable by undelete for 7 days, and the
role ID stays unusable until the definition is purged (~37 days). A
"delete and recreate with the same ID" plan therefore does not work on
the timeline you expect — prefer editing in place or new IDs.

## Destroy semantics

`deletionPolicy: DELETE` (default) soft-deletes as above. `PREVENT`
suits roles that production bindings depend on — the destroy fails
before any grant stops working. `ABANDON` leaves the role active and
every binding working, unmanaged; because custom-role deletion is soft
anyway, ABANDON here is genuinely low-drama — the definition simply
outlives its management.
