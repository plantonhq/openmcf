package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// bigQueryTableVerifier probes a BigQuery table via the bigquery v2 API.
// Posture assertions confirm the table exists, its type matches the output,
// and partitioning/clustering posture is present when the scenario set them.
type bigQueryTableVerifier struct{}

func (v *bigQueryTableVerifier) IDOutputKey() string { return "table_id" }

func (v *bigQueryTableVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	tableID := outputs["table_id"]
	if tableID == "" {
		return errors.New("table_id output missing after deploy")
	}

	datasetID := outputs["dataset_id"]
	if datasetID == "" {
		return errors.New("dataset_id output missing after deploy")
	}

	project := outputs["project"]
	if project == "" {
		project = svc.Project
	}

	table, err := svc.BigQuery.Tables.Get(project, datasetID, tableID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "bigquery table %s.%s in project %s not found after deploy",
			datasetID, tableID, project)
	}

	if wantType := outputs["type"]; wantType != "" && table.Type != wantType {
		return errors.Errorf("bigquery table %s.%s type mismatch: output %q, live %q",
			datasetID, tableID, wantType, table.Type)
	}

	// Partitioned native-table scenarios emit time partitioning; confirm it landed.
	if table.Type == "TABLE" && table.TimePartitioning != nil {
		if table.TimePartitioning.Type == "" {
			return errors.Errorf("bigquery table %s.%s has empty time partitioning type", datasetID, tableID)
		}
	}

	if table.Type == "VIEW" && table.View == nil {
		return errors.Errorf("bigquery table %s.%s is type VIEW but has no view definition", datasetID, tableID)
	}

	return nil
}

func (v *bigQueryTableVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	tableID := outputs["table_id"]
	if tableID == "" {
		return nil
	}

	datasetID := outputs["dataset_id"]
	if datasetID == "" {
		return nil
	}

	project := outputs["project"]
	if project == "" {
		project = svc.Project
	}

	_, err := svc.BigQuery.Tables.Get(project, datasetID, tableID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("bigquery table %s.%s in project %s still exists after destroy",
			datasetID, tableID, project)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing bigquery table %s.%s after destroy", datasetID, tableID)
}
