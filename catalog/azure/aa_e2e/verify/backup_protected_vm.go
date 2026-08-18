package verify

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// recoveryServicesProtectedItemAPIVersion pins the
// Microsoft.RecoveryServices backup GA line the protected-item
// verifier reads at -- the same line the pinned azurerm provider
// vendors for protected items (recoveryservicesbackup/2023-02-01).
const recoveryServicesProtectedItemAPIVersion = "2023-02-01"

// backupProtectedVmVerifier verifies an AzureBackupProtectedVm via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// protected item's full ARM ID (.../backupFabrics/Azure/
// protectionContainers/.../protectedItems/VM;iaasvmcontainerv2;...).
// Existence is the honest bar: creation only REGISTERS protection (the
// first backup runs on the policy's schedule), so the registration
// object is what a smoke lane can verify.
//
// ABSENCE-AFTER-DESTROY (measured live, proof session 055): the
// protected-item delete is asynchronous BEYOND the provider's poller
// -- ARM reads keep answering 200 for minutes after a destroy the
// engine already reported successful (a smoke item has ZERO recovery
// points, so a landed delete removes it outright; the 14-day
// soft-delete ghost class only applies to items WITH recovery
// points). So absence is verified with a bounded poll: 404 is
// absence, a 200 whose properties.isScheduledForDeferredDelete is
// true is ALSO absence (the bar the azurerm provider itself uses for
// ghosts), and an item still ACTIVE at the deadline fails honestly.
// The honest failure is load-bearing: one measured run had the
// provider's DeleteThenPoll return success while Azure ran NO
// DeleteBackupData job at all and the item survived active
// indefinitely -- a service-side drop this poll must not absorb.
// Recovery for that flake: `az backup protection disable
// --delete-backup-data true` (lands in seconds), then re-run.
type backupProtectedVmVerifier struct{}

// recoveryServicesProtectedItemAbsencePollTimeout bounds the absence
// poll (shared by the VM and file-share protected-item verifiers):
// long enough to absorb the measured minutes-scale read lag after a
// landed delete, short enough to surface a dropped delete within the
// lane.
const recoveryServicesProtectedItemAbsencePollTimeout = 10 * time.Minute

const recoveryServicesProtectedItemAbsencePollInterval = 15 * time.Second

// IDOutputKey is the protected item's full ARM ID.
func (*backupProtectedVmVerifier) IDOutputKey() string {
	return "backup_protected_vm_id"
}

func (*backupProtectedVmVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectedItemAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedvm verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackupprotectedvm %q not found after deploy", id)
	}
	return nil
}

func (*backupProtectedVmVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedvm verify-absent failed for %q", id)
	}
	deadline := time.Now().Add(recoveryServicesProtectedItemAbsencePollTimeout)
	for {
		resp, err := client.GetByID(ctx, id, recoveryServicesProtectedItemAPIVersion, nil)
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) && respErr.StatusCode == 404 {
				return nil
			}
			return pkgerrors.Wrapf(err, "azurebackupprotectedvm verify-absent failed for %q", id)
		}
		// The read answered: a deferred-delete ghost is absence; an
		// active item may just be the async delete's read lag -- keep
		// polling until the deadline before failing.
		if properties, ok := resp.Properties.(map[string]interface{}); ok {
			if deferred, _ := properties["isScheduledForDeferredDelete"].(bool); deferred {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return pkgerrors.Errorf("azurebackupprotectedvm %q still exists (active, not soft-delete-ghosted) %s after destroy -- the engine-reported delete never landed", id, recoveryServicesProtectedItemAbsencePollTimeout)
		}
		select {
		case <-ctx.Done():
			return pkgerrors.Wrapf(ctx.Err(), "azurebackupprotectedvm verify-absent cancelled for %q", id)
		case <-time.After(recoveryServicesProtectedItemAbsencePollInterval):
		}
	}
}
