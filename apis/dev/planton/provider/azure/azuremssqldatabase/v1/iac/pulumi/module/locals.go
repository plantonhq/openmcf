package module

import (
	"strings"

	azuremssqldatabasev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremssqldatabase/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMssqlDatabase *azuremssqldatabasev1.AzureMssqlDatabase
	AzureTags          map[string]string
}

// createModeStrings maps the spec's create-mode enum to ARM's values. An
// unspecified mode means a fresh (DEFAULT) database and is not sent at
// all, mirroring the Terraform module's null -- so the two engines
// produce the same ARM payload.
var createModeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseCreateMode]string{
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_DEFAULT:                            "Default",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_COPY:                               "Copy",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_SECONDARY:                          "Secondary",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_ONLINE_SECONDARY:                   "OnlineSecondary",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_POINT_IN_TIME_RESTORE:              "PointInTimeRestore",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_RECOVERY:                           "Recovery",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_RESTORE:                            "Restore",
	azuremssqldatabasev1.AzureMssqlDatabaseCreateMode_RESTORE_LONG_TERM_RETENTION_BACKUP: "RestoreLongTermRetentionBackup",
}

// secondaryTypeStrings maps the spec's secondary-type enum to ARM's
// values.
var secondaryTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseSecondaryType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseSecondaryType_GEO:   "Geo",
	azuremssqldatabasev1.AzureMssqlDatabaseSecondaryType_NAMED: "Named",
}

// licenseTypeStrings maps the Azure Hybrid Benefit enum to ARM's values.
var licenseTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseLicenseType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseLicenseType_BASE_PRICE:       "BasePrice",
	azuremssqldatabasev1.AzureMssqlDatabaseLicenseType_LICENSE_INCLUDED: "LicenseIncluded",
}

// enclaveTypeStrings maps the confidential-computing enclave enum to ARM's
// values.
var enclaveTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseEnclaveType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseEnclaveType_VBS:             "VBS",
	azuremssqldatabasev1.AzureMssqlDatabaseEnclaveType_DEFAULT_ENCLAVE: "Default",
}

// storageAccountTypeStrings maps the backup-redundancy enum to ARM's
// values.
var storageAccountTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseBackupStorageAccountType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseBackupStorageAccountType_GEO_REDUNDANT:      "Geo",
	azuremssqldatabasev1.AzureMssqlDatabaseBackupStorageAccountType_GEO_ZONE_REDUNDANT: "GeoZone",
	azuremssqldatabasev1.AzureMssqlDatabaseBackupStorageAccountType_LOCALLY_REDUNDANT:  "Local",
	azuremssqldatabasev1.AzureMssqlDatabaseBackupStorageAccountType_ZONE_REDUNDANT:     "Zone",
}

// importStorageKeyTypeStrings maps the bacpac credential-kind enum to
// ARM's values.
var importStorageKeyTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseImportStorageKeyType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseImportStorageKeyType_SHARED_ACCESS_KEY:  "SharedAccessKey",
	azuremssqldatabasev1.AzureMssqlDatabaseImportStorageKeyType_STORAGE_ACCESS_KEY: "StorageAccessKey",
}

// importAuthenticationTypeStrings maps the bacpac auth enum to ARM's
// values.
var importAuthenticationTypeStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseImportAuthenticationType]string{
	azuremssqldatabasev1.AzureMssqlDatabaseImportAuthenticationType_SQL:         "Sql",
	azuremssqldatabasev1.AzureMssqlDatabaseImportAuthenticationType_AD_PASSWORD: "ADPassword",
}

// threatDetectionStateStrings maps the Defender policy-state enum to ARM's
// values.
var threatDetectionStateStrings = map[azuremssqldatabasev1.AzureMssqlDatabaseThreatDetectionState]string{
	azuremssqldatabasev1.AzureMssqlDatabaseThreatDetectionState_ENABLED:  "Enabled",
	azuremssqldatabasev1.AzureMssqlDatabaseThreatDetectionState_DISABLED: "Disabled",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremssqldatabasev1.AzureMssqlDatabaseStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMssqlDatabase = stackInput.Target
	target := stackInput.Target

	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMssqlDatabase.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	// The user's spec tags merge over the metadata-derived tags -- user
	// tags deliberately win so an org's governance conventions can
	// override the derived values where they collide.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
