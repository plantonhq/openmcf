# Tiered Retention With DR

Two schedules, one policy: 24 rolling hourlies for oops-recovery, plus a Sunday weekly kept 90 days whose copies replicate encrypted to us-west-2 for 30 days of regional DR. The cost graph is explicit — every retention rule is its own meter.
