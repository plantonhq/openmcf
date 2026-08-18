package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appsync"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// datasources creates the API's data sources keyed by spec name.
//
// Lifecycle fact the render depends on: the EventBridge data source's
// update path silently DROPS its config upstream (recorded in _inbox)
// - treat EventBridge entries as replace-to-change by renaming them.
func datasources(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	api *createdApi) (map[string]*appsync.DataSource, error) {
	created := map[string]*appsync.DataSource{}

	for _, datasource := range locals.Spec.Datasources {
		datasourceArgs := &appsync.DataSourceArgs{
			ApiId: api.ApiId,
			Name:  pulumi.String(datasource.Name),
			Type:  pulumi.String(datasource.Type),
		}
		if datasource.Description != "" {
			datasourceArgs.Description = pulumi.String(datasource.Description)
		}
		if datasource.ServiceRoleArn.GetValue() != "" {
			datasourceArgs.ServiceRoleArn = pulumi.String(datasource.ServiceRoleArn.GetValue())
		}

		if dynamodb := datasource.Dynamodb; dynamodb != nil {
			dynamodbArgs := &appsync.DataSourceDynamodbConfigArgs{
				TableName: pulumi.String(dynamodb.TableName.GetValue()),
			}
			if dynamodb.Region != "" {
				dynamodbArgs.Region = pulumi.String(dynamodb.Region)
			}
			if dynamodb.UseCallerCredentials {
				dynamodbArgs.UseCallerCredentials = pulumi.Bool(true)
			}
			if dynamodb.Versioned {
				dynamodbArgs.Versioned = pulumi.Bool(true)
			}
			if deltaSync := dynamodb.DeltaSync; deltaSync != nil {
				deltaSyncArgs := &appsync.DataSourceDynamodbConfigDeltaSyncConfigArgs{
					DeltaSyncTableName: pulumi.String(deltaSync.DeltaSyncTableName),
				}
				if deltaSync.BaseTableTtl > 0 {
					deltaSyncArgs.BaseTableTtl = pulumi.Int(int(deltaSync.BaseTableTtl))
				}
				if deltaSync.DeltaSyncTableTtl > 0 {
					deltaSyncArgs.DeltaSyncTableTtl = pulumi.Int(int(deltaSync.DeltaSyncTableTtl))
				}
				dynamodbArgs.DeltaSyncConfig = deltaSyncArgs
			}
			datasourceArgs.DynamodbConfig = dynamodbArgs
		}

		if lambda := datasource.Lambda; lambda != nil {
			datasourceArgs.LambdaConfig = &appsync.DataSourceLambdaConfigArgs{
				FunctionArn: pulumi.String(lambda.FunctionArn.GetValue()),
			}
		}

		if http := datasource.Http; http != nil {
			httpArgs := &appsync.DataSourceHttpConfigArgs{
				Endpoint: pulumi.String(http.Endpoint),
			}
			// AWS_IAM is the only authorization type - pinned here; the
			// sigv4 block's presence selects signing.
			if sigv4 := http.Sigv4; sigv4 != nil {
				iamArgs := &appsync.DataSourceHttpConfigAuthorizationConfigAwsIamConfigArgs{}
				if sigv4.SigningRegion != "" {
					iamArgs.SigningRegion = pulumi.String(sigv4.SigningRegion)
				}
				if sigv4.SigningServiceName != "" {
					iamArgs.SigningServiceName = pulumi.String(sigv4.SigningServiceName)
				}
				httpArgs.AuthorizationConfig = &appsync.DataSourceHttpConfigAuthorizationConfigArgs{
					AuthorizationType: pulumi.String("AWS_IAM"),
					AwsIamConfig:      iamArgs,
				}
			}
			datasourceArgs.HttpConfig = httpArgs
		}

		if opensearch := datasource.Opensearch; opensearch != nil {
			opensearchArgs := &appsync.DataSourceOpensearchserviceConfigArgs{
				Endpoint: pulumi.String(opensearch.Endpoint.GetValue()),
			}
			if opensearch.Region != "" {
				opensearchArgs.Region = pulumi.String(opensearch.Region)
			}
			datasourceArgs.OpensearchserviceConfig = opensearchArgs
		}

		if elasticsearch := datasource.Elasticsearch; elasticsearch != nil {
			elasticsearchArgs := &appsync.DataSourceElasticsearchConfigArgs{
				Endpoint: pulumi.String(elasticsearch.Endpoint),
			}
			if elasticsearch.Region != "" {
				elasticsearchArgs.Region = pulumi.String(elasticsearch.Region)
			}
			datasourceArgs.ElasticsearchConfig = elasticsearchArgs
		}

		if eventbridge := datasource.Eventbridge; eventbridge != nil {
			datasourceArgs.EventBridgeConfig = &appsync.DataSourceEventBridgeConfigArgs{
				EventBusArn: pulumi.String(eventbridge.EventBusArn.GetValue()),
			}
		}

		if relationalDatabase := datasource.RelationalDatabase; relationalDatabase != nil {
			httpEndpointArgs := &appsync.DataSourceRelationalDatabaseConfigHttpEndpointConfigArgs{
				DbClusterIdentifier: pulumi.String(relationalDatabase.DbClusterIdentifier.GetValue()),
				AwsSecretStoreArn:   pulumi.String(relationalDatabase.AwsSecretStoreArn.GetValue()),
			}
			if relationalDatabase.DatabaseName != "" {
				httpEndpointArgs.DatabaseName = pulumi.String(relationalDatabase.DatabaseName)
			}
			if relationalDatabase.Schema != "" {
				httpEndpointArgs.Schema = pulumi.String(relationalDatabase.Schema)
			}
			if relationalDatabase.Region != "" {
				httpEndpointArgs.Region = pulumi.String(relationalDatabase.Region)
			}
			datasourceArgs.RelationalDatabaseConfig = &appsync.DataSourceRelationalDatabaseConfigArgs{
				// RDS_HTTP_ENDPOINT is the only source type - pinned here.
				SourceType:         pulumi.String("RDS_HTTP_ENDPOINT"),
				HttpEndpointConfig: httpEndpointArgs,
			}
		}

		createdDatasource, err := appsync.NewDataSource(ctx,
			fmt.Sprintf("datasource-%s", datasource.Name),
			datasourceArgs, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create datasource %s", datasource.Name)
		}
		created[datasource.Name] = createdDatasource
	}

	return created, nil
}
