# Azure Data Factory Dataset

Deploys one dataset inside an Azure Data Factory -- a named view of data telling pipelines and data flows what the data looks like and where it sits within a system a linked service already connects to, in any of 13 shapes plus a raw-JSON custom form. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one dataset of the shape the spec's variant block declares:

- **File formats** -- delimited text (CSV), JSON, Parquet, binary, and the flat Azure Blob form
- **HTTP** -- a file served by an HTTP endpoint, through a web linked service
- **Tables** -- Azure SQL Database, SQL Server, MySQL, PostgreSQL, Snowflake
- **Cosmos DB** -- a SQL API collection
- **Custom** -- any other Data Factory dataset type, as its ARM type name plus type-properties JSON

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory** -- the dataset lives in a factory; reference an AzureDataFactory's ID output (or provide a literal).
- **Linked service** -- every dataset reads through a connection; deploy the matching AzureDataFactoryLinkedService first and reference its name output (the Azure SQL table shape references its ARM ID instead).

### Azure Subscription

- **The linked service type must match the dataset shape** -- a delimited text dataset on a blob location needs a blob storage (or Data Lake Gen2) connection; a MySQL dataset needs a MySQL connection. Azure validates the pairing when a pipeline first USES the dataset, not at save time.
- **Saving does not read the data** -- Azure stores the definition without checking the file or table exists; a wrong path surfaces when a pipeline runs. Use Studio's Preview data button after deploy.
- **Declared columns are optional** -- Data Factory can infer or map columns at run time; declare them when downstream mappings need a stable contract.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Dataset**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CSV on Blob Storage** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-dataset.yaml
```

## After Deploy

The dataset appears in the factory's Studio under Author -> Datasets. Use **Preview data** there to verify the location and format parse correctly before pointing pipeline activities at it. The dataset stores metadata only -- deleting it never touches the underlying data.
