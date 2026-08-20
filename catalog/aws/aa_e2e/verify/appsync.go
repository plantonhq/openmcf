package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsynctypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"
	pkgerrors "github.com/pkg/errors"
)

// appSyncApiVerifier verifies an AwsAppSyncApi keyed on the API id
// (the provider's import ID for either arm). The kind is a
// graphql-XOR-events union: Exists accepts whichever pivot answers -
// GetGraphqlApi for the GraphQL arm, GetApi for the Events arm - and
// Absent demands BOTH answer NotFound (an id can never be live on the
// other pivot after destroy).
type appSyncApiVerifier struct{}

func (*appSyncApiVerifier) IDOutputKey() string { return "api_id" }

func (*appSyncApiVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := appsync.NewFromConfig(cfg)
	if _, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(id)}); err == nil {
		return nil
	} else if !isAppSyncNotFound(err) {
		return pkgerrors.Wrapf(err, "awsappsyncapi verify-exists (graphql) failed for %q", id)
	}
	if _, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(id)}); err != nil {
		if isAppSyncNotFound(err) {
			return pkgerrors.Errorf("awsappsyncapi %q not found on either pivot (graphql or events)", id)
		}
		return pkgerrors.Wrapf(err, "awsappsyncapi verify-exists (events) failed for %q", id)
	}
	return nil
}

func (*appSyncApiVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := appsync.NewFromConfig(cfg)
	if _, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(id)}); err == nil {
		return pkgerrors.Errorf("awsappsyncapi %q still exists after destroy (graphql pivot)", id)
	} else if !isAppSyncNotFound(err) {
		return pkgerrors.Wrapf(err, "awsappsyncapi verify-absent (graphql) failed for %q", id)
	}
	if _, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(id)}); err == nil {
		return pkgerrors.Errorf("awsappsyncapi %q still exists after destroy (events pivot)", id)
	} else if !isAppSyncNotFound(err) {
		return pkgerrors.Wrapf(err, "awsappsyncapi verify-absent (events) failed for %q", id)
	}
	return nil
}

// VerifyExistsFromOutputs additionally walks the declared data sources
// and functions (GraphQL arm) and channel namespaces (Events arm) -
// every satellite the outputs claim must answer on the live API.
func (v *appSyncApiVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	apiId, _ := outputs["api_id"].(string)
	if apiId == "" {
		return pkgerrors.New("awsappsyncapi outputs carry no api_id")
	}
	if err := v.VerifyExists(ctx, cfg, apiId, region); err != nil {
		return err
	}
	client := appsync.NewFromConfig(cfg)

	datasourceArns, _ := outputs["datasource_arns"].(map[string]interface{})
	for datasourceName := range datasourceArns {
		if _, err := client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(datasourceName),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awsappsyncapi data source %q read failed", datasourceName)
		}
	}

	functionIds, _ := outputs["function_ids"].(map[string]interface{})
	for functionName, functionId := range functionIds {
		id, _ := functionId.(string)
		if id == "" {
			return pkgerrors.Errorf("awsappsyncapi function %q carries an empty id output", functionName)
		}
		if _, err := client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(apiId),
			FunctionId: aws.String(id),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awsappsyncapi function %q read failed", functionName)
		}
	}

	channelNamespaceArns, _ := outputs["channel_namespace_arns"].(map[string]interface{})
	for namespaceName := range channelNamespaceArns {
		if _, err := client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(namespaceName),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awsappsyncapi channel namespace %q read failed", namespaceName)
		}
	}

	return nil
}

// isAppSyncNotFound reports whether err is AppSync's NotFoundException
// (the absent signal for both pivots).
func isAppSyncNotFound(err error) bool {
	var notFound *appsynctypes.NotFoundException
	return pkgerrors.As(err, &notFound)
}
