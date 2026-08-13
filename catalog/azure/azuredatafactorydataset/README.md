# Overview

The **AzureDataFactoryDataset** component deploys one dataset inside an Azure Data Factory (AzureDataFactory) -- a named view of data that tells pipelines and data flows what the data looks like and where it sits within a system a linked service (AzureDataFactoryLinkedService) already connects to: which container and path, which table, which file format. One kind covers all 13 dataset shapes azurerm models as first-class resources, plus the raw-JSON custom form for every other Data Factory dataset type.

## Purpose

- **Data shapes as configuration**: the contract between pipelines and their data -- formats, locations, declared columns -- lives in the manifest, reviewed and versioned like everything else, not clicked together in a portal.
- **One kind, 13 shapes**: the variant block declares the dataset type; Azure stores every shape in one factory-scoped namespace, and so does the catalog.
- **Metadata only, secret-free by design**: a dataset names WHERE data lives; the linked service it references carries HOW to authenticate. No dataset field holds secret material.

## Key Features

- Full azurerm v5 surface across ALL THIRTEEN provider dataset resources: file formats (delimited text/CSV, JSON, Parquet, binary, the flat blob form), an HTTP file form, tables (Azure SQL, SQL Server, MySQL, PostgreSQL, Snowflake), a Cosmos DB (SQL API) collection, and the raw-JSON custom escape hatch (Excel, XML, Avro, ORC, REST, SAP, and dozens more).
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output, and every linked service reference -- the shared `linked_service_name`, the SQL table's `linked_service_id`, and the custom form's block -- defaults to AzureDataFactoryLinkedService's outputs.
- Exact contracts: per-variant validation mirrors the provider's own rules -- exactly one location block per file format, each format's own path/filename requiredness, the interim-type and Snowflake column vocabularies, and the compression codecs each format accepts.

## Use Cases

- **Lakehouse landing zones**: CSV/JSON datasets over blob storage or Data Lake Gen2 locations describe raw drops; Parquet datasets describe the curated side.
- **Warehouse loads**: table datasets (Azure SQL, Synapse via custom, Snowflake) give copy activities their source and sink shapes.
- **HTTP ingestion**: an HTTP dataset over a web linked service pulls files from external endpoints on a schedule.
- **Partner file exchange**: a binary dataset on an SFTP location moves opaque archives without parsing them.
- **Anything else Data Factory speaks**: the custom form carries any dataset type (Excel, XML, Avro, ORC, REST, ...) as typed JSON.

## Future Enhancements

- The integration runtime kind completes the Data Factory family; dataset-to-dataflow wiring deepens as the family's reference graph closes.
