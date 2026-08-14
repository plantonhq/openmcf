# Blob Landing Event

This preset creates the standard file-arrival trigger: fire the parse pipeline the moment a `.csv` lands under the landing path, passing the blob's location to the run.

## When to Use

- Processing files as they arrive instead of polling on a schedule
- Landing-zone patterns where upstream systems drop files at their own pace

## Key Configuration Choices

- **`activated: false`** -- deploy stopped, verify the pipeline's Debug run, then flip to true
- **Narrow path filters** -- `blobPathBeginsWith` + `blobPathEndsWith` are the only volume guard; an over-broad filter on a busy account fires the pipeline for every blob
- **`ignoreEmptyBlobs: true`** -- skips the zero-byte marker blobs some upstream tools write
- **Blob location as parameters** -- `@triggerBody().folderPath/fileName` tell the pipeline which file fired it

## Requirements

The **Microsoft.EventGrid resource provider must be registered** on the subscription -- Azure creates an Event Grid subscription on the storage account behind the scenes (that subscription is part of this trigger's failure domain).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-storage-account>` | The Planton name of the `AzureStorageAccount` watched for blobs | Planton console (or replace `valueFrom` with `value:` and the account's ARM ID) |
| `<your-pipeline>` | The Planton name of the `AzureDataFactoryPipeline` to run per file | Planton console (or replace `valueFrom` with `value:` and the pipeline's name) |
| `blobPathBeginsWith` | Your landing path | The container and prefix upstream systems drop files into |

## Related Presets

- **Custom Event** -- event-driven orchestration from an Event Grid topic.
- **Daily Schedule** -- clock-shaped work.
