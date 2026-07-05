package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// satellites provisions the function-scoped settings AWS models as
// standalone resources but that are honestly part of the function's own
// configuration -- each is keyed by the function (replace-on-change),
// owned by exactly one function, and referenced by nothing else:
// aliases (with per-alias provisioned concurrency), the function URL,
// resource-policy invoke permissions, the asynchronous-invocation
// config, recursion detection, and runtime-update management.
//
// Event source mappings are deliberately NOT here: a mapping has
// independent AWS identity and wires OTHER resources (queues, streams)
// into the function, so it is its own first-class kind.
func satellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdFunction *lambda.Function) (pulumi.StringMap, *lambda.FunctionUrl, error) {
	spec := locals.AwsLambda.Spec

	// Named pointers to published versions -- the stable invocation
	// targets clients reference. Materialized per-name so list edits
	// update in place; repointing an alias ships or rolls back without
	// touching callers.
	aliasArns := pulumi.StringMap{}
	for _, alias := range spec.Aliases {
		aliasArgs := &lambda.AliasArgs{
			Name:            pulumi.String(alias.Name),
			FunctionName:    createdFunction.Name,
			FunctionVersion: pulumi.String(alias.FunctionVersion),
		}
		if alias.Description != "" {
			aliasArgs.Description = pulumi.String(alias.Description)
		}
		// Canary routing: at most one additional version and its
		// traffic fraction (AWS's own constraint).
		if len(alias.RoutingAdditionalVersionWeights) > 0 {
			weights := pulumi.Float64Map{}
			for version, weight := range alias.RoutingAdditionalVersionWeights {
				weights[version] = pulumi.Float64(weight)
			}
			aliasArgs.RoutingConfig = &lambda.AliasRoutingConfigArgs{
				AdditionalVersionWeights: weights,
			}
		}
		createdAlias, err := lambda.NewAlias(ctx, fmt.Sprintf("alias-%s", alias.Name),
			aliasArgs, pulumi.Provider(provider))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create alias %s", alias.Name)
		}
		aliasArns[alias.Name] = createdAlias.Arn

		// Pre-warmed execution environments keyed by the alias
		// qualifier -- eliminates cold starts on the alias at the cost
		// of paying for idle warmth.
		if alias.ProvisionedConcurrentExecutions != nil {
			if _, err := lambda.NewProvisionedConcurrencyConfig(ctx,
				fmt.Sprintf("provisioned-concurrency-%s", alias.Name),
				&lambda.ProvisionedConcurrencyConfigArgs{
					FunctionName:                    createdFunction.Name,
					Qualifier:                       createdAlias.Name,
					ProvisionedConcurrentExecutions: pulumi.Int(int(*alias.ProvisionedConcurrentExecutions)),
				}, pulumi.Provider(provider)); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to configure provisioned concurrency for alias %s", alias.Name)
			}
		}
	}

	// The built-in HTTPS endpoint. One per function
	// ($LATEST-qualified); with authorization_type NONE, AWS
	// additionally requires a public invoke permission, which the
	// provider manages on this resource.
	var createdFunctionUrl *lambda.FunctionUrl
	if spec.FunctionUrl != nil {
		urlArgs := &lambda.FunctionUrlArgs{
			FunctionName:      createdFunction.Name,
			AuthorizationType: pulumi.String(spec.FunctionUrl.AuthorizationType),
		}
		if spec.FunctionUrl.InvokeMode != "" {
			urlArgs.InvokeMode = pulumi.String(spec.FunctionUrl.InvokeMode)
		}
		if cors := spec.FunctionUrl.Cors; cors != nil {
			corsArgs := &lambda.FunctionUrlCorsArgs{
				AllowCredentials: pulumi.Bool(cors.AllowCredentials),
			}
			if len(cors.AllowOrigins) > 0 {
				corsArgs.AllowOrigins = pulumi.ToStringArray(cors.AllowOrigins)
			}
			if len(cors.AllowMethods) > 0 {
				corsArgs.AllowMethods = pulumi.ToStringArray(cors.AllowMethods)
			}
			if len(cors.AllowHeaders) > 0 {
				corsArgs.AllowHeaders = pulumi.ToStringArray(cors.AllowHeaders)
			}
			if len(cors.ExposeHeaders) > 0 {
				corsArgs.ExposeHeaders = pulumi.ToStringArray(cors.ExposeHeaders)
			}
			if cors.MaxAgeSeconds != 0 {
				corsArgs.MaxAge = pulumi.Int(int(cors.MaxAgeSeconds))
			}
			urlArgs.Cors = corsArgs
		}
		var err error
		createdFunctionUrl, err = lambda.NewFunctionUrl(ctx, "function-url", urlArgs, pulumi.Provider(provider))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create function URL")
		}
	}

	// Resource-policy statements authorizing external principals and
	// AWS services to invoke. Statements are create/delete-only in AWS
	// -- any field change replaces the statement (harmless: a statement
	// carries no state). Materialized per statement_id so list edits
	// update in place.
	for _, permission := range spec.InvokePermissions {
		// Empty keeps the sensible default: plain function invocation.
		action := permission.Action
		if action == "" {
			action = "lambda:InvokeFunction"
		}
		permissionArgs := &lambda.PermissionArgs{
			StatementId: pulumi.String(permission.StatementId),
			Function:    createdFunction.Name,
			Principal:   pulumi.String(permission.Principal),
			Action:      pulumi.String(action),
		}
		// source_arn/source_account scope service-principal grants to
		// one resource/account -- the confused-deputy guard.
		if permission.SourceArn != "" {
			permissionArgs.SourceArn = pulumi.String(permission.SourceArn)
		}
		if permission.SourceAccount != "" {
			permissionArgs.SourceAccount = pulumi.String(permission.SourceAccount)
		}
		if permission.PrincipalOrgId != "" {
			permissionArgs.PrincipalOrgId = pulumi.String(permission.PrincipalOrgId)
		}
		if permission.FunctionUrlAuthType != "" {
			permissionArgs.FunctionUrlAuthType = pulumi.String(permission.FunctionUrlAuthType)
		}
		if _, err := lambda.NewPermission(ctx, fmt.Sprintf("permission-%s", permission.StatementId),
			permissionArgs, pulumi.Provider(provider)); err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create invoke permission %s", permission.StatementId)
		}
	}

	// Asynchronous-invocation shaping: retries, event age, and
	// on-success / on-failure destinations. One config per function
	// ($LATEST-qualified).
	if async := spec.AsyncInvokeConfig; async != nil {
		asyncArgs := &lambda.FunctionEventInvokeConfigArgs{
			FunctionName: createdFunction.Name,
		}
		if async.MaximumRetryAttempts != nil {
			asyncArgs.MaximumRetryAttempts = pulumi.Int(int(*async.MaximumRetryAttempts))
		}
		if async.MaximumEventAgeSeconds != 0 {
			asyncArgs.MaximumEventAgeInSeconds = pulumi.Int(int(async.MaximumEventAgeSeconds))
		}
		onSuccess := async.OnSuccessDestinationArn.GetValue()
		onFailure := async.OnFailureDestinationArn.GetValue()
		if onSuccess != "" || onFailure != "" {
			destinationArgs := &lambda.FunctionEventInvokeConfigDestinationConfigArgs{}
			if onSuccess != "" {
				destinationArgs.OnSuccess = &lambda.FunctionEventInvokeConfigDestinationConfigOnSuccessArgs{
					Destination: pulumi.String(onSuccess),
				}
			}
			if onFailure != "" {
				destinationArgs.OnFailure = &lambda.FunctionEventInvokeConfigDestinationConfigOnFailureArgs{
					Destination: pulumi.String(onFailure),
				}
			}
			asyncArgs.DestinationConfig = destinationArgs
		}
		if _, err := lambda.NewFunctionEventInvokeConfig(ctx, "async-invoke-config",
			asyncArgs, pulumi.Provider(provider)); err != nil {
			return nil, nil, errors.Wrap(err, "failed to configure asynchronous invocation")
		}
	}

	// Recursive-loop detection. Only materialized when opting OUT of
	// the AWS default (Terminate) -- deleting the resource restores the
	// default, so rendering the default would be a no-op resource.
	if spec.RecursiveLoop == "Allow" {
		if _, err := lambda.NewFunctionRecursionConfig(ctx, "recursion-config",
			&lambda.FunctionRecursionConfigArgs{
				FunctionName:  createdFunction.Name,
				RecursiveLoop: pulumi.String(spec.RecursiveLoop),
			}, pulumi.Provider(provider)); err != nil {
			return nil, nil, errors.Wrap(err, "failed to configure recursion detection")
		}
	}

	// Runtime-update management. Only materialized when configured --
	// deleting the resource reverts to Auto (the AWS default).
	if rm := spec.RuntimeManagement; rm != nil {
		rmArgs := &lambda.RuntimeManagementConfigArgs{
			FunctionName:    createdFunction.Name,
			UpdateRuntimeOn: pulumi.String(rm.UpdateRuntimeOn),
		}
		if rm.RuntimeVersionArn != "" {
			rmArgs.RuntimeVersionArn = pulumi.String(rm.RuntimeVersionArn)
		}
		if _, err := lambda.NewRuntimeManagementConfig(ctx, "runtime-management-config",
			rmArgs, pulumi.Provider(provider)); err != nil {
			return nil, nil, errors.Wrap(err, "failed to configure runtime management")
		}
	}

	return aliasArns, createdFunctionUrl, nil
}
