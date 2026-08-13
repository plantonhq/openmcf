# Azure Data Factory Data Flow -- Operational Guide

Judgment calls that matter when you run Data Factory data flows in production.

## Author in the Studio, ship through the catalog

The data flow script is not a language to hand-write: build the flow visually in the Data Factory Studio (a dev factory is ideal), open the "Script" view, and copy it into the manifest. The catalog deliberately does not re-model the script language -- Azure owns it and evolves it on its own schedule -- so the Studio is the authoring surface and the manifest is the shipping surface. The sources/sinks/transformations blocks must name the same streams the script names; a mismatch fails at deploy time.

## Deploy-time green is not run-time green

A data flow deploys the moment its script parses and its stream names line up -- but the datasets and linked services its endpoints bind to are validated when a pipeline actually RUNS the flow. After deploying, run the owning pipeline's Debug execution once before trusting the flow in schedules.

## Flowlets are your dedup lever -- name them like an API

A flowlet's name is its contract: every embedding flow references it by name, and renaming it breaks them all at their next deploy (the reference is resolved at save time, not tracked by ARM). Treat flowlet names like package names -- stable, versioned by suffix when logic changes shape (`scrub-pii-v2`), never recycled for different logic.

## Spark spins up per run -- batch accordingly

Every pipeline run that executes a data flow pays a Spark cluster spin-up (minutes, billed per vCore-hour from the factory's integration runtime settings). Ten small flows run separately cost ten spin-ups; the same logic as one flow with ten transformations costs one. Consolidate where the logic allows, and use the pipeline's Execute Data Flow activity settings (compute size, TTL) to tune the runtime -- those dials live on the ACTIVITY, not on this resource.

## Rejected-row routing beats failing the run

For flows ingesting third-party data, wire each sink's `rejected_linked_service` to a quarantine store: schema-violating rows divert instead of failing the whole window's run. The alternative -- a run that dies at 3 a.m. on one bad row -- costs an on-call page and a backfill; the quarantine costs a review queue.
