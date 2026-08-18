# Shared Canary Groups

The groups-only shape: one instance owns the fleet's console groups, and every canary instance joins them by name (`groupNames: [prod-critical-flows]`). Groups aggregate run results across canaries in the CloudWatch console — the fleet view on one screen.
