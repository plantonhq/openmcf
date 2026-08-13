# Azure Data Factory Dataset -- Operational Guide

Judgment for operating datasets in production, earned from how Azure actually behaves.

## Match the dataset shape to the linked service type

Azure saves any pairing without complaint -- a delimited text dataset over a MySQL connection stores fine -- and fails only when a pipeline first uses it. The pairing rule is yours to keep: file formats (CSV, JSON, Parquet, binary, blob) ride storage or web connections; table shapes ride their database's own connection type; the HTTP shape rides a web connection. Studio's Preview data button is the cheapest verification that the pairing and the path are both right.

## Saving is not reading

Creating a dataset never touches the data: a wrong container, a typo'd table name, or a missing file all save successfully and surface as run-time copy failures. Preview data after deploy, and treat the first pipeline run against a new dataset as part of its rollout.

## One name, one dataset -- switching shapes replaces it

All 13 shapes live in one factory-scoped namespace ({factory_id}/datasets/{name}). Changing a dataset's variant block replaces the object at the same ARM address, and every pipeline activity referencing the name picks up the new shape immediately -- there is no versioning. Rename rather than reshape when pipelines still depend on the old contract.

## Parameters make one dataset serve many runs

The parameters map plus the dynamic_*_enabled flags turn literal paths into run-time expressions (`@{dataset().runDate}`), so one dataset serves every partition of a feed instead of one dataset per day. Pipelines override parameter values per activity. Prefer this over generating dataset-per-partition manifests -- the factory stays navigable.

## Declared columns are a contract, not a requirement

Data Factory infers columns happily; declare schema_column only when downstream mappings need a stable contract, and expect to maintain it when the source evolves -- a stale declared schema fails runs the inference form would have survived. Snowflake declares in its own type vocabulary (with precision/scale); everything else uses Data Factory's interim types.

## The custom form trades validation for reach

The custom variant carries any dataset type Data Factory speaks as raw JSON -- Excel sheets, XML, Avro, ORC, REST payloads. Azure validates the JSON at save time against the declared type, but nothing validates it BEFORE deploy: keep custom type-properties JSON small, copy shapes from the Data Factory REST API reference exactly, and prefer a first-class variant whenever one exists.
