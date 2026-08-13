package module

import (
	"github.com/pkg/errors"
	awsrestapigatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapigateway/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createdTree carries what the deployment needs from the route tree:
// every method/integration resource, so the deployment waits for the
// whole definition.
type createdTree struct {
	definitionResources []pulumi.Resource
}

// routeTree derives the API Gateway resource tree from the route paths
// and creates the method, integration, and response resources per
// route.
//
// Lifecycle facts the renders below depend on:
//   - resources are created shallow-first (locals pre-sorts the derived
//     paths), so a parent always exists before its children;
//   - method responses are serialized upstream behind a global mutex
//     (concurrent writes on one API conflict) with a 2-minute retry;
//   - an integration response requires BOTH its method response and
//     its integration to exist - AWS returns NotFound otherwise; the
//     explicit DependsOn below carries that edge;
//   - request-parameter updates on an integration can mutate its cache
//     key parameters; the provider reconciles that with a re-read
//     after a 3-second settle.
func routeTree(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, api *apigateway.RestApi, satellites *createdSatellites) (*createdTree, error) {
	spec := locals.Spec
	tree := &createdTree{}

	// The derived resource nodes, parents first.
	resources := map[string]*apigateway.Resource{}
	resourceIds := pulumi.StringMap{}
	for _, path := range locals.ResourcePaths {
		var parent pulumi.StringInput
		if pp := parentPath(path); pp == "" {
			parent = api.RootResourceId
		} else {
			parent = resources[pp].ID()
		}
		resource, err := apigateway.NewResource(ctx, "resource-"+path, &apigateway.ResourceArgs{
			RestApi:  api.ID(),
			ParentId: parent,
			PathPart: pulumi.String(lastSegment(path)),
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create resource %q", path)
		}
		resources[path] = resource
		resourceIds[path] = resource.ID().ToStringOutput()
	}
	ctx.Export(OpResourceIds, resourceIds)

	// Methods, integrations, and responses per route.
	for _, r := range sortedRoutes(spec.Routes) {
		key := routeKey(r)

		var resourceId pulumi.StringInput
		if r.Path == "/" {
			resourceId = api.RootResourceId
		} else {
			resourceId = resources[r.Path].ID()
		}

		authorization := "NONE"
		if r.Authorization != "" {
			authorization = r.Authorization
		}
		methodArgs := &apigateway.MethodArgs{
			RestApi:       api.ID(),
			ResourceId:    resourceId,
			HttpMethod:    pulumi.String(r.Method),
			Authorization: pulumi.String(authorization),
		}
		if r.AuthorizerName != "" {
			methodArgs.AuthorizerId = satellites.authorizers[r.AuthorizerName].ID()
		}
		if len(r.AuthorizationScopes) > 0 {
			methodArgs.AuthorizationScopes = pulumi.ToStringArray(r.AuthorizationScopes)
		}
		if r.ApiKeyRequired {
			methodArgs.ApiKeyRequired = pulumi.Bool(true)
		}
		if r.OperationName != "" {
			methodArgs.OperationName = pulumi.String(r.OperationName)
		}
		if len(r.RequestParameters) > 0 {
			methodArgs.RequestParameters = pulumi.ToBoolMap(r.RequestParameters)
		}
		if len(r.RequestModels) > 0 {
			methodArgs.RequestModels = modelNameMap(r.RequestModels, satellites)
		}
		if r.RequestValidatorName != "" {
			methodArgs.RequestValidatorId = satellites.validators[r.RequestValidatorName].ID()
		}
		method, err := apigateway.NewMethod(ctx, "method-"+key, methodArgs, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create method %q", key)
		}
		tree.definitionResources = append(tree.definitionResources, method)

		integration, err := routeIntegration(ctx, locals, provider, api, resourceId, r, method)
		if err != nil {
			return nil, errors.Wrapf(err, "create integration %q", key)
		}
		tree.definitionResources = append(tree.definitionResources, integration)

		// Typed responses: the method response and its integration
		// response mapping travel as one spec entry.
		for _, resp := range r.Responses {
			methodResponseArgs := &apigateway.MethodResponseArgs{
				RestApi:    api.ID(),
				ResourceId: resourceId,
				HttpMethod: method.HttpMethod,
				StatusCode: pulumi.String(resp.StatusCode),
			}
			if len(resp.ResponseModels) > 0 {
				methodResponseArgs.ResponseModels = modelNameMap(resp.ResponseModels, satellites)
			}
			if len(resp.ResponseParameters) > 0 {
				methodResponseArgs.ResponseParameters = pulumi.ToBoolMap(resp.ResponseParameters)
			}
			methodResponse, err := apigateway.NewMethodResponse(ctx, "method-response-"+key+"-"+resp.StatusCode, methodResponseArgs, pulumi.Provider(provider))
			if err != nil {
				return nil, errors.Wrapf(err, "create method response %q %s", key, resp.StatusCode)
			}
			tree.definitionResources = append(tree.definitionResources, methodResponse)

			integrationResponseArgs := &apigateway.IntegrationResponseArgs{
				RestApi:    api.ID(),
				ResourceId: resourceId,
				HttpMethod: method.HttpMethod,
				StatusCode: methodResponse.StatusCode,
			}
			if resp.SelectionPattern != "" {
				integrationResponseArgs.SelectionPattern = pulumi.String(resp.SelectionPattern)
			}
			if len(resp.IntegrationResponseParameters) > 0 {
				integrationResponseArgs.ResponseParameters = pulumi.ToStringMap(resp.IntegrationResponseParameters)
			}
			if len(resp.IntegrationResponseTemplates) > 0 {
				integrationResponseArgs.ResponseTemplates = pulumi.ToStringMap(resp.IntegrationResponseTemplates)
			}
			if resp.ContentHandling != "" {
				integrationResponseArgs.ContentHandling = pulumi.String(resp.ContentHandling)
			}
			// AWS requires the integration AND the method response to
			// exist first (see the header comment).
			integrationResponse, err := apigateway.NewIntegrationResponse(ctx, "integration-response-"+key+"-"+resp.StatusCode,
				integrationResponseArgs, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{integration, methodResponse}))
			if err != nil {
				return nil, errors.Wrapf(err, "create integration response %q %s", key, resp.StatusCode)
			}
			tree.definitionResources = append(tree.definitionResources, integrationResponse)
		}
	}

	return tree, nil
}

// routeIntegration renders one route's backend integration.
func routeIntegration(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, api *apigateway.RestApi,
	resourceId pulumi.StringInput, r *awsrestapigatewayv1alpha1.AwsRestApiGatewayRoute, method *apigateway.Method) (*apigateway.Integration, error) {
	i := r.Integration

	args := &apigateway.IntegrationArgs{
		RestApi:    api.ID(),
		ResourceId: resourceId,
		HttpMethod: method.HttpMethod,
		Type:       pulumi.String(i.Type),
	}
	if i.Uri.GetValue() != "" {
		args.Uri = pulumi.String(i.Uri.GetValue())
	}
	// Lambda invocations are always POST; the spec's validation allows
	// omitting it on AWS_PROXY and the module fills it.
	if i.HttpMethod != "" {
		args.IntegrationHttpMethod = pulumi.String(i.HttpMethod)
	} else if i.Type == "AWS_PROXY" {
		args.IntegrationHttpMethod = pulumi.String("POST")
	}
	if i.CredentialsArn.GetValue() != "" {
		args.Credentials = pulumi.String(i.CredentialsArn.GetValue())
	}
	if i.ConnectionType != "" {
		args.ConnectionType = pulumi.String(i.ConnectionType)
	}
	if i.VpcLinkId.GetValue() != "" {
		args.ConnectionId = pulumi.String(i.VpcLinkId.GetValue())
	}
	if i.PassthroughBehavior != "" {
		args.PassthroughBehavior = pulumi.String(i.PassthroughBehavior)
	}
	if i.ContentHandling != "" {
		args.ContentHandling = pulumi.String(i.ContentHandling)
	}
	if len(i.CacheKeyParameters) > 0 {
		args.CacheKeyParameters = pulumi.ToStringArray(i.CacheKeyParameters)
	}
	if i.CacheNamespace != "" {
		args.CacheNamespace = pulumi.String(i.CacheNamespace)
	}
	if len(i.RequestParameters) > 0 {
		args.RequestParameters = pulumi.ToStringMap(i.RequestParameters)
	}
	if len(i.RequestTemplates) > 0 {
		args.RequestTemplates = pulumi.ToStringMap(i.RequestTemplates)
	}
	if i.TimeoutMilliseconds > 0 {
		args.TimeoutMilliseconds = pulumi.Int(int(i.TimeoutMilliseconds))
	}
	if i.ResponseTransferMode != "" {
		args.ResponseTransferMode = pulumi.String(i.ResponseTransferMode)
	}
	if i.TlsInsecureSkipVerification {
		args.TlsConfig = &apigateway.IntegrationTlsConfigArgs{
			InsecureSkipVerification: pulumi.Bool(true),
		}
	}

	return apigateway.NewIntegration(ctx, "integration-"+routeKey(r), args, pulumi.Provider(provider))
}

// modelNameMap resolves model references: names defined in this spec
// resolve through the created model (carrying the dependency edge);
// the AWS built-ins ("Empty"/"Error") pass through as literals.
func modelNameMap(refs map[string]string, satellites *createdSatellites) pulumi.StringMap {
	out := pulumi.StringMap{}
	for contentType, name := range refs {
		if model, ok := satellites.models[name]; ok {
			out[contentType] = model.Name
			continue
		}
		out[contentType] = pulumi.String(name).ToStringOutput()
	}
	return out
}

func svrsToStringArray(in []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	out := pulumi.StringArray{}
	for _, ref := range in {
		out = append(out, pulumi.String(ref.GetValue()))
	}
	return out
}
