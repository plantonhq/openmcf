package resources

// The front-door route table: the ONE statement of how a request that reaches
// a Planton platform at its public origin is routed to a component. Every
// front door renders from it -- the Ingress object, the Gateway API HTTPRoute,
// and the built-in nginx gateway that serves over kubectl port-forward -- so
// the three doors are the same architecture at different addresses, and a
// route added here appears on all of them in one change.
//
// Every rule is a plain path prefix with the segment-wise meaning Kubernetes
// defines as portable (Ingress pathType Prefix, Gateway API PathPrefix). That
// is why the browser API has a path namespace of its own (APIPathPrefix): a
// raw gRPC request path is a single segment, /ai.planton.<Service>/<Method>,
// which no portable prefix rule can tell from a console page -- it took a
// controller-specific regex to route, and a controller that does not know the
// regex sends every API call to the console. The control plane serves its
// gRPC-Web door under the same prefix, and the console is handed a base URL
// that ends with it, so no client composes the shape.
//
// Order is most-specific first. Renderers that match by longest prefix (the
// Ingress, the Gateway API) do not depend on it; nginx does not either
// (longest-prefix location wins), but the rendered config reads in the same
// order as this table.

const (
	// APIPathPrefix is the path namespace of the browser-facing gRPC-Web API.
	// Mirrors the control plane's GrpcWebServerConfig.API_PATH_PREFIX -- the
	// two cannot import each other, and the boot-contract floor is where a
	// change to one is reconciled with the other.
	APIPathPrefix = "/rpc"

	// StoragePathPrefix routes the storage relay: expiring transfer URLs
	// (state-file download/upload) served by the control plane on its
	// browser-API port. The path shape is the control plane's
	// BlobRelayService.RELAY_PATH_PREFIX contract ("/storage/v1/relay/...");
	// routing the "/storage" root keeps room for future storage surfaces
	// without another edge change.
	StoragePathPrefix = "/storage"

	// ConsolePathPrefix is the catch-all: everything no other rule claims is
	// a web console page.
	ConsolePathPrefix = "/"
)

// FrontDoorBackend names the component a route delivers to.
type FrontDoorBackend int

const (
	// BackendControlPlane is the control plane's browser-API port (gRPC-Web
	// plus the storage relay).
	BackendControlPlane FrontDoorBackend = iota
	// BackendIdentity is the identity server (sign-in pages, OIDC endpoints).
	BackendIdentity
	// BackendConsole is the web console.
	BackendConsole
)

// FrontDoorRoute is one rule of the table.
type FrontDoorRoute struct {
	// PathPrefix is matched segment-wise against the request path.
	PathPrefix string
	Backend    FrontDoorBackend
}

// FrontDoorRoutes returns the table, most-specific rule first.
func FrontDoorRoutes() []FrontDoorRoute {
	return []FrontDoorRoute{
		{PathPrefix: APIPathPrefix, Backend: BackendControlPlane},
		{PathPrefix: StoragePathPrefix, Backend: BackendControlPlane},
		{PathPrefix: IdentityPathPrefix, Backend: BackendIdentity},
		{PathPrefix: ConsolePathPrefix, Backend: BackendConsole},
	}
}

// ServiceName is the Kubernetes Service the route's backend is reached at.
func (r FrontDoorRoute) ServiceName(crName string) string {
	switch r.Backend {
	case BackendIdentity:
		return IdentityServiceName(crName)
	case BackendConsole:
		return ConsoleServiceName(crName)
	default:
		return ControlPlaneServiceName(crName)
	}
}

// ServicePortName is the named Service port the route targets; Ingress and
// HTTPRoute backends reference ports by name so a port number change never
// touches the doors.
func (r FrontDoorRoute) ServicePortName() string {
	if r.Backend == BackendControlPlane {
		return "grpc-web"
	}
	return "http"
}

// ServicePort is the numeric Service port, for the one door (nginx) that
// dials upstreams by address.
func (r FrontDoorRoute) ServicePort() int {
	switch r.Backend {
	case BackendIdentity:
		return identityServicePort
	case BackendConsole:
		return consoleServicePort
	default:
		return controlPlaneGrpcWebPort
	}
}

// APIURL renders the base URL a browser client is given for the API: the
// front-door origin plus the API path namespace. Clients append
// "/<service>/<method>" to it.
func APIURL(frontDoorURL string) string {
	return frontDoorURL + APIPathPrefix
}
