package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// bigtableTableVerifier probes a Bigtable table via the Bigtable Admin
// API. The SCHEMA view returns column families, so the posture assertion
// confirms the table's declared families actually exist — the part of the
// table an application depends on.
type bigtableTableVerifier struct{}

func (v *bigtableTableVerifier) IDOutputKey() string { return "table_id" }

func (v *bigtableTableVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	tableId := outputs["table_id"]
	if tableId == "" {
		return errors.New("table_id output missing after deploy")
	}

	table, err := svc.BigtableAdmin.Projects.Instances.Tables.Get(tableId).View("SCHEMA_VIEW").Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "bigtable table %s not found after deploy", tableId)
	}
	if len(table.ColumnFamilies) == 0 {
		return errors.Errorf("bigtable table %s has no column families — applications cannot write to it", tableId)
	}
	return nil
}

func (v *bigtableTableVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	tableId := outputs["table_id"]
	if tableId == "" {
		return nil
	}

	_, err := svc.BigtableAdmin.Projects.Instances.Tables.Get(tableId).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing bigtable table %s after destroy", tableId)
	}
	return errors.Errorf("bigtable table %s still exists after destroy", tableId)
}
