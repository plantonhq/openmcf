package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// routes creates API Gateway routes, wiring each route to its deduplicated
// integration and optional authorizer. Returns the created routes so the
// stage (whose per-route settings reference route keys) can depend on them.
func routes(
	ctx *pulumi.Context,
	locals *Locals,
	createdApi *apigatewayv2.Api,
	createdIntegrations []*apigatewayv2.Integration,
	authorizerMap map[string]*apigatewayv2.Authorizer,
	provider *aws.Provider,
) ([]pulumi.Resource, error) {
	_, routeToIndex := dedupIntegrations(locals.Spec.Routes)
	created := make([]pulumi.Resource, 0, len(locals.Spec.Routes))

	for i, route := range locals.Spec.Routes {
		// Build a safe resource name from the route key.
		// "GET /users" -> "get-users", "$default" -> "default"
		safeName := sanitizeRouteKey(route.RouteKey)
		resourceName := fmt.Sprintf("%s-route-%s", locals.ApiName, safeName)

		integration := createdIntegrations[routeToIndex[i]]

		args := &apigatewayv2.RouteArgs{
			ApiId:    createdApi.ID(),
			RouteKey: pulumi.String(route.RouteKey),
			// Target format: "integrations/{integrationId}"
			Target: integration.ID().ApplyT(func(id string) string {
				return fmt.Sprintf("integrations/%s", id)
			}).(pulumi.StringOutput),
		}

		// NONE is the AWS default; only real authorization modes are sent.
		if route.AuthorizationType != "" && route.AuthorizationType != "NONE" {
			args.AuthorizationType = pulumi.StringPtr(route.AuthorizationType)

			// JWT routes bind JWT authorizers; CUSTOM routes bind REQUEST
			// (Lambda) authorizers -- the spec CEL guarantees the referenced
			// authorizer's type matches.
			if route.AuthorizerName != "" {
				if authorizer, ok := authorizerMap[route.AuthorizerName]; ok {
					args.AuthorizerId = authorizer.ID().ToStringOutput()
				}
			}

			// Scopes only apply to JWT authorization.
			if route.AuthorizationType == "JWT" && len(route.AuthorizationScopes) > 0 {
				args.AuthorizationScopes = pulumi.ToStringArray(route.AuthorizationScopes)
			}
		}

		// Stable operationId for OpenAPI exports and generated clients.
		if route.OperationName != "" {
			args.OperationName = pulumi.StringPtr(route.OperationName)
		}

		createdRoute, err := apigatewayv2.NewRoute(ctx, resourceName, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create route %q (index %d)", route.RouteKey, i)
		}
		created = append(created, createdRoute)
	}

	return created, nil
}

// sanitizeRouteKey converts a route key into a safe resource name component.
// "GET /users/{id}" -> "get-users-id"
// "$default" -> "default"
func sanitizeRouteKey(routeKey string) string {
	s := strings.ToLower(routeKey)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "$", "")
	// Collapse multiple hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}
