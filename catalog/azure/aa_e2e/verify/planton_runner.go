package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// plantonRunnerVerifier verifies an AzurePlantonRunner via the generic ARM
// resources GetByID (see armResourceExists), keyed on the runner's
// Container App ARM ID. The appliance is a single-replica Container App;
// Container Apps reports the app provisioned independently of replica
// health, so existence is the honest provisioning-level assertion (the
// runner's own join is proven at the control plane, not here).
type plantonRunnerVerifier struct{}

// IDOutputKey is the app's full ARM ID.
func (*plantonRunnerVerifier) IDOutputKey() string {
	return "container_app_id"
}

func (*plantonRunnerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureplantonrunner verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureplantonrunner %q not found after deploy", id)
	}
	return nil
}

func (*plantonRunnerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureplantonrunner verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureplantonrunner %q still exists after destroy", id)
	}
	return nil
}

// VerifyRuntimeFailureCause pins the fake-token app's designed failure as far
// as ARM's management plane can see. Reaching this phase at all proves the
// image PULLED: Container Apps validates image pullability during the app
// create (proven live -- an unpullable manifest fails the whole deployment),
// and this phase only runs after a successful deploy. The revision state then
// attests the container is failing by design: ARM marks a revision whose
// container runs-and-exits (the refused join; the startup probe never passes)
// as provisioningState=Failed with healthState=Unhealthy / runningState=Failed
// -- the exact live-observed triple. The join-refusal SPECIFICITY -- that the
// process exited on the join step and not something else -- is not readable
// through ARM without a Log Analytics workspace (the consumption-only test
// environment streams logs only), and is pinned by the Kubernetes and Cloud
// Run lanes running the SAME binary and image; this assertion carries Azure's
// share of the evidence honestly rather than claiming more than ARM exposes.
func (*plantonRunnerVerifier) VerifyRuntimeFailureCause(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id, cause string) error {
	if cause != "refused-join" {
		return pkgerrors.Errorf("unsupported runtime failure cause %q for the runner (supported: refused-join)", cause)
	}

	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	deadline := time.Now().Add(3 * time.Minute)
	lastState := "no revision state read yet"
	for {
		app, err := client.GetByID(ctx, id, containerAppsAPIVersion, nil)
		if err != nil {
			return pkgerrors.Wrapf(err, "reading the app for revision state (%q)", id)
		}
		if props, ok := app.Properties.(map[string]interface{}); ok {
			if rev, ok := props["latestRevisionName"].(string); ok && rev != "" {
				revResp, revErr := client.GetByID(ctx, id+"/revisions/"+rev, containerAppsAPIVersion, nil)
				if revErr == nil {
					if rp, ok := revResp.Properties.(map[string]interface{}); ok {
						provisioning, _ := rp["provisioningState"].(string)
						health, _ := rp["healthState"].(string)
						running, _ := rp["runningState"].(string)
						lastState = fmt.Sprintf("revision %s: provisioningState=%s healthState=%s runningState=%s",
							rev, provisioning, health, running)
						if health == "Unhealthy" || running == "Degraded" || running == "Failed" {
							fmt.Printf("  [verify] CAUSE: %s -- the image pulled (create validates pullability; this phase presupposes a successful deploy) and the container is failing by design\n", lastState)
							return nil
						}
					}
				}
			}
		}

		if time.Now().After(deadline) {
			return pkgerrors.Errorf("the runner revision never attested the designed running-and-failing state within the window; last: %s", lastState)
		}
		time.Sleep(10 * time.Second)
	}
}
