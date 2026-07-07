package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// dataprocAutoscalingPolicyVerifier probes a Dataproc autoscaling policy
// via the Dataproc API. The name output is the fully qualified resource
// name (projects/{p}/locations/{l}/autoscalingPolicies/{id}) — exactly
// the handle a cluster's autoscaling_policy_uri reference consumes, so
// verifying with it doubles as proof the composition handle is honest.
// Posture assertions confirm the policy carries worker bounds and the
// YARN algorithm — proof GCP accepted the policy as configured.
type dataprocAutoscalingPolicyVerifier struct{}

func (v *dataprocAutoscalingPolicyVerifier) IDOutputKey() string { return "name" }

func (v *dataprocAutoscalingPolicyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return errors.New("name output missing after deploy")
	}

	policy, err := svc.Dataproc.Projects.Locations.AutoscalingPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dataproc autoscaling policy %s not found after deploy", name)
	}

	if policy.WorkerConfig == nil || policy.WorkerConfig.MaxInstances == 0 {
		return errors.Errorf("dataproc autoscaling policy %s has no primary worker bounds after deploy", name)
	}
	if policy.BasicAlgorithm == nil || policy.BasicAlgorithm.YarnConfig == nil {
		return errors.Errorf("dataproc autoscaling policy %s has no YARN algorithm after deploy", name)
	}
	if outputs["policy_id"] != "" && policy.Id != outputs["policy_id"] {
		return errors.Errorf("dataproc autoscaling policy id mismatch: cloud %q vs output %q",
			policy.Id, outputs["policy_id"])
	}
	return nil
}

func (v *dataprocAutoscalingPolicyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return nil
	}

	_, err := svc.Dataproc.Projects.Locations.AutoscalingPolicies.Get(name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("dataproc autoscaling policy %s still exists after destroy", name)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing dataproc autoscaling policy %s after destroy", name)
}
