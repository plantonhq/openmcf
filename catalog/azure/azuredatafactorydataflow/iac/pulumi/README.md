# AzureDataFactoryDataFlow Pulumi Module

## Overview

Creates one data flow inside an Azure Data Factory -- either a mapping data flow (a complete, runnable transformation) or a flowlet (a reusable snippet other data flows embed). The spec's `flowlet` flag selects the form; the module creates the matching classic-SDK resource.

## Resources Created

Exactly one of:

- `datafactory.DataFlow` -- the mapping form (`flowlet: false`, the default); requires at least one source and one sink
- `datafactory.FlowletDataFlow` -- the flowlet form (`flowlet: true`); sources and sinks are optional (the embedding flow supplies them)

## Outputs

- `data_flow_id` -- the data flow's ARM resource ID (`{factory_id}/dataflows/{name}`, the same shape for both forms)
- `data_flow_name` -- the data flow's name (what flowlet references resolve against)

## Behavior Notes

- **The two forms share one name namespace** inside the factory, so flipping `flowlet` replaces the object (a different resource is created at the same ARM address).
- **The script is the transformation** -- `script`/`script_lines` carry the data flow script language Azure owns; the source/sink/transformation blocks only declare the named endpoints the script references.
- **The SDK generates a parallel nested-type set per resource** (`DataFlow*` / `FlowletDataFlow*`), so the module carries twin block builders in lockstep -- both write the exact same ARM wire shape.
- **Rejected-row routing exists on sinks only** -- Azure's model; the provider silently drops it on sources, so the spec does not offer it there.
- **No tags**: data flows are ARM sub-resources of the factory and expose none.
