package module

import (
	gcpbigquerytablev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbigquerytable/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/bigquery"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildExternalDataConfiguration maps the spec's external-table arm to the
// provider's block: the data stays in the source (GCS, Sheets, Bigtable,
// ...) and BigQuery reads it at query time.
func buildExternalDataConfiguration(
	config *gcpbigquerytablev1alpha1.GcpBigQueryTableExternalDataConfiguration,
) *bigquery.TableExternalDataConfigurationArgs {
	// Autodetect is always sent explicitly (the API requires the field).
	configArgs := &bigquery.TableExternalDataConfigurationArgs{
		Autodetect: pulumi.Bool(config.Autodetect),
		SourceUris: pulumi.ToStringArray(config.SourceUris),
	}

	if config.SourceFormat != "" {
		configArgs.SourceFormat = pulumi.StringPtr(config.SourceFormat)
	}
	if config.ObjectMetadata != "" {
		configArgs.ObjectMetadata = pulumi.StringPtr(config.ObjectMetadata)
	}
	if config.Compression != "" {
		configArgs.Compression = pulumi.StringPtr(config.Compression)
	}
	if config.Schema != "" {
		configArgs.Schema = pulumi.StringPtr(config.Schema)
	}
	if config.IgnoreUnknownValues {
		configArgs.IgnoreUnknownValues = pulumi.BoolPtr(true)
	}
	if config.MaxBadRecords > 0 {
		configArgs.MaxBadRecords = pulumi.IntPtr(int(config.MaxBadRecords))
	}
	if config.ConnectionId != "" {
		configArgs.ConnectionId = pulumi.StringPtr(config.ConnectionId)
	}
	if config.ReferenceFileSchemaUri != "" {
		configArgs.ReferenceFileSchemaUri = pulumi.StringPtr(config.ReferenceFileSchemaUri)
	}
	if config.MetadataCacheMode != "" {
		configArgs.MetadataCacheMode = pulumi.StringPtr(config.MetadataCacheMode)
	}
	if config.FileSetSpecType != "" {
		configArgs.FileSetSpecType = pulumi.StringPtr(config.FileSetSpecType)
	}
	if config.JsonExtension != "" {
		configArgs.JsonExtension = pulumi.StringPtr(config.JsonExtension)
	}

	if config.CsvOptions != nil {
		// The provider requires quote; an unset spec value means the API
		// default double-quote, while an explicit "" means unquoted data
		// (why the spec field is presence-tracked).
		quote := "\""
		if config.CsvOptions.Quote != nil {
			quote = config.CsvOptions.GetQuote()
		}
		csvArgs := &bigquery.TableExternalDataConfigurationCsvOptionsArgs{
			Quote: pulumi.String(quote),
		}
		if config.CsvOptions.AllowJaggedRows {
			csvArgs.AllowJaggedRows = pulumi.BoolPtr(true)
		}
		if config.CsvOptions.AllowQuotedNewlines {
			csvArgs.AllowQuotedNewlines = pulumi.BoolPtr(true)
		}
		if config.CsvOptions.Encoding != "" {
			csvArgs.Encoding = pulumi.StringPtr(config.CsvOptions.Encoding)
		}
		if config.CsvOptions.FieldDelimiter != "" {
			csvArgs.FieldDelimiter = pulumi.StringPtr(config.CsvOptions.FieldDelimiter)
		}
		if config.CsvOptions.SkipLeadingRows > 0 {
			csvArgs.SkipLeadingRows = pulumi.IntPtr(int(config.CsvOptions.SkipLeadingRows))
		}
		configArgs.CsvOptions = csvArgs
	}

	if config.JsonOptions != nil {
		jsonArgs := &bigquery.TableExternalDataConfigurationJsonOptionsArgs{}
		if config.JsonOptions.Encoding != "" {
			jsonArgs.Encoding = pulumi.StringPtr(config.JsonOptions.Encoding)
		}
		configArgs.JsonOptions = jsonArgs
	}

	if config.GoogleSheetsOptions != nil {
		sheetsArgs := &bigquery.TableExternalDataConfigurationGoogleSheetsOptionsArgs{}
		if config.GoogleSheetsOptions.Range != "" {
			sheetsArgs.Range = pulumi.StringPtr(config.GoogleSheetsOptions.Range)
		}
		if config.GoogleSheetsOptions.SkipLeadingRows > 0 {
			sheetsArgs.SkipLeadingRows = pulumi.IntPtr(int(config.GoogleSheetsOptions.SkipLeadingRows))
		}
		configArgs.GoogleSheetsOptions = sheetsArgs
	}

	if config.HivePartitioningOptions != nil {
		hiveArgs := &bigquery.TableExternalDataConfigurationHivePartitioningOptionsArgs{}
		if config.HivePartitioningOptions.Mode != "" {
			hiveArgs.Mode = pulumi.StringPtr(config.HivePartitioningOptions.Mode)
		}
		if config.HivePartitioningOptions.RequirePartitionFilter {
			hiveArgs.RequirePartitionFilter = pulumi.BoolPtr(true)
		}
		if config.HivePartitioningOptions.SourceUriPrefix != "" {
			hiveArgs.SourceUriPrefix = pulumi.StringPtr(config.HivePartitioningOptions.SourceUriPrefix)
		}
		configArgs.HivePartitioningOptions = hiveArgs
	}

	if config.AvroOptions != nil {
		// Always sent explicitly (the provider requires the field).
		configArgs.AvroOptions = &bigquery.TableExternalDataConfigurationAvroOptionsArgs{
			UseAvroLogicalTypes: pulumi.Bool(config.AvroOptions.UseAvroLogicalTypes),
		}
	}

	if config.ParquetOptions != nil {
		parquetArgs := &bigquery.TableExternalDataConfigurationParquetOptionsArgs{}
		if config.ParquetOptions.EnumAsString {
			parquetArgs.EnumAsString = pulumi.BoolPtr(true)
		}
		if config.ParquetOptions.EnableListInference {
			parquetArgs.EnableListInference = pulumi.BoolPtr(true)
		}
		configArgs.ParquetOptions = parquetArgs
	}

	if config.BigtableOptions != nil {
		bigtableArgs := &bigquery.TableExternalDataConfigurationBigtableOptionsArgs{}
		if config.BigtableOptions.IgnoreUnspecifiedColumnFamilies {
			bigtableArgs.IgnoreUnspecifiedColumnFamilies = pulumi.BoolPtr(true)
		}
		if config.BigtableOptions.ReadRowkeyAsString {
			bigtableArgs.ReadRowkeyAsString = pulumi.BoolPtr(true)
		}
		if config.BigtableOptions.OutputColumnFamiliesAsJson {
			bigtableArgs.OutputColumnFamiliesAsJson = pulumi.BoolPtr(true)
		}
		if len(config.BigtableOptions.ColumnFamilies) > 0 {
			families := make(bigquery.TableExternalDataConfigurationBigtableOptionsColumnFamilyArray, 0,
				len(config.BigtableOptions.ColumnFamilies))
			for _, family := range config.BigtableOptions.ColumnFamilies {
				familyArgs := &bigquery.TableExternalDataConfigurationBigtableOptionsColumnFamilyArgs{}
				if family.FamilyId != "" {
					familyArgs.FamilyId = pulumi.StringPtr(family.FamilyId)
				}
				if family.Type != "" {
					familyArgs.Type = pulumi.StringPtr(family.Type)
				}
				if family.Encoding != "" {
					familyArgs.Encoding = pulumi.StringPtr(family.Encoding)
				}
				if family.OnlyReadLatest {
					familyArgs.OnlyReadLatest = pulumi.BoolPtr(true)
				}
				if len(family.Columns) > 0 {
					columns := make(bigquery.TableExternalDataConfigurationBigtableOptionsColumnFamilyColumnArray, 0,
						len(family.Columns))
					for _, column := range family.Columns {
						columnArgs := &bigquery.TableExternalDataConfigurationBigtableOptionsColumnFamilyColumnArgs{}
						if column.QualifierEncoded != "" {
							columnArgs.QualifierEncoded = pulumi.StringPtr(column.QualifierEncoded)
						}
						if column.QualifierString != "" {
							columnArgs.QualifierString = pulumi.StringPtr(column.QualifierString)
						}
						if column.FieldName != "" {
							columnArgs.FieldName = pulumi.StringPtr(column.FieldName)
						}
						if column.Type != "" {
							columnArgs.Type = pulumi.StringPtr(column.Type)
						}
						if column.Encoding != "" {
							columnArgs.Encoding = pulumi.StringPtr(column.Encoding)
						}
						if column.OnlyReadLatest {
							columnArgs.OnlyReadLatest = pulumi.BoolPtr(true)
						}
						columns = append(columns, columnArgs)
					}
					familyArgs.Columns = columns
				}
				families = append(families, familyArgs)
			}
			bigtableArgs.ColumnFamilies = families
		}
		configArgs.BigtableOptions = bigtableArgs
	}

	return configArgs
}
