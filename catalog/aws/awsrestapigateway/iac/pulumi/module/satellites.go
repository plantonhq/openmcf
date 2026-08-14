package module

import (
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createdSatellites carries what routes and the stage need from the
// named satellites: created resources keyed by their spec names.
type createdSatellites struct {
	authorizers map[string]*apigateway.Authorizer
	validators  map[string]*apigateway.RequestValidator
	models      map[string]*apigateway.Model
	// documentationParts orders the documentation version behind the
	// parts it publishes.
	documentationParts []pulumi.Resource
	// gatewayResponses orders the deployment behind the response
	// customizations - they take effect only in a deployed snapshot.
	gatewayResponses []pulumi.Resource
	// documentationVersion orders the stage behind the published
	// version it references (AWS rejects a stage naming a version that
	// does not exist yet).
	documentationVersion pulumi.Resource
}

// apiSatellites creates the API-scoped named satellites: models,
// request validators, authorizers, gateway responses, and
// documentation parts.
//
// Lifecycle facts the renders below depend on:
//   - a model's content type is immutable (the provider replaces the
//     model); routes reference models by NAME, which survives the
//     replacement;
//   - deleting an authorizer or validator a method still references is
//     swallowed upstream (ConflictException) -- the object dangles
//     until the API deletes; harmless for this module because methods
//     and satellites live in one plan;
//   - Cognito authorizers get their credentials via a post-create
//     PATCH upstream (the create API silently ignores them);
//   - a documentation part addresses its element by LOCATION, which is
//     immutable -- location edits replace the part.
func apiSatellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, api *apigateway.RestApi) (*createdSatellites, error) {
	spec := locals.Spec
	created := &createdSatellites{
		authorizers: map[string]*apigateway.Authorizer{},
		validators:  map[string]*apigateway.RequestValidator{},
		models:      map[string]*apigateway.Model{},
	}

	// JSON Schema payload models. Iteration is identity-sorted for
	// deterministic previews (names here, response types for gateway
	// responses; documentation parts keep declaration order - position
	// IS their identity).
	modelIds := pulumi.StringMap{}
	sortedModels := append([]int{}, indexRange(len(spec.Models))...)
	sort.Slice(sortedModels, func(i, j int) bool { return spec.Models[sortedModels[i]].Name < spec.Models[sortedModels[j]].Name })
	for _, i := range sortedModels {
		m := spec.Models[i]
		args := &apigateway.ModelArgs{
			RestApi:     api.ID(),
			Name:        pulumi.String(m.Name),
			ContentType: pulumi.String(m.ContentType),
		}
		if m.Description != "" {
			args.Description = pulumi.String(m.Description)
		}
		if m.Schema != "" {
			args.Schema = pulumi.String(m.Schema)
		}
		model, err := apigateway.NewModel(ctx, "model-"+m.Name, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create model %q", m.Name)
		}
		created.models[m.Name] = model
		modelIds[m.Name] = model.ID().ToStringOutput()
	}
	ctx.Export(OpModelIds, modelIds)

	// Request validators.
	validatorIds := pulumi.StringMap{}
	sortedValidators := append([]int{}, indexRange(len(spec.RequestValidators))...)
	sort.Slice(sortedValidators, func(i, j int) bool {
		return spec.RequestValidators[sortedValidators[i]].Name < spec.RequestValidators[sortedValidators[j]].Name
	})
	for _, i := range sortedValidators {
		v := spec.RequestValidators[i]
		validator, err := apigateway.NewRequestValidator(ctx, "validator-"+v.Name, &apigateway.RequestValidatorArgs{
			RestApi:                   api.ID(),
			Name:                      pulumi.String(v.Name),
			ValidateRequestBody:       pulumi.Bool(v.ValidateRequestBody),
			ValidateRequestParameters: pulumi.Bool(v.ValidateRequestParameters),
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create request validator %q", v.Name)
		}
		created.validators[v.Name] = validator
		validatorIds[v.Name] = validator.ID().ToStringOutput()
	}
	ctx.Export(OpRequestValidatorIds, validatorIds)

	// Authorizers: Lambda TOKEN/REQUEST or Cognito user pools.
	authorizerIds := pulumi.StringMap{}
	sortedAuthorizers := append([]int{}, indexRange(len(spec.Authorizers))...)
	sort.Slice(sortedAuthorizers, func(i, j int) bool {
		return spec.Authorizers[sortedAuthorizers[i]].Name < spec.Authorizers[sortedAuthorizers[j]].Name
	})
	for _, i := range sortedAuthorizers {
		a := spec.Authorizers[i]
		args := &apigateway.AuthorizerArgs{
			RestApi: api.ID(),
			Name:    pulumi.String(a.Name),
			Type:    pulumi.String(a.Type),
		}
		if a.LambdaInvokeUri.GetValue() != "" {
			args.AuthorizerUri = pulumi.String(a.LambdaInvokeUri.GetValue())
		}
		if a.CredentialsArn.GetValue() != "" {
			args.AuthorizerCredentials = pulumi.String(a.CredentialsArn.GetValue())
		}
		if len(a.ProviderArns) > 0 {
			args.ProviderArns = svrsToStringArray(a.ProviderArns)
		}
		if a.IdentitySource != "" {
			args.IdentitySource = pulumi.String(a.IdentitySource)
		}
		if a.IdentityValidationExpression != "" {
			args.IdentityValidationExpression = pulumi.String(a.IdentityValidationExpression)
		}
		// Presence-typed: an explicit 0 disables caching; unset keeps
		// AWS's 300-second default.
		if a.ResultTtlSeconds != nil {
			args.AuthorizerResultTtlInSeconds = pulumi.Int(int(*a.ResultTtlSeconds))
		}
		authorizer, err := apigateway.NewAuthorizer(ctx, "authorizer-"+a.Name, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create authorizer %q", a.Name)
		}
		created.authorizers[a.Name] = authorizer
		authorizerIds[a.Name] = authorizer.ID().ToStringOutput()
	}
	ctx.Export(OpAuthorizerIds, authorizerIds)

	// Gateway-generated response customizations, keyed by type.
	sortedResponses := append([]int{}, indexRange(len(spec.GatewayResponses))...)
	sort.Slice(sortedResponses, func(i, j int) bool {
		return spec.GatewayResponses[sortedResponses[i]].ResponseType < spec.GatewayResponses[sortedResponses[j]].ResponseType
	})
	for _, i := range sortedResponses {
		g := spec.GatewayResponses[i]
		args := &apigateway.ResponseArgs{
			RestApiId:    api.ID(),
			ResponseType: pulumi.String(g.ResponseType),
		}
		if g.StatusCode != "" {
			args.StatusCode = pulumi.String(g.StatusCode)
		}
		if len(g.ResponseParameters) > 0 {
			args.ResponseParameters = pulumi.ToStringMap(g.ResponseParameters)
		}
		if len(g.ResponseTemplates) > 0 {
			args.ResponseTemplates = pulumi.ToStringMap(g.ResponseTemplates)
		}
		response, err := apigateway.NewResponse(ctx, "gateway-response-"+g.ResponseType, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create gateway response %q", g.ResponseType)
		}
		created.gatewayResponses = append(created.gatewayResponses, response)
	}

	// Documentation parts, keyed by declaration position (locations are
	// composite; position is the stable cross-engine key). The part-ID
	// map exports unconditionally (empty without documentation) so both
	// engines emit the same output set.
	documentationPartIds := pulumi.StringMap{}
	if spec.Documentation != nil {
		for i, p := range spec.Documentation.Parts {
			key := strconv.Itoa(i)
			location := &apigateway.DocumentationPartLocationArgs{
				Type: pulumi.String(p.Location.Type),
			}
			if p.Location.Path != "" {
				location.Path = pulumi.String(p.Location.Path)
			}
			if p.Location.Method != "" {
				location.Method = pulumi.String(p.Location.Method)
			}
			if p.Location.Name != "" {
				location.Name = pulumi.String(p.Location.Name)
			}
			if p.Location.StatusCode != "" {
				location.StatusCode = pulumi.String(p.Location.StatusCode)
			}
			part, err := apigateway.NewDocumentationPart(ctx, "documentation-part-"+key, &apigateway.DocumentationPartArgs{
				RestApiId:  api.ID(),
				Location:   location,
				Properties: pulumi.String(p.Properties),
			}, pulumi.Provider(provider))
			if err != nil {
				return nil, errors.Wrapf(err, "create documentation part %s", key)
			}
			created.documentationParts = append(created.documentationParts, part)
			documentationPartIds[key] = part.ID().ToStringOutput()
		}

		// Publishing snapshots the parts - it must run after them.
		if spec.Documentation.PublishedVersion != nil {
			versionArgs := &apigateway.DocumentationVersionArgs{
				RestApiId: api.ID(),
				Version:   pulumi.String(spec.Documentation.PublishedVersion.Version),
			}
			if spec.Documentation.PublishedVersion.Description != "" {
				versionArgs.Description = pulumi.String(spec.Documentation.PublishedVersion.Description)
			}
			version, err := apigateway.NewDocumentationVersion(ctx, "documentation-version", versionArgs,
				pulumi.Provider(provider), pulumi.DependsOn(created.documentationParts))
			if err != nil {
				return nil, errors.Wrap(err, "create documentation version")
			}
			created.documentationVersion = version
		}
	}
	ctx.Export(OpDocumentationPartIds, documentationPartIds)

	return created, nil
}

func indexRange(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}
