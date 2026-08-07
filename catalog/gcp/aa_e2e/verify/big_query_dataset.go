package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// bigQueryDatasetVerifier probes a BigQuery dataset via the bigquery v2 API.
// Posture assertions confirm the dataset exists and its live location matches
// the location output.
type bigQueryDatasetVerifier struct{}

func (v *bigQueryDatasetVerifier) IDOutputKey() string { return "dataset_id" }

func (v *bigQueryDatasetVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	datasetID := outputs["dataset_id"]
	if datasetID == "" {
		return errors.New("dataset_id output missing after deploy")
	}

	project := outputs["project"]
	if project == "" {
		project = svc.Project
	}

	dataset, err := svc.BigQuery.Datasets.Get(project, datasetID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "bigquery dataset %s in project %s not found after deploy", datasetID, project)
	}

	if wantLoc := outputs["location"]; wantLoc != "" && dataset.Location != wantLoc {
		return errors.Errorf("bigquery dataset %s location mismatch: output %q, live %q",
			datasetID, wantLoc, dataset.Location)
	}
	return nil
}

func (v *bigQueryDatasetVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	datasetID := outputs["dataset_id"]
	if datasetID == "" {
		return nil
	}

	project := outputs["project"]
	if project == "" {
		project = svc.Project
	}

	_, err := svc.BigQuery.Datasets.Get(project, datasetID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("bigquery dataset %s in project %s still exists after destroy", datasetID, project)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing bigquery dataset %s after destroy", datasetID)
}
