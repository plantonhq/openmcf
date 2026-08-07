package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// containerAppCustomDomainVerifier verifies an
// AzureContainerAppCustomDomain, keyed on the binding's synthetic ID
// ({container-app-id}/customDomainName/{domain}).
//
// STATE-AWARE by necessity: the binding has no ARM object of its own --
// it is an entry in the parent Container App's ingress configuration, so
// a 404 probe on the synthetic ID can never work. This verifier GETs the
// parent app and reads properties.configuration.ingress.customDomains:
// the binding exists when an entry carries the bound hostname, and it is
// absent when no entry does (a missing app also counts as absent -- the
// binding cannot outlive its app).
type containerAppCustomDomainVerifier struct{}

// IDOutputKey is the binding's synthetic ID.
func (*containerAppCustomDomainVerifier) IDOutputKey() string {
	return "custom_domain_id"
}

// customDomainBindingExists parses the synthetic ID into the parent app
// ID and the hostname, GETs the app, and reports whether the app's
// ingress carries the hostname.
func customDomainBindingExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) (bool, error) {
	parts := strings.SplitN(id, "/customDomainName/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false, pkgerrors.Errorf("custom domain id %q does not match {container-app-id}/customDomainName/{domain}", id)
	}
	containerAppID, domainName := parts[0], parts[1]

	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.GetByID(ctx, containerAppID, containerAppsAPIVersion, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}

	properties, ok := resp.Properties.(map[string]interface{})
	if !ok {
		return false, nil
	}
	configuration, ok := properties["configuration"].(map[string]interface{})
	if !ok {
		return false, nil
	}
	ingress, ok := configuration["ingress"].(map[string]interface{})
	if !ok {
		return false, nil
	}
	customDomains, ok := ingress["customDomains"].([]interface{})
	if !ok {
		return false, nil
	}
	for _, entry := range customDomains {
		domain, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := domain["name"].(string); ok && strings.EqualFold(name, domainName) {
			return true, nil
		}
	}
	return false, nil
}

func (*containerAppCustomDomainVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := customDomainBindingExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappcustomdomain verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappcustomdomain %q not found on the app's ingress after deploy", id)
	}
	return nil
}

func (*containerAppCustomDomainVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := customDomainBindingExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappcustomdomain verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappcustomdomain %q still present on the app's ingress after destroy", id)
	}
	return nil
}
