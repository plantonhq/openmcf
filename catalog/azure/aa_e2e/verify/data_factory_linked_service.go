package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryLinkedServiceVerifier verifies an
// AzureDataFactoryLinkedService via the generic ARM resources GetByID
// (see armResourceExists), keyed on the linked service's ARM ID --
// one ID shape serves all 23 connection types
// ({factory_id}/linkedservices/{name}), on the same
// datafactory/2018-06-01 API line the pinned provider vendors.
type dataFactoryLinkedServiceVerifier struct{}

func (*dataFactoryLinkedServiceVerifier) IDOutputKey() string {
	return "linked_service_id"
}

func (*dataFactoryLinkedServiceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorylinkedservice verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactorylinkedservice %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryLinkedServiceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorylinkedservice verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactorylinkedservice %q still exists after destroy", id)
	}
	return nil
}
