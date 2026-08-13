package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// networkWatcherFlowLogAPIVersion is the stable Microsoft.Network API
// version the existence probe is pinned to -- the line the pinned
// provider vendors for flowLogs.
const networkWatcherFlowLogAPIVersion = "2025-01-01"

// networkWatcherFlowLogVerifier verifies an AzureNetworkWatcherFlowLog
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the flow log's full ARM ID (a child of the regional Network
// Watcher: .../networkWatchers/{watcher}/flowLogs/{name}).
type networkWatcherFlowLogVerifier struct{}

// IDOutputKey is the flow log's full ARM ID.
func (*networkWatcherFlowLogVerifier) IDOutputKey() string {
	return "flow_log_id"
}

func (*networkWatcherFlowLogVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkWatcherFlowLogAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworkwatcherflowlog verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurenetworkwatcherflowlog %q not found after deploy", id)
	}
	return nil
}

func (*networkWatcherFlowLogVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkWatcherFlowLogAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworkwatcherflowlog verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurenetworkwatcherflowlog %q still exists after destroy", id)
	}
	return nil
}
