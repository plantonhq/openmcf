// Package aa_e2e implements the E2E provider harness for GCP. Like AWS (a real
// cloud account), Setup validates that the ambient Application Default
// Credentials chain can reach the test project, and resource verification runs
// through the Google Cloud REST APIs.
//
// Credentials are intentionally NOT plumbed through the stack input. The E2E
// framework builds every stack input with a nil provider config, so the IaC
// modules resolve credentials from the ambient ADC chain (locally:
// `gcloud auth application-default login`; in CI: workload identity
// federation). No static secret is ever stored on disk or in CI.
//
// The test project is resolved once at Setup (E2E_GCP_PROJECT, then
// GOOGLE_PROJECT, then the ADC credential's project) and exported as
// GOOGLE_PROJECT so both engines' subprocesses — and therefore both providers'
// default-project resolution — agree with the harness. Scenario manifests omit
// spec.project_id and ride this ambient project, so no manifest in the repo
// ever hardcodes a project id.
package aa_e2e

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/provider/gcp/aa_e2e/verify"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/redis/v1"
	run "google.golang.org/api/run/v2"
	"google.golang.org/api/spanner/v1"
	"google.golang.org/api/sqladmin/v1"
	"google.golang.org/api/storage/v1"
)

// Harness manages the GCP E2E test lifecycle.
type Harness struct {
	services *verify.Services

	// mu guards deployed, written by VerifyDeployed and read by VerifyDestroyed.
	mu       sync.Mutex
	deployed map[string]map[string]string
}

// NewHarness creates a GCP test harness. Credentials come from the ambient ADC
// chain (see the package doc); none are passed here.
func NewHarness() *Harness {
	return &Harness{deployed: make(map[string]map[string]string)}
}

// Setup resolves the test project, exports GOOGLE_PROJECT for the IaC
// subprocesses, loads ADC, and confirms the project is reachable via a
// side-effect-free cloudresourcemanager projects.get call.
func (h *Harness) Setup(ctx context.Context) error {
	creds, err := google.FindDefaultCredentials(ctx, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return errors.Wrap(err, "failed to load GCP Application Default Credentials "+
			"(locally: `gcloud auth application-default login`; in CI: workload identity federation)")
	}

	project := firstNonEmpty(os.Getenv("E2E_GCP_PROJECT"), os.Getenv("GOOGLE_PROJECT"), creds.ProjectID)
	if project == "" {
		return errors.New("no GCP test project resolved: set E2E_GCP_PROJECT (or GOOGLE_PROJECT), " +
			"or use ADC credentials that carry a project")
	}

	// Both the Terraform google provider and the Pulumi gcp provider resolve
	// their default project from GOOGLE_PROJECT, and both E2E spawners rebuild
	// the subprocess environment from this process at call time — so this single
	// export is what lets scenario manifests omit spec.project_id.
	if err := os.Setenv("GOOGLE_PROJECT", project); err != nil {
		return errors.Wrap(err, "failed to export GOOGLE_PROJECT")
	}

	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create cloudresourcemanager client")
	}
	iamService, err := iam.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create iam client")
	}
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create compute client")
	}
	storageService, err := storage.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create storage client")
	}
	sqlAdminService, err := sqladmin.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create sqladmin client")
	}
	redisService, err := redis.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create redis client")
	}
	containerService, err := container.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create container client")
	}
	runService, err := run.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create run client")
	}
	alloyDBService, err := alloydb.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create alloydb client")
	}
	dnsService, err := dns.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create dns client")
	}
	spannerService, err := spanner.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create spanner client")
	}
	bigQueryService, err := bigquery.NewService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create bigquery client")
	}

	gotProject, err := crmService.Projects.Get(project).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "GCP credential validation failed: cannot reach test project %q "+
			"(cloudresourcemanager projects.get)", project)
	}

	fmt.Printf("  [gcp] authenticated against project %s (%s)\n", gotProject.ProjectId, gotProject.LifecycleState)

	h.services = &verify.Services{
		Project:   project,
		Crm:       crmService,
		Iam:       iamService,
		Compute:   computeService,
		Storage:   storageService,
		SqlAdmin:  sqlAdminService,
		Redis:     redisService,
		Container: containerService,
		Run:       runService,
		AlloyDB:   alloyDBService,
		DNS:       dnsService,
		Spanner:   spannerService,
		BigQuery:  bigQueryService,
	}
	return nil
}

// Teardown is a no-op. Each scenario destroys its own resources in the DESTROY
// phase and confirms removal in VERIFY-CLN.
func (h *Harness) Teardown(ctx context.Context) error {
	return nil
}

// VerifyDeployed confirms the component's resource exists via its registered
// verifier. GCP identifiers are frequently compound (an IAM grant is a
// project+role+member tuple), so the whole string-ified output set is stored
// and handed to the verifier rather than a single id.
func (h *Harness) VerifyDeployed(ctx context.Context, component string, outputs map[string]interface{}) error {
	v, err := verify.GetVerifier(component)
	if err != nil {
		return err
	}

	strOutputs := stringOutputs(outputs)
	if strOutputs[v.IDOutputKey()] == "" {
		return errors.Errorf("no %q in outputs for %s -- cannot verify", v.IDOutputKey(), component)
	}

	h.mu.Lock()
	h.deployed[componentKey(ctx, component)] = strOutputs
	h.mu.Unlock()

	return v.VerifyExists(ctx, h.services, strOutputs)
}

// VerifyDestroyed confirms the previously deployed resource no longer exists.
func (h *Harness) VerifyDestroyed(ctx context.Context, component string) error {
	v, err := verify.GetVerifier(component)
	if err != nil {
		return err
	}

	h.mu.Lock()
	outputs := h.deployed[componentKey(ctx, component)]
	h.mu.Unlock()

	if len(outputs) == 0 {
		return errors.Errorf("no stored outputs for %s -- VerifyDeployed may not have run", component)
	}
	return v.VerifyAbsent(ctx, h.services, outputs)
}

// stringOutputs flattens stack outputs to strings, tolerating non-string scalars.
func stringOutputs(outputs map[string]interface{}) map[string]string {
	result := make(map[string]string, len(outputs))
	for key, value := range outputs {
		if s, ok := value.(string); ok {
			result[key] = s
			continue
		}
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

// componentKey combines the manifest path (from context) with the component name
// so concurrent scenarios of the same component type do not collide in the map.
func componentKey(ctx context.Context, component string) string {
	if mp, ok := ctx.Value(provider.ManifestPathKey{}).(string); ok && mp != "" {
		return mp + "::" + component
	}
	return component
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
