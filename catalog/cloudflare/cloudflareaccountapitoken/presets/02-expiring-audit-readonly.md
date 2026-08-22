# Expiring read-only audit token

The time-boxed shape: read-only access for an external auditor, alive for exactly one quarter and then dead on its own. The validity window is the control here -- `not_before` and `expires_on` mean nobody has to remember to revoke it, and the spec checks that the window runs forwards. This uses the whole-account grant form (`permission: "*"`), which is appropriate for read-only groups and deliberately broader than the zone-scoped form; keep it paired with read permission groups only.
