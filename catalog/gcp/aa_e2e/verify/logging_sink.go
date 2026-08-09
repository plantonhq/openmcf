package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// loggingSinkVerifier probes a Cloud Logging sink and confirms the two
// facts the kind exists to guarantee: the rendered destination URI reached
// the API, and GCP minted a writer identity (the grant handle the sink's
// output contract promises). Live lanes exercise the project scope only —
// folder/org sinks need org credentials (the recorded deferral) and the
// test tenant has no billing-account fixture — so the probe addresses
// projects/{project}/sinks/{name} via the generic v2 Sinks service.
type loggingSinkVerifier struct{}

func (v *loggingSinkVerifier) IDOutputKey() string { return "sink_name" }

func (v *loggingSinkVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	sinkName := fmt.Sprintf("projects/%s/sinks/%s", svc.Project, outputs["sink_name"])
	sink, err := svc.Logging.Sinks.Get(sinkName).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "logging sink %s not found after deploy", sinkName)
	}
	if sink.Destination == "" {
		return errors.Errorf("logging sink %s reports no destination", sinkName)
	}
	if sink.WriterIdentity == "" {
		return errors.Errorf("logging sink %s reports no writer identity", sinkName)
	}
	return nil
}

func (v *loggingSinkVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	sinkName := fmt.Sprintf("projects/%s/sinks/%s", svc.Project, outputs["sink_name"])
	_, err := svc.Logging.Sinks.Get(sinkName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing logging sink %s after destroy", sinkName)
	}
	return errors.Errorf("logging sink %s still exists after destroy", sinkName)
}
